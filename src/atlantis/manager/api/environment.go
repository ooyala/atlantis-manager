package api

import (
	. "atlantis/manager/rpc/types"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

func ListEnvs(w http.ResponseWriter, r *http.Request) {
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerListEnvsArg{auth, "", ""}
	var reply ManagerListEnvsReply
	err := manager.ListEnvs(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Envs": reply.Envs, "Status": reply.Status}, err))
}

func ResolveDeps(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	dnames := strings.Split(vars["DepNames"], ",")
	dArg := ManagerResolveDepsArg{ManagerAuthArg: auth, Env: vars["Env"], DepNames: dnames}
	var reply ManagerResolveDepsReply
	err := manager.ResolveDeps(dArg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Deps": reply.Deps, "Status": reply.Status}, err))
}

func GetDep(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	dArg := ManagerDepArg{auth, vars["Env"], vars["DepName"], ""}
	var reply ManagerDepReply
	err := manager.GetDep(dArg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Value": reply.Value, "Status": reply.Status}, err))
}

func UpdateDep(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	dArg := ManagerDepArg{auth, vars["Env"], vars["DepName"], r.FormValue("Value")}
	var reply ManagerDepReply
	err := manager.UpdateDep(dArg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Value": reply.Value, "Status": reply.Status}, err))
}

func DeleteDep(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	dArg := ManagerDepArg{auth, vars["Env"], vars["DepName"], r.FormValue("Value")}
	var reply ManagerDepReply
	err := manager.DeleteDep(dArg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Value": reply.Value, "Status": reply.Status}, err))
}

func GetEnv(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	dArg := ManagerEnvArg{auth, vars["Env"], ""}
	var reply ManagerEnvReply
	err := manager.GetEnv(dArg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Parent": reply.Parent, "Deps": reply.Deps,
		"ResolvedDeps": reply.ResolvedDeps, "Status": reply.Status}, err))
}

func UpdateEnv(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	dArg := ManagerEnvArg{auth, vars["Env"], r.FormValue("Parent")}
	var reply ManagerEnvReply
	err := manager.UpdateEnv(dArg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Parent": reply.Parent, "Deps": reply.Deps,
		"ResolvedDeps": reply.ResolvedDeps, "Status": reply.Status}, err))
}

func DeleteEnv(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	dArg := ManagerEnvArg{auth, vars["Env"], ""}
	var reply ManagerEnvReply
	err := manager.DeleteEnv(dArg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Parent": reply.Parent, "Deps": reply.Deps,
		"ResolvedDeps": reply.ResolvedDeps, "Status": reply.Status}, err))
}
