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
	"errors"
	"log"
)

type ZkTeamapps struct {
	Team string
	Apps []string
}

func GetTeamapps(team string) (*ZkTeamapps) {
	e := &ZkTeamapps{team, []string{}}
	err := e.Get()
	if err != nil {
		log.Println("Warning:Cannot find team apps for " + team)
	}
	return e
}

func Teamapps(name string, apps []string) *ZkTeamapps {
	return &ZkTeamapps{name, apps}
}

func (e *ZkTeamapps) Save() error {
	return setJson(e.path(), e)
}

/*
func (e *ZkTeamapps) Delete() error {
	return Zk.RecursiveDelete(e.path())
}*/

func (e *ZkTeamapps) AddApp(app string) error {
	for _, elements := range e.Apps {
		if elements == app {
			return errors.New("Unalbe to add app: app " + app + " already in list")
		}
	}
	e.Apps = append(e.Apps, app)
	return setJson(e.path(), e)
}

func (e *ZkTeamapps) DeleteApp(app string) error {
	for index, elements := range e.Apps {
		if elements == app {
			e.Apps = remove(e.Apps, index)
			
			//remove team from zk if it own no apps
			if len(e.Apps) == 0 {
				return Zk.RecursiveDelete(e.path())
			}
			return setJson(e.path(), e)
		}
	}
	return errors.New("unable to delete: app " + app + " not in list")
}


//helper function remove an element from array 
func remove(slice []string, index int) []string {
    return append(slice[:index], slice[index+1:]...)
}

func (e *ZkTeamapps) Get() error {
	return getJson(e.path(), e)
}

func (e *ZkTeamapps) path() string {
	return helper.GetBaseTeamappsPath(e.Team)
}


