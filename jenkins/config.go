package jenkins

import (
	"io/ioutil"

	jenkins "github.com/bndr/gojenkins"
)

type jenkinsClient interface {
	CreateJobInFolder(config string, jobName string, parentIDs ...string) (*jenkins.Job, error)
	Credentials() *jenkins.CredentialsManager
	DeleteJob(name string) (bool, error)
	GetJob(id string, parentIDs ...string) (*jenkins.Job, error)
}

type jenkinsAdapter struct {
	*jenkins.Jenkins
}

// Config is the set of parameters needed to configure the Jenkins provider.
type Config struct {
	ServerURL string
	CACert    string
	Username  string
	Password  string
}

func newJenkinsClient(c *Config) (*jenkinsAdapter, error) {
	client := jenkins.CreateJenkins(nil, c.ServerURL, c.Username, c.Password)
	if c.CACert != "" {
		// provide CA certificate if server is using self-signed certificate
		client.Requester.CACert, _ = ioutil.ReadFile(c.CACert)
	}
	_, err := client.Init()
	if err != nil {
		return nil, err
	}

	// return the Jenkins API client
	return &jenkinsAdapter{Jenkins: client}, nil
}

func (j *jenkinsAdapter) Credentials() *jenkins.CredentialsManager {
	return &jenkins.CredentialsManager{
		J: j.Jenkins,
	}
}
