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
	. "launchpad.net/gocheck"
)

type TestJsonStruct struct {
	String  string
	UInt16  uint16
	Int     int
	Bool    bool
	Float64 float64
	Slice   []string
	Map     map[string]string
}

func (s *DatamodelSuite) TestJson(c *C) {
	Zk.RecursiveDelete("/testjson")
	var getObj TestJsonStruct
	c.Assert(getJson("/testjson", &getObj), Not(IsNil))
	testObj := &TestJsonStruct{"string", 5, -2, true, 0.1337, []string{"leet"},
		map[string]string{"jigish": "winning"}}
	c.Assert(setJson("/testjson", testObj), IsNil)
	c.Assert(getJson("/testjson", &getObj), IsNil)
	c.Assert(getObj, DeepEquals, *testObj)
}
