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
	. "atlantis/manager/constant"
	"atlantis/manager/helper"
	routerzk "atlantis/router/zk"
	zookeeper "github.com/jigish/gozk-recipes"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/coopernurse/gorp"
)

var Zk *zookeeper.ZkConn
var DbConn *sql.DB
var DbMap *gorp.DbMap

func CreateRouterPortsPaths() {
	Zk.Touch(helper.GetBaseRouterPortsPath(true))
	Zk.Touch(helper.GetBaseRouterPortsPath(false))
}

func CreateRouterPaths() {
	helper.SetRouterRoot(true)
	for _, path := range routerzk.ZkPaths {
		Zk.Touch(path)
	}
	helper.SetRouterRoot(false)
	for _, path := range routerzk.ZkPaths {
		Zk.Touch(path)
	}
	for _, zone := range AvailableZones {
		Zk.Touch(helper.GetBaseRouterPath(true, zone))
		Zk.Touch(helper.GetBaseRouterPath(false, zone))
	}
}

func CreateLockPaths() {
	Zk.Touch(helper.GetBaseLockPath("deploy"))
	Zk.Touch(helper.GetBaseLockPath("router_ports_internal"))
	Zk.Touch(helper.GetBaseLockPath("router_ports_external"))
}

func CreateAppPath() {
	Zk.Touch(helper.GetBaseAppPath())
}

func CreateInstancePaths() {
	Zk.Touch(helper.GetBaseInstancePath())
	Zk.Touch(helper.GetBaseInstanceDataPath())
}

func CreateSupervisorPath() {
	Zk.Touch(helper.GetBaseSupervisorPath())
}

func CreateManagerPath() {
	Zk.Touch(helper.GetBaseManagerPath())
}

func CreateEnvPath() {
	Zk.Touch(helper.GetBaseEnvPath())
}

func CreatePaths() {
	CreateRouterPortsPaths()
	CreateRouterPaths()
	CreateLockPaths()
	CreateInstancePaths()
	CreateAppPath()
	CreateSupervisorPath()
	CreateManagerPath()
	CreateEnvPath()
}

func Init(zkUri string) {
	Zk = zookeeper.GetPanicingZk(zkUri)
	CreatePaths()

	DbConn, _ := sql.Open("mysql", "root@/manager")
	DbMap = &gorp.DbMap{Db: DbConn, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

	DbMap.AddTableWithName(Instance{}, "instance").SetKeys(false, "name") 
	DbMap.AddTableWithName(Enviroment{}, "enviroment").SetKeys(false, "name")
	DbMap.AddTableWithName(App{}, "apps").SetKeys(false, "name")
	DbMap.AddTableWithName(Manifest{}, "manifest").SetKeys(true, "id")
	DbMap.AddTableWithName(SqlManager{}, "manager").SetKeys(false, "host")
	DbMap.AddTableWithName(Role{}, "roles").SetKeys(true, "id")
	DbMap.AddTableWithName(SupervisorSql{}, "supervisor").SetKeys(false, "name")
	DbMap.AddTableWithName(PortMap{}, "portmap").SetKeys(true, "id")
	DbMap.AddTableWithName(EnvDepData{}, "envdepdata").SetKeys(true, "id")
	DbMap.AddTableWithName(AppDepData{}, "appdepdata").SetKeys(true, "id")
	//should really never create tables unless fresh install
	err := DbMap.CreateTablesIfNotExists()
	if err != nil {
	}	
}
