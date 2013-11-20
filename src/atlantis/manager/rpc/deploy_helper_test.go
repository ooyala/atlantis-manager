package rpc

import (
	"atlantis/crypto"
	. "atlantis/manager/constant"
	"atlantis/manager/datamodel"
	"atlantis/manager/dns"
	"atlantis/manager/helper"
	zookeeper "github.com/jigish/gozk-recipes"
	. "launchpad.net/gocheck"
)

type DeployHelperSuite struct{}

var _ = Suite(&DeployHelperSuite{})

type FakeDNSProvider bool

func (f FakeDNSProvider) CreateAliases(comment string, aliases []dns.Alias) (error, chan error) {
	errChan := make(chan error)
	go func(ch chan error) {
		ch <- nil
	}(errChan)
	return nil, errChan
}

func (f FakeDNSProvider) CreateCNames(comment string, cnames []dns.CName) (error, chan error) {
	errChan := make(chan error)
	go func(ch chan error) {
		ch <- nil
	}(errChan)
	return nil, errChan
}

func (f FakeDNSProvider) DeleteRecords(comment string, ids ...string) (error, chan error) {
	errChan := make(chan error)
	go func(ch chan error) {
		ch <- nil
	}(errChan)
	return nil, errChan
}

func (f FakeDNSProvider) CreateHealthCheck(ip string) (string, error) {
	return "", nil
}

func (f FakeDNSProvider) DeleteHealthCheck(id string) error {
	return nil
}

func (f FakeDNSProvider) Suffix() string {
	return "suffix.com"
}

var (
	zkTestServer *zookeeper.ZkTestServer
)

func (s *DeployHelperSuite) SetUpSuite(c *C) {
	crypto.Init()
	dns.Provider = FakeDNSProvider(true)
	zkTestServer = zookeeper.NewZkTestServer()
	c.Assert(zkTestServer.Init(), IsNil)
	datamodel.Zk = zkTestServer.Zk
}

func (s *DeployHelperSuite) TearDownSuite(c *C) {
	err := zkTestServer.Destroy()
	c.Assert(err, IsNil)
}

func (s *DeployHelperSuite) TestResolveDepValues(c *C) {
	datamodel.Zk.RecursiveDelete(helper.GetBaseEnvPath())
	datamodel.CreateEnvPath()
	zkEnv := datamodel.Env("root", "")
	err := zkEnv.Save()
	c.Assert(err, IsNil)
	deps, err := ResolveDepValues(zkEnv, []string{"somedep"}, false)
	c.Assert(err, Not(IsNil))
	zkEnv.UpdateDep("somedep", string(crypto.Encrypt([]byte("somevalue"))))
	deps, err = ResolveDepValues(zkEnv, []string{"somedep"}, false)
	c.Assert(err, IsNil)
	c.Assert(deps["somedep"], Equals, "somevalue")
	deps, err = ResolveDepValues(zkEnv, []string{"somedep", "hello-go"}, false)
	c.Assert(err, Not(IsNil))
	_, err = datamodel.CreateInstance(true, "hello-go", "1234567890", "root", "myhost")
	c.Assert(err, IsNil)
	deps, err = ResolveDepValues(zkEnv, []string{"somedep", "hello-go"}, false)
	c.Assert(err, IsNil)
	c.Assert(deps["somedep"], Equals, "somevalue")
	c.Assert(deps["hello-go"], Equals, "hello-go.root."+Region+".suffix.com")
}
