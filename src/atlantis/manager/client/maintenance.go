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

type ContainerMaintenanceCommand struct {
	Container   string `short:"c" long:"container" description:"the container to set maintenance for"`
	Maintenance bool   `short:"m" long:"maintenance" description:"true to set maintenance mode"`
}

func (c *ContainerMaintenanceCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	args = ExtractArgs([]*string{&c.Container}, args)
	Log("ContainerMaintenance ...")
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	arg := ManagerContainerMaintenanceArg{ManagerAuthArg: ManagerAuthArg{user, "", secret},
		ContainerID: c.Container, Maintenance: c.Maintenance}
	var reply ManagerContainerMaintenanceReply
	err = rpcClient.Call("ContainerMaintenance", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("ContainerMaintenance %s.", reply.Status)
	return Output(map[string]interface{}{"status": reply.Status}, nil, nil)
}

type IdleCommand struct {
}

func (c *IdleCommand) Execute(args []string) error {
	InitNoLogin()
	Log("Idle ...")
	arg := ManagerIdleArg{}
	var reply ManagerIdleReply
	err := rpcClient.Call("Idle", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("Idle %t.", reply.Idle)
	return Output(map[string]interface{}{"status": reply.Status, "idle": reply.Idle}, reply.Idle, nil)
}
