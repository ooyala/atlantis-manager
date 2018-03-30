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

func ListTeamAdmins(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"TeamAdmins": ""}, nil))
}


func ListTeamMembers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerListTeamMembersArg{auth, vars["Team"]}
	var reply ManagerListTeamMembersReply
	err := manager.ListTeamMembers(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"TeamMembers": reply.TeamMembers}, err))
}

func GetPermissions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{vars["User"], "", r.FormValue("Secret")}
	arg := ManagerSuperUserArg{auth}
	var reply ManagerSuperUserReply
	err := manager.IsSuperUser(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"SuperUser": reply.IsSuperUser}, err))
}
