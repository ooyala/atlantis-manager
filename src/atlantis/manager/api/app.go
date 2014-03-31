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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
)

func RequestAppDependency(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	envsString := r.FormValue("Envs")
	envs := strings.Split(envsString, ",")
	arg := ManagerRequestAppDependencyArg{
		ManagerAuthArg: auth,
		App:            vars["Depender"],
		Dependency:     vars["App"],
		Envs:           envs,
	}
	var reply ManagerRequestAppDependencyReply
	err = manager.RequestAppDependency(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": reply.Status}, err))
}

// ----------------------------------------------------------------------------------------------------------
// Depender App Data Methods
// ----------------------------------------------------------------------------------------------------------

func AddDependerAppData(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	depEnvData := map[string]*DependerEnvData{}
	if r.FormValue("DependerEnvData") != "" {
		err = json.Unmarshal([]byte(r.FormValue("DependerEnvData")), &depEnvData)
		if err != nil {
			fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": StatusError}, err))
			return
		}
	}
	arg := ManagerAddDependerAppDataArg{
		ManagerAuthArg: auth,
		App:            vars["App"],
		DependerAppData: &DependerAppData{
			Name:            vars["Depender"],
			DependerEnvData: depEnvData,
		},
	}
	var reply ManagerAddDependerAppDataReply
	err = manager.AddDependerAppData(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": reply.Status, "App": reply.App}, err))
}

func RemoveDependerAppData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerRemoveDependerAppDataArg{
		ManagerAuthArg: auth,
		App:            vars["App"],
		Depender:       vars["Depender"],
	}
	var reply ManagerRemoveDependerAppDataReply
	err := manager.RemoveDependerAppData(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": reply.Status, "App": reply.App}, err))
}

func GetDependerAppData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerGetDependerAppDataArg{
		ManagerAuthArg: auth,
		App:            vars["App"],
		Depender:       vars["Depender"],
	}
	var reply ManagerGetDependerAppDataReply
	err := manager.GetDependerAppData(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{
		"Status":          reply.Status,
		"DependerAppData": reply.DependerAppData,
	}, err))
}

// ----------------------------------------------------------------------------------------------------------
// Depender Env Data Methods
// ----------------------------------------------------------------------------------------------------------

func AddDependerEnvData(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	data := map[string]interface{}{}
	if r.FormValue("Data") != "" {
		err = json.Unmarshal([]byte(r.FormValue("Data")), &data)
		if err != nil {
			fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": StatusError}, err))
			return
		}
	}
	sgRaw := []string{}
	if r.FormValue("SecurityGroup") != "" {
		err = json.Unmarshal([]byte(r.FormValue("SecurityGroup")), &sgRaw)
		if err != nil {
			fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": StatusError}, err))
			return
		}
	}
	// convert []string -> map[string][]uint16
	sg := map[string][]uint16{}
	for _, groupAndPort := range sgRaw {
		parts := strings.SplitN(groupAndPort, ":", 2)
		if len(parts) != 2 {
			fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": StatusError},
				errors.New("Invalid Security Group entry: "+groupAndPort)))
			return
		}
		group := parts[0]
		port, err := strconv.ParseUint(parts[1], 10, 16)
		if err != nil {
			fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": StatusError},
				errors.New("Invalid Security Group entry: "+groupAndPort)))
			return
		}
		if existing, exists := sg[group]; !exists {
			sg[group] = []uint16{uint16(port)}
		} else {
			sg[group] = append(existing, uint16(port))
		}
	}
	arg := ManagerAddDependerEnvDataArg{
		ManagerAuthArg: auth,
		App:            vars["App"],
		DependerEnvData: &DependerEnvData{
			Name:          vars["Env"],
			SecurityGroup: sg,
			DataMap:       data,
		},
	}
	var reply ManagerAddDependerEnvDataReply
	err = manager.AddDependerEnvData(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": reply.Status, "App": reply.App}, err))
}

func RemoveDependerEnvData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerRemoveDependerEnvDataArg{
		ManagerAuthArg: auth,
		App:            vars["App"],
		Env:            vars["Env"],
	}
	var reply ManagerRemoveDependerEnvDataReply
	err := manager.RemoveDependerEnvData(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": reply.Status, "App": reply.App}, err))
}

func GetDependerEnvData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerGetDependerEnvDataArg{
		ManagerAuthArg: auth,
		App:            vars["App"],
		Env:            vars["Env"],
	}
	var reply ManagerGetDependerEnvDataReply
	err := manager.GetDependerEnvData(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{
		"Status":          reply.Status,
		"DependerEnvData": reply.DependerEnvData,
	}, err))
}

// ----------------------------------------------------------------------------------------------------------
// Depender Env Data For Depender App Methods
// ----------------------------------------------------------------------------------------------------------

func AddDependerEnvDataForDependerApp(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	data := map[string]interface{}{}
	if r.FormValue("Data") != "" {
		err = json.Unmarshal([]byte(r.FormValue("Data")), &data)
		if err != nil {
			fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": StatusError}, err))
			return
		}
	}
	sgRaw := []string{}
	if r.FormValue("SecurityGroup") != "" {
		err = json.Unmarshal([]byte(r.FormValue("SecurityGroup")), &sgRaw)
		if err != nil {
			fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": StatusError}, err))
			return
		}
	}
	// convert []string -> map[string][]uint16
	sg := map[string][]uint16{}
	for _, groupAndPort := range sgRaw {
		parts := strings.SplitN(groupAndPort, ":", 2)
		if len(parts) != 2 {
			fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": StatusError},
				errors.New("Invalid Security Group entry: "+groupAndPort)))
			return
		}
		group := parts[0]
		port, err := strconv.ParseUint(parts[1], 10, 16)
		if err != nil {
			fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": StatusError},
				errors.New("Invalid Security Group entry: "+groupAndPort)))
			return
		}
		if existing, exists := sg[group]; !exists {
			sg[group] = []uint16{uint16(port)}
		} else {
			sg[group] = append(existing, uint16(port))
		}
	}
	arg := ManagerAddDependerEnvDataForDependerAppArg{
		ManagerAuthArg: auth,
		App:            vars["App"],
		Depender:       vars["Depender"],
		DependerEnvData: &DependerEnvData{
			Name:          vars["Env"],
			SecurityGroup: sg,
			DataMap:       data,
		},
	}
	var reply ManagerAddDependerEnvDataForDependerAppReply
	err = manager.AddDependerEnvDataForDependerApp(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": reply.Status, "App": reply.App}, err))
}

func RemoveDependerEnvDataForDependerApp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerRemoveDependerEnvDataForDependerAppArg{
		ManagerAuthArg: auth,
		App:            vars["App"],
		Depender:       vars["Depender"],
		Env:            vars["Env"],
	}
	var reply ManagerRemoveDependerEnvDataForDependerAppReply
	err := manager.RemoveDependerEnvDataForDependerApp(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": reply.Status, "App": reply.App}, err))
}

func GetDependerEnvDataForDependerApp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	arg := ManagerGetDependerEnvDataForDependerAppArg{
		ManagerAuthArg: auth,
		App:            vars["App"],
		Depender:       vars["Depender"],
		Env:            vars["Env"],
	}
	var reply ManagerGetDependerEnvDataForDependerAppReply
	err := manager.GetDependerEnvDataForDependerApp(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{
		"Status":          reply.Status,
		"DependerEnvData": reply.DependerEnvData,
	}, err))
}
