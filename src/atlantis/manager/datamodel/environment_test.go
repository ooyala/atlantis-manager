package datamodel

import (
	"atlantis/manager/helper"
	. "launchpad.net/gocheck"
)

func (s *DatamodelSuite) TestEnvironment(c *C) {
	path := helper.GetBaseEnvPath()
	Zk.RecursiveDelete(path)
	// Create root environment
	root := Env("root", "")
	c.Assert(root.Get(), Not(IsNil))
	_, err := GetEnv("root")
	c.Assert(err, Not(IsNil))
	c.Assert(root.Save(), IsNil)

	// test get root env
	gRoot, err := GetEnv("root")
	c.Assert(err, IsNil)
	c.Assert(gRoot, DeepEquals, root)
	gRoot = Env("root", "")
	c.Assert(gRoot.Get(), IsNil)
	c.Assert(gRoot, DeepEquals, root)

	// create children
	branch := Env("branch", "root")
	c.Assert(branch.Save(), IsNil)
	leaf := Env("leaf", "branch")
	c.Assert(leaf.Save(), IsNil)

	// test get child environment
	gBranch := Env("branch", "")
	c.Assert(gBranch.Get(), IsNil)
	c.Assert(gBranch, DeepEquals, branch)

	c.Assert(root.UpdateDep("jigish", "leet"), IsNil)
	c.Assert(root.UpdateDep("rofl", "copter"), IsNil)
	c.Assert(root.UpdateDep("hello", "goodbye"), IsNil)
	c.Assert(root.UpdateDep("wtf", "bbq"), IsNil)

	c.Assert(branch.UpdateDep("hello", "sup"), IsNil)
	c.Assert(branch.UpdateDep("wtf", "omg"), IsNil)
	c.Assert(branch.UpdateDep("bah", "humbug"), IsNil)

	c.Assert(leaf.UpdateDep("rofl", "mao"), IsNil)
	c.Assert(leaf.UpdateDep("wtf", "lol"), IsNil)

	deps, err := leaf.ResolveDepValues([]string{"jigish", "rofl", "hello", "wtf", "bah"})
	c.Assert(err, IsNil)
	c.Assert(deps, DeepEquals, map[string]string{"jigish": "leet", "rofl": "mao", "hello": "sup",
		"wtf": "lol", "bah": "humbug"})
	_, err = leaf.ResolveDepValues([]string{"jjjigish", "rofl", "hello", "wtf", "bah"})
	c.Assert(err, Not(IsNil))
	deps, err = leaf.ResolveAllDepValues()
	c.Assert(err, IsNil)
	c.Assert(deps, DeepEquals, map[string]string{"jigish": "leet", "rofl": "mao", "hello": "sup",
		"wtf": "lol", "bah": "humbug"})

	c.Assert(branch.DeleteDep("hello"), IsNil)
	deps, err = leaf.ResolveDepValues([]string{"jigish", "rofl", "hello", "wtf", "bah"})
	c.Assert(err, IsNil)
	c.Assert(deps, DeepEquals, map[string]string{"jigish": "leet", "rofl": "mao", "hello": "goodbye",
		"wtf": "lol", "bah": "humbug"})

	c.Assert(branch.Delete(), IsNil)
	_, err = leaf.ResolveDepValues([]string{"jigish", "rofl", "hello", "wtf", "bah"})
	c.Assert(err, Not(IsNil))
	leaf.Parent = "root"
	c.Assert(leaf.Save, Not(IsNil))
	deps, err = leaf.ResolveDepValues([]string{"jigish", "rofl", "hello", "wtf"})
	c.Assert(err, IsNil)
	c.Assert(deps, DeepEquals, map[string]string{"jigish": "leet", "rofl": "mao", "hello": "goodbye",
		"wtf": "lol"})

	c.Assert(root.Delete(), IsNil)
	c.Assert(leaf.Delete(), IsNil)
}
