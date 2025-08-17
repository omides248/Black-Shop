package validation

import (
	"errors"
	"fmt"
	ozzo "github.com/go-ozzo/ozzo-validation/v4"
	"mime/multipart"
	"path/filepath"
	"strings"
)

func ImageRule(maxSize int64, allowedExtensions []string) ozzo.RuleFunc {
	return func(value interface{}) error {
		fileHeader, ok := value.(*multipart.FileHeader)
		if !ok || fileHeader == nil {
			return nil
		}

		if fileHeader.Size > maxSize {
			return fmt.Errorf("file size exceeds the limit of %d MB", maxSize/(1024*1024))
		}

		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		validExt := false
		for _, e := range allowedExtensions {
			if ext == e {
				validExt = true
				break
			}
		}
		if !validExt {
			return errors.New(fmt.Sprintf("invalid file extension. Only %v are allowed", allowedExtensions))
		}

		return nil
	}
}
