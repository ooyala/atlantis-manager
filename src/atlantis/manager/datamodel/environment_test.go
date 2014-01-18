package datamodel

import (
	"atlantis/manager/helper"
	. "launchpad.net/gocheck"
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
