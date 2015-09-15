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
	. "github.com/adjust/gocheck"
	"sort"
)

func (s *DatamodelSuite) TestSupervisor(c *C) {
	h := Supervisor(host)
	// test delete of non-existant node
	err := h.Delete()
	c.Assert(err, Not(IsNil))
	// create a host
	c.Assert(h.Touch(), IsNil)
	// set the container and port
	inst, err := CreateInstance(app, sha, env, host)
	c.Assert(err, IsNil)
	c.Assert(h.SetContainerAndPort(inst.ID, 1337), IsNil)
	// test Info()
	data, err := h.Info()
	c.Assert(err, IsNil)
	c.Assert(data.HasAppShaEnv(app, sha, env), Equals, true)
	c.Assert(data.CountAppShaEnv(app, sha, env), Equals, 1)
	dupinst, err := CreateInstance(app, sha, env, host)
	c.Assert(err, IsNil)
	c.Assert(h.SetContainerAndPort(dupinst.ID, 1338), IsNil)
	data, err = h.Info()
	c.Assert(err, IsNil)
	c.Assert(data.CountAppShaEnv(app, sha, env), Equals, 2)
	c.Assert(data.HasAppShaEnv(app+"1", sha, env), Equals, false)
	c.Assert(data.CountAppShaEnv(app+"1", sha, env), Equals, 0)
	c.Assert(data.HasAppShaEnv(app, sha+"1", env), Equals, false)
	c.Assert(data.HasAppShaEnv(app, sha, env+"1"), Equals, false)
	// remove the container
	c.Assert(h.RemoveContainer(inst.ID), IsNil)
	inst.Delete()
	h2 := Supervisor(host + "1")
	// create a new host
	c.Assert(h2.Touch(), IsNil)
	// list the hosts to make sure they are all there
	hosts, err := ListSupervisors()
	c.Assert(err, IsNil)
	sort.Strings(hosts) // sort so DeepEquals works
	c.Assert(hosts, DeepEquals, []string{host, host + "1"})
	// test ListSupervisorsForApp
	hosts, err = ListSupervisorsForApp(app)
	c.Assert(err, IsNil)
	sort.Strings(hosts)
	c.Assert(hosts, DeepEquals, []string{host, host + "1"})
	// delete first host
	err = h.Delete()
	c.Assert(err, IsNil)
	// test to make sure host was deleted in ListSupervisors
	hosts, err = ListSupervisors()
	c.Assert(err, IsNil)
	c.Assert(hosts, DeepEquals, []string{host + "1"})
	// test to make sure host was deleted in ListSupervisorsForApp
	hosts, err = ListSupervisorsForApp(app)
	c.Assert(err, IsNil)
	c.Assert(hosts, DeepEquals, []string{host + "1"})
	// delete second host
	err = h2.Delete()
	c.Assert(err, IsNil)
	// Test to make sure there are no hosts left
	hosts, err = ListSupervisors()
	c.Assert(err, IsNil)
	c.Assert(hosts, DeepEquals, []string{})
	hosts, err = ListSupervisorsForApp(app)
	c.Assert(err, IsNil)
	c.Assert(hosts, DeepEquals, []string{})
}
