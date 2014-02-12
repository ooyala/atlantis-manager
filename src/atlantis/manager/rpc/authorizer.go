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
	"atlantis/manager/ldap"
	. "atlantis/manager/rpc/types"
	"errors"
)

type Authorizer struct {
	User     string
	Password string
	Secret   string
}

func (a *Authorizer) Authenticate() (err error) {
	a.Secret, err = ldap.Login(a.User, a.Password, a.Secret)
	return err
}

func SimpleAuthorize(AuthArg *ManagerAuthArg) error {
	user, password, secret := AuthArg.Credentials()
	auther := Authorizer{user, password, secret}
	if err := auther.Authenticate(); err != nil {
		return err
	}
	return nil
}

func AuthorizeTeamAdmin(AuthArg *ManagerAuthArg, team string) error {
	if err := SimpleAuthorize(AuthArg); err != nil {
		return err
	}
	req := ManagerTeamAdminArg{*AuthArg, team}
	var res ManagerTeamAdminReply
	err := NewTask("Authorizer-IsTeamAdmin", &IsTeamAdminExecutor{req, &res}).Run()
	if err != nil || !res.IsAdmin {
		return errors.New("Not a Team Admin")
	}
	return nil
}

func AuthorizeApp(AuthArg *ManagerAuthArg, app string) error {
	if err := SimpleAuthorize(AuthArg); err != nil {
		return err
	}
	var reply ManagerIsAppAllowedReply
	arg := ManagerIsAppAllowedArg{ManagerAuthArg: *AuthArg, App: app, User: AuthArg.User}
	err := NewTask("Authorizer-IsAppAllowed", &IsAppAllowedExecutor{arg, &reply}).Run()
	if err != nil || !reply.IsAllowed {
		return errors.New("Not Authorized to Deploy App")
	}
	return nil
}

func AuthorizeSuperUser(AuthArg *ManagerAuthArg) error {
	if err := SimpleAuthorize(AuthArg); err != nil {
		return err
	}
	var reply ManagerSuperUserReply
	arg := ManagerSuperUserArg{*AuthArg}
	err := NewTask("Authorizer-IsSuperUser", &IsSuperUserExecutor{arg, &reply}).Run()
	if err != nil || !reply.IsSuperUser {
		return errors.New("Not a Super User")
	}
	return nil
}
