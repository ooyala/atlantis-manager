package datamodel

import (
	"atlantis/manager/helper"
	"atlantis/manager/rpc/types"
	"log"
)

type ZkApp types.App

func GetApp(name string) (za *ZkApp, err error) {
	za = &ZkApp{}
	err = getJson(helper.GetBaseAppPath(name), za)
	if za.AllowedDependerApps == nil {
		za.AllowedDependerApps = map[string]bool{}
		za.Save()
	}
	return
}

func CreateOrUpdateApp(nonAtlantis bool, name, repo, root, email string) (*ZkApp, error) {
	za, err := GetApp(name)
	if err != nil {
		za = &ZkApp{
			NonAtlantis: nonAtlantis,
			Name: name,
			Repo: repo,
			Root: root,
			Email: email,
			AllowedDependerApps: map[string]bool{},
		}
	} else {
		za.Name = name
		za.Repo = repo
		za.Root = root
		za.Email = email
		za.NonAtlantis = nonAtlantis
	}
	return za, za.Save()
}

func (za *ZkApp) Delete() error {
	// this just deletes the registration. no need to clean up already deployed instances
	return Zk.RecursiveDelete(za.path())
}

func (za *ZkApp) path() string {
	return helper.GetBaseAppPath(za.Name)
}

func (za *ZkApp) Save() error {
	if err := setJson(za.path(), za); err != nil {
		return err
	}
	return nil
}

func (za *ZkApp) AddDepender(app string) error {
	if za.AllowedDependerApps == nil {
		za.AllowedDependerApps = map[string]bool{}
		za.Save()
	}
	if _, err := GetApp(app); err != nil {
		return err
	}
	za.AllowedDependerApps[app] = true
	return za.Save()
}

func (za *ZkApp) RemoveDepender(app string) error {
	if za.AllowedDependerApps == nil {
		za.AllowedDependerApps = map[string]bool{}
		za.Save()
	}
	delete(za.AllowedDependerApps, app)
	return za.Save()
}

func (za *ZkApp) HasDepender(app string) bool {
	if za.AllowedDependerApps == nil {
		za.AllowedDependerApps = map[string]bool{}
		za.Save()
	}
	return za.AllowedDependerApps[app]
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
