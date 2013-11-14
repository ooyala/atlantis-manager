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
	var getObj TestJsonStruct
	c.Assert(getJson("/testjson", &getObj), Not(IsNil))
	testObj := &TestJsonStruct{"string", 5, -2, true, 0.1337, []string{"leet"},
		map[string]string{"jigish": "winning"}}
	c.Assert(setJson("/testjson", testObj), IsNil)
	c.Assert(getJson("/testjson", &getObj), IsNil)
	c.Assert(getObj, DeepEquals, *testObj)
}
