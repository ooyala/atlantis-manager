package datamodel

import (
	"atlantis/manager/helper"
	. "launchpad.net/gocheck"
)

func (s *DatamodelSuite) TestDNSModel(c *C) {
	Zk.RecursiveDelete(helper.GetBaseDNSPath())
	zkDNS := DNS(app, env)
	err := zkDNS.Save()
	c.Assert(err, IsNil)
	fetchedDNS, err := GetDNS(app, env)
	c.Assert(err, IsNil)
	c.Assert(zkDNS, DeepEquals, fetchedDNS)
	zkDNS.Shas = map[string]bool{"sha1": true, "sha2": true}
	zkDNS.RecordIDs = []string{"rid1", "rid2"}
	zkDNS.Save()
	fetchedDNS, err = GetDNS(app, env)
	c.Assert(err, IsNil)
	c.Assert(zkDNS, DeepEquals, fetchedDNS)
	err = zkDNS.Delete()
	c.Assert(err, IsNil)
	_, err = GetDNS(app, env)
	c.Assert(err, Not(IsNil))
}
