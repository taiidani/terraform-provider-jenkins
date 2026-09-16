package jenkins

import (
	"context"
	"io"
	"net/http"
	"strings"

	jenkins "github.com/bndr/gojenkins"
)

// htmlBodyDiscardTransport works around a behavior change in gojenkins v1.2.0:
// its response readers now error on non-JSON bodies that earlier versions
// ignored. Jenkins serves error pages (e.g. a 404) and post-redirect landing
// pages (e.g. after /doDelete) as text/html, which the client never reads.
// Emptying those bodies lets gojenkins see EOF and surface the HTTP status
// (e.g. 404) instead of a spurious "invalid character '<'" error. JSON and
// config.xml (application/xml) responses are untouched.
type htmlBodyDiscardTransport struct {
	base http.RoundTripper
}

func (t *htmlBodyDiscardTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.base.RoundTrip(req)
	if err != nil {
		return resp, err
	}

	if strings.HasPrefix(resp.Header.Get("Content-Type"), "text/html") {
		_ = resp.Body.Close()
		resp.Body = io.NopCloser(strings.NewReader(""))
		resp.ContentLength = 0
	}

	return resp, nil
}

type jenkinsClient interface {
	CreateJobInFolder(ctx context.Context, config string, jobName string, parentIDs ...string) (*jenkins.Job, error)
	Credentials() *jenkins.CredentialsManager
	DeleteJobInFolder(ctx context.Context, name string, parentIDs ...string) (bool, error)
	GetJob(ctx context.Context, id string, parentIDs ...string) (*jenkins.Job, error)
	GetFolder(ctx context.Context, id string, parents ...string) (*jenkins.Folder, error)
	GetView(ctx context.Context, name string) (*jenkins.View, error)
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
	httpClient := &http.Client{Transport: &htmlBodyDiscardTransport{base: http.DefaultTransport}}
	client := jenkins.CreateJenkins(httpClient, c.ServerURL, c.Username, c.Password)
	if c.CACert != nil {
		// provide CA certificate if server is using self-signed certificate
		if requester, ok := client.Requester.(*jenkins.Requester); ok {
			requester.CACert, _ = io.ReadAll(c.CACert)
		}
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
func (j *jenkinsAdapter) DeleteJobInFolder(ctx context.Context, name string, parentIDs ...string) (bool, error) {
	return j.DeleteJob(ctx, strings.Join(append(parentIDs, name), "/job/"))
}
