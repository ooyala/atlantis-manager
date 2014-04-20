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
	"atlantis/builder/api/types"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type SimpleBuilder struct {
	URL          string
	BuildTimeout time.Duration
}

func NewSimpleBuilder(url string) *SimpleBuilder {
	return &SimpleBuilder{URL: url}
}

func decodeBuildResp(resp *http.Response) (*types.Build, error) {
	if resp.StatusCode != 200 {
		// non-200 status is an error
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}
		return nil, errors.New("Bad Status: "+string(bodyBytes))
	}
	var build types.Build
	if err := json.NewDecoder(resp.Body).Decode(&build); err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()
	return &build, nil
}

func (b *SimpleBuilder) Build(t *Task, repo, root, sha string) (io.ReadCloser, error) {
	t.LogStatus("Triggering Simple Build")
	jsonBytes, err := json.Marshal(&types.Build{URL: repo, RelPath: root, Sha: sha})
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(b.URL+"/build", "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, errors.New("Builder Error: "+err.Error())
	}
	build, err := decodeBuildResp(resp)
	if err != nil {
		return nil, err
	}
	begin := time.Now()
	end := begin.Add(b.BuildTimeout)
	for build.Status != types.StatusDone && build.Status != types.StatusError {
		if time.Now().After(end) {
			// timeout
			return nil, errors.New("Build Timeout.")
		}
		resp, err := http.Get(b.URL+"/build/"+build.ID)
		if err != nil {
			return nil, errors.New("Builder Error: "+err.Error())
		}
		if build, err = decodeBuildResp(resp); err != nil {
			return nil, err
		}
	}
	if build.Status == types.StatusError {
		return nil, errors.New(fmt.Sprintf("%v", build.Error))
	}

	// must be status done, fetch manifest
	resp, err = http.Get(b.URL+"/build/"+build.ID+"/manifest")
	if err != nil {
		return nil, errors.New("Builder Error: "+err.Error())
	}
	return resp.Body, nil
}
