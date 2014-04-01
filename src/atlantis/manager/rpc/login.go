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
)

type LoginExecutor struct {
	arg   ManagerLoginArg
	reply *ManagerLoginReply
}

func (e *LoginExecutor) Request() interface{} {
	return e.arg
}

func (e *LoginExecutor) Result() interface{} {
	return e.reply
}

func (e *LoginExecutor) Description() string {
	return "[" + e.arg.User + "] Login"
}

func (e *LoginExecutor) Execute(t *Task) error {
	e.reply.Secret = e.arg.Secret
	e.reply.LoggedIn = false
	auther := Authorizer{e.arg.User, e.arg.Pass, e.arg.Secret}
	err := auther.Authenticate()
	if err != nil {
		return err
	}
	e.reply.Secret = auther.Secret
	// reply.Secret should be blank if the username/pass is incorrect
	if e.reply.Secret != "" {
		e.reply.LoggedIn = true
	}
	return nil
}

func (e *LoginExecutor) Authorize() error {
	return nil
}

func (m *ManagerRPC) Login(arg ManagerLoginArg, reply *ManagerLoginReply) error {
	return NewTask("Login", &LoginExecutor{arg, reply}).Run()
}
