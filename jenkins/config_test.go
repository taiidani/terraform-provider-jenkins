package jenkins

import (
	"bytes"
	"testing"

	jenkins "github.com/bndr/gojenkins"
)

type mockJenkinsClient struct {
	mockCreateJobInFolder func(config string, jobName string, parentIDs ...string) (*jenkins.Job, error)
	mockDeleteJob         func(name string) (bool, error)
	mockGetJob            func(id string, parentIDs ...string) (*jenkins.Job, error)
}

func (m *mockJenkinsClient) CreateJobInFolder(config string, jobName string, parentIDs ...string) (*jenkins.Job, error) {
	return m.mockCreateJobInFolder(config, jobName, parentIDs...)
}

func (m *mockJenkinsClient) Credentials() *jenkins.CredentialsManager {
	return &jenkins.CredentialsManager{}
}

func (m *mockJenkinsClient) DeleteJob(name string) (bool, error) {
	return m.mockDeleteJob(name)
}

func (m *mockJenkinsClient) GetJob(id string, parentIDs ...string) (*jenkins.Job, error) {
	return m.mockGetJob(id, parentIDs...)
}

func TestNewJenkinsClient(t *testing.T) {
	c := newJenkinsClient(&Config{})
	if c == nil {
		t.Errorf("Expected populated client")
	}

	c = newJenkinsClient(&Config{
		CACert: bytes.NewBufferString("certificate"),
	})
	if string(c.Requester.CACert) != "certificate" {
		t.Errorf("Initialization did not extract certificate data")
	}
}

func TestJenkinsAdapter_Credentials(t *testing.T) {
	c := newJenkinsClient(&Config{})
	cm := c.Credentials()

	if cm == nil {
		t.Errorf("Expected populated client")
	} else if cm.J != c.Jenkins {
		t.Error("Expected credentials client to match client")
	}
}
