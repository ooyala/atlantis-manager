package datamodel

import (
	"atlantis/manager/helper"
	routerzk "atlantis/router/zk"
	zookeeper "github.com/jigish/gozk-recipes"
)

var Zk *zookeeper.ZkConn

func CreateRouterPaths() {
	helper.SetRouterRoot(true)
	for _, path := range routerzk.ZkPaths {
		Zk.Touch(path)
	}
	helper.SetRouterRoot(false)
	for _, path := range routerzk.ZkPaths {
		Zk.Touch(path)
	}
}

func CreateDeployLockPath() {
	Zk.Touch(helper.GetBaseLockPath("deploy"))
}

func CreateAppPath() {
	Zk.Touch(helper.GetBaseAppPath())
}

func CreateInstancePaths() {
	Zk.Touch(helper.GetBaseInstancePath())
	Zk.Touch(helper.GetBaseInstanceDataPath())
}

func CreateSupervisorPath() {
	Zk.Touch(helper.GetBaseHostPath())
}

func CreateManagerPath() {
	Zk.Touch(helper.GetBaseManagerPath())
}

func CreateEnvPath() {
	Zk.Touch(helper.GetBaseEnvPath())
}

func CreatePaths() {
	CreateRouterPaths()
	CreateDeployLockPath()
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
