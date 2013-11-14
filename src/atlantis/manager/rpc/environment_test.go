package rpc

import (
	. "launchpad.net/gocheck"
)

type EnvironmentSuite struct{}

var _ = Suite(&EnvironmentSuite{})

func (s *EnvironmentSuite) TestEnvironmentTracker(c *C) {
	c.Check(IsEnvInUse("root"), Equals, false)
	c.Check(HasChildren("root"), Equals, false)
	UpdateEnv("root", "", "")
	UpdateEnv("branch", "", "root")
	UpdateEnv("leaf", "", "branch")
	c.Check(IsEnvInUse("root"), Equals, false)
	c.Check(HasChildren("root"), Equals, true)
	c.Check(HasChildren("branch"), Equals, true)
	c.Check(HasChildren("leaf"), Equals, false)
	UpdateEnv("leaf", "branch", "root")
	c.Check(HasChildren("root"), Equals, true)
	c.Check(HasChildren("branch"), Equals, false)
	c.Check(HasChildren("leaf"), Equals, false)

	AddAppShaToEnv("app", "sha", "branch")
	AddAppShaToEnv("app", "sha", "branch")
	c.Check(IsEnvInUse("branch"), Equals, true)
	c.Check(IsEnvInUse("root"), Equals, true)
	UpdateEnv("leaf", "root", "branch")
	c.Check(IsEnvInUse("leaf"), Equals, false)
	DeleteAppShaFromEnv("app", "sha", "branch")
	c.Check(IsEnvInUse("branch"), Equals, true)
	c.Check(IsEnvInUse("leaf"), Equals, false)
	c.Check(IsEnvInUse("root"), Equals, true)
	DeleteAppShaFromEnv("app", "sha", "branch")
	c.Check(IsEnvInUse("branch"), Equals, false)
	c.Check(IsEnvInUse("leaf"), Equals, false)
	c.Check(IsEnvInUse("root"), Equals, false)

	DeleteEnv("banch")
	DeleteEnv("root")
	c.Check(HasChildren("root"), Equals, false)
}
