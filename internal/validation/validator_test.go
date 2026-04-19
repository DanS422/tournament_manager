package validation_test

import (
	"reflect"
	"testing"

	"tournament_manager/internal/validation"
)

type testForm struct {
	Name     string `json:"name" validate:"required"`
	Location string `json:"location" validate:"min=2"`
}

func TestValidateStruct_Valid(t *testing.T) {
	form := testForm{
		Name:     "Test Cup",
		Location: "SG",
	}

	errs := validation.ValidateStruct(form)

	if errs != nil {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestValidateStruct_Required(t *testing.T) {
	form := testForm{
		Location: "SG",
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
	}

	errs := validation.ValidateStruct(form)
	want := map[string]string{
		"location": "is invalid",
	}

	if !reflect.DeepEqual(errs, want) {
		t.Fatalf("expected %v, got %v", want, errs)
	}
}
