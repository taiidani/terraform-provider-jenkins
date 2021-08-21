package jenkins

import (
	"testing"
)

func TestRenderTemplate(t *testing.T) {
	// Set up inputs
	input := "<root>{{ .Name }}</root>"
	expected := "<root>Test Name</root>"

	// Set up Job
	job := resourceJenkinsJob()
	d := job.TestResourceData()
	_ = d.Set("name", "Test Name")
	_ = d.Set("parameters", map[string]string{"Param": "Test"})

	// Test simple
	if actual, err := renderTemplate(input, d); err != nil {
		t.Fatal(err)
	} else if actual != expected {
		t.Errorf("Expected %s to be considered equal to %s", actual, expected)
	}

	// Now with a fully populated template
	input = `<root>
	<name>{{ .Name }}</name>
	<parameters>
		{{ range $key, $value := .Parameters -}}
		<parameter>{{ $key }}: {{ $value }}</parameter>
		{{- end }}
	</parameters>
</root>
`

	expected = `<root>
	<name>Test Name</name>
	<parameters>
		<parameter>Param: Test</parameter>
	</parameters>
</root>
`

	if actual, err := renderTemplate(input, d); err != nil {
		t.Fatal(err)
	} else if actual != expected {
		t.Errorf("Expected %s to be considered equal to %s", actual, expected)
	}
}

func TestRenderTemplateInvalid(t *testing.T) {
	// Set up Job
	job := resourceJenkinsJob()
	d := job.TestResourceData()
	_ = d.Set("name", "Test Name")
	_ = d.Set("parameters", map[string]string{"Param": "Test"})

	// Now an invalid template
	input := "i am invalid{{ end }}"
	if _, err := renderTemplate(input, d); err == nil {
		t.Errorf("Expected an error to be emitted with an invalid template: %s", err)
	}

	// Now valid but with an unbound parameter
	input = "i am {{ .Mostly }} valid"
	if _, err := renderTemplate(input, d); err == nil {
		t.Errorf("Expected an error to be emitted with an invalid template: %s", err)
	}
}
