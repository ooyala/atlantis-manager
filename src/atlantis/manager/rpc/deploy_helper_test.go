package rpc

import (
	. "atlantis/common"
	"atlantis/crypto"
	. "atlantis/manager/constant"
	"atlantis/manager/datamodel"
	"atlantis/manager/dns"
	"atlantis/manager/helper"
	. "atlantis/manager/rpc/types"
	scrypto "atlantis/supervisor/crypto"
	"fmt"
	zookeeper "github.com/jigish/gozk-recipes"
	. "launchpad.net/gocheck"
)

type DeployHelperSuite struct{}

var _ = Suite(&DeployHelperSuite{})

type FakeDNSProvider bool

func (f FakeDNSProvider) CreateRecords(region, comment string, arecords []dns.Record) error {
	return nil
}

func (f FakeDNSProvider) DeleteRecords(region, comment string, ids ...string) (error, chan error) {
	errChan := make(chan error)
	go func(ch chan error) {
		ch <- nil
	}(errChan)
	return nil, errChan
}

func (f FakeDNSProvider) GetRecordsForValue(region, value string) ([]string, error) {
	return []string{}, nil
}

func (f FakeDNSProvider) CreateHealthCheck(ip string, port uint16) (string, error) {
	return "", nil
}

func (f FakeDNSProvider) DeleteHealthCheck(id string) error {
	return nil
}

func (f FakeDNSProvider) Suffix(region string) (string, error) {
	return Region + ".suffix.com", nil
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
	datamodel.Zk.RecursiveDelete(helper.GetBaseRouterPath(true))
	datamodel.Zk.RecursiveDelete(helper.GetBaseRouterPath(false))
	datamodel.Zk.RecursiveDelete(helper.GetBaseRouterPortsPath(true))
	datamodel.Zk.RecursiveDelete(helper.GetBaseRouterPortsPath(false))
	datamodel.CreateEnvPath()
	datamodel.CreateRouterPaths()
	datamodel.Router(true, "dev", "somehost", "1.2.3.4").Save()
	zkEnv := datamodel.Env("root")
	err := zkEnv.Save()
	c.Assert(err, IsNil)
	deps, err := ResolveDepValues("app", zkEnv, []string{"hello-go"}, false, &Task{})
	c.Assert(err, Not(IsNil))
	_, err = datamodel.CreateInstance("hello-go", "1234567890", "root", "myhost")
	c.Assert(err, IsNil)
	_, err = datamodel.CreateOrUpdateApp(false, false, "app", "ssh://github.com/ooyala/apo", "/", "jigish@ooyala.com")
	c.Assert(err, IsNil)
	zkApp, err := datamodel.CreateOrUpdateApp(false, true, "hello-go", "ssh://github.com/ooyala/hello-go", "/", "jigish@ooyala.com")
	c.Assert(err, IsNil)
	c.Assert(zkApp.AddDependerAppData(&DependerAppData{Name: "app", DependerEnvData: map[string]*DependerEnvData{"root": &DependerEnvData{Name: "root"}}}), IsNil)
	deps, err = ResolveDepValues("app", zkEnv, []string{"hello-go"}, false, &Task{})
	c.Assert(err, IsNil)
	c.Assert(deps["dev1"]["hello-go"].DataMap["address"], Equals, fmt.Sprintf("internal-router.1.%s.suffix.com:%d", Region, datamodel.MinRouterPort))
	deps, err = ResolveDepValues("app", zkEnv, []string{"hello-go"}, true, &Task{})
	c.Assert(err, IsNil)
	c.Assert(deps["dev1"]["hello-go"].EncryptedData, Not(Equals), "")
	c.Assert(deps["dev1"]["hello-go"].DataMap, IsNil)
	scrypto.DecryptAppDep(deps["dev1"]["hello-go"])
	c.Assert(deps["dev1"]["hello-go"].DataMap, Not(IsNil))
	c.Assert(deps["dev1"]["hello-go"].DataMap["address"], Equals, fmt.Sprintf("internal-router.1.%s.suffix.com:%d", Region, datamodel.MinRouterPort))
}
