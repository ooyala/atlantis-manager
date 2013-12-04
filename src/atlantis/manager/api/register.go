package api

import (
	. "atlantis/common"
	. "atlantis/manager/rpc/types"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func ListRouters(w http.ResponseWriter, r *http.Request) {
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	internal := false
	internalStr := r.FormValue("Internal")
	if internalStr != "" {
		internal, _ = strconv.ParseBool(internalStr)
	}
	arg := ManagerListRoutersArg{ManagerAuthArg: auth, Internal: internal}
	var reply ManagerListRoutersReply
	err := manager.ListRouters(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Routers": reply.Routers, "Status": reply.Status}, err))
}

func RegisterRouter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	internal := false
	internalStr := r.FormValue("Internal")
	if internalStr != "" {
		internal, _ = strconv.ParseBool(internalStr)
	}
	arg := ManagerRegisterRouterArg{ManagerAuthArg: auth, Internal: internal, Zone: vars["Zone"], IP: vars["IP"]}
	var reply AsyncReply
	err := manager.RegisterRouter(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Id": reply.Id}, err))
}

func UnregisterRouter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	internal := false
	internalStr := r.FormValue("Internal")
	if internalStr != "" {
		internal, _ = strconv.ParseBool(internalStr)
	}
	arg := ManagerRegisterRouterArg{ManagerAuthArg: auth, Internal: internal, Zone: vars["Zone"], IP: vars["IP"]}
	var reply AsyncReply
	err := manager.UnregisterRouter(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Id": reply.Id}, err))
}

func GetRouter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	internal := false
	internalStr := r.FormValue("Internal")
	if internalStr != "" {
		internal, _ = strconv.ParseBool(internalStr)
	}
	arg := ManagerGetRouterArg{ManagerAuthArg: auth, Internal: internal, Zone: vars["Zone"], IP: vars["IP"]}
	var reply ManagerGetRouterReply
	err := manager.GetRouter(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": reply.Status, "Router": reply.Router}, err))
}

func ListRegisteredApps(w http.ResponseWriter, r *http.Request) {
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerListRegisteredAppsArg{auth}
	var reply ManagerListRegisteredAppsReply
	err := manager.ListRegisteredApps(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Apps": reply.Apps, "Status": reply.Status}, err))
}

func RegisterApp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerRegisterAppArg{
		ManagerAuthArg: auth,
		Name:           vars["App"],
		Repo:           r.FormValue("Repo"),
		Root:           r.FormValue("Root"),
	}
	var reply ManagerRegisterAppReply
	err := manager.RegisterApp(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": reply.Status}, err))
}

func UnregisterApp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerRegisterAppArg{ManagerAuthArg: auth, Name: vars["App"]}
	var reply ManagerRegisterAppReply
	err := manager.UnregisterApp(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": reply.Status}, err))
}

func GetApp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerGetAppArg{ManagerAuthArg: auth, Name: vars["App"]}
	var reply ManagerGetAppReply
	err := manager.GetApp(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": reply.Status, "App": reply.App}, err))
}

func ListSupervisors(w http.ResponseWriter, r *http.Request) {
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerListSupervisorsArg{auth}
	var reply ManagerListSupervisorsReply
	err := manager.ListSupervisors(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Supervisors": reply.Supervisors, "Status": reply.Status}, err))
}

func RegisterSupervisor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerRegisterSupervisorArg{auth, vars["Host"]}
	var reply ManagerRegisterSupervisorReply
	err := manager.RegisterSupervisor(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": reply.Status}, err))
}

func UnregisterSupervisor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerRegisterSupervisorArg{auth, vars["Host"]}
	var reply ManagerRegisterSupervisorReply
	err := manager.UnregisterSupervisor(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": reply.Status}, err))
}

func ListManagers(w http.ResponseWriter, r *http.Request) {
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerListManagersArg{auth}
	var reply ManagerListManagersReply
	err := manager.ListManagers(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Managers": reply.Managers, "Status": reply.Status}, err))
}

func RegisterManager(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerRegisterManagerArg{
		ManagerAuthArg: auth,
		IP:             vars["IP"],
		Region:         vars["Region"],
		ManagerCName:   r.FormValue("ManagerCName"),
		RegistryCName:  r.FormValue("RegistryCName"),
	}
	var reply AsyncReply
	err := manager.RegisterManager(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Id": reply.Id}, err))
}

func UnregisterManager(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerRegisterManagerArg{ManagerAuthArg: auth, IP: vars["IP"], Region: vars["Region"]}
	var reply AsyncReply
	err := manager.UnregisterManager(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Id": reply.Id}, err))
}
