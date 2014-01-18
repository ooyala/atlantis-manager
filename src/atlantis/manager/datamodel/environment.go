package datamodel

import (
	"atlantis/manager/helper"
)

type ZkEnv struct {
	Name string
}

func GetEnv(name string) (*ZkEnv, error) {
	e := &ZkEnv{name}
	err := e.Get()
	return e, err
}

func Env(name string) *ZkEnv {
	return &ZkEnv{name}
}

func (e *ZkEnv) Save() error {
	return setJson(e.path(), e)
}

func (e *ZkEnv) Delete() error {
	if err := ReclaimRouterPortsForEnv(true, e.Name); err != nil {
		return err
	}
	if err := ReclaimRouterPortsForEnv(false, e.Name); err != nil {
		return err
	}
	return Zk.RecursiveDelete(e.path())
}

func (e *ZkEnv) Get() error {
	return getJson(e.path(), e)
}

func (e *ZkEnv) path() string {
	return helper.GetBaseEnvPath(e.Name)
}

func ListEnvs() (envs []string, err error) {
	envs, _, err = Zk.Children(helper.GetBaseEnvPath())
	if envs == nil {
		return []string{}, err
	}
	return
}
