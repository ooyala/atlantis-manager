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
	. "atlantis/common"
	. "atlantis/manager/rpc/types"
	"errors"
	"time"
)

const waitPollInterval = 3 * time.Second

type StatusCommand struct {
	ID string `short:"i" long:"id" description:"the task ID to fetch the status for"`
}

func (c *StatusCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	args = ExtractArgs([]*string{&c.ID}, args)
	Log("Task Status...")
	arg := c.ID
	var reply TaskStatus
	err = rpcClient.Call("Status", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> %s", reply.String())
	return Output(reply.Map(), reply.Status, nil)
}

type ResultCommand struct {
	ID string `short:"i" long:"id" description:"the task ID to fetch the result for"`
}

func (c *ResultCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	args = ExtractArgs([]*string{&c.ID}, args)
	arg := c.ID
	var reply TaskStatus
	err = rpcClient.Call("Status", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	switch reply.Name {
	case "Deploy":
		return (&DeployResultCommand{c.ID}).Execute(args)
	case "Teardown":
		return (&TeardownResultCommand{c.ID}).Execute(args)
	case "RegisterManager":
		return (&RegisterManagerResultCommand{c.ID}).Execute(args)
	case "UnregisterManager":
		return (&UnregisterManagerResultCommand{c.ID}).Execute(args)
	case "RegisterRouter":
		return (&RegisterRouterResultCommand{c.ID}).Execute(args)
	case "UnregisterRouter":
		return (&UnregisterRouterResultCommand{c.ID}).Execute(args)
	case "RegisterSupervisor":
		return (&RegisterSupervisorResultCommand{c.ID}).Execute(args)
	case "UnregisterSupervisor":
		return (&UnregisterSupervisorResultCommand{c.ID}).Execute(args)
	default:
		return OutputError(errors.New("Invalid Task Name: " + reply.Name))
	}
}

type WaitCommand struct {
	ID string `short:"i" long:"id" description:"the task ID to wait on"`
}

func (c *WaitCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	args = ExtractArgs([]*string{&c.ID}, args)
	Log("Waiting...")
	arg := c.ID
	var statusReply TaskStatus
	var currentStatus string
	if err := rpcClient.Call("Status", arg, &statusReply); err != nil {
		return OutputError(err)
	}
	for !statusReply.Done {
		time.Sleep(waitPollInterval)
		if currentStatus != statusReply.Status {
			currentStatus = statusReply.Status
			Log(currentStatus)
		}
		if err := rpcClient.Call("Status", c.ID, &statusReply); err != nil {
			return OutputError(err)
		}
	}
	return (&ResultCommand{c.ID}).Execute(args)
}

type ListTaskIDsCommand struct {
}

func (c *ListTaskIDsCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	Log("List Task IDs...")
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	var ids []string
	if err := rpcClient.Call("ListTaskIDs", authArg, &ids); err != nil {
		return OutputError(err)
	}
	return Output(map[string]interface{}{"ids": ids}, ids, nil)
}
