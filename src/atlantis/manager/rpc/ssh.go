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
	"atlantis/manager/datamodel"
	. "atlantis/manager/rpc/types"
	"atlantis/manager/supervisor"
	"errors"
	"fmt"
)

// NOTE[jigish]: these are simple pass-throughs to supervisor.

type AuthorizeSSHExecutor struct {
	arg   ManagerAuthorizeSSHArg
	reply *ManagerAuthorizeSSHReply
}

func (e *AuthorizeSSHExecutor) Request() interface{} {
	return e.arg
}

func (e *AuthorizeSSHExecutor) Result() interface{} {
	return e.reply
}

func (e *AuthorizeSSHExecutor) Description() string {
	return fmt.Sprintf("%s @ %s : \n%s", e.arg.User, e.arg.ContainerID, e.arg.PublicKey)
}

func (e *AuthorizeSSHExecutor) Execute(t *Task) error {
	if e.arg.PublicKey == "" {
		return errors.New("Please specify an SSH public key.")
	}
	if e.arg.ContainerID == "" {
		return errors.New("Please specify a container id.")
	}
	if e.arg.User == "" {
		return errors.New("Please specify a user.")
	}
	instance, err := datamodel.GetInstance(e.arg.ContainerID)
	if err != nil {
		return err
	}
	ihReply, err := supervisor.AuthorizeSSH(instance.Host, e.arg.ContainerID, e.arg.User, e.arg.PublicKey)
	if err != nil {
		e.reply.Status = StatusError
		return err
	}
	e.reply.Host = instance.Host
	e.reply.Port = ihReply.Port
	e.reply.Status = ihReply.Status
	return err
}

func (e *AuthorizeSSHExecutor) Authorize() error {
	return nil
}

type DeauthorizeSSHExecutor struct {
	arg   ManagerAuthorizeSSHArg
	reply *ManagerAuthorizeSSHReply
}

func (e *DeauthorizeSSHExecutor) Request() interface{} {
	return e.arg
}

func (e *DeauthorizeSSHExecutor) Result() interface{} {
	return e.reply
}

func (e *DeauthorizeSSHExecutor) Description() string {
	return fmt.Sprintf("%s @ %s", e.arg.User, e.arg.ContainerID)
}

func (e *DeauthorizeSSHExecutor) Execute(t *Task) error {
	if e.arg.ContainerID == "" {
		return errors.New("Please specify a container id.")
	}
	if e.arg.User == "" {
		return errors.New("Please specify a user.")
	}
	instance, err := datamodel.GetInstance(e.arg.ContainerID)
	if err != nil {
		return err
	}
	ihReply, err := supervisor.DeauthorizeSSH(instance.Host, e.arg.ContainerID, e.arg.User)
	e.reply.Status = ihReply.Status
	return err
}

func (e *DeauthorizeSSHExecutor) Authorize() error {
	return nil
}

func (m *ManagerRPC) AuthorizeSSH(arg ManagerAuthorizeSSHArg, reply *ManagerAuthorizeSSHReply) error {
	return NewTask("AuthorizeSSH", &AuthorizeSSHExecutor{arg, reply}).Run()
}

func (m *ManagerRPC) DeauthorizeSSH(arg ManagerAuthorizeSSHArg, reply *ManagerAuthorizeSSHReply) error {
	return NewTask("DeauthorizeSSH", &DeauthorizeSSHExecutor{arg, reply}).Run()
}
