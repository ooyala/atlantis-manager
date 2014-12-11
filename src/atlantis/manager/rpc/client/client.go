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
	atlantis "atlantis/common"
	. "atlantis/manager/constant"
)

type ManagerRPCClient struct {
	atlantis.RPCClient
	User    string
	Secrets map[string]string
}

type AuthedArg interface {
	SetCredentials(string, string)
}

func (r *ManagerRPCClient) CallAuthed(name string, arg AuthedArg, reply interface{}) error {
	arg.SetCredentials(r.User, r.Secrets[r.Opts.RPCHostAndPort()])

	return r.RPCClient.Call(name, arg, reply)
}

func NewManagerRPCClient(hostAndPort string) *atlantis.RPCClient {
	return atlantis.NewRPCClient(hostAndPort, "ManagerRPC", ManagerRPCVersion, true)
}

func NewManagerRPCClientWithConfig(cfg atlantis.RPCServerOpts) *atlantis.RPCClient {
	return atlantis.NewRPCClientWithConfig(cfg, "ManagerRPC", ManagerRPCVersion, true)
}
