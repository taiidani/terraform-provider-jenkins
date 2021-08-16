package jenkins

import (
	"testing"

	"github.com/hashicorp/go-cty/cty"
)

func TestValidateJobName(t *testing.T) {

	input, ctyPath := "job_name", make(cty.Path, 0)
	actual := validateJobName(input, ctyPath)

	if actual.HasError() {
		t.Errorf("Error, validation failed for input: %s", input)
	}

	// Test if we fail when we should
	input = "job_name/second_level"
	actual = validateJobName(input, ctyPath)
	if !actual.HasError() {
		t.Errorf("Error, validation failed for input: %s", input)
	}
}

func TestValidateFolderName(t *testing.T) {

	input, ctyPath := "folder_name", make(cty.Path, 0)
	actual := validateFolderName(input, ctyPath)

	if actual.HasError() {
		t.Errorf("Error, validation failed for input: %s", input)
	}
}

func TestValidateCredentialScope(t *testing.T) {

	input, ctyPath := "GLOBAL", make(cty.Path, 0)
	actual := validateCredentialScope(input, ctyPath)
	if actual.HasError() {
		t.Errorf("Error, validation failed for input: %s", input)
	}

	input = "SYSTEM"
	actual = validateCredentialScope(input, ctyPath)
	if actual.HasError() {
		t.Errorf("Error, validation failed for input: %s", input)
	}

	// Test if we fail when we should
	input = "WRONG_INPUT"
	actual = validateCredentialScope(input, ctyPath)
	if !actual.HasError() {
		t.Errorf("Error, negative validation failed for input: %s", input)
	}
}
