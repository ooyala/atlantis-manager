package datamodel

import (
	"atlantis/manager/helper"
	. "launchpad.net/gocheck"
	"sort"
)

func (s *DatamodelSuite) TestManagerPath(c *C) {
	opt := Manager("dev", "myhost")
	c.Assert(opt.path(), Equals, helper.GetBaseManagerPath("dev", "myhost"))
}

func (s *DatamodelSuite) TestManagerTouchAndDelete(c *C) {
	opt := Manager("dev", "myhost")
	err := opt.Delete()
	c.Assert(err, Not(IsNil))
	c.Assert(opt.Touch(), IsNil)
	err = opt.Delete()
	c.Assert(err, IsNil)
}

func (s *DatamodelSuite) TestManagerListers(c *C) {
	regions, err := ListRegions()
	c.Assert(err, Not(IsNil))
	c.Assert(regions, DeepEquals, []string{})
	regionManagers, err := ListManagersInRegion("dev")
	c.Assert(err, Not(IsNil))
	c.Assert(regionManagers, DeepEquals, []string{})
	managers, err := ListManagers()
	c.Assert(err, Not(IsNil))
	c.Assert(managers, DeepEquals, map[string][]string{})
	devOpt := Manager("dev", "mydevhost")
	c.Assert(devOpt.Touch(), IsNil)
	devOtherOpt := Manager("dev", "myotherdevhost")
	c.Assert(devOtherOpt.Touch(), IsNil)
	omgOpt := Manager("omg", "myomghost")
	c.Assert(omgOpt.Touch(), IsNil)
	bbqOpt := Manager("bbq", "mybbqhost")
	c.Assert(bbqOpt.Touch(), IsNil)
	regions, err = ListRegions()
	c.Assert(err, IsNil)
	sort.Strings(regions)
	c.Assert(regions, DeepEquals, []string{"bbq", "dev", "omg"})
	devManagers, err := ListManagersInRegion("dev")
	c.Assert(err, IsNil)
	sort.Strings(devManagers)
	c.Assert(devManagers, DeepEquals, []string{"mydevhost", "myotherdevhost"})
	managers, err = ListManagers()
	c.Assert(err, IsNil)
	for _, regionManagers := range managers {
		sort.Strings(regionManagers)
	}
	c.Assert(managers, DeepEquals, map[string][]string{"bbq": []string{"mybbqhost"},
		"dev": []string{"mydevhost", "myotherdevhost"}, "omg": []string{"myomghost"}})
}
