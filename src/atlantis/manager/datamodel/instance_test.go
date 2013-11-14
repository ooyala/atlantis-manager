package datamodel

import (
	"atlantis/manager/helper"
	"fmt"
	. "launchpad.net/gocheck"
	"sort"
)

func (s *DatamodelSuite) TestInstance(c *C) {
	Zk.RecursiveDelete(helper.GetBaseInstancePath())
	externalInst, err := CreateInstance(false, app, sha, env, host)
	c.Assert(err, IsNil)
	c.Assert(externalInst.Id, Not(Equals), "")
	c.Assert(externalInst.App, Equals, app)
	c.Assert(externalInst.Sha, Equals, sha)
	c.Assert(externalInst.Host, Equals, host)
	c.Assert(externalInst.Internal, Equals, false)
	c.Assert(externalInst.SetPort(uint16(1337)), IsNil)
	otherExtInst, err := GetInstance(externalInst.Id)
	c.Assert(*otherExtInst, Equals, *externalInst)
	err = externalInst.Delete()
	c.Assert(err, IsNil)
	internalInst, err := CreateInstance(true, app, sha, env, host)
	c.Assert(err, IsNil)
	c.Assert(internalInst.Id, Not(Equals), "")
	c.Assert(internalInst.App, Equals, app)
	c.Assert(internalInst.Sha, Equals, sha)
	c.Assert(internalInst.Host, Equals, host)
	c.Assert(internalInst.Internal, Equals, true)
	otherInttInst, err := GetInstance(internalInst.Id)
	c.Assert(*otherInttInst, Equals, *internalInst)
	err = internalInst.Delete()
	c.Assert(err, IsNil)
}

func (s *DatamodelSuite) TestInstanceListers(c *C) {
	Zk.RecursiveDelete(helper.GetBaseInstancePath())
	apps := []string{"app1", "app2", "app3"}
	appsMap := map[string][]string{}
	appsMap["app1"] = []string{"app1sha1", "app1sha2", "app1sha3"}
	appsMap["app2"] = []string{"app2sha1"}
	appsMap["app3"] = []string{"app3sha1"}
	shas := map[string]map[string][]string{}
	shas["app1sha1"] = map[string][]string{"env1": []string{"", ""}}
	shas["app1sha2"] = map[string][]string{"env2": []string{""}}
	shas["app1sha3"] = map[string][]string{"env1": []string{"", "", ""}, "env2": []string{""}}
	shas["app2sha1"] = map[string][]string{"env2": []string{""}}
	shas["app3sha1"] = map[string][]string{"env3": []string{""}}
	envs := map[string][]string{}
	envs["app1sha1"] = []string{"env1"}
	envs["app1sha2"] = []string{"env2"}
	envs["app1sha3"] = []string{"env1", "env2"}
	envs["app2sha1"] = []string{"env2"}
	envs["app3sha1"] = []string{"env3"}
	instances := []*ZkInstance{}
	for _, app := range apps {
		for _, sha := range appsMap[app] {
			for env, hosts := range shas[sha] {
				for i, _ := range hosts {
					inst, err := CreateInstance(false, app, sha, env, fmt.Sprintf("%s-%d", host, i))
					c.Assert(err, IsNil)
					shas[sha][env][i] = inst.Id
					instances = append(instances, inst)
				}
				sort.Strings(shas[sha][env])
			}
		}
	}

	getApps, err := ListApps()
	c.Assert(err, IsNil)
	sort.Strings(getApps)
	c.Assert(getApps, DeepEquals, apps)
	for _, app := range getApps {
		getShas, err := ListShas(app)
		c.Assert(err, IsNil)
		sort.Strings(getShas)
		c.Assert(getShas, DeepEquals, appsMap[app])
		for _, sha := range getShas {
			getEnvs, err := ListAppEnvs(app, sha)
			c.Assert(err, IsNil)
			sort.Strings(getEnvs)
			c.Assert(getEnvs, DeepEquals, envs[sha])
			for _, env := range getEnvs {
				getInsts, err := ListInstances(app, sha, env)
				c.Assert(err, IsNil)
				sort.Strings(getInsts)
				c.Assert(getInsts, DeepEquals, shas[sha][env])
			}
		}
	}

	for _, inst := range instances {
		err = inst.Delete()
		c.Assert(err, IsNil)
	}

	getApps, err = ListApps()
	c.Assert(err, IsNil)
	sort.Strings(getApps)
	c.Assert(getApps, DeepEquals, []string{})
}
