package file

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/markbates/validate"
	"github.com/pkg/errors"
)

type Bucket struct {
	Dir  string
	Perm os.FileMode
}

func New(dir string, perm os.FileMode) (*Bucket, error) {
	err := os.MkdirAll(dir, perm)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &Bucket{Dir: dir, Perm: perm}, nil
}

func (c Bucket) FieldName() string {
	return "File"
}

func (c Bucket) Path(h *multipart.FileHeader) string {
	return h.Filename
}

func (c Bucket) Validate(h *multipart.FileHeader) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (c Bucket) Put(path string, r io.Reader, size int64, mt string) error {
	var err error
	path = filepath.Join(c.Dir, path)
	path, err = filepath.Abs(path)
	if err != nil {
		return errors.WithStack(err)
	}
	err = os.MkdirAll(filepath.Dir(path), c.Perm)
	if err != nil {
		return errors.WithStack(err)
	}

	f, err := os.Create(path)
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = io.Copy(f, r)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
