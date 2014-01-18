package rpc

import (
	. "launchpad.net/gocheck"
)

type EnvironmentSuite struct{}

var _ = Suite(&EnvironmentSuite{})

func (s *EnvironmentSuite) TestEnvironmentTracker(c *C) {
	c.Assert(IsEnvInUse("root"), Equals, false)

	AddAppShaToEnv("app", "sha", "root")
	c.Assert(IsEnvInUse("root"), Equals, true)
	AddAppShaToEnv("app", "sha", "root")
	c.Assert(IsEnvInUse("root"), Equals, true)
	DeleteAppShaFromEnv("app", "sha", "root")
	c.Assert(IsEnvInUse("root"), Equals, true)
	DeleteAppShaFromEnv("app", "sha", "root")
	c.Assert(IsEnvInUse("root"), Equals, false)
	DeleteEnv("root")
	c.Assert(IsEnvInUse("root"), Equals, false)
}
