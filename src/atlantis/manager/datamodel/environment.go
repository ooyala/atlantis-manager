package datamodel

import (
	. "atlantis/common"
	"atlantis/manager/helper"
	"errors"
	"fmt"
	"log"
)

// NOTE[jigish]: all methods in this file expect that the value of the dependencies are already encrypted.
//               also that environment parents actually exist

type ZkEnv struct {
	Name   string
	Parent string
}

func GetEnv(name string) (*ZkEnv, error) {
	e := &ZkEnv{name, ""}
	err := e.Get()
	return e, err
}

func Env(name, parent string) *ZkEnv {
	return &ZkEnv{name, parent}
}

func (e *ZkEnv) Save() error {
	return setJson(e.path(), e)
}

func (e *ZkEnv) Delete() error {
	return Zk.RecursiveDelete(e.path())
}

func (e *ZkEnv) Get() error {
	return getJson(e.path(), e)
}

func (e *ZkEnv) path() string {
	return helper.GetBaseEnvPath(e.Name)
}

func (e *ZkEnv) depPath(name string) string {
	return helper.GetBaseDepPath(e.Name, name)
}

type ZkDep struct {
	Name  string
	Value string
}

func (e *ZkEnv) UpdateDep(name, value string) error {
	return setJson(e.depPath(name), ZkDep{name, value})
}

func (e *ZkEnv) DeleteDep(name string) error {
	return Zk.Delete(e.depPath(name), -1)
}

func (e *ZkEnv) GetDepValue(name string) (value string, err error) {
	dep := &ZkDep{}
	err = getJson(e.depPath(name), dep)
	return dep.Value, err
}

func (e *ZkEnv) AllDepValues() (map[string]string, error) {
	deps, err := ListDeps(e.Name)
	if err != nil {
		return nil, err
	}
	return e.DepValues(deps)
}

func (e *ZkEnv) DepValues(deps []string) (map[string]string, error) {
	values := map[string]string{}
	for _, dep := range deps {
		value, err := e.GetDepValue(dep)
		if err != nil {
			continue // get all we can
		}
		values[dep] = value
	}
	return values, nil
}

func (e *ZkEnv) ResolveAllDepValues() (map[string]string, error) {
	zkEnv := e
	values, err := zkEnv.AllDepValues()
	if err != nil {
		return values, err
	}
	seenParents := map[string]bool{}
	for zkEnv.Parent != "" {
		if seenParents[zkEnv.Parent] {
			log.Printf("ERROR: Cyclical parent (%s) for environment (%s).", zkEnv.Parent, e.Name)
			break
		}
		seenParents[zkEnv.Parent] = true
		zkEnv, err = GetEnv(zkEnv.Parent)
		if err != nil {
			return values, err
		}
		newValues, err := zkEnv.AllDepValues()
		if err != nil {
			return values, err
		}
		for k, v := range values {
			newValues[k] = v
		}
		values = newValues
	}
	valueKeys := make([]string, len(values))
	i := 0
	for k, _ := range values {
		valueKeys[i] = k
		i++
	}
	return values, nil
}

func (e *ZkEnv) ResolveDepValues(deps []string) (map[string]string, error) {
	zkEnv := e
	values, err := zkEnv.DepValues(deps)
	if err != nil {
		return values, err
	}
	seenParents := map[string]bool{}
	for zkEnv.Parent != "" {
		if seenParents[zkEnv.Parent] {
			log.Printf("ERROR: Cyclical parent (%s) for environment (%s).", zkEnv.Parent, e.Name)
			break
		}
		seenParents[zkEnv.Parent] = true
		zkEnv, err = GetEnv(zkEnv.Parent)
		if err != nil {
			return values, err
		}
		newValues, err := zkEnv.DepValues(deps)
		if err != nil {
			return values, err
		}
		for k, v := range values {
			newValues[k] = v
		}
		values = newValues
	}
	valueKeys := make([]string, len(values))
	i := 0
	for k, _ := range values {
		valueKeys[i] = k
		i++
	}
	only1, _ := DiffSlices(deps, valueKeys)
	if len(only1) > 0 {
		return values, errors.New(fmt.Sprintf("Could not fine all deps in %s. Missing %v.", e.Name, only1))
	}
	return values, nil
}

func ListEnvs() (envs []string, err error) {
	envs, _, err = Zk.Children(helper.GetBaseEnvPath())
	if envs == nil {
		return []string{}, err
	}
	return
}

func ListDeps(env string) (deps []string, err error) {
	deps, _, err = Zk.Children(helper.GetBaseEnvPath(env))
	if deps == nil {
		return []string{}, err
	}
	return
}
