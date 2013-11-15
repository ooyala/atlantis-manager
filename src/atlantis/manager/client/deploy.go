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
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerDeployArg{
		ManagerAuthArg: authArg,
		App:            c.App,
		Sha:            c.Sha,
		Env:            c.Env,
		Instances:      c.Instances,
		CPUShares:      c.CPUShares,
		MemoryLimit:    c.MemoryLimit,
		Dev:            c.Dev,
	}
	var reply atlantis.AsyncReply
	if err := rpcClient.Call("Deploy", arg, &reply); err != nil {
		return OutputError(err)
	}
	Log("-> ID: %s", reply.Id)
	if !c.Wait {
		return Output(map[string]interface{}{"id": reply.Id}, reply.Id, nil)
	}
	return (&WaitCommand{reply.Id}).Execute(args)
}

type CopyContainerCommand struct {
	ContainerId string `short:"c" long:"container" description:"the id of the container to copy"`
	Instances   uint   `short:"i" long:"instances" default:"1" description:"the number of instances to deploy in each AZ"`
	Wait        bool   `long:"wait" description:"wait until the deploy is done before exiting"`
}

func (c *CopyContainerCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	Log("CopyContainer...")
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerCopyContainerArg{ManagerAuthArg: authArg, ContainerId: c.ContainerId, Instances: c.Instances}
	var reply atlantis.AsyncReply
	if err := rpcClient.Call("CopyContainer", arg, &reply); err != nil {
		return OutputError(err)
	}
	Log("-> ID: %s", reply.Id)
	if !c.Wait {
		return Output(map[string]interface{}{"id": reply.Id}, reply.Id, nil)
	}
	return (&WaitCommand{reply.Id}).Execute(args)
}

type MoveContainerCommand struct {
	ContainerId string `short:"c" long:"container" description:"the id of the container to move"`
	Wait        bool   `long:"wait" description:"wait until the deploy is done before exiting"`
}

func (c *MoveContainerCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	Log("MoveContainer...")
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerMoveContainerArg{ManagerAuthArg: authArg, ContainerId: c.ContainerId}
	var reply atlantis.AsyncReply
	if err := rpcClient.Call("MoveContainer", arg, &reply); err != nil {
		return OutputError(err)
	}
	Log("-> ID: %s", reply.Id)
	if !c.Wait {
		return Output(map[string]interface{}{"id": reply.Id}, reply.Id, nil)
	}
	return (&WaitCommand{reply.Id}).Execute(args)
}

func OutputDeployReply(reply *ManagerDeployReply) error {
	Log("-> Status: %s", reply.Status)
	Log("-> Deployed Containers:")
	quietContainerIds := make([]string, len(reply.Containers))
	for i, cont := range reply.Containers {
		Log("->   %s", cont.String())
		quietContainerIds[i] = cont.Id
	}
	return Output(map[string]interface{}{"status": reply.Status, "containers": reply.Containers},
		quietContainerIds, nil)
}

type DeployResultCommand struct {
	Id string `short:"i" long:"id" description:"the task ID to fetch the result for"`
}

func (c *DeployResultCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	args = ExtractArgs([]*string{&c.Id}, args)
	Log("Deploy Result...")
	arg := c.Id
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
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerTeardownArg{authArg, c.App, c.Sha, c.Env, c.Container, c.All}
	var reply atlantis.AsyncReply
	if err := rpcClient.Call("Teardown", arg, &reply); err != nil {
		return OutputError(err)
	}
	Log("-> ID: %s", reply.Id)
	if !c.Wait {
		return Output(map[string]interface{}{"id": reply.Id}, reply.Id, nil)
	}
	return (&WaitCommand{reply.Id}).Execute(args)
}

func OutputTeardownReply(reply *ManagerTeardownReply) error {
	Log("-> Status: %s", reply.Status)
	Log("-> Torn Containers:")
	for _, cont := range reply.ContainerIds {
		Log("->   %s", cont)
	}
	return Output(map[string]interface{}{"status": reply.Status, "containerIds": reply.ContainerIds},
		reply.ContainerIds, nil)
}

type TeardownResultCommand struct {
	Id string `short:"i" long:"id" description:"the task ID to fetch the result for"`
}

func (c *TeardownResultCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	args = ExtractArgs([]*string{&c.Id}, args)
	Log("Teardown Result...")
	arg := c.Id
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
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerGetContainerArg{authArg, c.Container}
	var reply ManagerGetContainerReply
	if err := rpcClient.Call("GetContainer", arg, &reply); err != nil {
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
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerListContainersArg{authArg, c.App, c.Sha, c.Env}
	var reply ManagerListContainersReply
	if err := rpcClient.Call("ListContainers", arg, &reply); err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	Log("-> containers:")
	for _, cont := range reply.ContainerIds {
		Log("->   %s", cont)
	}
	return Output(map[string]interface{}{"status": reply.Status, "containerIds": reply.ContainerIds},
		reply.ContainerIds, nil)
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
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerListEnvsArg{authArg, c.App, c.Sha}
	var reply ManagerListEnvsReply
	if err := rpcClient.Call("ListEnvs", arg, &reply); err != nil {
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
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerListShasArg{authArg, c.App}
	var reply ManagerListShasReply
	if err := rpcClient.Call("ListShas", arg, &reply); err != nil {
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
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerListAppsArg{authArg}
	var reply ManagerListAppsReply
	if err := rpcClient.Call("ListApps", arg, &reply); err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	Log("-> apps:")
	for _, app := range reply.Apps {
		Log("->   %s", app)
	}
	return Output(map[string]interface{}{"status": reply.Status, "apps": reply.Apps}, reply.Apps, nil)
}
