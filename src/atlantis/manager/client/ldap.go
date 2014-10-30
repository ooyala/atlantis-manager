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
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/mewpkg/gopass"
	"io/ioutil"
	"os"
	"strings"
)

type LoginCommand struct {
	User     string `short:"u" long:"user" description:"the username of the user"`
	Password string `short:"p" long:"password" description:"the password of the user"`
}

const tokenFile = "manager-tokens"

func (c *LoginCommand) Execute(args []string) error {
	overlayConfig()
	var err error
	Log("Logging in over RPC")
	_, err = AutoLogin(c.User, c.Password)
	for err != nil {
		// Fall back to manual login if anything goes wrong
		PromptUsername(&c.User)
		err = PromptPassword(&c.Password)
		if err != nil {
			return err
		}
		_, err = AutoLogin(c.User, c.Password)
		if err != nil {
			Log("Login unsuccessful")
		}
	}
	return OutputEmpty()
}

// ----------------------------------------------------------------------------------------------------------
// User and Application Authorization
// ----------------------------------------------------------------------------------------------------------

var ldapOperationsEnabled bool = true

type CreateTeamCommand struct {
	Team string `short:"t" long:"team" description:"the name of the team"`
}

func (c *CreateTeamCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	if !ldapOperationsEnabled {
		Log("LDAP Operations are not enabled.")
		return OutputEmpty()
	}
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	if c.Team != "" {
		auth := ManagerAuthArg{user, "", secret}
		arg := ManagerTeamArg{auth, c.Team}
		var reply ManagerTeamReply
		if err := rpcClient.Call("CreateTeam", arg, &reply); err != nil {
			return OutputError(err)
		}
		Log(c.Team + " created as a team.")
		return OutputEmpty()
	} else {
		return Output(map[string]interface{}{}, nil, errors.New("Missing Team Argument"))
	}
}

type DeleteTeamCommand struct {
	Team string `short:"t" long:"team" description:"the name of the team"`
}

func (c *DeleteTeamCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	if !ldapOperationsEnabled {
		Log("LDAP Operations are not enabled.")
		return OutputEmpty()
	}
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	if err := IsTeamAdmin(c.Team); err != nil {
		return err
	}
	// We're not using the user argument for anything when it comes to removal
	// since it is already binded with the current LDAP session
	if c.Team == "" {
		return Output(map[string]interface{}{}, nil, errors.New("Missing Team Argument"))
	}
	auth := ManagerAuthArg{user, "", secret}
	arg := ManagerTeamArg{auth, c.Team}
	var reply ManagerTeamReply
	if err := rpcClient.Call("DeleteTeam", arg, &reply); err != nil {
		return OutputError(err)
	}
	Log(c.Team + " deleted as a team.")
	return OutputEmpty()
}

type ListTeamsCommand struct {
}

func (c *ListTeamsCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	if !ldapOperationsEnabled {
		Log("LDAP Operations are not enabled.")
		return OutputEmpty()
	}
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	auth := ManagerAuthArg{user, "", secret}
	arg := ManagerListTeamsArg{auth}
	var reply ManagerListTeamsReply
	if err := rpcClient.Call("ListTeams", arg, &reply); err != nil {
		return OutputError(err)
	}
	Log("-> teams:")
	for _, team := range reply.Teams {
		Log("->   %s", team)
	}
	return Output(map[string]interface{}{"teams": reply.Teams}, reply.Teams, nil)
}

type ListTeamEmailsCommand struct {
	Team string `short:"t" long:"team" description:"the name of the team"`
}

func (c *ListTeamEmailsCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	if !ldapOperationsEnabled {
		Log("LDAP Operations are not enabled.")
		return OutputEmpty()
	}
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	auth := ManagerAuthArg{user, "", secret}
	arg := ManagerListTeamEmailsArg{auth, c.Team}
	var reply ManagerListTeamEmailsReply
	if err := rpcClient.Call("ListTeamEmails", arg, &reply); err != nil {
		return OutputError(err)
	}
	Log("-> Emails:")
	for _, val := range reply.TeamEmails {
		Log("->   %s", val)
	}
	return Output(map[string]interface{}{"emails": reply.TeamEmails}, reply.TeamEmails, nil)
}

type ListTeamAdminsCommand struct {
	Team string `short:"t" long:"team" description:"the name of the team"`
}

func (c *ListTeamAdminsCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	if !ldapOperationsEnabled {
		Log("LDAP Operations are not enabled.")
		return OutputEmpty()
	}
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	auth := ManagerAuthArg{user, "", secret}
	arg := ManagerListTeamAdminsArg{auth, c.Team}
	var reply ManagerListTeamAdminsReply
	if err := rpcClient.Call("ListTeamAdmins", arg, &reply); err != nil {
		return OutputError(err)
	}
	Log("-> Admins:")
	for _, val := range reply.TeamAdmins {
		Log("->   %s", val)
	}
	return Output(map[string]interface{}{"admins": reply.TeamAdmins}, reply.TeamAdmins, nil)
}

type ListTeamMembersCommand struct {
	Team string `short:"t" long:"team" description:"the name of the team"`
}

func (c *ListTeamMembersCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	if !ldapOperationsEnabled {
		Log("LDAP Operations are not enabled.")
		return OutputEmpty()
	}
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	auth := ManagerAuthArg{user, "", secret}
	arg := ManagerListTeamMembersArg{auth, c.Team}
	var reply ManagerListTeamMembersReply
	if err := rpcClient.Call("ListTeamMembers", arg, &reply); err != nil {
		return OutputError(err)
	}
	Log("-> Members:")
	for _, val := range reply.TeamMembers {
		Log("->   %s", val)
	}
	return Output(map[string]interface{}{"members": reply.TeamMembers}, reply.TeamMembers, nil)
}

type ListTeamAppsCommand struct {
	Team string `short:"t" long:"team" description:"the name of the team"`
}

func (c *ListTeamAppsCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	if !ldapOperationsEnabled {
		Log("LDAP Operations are not enabled.")
		return OutputEmpty()
	}
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	auth := ManagerAuthArg{user, "", secret}
	arg := ManagerListTeamAppsArg{auth, c.Team}
	var reply ManagerListTeamAppsReply
	if err := rpcClient.Call("ListTeamApps", arg, &reply); err != nil {
		return OutputError(err)
	}
	Log("-> Apps:")
	for _, val := range reply.TeamApps {
		Log("->   %s", val)
	}
	return Output(map[string]interface{}{"apps": reply.TeamApps}, reply.TeamApps, nil)
}

type AddTeamMemberCommand struct {
	User string `short:"u" long:"user" description:"the name of the user"`
	Team string `short:"t" long:"team" description:"the name of the team"`
}

func (c *AddTeamMemberCommand) Execute(args []string) error {
	if err := ModifyTeamMember("AddTeamMember", c.Team, c.User); err != nil {
		return OutputError(err)
	}
	Log(c.User + " added to team : " + c.Team)
	return OutputEmpty()

}

type RemoveTeamMemberCommand struct {
	User string `short:"u" long:"user" description:"the name of the user"`
	Team string `short:"t" long:"team" description:"the name of the team"`
}

func (c *RemoveTeamMemberCommand) Execute(args []string) error {
	if err := ModifyTeamMember("RemoveTeamMember", c.Team, c.User); err != nil {
		return OutputError(err)
	}
	Log(c.User + " removed to team : " + c.Team)
	return OutputEmpty()
}

func ModifyTeamMember(action, team, user string) error {
	if err := Init(); err != nil {
		return err
	}
	if !ldapOperationsEnabled {
		Log("LDAP Operations are not enabled.")
		return OutputEmpty()
	}
	cuser, secret, err := GetSecret()
	if err != nil {
		return err
	}
	if err := IsTeamAdmin(team); err != nil {
		return err
	}
	if user == "" && team == "" {
		return errors.New("Missing User/Team Arguments")
	}
	auth := ManagerAuthArg{cuser, "", secret}
	arg := ManagerTeamMemberArg{auth, team, user}
	var reply ManagerTeamMemberReply
	if err := rpcClient.Call(action, arg, &reply); err != nil {
		return err
	}
	return nil
}

type AddTeamAdminCommand struct {
	User string `short:"u" long:"user" description:"the name of the user"`
	Team string `short:"t" long:"team" description:"the name of the team"`
}

func (c *AddTeamAdminCommand) Execute(args []string) error {
	if err := ModifyTeamAdmin("AddTeamAdmin", c.User, c.Team); err != nil {
		return OutputError(err)
	}
	Log(c.User + " is now a team admin of : " + c.Team)
	return OutputEmpty()
}

type RemoveTeamAdminCommand struct {
	User string `short:"u" long:"user" description:"the name of the user"`
	Team string `short:"t" long:"team" description:"the name of the team"`
}

func (c *RemoveTeamAdminCommand) Execute(args []string) error {
	if err := ModifyTeamAdmin("RemoveTeamAdmin", c.User, c.Team); err != nil {
		return OutputError(err)
	}
	Log(c.User + " is no longer a the team admin of : " + c.Team)
	return OutputEmpty()
}

func ModifyTeamAdmin(action, user, team string) error {
	if err := Init(); err != nil {
		return err
	}
	if !ldapOperationsEnabled {
		Log("LDAP Operations are not enabled.")
		return OutputEmpty()
	}
	cuser, secret, err := GetSecret()
	if err != nil {
		return err
	}
	if err := IsTeamAdmin(team); err != nil {
		return err
	}
	if user == "" && team == "" {
		return errors.New("Missing User/Team Arguments")
	}
	auth := ManagerAuthArg{cuser, "", secret}
	arg := ManagerModifyTeamAdminArg{auth, team, user}
	var reply ManagerAppReply
	if err := rpcClient.Call(action, arg, &reply); err != nil {
		return err
	}
	return nil
}

type AllowTeamAppCommand struct {
	App  string `short:"a" long:"app" description:"the name of app"`
	Team string `short:"t" long:"team" description:"the name of the team"`
}

func (c *AllowTeamAppCommand) Execute(args []string) error {
	if err := ModifyAllowedTeamApp("AllowApp", c.Team, c.App); err != nil {
		return OutputError(err)
	}
	Log("Team: " + c.Team + " has permissions for APP: " + c.App)
	return OutputEmpty()
}

type DisallowTeamAppCommand struct {
	App  string `short:"a" long:"app" description:"the name of app"`
	Team string `short:"t" long:"team" description:"the name of the team"`
}

func (c *DisallowTeamAppCommand) Execute(args []string) error {
	if err := ModifyAllowedTeamApp("DisallowApp", c.Team, c.App); err != nil {
		return OutputError(err)
	}
	Log("Team: " + c.Team + " no longer has permissions for APP: " + c.App)
	return OutputEmpty()
}

func ModifyAllowedTeamApp(action, team, app string) error {
	if err := Init(); err != nil {
		return err
	}
	if !ldapOperationsEnabled {
		Log("LDAP Operations are not enabled.")
		return OutputEmpty()
	}
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	if err := IsTeamAdmin(team); err != nil {
		return err
	}
	if app == "" && team == "" {
		return errors.New("Missing App/Team Arguments")
	}
	auth := ManagerAuthArg{user, "", secret}
	arg := ManagerAppArg{auth, app, team}
	var reply ManagerAppReply
	if err := rpcClient.Call(action, arg, &reply); err != nil {
		return err
	}
	return nil
}

type AddTeamEmailCommand struct {
	Email string `short:"e" long:"email" description:"email to add"`
	Team  string `short:"t" long:"team" description:"the name of the team"`
}

func (c *AddTeamEmailCommand) Execute(args []string) error {
	if err := ModifyTeamEmail("AddTeamEmail", c.Email, c.Team); err != nil {
		return OutputError(err)
	}
	Log("Email: " + c.Email + " added to team - " + c.Team)
	return OutputEmpty()
}

type RemoveTeamEmailCommand struct {
	Email string `short:"e" long:"email" description:"email to add"`
	Team  string `short:"t" long:"team" description:"the name of the team"`
}

func (c *RemoveTeamEmailCommand) Execute(args []string) error {
	if err := ModifyTeamEmail("RemoveTeamEmail", c.Email, c.Team); err != nil {
		return OutputError(err)
	}
	Log("Email: " + c.Email + " removed from team - " + c.Team)
	return OutputEmpty()
}

func ModifyTeamEmail(action, email, team string) error {
	if err := Init(); err != nil {
		return err
	}
	if !ldapOperationsEnabled {
		Log("LDAP Operations are not enabled.")
		return OutputEmpty()
	}
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	if err := IsTeamAdmin(team); err != nil {
		return err
	}
	if email == "" && team == "" {
		return errors.New("Missing Email/Team Arguments")
	}
	auth := ManagerAuthArg{user, "", secret}
	arg := ManagerEmailArg{auth, team, email}
	var reply ManagerEmailReply
	if err := rpcClient.Call(action, arg, &reply); err != nil {
		return err
	}
	return nil
}

type IsAppAllowedCommand struct {
	App  string `short:"a" long:"app" description:"the app to check permissions for"`
	User string `short:"u" long:"user" description:"[superuser only] the user to check permissions for"`
}

func (c *IsAppAllowedCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return err
	}
	if !ldapOperationsEnabled {
		Log("LDAP Operations are not enabled. All apps are allowed.")
		return OutputEmpty()
	}
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	auth := ManagerAuthArg{user, "", secret}
	arg := ManagerIsAppAllowedArg{ManagerAuthArg: auth, User: c.User, App: c.App}
	var reply ManagerIsAppAllowedReply
	if err = rpcClient.Call("IsAppAllowed", arg, &reply); err == nil {
		Log("-> %t", reply.IsAllowed)
	}
	return Output(map[string]interface{}{"isAllowed": reply.IsAllowed}, reply.IsAllowed, err)
}

type ListAllowedAppsCommand struct {
	User string `short:"u" long:"user" description:"[superuser only] the user to check permissions for"`
}

func (c *ListAllowedAppsCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return err
	}
	if !ldapOperationsEnabled {
		Log("LDAP Operations are not enabled. All apps are allowed.")
		return OutputEmpty()
	}
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	auth := ManagerAuthArg{user, "", secret}
	arg := ManagerListAllowedAppsArg{ManagerAuthArg: auth, User: c.User}
	var reply ManagerListAllowedAppsReply
	if err = rpcClient.Call("ListAllowedApps", arg, &reply); err == nil {
		Log("-> apps:")
		for _, app := range reply.Apps {
			Log("->   %s", app)
		}
	}
	return Output(map[string]interface{}{"apps": reply.Apps}, reply.Apps, err)
}

// ----------------------------------------------------------------------------------------------------------
// Login Related Functions
// ----------------------------------------------------------------------------------------------------------

func AutoLoginDefault() error {
	secret, err := AutoLogin("", "")
	if secret.LoggedIn == false && err == nil {
		ldapOperationsEnabled = false
	}
	tries := 0
	for err != nil && tries < 2 {
		var user string
		var password string
		PromptUsername(&user)
		PromptPassword(&password)
		_, err = AutoLogin(user, password)
		tries++
	}
	return err
}

func AutoLogin(overrideUser, overridePassword string) (ManagerLoginReply, error) {
	var reply ManagerLoginReply
	password := ""
	user, secrets, err := GetSecrets()
	rpcClient.User = user
	rpcClient.Secrets = secrets
	secret := secrets[rpcClient.Opts.RPCHostAndPort()]
	if err != nil {
		return reply, err
	}
	if overrideUser != "" {
		user = overrideUser
	}
	if overridePassword != "" {
		password = overridePassword
	}
	if user == "" {
		PromptUsername(&user)
	}
	if secret == "" && overridePassword == "" {
		PromptPassword(&password)
	}
	// Attempt to Auto-Login assuming we have user/secret
	arg := ManagerLoginArg{user, password, secret}
	if err := rpcClient.Call("Login", arg, &reply); err != nil {
		return reply, err
	}
	if reply.LoggedIn {
		if err := SaveSecret(arg.User, reply.Secret); err != nil {
			return reply, err
		}
	}
	return reply, nil
}

func PromptUsername(user *string) {
	*user = ""
	for *user == "" {
		fmt.Printf("LDAP Username: ")
		fmt.Scanln(user)
	}
}

func PromptPassword(pass *string) error {
	*pass = ""
	var err error
	for *pass == "" {
		*pass, err = gopass.GetPass("LDAP Password: ")
		if err != nil {
			return err
		}
	}
	return nil
}

func GetSecret() (string, string, error) {
	user, secrets, _ := GetSecrets()
	return user, secrets[rpcClient.Opts.RPCHostAndPort()], nil
}

type TokenFileFormat struct {
	User    string
	Secrets map[string]string
}

func GetSecrets() (string, map[string]string, error) {
	path := strings.Replace(cfg.KeyPath, "~", os.Getenv("HOME"), 1)
	var filePath string = path + "/" + tokenFile
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", map[string]string{}, nil
	}
	tokenData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", map[string]string{}, err
	}
	var token TokenFileFormat
	if _, err := toml.Decode(string(tokenData), &token); err != nil {
		return "", map[string]string{}, err
	}
	return token.User, token.Secrets, nil
}

func SaveSecret(user string, secret string) error {
	_, secrets, _ := GetSecrets()
	secrets[rpcClient.Opts.RPCHostAndPort()] = secret

	rpcClient.User = user
	rpcClient.Secrets = secrets

	path := strings.Replace(cfg.KeyPath, "~", os.Getenv("HOME"), 1)
	if err := os.MkdirAll(path, 0700); err != nil {
		return err
	}
	file, err := os.Create(path + "/" + tokenFile)
	if err != nil {
		return err
	}
	token := TokenFileFormat{user, secrets}
	if err := toml.NewEncoder(file).Encode(token); err != nil {
		return err
	}
	return nil
}

func IsTeamAdmin(team string) error {
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	if user == "" {
		return errors.New("Not Logged In")
	}
	if team == "" {
		return errors.New("No team specified")
	}
	auth := ManagerAuthArg{user, "", secret}
	arg := ManagerTeamArg{auth, team}
	var reply ManagerTeamAdminReply

	if err := rpcClient.Call("IsTeamAdmin", arg, &reply); err != nil {
		return err
	}
	if reply.IsAdmin == false {
		return errors.New("TEAM: Permission Denied")
	}
	return nil
}

func IsSuperUser() error {
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	if user == "" {
		return errors.New("Not in SuperUsers Group")
	}
	auth := ManagerAuthArg{user, "", secret}
	arg := ManagerSuperUserArg{auth}
	var reply ManagerSuperUserReply
	if err := rpcClient.Call("IsSuperUser", arg, &reply); err != nil {
		return err
	}
	if reply.IsSuperUser == false {
		return errors.New("Permission Denied")
	}
	return nil
}
