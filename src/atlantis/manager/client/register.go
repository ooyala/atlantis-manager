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

type RegisterRouterCommand struct {
	Wait     bool   `long:"wait" description:"wait until done before exiting"`
	Internal bool   `long:"internal" description:"true to list internal routers"`
	Zone     string `short:"z" long:"zone" description:"the zone to register in"`
	Host     string `short:"H" long:"host" description:"the host to register"`
	IP       string `short:"i" long:"ip" description:"the ip to register"`
	Arg      ManagerRegisterRouterArg
	Reply    ManagerRegisterRouterReply
}

type UnregisterRouterCommand struct {
	Wait     bool   `long:"wait" description:"wait until done before exiting"`
	Internal bool   `long:"internal" description:"true to list internal routers"`
	Zone     string `short:"z" long:"zone" description:"the zone to register in"`
	Host     string `short:"H" long:"host" description:"the host to unregister"`
	Arg      ManagerRegisterRouterArg
	Reply    ManagerRegisterRouterReply
}

func OutputRegisterRouterReply(reply *ManagerRegisterRouterReply) error {
	Log("-> Status: %s", reply.Status)
	if reply.Router != nil {
		Log("-> Router:")
		Log("->   Internal : %t", reply.Router.Internal)
		Log("->   Zone     : %s", reply.Router.Zone)
		Log("->   Host     : %s", reply.Router.Host)
		Log("->   CName    : %s", reply.Router.CName)
	}
	return Output(map[string]interface{}{"status": reply.Status, "router": reply.Router}, nil, nil)
}

type RegisterRouterResultCommand struct {
	ID string `short:"i" long:"id" description:"the task ID to fetch the result for"`
}

func (c *RegisterRouterResultCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	args = ExtractArgs([]*string{&c.ID}, args)
	Log("RegisterRouter Result...")
	arg := c.ID
	var reply ManagerRegisterRouterReply
	if err := rpcClient.Call("RegisterRouterResult", arg, &reply); err != nil {
		return OutputError(err)
	}
	return OutputRegisterRouterReply(&reply)
}

type UnregisterRouterResultCommand struct {
	ID string `short:"i" long:"id" description:"the task ID to fetch the result for"`
}

func (c *UnregisterRouterResultCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	args = ExtractArgs([]*string{&c.ID}, args)
	Log("UnregisterRouter Result...")
	arg := c.ID
	var reply ManagerRegisterRouterReply
	if err := rpcClient.Call("UnregisterRouterResult", arg, &reply); err != nil {
		return OutputError(err)
	}
	return OutputRegisterRouterReply(&reply)
}

type GetRouterCommand struct {
	Internal bool   `long:"internal" description:"true to list internal routers"`
	Zone     string `short:"z" long:"zone" description:"the zone to register in"`
	Host     string `short:"H" long:"host" description:"the host to get"`
	Arg      ManagerGetRouterArg
	Reply    ManagerGetRouterReply
}

type ListRoutersCommand struct {
	Internal bool `long:"internal" description:"true to list internal routers"`
	Arg      ManagerListRoutersArg
	Reply    ManagerListRoutersReply
}

type RegisterAppCommand struct {
	Name        string `short:"a" long:"app" description:"the app to register"`
	NonAtlantis bool   `short:"n" long:"non-atlantis" description:"true if this is a non-atlantis app"`
	Internal    bool   `short:"i" long:"internal" description:"true if this is an internal app"`
	Repo        string `short:"g" long:"git" description:"the app's git repository"`
	Root        string `short:"r" long:"root" description:"the app's root within the repo"`
	Email       string `short:"e" long:"email" description"the email of the app's owner"`
	Arg         ManagerRegisterAppArg
	Reply       ManagerRegisterAppReply
}

type UpdateAppCommand struct {
	App         string `short:"a" long:"app" description:"the app to update"`
	NonAtlantis bool   `short:"n" long:"non-atlantis" description:"true if this is a non-atlantis app"`
	Internal    bool   `short:"i" long:"internal" description:"true if this is an internal app"`
	Repo        string `short:"g" long:"git" description:"the app's git repository (or host:port for non-atlantis apps)"`
	Root        string `short:"r" long:"root" description:"the app's root within the repo"`
	Email       string `short:"e" long:"email" description"the email of the app's owner"`
}

func (c *UpdateAppCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Update App...")
	args = ExtractArgs([]*string{&c.App, &c.Repo, &c.Root}, args)
	arg := ManagerRegisterAppArg{
		ManagerAuthArg: dummyAuthArg,
		NonAtlantis:    c.NonAtlantis,
		Internal:       c.Internal,
		Name:           c.App,
		Repo:           c.Repo,
		Root:           c.Root,
		Email:          c.Email,
	}
	var reply ManagerRegisterAppReply
	err = rpcClient.CallAuthed("UpdateApp", &arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	return Output(map[string]interface{}{"status": reply.Status}, nil, nil)
}

type UnregisterAppCommand struct {
	Name  string `short:"a" long:"app" description:"the app to unregister"`
	Arg   ManagerRegisterAppArg
	Reply ManagerRegisterAppReply
}

type GetAppCommand struct {
	Name  string `short:"a" long:"app" description:"the app to get"`
	Arg   ManagerGetAppArg
	Reply ManagerGetAppReply
}

type ListRegisteredAppsCommand struct {
	Arg   ManagerListRegisteredAppsArg
	Reply ManagerListRegisteredAppsReply
}

type ListAuthorizedRegisteredAppsCommand struct {
	Arg   ManagerListRegisteredAppsArg
	Reply ManagerListRegisteredAppsReply
}

type HealthCommand struct {
}

func (c *HealthCommand) Execute(args []string) error {
	InitNoLogin()
	Log("Manager Health Check...")
	arg := ManagerHealthCheckArg{}
	var reply ManagerHealthCheckReply
	err := rpcClient.CallWithTimeout("HealthCheck", arg, &reply, 5)
	if err != nil {
		return OutputError(err)
	}
	Log("-> region: %s", reply.Region)
	Log("-> status: %s", reply.Status)
	return Output(map[string]interface{}{"status": reply.Status, "region": reply.Region}, reply.Region, nil)
}

type RegisterManagerCommand struct {
	Wait          bool   `long:"wait" description:"wait until done before exiting"`
	Host          string `short:"H" long:"host" description:"the host to register"`
	Region        string `short:"r" long:"region" description:"the region to register"`
	ManagerCName  string `long:"manager-cname" description:"the manager's cname if it already has one"`
	RegistryCName string `long:"registry-cname" description:"the registry's cname if it already has one"`
	Arg           ManagerRegisterManagerArg
	Reply         ManagerRegisterManagerReply
}

type UnregisterManagerCommand struct {
	Wait   bool   `long:"wait" description:"wait until done before exiting"`
	Host   string `short:"H" long:"host" description:"the host to register"`
	Region string `short:"r" long:"region" description:"the region to unregister"`
	Arg    ManagerRegisterManagerArg
	Reply  ManagerRegisterManagerReply
}

func OutputRegisterManagerReply(reply *ManagerRegisterManagerReply) error {
	Log("-> Status: %s", reply.Status)
	if reply.Manager != nil {
		Log("-> Manager:")
		Log("->   Region:         %s", reply.Manager.Region)
		Log("->   Host:           %s", reply.Manager.Host)
		Log("->   Registry CName: %s", reply.Manager.RegistryCName)
		Log("->   Manager CName:  %s", reply.Manager.ManagerCName)
	}
	return Output(map[string]interface{}{"status": reply.Status, "manager": reply.Manager}, nil, nil)
}

type RegisterManagerResultCommand struct {
	ID string `short:"i" long:"id" description:"the task ID to fetch the result for"`
}

func (c *RegisterManagerResultCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	args = ExtractArgs([]*string{&c.ID}, args)
	Log("RegisterManager Result...")
	arg := c.ID
	var reply ManagerRegisterManagerReply
	if err := rpcClient.Call("RegisterManagerResult", arg, &reply); err != nil {
		return OutputError(err)
	}
	return OutputRegisterManagerReply(&reply)
}

type UnregisterManagerResultCommand struct {
	ID string `short:"i" long:"id" description:"the task ID to fetch the result for"`
}

func (c *UnregisterManagerResultCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	args = ExtractArgs([]*string{&c.ID}, args)
	Log("UnregisterManager Result...")
	arg := c.ID
	var reply ManagerRegisterManagerReply
	if err := rpcClient.Call("UnregisterManagerResult", arg, &reply); err != nil {
		return OutputError(err)
	}
	return OutputRegisterManagerReply(&reply)
}

type ListManagersCommand struct {
	Arg   ManagerListManagersArg
	Reply ManagerListManagersReply
}

func OutputGetManagerReply(reply *ManagerGetManagerReply) error {
	Log("-> Status: %s", reply.Status)
	if reply.Manager != nil {
		Log("-> Manager:")
		Log("->   Region:         %s", reply.Manager.Region)
		Log("->   Host:           %s", reply.Manager.Host)
		Log("->   Registry CName: %s", reply.Manager.RegistryCName)
		Log("->   Manager CName:  %s", reply.Manager.ManagerCName)
		Log("->   Roles:")
		for role, typeMap := range reply.Manager.Roles {
			Log("->     %s", role)
			for typeName, val := range typeMap {
				Log("->       %s : %t", typeName, val)
			}
		}
	}
	return Output(map[string]interface{}{"status": reply.Status, "manager": reply.Manager}, nil, nil)
}

type GetManagerCommand struct {
	Region string `short:"r" long:"region" description:"the region to get"`
	Host   string `short:"H" long:"host" description:"the host to get"`
	Arg    ManagerGetManagerArg
	Reply  ManagerGetManagerReply
}

type GetSelfCommand struct {
	Arg   ManagerGetSelfArg
	Reply ManagerGetManagerReply
}

type RegisterSupervisorCommand struct {
	Wait  bool   `long:"wait" description:"wait until done before exiting"`
	Host  string `short:"H" long:"host" description:"the supervisor host to register"`
	Arg   ManagerRegisterSupervisorArg
	Reply ManagerRegisterSupervisorReply
}

type UnregisterSupervisorCommand struct {
	Wait  bool   `long:"wait" description:"wait until done before exiting"`
	Host  string `short:"H" long:"host" description:"the supervisor host to register"`
	Arg   ManagerRegisterSupervisorArg
	Reply ManagerRegisterSupervisorReply
}

func OutputRegisterSupervisorReply(reply *ManagerRegisterSupervisorReply) error {
	Log("-> Status: %s", reply.Status)
	return Output(map[string]interface{}{"status": reply.Status}, nil, nil)
}

type RegisterSupervisorResultCommand struct {
	ID string `short:"i" long:"id" description:"the task ID to fetch the result for"`
}

func (c *RegisterSupervisorResultCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	args = ExtractArgs([]*string{&c.ID}, args)
	Log("RegisterSupervisor Result...")
	arg := c.ID
	var reply ManagerRegisterSupervisorReply
	if err := rpcClient.Call("RegisterSupervisorResult", arg, &reply); err != nil {
		return OutputError(err)
	}
	return OutputRegisterSupervisorReply(&reply)
}

type UnregisterSupervisorResultCommand struct {
	ID string `short:"i" long:"id" description:"the task ID to fetch the result for"`
}

func (c *UnregisterSupervisorResultCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	args = ExtractArgs([]*string{&c.ID}, args)
	Log("UnregisterSupervisor Result...")
	arg := c.ID
	var reply ManagerRegisterSupervisorReply
	if err := rpcClient.Call("UnregisterSupervisorResult", arg, &reply); err != nil {
		return OutputError(err)
	}
	return OutputRegisterSupervisorReply(&reply)
}

type ListSupervisorsCommand struct {
	Arg   ManagerListSupervisorsArg
	Reply ManagerListSupervisorsReply
}
