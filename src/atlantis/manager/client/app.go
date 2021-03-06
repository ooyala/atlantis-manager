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

package client

import (
	. "atlantis/manager/rpc/types"
)

type RequestAppDependencyCommand struct {
	App        string   `short:"a" long:"app" description:"the app to request a dependency for"`
	Dependency string   `short:"d" long:"dependency" description:"the dependency to request"`
	Envs       []string `short:"e" long:"env" description:"the envs to request the dependency in"`
	Arg        ManagerRequestAppDependencyArg
	Reply      ManagerRequestAppDependencyReply
}

// ----------------------------------------------------------------------------------------------------------
// Depender App Data Methods
// ----------------------------------------------------------------------------------------------------------

type AddDependerAppDataCommand struct {
	App      string `short:"a" long:"app" description:"the app to add a depender for"`
	FromFile string `short:"f" long:"file" description:"the file to pull the data from"`
	Arg      ManagerRequestAppDependencyArg
	Reply    ManagerRequestAppDependencyReply
}

type RemoveDependerAppDataCommand struct {
	App      string `short:"a" long:"app" description:"the app to remove a depender from"`
	Depender string `short:"r" long:"depender" description:"the depender app to remove"`
	Arg      ManagerRemoveDependerAppDataArg
	Reply    ManagerRemoveDependerAppDataReply
}

type GetDependerAppDataCommand struct {
	App        string `short:"a" long:"app" description:"the app to get a depender from"`
	Depender   string `short:"r" long:"depender" description:"the depender app to get"`
	Properties string `message:"Get Depender App" field:"DependerAppData" name:"DependerAppData"`
	Arg        ManagerGetDependerAppDataArg
	Reply      ManagerGetDependerAppDataReply
}

// ----------------------------------------------------------------------------------------------------------
// Depender Env Data Methods
// ----------------------------------------------------------------------------------------------------------

type AddDependerEnvDataCommand struct {
	App        string `short:"a" long:"app" description:"the app to add an env for"`
	FromFile   string `short:"f" long:"file" description:"the file to pull the data from"`
	Properties string `field:"App"`
	Arg        ManagerAddDependerEnvDataArg
	Reply      ManagerAddDependerEnvDataReply
	FileData   DependerEnvData
}

type RemoveDependerEnvDataCommand struct {
	App        string `short:"a" long:"app" description:"the app to remove an env from"`
	Env        string `short:"e" long:"env" description:"the env to remove"`
	Properties string `message:"Remove Depender Env" field:"App" name:"app"`
	Arg        ManagerRemoveDependerEnvDataArg
	Reply      ManagerRemoveDependerEnvDataReply
}

type GetDependerEnvDataCommand struct {
	App        string `short:"a" long:"app" description:"the app to get an env from"`
	Env        string `short:"e" long:"depender" description:"the env to get"`
	Properties string `field:"DependerEnvData"`
	Arg        ManagerGetDependerEnvDataArg
	Reply      ManagerGetDependerEnvDataReply
}

// ----------------------------------------------------------------------------------------------------------
// Depender Env Data For Depender App Methods
// ----------------------------------------------------------------------------------------------------------

type AddDependerEnvDataForDependerAppCommand struct {
	App      string `short:"a" long:"app" description:"the app to add an env for"`
	Depender string `short:"r" long:"depender" description:"the depender to add an env for"`
	FromFile string `short:"f" long:"file" description:"the file to pull the data from"`
	Arg      ManagerAddDependerEnvDataForDependerAppArg
	Reply    ManagerAddDependerEnvDataForDependerAppReply
	FileData DependerEnvData
}

type RemoveDependerEnvDataForDependerAppCommand struct {
	App      string `short:"a" long:"app" description:"the app to remove an env from"`
	Depender string `short:"r" long:"depender" description:"the depender to add an env for"`
	Env      string `short:"e" long:"env" description:"the env to remove"`
	Arg      ManagerRemoveDependerEnvDataForDependerAppArg
	Reply    ManagerRemoveDependerEnvDataForDependerAppReply
}

type GetDependerEnvDataForDependerAppCommand struct {
	App        string `short:"a" long:"app" description:"the app to get an env from"`
	Depender   string `short:"r" long:"depender" description:"the depender to add an env for"`
	Env        string `short:"e" long:"env" description:"the env to get"`
	Properties string `field:"DependerEnvData"`
	Arg        ManagerGetDependerEnvDataForDependerAppArg
	Reply      ManagerGetDependerEnvDataForDependerAppReply
}
