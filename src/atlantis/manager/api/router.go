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
	"atlantis/router/config"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
)

func GetAppEnvPort(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	pArg := ManagerGetAppEnvPortArg{
		ManagerAuthArg: auth,
		App:            vars["App"],
		Env:            vars["Env"],
	}
	var reply ManagerGetAppEnvPortReply
	err := manager.GetAppEnvPort(pArg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Port": reply.Port, "Status": reply.Status}, err))
}

func ListAppEnvsWithPort(w http.ResponseWriter, r *http.Request) {
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	internal, err := strconv.ParseBool(r.FormValue("Internal"))
	if err != nil {
		fmt.Fprintf(w, "%s", Output(map[string]interface{}{}, err))
		return
	}
	pArg := ManagerListAppEnvsWithPortArg{
		ManagerAuthArg: auth,
		Internal:       internal,
	}
	var reply ManagerListAppEnvsWithPortReply
	err = manager.ListAppEnvsWithPort(pArg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"AppEnvs": reply.AppEnvs, "Status": reply.Status}, err))
}

func UpdatePort(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	internal, err := strconv.ParseBool(r.FormValue("Internal"))
	if err != nil {
		fmt.Fprintf(w, "%s", Output(map[string]interface{}{}, err))
		return
	}
	port, err := strconv.ParseUint(vars["Port"], 10, 16)
	if err != nil {
		fmt.Fprintf(w, "%s", Output(map[string]interface{}{}, err))
		return
	}
	pArg := ManagerUpdatePortArg{
		ManagerAuthArg: auth,
		Port: config.Port{
			Port:     uint16(port),
			Trie:     r.FormValue("Trie"),
			Internal: internal,
		},
	}
	var reply ManagerUpdatePortReply
	err = manager.UpdatePort(pArg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": reply.Status}, err))
}

func DeletePort(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	internal, err := strconv.ParseBool(r.FormValue("Internal"))
	if err != nil {
		fmt.Fprintf(w, "%s", Output(map[string]interface{}{}, err))
		return
	}
	port, err := strconv.ParseUint(vars["Port"], 10, 16)
	if err != nil {
		fmt.Fprintf(w, "%s", Output(map[string]interface{}{}, err))
		return
	}
	pArg := ManagerDeletePortArg{
		ManagerAuthArg: auth,
		Port:           uint16(port),
		Internal:       internal,
	}
	var reply ManagerDeletePortReply
	err = manager.DeletePort(pArg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": reply.Status}, err))
}

func GetPort(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	internal, err := strconv.ParseBool(r.FormValue("Internal"))
	if err != nil {
		fmt.Fprintf(w, "%s", Output(map[string]interface{}{}, err))
		return
	}
	port, err := strconv.ParseUint(vars["Port"], 10, 16)
	if err != nil {
		fmt.Fprintf(w, "%s", Output(map[string]interface{}{}, err))
		return
	}
	pArg := ManagerGetPortArg{
		ManagerAuthArg: auth,
		Port:           uint16(port),
		Internal:       internal,
	}
	var reply ManagerGetPortReply
	err = manager.GetPort(pArg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Port": reply.Port, "Status": reply.Status}, err))
}

func ListPorts(w http.ResponseWriter, r *http.Request) {
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	internal, err := strconv.ParseBool(r.FormValue("Internal"))
	if err != nil {
		fmt.Fprintf(w, "%s", Output(map[string]interface{}{}, err))
		return
	}
	arg := ManagerListPortsArg{auth, internal}
	var reply ManagerListPortsReply
	err = manager.ListPorts(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Ports": reply.Ports, "Status": reply.Status}, err))
}

func GetPool(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	internal, err := strconv.ParseBool(r.FormValue("Internal"))
	if err != nil {
		fmt.Fprintf(w, "%s", Output(map[string]interface{}{}, err))
		return
	}
	pArg := ManagerGetPoolArg{auth, vars["PoolName"], internal}
	var reply ManagerGetPoolReply
	err = manager.GetPool(pArg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Pool": reply.Pool, "Status": reply.Status}, err))
}

func UpdatePool(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	internal, err := strconv.ParseBool(r.FormValue("Internal"))
	if err != nil {
		fmt.Fprintf(w, "%s", Output(map[string]interface{}{}, err))
		return
	}
	hosts := map[string]config.Host{}
	for _, host := range strings.Split(r.FormValue("Hosts"), ", ") {
		hosts[host] = config.Host{Address: host}
	}
	pArg := ManagerUpdatePoolArg{auth, config.Pool{Name: vars["PoolName"],
		Config: config.PoolConfig{HealthzEvery: r.FormValue("HealthCheckEvery"),
			HealthzTimeout: r.FormValue("HealthzTimeout"), RequestTimeout: r.FormValue("RequestTimeout")},
		Hosts: hosts, Internal: internal}}
	var reply ManagerUpdatePoolReply
	err = manager.UpdatePool(pArg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": reply.Status}, err))
}

func DeletePool(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	internal, err := strconv.ParseBool(r.FormValue("Internal"))
	if err != nil {
		fmt.Fprintf(w, "%s", Output(map[string]interface{}{}, err))
		return
	}
	pArg := ManagerDeletePoolArg{auth, vars["PoolName"], internal}
	var reply ManagerDeletePoolReply
	err = manager.DeletePool(pArg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": reply.Status}, err))
}

func ListPools(w http.ResponseWriter, r *http.Request) {
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	internal, err := strconv.ParseBool(r.FormValue("Internal"))
	if err != nil {
		fmt.Fprintf(w, "%s", Output(map[string]interface{}{}, err))
		return
	}
	arg := ManagerListPoolsArg{auth, internal}
	var reply ManagerListPoolsReply
	err = manager.ListPools(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Pools": reply.Pools, "Status": reply.Status}, err))
}

func GetRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	internal, err := strconv.ParseBool(r.FormValue("Internal"))
	if err != nil {
		fmt.Fprintf(w, "%s", Output(map[string]interface{}{}, err))
		return
	}
	rArg := ManagerGetRuleArg{auth, vars["RuleName"], internal}
	var reply ManagerGetRuleReply
	err = manager.GetRule(rArg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Rule": reply.Rule, "Status": reply.Status}, err))
}

func UpdateRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	internal, err := strconv.ParseBool(r.FormValue("Internal"))
	if err != nil {
		fmt.Fprintf(w, "%s", Output(map[string]interface{}{}, err))
		return
	}
	rArg := ManagerUpdateRuleArg{auth, config.Rule{Name: vars["RuleName"], Type: r.FormValue("Type"),
		Value: r.FormValue("Value"), Next: r.FormValue("Next"), Pool: r.FormValue("Pool"), Internal: internal}}
	var reply ManagerUpdateRuleReply
	err = manager.UpdateRule(rArg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": reply.Status}, err))
}

func DeleteRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	internal, err := strconv.ParseBool(r.FormValue("Internal"))
	if err != nil {
		fmt.Fprintf(w, "%s", Output(map[string]interface{}{}, err))
		return
	}
	rArg := ManagerDeleteRuleArg{auth, vars["RuleName"], internal}
	var reply ManagerDeleteRuleReply
	err = manager.DeleteRule(rArg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": reply.Status}, err))
}

func ListRules(w http.ResponseWriter, r *http.Request) {
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	internal, err := strconv.ParseBool(r.FormValue("Internal"))
	if err != nil {
		fmt.Fprintf(w, "%s", Output(map[string]interface{}{}, err))
		return
	}
	arg := ManagerListRulesArg{auth, internal}
	var reply ManagerListRulesReply
	err = manager.ListRules(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Rules": reply.Rules, "Status": reply.Status}, err))
}

func GetTrie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	internal, err := strconv.ParseBool(r.FormValue("Internal"))
	if err != nil {
		fmt.Fprintf(w, "%s", Output(map[string]interface{}{}, err))
		return
	}
	tArg := ManagerGetTrieArg{auth, vars["TrieName"], internal}
	var reply ManagerGetTrieReply
	err = manager.GetTrie(tArg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Trie": reply.Trie, "Status": reply.Status}, err))
}

func UpdateTrie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	internal, err := strconv.ParseBool(r.FormValue("Internal"))
	if err != nil {
		fmt.Fprintf(w, "%s", Output(map[string]interface{}{}, err))
		return
	}
	var temp []string
	err = json.Unmarshal([]byte(r.Form["Rules"][0]), &temp)
	tArg := ManagerUpdateTrieArg{auth, config.Trie{Name: vars["TrieName"], Rules: temp, Internal: internal}}
	var reply ManagerUpdateTrieReply
	err = manager.UpdateTrie(tArg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": reply.Status}, err))
}

func DeleteTrie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	internal, err := strconv.ParseBool(r.FormValue("Internal"))
	if err != nil {
		fmt.Fprintf(w, "%s", Output(map[string]interface{}{}, err))
		return
	}
	tArg := ManagerDeleteTrieArg{auth, vars["TrieName"], internal}
	var reply ManagerDeleteTrieReply
	err = manager.DeleteTrie(tArg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": reply.Status}, err))
}

func ListTries(w http.ResponseWriter, r *http.Request) {
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	internal, err := strconv.ParseBool(r.FormValue("Internal"))
	if err != nil {
		fmt.Fprintf(w, "%s", Output(map[string]interface{}{}, err))
		return
	}
	arg := ManagerListTriesArg{auth, internal}
	var reply ManagerListTriesReply
	err = manager.ListTries(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Tries": reply.Tries, "Status": reply.Status}, err))
}
