package rpc

import (
	aldap "atlantis/manager/ldap"
	"github.com/mavricknz/ldap"
	. "launchpad.net/gocheck"
)

type LDAPSuite struct{}

var _ = Suite(&LDAPSuite{})

func (s *LDAPSuite) TestLDAPAuthorization(c *C) {
	// Has App Permissions
	posAllowedApp := map[string]string{"testApp": "testApp"}
	negAllowedApp := map[string]string{"noop": "noop"}
	c.Assert(ProcessAppPermission("testApp", posAllowedApp), Equals, true)
	c.Assert(ProcessAppPermission("testApp", negAllowedApp), Equals, false)
	return
}

func (s *LDAPSuite) TestIsTeamAdmin(c *C) {
	posTeamAdmin := ldap.EntryAttribute{aldap.TeamAdminAttr, []string{"cn=user"}}
	negTeamAdmin := ldap.EntryAttribute{aldap.TeamAdminAttr, []string{"cn=blah"}}
	negEntry := ldap.Entry{"cn=team", []*ldap.EntryAttribute{&negTeamAdmin}}
	entry := ldap.Entry{"cn=team", []*ldap.EntryAttribute{&posTeamAdmin}}
	pentries := []*ldap.Entry{&entry}
	nentries := []*ldap.Entry{&negEntry}
	psr := ldap.SearchResult{pentries, []string{}, []ldap.Control{}}
	nsr := ldap.SearchResult{nentries, []string{}, []ldap.Control{}}
	c.Assert(ProcessTeamAdmin("cn=user", &psr), Equals, true)
	c.Assert(ProcessTeamAdmin("cn=user", &nsr), Equals, false)
	return
}

func (s *LDAPSuite) TestExistence(c *C) {
	attr := ldap.EntryAttribute{"something", []string{"something"}}
	entry := ldap.Entry{"cn=something", []*ldap.EntryAttribute{&attr}}
	entries := []*ldap.Entry{&entry}
	noEntries := []*ldap.Entry{}
	sr := ldap.SearchResult{entries, []string{}, []ldap.Control{}}
	ss := ldap.SearchResult{noEntries, []string{}, []ldap.Control{}}
	c.Assert(ExistenceTest(sr), Equals, true)
	c.Assert(ExistenceTest(ss), Equals, false)
	return
}

func ExistenceTest(sr ldap.SearchResult) bool {
	if len(sr.Entries) > 0 {
		return true
	}
	return false
}
