package jenkins

import (
	"io"
	"io/ioutil"
	"strings"

	jenkins "github.com/bndr/gojenkins"
)

type jenkinsClient interface {
	CreateJobInFolder(config string, jobName string, parentIDs ...string) (*jenkins.Job, error)
	Credentials() *jenkins.CredentialsManager
	DeleteJobInFolder(name string, parentIDs ...string) (bool, error)
	GetJob(id string, parentIDs ...string) (*jenkins.Job, error)
	GetFolder(id string, parents ...string) (*jenkins.Folder, error)
}

// jenkinsAdapter wraps the Jenkins client, enabling additional functionality
type jenkinsAdapter struct {
	*jenkins.Jenkins
}

// Config is the set of parameters needed to configure the Jenkins provider.
type Config struct {
	ServerURL string
	CACert    io.Reader
	Username  string
	Password  string
}

func newJenkinsClient(c *Config) *jenkinsAdapter {
	client := jenkins.CreateJenkins(nil, c.ServerURL, c.Username, c.Password)
	if c.CACert != nil {
		// provide CA certificate if server is using self-signed certificate
		client.Requester.CACert, _ = ioutil.ReadAll(c.CACert)
	}

	// return the Jenkins API client
	return &jenkinsAdapter{Jenkins: client}
}

func (j *jenkinsAdapter) Credentials() *jenkins.CredentialsManager {
	return &jenkins.CredentialsManager{
		J: j.Jenkins,
	}
}

// DeleteJobInFolder assists in running DeleteJob funcs, as DeleteJob is not folder aware
// and cannot take a canonical job ID without mishandling it.
func (j *jenkinsAdapter) DeleteJobInFolder(name string, parentIDs ...string) (bool, error) {
	return j.DeleteJob(strings.Join(append(parentIDs, name), "/job/"))
}
