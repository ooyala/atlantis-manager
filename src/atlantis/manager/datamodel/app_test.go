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
	"atlantis/crypto"
	"atlantis/manager/helper"
	"atlantis/manager/rpc/types"
	. "launchpad.net/gocheck"
)

func (s *DatamodelSuite) TestApp(c *C) {
	crypto.Init()
	Zk.RecursiveDelete(helper.GetBaseAppPath())
	Zk.RecursiveDelete(helper.GetBaseEnvPath())
	Env("prod").Save()
	Env("staging").Save()
	apps, err := ListRegisteredApps()
	c.Assert(err, Not(IsNil)) // the path doesn't exist. this is an error
	c.Assert(len(apps), Equals, 0)
	app1, err := GetApp(app)
	c.Assert(err, Not(IsNil))
	app1, err = CreateOrUpdateApp(true, true, "app1", "", "", "jigish@ooyala.com")
	c.Assert(err, IsNil)
	c.Assert(app1.NonAtlantis, Equals, true)
	c.Assert(app1.Internal, Equals, true)
	c.Assert(app1.Name, Equals, "app1")
	c.Assert(app1.Repo, Equals, "")
	c.Assert(app1.Root, Equals, "")
	c.Assert(app1.Email, Equals, "jigish@ooyala.com")
	app2, err := CreateOrUpdateApp(false, true, "app2", repo, root, "jigish@ooyala.com")
	c.Assert(err, IsNil)
	c.Assert(app2.NonAtlantis, Equals, false)
	c.Assert(app2.Internal, Equals, true)
	c.Assert(app2.Name, Equals, "app2")
	c.Assert(app2.Repo, Equals, repo)
	c.Assert(app2.Root, Equals, root)
	c.Assert(app2.Email, Equals, "jigish@ooyala.com")

	// verify ListRegisteredApps
	apps, err = ListRegisteredApps()
	c.Assert(err, IsNil)
	c.Assert(len(apps), Equals, 2)

	// attempt to set env/app data
	c.Assert(app1.AddDependerEnvData(&types.DependerEnvData{
		Name:          "prod",
		SecurityGroup: map[string][]uint16{"1.1.1.1": []uint16{1111}, "1.1.1.2": []uint16{1111}},
		DataMap: map[string]interface{}{
			"dep1": "prodvalue1",
		},
	}), IsNil)
	c.Assert(app1.AddDependerAppData(&types.DependerAppData{
		Name: "app2",
		DependerEnvData: map[string]*types.DependerEnvData{
			"prod": &types.DependerEnvData{
				Name: "prod",
				DataMap: map[string]interface{}{
					"dep2": "prodvalue2",
				},
			},
		},
	}), IsNil)
	c.Assert(app1.AddDependerEnvDataForDependerApp("app2", &types.DependerEnvData{
		Name:          "staging",
		SecurityGroup: map[string][]uint16{"1.1.2.1": []uint16{1121}, "1.1.2.2": []uint16{1121}},
		DataMap: map[string]interface{}{
			"dep1": "stagingvalue1",
		},
	}), IsNil)

	// ensure data is encrypted
	app1, err = GetApp("app1")
	c.Assert(err, IsNil)
	c.Assert(app1.NonAtlantis, Equals, true)
	c.Assert(app1.Internal, Equals, true)
	c.Assert(app1.Name, Equals, "app1")
	c.Assert(app1.Repo, Equals, "")
	c.Assert(app1.Root, Equals, "")
	c.Assert(app1.Email, Equals, "jigish@ooyala.com")
	c.Assert(app1.GetDependerEnvData("somethingthatdoesntexist", false), IsNil)
	prodEnvData := app1.GetDependerEnvData("prod", false)
	c.Assert(prodEnvData, Not(IsNil))
	c.Assert(prodEnvData.EncryptedData, Not(Equals), "")
	c.Assert(prodEnvData.DataMap, IsNil)
	appData := app1.GetDependerAppData("app2", false)
	c.Assert(appData, Not(IsNil))
	for _, envData := range appData.DependerEnvData {
		c.Assert(envData.EncryptedData, Not(Equals), "")
		c.Assert(envData.DataMap, IsNil)
	}
	prodEnvData = app1.GetDependerEnvDataForDependerApp("app2", "staging", false)
	c.Assert(prodEnvData, Not(IsNil))
	c.Assert(prodEnvData.EncryptedData, Not(Equals), "")
	c.Assert(prodEnvData.DataMap, IsNil)

	// check that data is decrypted
	app1, err = GetApp("app1")
	c.Assert(err, IsNil)
	c.Assert(app1.NonAtlantis, Equals, true)
	c.Assert(app1.Internal, Equals, true)
	c.Assert(app1.Name, Equals, "app1")
	c.Assert(app1.Repo, Equals, "")
	c.Assert(app1.Root, Equals, "")
	c.Assert(app1.Email, Equals, "jigish@ooyala.com")
	c.Assert(app1.GetDependerEnvData("somethingthatdoesntexist", true), IsNil)
	prodEnvData = app1.GetDependerEnvData("prod", true)
	c.Assert(prodEnvData, Not(IsNil))
	c.Assert(prodEnvData.DataMap, Not(IsNil))
	appData = app1.GetDependerAppData("app2", true)
	c.Assert(appData, Not(IsNil))
	for _, envData := range appData.DependerEnvData {
		c.Assert(envData.DataMap, Not(IsNil))
	}
	prodEnvData = app1.GetDependerEnvDataForDependerApp("app2", "staging", true)
	c.Assert(prodEnvData, Not(IsNil))
	c.Assert(prodEnvData.DataMap, Not(IsNil))
}
