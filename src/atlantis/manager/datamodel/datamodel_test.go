package datamodel

import (
	zookeeper "github.com/jigish/gozk-recipes"
	. "launchpad.net/gocheck"
	"testing"
)

const (
	app  = "my-app"
	sha  = "mysha"
	env  = "myenv"
	host = "my-host"
	pool = "my-pool"
	rule = "my-rule"
	trie = "my-trie"
	dep  = "my-dep"
	opt  = "my-opt"
	repo = "my-repo"
	root = "my-root"
)

func TestDatamodel(t *testing.T) { TestingT(t) }

type DatamodelSuite struct{}

var _ = Suite(&DatamodelSuite{})

var (
	zkTestServer *zookeeper.ZkTestServer
)

func (s *DatamodelSuite) SetUpSuite(c *C) {
	zkTestServer = zookeeper.NewZkTestServer()
	c.Assert(zkTestServer.Init(), IsNil)
	Zk = zkTestServer.Zk
}

func (s *DatamodelSuite) TearDownSuite(c *C) {
	err := zkTestServer.Destroy()
	c.Assert(err, IsNil)
}
