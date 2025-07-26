package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	v *validator.Validate
}

func NewCustomValidator() *CustomValidator {
	v := validator.New()
	cv := &CustomValidator{v: v}

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return cv
}

func (cv *CustomValidator) Validate(i any) error {
	err := cv.v.Struct(i)
	if err != nil {
		var validationErr validator.ValidationErrors
		if errors.As(err, &validationErr) {
			return cv.newValidationError(validationErr[0].Field(), validationErr[0].Tag(), validationErr[0].Param())
		}
		return err
	}
	return nil
}

func (cv *CustomValidator) newValidationError(field string, tag string, param string) error {
	switch tag {
	case "required":
		return fmt.Errorf("field %s is required", field)
	case "len":
		return fmt.Errorf("field %s must be %s characters length", field, param)
	case "uri":
		return fmt.Errorf("field %s must be a valid URI", field)
	case "email":
		return fmt.Errorf("field %s must be a valid email address", field)
	case "min":
		return fmt.Errorf("field %s must be at least %s characters", field, param)
	case "max":
		return fmt.Errorf("field %s must be at most %s characters", field, param)
	default:
		return fmt.Errorf("field %s is invalid", field)
	}
}
