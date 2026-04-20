package validation

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(v any) map[string]string {
	err := validate.Struct(v)

	if err == nil {
		return nil
	}

	return formErrors(v, err)
}

func formErrors(s any, err error) map[string]string {
	out := make(map[string]string)
	errors, ok := err.(validator.ValidationErrors)

	if !ok {
		return out
	}

	for _, e := range errors {
		field := jsonFileName(s, e.Field())

		switch e.Tag() {
		case "required":
			out[field] = "is required"
		case "oneof":
			out[field] = "must be one of: " + e.Param()
		default:
			out[field] = "is invalid"
		}
	}

	return out
}

func jsonFileName(s any, field string) string {
	t := reflect.TypeOf(s)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	f, ok := t.FieldByName(field)
	if !ok {
		return strings.ToLower(field)
	}

	tag := f.Tag.Get("json")

	if tag == "" {
		return strings.ToLower(field)
	}
	// handle json:"name,omitempty"
	name := strings.Split(tag, ",")[0]

	if name == "" {
		return strings.ToLower(field)
	}

	return name
}
