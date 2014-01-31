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
	. "atlantis/manager/constant"
)

type VersionCommand struct {
}

func (c *VersionCommand) Execute(args []string) error {
	InitNoLogin()
	Log("Manager Version Check...")
	arg := VersionArg{}
	var reply VersionReply
	defer func() {
		if err := recover(); err != nil {
			reply.RPCVersion = "unknown"
			reply.APIVersion = "unknown"
		}
		Log("-> client rpc: %s", ManagerRPCVersion)
		Log("-> server rpc: %s", reply.RPCVersion)
		Log("-> server api: %s", reply.APIVersion)
		Output(map[string]interface{}{"client": map[string]string{"rpc": ManagerRPCVersion},
			"server": map[string]string{"rpc": reply.RPCVersion, "api": reply.APIVersion}},
			map[string]string{"client rpc": ManagerRPCVersion, "server rpc": reply.RPCVersion,
				"server api": reply.APIVersion}, nil)
	}()
	err := rpcClient.Call("Version", arg, &reply)
	return Output(nil, nil, err)
}
