package api

import (
	. "atlantis/manager/rpc/types"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func ListEnvs(w http.ResponseWriter, r *http.Request) {
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerListEnvsArg{auth, "", ""}
	var reply ManagerListEnvsReply
	err := manager.ListEnvs(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Envs": reply.Envs, "Status": reply.Status}, err))
}

func UpdateEnv(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	dArg := ManagerEnvArg{auth, vars["Env"]}
	var reply ManagerEnvReply
	err := manager.UpdateEnv(dArg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": reply.Status}, err))
}

func DeleteEnv(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	dArg := ManagerEnvArg{auth, vars["Env"]}
	var reply ManagerEnvReply
	err := manager.DeleteEnv(dArg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": reply.Status}, err))
}
