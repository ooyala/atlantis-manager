package api

import (
	. "atlantis/manager/rpc/types"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

// TODO(edanaher): These functions are so similar...  what a waste of space.

func ListTeamApps(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerListTeamAppsArg{auth, vars["Team"]}
	var reply ManagerListTeamAppsReply
	err := manager.ListTeamApps(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"TeamApps": reply.TeamApps}, err))
}

func AllowApp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerAppArg{auth, vars["App"], vars["Team"]}
	var reply ManagerAppReply
	err := manager.AllowApp(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{}, err))
}

func DisallowApp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerAppArg{auth, vars["App"], vars["Team"]}
	var reply ManagerAppReply
	err := manager.DisallowApp(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{}, err))
}

func ListTeams(w http.ResponseWriter, r *http.Request) {
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerListTeamsArg{auth}
	var reply ManagerListTeamsReply
	err := manager.ListTeams(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Teams": reply.Teams}, err))
}

func CreateTeam(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerTeamArg{auth, vars["Team"]}
	var reply ManagerTeamReply
	err := manager.CreateTeam(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{}, err))
}

func DeleteTeam(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerTeamArg{auth, vars["Team"]}
	var reply ManagerTeamReply
	err := manager.DeleteTeam(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{}, err))
}

func ListTeamAdmins(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerListTeamAdminsArg{auth, vars["Team"]}
	var reply ManagerListTeamAdminsReply
	err := manager.ListTeamAdmins(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"TeamAdmins": reply.TeamAdmins}, err))
}

func AddTeamAdmin(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerModifyTeamAdminArg{auth, vars["Team"], vars["Admin"]}
	var reply ManagerModifyTeamAdminReply
	err := manager.AddTeamAdmin(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{}, err))
}

func RemoveTeamAdmin(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerModifyTeamAdminArg{auth, vars["Team"], vars["Admin"]}
	var reply ManagerModifyTeamAdminReply
	err := manager.RemoveTeamAdmin(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{}, err))
}

func ListTeamMembers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerListTeamMembersArg{auth, vars["Team"]}
	var reply ManagerListTeamMembersReply
	err := manager.ListTeamMembers(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"TeamMembers": reply.TeamMembers}, err))
}

func AddTeamMember(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerTeamMemberArg{auth, vars["Team"], vars["Member"]}
	var reply ManagerTeamMemberReply
	err := manager.AddTeamMember(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{}, err))
}

func RemoveTeamMember(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerTeamMemberArg{auth, vars["Team"], vars["Member"]}
	var reply ManagerTeamMemberReply
	err := manager.RemoveTeamMember(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{}, err))
}

func AddTeamEmail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerEmailArg{auth, vars["Team"], vars["Email"]}
	var reply ManagerEmailReply
	err := manager.AddTeamEmail(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{}, err))
}

func RemoveTeamEmail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerEmailArg{auth, vars["Team"], vars["Email"]}
	var reply ManagerEmailReply
	err := manager.RemoveTeamEmail(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{}, err))
}

func GetPermissions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{vars["User"], "", r.FormValue("Secret")}
	arg := ManagerSuperUserArg{auth}
	var reply ManagerSuperUserReply
	err := manager.IsSuperUser(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"SuperUser": reply.IsSuperUser}, err))
}
