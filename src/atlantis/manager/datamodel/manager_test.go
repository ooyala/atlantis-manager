package datamodel

import (
	"atlantis/manager/helper"
	. "launchpad.net/gocheck"
	"sort"
)

func (s *DatamodelSuite) TestManagerPath(c *C) {
	Zk.RecursiveDelete(helper.GetBaseManagerPath())
	CreateManagerPath()
	opt := Manager("dev", "2.1.1.1")
	c.Assert(opt.path(), Equals, helper.GetBaseManagerPath("dev", "2.1.1.1"))
}

func (s *DatamodelSuite) TestManagerTouchAndDelete(c *C) {
	Zk.RecursiveDelete(helper.GetBaseManagerPath())
	CreateManagerPath()
	opt := Manager("dev", "2.1.1.1")
	err := opt.Delete()
	c.Assert(err, Not(IsNil))
	c.Assert(opt.Save(), IsNil)
	err = opt.Delete()
	c.Assert(err, IsNil)
}

func (s *DatamodelSuite) TestManagerListers(c *C) {
	Zk.RecursiveDelete(helper.GetBaseManagerPath())
	CreateManagerPath()
	regions, err := ListRegions()
	c.Assert(err, IsNil)
	c.Assert(regions, DeepEquals, []string{})
	regionManagers, err := ListManagersInRegion("dev")
	c.Assert(err, Not(IsNil))
	c.Assert(regionManagers, DeepEquals, []string{})
	managers, err := ListManagers()
	c.Assert(err, IsNil)
	c.Assert(managers, DeepEquals, map[string][]string{})
	devOpt := Manager("dev", "2.1.1.1")
	c.Assert(devOpt.Save(), IsNil)
	devOtherOpt := Manager("dev", "2.1.1.2")
	c.Assert(devOtherOpt.Save(), IsNil)
	omgOpt := Manager("omg", "2.1.1.3")
	c.Assert(omgOpt.Save(), IsNil)
	bbqOpt := Manager("bbq", "2.1.1.4")
	c.Assert(bbqOpt.Save(), IsNil)
	regions, err = ListRegions()
	c.Assert(err, IsNil)
	sort.Strings(regions)
	c.Assert(regions, DeepEquals, []string{"bbq", "dev", "omg"})
	devManagers, err := ListManagersInRegion("dev")
	c.Assert(err, IsNil)
	sort.Strings(devManagers)
	c.Assert(devManagers, DeepEquals, []string{"2.1.1.1", "2.1.1.2"})
	managers, err = ListManagers()
	c.Assert(err, IsNil)
	for _, regionManagers := range managers {
		sort.Strings(regionManagers)
	}
	c.Assert(managers, DeepEquals, map[string][]string{"bbq": []string{"2.1.1.4"},
		"dev": []string{"2.1.1.1", "2.1.1.2"}, "omg": []string{"2.1.1.3"}})
}
