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

func GetPool(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	r.ParseForm()
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	internal, err := strconv.ParseBool(r.FormValue("Internal"))
	pArg := ManagerGetPoolArg{auth, vars["PoolName"], internal}
	var reply ManagerGetPoolReply
	err = manager.GetPool(pArg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Pool": reply.Pool, "Status": reply.Status}, err))
}

func UpdatePool(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	internal, err := strconv.ParseBool(r.FormValue("Internal"))
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
	pArg := ManagerDeletePoolArg{auth, vars["PoolName"], internal}
	var reply ManagerDeletePoolReply
	err = manager.DeletePool(pArg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": reply.Status}, err))
}

func ListPools(w http.ResponseWriter, r *http.Request) {
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	internal, err := strconv.ParseBool(r.FormValue("Internal"))
	arg := ManagerListPoolsArg{auth, internal}
	var reply ManagerListPoolsReply
	err = manager.ListPools(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Pools": reply.Pools, "Status": reply.Status}, err))
}

func GetRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	r.ParseForm()
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	internal, err := strconv.ParseBool(r.FormValue("Internal"))
	rArg := ManagerGetRuleArg{auth, vars["RuleName"], internal}
	var reply ManagerGetRuleReply
	err = manager.GetRule(rArg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Rule": reply.Rule, "Status": reply.Status}, err))
}

func UpdateRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	internal, err := strconv.ParseBool(r.FormValue("Internal"))
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
	rArg := ManagerDeleteRuleArg{auth, vars["RuleName"], internal}
	var reply ManagerDeleteRuleReply
	err = manager.DeleteRule(rArg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": reply.Status}, err))
}

func ListRules(w http.ResponseWriter, r *http.Request) {
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	internal, err := strconv.ParseBool(r.FormValue("Internal"))
	arg := ManagerListRulesArg{auth, internal}
	var reply ManagerListRulesReply
	err = manager.ListRules(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Rules": reply.Rules, "Status": reply.Status}, err))
}

func GetTrie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	internal, err := strconv.ParseBool(r.FormValue("Internal"))
	tArg := ManagerGetTrieArg{auth, vars["TrieName"], internal}
	var reply ManagerGetTrieReply
	err = manager.GetTrie(tArg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Trie": reply.Trie, "Status": reply.Status}, err))
}

func UpdateTrie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	internal, err := strconv.ParseBool(r.FormValue("Internal"))
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
	tArg := ManagerDeleteTrieArg{auth, vars["TrieName"], internal}
	var reply ManagerDeleteTrieReply
	err = manager.DeleteTrie(tArg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": reply.Status}, err))
}

func ListTries(w http.ResponseWriter, r *http.Request) {
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	internal, err := strconv.ParseBool(r.FormValue("Internal"))
	arg := ManagerListTriesArg{auth, internal}
	var reply ManagerListTriesReply
	err = manager.ListTries(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Tries": reply.Tries, "Status": reply.Status}, err))
}
