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
	. "github.com/adjust/gocheck"
)

func (s *DatamodelSuite) TestDNSModel(c *C) {
	Zk.RecursiveDelete(helper.GetBaseDNSPath())
	zkDNS := DNS(app, env)
	err := zkDNS.Save()
	c.Assert(err, IsNil)
	fetchedDNS, err := GetDNS(app, env)
	c.Assert(err, IsNil)
	c.Assert(zkDNS, DeepEquals, fetchedDNS)
	zkDNS.Shas = map[string]bool{"sha1": true, "sha2": true}
	zkDNS.RecordIDs = []string{"rid1", "rid2"}
	zkDNS.Save()
	fetchedDNS, err = GetDNS(app, env)
	c.Assert(err, IsNil)
	c.Assert(zkDNS, DeepEquals, fetchedDNS)
	err = zkDNS.Delete()
	c.Assert(err, IsNil)
	_, err = GetDNS(app, env)
	c.Assert(err, Not(IsNil))
}
