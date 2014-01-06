package api

import (
	. "atlantis/common"
	. "atlantis/manager/rpc/types"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func Deploy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	cpushares, err := strconv.ParseUint(r.FormValue("CPUShares"), 10, 0)
	if err != nil {
		fmt.Fprintf(w, "{\"error\": \"%s\"}", err.Error())
		return
	}
	memlimit, err := strconv.ParseUint(r.FormValue("MemoryLimit"), 10, 0)
	if err != nil {
		fmt.Fprintf(w, "{\"error\": \"%s\"}", err.Error())
		return
	}
	instances, err := strconv.ParseUint(r.FormValue("Instances"), 10, 0)
	if err != nil {
		fmt.Fprintf(w, "{\"error\": \"%s\"}", err.Error())
		return
	}
	dev, err := strconv.ParseBool(r.FormValue("Dev"))
	if err != nil {
		fmt.Fprintf(w, "{\"error\": \"%s\"}", err.Error())
		return
	}
	dArg := ManagerDeployArg{
		ManagerAuthArg: auth,
		App:            vars["App"],
		Sha:            vars["Sha"],
		Env:            vars["Env"],
		Instances:      uint(instances),
		CPUShares:      uint(cpushares),
		MemoryLimit:    uint(memlimit),
		Dev:            bool(dev),
	}
	var reply AsyncReply
	err = manager.Deploy(dArg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"ID": reply.ID}, err))
}

func CopyContainer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	instances, err := strconv.ParseUint(r.FormValue("Instances"), 10, 0)
	if err != nil {
		fmt.Fprintf(w, "{\"error\": \"%s\"}", err.Error())
		return
	}
	ccArg := ManagerCopyContainerArg{ManagerAuthArg: auth, Instances: uint(instances),
		ContainerID: vars["ID"]}
	var reply AsyncReply
	err = manager.CopyContainer(ccArg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"ID": reply.ID}, err))
}

func MoveContainer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	ccArg := ManagerMoveContainerArg{ManagerAuthArg: auth, ContainerID: vars["ID"]}
	var reply AsyncReply
	err := manager.MoveContainer(ccArg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"ID": reply.ID}, err))
}

func CopyOrphaned(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	cleanup, err := strconv.ParseBool(r.FormValue("CleanupZk"))
	if err != nil {
		fmt.Fprintf(w, "{\"error\": \"%s\"}", err.Error())
		return
	}
	ccArg := ManagerCopyOrphanedArg{ManagerAuthArg: auth, ContainerID: vars["ID"], Host: r.FormValue("Host"),
		CleanupZk: cleanup}
	var reply AsyncReply
	err = manager.CopyOrphaned(ccArg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"ID": reply.ID}, err))
}

func Teardown(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerTeardownArg{auth, vars["App"], vars["Sha"], vars["Env"], "", false}
	var reply AsyncReply
	err := manager.Teardown(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"ID": reply.ID}, err))
}

func TeardownContainerID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	cArg := ManagerTeardownArg{auth, "", "", "", vars["ID"], false}
	var reply AsyncReply
	err := manager.Teardown(cArg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"ID": reply.ID}, err))
}

func TeardownContainers(w http.ResponseWriter, r *http.Request) {
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	all, err := strconv.ParseBool(r.FormValue("All"))
	if err != nil {
		fmt.Fprintf(w, "{\"error\": \"%s\"}", err.Error())
		return
	}
	tArg := ManagerTeardownArg{auth, r.FormValue("App"), r.FormValue("Sha"), r.FormValue("Env"), r.FormValue("ContainerID"), all}
	var reply AsyncReply
	err = manager.Teardown(tArg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"ID": reply.ID}, err))
}
