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

func (s *DatamodelSuite) TestEnvironment(c *C) {
	path := helper.GetBaseEnvPath()
	Zk.RecursiveDelete(path)
	// Create root environment
	root := Env("root")
	c.Assert(root.Get(), Not(IsNil))
	_, err := GetEnv("root")
	c.Assert(err, Not(IsNil))
	c.Assert(root.Save(), IsNil)

	// test get root env
	gRoot, err := GetEnv("root")
	c.Assert(err, IsNil)
	c.Assert(gRoot, DeepEquals, root)
	gRoot = Env("root")
	c.Assert(gRoot.Get(), IsNil)
	c.Assert(gRoot, DeepEquals, root)
}
