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

func ContainerIDGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	cArg := ManagerGetContainerArg{auth, vars["ID"]}
	var reply ManagerGetContainerReply
	err := manager.GetContainer(cArg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Container": reply.Container, "Status": reply.Status}, err))
}

func ContainerHealthzGet(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("Host") == "" || r.FormValue("Port") == "" {
		fmt.Fprintf(w, "%s", "No Params Entered")
		return
	}
	host := "http://" + r.FormValue("Host")
	port := r.FormValue("Port")
	resp, err := http.Get(host + ":" + port + "/healthz")
	serverStatus := ""
	if err == nil {
		serverStatus = resp.Header.Get("Server-Status")
		if serverStatus == "" {
			serverStatus = "Unknown"
		}
	}
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": serverStatus}, err))
}

func ListApps(w http.ResponseWriter, r *http.Request) {
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerListAppsArg{auth}
	var reply ManagerListAppsReply
	err := manager.ListApps(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Apps": reply.Apps, "Status": reply.Status}, err))
}

func ListShas(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerListShasArg{auth, vars["App"]}
	var reply ManagerListShasReply
	err := manager.ListShas(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Shas": reply.Shas, "Status": reply.Status}, err))
}

func DeployListEnvs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerListEnvsArg{auth, vars["App"], vars["Sha"]}
	var reply ManagerListEnvsReply
	err := manager.ListEnvs(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Envs": reply.Envs, "Status": reply.Status}, err))
}

func ListContainers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	cArg := ManagerListContainersArg{auth, vars["App"], vars["Sha"], vars["Env"]}
	var reply ManagerListContainersReply
	err := manager.ListContainers(cArg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"ContainerIDs": reply.ContainerIDs, "Status": reply.Status}, err))
}
