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
	"fmt"
	"atlantis/manager/helper"
	_ "github.com/go-sql-driver/mysql"
	"github.com/coopernurse/gorp"
)

type Enviroment struct {
	Name string 	`db:"name"`
}

type ZkEnv struct {
	Name string
}

func GetEnv(name string) (*ZkEnv, error) {
	e := &ZkEnv{name}
	err := e.Get()

	/////////////////////// SQL /////////////////
	obj, err := DbMap.Get(Enviroment{}, name)
	if err != nil {

	}
	if obj == nil {

	}
	env := obj.(*Enviroment)	
	////////////////////////////////////////////
	return e, err
}

func Env(name string) *ZkEnv {
	return &ZkEnv{name}
}

func (e *ZkEnv) Save() error {

	///////////////// SQL /////////////////////////	
	env := Enviroment{e.Name}
	DbMap.Insert(env)
	///////////////////////////////////////////////

	
	return setJson(e.path(), e)
}

func (e *ZkEnv) Delete() error {
	if err := ReclaimRouterPortsForEnv(true, e.Name); err != nil {
		return err
	}
	if err := ReclaimRouterPortsForEnv(false, e.Name); err != nil {
		return err
	}

	/////////////////// SQL //////////////////////
	env := Enviroment{e.Name}
	DbMap.Delete(env)
	//if not
	//DbMap.Exec("delete from enviroment where name=?", env.Name)
	/////////////////////////////////////////////
	
	return Zk.RecursiveDelete(e.path())
}

func (e *ZkEnv) Get() error {
	return getJson(e.path(), e)
}

func (e *ZkEnv) path() string {
	return helper.GetBaseEnvPath(e.Name)
}

func ListEnvs() (envs []string, err error) {
	envs, _, err = Zk.Children(helper.GetBaseEnvPath())
	if envs == nil {
		return []string{}, err
	}

	///////////////// SQL /////////////////////
	var envs []Enviroment
	_, err := DbMap.Select(&envs, "select * from enviroment")
	if err != nil {
		//
	}
	//////////////////////////////////////////
	
	return
}
