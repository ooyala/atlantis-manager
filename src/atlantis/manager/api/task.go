package api

import (
	. "atlantis/common"
	. "atlantis/manager/rpc/types"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func GetTaskStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var statusReply TaskStatus
	err := manager.Status(vars["ID"], &statusReply)
	output := map[string]interface{}{
		"Name":        statusReply.Name,
		"Status":      statusReply.Status,
		"Warnings":    statusReply.Warnings,
		"Description": statusReply.Description,
		"Done":        statusReply.Done,
	}
	if !statusReply.Done {
		fmt.Fprintf(w, "%s", Output(output, err))
		return
	}
	if statusReply.Name == "Deploy" {
		var reply ManagerDeployReply
		err = manager.DeployResult(vars["ID"], &reply)
		output["Containers"] = reply.Containers
	} else if statusReply.Name == "Teardown" {
		var reply ManagerTeardownReply
		err = manager.TeardownResult(vars["ID"], &reply)
		output["Containers"] = reply.ContainerIDs
	} else if statusReply.Name == "RegisterRouter" {
		var reply ManagerRegisterRouterReply
		err = manager.RegisterRouterResult(vars["ID"], &reply)
		output["Router"] = reply.Router
	} else if statusReply.Name == "UnregisterRouter" {
		var reply ManagerRegisterRouterReply
		err = manager.UnregisterRouterResult(vars["ID"], &reply)
	} else if statusReply.Name == "RegisterManager" {
		var reply ManagerRegisterManagerReply
		err = manager.RegisterManagerResult(vars["ID"], &reply)
		output["Manager"] = reply.Manager
	} else if statusReply.Name == "UnregisterManager" {
		var reply ManagerRegisterManagerReply
		err = manager.UnregisterManagerResult(vars["ID"], &reply)
	} else if statusReply.Name == "RegisterSupervisor" {
		var reply ManagerRegisterSupervisorReply
		err = manager.RegisterSupervisorResult(vars["ID"], &reply)
	} else if statusReply.Name == "UnregisterSupervisor" {
		var reply ManagerRegisterSupervisorReply
		err = manager.UnregisterSupervisorResult(vars["ID"], &reply)
	}

	fmt.Fprintf(w, "%s", Output(output, err))
}

func ListTaskIDs(w http.ResponseWriter, r *http.Request) {
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	var ids []string
	err := manager.ListTaskIDs(auth, &ids)
	output := map[string]interface{}{"IDs": ids}
	fmt.Fprintf(w, "%s", Output(output, err))
}
