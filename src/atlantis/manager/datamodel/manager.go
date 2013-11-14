package datamodel

import (
	"atlantis/manager/helper"
	"log"
)

type ZkManager struct {
	region string
	host   string
}

func Manager(region, host string) *ZkManager {
	return &ZkManager{region, host}
}

func (o *ZkManager) Touch() error {
	_, err := Zk.Touch(o.path())
	return err
}

// Delete the manager node and all children (don't realy need DelDir here but there isn't much overhead)
func (o *ZkManager) Delete() error {
	if err := Zk.RecursiveDelete(o.path()); err != nil {
		return err
	}
	managers, err := ListManagersInRegion(o.region)
	if err == nil && managers != nil && len(managers) == 0 {
		Zk.RecursiveDelete(helper.GetBaseManagerPath(o.region))
	} else if err != nil {
		log.Printf("Warning: clean up fail during managers delete: %s", err)
		// this is extra, no need to return the error if we couldn't get them
	}
	return nil
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

func (o *ZkManager) path() string {
	return helper.GetBaseManagerPath(o.region, o.host)
}
