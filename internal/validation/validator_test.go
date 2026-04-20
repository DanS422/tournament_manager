package validation_test

import (
	"reflect"
	"testing"

	"tournament_manager/internal/validation"
)

type testForm struct {
	Name     string `json:"name" validate:"required"`
	Location string `json:"location" validate:"min=2"`
	Gender   string `json:"gender" validate:"oneof=male female diverse"`
}

func TestValidateStruct_Valid(t *testing.T) {
	form := testForm{
		Name:     "Test Cup",
		Location: "SG",
		Gender:   "male",
	}

	errs := validation.ValidateStruct(form)

	if errs != nil {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestValidateStruct_Required(t *testing.T) {
	form := testForm{
		Location: "SG",
		Gender:   "female",
	}

	errs := validation.ValidateStruct(form)
	want := map[string]string{
		"name": "is required",
	}

	if !reflect.DeepEqual(errs, want) {
		t.Fatalf("expected %v, got %v", want, errs)
	}
}

func TestValidateStruct_Invalid(t *testing.T) {
	form := testForm{
		Name:     "Test Cup",
		Location: "S",
		Gender:   "diverse",
	}

	errs := validation.ValidateStruct(form)
	want := map[string]string{
		"location": "is invalid",
	}

	if !reflect.DeepEqual(errs, want) {
		t.Fatalf("expected %v, got %v", want, errs)
	}
}

func TestValidateStruct_OneOf(t *testing.T) {
	form := testForm{
		Name:     "Test Cup",
		Location: "SG",
		Gender:   "invalid",
	}

	errs := validation.ValidateStruct(form)
	want := map[string]string{
		"gender": "must be one of: male female diverse",
	}

	if !reflect.DeepEqual(errs, want) {
		t.Fatalf("expected %v, got %v", want, errs)
	}
}
