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
	. "atlantis/common"
	. "atlantis/manager/rpc/types"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
)

func ResolveDeps(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerResolveDepsArg{
		ManagerAuthArg: auth,
		App:            vars["App"],
		Env:            vars["Env"],
		DepNames:       strings.Split(vars["DepNames"], ","),
	}
	var reply ManagerResolveDepsReply
	err := manager.ResolveDeps(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Deps": reply.Deps}, err))
}

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
	skipBld, err := strconv.ParseBool(r.FormValue("SkipBuild"))
	if err != nil {
		//use false as default value 
		skipBld = false
	}
	manifest := r.FormValue("Manifest")
	
	dArg := ManagerDeployArg{
		ManagerAuthArg: auth,
		App:            vars["App"],
		Sha:            vars["Sha"],
		Env:            vars["Env"],
		Instances:      uint(instances),
		CPUShares:      uint(cpushares),
		MemoryLimit:    uint(memlimit),
		Dev:            bool(dev),
		SkipBuild:      bool(skipBld),
		Manifest:       manifest,
	}
	var reply AsyncReply
	err = manager.Deploy(dArg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"ID": reply.ID}, err))
}

func DeployContainer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	instances, err := strconv.ParseUint(r.FormValue("Instances"), 10, 0)
	if err != nil {
		fmt.Fprintf(w, "{\"error\": \"%s\"}", err.Error())
		return
	}
	ccArg := ManagerDeployContainerArg{ManagerAuthArg: auth, Instances: uint(instances),
		ContainerID: vars["ID"]}
	var reply AsyncReply
	err = manager.DeployContainer(ccArg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"ID": reply.ID}, err))
}

func CopyContainer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	ccArg := ManagerCopyContainerArg{
		ManagerAuthArg: auth,
		ContainerID:    vars["ID"],
		ToHost:         r.FormValue("ToHost"),
	}
	var reply AsyncReply
	err := manager.CopyContainer(ccArg, &reply)
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
