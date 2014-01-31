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

package datamodel

import (
	"atlantis/manager/helper"
)

type ZkDNS struct {
	App       string
	Env       string
	Shas      map[string]bool
	RecordIDs []string
}

func DNS(app, env string) *ZkDNS {
	return &ZkDNS{App: app, Env: env}
}

func GetDNS(app, env string) (zd *ZkDNS, err error) {
	zd = &ZkDNS{}
	err = getJson(helper.GetBaseDNSPath(app, env), zd)
	return
}

func (d *ZkDNS) Save() error {
	return setJson(d.path(), d)
}

func (d *ZkDNS) Delete() error {
	return Zk.RecursiveDelete(d.path())
}

func (r *ZkDNS) path() string {
	return helper.GetBaseDNSPath(r.App, r.Env)
}
