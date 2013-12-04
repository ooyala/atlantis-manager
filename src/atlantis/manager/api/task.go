package api

import (
	. "atlantis/common"
	. "atlantis/manager/rpc/types"
	"fmt"
	"net/http"
)

func Wait(w http.ResponseWriter, r *http.Request) {
	var statusReply TaskStatus
	err := manager.Status(r.FormValue("Id"), &statusReply)
	output := map[string]interface{}{
		"Name":        statusReply.Name,
		"Status":      statusReply.Status,
		"Description": statusReply.Description,
		"Done":        statusReply.Done,
	}
	if !statusReply.Done {
		fmt.Fprintf(w, "%s", Output(output, err))
		return
	}
	if statusReply.Name == "Deploy" {
		var reply ManagerDeployReply
		err = manager.DeployResult(r.FormValue("Id"), &reply)
		output["Containers"] = reply.Containers
	} else if statusReply.Name == "Teardown" {
		var reply ManagerTeardownReply
		err = manager.TeardownResult(r.FormValue("Id"), &reply)
		output["Containers"] = reply.ContainerIds
	} else if statusReply.Name == "RegisterRouter" {
		var reply ManagerRegisterRouterReply
		err = manager.RegisterRouterResult(r.FormValue("Id"), &reply)
		output["Router"] = reply.Router
	} else if statusReply.Name == "UnregisterRouter" {
		var reply ManagerRegisterRouterReply
		err = manager.UnregisterRouterResult(r.FormValue("Id"), &reply)
	} else if statusReply.Name == "RegisterManager" {
		var reply ManagerRegisterManagerReply
		err = manager.RegisterManagerResult(r.FormValue("Id"), &reply)
		output["Manager"] = reply.Manager
	} else if statusReply.Name == "UnregisterManager" {
		var reply ManagerRegisterManagerReply
		err = manager.UnregisterManagerResult(r.FormValue("Id"), &reply)
	}

	fmt.Fprintf(w, "%s", Output(output, err))
}
