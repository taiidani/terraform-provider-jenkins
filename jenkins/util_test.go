package jenkins

import (
	"testing"
)

func TestFormatFolderName(t *testing.T) {
	inputSimple, inputFolder, inputNested, inputDuped := "job-name", "folder/job-name", "parent/child/job-name", "parent/job/child/job/job-name"

	// Simple
	actual := formatFolderName(inputSimple)
	if actual != inputSimple {
		t.Errorf("Expected %s but received %s", inputSimple, actual)
	}

	// Folder
	actual = formatFolderName(inputFolder)
	if actual != "folder/job/job-name" {
		t.Errorf("Expected %s but received %s", inputFolder, actual)
	}

	// Nested
	actual = formatFolderName(inputNested)
	if actual != "parent/job/child/job/job-name" {
		t.Errorf("Expected %s but received %s", inputNested, actual)
	}

	// Deduplicate
	actual = formatFolderName(inputDuped)
	if actual != "parent/job/child/job/job-name" {
		t.Errorf("Expected %s but received %s", inputDuped, actual)
	}
}

func TestFormatFolderID(t *testing.T) {
	inputSimple := []string{"folder-id"}
	inputNested := []string{"folder-parent", "folder-id"}
	inputDuped := []string{"folder-parent", "job", "folder-id"}

	// Simple
	actual := formatFolderID(inputSimple)
	if actual != "/job/folder-id" {
		t.Errorf("Expected /job/folder-id but received %s", actual)
	}

	// Nested
	actual = formatFolderID(inputNested)
	if actual != "/job/folder-parent/job/folder-id" {
		t.Errorf("Expected /job/folder-parent/job/folder-id but received %s", actual)
	}

	// Deduplicate
	actual = formatFolderID(inputDuped)
	if actual != "/job/folder-parent/job/folder-id" {
		t.Errorf("Expected /job/folder-parent/job/folder-id but received %s", actual)
	}
}

func TestParseCanonicalJobID(t *testing.T) {
	inputSimple, inputFolder, inputNested := "job-name", "folder/job-name", "parent/child/job-name"

	// Simple
	actual, actualFolders := parseCanonicalJobID(inputSimple)
	if actual != inputSimple || len(actualFolders) != 0 {
		t.Errorf("Expected %s with empty folder array but received %s %s", inputSimple, actual, actualFolders)
	}

	// Folder
	actual, actualFolders = parseCanonicalJobID(inputFolder)
	if actual != inputSimple || len(actualFolders) != 1 || actualFolders[0] != "folder" {
		t.Errorf("Expected %s with single folder array but received %s %s", inputSimple, actual, actualFolders)
	}

	// Nested
	actual, actualFolders = parseCanonicalJobID(inputNested)
	if actual != inputSimple || len(actualFolders) != 2 || actualFolders[0] != "parent" || actualFolders[1] != "child" {
		t.Errorf("Expected %s with double folder array but received %s %s", inputSimple, actual, actualFolders)
	}
}

func TestGenerateCredentialID(t *testing.T) {
	inputFolder, inputName := "test-folder", "test-name"
	actual := generateCredentialID(inputFolder, inputName)
	if actual != "test-folder/test-name" {
		t.Errorf("Expected %s/%s but got: %s", inputFolder, inputName, actual)
	}
}
