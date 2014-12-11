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
	Arg         ManagerDeployArg
	Reply       atlantis.AsyncReply
}

func (c *DeployCommand) Execute(args []string) error {
	return genericExecuter(c, args)
}

type DeployContainerCommand struct {
	ContainerID string `short:"c" long:"container" description:"the id of the container to replicate"`
	Instances   uint   `short:"i" long:"instances" default:"1" description:"the number of instances to deploy in each AZ"`
	Wait        bool   `long:"wait" description:"wait until the deploy is done before exiting"`
	Arg         ManagerDeployArg
	Reply       atlantis.AsyncReply
}

func (c *DeployContainerCommand) Execute(args []string) error {
	return genericExecuter(c, args)
}

type CopyContainerCommand struct {
	ContainerID string `short:"c" long:"container" description:"the id of the container to copy"`
	ToHost      string `short:"H" long:"host" description:"the host to copy to"`
	PostCopy    int    `short:"p" long:"post" description:"what to do after the copy. (0 = nothing, 1 = cleanup datamodel only, 2 = teardown)"`
	Wait        bool   `long:"wait" description:"wait until the deploy is done before exiting"`
	Arg         ManagerCopyContainerArg
	Reply       atlantis.AsyncReply
}

func (c *CopyContainerCommand) Execute(args []string) error {
	return genericExecuter(c, args)
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
	App         string `short:"a" long:"app" description:"the app to teardown"`
	Sha         string `short:"s" long:"sha" description:"the sha to teardown"`
	Env         string `short:"e" long:"env" description:"the environment to teardown"`
	ContainerID string `short:"c" long:"container" description:"the container to teardown"`
	All         bool   `long:"all" description:"teardown all containers in every supervisor"`
	Wait        bool   `long:"wait" description:"wait until the teardown is done before exiting"`
	Arg         ManagerTeardownArg
	Reply       atlantis.AsyncReply
}

func (c *TeardownCommand) Execute(args []string) error {
	return genericExecuter(c, args)
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
	ContainerID string `short:"c" long:"container" description:"the container to get"`
	Arg         ManagerGetContainerArg
	Reply       ManagerGetContainerReply
}

func (c *GetContainerCommand) Execute(args []string) error {
	return genericExecuter(c, args)
}

type ListContainersCommand struct {
	App        string `short:"a" long:"app" description:"the app to list"`
	Sha        string `short:"s" long:"sha" description:"the sha to list"`
	Env        string `short:"e" long:"env" description:"the environment to list"`
	Properties string `field:"ContainerIDs" name:"containers"`
	Arg        ManagerListContainersArg
	Reply      ManagerListContainersReply
}

func (c *ListContainersCommand) Execute(args []string) error {
	return genericExecuter(c, args)
}

type ListEnvsCommand struct {
	App   string `short:"a" long:"app" description:"the app to list (empty for all available envs)"`
	Sha   string `short:"s" long:"sha" description:"the sha to list (empty for all available envs)"`
	Arg   ManagerListEnvsArg
	Reply ManagerListEnvsReply
}

func (c *ListEnvsCommand) Execute(args []string) error {
	return genericExecuter(c, args)
}

type ListShasCommand struct {
	App   string `short:"a" long:"app" description:"the app to list"`
	Arg   ManagerListShasArg
	Reply ManagerListShasReply
}

func (c *ListShasCommand) Execute(args []string) error {
	return genericExecuter(c, args)
}

type ListAppsCommand struct {
	Arg   ManagerListAppsArg
	Reply ManagerListAppsReply
}

func (c *ListAppsCommand) Execute(args []string) error {
	return genericExecuter(c, args)
}
