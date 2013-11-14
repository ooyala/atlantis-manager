package datamodel

import (
	"atlantis/manager/helper"
	"log"
)

type ZkApp struct {
	Name string
	Repo string
	Root string
}

func GetApp(name string) (za *ZkApp, err error) {
	za = &ZkApp{}
	err = getJson(helper.GetBaseAppPath(name), za)
	return
}

func CreateOrUpdateApp(name, repo, root string) (*ZkApp, error) {
	za := &ZkApp{Name: name, Repo: repo, Root: root}
	if _, err := Zk.Touch(za.path()); err != nil {
		Zk.RecursiveDelete(za.path())
		return za, err
	}
	if err := setJson(za.path(), za); err != nil {
		// clean up if error
		Zk.RecursiveDelete(za.path())
		return za, err
	}
	return za, nil
}

func (za *ZkApp) Delete() error {
	// this just deletes the registration. no need to clean up already deployed instances
	return Zk.RecursiveDelete(za.path())
}

func (za *ZkApp) path() string {
	return helper.GetBaseAppPath(za.Name)
}

func ListRegisteredApps() (apps []string, err error) {
	apps, _, err = Zk.VisibleChildren(helper.GetBaseAppPath())
	if err != nil {
		log.Printf("Error getting list of registered apps. Error: %s.", err.Error())
	}
	if apps == nil {
		log.Println("No registered apps found. Returning empty list.")
		apps = []string{}
	}
	return
}
