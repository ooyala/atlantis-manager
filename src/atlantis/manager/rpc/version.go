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

package rpc

import (
	. "atlantis/common"
	. "atlantis/manager/constant"
)

type VersionExecutor struct {
	arg   VersionArg
	reply *VersionReply
}

func (e *VersionExecutor) Request() interface{} {
	return e.arg
}

func (e *VersionExecutor) Result() interface{} {
	return e.reply
}

func (e *VersionExecutor) Description() string {
	return "Version"
}

func (e *VersionExecutor) Execute(t *Task) error {
	e.reply.RPCVersion = ManagerRPCVersion
	e.reply.APIVersion = ManagerAPIVersion
	t.Log("-> RPC:%s API:%s", ManagerRPCVersion, ManagerAPIVersion)
	return nil
}

func (e *VersionExecutor) Authorize() error {
	return nil
}

func (e *VersionExecutor) AllowDuringMaintenance() bool {
	return true
}

func (m *ManagerRPC) Version(arg VersionArg, reply *VersionReply) error {
	return NewTask("Version", &VersionExecutor{arg, reply}).Run()
}
