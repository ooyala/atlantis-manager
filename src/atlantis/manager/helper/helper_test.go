package helper

import (
	. "atlantis/manager/constant"
	routerzk "atlantis/router/zk"
	. "launchpad.net/gocheck"
	"testing"
)

const (
	app  = "my-app"
	sha  = "mysha"
	env  = "myenv"
	host = "my-host"
	cont = "my-app-mysha-myenv-123556"
	pool = "my-pool"
	rule = "my-rule"
	trie = "my-trie"
	dep  = "my-dep"
	opt  = "my-opt"
	zone = "my-zone"
)

func TestDatamodel(t *testing.T) { TestingT(t) }

type HelperSuite struct{}

var _ = Suite(&HelperSuite{})

func (s *HelperSuite) TestHelperContainerNaming(c *C) {
	name := CreateContainerId(app, sha, env)
	c.Assert(name, Matches, app+"."+sha+"."+env+".*")
	name = CreateContainerId(app, "1234567", env)
	c.Assert(name, Matches, app+".123456."+env+".*")
}

func (s *HelperSuite) TestHelperAppPath(c *C) {
	c.Assert(GetBaseAppPath(), Equals, "/atlantis/apps/"+Region)
	c.Assert(GetBaseAppPath(app), Equals, "/atlantis/apps/"+Region+"/"+app)
}

func (s *HelperSuite) TestHelperInstancePath(c *C) {
	c.Assert(GetBaseInstancePath(), Equals, "/atlantis/instances/"+Region)
	c.Assert(GetBaseInstancePath(app), Equals, "/atlantis/instances/"+Region+"/"+app)
	c.Assert(GetBaseInstancePath(app, sha), Equals, "/atlantis/instances/"+Region+"/"+app+"/"+sha)
	c.Assert(GetBaseInstancePath(app, sha, env), Equals,
		"/atlantis/instances/"+Region+"/"+app+"/"+sha+"/"+env)
	c.Assert(GetBaseInstancePath(app, sha, env, cont), Equals,
		"/atlantis/instances/"+Region+"/"+app+"/"+sha+"/"+env+"/"+cont)
}

func (s *HelperSuite) TestHelperInstanceDataPath(c *C) {
	c.Assert(GetBaseInstanceDataPath(), Equals, "/atlantis/instance_data/"+Region)
	c.Assert(GetBaseInstanceDataPath(cont), Equals, "/atlantis/instance_data/"+Region+"/"+cont)
}

func (s *HelperSuite) TestHelperHostPath(c *C) {
	c.Assert(GetBaseHostPath(), Equals, "/atlantis/hosts/"+Region)
	c.Assert(GetBaseHostPath(host), Equals, "/atlantis/hosts/"+Region+"/"+host)
	c.Assert(GetBaseHostPath(host, cont), Equals, "/atlantis/hosts/"+Region+"/"+host+"/"+cont)
}

func (s *HelperSuite) TestHelperPoolName(c *C) {
	c.Assert(CreatePoolName(app, sha, env), Matches, app+"."+sha+"."+env)
}

func (s *HelperSuite) TestHelperRouterRoot(c *C) {
	SetRouterRoot(true)
	c.Assert(routerzk.ZkPaths["pools"], Equals, "/atlantis/router/"+Region+"/internal/pools")
	c.Assert(routerzk.ZkPaths["rules"], Equals, "/atlantis/router/"+Region+"/internal/rules")
	c.Assert(routerzk.ZkPaths["tries"], Equals, "/atlantis/router/"+Region+"/internal/tries")
	SetRouterRoot(false)
	c.Assert(routerzk.ZkPaths["pools"], Equals, "/atlantis/router/"+Region+"/external/pools")
	c.Assert(routerzk.ZkPaths["rules"], Equals, "/atlantis/router/"+Region+"/external/rules")
	c.Assert(routerzk.ZkPaths["tries"], Equals, "/atlantis/router/"+Region+"/external/tries")
}

func (s *HelperSuite) TestHelperRouterPath(c *C) {
	c.Assert(GetBaseRouterPath(), Equals, "/atlantis/routers/"+Region)
	c.Assert(GetBaseRouterPath(zone), Equals, "/atlantis/routers/"+Region+"/"+zone)
	c.Assert(GetBaseRouterPath(zone, host), Equals, "/atlantis/routers/"+Region+"/"+zone+"/"+host)
}

func (s *HelperSuite) TestHelperDNSPath(c *C) {
	c.Assert(GetBaseDNSPath(), Equals, "/atlantis/dns/"+Region)
	c.Assert(GetBaseDNSPath(app), Equals, "/atlantis/dns/"+Region+"/"+app)
	c.Assert(GetBaseDNSPath(app, env), Equals, "/atlantis/dns/"+Region+"/"+app+"/"+env)
}

func (s *HelperSuite) TestHelperManagerPath(c *C) {
	c.Assert(GetBaseManagerPath(), Equals, "/atlantis/managers")
	c.Assert(GetBaseManagerPath(Region), Equals, "/atlantis/managers/"+Region)
	c.Assert(GetBaseManagerPath(Region, host), Equals, "/atlantis/managers/"+Region+"/"+host)
}

func (s *HelperSuite) TestHelperDepPath(c *C) {
	c.Assert(GetBaseDepPath(env, dep), Equals, "/atlantis/environments/"+Region+"/"+env+"/"+dep)
}

func (s *HelperSuite) TestHelperEnvPath(c *C) {
	c.Assert(GetBaseEnvPath(), Equals, "/atlantis/environments/"+Region)
	c.Assert(GetBaseEnvPath(env), Equals, "/atlantis/environments/"+Region+"/"+env)
}

func (s *HelperSuite) TestHelperLockPath(c *C) {
	c.Assert(GetBaseLockPath(), Equals, "/atlantis/lock/"+Region)
	c.Assert(GetBaseLockPath("deploy"), Equals, "/atlantis/lock/"+Region+"/deploy")
}

func (s *HelperSuite) TestGetRegionRouterCName(c *C) {
	c.Assert(GetRegionRouterCName("atlantis.com"), Equals, "router."+Region+".atlantis.com")
}

func (s *HelperSuite) TestGetZoneRouterCName(c *C) {
	c.Assert(GetZoneRouterCName(Region+"1", "atlantis.com"), Equals, "router."+Region+"1.atlantis.com")
}

func (s *HelperSuite) TestGetRouterCName(c *C) {
	c.Assert(GetRouterCName(1, Region+"1", "atlantis.com"), Equals, "router1."+Region+"1.atlantis.com")
}

func (s *HelperSuite) TestGetRegionAppAlias(c *C) {
	c.Assert(GetRegionAppAlias(app, env, "atlantis.com"), Equals, app+"."+env+"."+Region+".atlantis.com")
}

func (s *HelperSuite) TestGetZoneAppAlias(c *C) {
	c.Assert(GetZoneAppAlias(app, env, Region+"1",  "atlantis.com"), Equals,
		app+"."+env+"."+Region+"1.atlantis.com")
}

func (s *HelperSuite) TestHelperRegionAndZone(c *C) {
	c.Assert(RegionAndZone(Region+"1"), Equals, Region+"1")
	c.Assert(RegionAndZone("1"), Equals, Region+"1")
}
