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

package client

import (
	atlantis "atlantis/common"
	. "atlantis/manager/rpc/types"
	"fmt"
)

type DeployCommand struct {
	App         string `short:"a" long:"app" description:"the app to deploy"`
	Sha         string `short:"s" long:"sha" description:"the sha to deploy"`
	Env         string `short:"e" long:"env" description:"the environment to deploy"`
	Instances   uint   `short:"i" long:"instances" default:"1" description:"the number of instances to deploy in each AZ"`
	CPUShares   uint   `short:"c" long:"cpu-shares" default:"0" description:"the number of CPU shares per instance"`
	MemoryLimit uint   `short:"m" long:"memory-limit" default:"0" description:"the MBytes of memory per instance"`
	Dev         bool   `long:"dev" description:"only deploy 1 instance in 1 AZ"`
	Wait        bool   `long:"wait" description:"wait until the deploy is done before exiting"`
}

func (c *DeployCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	Log("Deploy...")
	arg := ManagerDeployArg{
		ManagerAuthArg: dummyAuthArg,
		App:            c.App,
		Sha:            c.Sha,
		Env:            c.Env,
		Instances:      c.Instances,
		CPUShares:      c.CPUShares,
		MemoryLimit:    c.MemoryLimit,
		Dev:            c.Dev,
	}
	var reply atlantis.AsyncReply
	if err := rpcClient.CallAuthed("Deploy", &arg, &reply); err != nil {
		return OutputError(err)
	}
	Log("-> ID: %s", reply.ID)
	if !c.Wait {
		return Output(map[string]interface{}{"id": reply.ID}, reply.ID, nil)
	}
	return (&WaitCommand{reply.ID}).Execute(args)
}

type DeployContainerCommand struct {
	ContainerID string `short:"c" long:"container" description:"the id of the container to replicate"`
	Instances   uint   `short:"i" long:"instances" default:"1" description:"the number of instances to deploy in each AZ"`
	Wait        bool   `long:"wait" description:"wait until the deploy is done before exiting"`
}

func (c *DeployContainerCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	Log("DeployContainer...")
	arg := ManagerDeployContainerArg{ManagerAuthArg: dummyAuthArg, ContainerID: c.ContainerID, Instances: c.Instances}
	var reply atlantis.AsyncReply
	if err := rpcClient.CallAuthed("DeployContainer", &arg, &reply); err != nil {
		return OutputError(err)
	}
	Log("-> ID: %s", reply.ID)
	if !c.Wait {
		return Output(map[string]interface{}{"id": reply.ID}, reply.ID, nil)
	}
	return (&WaitCommand{reply.ID}).Execute(args)
}

type CopyContainerCommand struct {
	ContainerID string `short:"c" long:"container" description:"the id of the container to copy"`
	ToHost      string `short:"H" long:"host" description:"the host to copy to"`
	PostCopy    int    `short:"p" long:"post" description:"what to do after the copy. (0 = nothing, 1 = cleanup datamodel only, 2 = teardown)"`
	Wait        bool   `long:"wait" description:"wait until the deploy is done before exiting"`
}

func (c *CopyContainerCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	Log("CopyContainer...")
	arg := ManagerCopyContainerArg{
		ManagerAuthArg: dummyAuthArg,
		ContainerID:    c.ContainerID,
		ToHost:         c.ToHost,
		PostCopy:       c.PostCopy,
	}
	var reply atlantis.AsyncReply
	if err := rpcClient.CallAuthed("CopyContainer", &arg, &reply); err != nil {
		return OutputError(err)
	}
	Log("-> ID: %s", reply.ID)
	if !c.Wait {
		return Output(map[string]interface{}{"id": reply.ID}, reply.ID, nil)
	}
	return (&WaitCommand{reply.ID}).Execute(args)
}

func OutputDeployReply(reply *ManagerDeployReply) error {
	Log("-> Status: %s", reply.Status)
	Log("-> Deployed Containers:")
	quietContainerIDs := make([]string, len(reply.Containers))
	for i, cont := range reply.Containers {
		Log("->   %s", cont.String())
		quietContainerIDs[i] = cont.ID
	}
	return Output(map[string]interface{}{"status": reply.Status, "containers": reply.Containers},
		quietContainerIDs, nil)
}

type DeployResultCommand struct {
	ID string `short:"i" long:"id" description:"the task ID to fetch the result for"`
}

func (c *DeployResultCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	args = ExtractArgs([]*string{&c.ID}, args)
	Log("Deploy Result...")
	arg := c.ID
	var reply ManagerDeployReply
	if err := rpcClient.Call("DeployResult", arg, &reply); err != nil {
		return OutputError(err)
	}
	return OutputDeployReply(&reply)
}

type TeardownCommand struct {
	App       string `short:"a" long:"app" description:"the app to teardown"`
	Sha       string `short:"s" long:"sha" description:"the sha to teardown"`
	Env       string `short:"e" long:"env" description:"the environment to teardown"`
	Container string `short:"c" long:"container" description:"the container to teardown"`
	All       bool   `long:"all" description:"teardown all containers in every supervisor"`
	Wait      bool   `long:"wait" description:"wait until the teardown is done before exiting"`
}

func (c *TeardownCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	Log("Teardown...")
	arg := ManagerTeardownArg{dummyAuthArg, c.App, c.Sha, c.Env, c.Container, c.All}
	var reply atlantis.AsyncReply
	if err := rpcClient.CallAuthed("Teardown", &arg, &reply); err != nil {
		return OutputError(err)
	}
	Log("-> ID: %s", reply.ID)
	if !c.Wait {
		return Output(map[string]interface{}{"id": reply.ID}, reply.ID, nil)
	}
	return (&WaitCommand{reply.ID}).Execute(args)
}

func OutputTeardownReply(reply *ManagerTeardownReply) error {
	Log("-> Status: %s", reply.Status)
	Log("-> Torn Containers:")
	for _, cont := range reply.ContainerIDs {
		Log("->   %s", cont)
	}
	return Output(map[string]interface{}{"status": reply.Status, "containerIDs": reply.ContainerIDs},
		reply.ContainerIDs, nil)
}

type TeardownResultCommand struct {
	ID string `short:"i" long:"id" description:"the task ID to fetch the result for"`
}

func (c *TeardownResultCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	args = ExtractArgs([]*string{&c.ID}, args)
	Log("Teardown Result...")
	arg := c.ID
	var reply ManagerTeardownReply
	if err := rpcClient.Call("TeardownResult", arg, &reply); err != nil {
		return OutputError(err)
	}
	return OutputTeardownReply(&reply)
}

type GetContainerCommand struct {
	Container string `short:"c" long:"container" description:"the container to get"`
}

func (c *GetContainerCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	args = ExtractArgs([]*string{&c.Container}, args)
	Log("Get Container...")
	arg := ManagerGetContainerArg{dummyAuthArg, c.Container}
	var reply ManagerGetContainerReply
	if err := rpcClient.CallAuthed("GetContainer", &arg, &reply); err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	Log("-> %s", reply.Container.String())
	return Output(map[string]interface{}{"status": reply.Status, "container": reply.Container},
		fmt.Sprintf("%s:%d", reply.Container.Host, reply.Container.PrimaryPort), nil)
}

type ListContainersCommand struct {
	App string `short:"a" long:"app" description:"the app to list"`
	Sha string `short:"s" long:"sha" description:"the sha to list"`
	Env string `short:"e" long:"env" description:"the environment to list"`
}

func (c *ListContainersCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	args = ExtractArgs([]*string{&c.App, &c.Sha, &c.Env}, args)
	Log("List Containers...")
	arg := ManagerListContainersArg{dummyAuthArg, c.App, c.Sha, c.Env}
	var reply ManagerListContainersReply
	if err := rpcClient.CallAuthed("ListContainers", &arg, &reply); err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	Log("-> containers:")
	for _, cont := range reply.ContainerIDs {
		Log("->   %s", cont)
	}
	return Output(map[string]interface{}{"status": reply.Status, "containerIDs": reply.ContainerIDs},
		reply.ContainerIDs, nil)
}

type ListEnvsCommand struct {
	App string `short:"a" long:"app" description:"the app to list (empty for all available envs)"`
	Sha string `short:"s" long:"sha" description:"the sha to list (empty for all available envs)"`
}

func (c *ListEnvsCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	args = ExtractArgs([]*string{&c.App, &c.Sha}, args)
	Log("List Envs...")
	arg := ManagerListEnvsArg{dummyAuthArg, c.App, c.Sha}
	var reply ManagerListEnvsReply
	if err := rpcClient.CallAuthed("ListEnvs", &arg, &reply); err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	Log("-> envs:")
	for _, env := range reply.Envs {
		Log("->   %s", env)
	}
	return Output(map[string]interface{}{"status": reply.Status, "envs": reply.Envs}, reply.Envs, nil)
}

type ListShasCommand struct {
	App string `short:"a" long:"app" description:"the app to list"`
}

func (c *ListShasCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	args = ExtractArgs([]*string{&c.App}, args)
	Log("List Shas...")
	arg := ManagerListShasArg{dummyAuthArg, c.App}
	var reply ManagerListShasReply
	if err := rpcClient.CallAuthed("ListShas", &arg, &reply); err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	Log("-> shas:")
	for _, sha := range reply.Shas {
		Log("->   %s", sha)
	}
	return Output(map[string]interface{}{"status": reply.Status, "shas": reply.Shas}, reply.Shas, nil)
}

type ListAppsCommand struct {
}

func (c *ListAppsCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	Log("List Apps...")
	arg := ManagerListAppsArg{dummyAuthArg}
	var reply ManagerListAppsReply
	if err := rpcClient.CallAuthed("ListApps", &arg, &reply); err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	Log("-> apps:")
	for _, app := range reply.Apps {
		Log("->   %s", app)
	}
	return Output(map[string]interface{}{"status": reply.Status, "apps": reply.Apps}, reply.Apps, nil)
}
