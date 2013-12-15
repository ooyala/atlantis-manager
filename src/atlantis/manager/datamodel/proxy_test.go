package datamodel

import (
	"atlantis/manager/helper"
	. "launchpad.net/gocheck"
)

func (s *DatamodelSuite) TestProxy(c *C) {
	MinProxyPort = uint16(65533)
	MaxProxyPort = uint16(65535)
	Zk.RecursiveDelete(helper.GetBaseProxyPath())
	zp := GetProxy()
	c.Assert(zp.AppMap, DeepEquals, map[string]string{})
	c.Assert(zp.PortMap, DeepEquals, map[string]ZkProxyAppEnv{})
	c.Assert(zp.AddAppEnv("app1", "prod"), IsNil)
	c.Assert(zp.AddAppEnv("app1", "prod"), IsNil)
	c.Assert(zp.AddAppEnv("app2", "prod"), IsNil)
	c.Assert(zp.AddAppEnv("app2", "staging"), IsNil)
	c.Assert(zp.AddAppEnv("app3", "staging"), Not(IsNil)) // no more ports
	c.Assert(len(zp.AppMap), Equals, 3)
	c.Assert(len(zp.PortMap), Equals, 3)
	c.Assert(zp.RemoveAppEnv("app1", "prod"), IsNil)
	c.Assert(len(zp.AppMap), Equals, 2)
	c.Assert(len(zp.PortMap), Equals, 2)
	c.Assert(zp.RemoveAppEnv("app1", "prod"), IsNil)
	c.Assert(len(zp.AppMap), Equals, 2)
	c.Assert(len(zp.PortMap), Equals, 2)
	c.Assert(zp.RemoveAppEnv("app2", "prod"), IsNil)
	c.Assert(len(zp.AppMap), Equals, 1)
	c.Assert(len(zp.PortMap), Equals, 1)
	c.Assert(zp.RemoveAppEnv("app2", "staging"), IsNil)
	c.Assert(len(zp.AppMap), Equals, 0)
	c.Assert(len(zp.PortMap), Equals, 0)
	c.Assert(zp.AddAll("app1", []string{"prod", "staging"}), IsNil)
	c.Assert(len(zp.AppMap), Equals, 2)
	c.Assert(len(zp.PortMap), Equals, 2)
	c.Assert(zp.RemoveAll("app1", []string{"prod", "staging"}), IsNil)
	c.Assert(len(zp.AppMap), Equals, 0)
	c.Assert(len(zp.PortMap), Equals, 0)
}
