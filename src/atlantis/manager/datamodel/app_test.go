package datamodel

import (
	"atlantis/manager/helper"
	. "launchpad.net/gocheck"
)

func (s *DatamodelSuite) TestApp(c *C) {
	Zk.RecursiveDelete(helper.GetBaseAppPath())
	apps, err := ListRegisteredApps()
	c.Assert(err, Not(IsNil)) // the path doesn't exist. this is an error
	c.Assert(len(apps), Equals, 0)
	app1, err := GetApp(app)
	c.Assert(err, Not(IsNil))
	app1, err = CreateOrUpdateApp(app, repo, root)
	c.Assert(err, IsNil)
	c.Assert(app1.Name, Equals, app)
	c.Assert(app1.Repo, Equals, repo)
	c.Assert(app1.Root, Equals, root)
	apps, err = ListRegisteredApps()
	c.Assert(err, IsNil)
	c.Assert(len(apps), Equals, 1)
	c.Assert(apps[0], Equals, app)
	app1, err = GetApp(app)
	c.Assert(err, IsNil)
	c.Assert(app1.Name, Equals, app)
	c.Assert(app1.Repo, Equals, repo)
	c.Assert(app1.Root, Equals, root)
	app1, err = CreateOrUpdateApp(app, repo+"2", root+"2")
	c.Assert(err, IsNil)
	c.Assert(app1.Name, Equals, app)
	c.Assert(app1.Repo, Equals, repo+"2")
	c.Assert(app1.Root, Equals, root+"2")
	apps, err = ListRegisteredApps()
	c.Assert(err, IsNil)
	c.Assert(len(apps), Equals, 1)
	c.Assert(apps[0], Equals, app)
	app1, err = GetApp(app)
	c.Assert(err, IsNil)
	c.Assert(app1.Name, Equals, app)
	c.Assert(app1.Repo, Equals, repo+"2")
	c.Assert(app1.Root, Equals, root+"2")
	err = app1.Delete()
	c.Assert(err, IsNil)
	apps, err = ListRegisteredApps()
	c.Assert(err, IsNil) // the path exists, not an error
	c.Assert(len(apps), Equals, 0)
	app1, err = GetApp(app)
	c.Assert(err, Not(IsNil))
}
