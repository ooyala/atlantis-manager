/* Copyright 2014 Ooyala, Inc. All rights reserved.
 *
 * This file is licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
 * except in compliance with the License. You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is
 * distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

package datamodel

import (
	"atlantis/manager/helper"
	"atlantis/manager/rpc/types"
	"log"
)

type ZkManager types.Manager

func Manager(region, value string) *ZkManager {
	return &ZkManager{Region: region, Host: value, Roles: map[string]map[string]bool{}}
}

func (m *ZkManager) Save() error {
	if m.Roles == nil {
		m.Roles = map[string]map[string]bool{}
	}
	return setJson(m.path(), m)
}

func (m *ZkManager) AddRole(name, roleType string) error {
	if m.Roles == nil {
		m.Roles = map[string]map[string]bool{}
	}
	if m.Roles[name] == nil {
		m.Roles[name] = map[string]bool{}
	}
	if roleType != "" {
		m.Roles[name][roleType] = true
	}
	return m.Save()
}

func (m *ZkManager) HasRole(name, roleType string) bool {
	if m.Roles[name] == nil {
		return false
	}
	return m.Roles[name][roleType]
}

func (m *ZkManager) RemoveRole(name, roleType string) error {
	if m.Roles == nil {
		m.Roles = map[string]map[string]bool{}
		return m.Save()
	}
	if m.Roles[name] == nil {
		return nil
	}
	if roleType != "" {
		delete(m.Roles[name], roleType)
	} else {
		delete(m.Roles, name)
	}
	return m.Save()
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
	if err == nil && zm.Roles == nil {
		err = zm.Save()
	}
	return
}

func ManagerHasRole(region, value, name, roleType string) (bool, error) {
	zm, err := GetManager(region, value)
	if err != nil {
		return false, err
	}
	return zm.HasRole(name, roleType), nil
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
