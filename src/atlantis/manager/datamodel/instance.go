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
	"errors"
	"log"
)

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
	return false
}

func GetInstance(id string) (zi *ZkInstance, err error) {
	zi = &ZkInstance{}
	err = getJson(helper.GetBaseInstanceDataPath(id), zi)
	if err != nil {
		/* NOTE(edanaher): The default error message here is typically from zookeeper, which is not very user
		 * friendly.  And at Ooyala, we've had a number of users run into this because they're looking at the wrong
		 * region.  So as long as that is the dominant source of this error, it makes sense to suggest this as the
		 * cause, much like compilers will attempt to guess at missing semicolons. */
		err = errors.New(err.Error() + "\nContainer " + id + " not found; is this the right region?")
	}
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

func (zi *ZkInstance) SetManifest(m *types.Manifest) error {
	zi.Manifest = m
	return setJson(zi.dataPath(), zi)
}

func (zi *ZkInstance) path() string {
	return helper.GetBaseInstancePath(zi.App, zi.Sha, zi.Env, zi.ID)
}

func (zi *ZkInstance) dataPath() string {
	return helper.GetBaseInstanceDataPath(zi.ID)
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
