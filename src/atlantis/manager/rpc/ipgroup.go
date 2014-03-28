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
	"atlantis/manager/netsec"
	. "atlantis/manager/rpc/types"
	"errors"
	"fmt"
)

type UpdateIPGroupExecutor struct {
	arg   ManagerUpdateIPGroupArg
	reply *ManagerUpdateIPGroupReply
}

func (e *UpdateIPGroupExecutor) Request() interface{} {
	return e.arg
}

func (e *UpdateIPGroupExecutor) Result() interface{} {
	return e.reply
}

func (e *UpdateIPGroupExecutor) Description() string {
	return fmt.Sprintf("%s -> %v", e.arg.Name, e.arg.IPs)
}

func (e *UpdateIPGroupExecutor) Authorize() error {
	return AuthorizeSuperUser(&e.arg.ManagerAuthArg)
}

func (e *UpdateIPGroupExecutor) Execute(t *Task) error {
	if e.arg.Name == "" {
		return errors.New("Please specify a Name.")
	}
	if e.arg.IPs == nil {
		return errors.New("Please specify a list of IPs.")
	}
	if err := netsec.UpdateIPGroup(e.arg.Name, e.arg.IPs); err != nil {
		e.reply.Status = StatusError
		return err
	}
	e.reply.Status = StatusOk
	return nil
}

func (m *ManagerRPC) UpdateIPGroup(arg ManagerUpdateIPGroupArg, reply *ManagerUpdateIPGroupReply) error {
	return NewTask("UpdateIPGroup", &UpdateIPGroupExecutor{arg, reply}).Run()
}

type DeleteIPGroupExecutor struct {
	arg   ManagerDeleteIPGroupArg
	reply *ManagerDeleteIPGroupReply
}

func (e *DeleteIPGroupExecutor) Request() interface{} {
	return e.arg
}

func (e *DeleteIPGroupExecutor) Result() interface{} {
	return e.reply
}

func (e *DeleteIPGroupExecutor) Description() string {
	return fmt.Sprintf("%s", e.arg.Name)
}

func (e *DeleteIPGroupExecutor) Authorize() error {
	return AuthorizeSuperUser(&e.arg.ManagerAuthArg)
}

func (e *DeleteIPGroupExecutor) Execute(t *Task) error {
	if e.arg.Name == "" {
		return errors.New("Please specify a Name.")
	}
	if err := netsec.DeleteIPGroup(e.arg.Name); err != nil {
		e.reply.Status = StatusError
		return err
	}
	e.reply.Status = StatusOk
	return nil
}

func (m *ManagerRPC) DeleteIPGroup(arg ManagerDeleteIPGroupArg, reply *ManagerDeleteIPGroupReply) error {
	return NewTask("DeleteIPGroup", &DeleteIPGroupExecutor{arg, reply}).Run()
}

type GetIPGroupExecutor struct {
	arg   ManagerGetIPGroupArg
	reply *ManagerGetIPGroupReply
}

func (e *GetIPGroupExecutor) Request() interface{} {
	return e.arg
}

func (e *GetIPGroupExecutor) Result() interface{} {
	return e.reply
}

func (e *GetIPGroupExecutor) Description() string {
	return fmt.Sprintf("%s", e.arg.Name)
}

func (e *GetIPGroupExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

func (e *GetIPGroupExecutor) Execute(t *Task) (err error) {
	if e.arg.Name == "" {
		return errors.New("Please specify a Name.")
	}
	e.reply.IPGroup, err = netsec.GetIPGroup(e.arg.Name)
	if err != nil {
		e.reply.Status = StatusError
		return err
	}
	e.reply.Status = StatusOk
	return nil
}

func (m *ManagerRPC) GetIPGroup(arg ManagerGetIPGroupArg, reply *ManagerGetIPGroupReply) error {
	return NewTask("GetIPGroup", &GetIPGroupExecutor{arg, reply}).Run()
}
