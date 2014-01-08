package datamodel

import (
	. "atlantis/manager/constant"
	"atlantis/manager/helper"
	routerzk "atlantis/router/zk"
	zookeeper "github.com/jigish/gozk-recipes"
)

var Zk *zookeeper.ZkConn

func CreateRouterPortsPaths() {
	Zk.Touch(helper.GetBaseRouterPortsPath(true))
	Zk.Touch(helper.GetBaseRouterPortsPath(false))
}

func CreateRouterPaths() {
	helper.SetRouterRoot(true)
	for _, path := range routerzk.ZkPaths {
		Zk.Touch(path)
	}
	helper.SetRouterRoot(false)
	for _, path := range routerzk.ZkPaths {
		Zk.Touch(path)
	}
	for _, zone := range AvailableZones {
		Zk.Touch(helper.GetBaseRouterPath(true, zone))
		Zk.Touch(helper.GetBaseRouterPath(false, zone))
	}
}

func CreateLockPaths() {
	Zk.Touch(helper.GetBaseLockPath("deploy"))
	Zk.Touch(helper.GetBaseLockPath("router_ports_internal"))
	Zk.Touch(helper.GetBaseLockPath("router_ports_external"))
}

func CreateAppPath() {
	Zk.Touch(helper.GetBaseAppPath())
}

func CreateInstancePaths() {
	Zk.Touch(helper.GetBaseInstancePath())
	Zk.Touch(helper.GetBaseInstanceDataPath())
}

func CreateSupervisorPath() {
	Zk.Touch(helper.GetBaseSupervisorPath())
}

func CreateManagerPath() {
	Zk.Touch(helper.GetBaseManagerPath())
}

func CreateEnvPath() {
	Zk.Touch(helper.GetBaseEnvPath())
}

func CreatePaths() {
	CreateRouterPortsPaths()
	CreateRouterPaths()
	CreateLockPaths()
	CreateInstancePaths()
	CreateAppPath()
	CreateSupervisorPath()
	CreateManagerPath()
	CreateEnvPath()
}

func Init(zkUri string) {
	Zk = zookeeper.GetPanicingZk(zkUri)
	CreatePaths()
}
