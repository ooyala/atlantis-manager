package builder

import (
	. "atlantis/common"
	"errors"
	"github.com/ooyala/go-jenkins-cli"
	"io"
)

type JenkinsBuilder struct {
	URL string
	Job string
}

func NewJenkinsBuilder(url, job string) *JenkinsBuilder {
	return &JenkinsBuilder{URL: url, Job: job}
}

func (b *JenkinsBuilder) Build(t *Task, repo, root, sha string) (io.ReadCloser, error) {
	jenkins.JENKINS_SERVER = b.URL
	t.LogStatus("Triggering Jenkins Build")
	info, err := jenkins.DoBuild(b.Job, "app_repo="+repo+"&app_root="+root+"&app_commit="+sha, true)
	if err != nil {
		return nil, errors.New("Jenkins Error: " + err.Error())
	}
	if info.Result != "SUCCESS" {
		return nil, errors.New("Jenkins Build " + info.Url + " " + info.Result)
	}
	return jenkins.GetArtifactReader(b.Job, info.Id, ManifestFile)
}
