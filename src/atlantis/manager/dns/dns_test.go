package dns

import (
	. "launchpad.net/gocheck"
	"testing"
)

func TestDatamodel(t *testing.T) { TestingT(t) }

type DNSSuite struct{}

var _ = Suite(&DNSSuite{})

func (s *DNSSuite) TestIPRegexp(c *C) {
	c.Assert(IPRegexp.MatchString("ec2-75-101-190-248.compute-1.amazonaws.com"), Equals, false)
	c.Assert(IPRegexp.MatchString("1.1.1.1"), Equals, true)
	c.Assert(IPRegexp.MatchString("a.1.1.1"), Equals, false)
	c.Assert(IPRegexp.MatchString("1.a.1.1"), Equals, false)
	c.Assert(IPRegexp.MatchString("1.1.a.1"), Equals, false)
	c.Assert(IPRegexp.MatchString("1.1.1.a"), Equals, false)
	c.Assert(IPRegexp.MatchString("12.12.12.12"), Equals, true)
	c.Assert(IPRegexp.MatchString("123.123.123.123"), Equals, true)
	c.Assert(IPRegexp.MatchString("q123.123.123.123"), Equals, false)
	c.Assert(IPRegexp.MatchString("123.123.123.123q"), Equals, false)
}
