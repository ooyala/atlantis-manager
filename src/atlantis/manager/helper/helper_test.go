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

package helper

import (
	. "atlantis/manager/constant"
	routerzk "atlantis/router/zk"
	. "launchpad.net/gocheck"
	"sort"
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
	name := CreateContainerID(app, sha, env)
	c.Assert(name, Matches, app+"."+sha+"."+env+".*")
	name = CreateContainerID(app, "1234567", env)
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

func (s *HelperSuite) TestHelperSupervisorPath(c *C) {
	c.Assert(GetBaseSupervisorPath(), Equals, "/atlantis/supervisors/"+Region)
	c.Assert(GetBaseSupervisorPath(host), Equals, "/atlantis/supervisors/"+Region+"/"+host)
	c.Assert(GetBaseSupervisorPath(host, cont), Equals, "/atlantis/supervisors/"+Region+"/"+host+"/"+cont)
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
	c.Assert(GetBaseRouterPath(true), Equals, "/atlantis/routers/"+Region+"/internal")
	c.Assert(GetBaseRouterPath(true, zone), Equals, "/atlantis/routers/"+Region+"/internal/"+zone)
	c.Assert(GetBaseRouterPath(true, zone, host), Equals, "/atlantis/routers/"+Region+"/internal/"+zone+"/"+host)
	c.Assert(GetBaseRouterPath(false), Equals, "/atlantis/routers/"+Region+"/external")
	c.Assert(GetBaseRouterPath(false, zone), Equals, "/atlantis/routers/"+Region+"/external/"+zone)
	c.Assert(GetBaseRouterPath(false, zone, host), Equals, "/atlantis/routers/"+Region+"/external/"+zone+"/"+host)
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

func (s *HelperSuite) TestGetManagerCName(c *C) {
	c.Assert(GetManagerCName(1, "us-east-1.atlantis.com"), Equals, "manager1.us-east-1.atlantis.com")
}

func (s *HelperSuite) TestGetRegistryCName(c *C) {
	c.Assert(GetRegistryCName(1, "us-east-1.atlantis.com"), Equals, "registry1.us-east-1.atlantis.com")
}

func (s *HelperSuite) TestGetRegionRouterCName(c *C) {
	c.Assert(GetRegionRouterCName(true, "us-east-1.atlantis.com"), Equals, "internal-router.us-east-1.atlantis.com")
	c.Assert(GetRegionRouterCName(false, "us-east-1.atlantis.com"), Equals, "router.us-east-1.atlantis.com")
}

func (s *HelperSuite) TestGetZoneRouterCName(c *C) {
	c.Assert(GetZoneRouterCName(true, Region+"1", Region+".atlantis.com"), Equals, "internal-router.1."+Region+".atlantis.com")
	c.Assert(GetZoneRouterCName(false, Region+"1", Region+".atlantis.com"), Equals, "router.1."+Region+".atlantis.com")
}

func (s *HelperSuite) TestGetRouterCName(c *C) {
	c.Assert(GetRouterCName(true, 1, Region+"1", Region+".atlantis.com"), Equals, "internal-router1.1."+Region+".atlantis.com")
	c.Assert(GetRouterCName(false, 1, Region+"1", Region+".atlantis.com"), Equals, "router1.1."+Region+".atlantis.com")
}

func (s *HelperSuite) TestGetAppCNameSuffixes(c *C) {
	obtain := GetAppCNameSuffixes("us-east-1.atlantis.com")
	sort.Strings(obtain)
	expect := []string{"1.us-east-1.atlantis.com", "us-east-1.atlantis.com"}
	sort.Strings(expect)
	c.Assert(obtain[0], Equals, expect[0])
	c.Assert(obtain[1], Equals, expect[1])
}

func (s *HelperSuite) TestGetRegionAppCName(c *C) {
	c.Assert(GetRegionAppCName(app, env, "us-east-1.atlantis.com"), Equals, app+"."+env+".us-east-1.atlantis.com")
}

func (s *HelperSuite) TestGetZoneAppCName(c *C) {
	env := "oogabooga"
	c.Assert(GetZoneAppCName(app, env, Region+"1", Region+".atlantis.com"), Equals,
		app+"."+env+".1."+Region+".atlantis.com")
	env = "ooga-booga"
	c.Assert(GetZoneAppCName(app, env, Region+"1", Region+".atlantis.com"), Equals,
		app+"."+env+".1."+Region+".atlantis.com")
	env = "ooga_booga"
	c.Assert(GetZoneAppCName(app, env, Region+"1", Region+".atlantis.com"), Equals,
		app+"."+env+".1."+Region+".atlantis.com")
	env = "prodooga"
	c.Assert(GetZoneAppCName(app, env, Region+"1", Region+".atlantis.com"), Equals,
		app+"."+env+".1."+Region+".atlantis.com")
	env = "productionooga"
	c.Assert(GetZoneAppCName(app, env, Region+"1", Region+".atlantis.com"), Equals,
		app+"."+env+".1."+Region+".atlantis.com")
	env = "prod"
	c.Assert(GetZoneAppCName(app, env, Region+"1", Region+".atlantis.com"), Equals,
		app+".1."+Region+".atlantis.com")
	env = "production"
	c.Assert(GetZoneAppCName(app, env, Region+"1", Region+".atlantis.com"), Equals,
		app+".1."+Region+".atlantis.com")
	env = "prod-ooga"
	c.Assert(GetZoneAppCName(app, env, Region+"1", Region+".atlantis.com"), Equals,
		app+".1."+Region+".atlantis.com")
	env = "prod_ooga"
	c.Assert(GetZoneAppCName(app, env, Region+"1", Region+".atlantis.com"), Equals,
		app+".1."+Region+".atlantis.com")
	env = "production-ooga"
	c.Assert(GetZoneAppCName(app, env, Region+"1", Region+".atlantis.com"), Equals,
		app+".1."+Region+".atlantis.com")
	env = "production_ooga"
	c.Assert(GetZoneAppCName(app, env, Region+"1", Region+".atlantis.com"), Equals,
		app+".1."+Region+".atlantis.com")
}

func (s *HelperSuite) TestHelperRegionAndZone(c *C) {
	c.Assert(RegionAndZone(Region+"1"), Equals, Region+"1")
	c.Assert(RegionAndZone("1"), Equals, Region+"1")
}

func (s *HelperSuite) TestHelperZoneMinusRegion(c *C) {
	c.Assert(ZoneMinusRegion(Region+"1"), Equals, "1")
	c.Assert(ZoneMinusRegion("asdas1"), Equals, "asdas1")
}

func (s *HelperSuite) TestHelperEmptyIfProdPrefix(c *C) {
	c.Assert(EmptyIfProdPrefix("oogabooga"), Equals, ".oogabooga")
	c.Assert(EmptyIfProdPrefix("ooga-booga"), Equals, ".ooga-booga")
	c.Assert(EmptyIfProdPrefix("ooga_booga"), Equals, ".ooga_booga")
	c.Assert(EmptyIfProdPrefix("prodooga"), Equals, ".prodooga")
	c.Assert(EmptyIfProdPrefix("productionooga"), Equals, ".productionooga")
	c.Assert(EmptyIfProdPrefix("prod"), Equals, "")
	c.Assert(EmptyIfProdPrefix("production"), Equals, "")
	c.Assert(EmptyIfProdPrefix("prod-ooga"), Equals, "")
	c.Assert(EmptyIfProdPrefix("prod_ooga"), Equals, "")
	c.Assert(EmptyIfProdPrefix("production-ooga-booga"), Equals, "")
	c.Assert(EmptyIfProdPrefix("production_ooga-booga"), Equals, "")
}
