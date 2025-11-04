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

func TestNormalizeXMLPluginVersions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single plugin",
			input:    `<flow-definition plugin="workflow-job@2.25">`,
			expected: `<flow-definition plugin="workflow-job">`,
		},
		{
			name:     "multiple plugins",
			input:    `<flow-definition plugin="workflow-job@2.25"><definition plugin="workflow-cps@2.59">`,
			expected: `<flow-definition plugin="workflow-job"><definition plugin="workflow-cps">`,
		},
		{
			name:     "no plugins",
			input:    `<project><description>test</description></project>`,
			expected: `<project><description>test</description></project>`,
		},
		{
			name:     "mixed content",
			input:    `<scm class="hudson.plugins.git.GitSCM" plugin="git@3.9.1"><url>test</url></scm>`,
			expected: `<scm class="hudson.plugins.git.GitSCM" plugin="git"><url>test</url></scm>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := NormalizeXMLPluginVersions(tt.input)
			if actual != tt.expected {
				t.Errorf("Expected %s but got: %s", tt.expected, actual)
			}
		})
	}
}

func TestTemplateDiff_SkipPluginVersions(t *testing.T) {
	// Set up Job
	job := resourceJenkinsJob()

	// Test case 1: skip_plugins_version_update = true, should ignore version differences
	bagWithSkip := job.TestResourceData()
	_ = bagWithSkip.Set("skip_plugins_version_update", true)

	inputOld := `<flow-definition plugin="workflow-job@2.25"><definition plugin="workflow-cps@2.59"></definition></flow-definition>`
	inputNew := `<flow-definition plugin="workflow-job@3.0"><definition plugin="workflow-cps@3.0"></definition></flow-definition>`

	if actual := templateDiff("", inputOld, inputNew, bagWithSkip); !actual {
		t.Errorf("Expected XMLs with different plugin versions to be considered equal when skip_plugins_version_update=true")
	}

	// Test case 2: skip_plugins_version_update = false, should detect version differences
	bagWithoutSkip := job.TestResourceData()
	_ = bagWithoutSkip.Set("skip_plugins_version_update", false)

	if actual := templateDiff("", inputOld, inputNew, bagWithoutSkip); actual {
		t.Errorf("Expected XMLs with different plugin versions to be considered inequal when skip_plugins_version_update=false")
	}

	// Test case 3: skip_plugins_version_update = true but old is empty (create), should not skip
	inputOldEmpty := ""
	inputNewCreate := `<flow-definition plugin="workflow-job@2.25"></flow-definition>`

	if actual := templateDiff("", inputOldEmpty, inputNewCreate, bagWithSkip); actual {
		t.Errorf("Expected comparison with empty old value to work normally even with skip_plugins_version_update=true")
	}

	// Test case 4: skip_plugins_version_update = true, but actual content differs
	inputOldDifferent := `<flow-definition plugin="workflow-job@2.25"><description>old</description></flow-definition>`
	inputNewDifferent := `<flow-definition plugin="workflow-job@3.0"><description>new</description></flow-definition>`

	if actual := templateDiff("", inputOldDifferent, inputNewDifferent, bagWithSkip); actual {
		t.Errorf("Expected XMLs with different content to be considered inequal even when skip_plugins_version_update=true")
	}
}
