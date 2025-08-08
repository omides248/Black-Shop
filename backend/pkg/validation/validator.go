package validation

import "fmt"

type OZZOValidator interface {
	Validate() error
}

type CustomValidator struct{}

func (cv *CustomValidator) Validate(i interface{}) error {
	if v, ok := i.(OZZOValidator); ok {
		return v.Validate()
	}
	return fmt.Errorf("model does not implement OZZOValidator interface")
}

func New() *CustomValidator {
	return &CustomValidator{}
}
