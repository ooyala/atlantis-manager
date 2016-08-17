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
	. "atlantis/manager/rpc/types"
)

type HealthCheckExecutor struct {
	arg   ManagerHealthCheckArg
	reply *ManagerHealthCheckReply
}

func (e *HealthCheckExecutor) Request() interface{} {
	return e.arg
}

func (e *HealthCheckExecutor) Result() interface{} {
	return e.reply
}

func (e *HealthCheckExecutor) Description() string {
	return "HealthCheck"
}

func (e *HealthCheckExecutor) Execute(t *Task) error {
	e.reply.Region = Region
	e.reply.Zone = Zone
	if Tracker.UnderMaintenance() {
		e.reply.Status = StatusMaintenance
	} else {
		e.reply.Status = StatusOk
	}
	t.Log("[RPC][HealthCheck] region:%s zone:%s status:%s", e.reply.Region, e.reply.Zone, e.reply.Status)
	return nil
}

func (e *HealthCheckExecutor) Authorize() error {
	return nil // allow anyone to check health
}

func (e *HealthCheckExecutor) AllowDuringMaintenance() bool {
	return true // allow checking health during maintenance.
}

func (m *ManagerRPC) HealthCheck(arg ManagerHealthCheckArg, reply *ManagerHealthCheckReply) error {
	return NewTask("HealthCheck", &HealthCheckExecutor{arg, reply}).Run()
}
