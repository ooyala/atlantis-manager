package datamodel

import (
	"atlantis/manager/helper"
	"atlantis/manager/rpc/types"
	"log"
)

type ZkManager types.Manager

func Manager(region, value string) *ZkManager {
	return &ZkManager{Region: region, Host: value}
}

func (m *ZkManager) Save() error {
	return setJson(m.path(), m)
}

// Delete the manager node and all children (don't realy need DelDir here but there isn't much overhead)
func (m *ZkManager) Delete() error {
	if err := Zk.RecursiveDelete(m.path()); err != nil {
		return err
	}
	managers, err := ListManagersInRegion(m.Region)
	if err == nil && managers != nil && len(managers) == 0 {
		Zk.RecursiveDelete(helper.GetBaseManagerPath(m.Region))
	} else if err != nil {
		log.Printf("Warning: clean up fail during managers delete: %s", err)
		// this is extra, no need to return the error if we couldn't get them
	}
	return nil
}

func GetManager(region, value string) (zm *ZkManager, err error) {
	zm = &ZkManager{}
	err = getJson(helper.GetBaseManagerPath(region, value), zm)
	return
}

func (m *ZkManager) path() string {
	return helper.GetBaseManagerPath(m.Region, m.Host)
}

func ListRegions() (regions []string, err error) {
	basePath := helper.GetBaseManagerPath()
	regions, _, err = Zk.Children(basePath)
	if err != nil {
		log.Printf("Error getting list of regions. Error: %s.", err.Error())
	}
	if regions == nil {
		log.Println("No regions found. Returning empty list.")
		regions = []string{}
	}
	return
}

func ListManagersInRegion(region string) (managers []string, err error) {
	basePath := helper.GetBaseManagerPath(region)
	managers, _, err = Zk.Children(basePath)
	if err != nil {
		log.Printf("Error getting list of managers for region %s. Error: %s.", region, err.Error())
	}
	if managers == nil {
		log.Printf("No managers found in region %s", region)
		managers = []string{}
	}
	return
}

func ListManagers() (managers map[string][]string, err error) {
	managers = map[string][]string{}
	regions, err := ListRegions()
	if err != nil {
		return
	}
	for _, region := range regions {
		managers[region], err = ListManagersInRegion(region)
		if err != nil {
			return
		}
	}
	return
}
