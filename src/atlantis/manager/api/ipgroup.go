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
	"strings"
)

func UpdateIPGroup(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	ipsString := r.FormValue("IPs")
	ips := strings.Split(ipsString, ",")
	arg := ManagerUpdateIPGroupArg{
		ManagerAuthArg: auth,
		Name:           vars["Name"],
		IPs:            ips,
	}
	var reply ManagerUpdateIPGroupReply
	err = manager.UpdateIPGroup(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": reply.Status}, err))
}

func DeleteIPGroup(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerDeleteIPGroupArg{
		ManagerAuthArg: auth,
		Name:           vars["Name"],
	}
	var reply ManagerDeleteIPGroupReply
	err = manager.DeleteIPGroup(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": reply.Status}, err))
}

func GetIPGroup(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerGetIPGroupArg{
		ManagerAuthArg: auth,
		Name:           vars["Name"],
	}
	var reply ManagerGetIPGroupReply
	err = manager.GetIPGroup(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": reply.Status, "IPGroup": reply.IPGroup}, err))
}
