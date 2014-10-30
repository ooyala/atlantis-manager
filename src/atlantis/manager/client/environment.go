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

type UpdateDepCommand struct {
	Name  string `short:"n" long:"name" description:"the name of the dependency"`
	Value string `short:"v" long:"value" description:"the value of the dependency"`
	Env   string `short:"e" long:"env" description:"the environment of the dependency"`
}

func (c *UpdateDepCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	Log("Update Dep...")
	arg := ManagerDepArg{dummyAuthArg, c.Env, c.Name, c.Value}
	var reply ManagerDepReply
	if err := rpcClient.CallAuthed("UpdateDep", &arg, &reply); err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	Log("-> value: %s", reply.Value)
	return Output(map[string]interface{}{"status": reply.Status, "value": reply.Value}, reply.Value, nil)
}

type ResolveDepsCommand struct {
	App      string   `short:"a" long:"app" description:"the app the resolve dependencies for"`
	Env      string   `short:"e" long:"env" description:"the environment of the dependencies to resolve"`
	DepNames []string `short:"d" long:"dep" description:"the dep names to resolve"`
}

func (c *ResolveDepsCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	Log("Resolve Deps...")
	arg := ManagerResolveDepsArg{ManagerAuthArg: dummyAuthArg, App: c.App, Env: c.Env, DepNames: c.DepNames}
	var reply ManagerResolveDepsReply
	if err := rpcClient.CallAuthed("ResolveDeps", &arg, &reply); err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	Log("-> deps:")
	for zone, deps := range reply.Deps {
		Log("->   %s", zone)
		for name, value := range deps {
			Log("->     %s : %s", name, value)
		}
	}
	return Output(map[string]interface{}{"status": reply.Status, "Deps": reply.Deps}, reply.Deps, nil)
}

type GetDepCommand struct {
	Name string `short:"n" long:"name" description:"the name of the dependency"`
	Env  string `short:"e" long:"env" description:"the environment of the dependency"`
}

func (c *GetDepCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	Log("Get Dep...")
	arg := ManagerDepArg{dummyAuthArg, c.Env, c.Name, ""}
	var reply ManagerDepReply
	if err := rpcClient.CallAuthed("GetDep", &arg, &reply); err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	Log("-> value: %s", reply.Value)
	return Output(map[string]interface{}{"status": reply.Status, "value": reply.Value}, reply.Value, nil)
}

type DeleteDepCommand struct {
	Name string `short:"n" long:"name" description:"the name of the dependency"`
	Env  string `short:"e" long:"env" description:"the environment of the dependency"`
}

func (c *DeleteDepCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	Log("Delete Dep...")
	arg := ManagerDepArg{dummyAuthArg, c.Env, c.Name, ""}
	var reply ManagerDepReply
	if err := rpcClient.CallAuthed("DeleteDep", &arg, &reply); err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	Log("-> value: %s", reply.Value)
	return Output(map[string]interface{}{"status": reply.Status, "value": reply.Value}, reply.Value, nil)
}

type UpdateEnvCommand struct {
	Name string `short:"n" long:"name" description:"the name of the environment"`
}

func (c *UpdateEnvCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	Log("Update Env...")
	arg := ManagerEnvArg{dummyAuthArg, c.Name}
	var reply ManagerEnvReply
	if err := rpcClient.CallAuthed("UpdateEnv", &arg, &reply); err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	return Output(map[string]interface{}{"status": reply.Status}, nil, nil)
}

type DeleteEnvCommand struct {
	Name string `short:"n" long:"name" description:"the name of the environment"`
}

func (c *DeleteEnvCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	Log("Delete Env...")
	arg := ManagerEnvArg{dummyAuthArg, c.Name}
	var reply ManagerEnvReply
	if err := rpcClient.CallAuthed("DeleteEnv", &arg, &reply); err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	return Output(map[string]interface{}{"status": reply.Status}, nil, nil)
}
