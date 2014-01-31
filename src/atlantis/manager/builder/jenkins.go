/* Copyright 2014 Ooyala, Inc. All rights reserved.
 *
 * This file is licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
 * except in compliance with the License. You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is
 * distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

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
	return jenkins.GetArtifactReader(b.Job, info.ID, ManifestFile)
}
