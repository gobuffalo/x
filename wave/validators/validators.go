package validators

import (
	"fmt"
	"mime/multipart"
	"strconv"

	humanize "github.com/dustin/go-humanize"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
)

type FileTypeValidator struct {
	Field        string
	AllowedTypes map[string]bool
	Headers      *multipart.FileHeader
}

func (f FileTypeValidator) IsValid(errors *validate.Errors) {
	key := validators.GenerateKey(f.Field)
	if !f.AllowedTypes[f.Headers.Header.Get("Content-Type")] {
		errors.Add(key, "not an allowed type")
	}
}

type MaxFileSizeValidator struct {
	Field   string
	MaxSize int
	Headers *multipart.FileHeader
}

func (m MaxFileSizeValidator) IsValid(errors *validate.Errors) {
	key := validators.GenerateKey(m.Field)
	s, err := strconv.Atoi(m.Headers.Header.Get("Content-Length"))
	if err != nil {
		errors.Add(key, "couldn't parse content length")
		return
	}
	if s > m.MaxSize {
		errors.Add(key, fmt.Sprintf("is too big %s", humanize.Bytes(uint64(s))))
	}
}

const MaxImageSize = 1024 * 1024 * 5

var AllowedImages = map[string]bool{
	"image/png":  true,
	"image/jpg":  true,
	"image/jpeg": true,
	"image/svg":  true,
	"image/gif":  true,
}
