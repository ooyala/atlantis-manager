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
	internal, err := strconv.ParseBool(r.FormValue("Internal"))
	if err != nil {
		fmt.Fprintf(w, "%s", Output(map[string]interface{}{}, err))
		return
	}
	arg := ManagerListRoutersArg{ManagerAuthArg: auth, Internal: internal}
	var reply ManagerListRoutersReply
	err = manager.ListRouters(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Routers": reply.Routers, "Status": reply.Status}, err))
}

func RegisterRouter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	internal, err := strconv.ParseBool(r.FormValue("Internal"))
	if err != nil {
		fmt.Fprintf(w, "%s", Output(map[string]interface{}{}, err))
		return
	}
	arg := ManagerRegisterRouterArg{
		ManagerAuthArg: auth,
		Internal:       internal,
		Zone:           vars["Zone"],
		Host:           vars["Host"],
	}
	var reply AsyncReply
	err = manager.RegisterRouter(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"ID": reply.ID}, err))
}

func UnregisterRouter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	internal, err := strconv.ParseBool(r.FormValue("Internal"))
	if err != nil {
		fmt.Fprintf(w, "%s", Output(map[string]interface{}{}, err))
		return
	}
	arg := ManagerRegisterRouterArg{
		ManagerAuthArg: auth,
		Internal:       internal,
		Zone:           vars["Zone"],
		Host:           vars["Host"],
	}
	var reply AsyncReply
	err = manager.UnregisterRouter(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"ID": reply.ID}, err))
}

func GetRouter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	internal, err := strconv.ParseBool(r.FormValue("Internal"))
	if err != nil {
		fmt.Fprintf(w, "%s", Output(map[string]interface{}{}, err))
		return
	}
	arg := ManagerGetRouterArg{
		ManagerAuthArg: auth,
		Internal:       internal,
		Zone:           vars["Zone"],
		Host:           vars["Host"],
	}
	var reply ManagerGetRouterReply
	err = manager.GetRouter(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": reply.Status, "Router": reply.Router}, err))
}

func ListRegisteredApps(w http.ResponseWriter, r *http.Request) {
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	authorizedOnly, _ := strconv.ParseBool(r.FormValue("AuthorizedOnly"))
	if authorizedOnly {
		arg := ManagerListRegisteredAppsArg{auth}
		var reply ManagerListRegisteredAppsReply
		err := manager.ListAuthorizedRegisteredApps(arg, &reply)
		fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Apps": reply.Apps, "Status": reply.Status}, err))
	} else {
		arg := ManagerListRegisteredAppsArg{auth}
		var reply ManagerListRegisteredAppsReply
		err := manager.ListRegisteredApps(arg, &reply)
		fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Apps": reply.Apps, "Status": reply.Status}, err))
	}
}

func RegisterApp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerRegisterAppArg{
		ManagerAuthArg: auth,
		Name:           vars["App"],
		Repo:           r.FormValue("Repo"),
		Root:           r.FormValue("Root"),
		Email:          r.FormValue("Email"),
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

func AddDependerApp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerDependerAppArg{
		ManagerAuthArg: auth,
		Dependee:       vars["Dependee"],
		Depender:       vars["Depender"],
	}
	var reply ManagerDependerAppReply
	err := manager.AddDependerApp(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": reply.Status, "Dependee": reply.Dependee}, err))
}

func RemoveDependerApp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerDependerAppArg{
		ManagerAuthArg: auth,
		Dependee:       vars["Dependee"],
		Depender:       vars["Depender"],
	}
	var reply ManagerDependerAppReply
	err := manager.RemoveDependerApp(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": reply.Status, "Dependee": reply.Dependee}, err))
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
		Host:           vars["Host"],
		Region:         vars["Region"],
		ManagerCName:   r.FormValue("ManagerCName"),
		RegistryCName:  r.FormValue("RegistryCName"),
	}
	var reply AsyncReply
	err := manager.RegisterManager(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"ID": reply.ID}, err))
}

func UnregisterManager(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerRegisterManagerArg{ManagerAuthArg: auth, Host: vars["Host"], Region: vars["Region"]}
	var reply AsyncReply
	err := manager.UnregisterManager(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"ID": reply.ID}, err))
}

func GetManager(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerGetManagerArg{
		ManagerAuthArg: auth,
		Region:         vars["Region"],
		Host:           vars["Host"],
	}
	var reply ManagerGetManagerReply
	err := manager.GetManager(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": reply.Status, "Manager": reply.Manager}, err))
}

func GetSelf(w http.ResponseWriter, r *http.Request) {
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerGetSelfArg{ManagerAuthArg: auth}
	var reply ManagerGetManagerReply
	err := manager.GetSelf(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": reply.Status, "Manager": reply.Manager}, err))
}

func AddRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerRoleArg{
		ManagerAuthArg: auth,
		Host:           vars["Host"],
		Region:         vars["Region"],
		Role:           vars["Role"],
	}
	var reply ManagerRoleReply
	err := manager.AddRole(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": reply.Status, "Manager": reply.Manager}, err))
}

func RemoveRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerRoleArg{
		ManagerAuthArg: auth,
		Host:           vars["Host"],
		Region:         vars["Region"],
		Role:           vars["Role"],
	}
	var reply ManagerRoleReply
	err := manager.RemoveRole(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": reply.Status, "Manager": reply.Manager}, err))
}

func AddRoleType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerRoleArg{
		ManagerAuthArg: auth,
		Host:           vars["Host"],
		Region:         vars["Region"],
		Role:           vars["Role"],
		Type:           vars["Type"],
	}
	var reply ManagerRoleReply
	err := manager.AddRole(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": reply.Status, "Manager": reply.Manager}, err))
}

func RemoveRoleType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerRoleArg{
		ManagerAuthArg: auth,
		Host:           vars["Host"],
		Region:         vars["Region"],
		Role:           vars["Role"],
		Type:           vars["Type"],
	}
	var reply ManagerRoleReply
	err := manager.RemoveRole(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": reply.Status, "Manager": reply.Manager}, err))
}
