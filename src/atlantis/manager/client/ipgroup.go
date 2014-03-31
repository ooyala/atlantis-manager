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

type UpdateIPGroupCommand struct {
	Name string   `short:"n" long:"name" description:"the name of the ip group"`
	IPs  []string `short:"i" long:"ips" description:"the the ip(s) in the group"`
}

func (c *UpdateIPGroupCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Update IP Group...")
	args = ExtractArgs([]*string{&c.Name}, args)
	user, secret, err := GetSecret()
	if err != nil {
		return OutputError(err)
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerUpdateIPGroupArg{
		ManagerAuthArg: authArg,
		Name:           c.Name,
		IPs:            c.IPs,
	}
	var reply ManagerUpdateIPGroupReply
	err = rpcClient.Call("UpdateIPGroup", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> Status: %s", reply.Status)
	return Output(map[string]interface{}{"status": reply.Status}, nil, nil)
}

type DeleteIPGroupCommand struct {
	Name string `short:"n" long:"name" description:"the name of the ip group"`
}

func (c *DeleteIPGroupCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Delete IP Group...")
	args = ExtractArgs([]*string{&c.Name}, args)
	user, secret, err := GetSecret()
	if err != nil {
		return OutputError(err)
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerDeleteIPGroupArg{
		ManagerAuthArg: authArg,
		Name:           c.Name,
	}
	var reply ManagerDeleteIPGroupReply
	err = rpcClient.Call("DeleteIPGroup", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> Status: %s", reply.Status)
	return Output(map[string]interface{}{"status": reply.Status}, nil, nil)
}

type GetIPGroupCommand struct {
	Name string `short:"n" long:"name" description:"the name of the ip group"`
}

func (c *GetIPGroupCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Get IP Group...")
	args = ExtractArgs([]*string{&c.Name}, args)
	user, secret, err := GetSecret()
	if err != nil {
		return OutputError(err)
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerGetIPGroupArg{
		ManagerAuthArg: authArg,
		Name:           c.Name,
	}
	var reply ManagerGetIPGroupReply
	err = rpcClient.Call("GetIPGroup", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> Status: %s", reply.Status)
	Log("-> IP Group:")
	Log("->   Name: %s", reply.IPGroup.Name)
	Log("->   IPs:")
	for _, ip := range reply.IPGroup.IPs {
		Log("->     %s", ip)
	}
	return Output(map[string]interface{}{"status": reply.Status, "ipGroup": reply.IPGroup}, nil, nil)
}
