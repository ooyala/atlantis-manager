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
	aldap "atlantis/manager/ldap"
	. "atlantis/manager/rpc/types"
	"errors"
	"fmt"
	"github.com/mavricknz/ldap"
	"strings"
)

// ----------------------------------------------------------------------------------------------------------
// Teams
// ----------------------------------------------------------------------------------------------------------

type CreateTeamExecutor struct {
	arg   ManagerTeamArg
	reply *ManagerTeamReply
}

func (e *CreateTeamExecutor) Request() interface{} {
	return e.arg
}

func (e *CreateTeamExecutor) Result() interface{} {
	return e.reply
}

func (e *CreateTeamExecutor) Description() string {
	return fmt.Sprintf("%s : %s", e.arg.User, e.arg.Team)
}

func (e *CreateTeamExecutor) Execute(t *Task) error {
	conn, err := InitConnection(&e.arg.ManagerAuthArg)
	if err != nil {
		return err
	}

	if TeamExists(e.arg.Team, &e.arg.ManagerAuthArg) {
		return errors.New("Team Already Exists")
	}

	var addDNs []string = []string{aldap.TeamCommonName + "=" + e.arg.Team + "," + aldap.TeamOu}
	var Attrs []ldap.EntryAttribute = []ldap.EntryAttribute{
		ldap.EntryAttribute{
			Name:   "objectclass",
			Values: []string{aldap.TeamClass, "groupOfNames", "top"},
		},
		ldap.EntryAttribute{
			Name:   aldap.TeamAdminAttr,
			Values: []string{aldap.UserCommonName + "=" + e.arg.User + "," + aldap.UserOu},
		},
		ldap.EntryAttribute{
			Name:   aldap.TeamCommonName,
			Values: []string{e.arg.Team},
		},
		ldap.EntryAttribute{
			Name:   aldap.UsernameAttr,
			Values: []string{aldap.UserCommonName + "=" + e.arg.User + "," + aldap.UserOu},
		},
	}
	addReq := ldap.NewAddRequest(addDNs[0])
	for _, attr := range Attrs {
		addReq.AddAttribute(&attr)
	}
	if err := conn.Add(addReq); err != nil {
		return err
	}
	return nil
}

func (e *CreateTeamExecutor) Authorize() error {
	if err := checkRole("permissions", "write"); err != nil {
		return err
	}
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

type DeleteTeamExecutor struct {
	arg   ManagerTeamArg
	reply *ManagerTeamReply
}

func (e *DeleteTeamExecutor) Request() interface{} {
	return e.arg
}

func (e *DeleteTeamExecutor) Result() interface{} {
	return e.reply
}

func (e *DeleteTeamExecutor) Description() string {
	return fmt.Sprintf("%s", e.arg.Team)
}

func (e *DeleteTeamExecutor) Execute(t *Task) error {
	conn, err := InitConnection(&e.arg.ManagerAuthArg)
	if err != nil {
		return err
	}

	if !TeamExists(e.arg.Team, &e.arg.ManagerAuthArg) {
		return errors.New("Team Does Not Exist")
	}

	delReq := ldap.NewDeleteRequest(aldap.TeamCommonName + "=" + e.arg.Team + "," + aldap.TeamOu)
	if err := conn.Delete(delReq); err != nil {
		return err
	}

	return nil
}

func (e *DeleteTeamExecutor) Authorize() error {
	if err := checkRole("permissions", "write"); err != nil {
		return err
	}
	return AuthorizeTeamAdmin(&e.arg.ManagerAuthArg, e.arg.Team)
}

func (m *ManagerRPC) CreateTeam(arg ManagerTeamArg, reply *ManagerTeamReply) error {
	return NewTask("CreateTeam", &CreateTeamExecutor{arg, reply}).Run()
}

func (m *ManagerRPC) DeleteTeam(arg ManagerTeamArg, reply *ManagerTeamReply) error {
	return NewTask("DeleteTeam", &DeleteTeamExecutor{arg, reply}).Run()
}

type AddTeamEmailExecutor struct {
	arg   ManagerEmailArg
	reply *ManagerEmailReply
}

func (e *AddTeamEmailExecutor) Request() interface{} {
	return e.arg
}

func (e *AddTeamEmailExecutor) Result() interface{} {
	return e.reply
}

func (e *AddTeamEmailExecutor) Description() string {
	return fmt.Sprintf("%s : %s", e.arg.Email, e.arg.Team)
}

func (e *AddTeamEmailExecutor) Execute(t *Task) error {
	return ModifyTeamEmail(ldap.ModAdd, e.arg, e.reply)
}

func (e *AddTeamEmailExecutor) Authorize() error {
	if err := checkRole("permissions", "write"); err != nil {
		return err
	}
	return AuthorizeTeamAdmin(&e.arg.ManagerAuthArg, e.arg.Team)
}

type RemoveTeamEmailExecutor struct {
	arg   ManagerEmailArg
	reply *ManagerEmailReply
}

func (e *RemoveTeamEmailExecutor) Request() interface{} {
	return e.arg
}

func (e *RemoveTeamEmailExecutor) Result() interface{} {
	return e.reply
}

func (e *RemoveTeamEmailExecutor) Description() string {
	return fmt.Sprintf("%s : %s", e.arg.Email, e.arg.Team)
}

func (e *RemoveTeamEmailExecutor) Execute(t *Task) error {
	return ModifyTeamEmail(ldap.ModDelete, e.arg, e.reply)
}

func (e *RemoveTeamEmailExecutor) Authorize() error {
	if err := checkRole("permissions", "write"); err != nil {
		return err
	}
	return AuthorizeTeamAdmin(&e.arg.ManagerAuthArg, e.arg.Team)
}

func ModifyTeamEmail(action int, arg ManagerEmailArg, reply *ManagerEmailReply) error {
	conn, err := InitConnection(&arg.ManagerAuthArg)
	if err != nil {
		return err
	}

	if !TeamExists(arg.Team, &arg.ManagerAuthArg) {
		return errors.New("Team Does Not Exist")
	}

	if action == ldap.ModDelete && !EmailExists(arg.Email, arg.Team, &arg.ManagerAuthArg) {
		return errors.New("Email does not exist.")
	} else if action == ldap.ModAdd && EmailExists(arg.Email, arg.Team, &arg.ManagerAuthArg) {
		return errors.New("Email already exists.")
	}

	var modDNs []string = []string{aldap.TeamCommonName + "=" + arg.Team + "," + aldap.TeamOu}
	var Attrs []string = []string{"email"}
	var vals []string = []string{arg.Email}
	modReq := ldap.NewModifyRequest(modDNs[0])
	mod := ldap.NewMod(uint8(action), Attrs[0], vals)
	modReq.AddMod(mod)
	if err := conn.Modify(modReq); err != nil {
		return err
	}

	return nil
}

func (m *ManagerRPC) AddTeamEmail(arg ManagerEmailArg, reply *ManagerEmailReply) error {
	return NewTask("AddTeamEmail", &AddTeamEmailExecutor{arg, reply}).Run()
}

func (m *ManagerRPC) RemoveTeamEmail(arg ManagerEmailArg, reply *ManagerEmailReply) error {
	return NewTask("RemoveTeamEmail", &RemoveTeamEmailExecutor{arg, reply}).Run()
}

type AddTeamAdminExecutor struct {
	arg   ManagerModifyTeamAdminArg
	reply *ManagerModifyTeamAdminReply
}

func (e *AddTeamAdminExecutor) Request() interface{} {
	return e.arg
}

func (e *AddTeamAdminExecutor) Result() interface{} {
	return e.reply
}

func (e *AddTeamAdminExecutor) Description() string {
	return fmt.Sprintf("%s : %s", e.arg.User, e.arg.Team)
}

func (e *AddTeamAdminExecutor) Execute(t *Task) error {
	return ModifyTeamAdmin(ldap.ModAdd, e.arg, e.reply)
}

func (e *AddTeamAdminExecutor) Authorize() error {
	if err := checkRole("permissions", "write"); err != nil {
		return err
	}
	req := ManagerTeamAdminArg{e.arg.ManagerAuthArg, e.arg.Team}
	var res ManagerTeamAdminReply
	err := NewTask("AddTeamAdmin-IsTeamAdmin", &IsTeamAdminExecutor{req, &res}).Run()
	if err != nil || !res.IsAdmin {
		return errors.New("Permission denied")
	}
	return nil
}

type RemoveTeamAdminExecutor struct {
	arg   ManagerModifyTeamAdminArg
	reply *ManagerModifyTeamAdminReply
}

func (e *RemoveTeamAdminExecutor) Request() interface{} {
	return e.arg
}

func (e *RemoveTeamAdminExecutor) Result() interface{} {
	return e.reply
}

func (e *RemoveTeamAdminExecutor) Description() string {
	return fmt.Sprintf("%s : %s", e.arg.User, e.arg.Team)
}

func (e *RemoveTeamAdminExecutor) Execute(t *Task) error {
	return ModifyTeamAdmin(ldap.ModDelete, e.arg, e.reply)
}

func (e *RemoveTeamAdminExecutor) Authorize() error {
	if err := checkRole("permissions", "write"); err != nil {
		return err
	}
	req := ManagerTeamAdminArg{e.arg.ManagerAuthArg, e.arg.Team}
	var res ManagerTeamAdminReply
	err := NewTask("RemoveTeamAdmin-IsTeamAdmin", &IsTeamAdminExecutor{req, &res}).Run()
	if err != nil || !res.IsAdmin {
		return errors.New("Permission denied")
	}
	return nil
}

func ModifyTeamAdmin(action int, arg ManagerModifyTeamAdminArg, reply *ManagerModifyTeamAdminReply) error {
	conn, err := InitConnection(&arg.ManagerAuthArg)
	if err != nil {
		return err
	}

	if !UserExists(arg.User, &arg.ManagerAuthArg) {
		return errors.New("User does not exist")
	}

	if !TeamExists(arg.Team, &arg.ManagerAuthArg) {
		return errors.New("Team Does Not Exist")
	}

	var modDNs []string = []string{aldap.TeamCommonName + "=" + arg.Team + "," + aldap.TeamOu}
	var Attrs []string = []string{aldap.TeamAdminAttr}
	var vals []string = []string{aldap.UserCommonName + "=" + arg.User + "," + aldap.UserOu}
	modReq := ldap.NewModifyRequest(modDNs[0])
	mod := ldap.NewMod(uint8(action), Attrs[0], vals)
	modReq.AddMod(mod)
	if err := conn.Modify(modReq); err != nil {
		return err
	}

	return nil
}

func (m *ManagerRPC) AddTeamAdmin(arg ManagerModifyTeamAdminArg, reply *ManagerModifyTeamAdminReply) error {
	return NewTask("AddTeamAdmin", &AddTeamAdminExecutor{arg, reply}).Run()
}

func (m *ManagerRPC) RemoveTeamAdmin(arg ManagerModifyTeamAdminArg, reply *ManagerModifyTeamAdminReply) error {
	return NewTask("RemoveTeamAdmin", &RemoveTeamAdminExecutor{arg, reply}).Run()
}

type AddTeamMemberExecutor struct {
	arg   ManagerTeamMemberArg
	reply *ManagerTeamMemberReply
}

func (e *AddTeamMemberExecutor) Request() interface{} {
	return e.arg
}

func (e *AddTeamMemberExecutor) Result() interface{} {
	return e.reply
}

func (e *AddTeamMemberExecutor) Description() string {
	return fmt.Sprintf("%s : %s", e.arg.User, e.arg.Team)
}

func (e *AddTeamMemberExecutor) Execute(t *Task) error {
	return ModifyTeamMember(ldap.ModAdd, e.arg, e.reply)
}

func (e *AddTeamMemberExecutor) Authorize() error {
	if err := checkRole("permissions", "write"); err != nil {
		return err
	}
	return AuthorizeTeamAdmin(&e.arg.ManagerAuthArg, e.arg.Team)
}

type RemoveTeamMemberExecutor struct {
	arg   ManagerTeamMemberArg
	reply *ManagerTeamMemberReply
}

func (e *RemoveTeamMemberExecutor) Request() interface{} {
	return e.arg
}

func (e *RemoveTeamMemberExecutor) Result() interface{} {
	return e.reply
}

func (e *RemoveTeamMemberExecutor) Description() string {
	return fmt.Sprintf("%s : %s", e.arg.User, e.arg.Team)
}

func (e *RemoveTeamMemberExecutor) Execute(t *Task) error {
	return ModifyTeamMember(ldap.ModDelete, e.arg, e.reply)
}

func (e *RemoveTeamMemberExecutor) Authorize() error {
	if err := checkRole("permissions", "write"); err != nil {
		return err
	}
	return AuthorizeTeamAdmin(&e.arg.ManagerAuthArg, e.arg.Team)
}

func ModifyTeamMember(action int, arg ManagerTeamMemberArg, reply *ManagerTeamMemberReply) error {
	conn, err := InitConnection(&arg.ManagerAuthArg)
	if err != nil {
		return err
	}
	if action != ldap.ModDelete && !UserExists(arg.User, &arg.ManagerAuthArg) {
		return errors.New("User does not exist")
	}

	if !TeamExists(arg.Team, &arg.ManagerAuthArg) {
		return errors.New("Team Does Not Exist")
	}
	var modDNs []string = []string{aldap.TeamCommonName + "=" + arg.Team + "," + aldap.TeamOu}
	var Attrs []string = []string{aldap.UsernameAttr}
	var vals []string = []string{aldap.UserCommonName + "=" + arg.User + "," + aldap.UserOu}
	modReq := ldap.NewModifyRequest(modDNs[0])
	mod := ldap.NewMod(uint8(action), Attrs[0], vals)
	modReq.AddMod(mod)
	if err := conn.Modify(modReq); err != nil {
		return err
	}

	return nil
}

func (m *ManagerRPC) AddTeamMember(arg ManagerTeamMemberArg, reply *ManagerTeamMemberReply) error {
	return NewTask("AddTeamMember", &AddTeamMemberExecutor{arg, reply}).Run()
}

func (m *ManagerRPC) RemoveTeamMember(arg ManagerTeamMemberArg, reply *ManagerTeamMemberReply) error {
	return NewTask("RemoveTeamMember", &RemoveTeamMemberExecutor{arg, reply}).Run()
}

// ----------------------------------------------------------------------------------------------------------
// Listers
// ----------------------------------------------------------------------------------------------------------

type ListTeamsExecutor struct {
	arg   ManagerListTeamsArg
	reply *ManagerListTeamsReply
}

func (e *ListTeamsExecutor) Request() interface{} {
	return e.arg
}

func (e *ListTeamsExecutor) Result() interface{} {
	return e.reply
}

func (e *ListTeamsExecutor) Description() string {
	return "ListTeams"
}

func (e *ListTeamsExecutor) Execute(t *Task) (err error) {
	e.reply.Teams, err = ListTeams(&e.arg.ManagerAuthArg)
	return err
}

func (e *ListTeamsExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

func (m *ManagerRPC) ListTeams(arg ManagerListTeamsArg, reply *ManagerListTeamsReply) error {
	return NewTask("ListTeams", &ListTeamsExecutor{arg, reply}).Run()
}

type ListTeamEmailsExecutor struct {
	arg   ManagerListTeamEmailsArg
	reply *ManagerListTeamEmailsReply
}

func (e *ListTeamEmailsExecutor) Request() interface{} {
	return e.arg
}

func (e *ListTeamEmailsExecutor) Result() interface{} {
	return e.reply
}

func (e *ListTeamEmailsExecutor) Description() string {
	return "ListTeamEmails"
}

func (e *ListTeamEmailsExecutor) Execute(t *Task) (err error) {
	e.reply.TeamEmails, err = ListTeamEmails(e.arg.Team, &e.arg.ManagerAuthArg)
	return err
}

func (e *ListTeamEmailsExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

func (m *ManagerRPC) ListTeamEmails(arg ManagerListTeamEmailsArg, reply *ManagerListTeamEmailsReply) error {
	return NewTask("ListTeamEmails", &ListTeamEmailsExecutor{arg, reply}).Run()
}

type ListTeamAdminsExecutor struct {
	arg   ManagerListTeamAdminsArg
	reply *ManagerListTeamAdminsReply
}

func (e *ListTeamAdminsExecutor) Request() interface{} {
	return e.arg
}

func (e *ListTeamAdminsExecutor) Result() interface{} {
	return e.reply
}

func (e *ListTeamAdminsExecutor) Description() string {
	return "ListTeamAdmins"
}

func (e *ListTeamAdminsExecutor) Execute(t *Task) (err error) {
	e.reply.TeamAdmins, err = ListTeamAdmins(e.arg.Team, &e.arg.ManagerAuthArg)
	return err
}

func (e *ListTeamAdminsExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

func (m *ManagerRPC) ListTeamAdmins(arg ManagerListTeamAdminsArg, reply *ManagerListTeamAdminsReply) error {
	return NewTask("ListTeamAdmins", &ListTeamAdminsExecutor{arg, reply}).Run()
}

type ListTeamMembersExecutor struct {
	arg   ManagerListTeamMembersArg
	reply *ManagerListTeamMembersReply
}

func (e *ListTeamMembersExecutor) Request() interface{} {
	return e.arg
}

func (e *ListTeamMembersExecutor) Result() interface{} {
	return e.reply
}

func (e *ListTeamMembersExecutor) Description() string {
	return "ListTeamMembers"
}

func (e *ListTeamMembersExecutor) Execute(t *Task) (err error) {
	e.reply.TeamMembers, err = ListTeamMembers(e.arg.Team, &e.arg.ManagerAuthArg)
	return err
}

func (e *ListTeamMembersExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

func (m *ManagerRPC) ListTeamMembers(arg ManagerListTeamMembersArg, reply *ManagerListTeamMembersReply) error {
	return NewTask("ListTeamMembers", &ListTeamMembersExecutor{arg, reply}).Run()
}

type ListTeamAppsExecutor struct {
	arg   ManagerListTeamAppsArg
	reply *ManagerListTeamAppsReply
}

func (e *ListTeamAppsExecutor) Request() interface{} {
	return e.arg
}

func (e *ListTeamAppsExecutor) Result() interface{} {
	return e.reply
}

func (e *ListTeamAppsExecutor) Description() string {
	return "ListTeamApps"
}

func (e *ListTeamAppsExecutor) Execute(t *Task) (err error) {
	e.reply.TeamApps, err = ListTeamApps(e.arg.Team, &e.arg.ManagerAuthArg)
	return err
}

func (e *ListTeamAppsExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

func (m *ManagerRPC) ListTeamApps(arg ManagerListTeamAppsArg, reply *ManagerListTeamAppsReply) error {
	return NewTask("ListTeamApps", &ListTeamAppsExecutor{arg, reply}).Run()
}

// ----------------------------------------------------------------------------------------------------------
// Apps
// ----------------------------------------------------------------------------------------------------------

type AllowAppExecutor struct {
	arg   ManagerAppArg
	reply *ManagerAppReply
}

func (e *AllowAppExecutor) Request() interface{} {
	return e.arg
}

func (e *AllowAppExecutor) Result() interface{} {
	return e.reply
}

func (e *AllowAppExecutor) Description() string {
	return fmt.Sprintf("%s : %s", e.arg.App, e.arg.Team)
}

func (e *AllowAppExecutor) Execute(t *Task) error {
	conn, err := InitConnection(&e.arg.ManagerAuthArg)
	if err != nil {
		return err
	}

	if !TeamExists(e.arg.Team, &e.arg.ManagerAuthArg) {
		return errors.New("Team Does Not Exist")
	}

	var addDNs []string = []string{aldap.AllowedAppAttr + "=" + e.arg.App + "," + aldap.TeamCommonName + "=" + e.arg.Team + "," + aldap.TeamOu}
	var Attrs []ldap.EntryAttribute = []ldap.EntryAttribute{
		ldap.EntryAttribute{
			Name:   "objectclass",
			Values: []string{aldap.AppClass, "top"},
		},
		ldap.EntryAttribute{
			Name:   aldap.AllowedAppAttr,
			Values: []string{e.arg.App},
		},
	}
	addReq := ldap.NewAddRequest(addDNs[0])
	for _, attr := range Attrs {
		addReq.AddAttribute(&attr)
	}
	if err := conn.Add(addReq); err != nil {
		return err
	}

	return nil
}

func (e *AllowAppExecutor) Authorize() error {
	if err := checkRole("permissions", "write"); err != nil {
		return err
	}
	return AuthorizeTeamAdmin(&e.arg.ManagerAuthArg, e.arg.Team)
}

type DisallowAppExecutor struct {
	arg   ManagerAppArg
	reply *ManagerAppReply
}

func (e *DisallowAppExecutor) Request() interface{} {
	return e.arg
}

func (e *DisallowAppExecutor) Result() interface{} {
	return e.reply
}

func (e *DisallowAppExecutor) Description() string {
	return fmt.Sprintf("%s : %s", e.arg.App, e.arg.Team)
}

func (e *DisallowAppExecutor) Execute(t *Task) error {
	if !TeamExists(e.arg.Team, &e.arg.ManagerAuthArg) {
		return errors.New("Team Does Not Exist")
	}
	conn, err := InitConnection(&e.arg.ManagerAuthArg)
	if err != nil {
		return err
	}
	delReq := ldap.NewDeleteRequest(aldap.AllowedAppAttr + "=" + e.arg.App + "," + aldap.TeamCommonName + "=" + e.arg.Team + "," + aldap.TeamOu)
	if err := conn.Delete(delReq); err != nil {
		return err
	}

	return nil
}

func (e *DisallowAppExecutor) Authorize() error {
	if err := checkRole("permissions", "write"); err != nil {
		return err
	}
	return AuthorizeTeamAdmin(&e.arg.ManagerAuthArg, e.arg.Team)
}

func (m *ManagerRPC) AllowApp(arg ManagerAppArg, reply *ManagerAppReply) error {
	return NewTask("AllowApp", &AllowAppExecutor{arg, reply}).Run()
}

func (m *ManagerRPC) DisallowApp(arg ManagerAppArg, reply *ManagerAppReply) error {
	return NewTask("DisallowApp", &DisallowAppExecutor{arg, reply}).Run()
}

type IsAppAllowedExecutor struct {
	arg   ManagerIsAppAllowedArg
	reply *ManagerIsAppAllowedReply
}

func (e *IsAppAllowedExecutor) Request() interface{} {
	return e.arg
}

func (e *IsAppAllowedExecutor) Result() interface{} {
	return e.reply
}

func (e *IsAppAllowedExecutor) Description() string {
	return fmt.Sprintf("%s : %s", e.arg.App, e.arg.User)
}

func (e *IsAppAllowedExecutor) Execute(t *Task) error {
	if aldap.SkipAuthorization {
		e.reply.IsAllowed = true
		return nil
	}

	var suReply ManagerSuperUserReply
	if err := NewTask("IsAppAllowed-IsSuperUser",
		&IsSuperUserExecutor{ManagerSuperUserArg{e.arg.ManagerAuthArg}, &suReply}).Run(); err != nil {
		return err
	}
	user := e.arg.ManagerAuthArg.User
	if suReply.IsSuperUser && e.arg.User != "" {
		user = e.arg.User
	}

	e.reply.IsAllowed = IsAppAllowed(&e.arg.ManagerAuthArg, user, e.arg.App)
	return nil
}

func (e *IsAppAllowedExecutor) Authorize() error {
	return nil
}

type ListAllowedAppsExecutor struct {
	arg   ManagerListAllowedAppsArg
	reply *ManagerListAllowedAppsReply
}

func (e *ListAllowedAppsExecutor) Request() interface{} {
	return e.arg
}

func (e *ListAllowedAppsExecutor) Result() interface{} {
	return e.reply
}

func (e *ListAllowedAppsExecutor) Description() string {
	return fmt.Sprintf("%s: %s", e.arg.ManagerAuthArg.User, e.arg.User)
}

func (e *ListAllowedAppsExecutor) Execute(t *Task) error {
	if aldap.SkipAuthorization {
		e.reply.Apps = []string{"all"}
		return nil
	}

	var suReply ManagerSuperUserReply
	if err := NewTask("ListAllowedApps-IsSuperUser",
		&IsSuperUserExecutor{ManagerSuperUserArg{e.arg.ManagerAuthArg}, &suReply}).Run(); err != nil {
		return err
	}
	user := e.arg.ManagerAuthArg.User
	if suReply.IsSuperUser && e.arg.User != "" {
		user = e.arg.User
	}
	appMap := GetAllowedApps(&e.arg.ManagerAuthArg, user)
	e.reply.Apps = []string{}
	for app, isAllowed := range appMap {
		if isAllowed {
			e.reply.Apps = append(e.reply.Apps, app)
		}
	}

	return nil
}

func (e *ListAllowedAppsExecutor) Authorize() error {
	return nil
}

func GetAllowedApps(auth *ManagerAuthArg, user string) map[string]bool {
	result := map[string]bool{}
	filterStr := "(&(objectClass=" + aldap.TeamClass + ")(" + aldap.UsernameAttr + "=" + aldap.UserCommonName + "=" + user +
		"," + aldap.UserOu + "))"
	sr, err := NewSearchReq(filterStr, []string{aldap.TeamCommonName}, auth)
	if err != nil {
		return result
	}
	for i := 0; i < len(sr.Entries); i++ {
		filterStr := "(&(objectClass=" + aldap.AppClass + ")(" + aldap.TeamCommonName + ":dn:=" +
			sr.Entries[i].GetAttributeValues(aldap.TeamCommonName)[0] + "))"
		ss, err := NewSearchReq(filterStr, []string{aldap.AllowedAppAttr}, auth)
		if err != nil {
			return result
		}
		appCount := len(ss.Entries)
		for j := 0; j < appCount; j++ {
			app := ss.Entries[j].GetAttributeValues(aldap.AllowedAppAttr)[0]
			result[app] = true
		}
	}
	return result
}

func IsAppAllowed(auth *ManagerAuthArg, user string, app string) bool {
	var suReply ManagerSuperUserReply
	if err := NewTask("IsAppAllowed-IsSuperUser",
		&IsSuperUserExecutor{ManagerSuperUserArg{*auth}, &suReply}).Run(); err != nil {
		return false
	}
	if suReply.IsSuperUser {
		// shortcut for superusers
		return true
	}
	return GetAllowedApps(auth, user)[app]
}

// ----------------------------------------------------------------------------------------------------------
// Permissions
// ----------------------------------------------------------------------------------------------------------

type IsTeamAdminExecutor struct {
	arg   ManagerTeamAdminArg
	reply *ManagerTeamAdminReply
}

func (e *IsTeamAdminExecutor) Request() interface{} {
	return e.arg
}

func (e *IsTeamAdminExecutor) Result() interface{} {
	return e.reply
}

func (e *IsTeamAdminExecutor) Description() string {
	return fmt.Sprintf("%s : %s", e.arg.User, e.arg.Team)
}

func (e *IsTeamAdminExecutor) Execute(t *Task) error {
	if aldap.SkipAuthorization {
		e.reply.IsAdmin = true
		return nil
	}

	var suReply ManagerSuperUserReply
	if err := NewTask("IsTeamAdmin-IsSuperUser", &IsSuperUserExecutor{ManagerSuperUserArg{e.arg.ManagerAuthArg},
		&suReply}).Run(); err == nil {
		e.reply.IsAdmin = suReply.IsSuperUser
		if e.reply.IsAdmin {
			return nil
		}
	} else {
		return err
	}

	if !TeamExists(e.arg.Team, &e.arg.ManagerAuthArg) {
		e.reply.IsAdmin = false
		return errors.New("Team Does Not Exist")
	}

	filterStr := "(&(objectClass=" + aldap.TeamClass + ")(" + aldap.TeamCommonName + "=" + e.arg.Team + "))"
	sr, err := NewSearchReq(filterStr, []string{aldap.TeamAdminAttr}, &e.arg.ManagerAuthArg)
	if err != nil {
		e.reply.IsAdmin = false
		return err
	} else if len(sr.Entries) == 0 {
		e.reply.IsAdmin = false
		return errors.New("Could not list team admin attribute")
	}
	e.reply.IsAdmin = ProcessTeamAdmin(aldap.UserCommonName+"="+e.arg.User+","+aldap.UserOu, sr)
	return nil
}

func ProcessTeamAdmin(userdn string, sr *ldap.SearchResult) bool {
	srTeamAdmin := sr.Entries[0].GetAttributeValues(aldap.TeamAdminAttr)
	teamAdminCount := len(srTeamAdmin)
	for i := 0; i < teamAdminCount; i++ {
		teamLead := sr.Entries[0].GetAttributeValues(aldap.TeamAdminAttr)[i]
		if teamLead == userdn {
			return true
		}
	}
	return false
}

func (e *IsTeamAdminExecutor) Authorize() error {
	return nil
}

type IsSuperUserExecutor struct {
	arg   ManagerSuperUserArg
	reply *ManagerSuperUserReply
}

func (e *IsSuperUserExecutor) Request() interface{} {
	return e.arg
}

func (e *IsSuperUserExecutor) Result() interface{} {
	return e.reply
}

func (e *IsSuperUserExecutor) Description() string {
	return fmt.Sprintf("%s", e.arg.User)
}

func (e *IsSuperUserExecutor) Execute(t *Task) error {
	if aldap.SkipAuthorization {
		e.reply.IsSuperUser = true
		return nil
	}
	if !UserExists(e.arg.User, &e.arg.ManagerAuthArg) {
		e.reply.IsSuperUser = false
		return nil
	}
	filterStr := "(&(objectClass=" + aldap.TeamClass + ")(" + aldap.SuperUserGroup + ")(" + aldap.UsernameAttr + "=" + aldap.UserCommonName + "=" + e.arg.User + "," + aldap.UserOu + "))"
	sr, err := NewSearchReq(filterStr, []string{aldap.TeamCommonName}, &e.arg.ManagerAuthArg)
	if err != nil {
		return err
	}
	if len(sr.Entries) > 0 {
		e.reply.IsSuperUser = true
	} else {
		e.reply.IsSuperUser = false
	}
	t.Log("-> %t", e.reply.IsSuperUser)
	return nil
}

func (e *IsSuperUserExecutor) Authorize() error {
	return nil
}

func (m *ManagerRPC) IsAppAllowed(arg ManagerIsAppAllowedArg, reply *ManagerIsAppAllowedReply) error {
	return NewTask("IsAppAllowed", &IsAppAllowedExecutor{arg, reply}).Run()
}

func (m *ManagerRPC) ListAllowedApps(arg ManagerListAllowedAppsArg, reply *ManagerListAllowedAppsReply) error {
	return NewTask("ListAllowedApps", &ListAllowedAppsExecutor{arg, reply}).Run()
}

func (m *ManagerRPC) IsTeamAdmin(arg ManagerTeamAdminArg, reply *ManagerTeamAdminReply) error {
	return NewTask("IsTeamAdmin", &IsTeamAdminExecutor{arg, reply}).Run()
}

func (m *ManagerRPC) IsSuperUser(arg ManagerSuperUserArg, reply *ManagerSuperUserReply) error {
	return NewTask("IsSuperUser", &IsSuperUserExecutor{arg, reply}).Run()
}

// ----------------------------------------------------------------------------------------------------------
// Helpers
// ----------------------------------------------------------------------------------------------------------

func TeamExists(name string, auth *ManagerAuthArg) bool {
	filterStr := "(&(objectClass=" + aldap.TeamClass + ")(" + aldap.TeamCommonName + "=" + name + "))"
	sr, err := NewSearchReq(filterStr, []string{aldap.TeamCommonName}, auth)
	if err != nil || len(sr.Entries) == 0 {
		return false
	}
	return true
}

func ListTeams(auth *ManagerAuthArg) ([]string, error) {
	filterStr := "(&(objectClass=" + aldap.TeamClass + ")"
	sr, err := NewSearchReq(filterStr, []string{aldap.TeamCommonName}, auth)
	ret := []string{}
	if err != nil || sr == nil {
		return ret, err
	}
	for _, entry := range sr.Entries {
		vals := entry.GetAttributeValues(aldap.TeamCommonName)
		if len(vals) > 0 {
			ret = append(ret, vals...)
		}
	}
	return ret, nil
}

func ListTeamAttributes(team, attribute string, auth *ManagerAuthArg) ([]string, error) {
	filterStr := "(&(objectClass=" + aldap.TeamClass + ")(" + aldap.TeamCommonName + "=" + team + "))"
	sr, err := NewSearchReq(filterStr, []string{attribute}, auth)
	ret := []string{}
	if err != nil || sr == nil {
		return ret, err
	}
	for _, entry := range sr.Entries {
		vals := entry.GetAttributeValues(attribute)
		if len(vals) > 0 {
			ret = append(ret, vals...)
		}
	}
	return ret, nil
}

func ListTeamEmails(team string, auth *ManagerAuthArg) ([]string, error) {
	return ListTeamAttributes(team, "email", auth)
}

func ListTeamAdmins(team string, auth *ManagerAuthArg) ([]string, error) {
	ret, err := ListTeamAttributes(team, aldap.TeamAdminAttr, auth)
	for i, val := range ret {
		ret[i] = ExtractUsername(val)
	}
	return ret, err
}

func ListTeamMembers(team string, auth *ManagerAuthArg) ([]string, error) {
	ret, err := ListTeamAttributes(team, aldap.UsernameAttr, auth)
	for i, val := range ret {
		ret[i] = ExtractUsername(val)
	}
	return ret, err
}

func ListTeamApps(team string, auth *ManagerAuthArg) ([]string, error) {
	result := []string{}
	filterStr := "(&(objectClass=" + aldap.AppClass + ")(" + aldap.TeamCommonName + ":dn:=" +
		team + "))"
	ss, err := NewSearchReq(filterStr, []string{aldap.AllowedAppAttr}, auth)
	if err != nil {
		return result, err
	}
	appCount := len(ss.Entries)
	for j := 0; j < appCount; j++ {
		app := ss.Entries[j].GetAttributeValues(aldap.AllowedAppAttr)[0]
		result = append(result, app)
	}
	return result, nil
}

func UserExists(name string, auth *ManagerAuthArg) bool {
	filterStr := "(&(objectClass=" + aldap.UserClass + ")(" + aldap.UserClassAttr + "=" + name + "))"
	sr, err := NewSearchReq(filterStr, []string{aldap.UserClassAttr}, auth)
	if err != nil || len(sr.Entries) == 0 {
		return false
	}
	return true
}

func EmailExists(email string, team string, auth *ManagerAuthArg) bool {
	filterStr := "(&(objectClass=" + aldap.TeamClass + ")(" + aldap.TeamCommonName + "=" + team + ")(email=" + email + "))"
	sr, err := NewSearchReq(filterStr, []string{"email"}, auth)
	if err != nil || len(sr.Entries) == 0 {
		return false
	}
	return true
}

func ExtractUsername(fullname string) string {
	return strings.Replace(strings.Replace(fullname, ","+aldap.UserOu, "", 1), aldap.UserCommonName+"=", "", 1)
}

// ----------------------------------------------------------------------------------------------------------
// Init
// ----------------------------------------------------------------------------------------------------------

func InitConnection(auth *ManagerAuthArg) (*ldap.LDAPConnection, error) {
	if conn := aldap.LookupConnection(auth.User, auth.Secret); conn != nil {
		return conn, nil
	}
	return nil, errors.New("No connection found.")
}

// ----------------------------------------------------------------------------------------------------------
// Search Requests
// ----------------------------------------------------------------------------------------------------------

func NewSearchReq(filter string, attributes []string, auth *ManagerAuthArg) (*ldap.SearchResult, error) {
	// 2 => Searching the Whole Subtree
	conn, err := InitConnection(auth)
	if err != nil {
		return nil, err
	}
	searchReq := ldap.NewSimpleSearchRequest(aldap.BaseDomain, 2, filter, attributes)
	sr, err := conn.Search(searchReq)
	return sr, err
}
