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

package dns

import (
	. "github.com/adjust/gocheck"
	"testing"
)

func TestDatamodel(t *testing.T) { TestingT(t) }

type DNSSuite struct{}

var _ = Suite(&DNSSuite{})

func (s *DNSSuite) TestIPRegexp(c *C) {
	c.Assert(IPRegexp.MatchString("ec2-75-101-190-248.compute-1.amazonaws.com"), Equals, false)
	c.Assert(IPRegexp.MatchString("1.1.1.1"), Equals, true)
	c.Assert(IPRegexp.MatchString("a.1.1.1"), Equals, false)
	c.Assert(IPRegexp.MatchString("1.a.1.1"), Equals, false)
	c.Assert(IPRegexp.MatchString("1.1.a.1"), Equals, false)
	c.Assert(IPRegexp.MatchString("1.1.1.a"), Equals, false)
	c.Assert(IPRegexp.MatchString("12.12.12.12"), Equals, true)
	c.Assert(IPRegexp.MatchString("123.123.123.123"), Equals, true)
	c.Assert(IPRegexp.MatchString("q123.123.123.123"), Equals, false)
	c.Assert(IPRegexp.MatchString("123.123.123.123q"), Equals, false)
}
