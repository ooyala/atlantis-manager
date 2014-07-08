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
	. "atlantis/manager/rpc/types"
	"atlantis/manager/status"
	"encoding/json"
)

type UsageExecutor struct {
	arg   ManagerUsageArg
	reply *ManagerUsageReply
}

func (e *UsageExecutor) Request() interface{} {
	return e.arg
}

func (e *UsageExecutor) Result() interface{} {
	return e.reply
}

func (e *UsageExecutor) Description() string {
	return "Usage"
}

func (e *UsageExecutor) Execute(t *Task) (err error) {
	e.reply.Usage, err = status.GetUsage()
	b, err := json.Marshal(e.reply.Usage)
	if err != nil {
		log.Println("error:", err)
		return
	}
	t.Log("[RPC][Usage] -> %s", b)
}

func (e *UsageExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

func (m *ManagerRPC) Usage(arg ManagerUsageArg, reply *ManagerUsageReply) error {
	return NewTask("Usage", &UsageExecutor{arg, reply}).Run()
}
