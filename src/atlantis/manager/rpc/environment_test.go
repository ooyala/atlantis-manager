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

package rpc

import (
	. "github.com/adjust/gocheck"
)

type EnvironmentSuite struct{}

var _ = Suite(&EnvironmentSuite{})

func (s *EnvironmentSuite) TestEnvironmentTracker(c *C) {
	c.Assert(IsEnvInUse("root"), Equals, false)

	AddAppShaToEnv("app", "sha", "root")
	c.Assert(IsEnvInUse("root"), Equals, true)
	AddAppShaToEnv("app", "sha", "root")
	c.Assert(IsEnvInUse("root"), Equals, true)
	DeleteAppShaFromEnv("app", "sha", "root")
	c.Assert(IsEnvInUse("root"), Equals, true)
	DeleteAppShaFromEnv("app", "sha", "root")
	c.Assert(IsEnvInUse("root"), Equals, false)
	DeleteEnv("root")
	c.Assert(IsEnvInUse("root"), Equals, false)
}
