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
	"atlantis/manager/crypto"
	crypt "atlantis/crypto"
	"atlantis/manager/helper"
	"atlantis/manager/rpc/types"
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

type App struct {
	Name string		`db:"name"`
	NonAtlantis bool	`db:"nonatlantis"`
	Internal bool		`db:"internal"`
	Email string		`db:"email"`
	Repo string		`db:"repo"`
	Root string		`db:"root"`
}

type DepsData struct {
	
	ID int64	`db:"id"`
	AppDepId int64	`db:"appdepid"`
	Enviroment string `db:"env"`
	SecGroup string	`db:"secgroup"`
	DataMap	string	`db:"datamap"`
	EncryptedData string	`db:"encdata"`
}

type AppDeps struct {
	ID int64 `db:"id"`
	App string `db:"app"`    
	Depender string `db:"depender"`
}

type ZkApp types.App

func GetApp(name string) (za *ZkApp, err error) {
	za = &ZkApp{}
	err = getJson(helper.GetBaseAppPath(name), za)
	if za.DependerEnvData == nil {
		za.DependerEnvData = map[string]*types.DependerEnvData{}
		za.Save()
	}
	if za.DependerAppData == nil {
		za.DependerAppData = map[string]*types.DependerAppData{}
		za.Save()
	}
	
	////////////////////////// SQL ////////////////////////////////
	obj, err := DbMap.Get(App{}, name)
	if err != nil {
		fmt.Printf("\n%v\n", err)	
	} else {
		if obj == nil {
			fmt.Println("App doesn't exists \n")
			err = errors.New("No app with name: " + name)
			za = nil	
		} else { 
			app := obj.(*App)
			err = nil
			za = &ZkApp{ app.NonAtlantis, app.Internal, app.Name, app.Email, app.Repo, app.Root, 
				map[string]*types.DependerEnvData{}, map[string]*types.DependerAppData{}}
		}
	}
	///////////////////////////////////////////////////////////////

	return
}

func CreateOrUpdateApp(nonAtlantis, internal bool, name, repo, root, email string) (*ZkApp, error) {
	za, err := GetApp(name)
	if err != nil {
		za = &ZkApp{
			NonAtlantis:     nonAtlantis,
			Internal:        internal,
			Name:            name,
			Repo:            repo,
			Root:            root,
			Email:           email,
			DependerEnvData: map[string]*types.DependerEnvData{},
			DependerAppData: map[string]*types.DependerAppData{},
		}
		if err := za.Save(); err != nil {
			return za, err
		}
	} else {
		za.Name = name
		za.Repo = repo
		za.Root = root
		za.Email = email
		if za.Internal != internal {
			return za, errors.New("apps may not change from internal to external (and visa versa). please unregister and reregister.")
		}
		if za.NonAtlantis != nonAtlantis {
			return za, errors.New("apps may not change from non-atlantis to atlantis (and visa versa). please unregister and reregister.")
		}
		if err := za.Save(); err != nil {
			return za, err
		}
	}
	return za, nil
}

func (za *ZkApp) Delete() error {
	// this just deletes the registration. no need to clean up already deployed instances
	if !za.NonAtlantis {
		if err := ReclaimRouterPortsForApp(za.Internal, za.Name); err != nil {
			return err
		}
	}
	
	////////////////////// SQL /////////////////////////////////
	app := App{}
	app.Name = za.Name
	DbMap.Delete(app)
	//if not
	//_, err := DbMap.Exec("delete from apps where name=?", app.Name)
	//////////////////////////////////////////////////////////
	
	
	return Zk.RecursiveDelete(za.path())
}

func (za *ZkApp) path() string {
	return helper.GetBaseAppPath(za.Name)
}

func (za *ZkApp) Save() error {
	if err := setJson(za.path(), za); err != nil {
		return err
	}


	///////////////////////// SQL ///////////////////////////////////////
	fmt.Println("YAY SQL SAVE APP STUFF")
	app := App{za.Name, za.NonAtlantis, za.Internal, za.Email, za.Repo, za.Root}
	obj, err := DbMap.Get(App{}, za.Name)	
	if err != nil {
		fmt.Printf("\n Failed trying to check if app exists \n")
	}
	if za.Name == "" {
		return errors.New("Trying to save app with no name")
	}
	//app doesnt exist, insert
	if obj == nil {

		fmt.Printf("\n App %s does not exists, insert \n", za.Name)
		err = DbMap.Insert(&app)
		if err != nil {
			fmt.Printf("\n %v \n", err)
		}
	} else {
		fmt.Printf("\n App %s DOES exist, update \n ", za.Name)
		//update
		_, err = DbMap.Update(&app)
		if err != nil {
			fmt.Printf("\n %v \n", err)
		}
	}	
	/////////////////////////////////////////////////////////////////////

	return nil
}

func (za *ZkApp) AddDependerEnvData(data *types.DependerEnvData) error {
	if za.DependerEnvData == nil {
		za.DependerEnvData = map[string]*types.DependerEnvData{}
	}
	if _, err := GetEnv(data.Name); err != nil {
		return err
	}
	crypto.EncryptDependerEnvData(data)
	za.DependerEnvData[data.Name] = data

	//////////////////////////// SQL ///////////////////////////////////
	secGroupStr, err := json.Marshal(data.SecurityGroup)
	if err != nil {

	}
	mapData, err := json.Marshal(data.DataMap)
	if err != nil {
	}	
	encDataStr := string(crypt.Encrypt(mapData)) 			
		
	var appDep AppDeps
	err = DbMap.SelectOne(&appDep, "select * from appdeps where app=? AND depender=?", za.Name, nil)
	if err != nil {
		fmt.Printf("\n No record for App : %s, creating \n", za.Name)
		appDep = AppDeps{App: za.Name}		
		err = DbMap.Insert(&appDep)
		if err != nil {
			fmt.Printf("\n Could not insert new plain app in dep : %v \n", err)
		}
	}
 
	dep := DepsData{AppDepId: appDep.ID, Enviroment: data.Name, SecGroup: string(secGroupStr), DataMap: string(mapData), 
				EncryptedData: encDataStr}	
	err = DbMap.Insert(&dep)
	if err != nil {
		fmt.Printf("\n Could not insert new env data for %s : %v \n", za.Name, err)	
	}	
	///////////////////////////////////////////////////////////////////

	return za.Save()
}

func (za *ZkApp) RemoveDependerEnvData(env string) error {
	if za.DependerEnvData == nil {
		za.DependerEnvData = map[string]*types.DependerEnvData{}
	}
	delete(za.DependerEnvData, env)

	///////////////////////////// SQL ///////////////////////////////
	var appDepId int
	err := DbMap.SelectOne(&appDepId, "select id from appdeps where app=? and depender=?", za.Name, nil)
	if err != nil {
		fmt.Printf("\n Cannot find dep ID to remove data for env: %s : %v\n", env, err)
	} else {
		_, err := DbMap.Exec("delete from depsdata where appdepid=?", appDepId)
		if err != nil {
			fmt.Printf("\n Error deleting dep data for env %s : %v \n", env, err)
		}
 	}
	////////////////////////////////////////////////////////////////

	return za.Save()
}

func (za *ZkApp) GetDependerEnvData(env string, decrypt bool) *types.DependerEnvData {
	if za.DependerEnvData == nil {
		za.DependerEnvData = map[string]*types.DependerEnvData{}
		za.Save()
	}
	ded := za.DependerEnvData[env]
	if ded == nil {
		return nil
	}
	if decrypt {
		crypto.DecryptDependerEnvData(ded)
	}

	/////////////////////// SQL /////////////////////////////////
	var appDepId int
	err := DbMap.SelectOne(&appDepId, "select id from appdeps where app=? and depender=?", za.Name, nil)
	if err != nil {
		fmt.Printf("\n Cannot find dep ID to get  data for env: %s : %v\n", env, err)
	} else { 
		var dep DepsData 
		err := DbMap.SelectOne(&dep, "select * from deps where appdepid=?", appDepId)
		if err != nil {
		}
		var secG map[string][]uint16
		var dMap map[string]interface{}
		err = json.Unmarshal([]byte(dep.SecGroup), secG)
		if err != nil {

		}
		err = json.Unmarshal([]byte(dep.DataMap), dMap)
		if err != nil{
		}
	}
	//ded := DependerEnvData{dep.Enviroment, secG, dep.EncryptedData, dep.DataMap}
	////////////////////////////////////////////////////////////


	return ded
}

func (za *ZkApp) AddDependerAppData(data *types.DependerAppData) error {
	if za.DependerAppData == nil {
		za.DependerAppData = map[string]*types.DependerAppData{}
	}
	if _, err := GetApp(data.Name); err != nil {
		return err
	}
	for _, ded := range data.DependerEnvData {
		if _, err := GetEnv(ded.Name); err != nil {
			return err
		}
		crypto.EncryptDependerEnvData(ded)
	}
	za.DependerAppData[data.Name] = data
	
	////////////////////////////// SQL //////////////////////////////////////////
	
	var appDep AppDeps
	err := DbMap.SelectOne(&appDep, "select * from appdeps where app=? AND depender=?", za.Name, data.Name)
	if err != nil {
		fmt.Printf("\n No record for App : %s, creating \n", za.Name)
		appDep = AppDeps{App: za.Name, Depender: data.Name}		
		err = DbMap.Insert(&appDep)
		if err != nil {
			fmt.Printf("\n Could not new app dep relation : %v \n", err)
		}
	}
	for _, ded := range data.DependerEnvData {
		secGroupStr, err := json.Marshal(ded.SecurityGroup)
		if err != nil {

		}
		mapData, err := json.Marshal(ded.DataMap)
		if err != nil {
		}	
		encDataStr := string(crypt.Encrypt(mapData)) 			

		dep := DepsData{AppDepId: appDep.ID, Enviroment: ded.Name, 
				 SecGroup: string(secGroupStr), DataMap: string(mapData), EncryptedData: encDataStr}
		err = DbMap.Insert(&dep)
		if err != nil {
			fmt.Printf("\n Error inserting env data for app,depender,env %s,%s,%s : %v \n", za.Name, data.Name, ded.Name, err)
		}
	}
	////////////////////////////////////////////////////////////////////////////

	return za.Save()
}

func (za *ZkApp) RemoveDependerAppData(app string) error {
	if za.DependerAppData == nil {
		za.DependerAppData = map[string]*types.DependerAppData{}
	}
	delete(za.DependerAppData, app)


	///////////////////////////// SQL ///////////////////////////////////////
	var appDepId int
	err := DbMap.SelectOne(&appDepId, "select id from appdeps where app=? and depender=?", za.Name, nil)
	if err != nil {
		fmt.Printf("\n Cannot find dep ID to delete data for app,depender: %s : %v\n", za.Name, app, err)
	} else { 
		_, err = DbMap.Exec("delete from depsdata where appdepid=?", appDepId)
		if err != nil {
			fmt.Printf("\n Trouble deleting depdata for app,dep %s,%s : %v\n", za.Name, app, err)	
		}		
		_, err = DbMap.Exec("delete from appdeps where id=?", appDepId)
		if err != nil{
			fmt.Printf("\n Error deleting relation in appdep table for app,dep %s,%s : %v \n", za.Name, app, err)
		}
	}
	////////////////////////////////////////////////////////////////////////	

	return za.Save()
}

func (za *ZkApp) GetDependerAppData(app string, decrypt bool) *types.DependerAppData {
	if za.DependerAppData == nil {
		za.DependerAppData = map[string]*types.DependerAppData{}
		za.Save()
	}
	dad := za.DependerAppData[app]
	if dad == nil {
		return nil
	}
	if decrypt {
		for _, ded := range dad.DependerEnvData {
			crypto.DecryptDependerEnvData(ded)
		}
	}
	//////////////////////////// SQL ///////////////////////////////////////
	var appDepId int
	err := DbMap.SelectOne(&appDepId, "select id from appdeps where app=? and depender=?", za.Name, app)
	if err != nil {
		fmt.Printf("\n Cannot find dep ID to delete data for app,depender: %s : %v\n", za.Name, app, err)
	} else {
		sqlDad := &types.DependerAppData{}
		sqlDad.Name = app 
		sqlDad.DependerEnvData = make(map[string]*types.DependerEnvData)
		var depDataList []DepsData
		_, err := DbMap.Select(&depDataList, "select * from depsdata where appdepid=?", appDepId)
		if err != nil {
			fmt.Printf("Could not find any dep data for app,dep %s,%s : %v \n", za.Name, app, err)
		}

		for _, dData := range depDataList {

			sqlDad.DependerEnvData[dData.Enviroment] = &types.DependerEnvData{Name: dData.Enviroment} 
			var secg map[string][]uint16
			err = json.Unmarshal([]byte(dData.SecGroup), &secg)
			if err != nil {
			}
			sqlDad.DependerEnvData[dData.Enviroment].SecurityGroup = secg
			if decrypt {
				dMapStr := crypt.Decrypt([]byte(dData.EncryptedData))
				var dMap map[string]interface{}
				err = json.Unmarshal([]byte(dMapStr), &dMap)	
				if err != nil {
				}
				sqlDad.DependerEnvData[dData.Enviroment].DataMap = dMap
			}
			sqlDad.DependerEnvData[dData.Enviroment].EncryptedData = dData.EncryptedData
		}
	} 
	///////////////////////////////////////////////////////////////////////

	return dad
}

func (za *ZkApp) AddDependerEnvDataForDependerApp(app string, data *types.DependerEnvData) error {
	if _, err := GetApp(app); err != nil {
		return err
	}
	if _, err := GetEnv(data.Name); err != nil {
		return err
	}
	dad := za.GetDependerAppData(app, false)
	if dad == nil {
		dad = &types.DependerAppData{Name: app, DependerEnvData: map[string]*types.DependerEnvData{}}
	}
	crypto.EncryptDependerEnvData(data)
	dad.DependerEnvData[data.Name] = data
	za.DependerAppData[app] = dad

	///////////////////////////////////////////////// SQL //////////////////////////////////////////////
	var appDep AppDeps
	err := DbMap.SelectOne(&appDep, "select * from appdeps where app=? AND depender=?", za.Name, app)
	if err != nil {
		fmt.Printf("\n No record for App : %s, creating \n", za.Name)
		appDep = AppDeps{App: za.Name, Depender: app}		
		err = DbMap.Insert(&appDep)
		if err != nil {
			fmt.Printf("\n Could not new app dep relation : %v \n", err)
		}
	}
	secGroupStr, err := json.Marshal(data.SecurityGroup)
	if err != nil {

	}
	mapData, err := json.Marshal(data.DataMap)
	if err != nil {
	}	
	encDataStr := string(crypt.Encrypt(mapData)) 			

	dep := DepsData{AppDepId: appDep.ID, Enviroment: data.Name, 
			 SecGroup: string(secGroupStr), DataMap: string(mapData), EncryptedData: encDataStr}
	err = DbMap.Insert(&dep)
	if err != nil {
		fmt.Printf("\n Error inserting env data for app,depender,env %s,%s,%s : %v \n", za.Name, app, data.Name, err)
	}
	/////////////////////////////////////////////////////////////////////////////////////////////////////

	return za.Save()
}

func (za *ZkApp) RemoveDependerEnvDataForDependerApp(app, env string) error {
	dad := za.GetDependerAppData(app, false)
	if dad == nil {
		return nil
	}
	delete(dad.DependerEnvData, env)
	za.DependerAppData[app] = dad

	//////////////////////////////////////////////// SQL //////////////////////////////////////////
	var appDepId int
	err := DbMap.SelectOne(&appDepId, "select id from appdeps where app=? and depender=?", za.Name, app) 
	if err != nil {
		fmt.Printf("\n Cannot find dep ID to delete data for app,depender: %s : %v\n", za.Name, app, err)
	} else {
		_, err = DbMap.Exec("delete from depsdata where appdepid=? and env=?", appDepId, env)
		if err != nil {
			fmt.Printf("\n Could not delete envdata for app,dep,env %s,%s,%s : %v \n", za.Name, app, env)
		}
	}

	///////////////////////////////////////////////////////////////////////////////////////////// 


	return za.Save()
}

func (za *ZkApp) GetDependerEnvDataForDependerApp(app, env string, decrypt bool) *types.DependerEnvData {
	dad := za.GetDependerAppData(app, false)
	if dad == nil {
		return nil
	}
	ded := dad.DependerEnvData[env]
	if ded == nil {
		return nil
	}
	if decrypt {
		crypto.DecryptDependerEnvData(ded)
	}

	///////////////////////////////////////////////// SQL ///////////////////////////////////////////////
	var appDepId int
	err := DbMap.SelectOne(&appDepId, "select id from appdeps where app=? and depender=?", za.Name, app) 
	if err != nil {
		fmt.Printf("\n Cannot find dep ID to delete data for app,depender: %s : %v\n", za.Name, app, err)
	}
	var depEnvData DepsData
	err = DbMap.SelectOne(&depEnvData, "select * from depsdata where appdepid=? and env=?", appDepId, env)
	if err != nil {
		return nil
	}

	sqlDed := &types.DependerEnvData{Name: depEnvData.Enviroment}
	var secg map[string][]uint16
	err = json.Unmarshal([]byte(depEnvData.SecGroup), &secg)
	if err != nil {
	}
	sqlDed.SecurityGroup = secg
	if decrypt {
		var dMap map[string]interface{}
		dMapStr := crypt.Decrypt([]byte(depEnvData.EncryptedData))
		err = json.Unmarshal(dMapStr, &dMap)	
		if err != nil {
		}
		sqlDed.DataMap = dMap
	}
	sqlDed.EncryptedData = depEnvData.EncryptedData

	////////////////////////////////////////////////////////////////////////////////////////////////////// 

	return ded
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
	
	/////////////////////////// SQL ///////////////////////////////////
	var appSql []App
	_, err = DbMap.Select(&appSql, "select * from apps")
	if err != nil {

	}
	////////////////////////////////////////////////////////////////////
	
	return
}
