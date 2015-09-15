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
	"fmt"
	. "github.com/adjust/gocheck"
)

func (s *DatamodelSuite) TestDeployAndTeardownLocking(c *C) {
	Zk.RecursiveDelete(helper.GetBaseLockPath("deploy"))
	CreateLockPaths()

	// Fire off a bunch of deploys
	dl0 := NewDeployLock("dl0", "app0", "sha0", "env0")
	c.Assert(dl0.Lock(), IsNil)
	dl1 := NewDeployLock("dl1", "app1", "sha1", "env1")
	c.Assert(dl1.Lock(), IsNil)
	dl2 := NewDeployLock("dl2", "app1", "sha2", "env1")
	c.Assert(dl2.Lock(), IsNil)
	dl3 := NewDeployLock("dl3", "app1", "sha1", "env2")
	c.Assert(dl3.Lock(), IsNil)
	dl4 := NewDeployLock("dl4", "app1", "sha1", "env2")
	err := dl4.Lock()
	c.Assert(err, Not(IsNil))
	c.Assert(err, FitsTypeOf, LockConflictError("dl3"))
	c.Assert(err, Equals, LockConflictError("dl3"))
	// Unlock and try to lock through previously failed one
	c.Assert(dl3.Unlock(), IsNil)
	c.Assert(dl4.Lock(), IsNil)

	// Try some Teardowns
	tl0 := NewTeardownLock("tl0", "app2", "sha2", "env2")
	c.Assert(tl0.Lock(), IsNil)
	tl1 := NewTeardownLock("tl1", "app0", "sha0", "env0")
	err = tl1.Lock()
	c.Assert(err, Not(IsNil))
	c.Assert(err, FitsTypeOf, LockConflictError("dl0"))
	c.Assert(err, Equals, LockConflictError("dl0"))
	tl2 := NewTeardownLock("tl2", "app0", "sha1")
	c.Assert(tl2.Lock(), IsNil)
	tl3 := NewTeardownLock("tl3", "app3")
	c.Assert(tl3.Lock(), IsNil)
	tl4 := NewTeardownLock("tl4", "app0")
	err = tl4.Lock()
	c.Assert(err, Not(IsNil))
	c.Assert(err, FitsTypeOf, LockConflictError("dl0"))
	c.Assert(err, Equals, LockConflictError("dl0"))
	tl5 := NewTeardownLock("tl5")
	err = tl5.Lock()
	c.Assert(err, Not(IsNil))
	c.Assert(err, FitsTypeOf, LockConflictError("dl0"))

	// Try a deploy while tearing down
	dl5 := NewDeployLock("dl5", "app3", "sha3", "env3")
	err = dl5.Lock()
	c.Assert(err, Not(IsNil))
	c.Assert(err, FitsTypeOf, LockConflictError("tl3"))
	c.Assert(err, Equals, LockConflictError("tl3"))
}

func (s *DatamodelSuite) TestLockPrint(c *C) {
	e := LockConflictError("hello")
	fmt.Sprintf("%s", e)
}
