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

type Manager struct {
	Host string	`db:"host"`
	Region string	`db:"region"`
	CName string	`db:"cname"`
	RecordID string	`db:"recordID"`
	RegistryCName string	`db:"registrycname"`
	RegistryRecordID string	`db:"registryrecID"`

}

type Role struct {
	Id int64	`db:"id"`
	Name string	`db:"name"`
	RoleType string	`db:"roleType"`
	Value bool	`db:"value"`
	Manager string	`db:"manager"`
}

func Manager(region, value string) *ZkManager {
	return &ZkManager{Region: region, Host: value, Roles: map[string]map[string]bool{}}
}

func (m *ZkManager) Save() error {
	if m.Roles == nil {
		m.Roles = map[string]map[string]bool{}
	}

	////////////////// SQL //////////////////////////////
	//TODO: handle roles
	//TODO: might need to check if insert or update
	sqlManager := Manager{m.Host, m.Region, m.ManagerCName, m.ManagerRecordID,
				m.RegistryCName, m.RegistryRecordID, 0}
	err := DbMap.Insert(&sqlManager)
	if err != nil {
	}
	////////////////////////////////////////////////////


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

	//////////////// SQL ////////////////////////
	var role Role
	err := DbMap.SelectOne(&role, "select * from roles where name=? AND roletype=? AND manager=?", name, roleType, m.Host)
	//if not found/not exists create and update
	if err != nil {
		role = role{Name: name, RoleType: roletype, Value: true, Manager: m.host}

		//update the role
		err = DbMap.Insert(&role)
		if err != nil {

		}	
	} 
	////////////////////////////////////////////


	return m.Save()
}

func (m *ZkManager) HasRole(name, roleType string) bool {
	if m.Roles[name] == nil {
		return false
	}

	////////////////////////// SQL ////////////////////
	//TODO: allow for any "RoleType" if not only read/write are supported.
	var role Role 
	err := DbMap.SelectOne(&role, "select * from roles where name=? AND roletype=? AND manager=?", name, roleType, m.Host)
	//no role with that name for this manager
	if err != nil {
		//return role.Value
	} else {
		//not exists
		//return false
	}
	//////////////////////////////////////////////////
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

	//////////////// SQL /////////////////////////
	_, err := DbMap.Exec("delete from roles where name=? AND roletype=? AND manager=?", name, roleType, m.Host)
	if err != nil {
		//do something
	}
	/////////////////////////////////////////////
	
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


	///////////////////////////// SQL //////////////////////////////////
	sqlManager := Manager{}
	sqlManager.Host = m.Host
	DbMap.Delete(sqlManager)
	//if not
	//_, err := DbMap.Exec("delete from manager where host=?", sqlManager.Host)	
	///////////////////////////////////////////////////////////////////

	return nil
}

func GetManager(region, value string) (zm *ZkManager, err error) {
	zm = &ZkManager{}
	err = getJson(helper.GetBaseManagerPath(region, value), zm)
	if err == nil && zm.Roles == nil {
		err = zm.Save()
	}

	//////////////////////// SQL /////////////////////////////
	//TODO: verify the value string is the host name
	obj, err := DbMap.Get(Manager{}, value)
	if err != nil {

	}
	if obj == nil {
		//manager does not exists
	}
	man := obj.(*Manager)
	/////////////////////////////////////////////////////////
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
