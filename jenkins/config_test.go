package jenkins

import jenkins "github.com/bndr/gojenkins"

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
