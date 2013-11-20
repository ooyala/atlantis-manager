package datamodel

import (
	"atlantis/manager/helper"
	"log"
)

type ZkInstance struct {
	Internal bool
	Id       string
	App      string
	Sha      string
	Env      string
	Host     string
	Port     uint16
}

func InstanceExists(id string) bool {
	if stat, err := Zk.Exists(helper.GetBaseInstanceDataPath(id)); err == nil && stat != nil {
		return true
	}
	return false
}

func GetInstance(id string) (zi *ZkInstance, err error) {
	zi = &ZkInstance{}
	err = getJson(helper.GetBaseInstanceDataPath(id), zi)
	return
}

func CreateInstance(internal bool, app, sha, env, host string) (*ZkInstance, error) {
	id := helper.CreateContainerId(app, sha, env)
	for InstanceExists(id) {
		id = helper.CreateContainerId(app, sha, env)
	}
	zi := &ZkInstance{internal, id, app, sha, env, host, 0}
	if _, err := Zk.Touch(zi.path()); err != nil {
		Zk.RecursiveDelete(zi.path())
		return zi, err
	}
	if err := setJson(zi.dataPath(), zi); err != nil {
		// clean up
		Zk.RecursiveDelete(zi.path())
		Zk.RecursiveDelete(zi.dataPath())
		return zi, err
	}
	return zi, nil
}

func (zi *ZkInstance) Delete() (bool, error) { // true if this was the last instance of app+sha+env
	var (
		last bool
		err  error
		err2 error
		list []string
	)
	// try to get the data (its ok if we can't)
	dataErr := getJson(zi.dataPath(), zi)
	err = Zk.RecursiveDelete(zi.dataPath())
	err2 = Zk.RecursiveDelete(zi.path())
	if err != nil {
		return last, err
	}
	if err2 != nil {
		return last, err2
	}
	// check if we can delete parent directories
	if list, err = ListInstances(zi.App, zi.Sha, zi.Env); err != nil {
		log.Printf("Warning: clean up fail during instance delete: %s", err)
		return last, nil // this is extra, no need to return the error if we couldn't get them
	} else if list != nil && len(list) > 0 {
		return last, nil
	}
	last = true
	Zk.RecursiveDelete(helper.GetBaseInstancePath(zi.App, zi.Sha, zi.Env))
	if dataErr != nil {
		log.Printf("Warning: could not fetch data to clean up pool: %s", err)
	}
	if list, err = ListAppEnvs(zi.App, zi.Sha); err != nil {
		log.Printf("Warning: clean up fail during instance delete: %s", err)
		return last, nil // this is extra, no need to return the error if we couldn't get them
	} else if list != nil && len(list) > 0 {
		return last, nil
	}
	Zk.RecursiveDelete(helper.GetBaseInstancePath(zi.App, zi.Sha))
	// no need to kill pools, they should have been cleaned up when we deleted the instances
	if list, err = ListShas(zi.App); err != nil {
		log.Printf("Warning: clean up fail during instance delete: %s", err)
		return last, nil // this is extra, no need to return the error if we couldn't get them
	} else if list != nil && len(list) > 0 {
		return last, nil
	}
	Zk.RecursiveDelete(helper.GetBaseInstancePath(zi.App))
	// no need to kill pools, they should have been cleaned up when we deleted the instances
	return last, nil
}

func (zi *ZkInstance) SetPort(port uint16) error {
	zi.Port = port
	return setJson(zi.dataPath(), zi)
}

func (zi *ZkInstance) path() string {
	return helper.GetBaseInstancePath(zi.App, zi.Sha, zi.Env, zi.Id)
}

func (zi *ZkInstance) dataPath() string {
	return helper.GetBaseInstanceDataPath(zi.Id)
}

func ListApps() (apps []string, err error) {
	apps, _, err = Zk.VisibleChildren(helper.GetBaseInstancePath())
	if err != nil {
		log.Printf("Error getting list of apps. Error: %s.", err.Error())
	}
	if apps == nil {
		log.Println("No apps found. Returning empty list.")
		apps = []string{}
	}
	return
}

func ListShas(app string) (shas []string, err error) {
	shas, _, err = Zk.VisibleChildren(helper.GetBaseInstancePath(app))
	if err != nil {
		log.Printf("Error getting list of shas. Error: %s.", err.Error())
	}
	if shas == nil {
		log.Println("No shas found. Returning empty list.")
		shas = []string{}
	}
	return
}

func ListAppEnvs(app, sha string) (envs []string, err error) {
	envs, _, err = Zk.VisibleChildren(helper.GetBaseInstancePath(app, sha))
	if err != nil {
		log.Printf("Error getting list of shas. Error: %s.", err.Error())
	}
	if envs == nil {
		log.Println("No shas found. Returning empty list.")
		envs = []string{}
	}
	return
}

func ListInstances(app, sha, env string) (instances []string, err error) {
	instances, _, err = Zk.VisibleChildren(helper.GetBaseInstancePath(app, sha, env))
	if err != nil {
		log.Printf("Error getting list of instances. Error: %s.", err.Error())
	}
	if instances == nil {
		log.Println("No instances found. Returning empty list.")
		instances = []string{}
	}
	return
}

func ListAllInstances() (instances []string, err error) {
	instances, _, err = Zk.VisibleChildren(helper.GetBaseInstanceDataPath())
	if err != nil {
		log.Printf("Error getting list of all instances. Error: %s.", err.Error())
	}
	if instances == nil {
		log.Println("No instances found. Returning empty list.")
		instances = []string{}
	}
	return
}
