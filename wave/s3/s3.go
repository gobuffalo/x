package s3

import (
	"io"
	"mime/multipart"
	"os"

	"github.com/markbates/going/defaults"
	"github.com/markbates/validate"
	"github.com/markbates/wave"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"github.com/pkg/errors"
)

type Bucket struct {
	*s3.Bucket
}

type Options struct {
	AWSKey     string
	AWSSecret  string
	Region     aws.Region
	Permission s3.ACL
}

func New(name string, opts Options) (*Bucket, error) {
	opts.AWSKey = defaults.String(opts.AWSKey, os.Getenv("S3_KEY"))
	opts.AWSSecret = defaults.String(opts.AWSSecret, os.Getenv("S3_SECRET"))
	if (opts.Region == aws.Region{}) {
		opts.Region = aws.USEast
	}
	if opts.Permission == s3.ACL("") {
		opts.Permission = s3.PublicRead
	}

	auth, err := aws.GetAuth(opts.AWSKey, opts.AWSSecret)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	client := s3.New(auth, opts.Region)
	bucket := client.Bucket(name)
	// check if the bucket exists:
	_, err = bucket.Head("/")
	if err != nil {
		// create the bucket
		err = bucket.PutBucket(opts.Permission)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return &Bucket{bucket}, nil
}

type s3Uploader struct {
	bucket *s3.Bucket
}

func (c s3Uploader) FieldName() string {
	return "File"
}

func (c s3Uploader) Path(h *multipart.FileHeader) string {
	return h.Filename
}

func (c s3Uploader) Validate(h *multipart.FileHeader) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (c s3Uploader) Put(path string, r io.Reader, size int64, mt string) error {
	err := c.bucket.PutReader(path, r, size, mt, s3.PublicRead)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (b *Bucket) Uploader() wave.Uploader {
	return s3Uploader{
		bucket: b.Bucket,
	}
}
