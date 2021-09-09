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

func TestTemplateDiff(t *testing.T) {
	// Set up inputs
	inputLeft := "<?xml version=\"1.0\" encoding=\"UTF-8\"?><root>Test Case</root>"
	inputRight := "<root>Test Case</root>"

	// Set up Job
	job := resourceJenkinsJob()
	bag := job.TestResourceData()

	if actual := templateDiff("", inputLeft, inputRight, bag); !actual {
		t.Errorf("Expected %s to be considered equal to %s", inputLeft, inputRight)
	}

	// Now try invalid inputs
	inputLeft = "<?xml version=\"1.0\" encoding=\"UTF-8\"?><root>Test Incorrect</root>"
	if actual := templateDiff("", inputLeft, inputRight, bag); actual {
		t.Errorf("Expected %s to be considered inequal to %s", inputLeft, inputRight)
	}

	inputRight = "<root>Test Incorrect</root>"
	if actual := templateDiff("", inputLeft, inputRight, bag); !actual {
		t.Errorf("Expected %s to be considered equal to %s", inputLeft, inputRight)
	}

	inputRight = "<root>Test Even More Incorrect</root>"
	if actual := templateDiff("", inputLeft, inputRight, bag); actual {
		t.Errorf("Expected %s to be considered inequal to %s", inputLeft, inputRight)
	}
}

func TestTemplateDiff_HTMLEntities(t *testing.T) {
	job := resourceJenkinsFolder()
	bag := job.TestResourceData()
	_ = bag.Set("description", "Case")

	inputLeft := "<root>&apos;/&apos;</root>"
	inputRight := "<root>'/'</root>"
	if actual := templateDiff("", inputLeft, inputRight, bag); !actual {
		t.Errorf("Expected %s to be considered equal to %s", inputLeft, inputRight)
	}

	inputLeft = "<root>'/'</root>"
	inputRight = "<root>&apos;/&apos;</root>"
	if actual := templateDiff("", inputLeft, inputRight, bag); !actual {
		t.Errorf("Expected %s to be considered equal to %s", inputLeft, inputRight)
	}
}

func TestGenerateCredentialID(t *testing.T) {
	inputFolder, inputName := "test-folder", "test-name"
	actual := generateCredentialID(inputFolder, inputName)
	if actual != "test-folder/test-name" {
		t.Errorf("Expected %s/%s but got: %s", inputFolder, inputName, actual)
	}
}
