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
	"atlantis/supervisor/rpc/types"
	"strings"
	"log"
)

type Instance struct {
	ID       string		`db:"name"`
	App      string		`db:"appId"`
	Sha      string		`db:"sha"`
	Env      string		`db:"envId"`
	Host     string		`db:"hostId"`
	Port     uint16		`db:"port"`
	Manifest int64		`db:"manifestId"`	
}

type Manifest struct {
	ID int64		`db:"id"`
	Name string		`db:"name"`
	Description string	`db:"description"`
	Instances int64 	`db:"instances"`
	CPUShares int64		`db:"cpushares"`
	MemoryLimit int64	`db:"memorylimit"`
	AppType	string		`db:"apptype"`
	JavaType string		`db:"javatype"`
	RunCommands string	`db:"runcommands"`
	Dependencies int64	`db:dependencies"`
}

type ZkInstance struct {
	ID       string		
	App      string		
	Sha      string		
	Env      string		
	Host     string	
	Port     uint16		
	Manifest *types.Manifest	
}

func InstanceExists(id string) bool {
	if stat, err := Zk.Exists(helper.GetBaseInstanceDataPath(id)); err == nil && stat != nil {
		return true
	}

	///////////// SQL /////////////////////
	obj, err := DbMap.Get(Instance{}, id)
	if err != nil {
		//fail	
	}
	if obj == nil {
		//not found
	}
	//////////////////////////////////////

	return false
}

func GetInstance(id string) (zi *ZkInstance, err error) {
	zi = &ZkInstance{}
	err = getJson(helper.GetBaseInstanceDataPath(id), zi)

	////////////// SQL /////////////
	obj, err := DbMap.Get(Instance{}, id)
	inst := obj.(*Instance)		
	if inst != nil {
	}
	//TODO: retrive manifest from DB and build instance obj
	//or require whoever uses instance object to manually retrieve
	//the manifest
	////////////////////////////////
	
	return
}

func CreateInstance(app, sha, env, host string) (*ZkInstance, error) {
	id := helper.CreateContainerID(app, sha, env)
	for InstanceExists(id) {
		id = helper.CreateContainerID(app, sha, env)
	}
	zi := &ZkInstance{ID: id, App: app, Sha: sha, Env: env, Host: host, Port: 0}
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

	/////////////// SQL //////////////////////////////

	//TODO: manifest Id FK needs to be set
	//eventually change methods to use Instance instead of ZkInstance also
	inst := Instance{id, app, sha, env, host, 0, 0}		
	DbMap.Insert(inst)


	//////////////////////////////////////////////////
	return zi, nil
}

func (zi *ZkInstance) Delete() (bool, error) { // true if this was the last instance of app+sha+env
	var (
		last bool
		err  error
		err2 error
		list []string
	)
	
	/////////// SQL //////////////////////////
	//TODO check to be sure this works
	inst := Instance{}
	inst.ID = zi.ID
	DbMap.Delete(inst)
	//if not
	//_, err := DbMap.Exec("delete from instance where name=?", inst.ID)	
	/////////////////////////////////////////


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

	//////////////// SQL /////////////////////////
	obj, err := DbMap.Get(Instance{}, zi.ID)
	if err != nil {
		//
	}
	inst := obj.(*Instance)
	inst.Port = port
	_, err = DbMap.Update(inst) 		
	if err != nil {
	}
	/////////////////////////////////////////////

	zi.Port = port
	return setJson(zi.dataPath(), zi)
}

func (zi *ZkInstance) SetManifest(m *types.Manifest) error {
	zi.Manifest = m

	////////////////////////// SQL /////////////////////////////
	//build sql manifest from ZK manifest
	sqlManifest := Manifest{
		Name: m.Name,
		Description: m.Description,
		Instances: int64(m.Instances),
		CPUShares: int64(m.CPUShares),
		MemoryLimit: int64(m.MemoryLimit),
		AppType: m.AppType,
		JavaType: m.JavaType,
		RunCommands: strings.Join(m.RunCommands, ","),
		Dependencies: 0,
	}

	//Insert it to DB
	//****NOTE**** this should populate the ID field of Manifest with
	//the PK of the row in the DB	
	err := DbMap.Insert(&sqlManifest)
	if err != nil {
	}	
	
	//Retrieve instance from DB
	obj, err := DbMap.Get(Instance{}, zi.ID)
	if err != nil {
	}
	inst := obj.(*Instance)
	inst.Manifest = sqlManifest.ID
	_, err = DbMap.Update(&inst)
	if err != nil {
	}
	///////////////////////////////////////////////////////////	

	return setJson(zi.dataPath(), zi)
}

func (zi *ZkInstance) path() string {
	return helper.GetBaseInstancePath(zi.App, zi.Sha, zi.Env, zi.ID)
}

func (zi *ZkInstance) dataPath() string {
	return helper.GetBaseInstanceDataPath(zi.ID)
}

//TODO: Figure out region issue (one db per region or add region column to tables)
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

	///////////////////// SQL //////////////////////////////
        var envList []string
        _, err = DbMap.Select(&envList, "select envId from instance where appId = :appid",
                                map[string]interface{}{
                                        "appid": app,
                                })
	//////////////////////////////////////////////////////
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
	
	///////////////////// SQL //////////////////////////////
	var envList []string
	_, err = DbMap.Select(&envList, "select envId from instance where appId = :appid and sha = :sha", 
				map[string]interface{}{
					"appid": app,
					"sha": sha,
				})
	if err != nil {

	}
	///////////////////////////////////////////////////////

	
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

	//////////////////// SQL //////////////////////////////
	var insts []Instance
	_, err = DbMap.Select(&insts,"select * from instance where appId = :appid and sha = :sha and env = :env",
				map[string]interface{}{
					"appid": app,
					"sha": sha,
					"env": env,
				})
	if err != nil {

	}
	//////////////////////////////////////////////////////

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


	///////////////////// SQL ////////////////////////////
	var insts []Instance
	_, err = DbMap.Select(&insts, "select * from instance")
	if err != nil {

	}
	/////////////////////////////////////////////////////


	return
}
