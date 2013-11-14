package client

import (
	. "atlantis/manager/rpc/types"
)

type RegisterAppCommand struct {
	App  string `short:"a" long:"app" description:"the app to register"`
	Repo string `short:"g" long:"git" description:"the app's git repository"`
	Root string `short:"r" long:"root" description:"the app's root within the repo"`
}

func (c *RegisterAppCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Register App...")
	args = ExtractArgs([]*string{&c.App, &c.Repo, &c.Root}, args)
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerRegisterAppArg{ManagerAuthArg: authArg, Name: c.App, Repo: c.Repo, Root: c.Root}
	var reply ManagerRegisterAppReply
	err = rpcClient.Call("RegisterApp", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	return Output(map[string]interface{}{"status": reply.Status}, nil, nil)
}

type UnregisterAppCommand struct {
	App string `short:"a" long:"app" description:"the app to unregister"`
}

func (c *UnregisterAppCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Unregister App...")
	args = ExtractArgs([]*string{&c.App}, args)
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerRegisterAppArg{ManagerAuthArg: authArg, Name: c.App}
	var reply ManagerRegisterAppReply
	err = rpcClient.Call("UnregisterApp", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	return Output(map[string]interface{}{"status": reply.Status}, nil, nil)
}

type GetAppCommand struct {
	App string `short:"a" long:"app" description:"the app to unregister"`
}

func (c *GetAppCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Get App...")
	args = ExtractArgs([]*string{&c.App}, args)
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerGetAppArg{ManagerAuthArg: authArg, Name: c.App}
	var reply ManagerGetAppReply
	err = rpcClient.Call("GetApp", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	Log("-> name: %s", reply.App.Name)
	Log("-> repo: %s", reply.App.Repo)
	Log("-> root: %s", reply.App.Root)
	return Output(map[string]interface{}{"status": reply.Status, "app": reply.App}, nil, nil)
}

type ListRegisteredAppsCommand struct {
}

func (c *ListRegisteredAppsCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("List Registered Apps..")
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerListRegisteredAppsArg{authArg}
	var reply ManagerListRegisteredAppsReply
	err = rpcClient.Call("ListRegisteredApps", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	for _, app := range reply.Apps {
		Log("->   %s", app)
	}
	return Output(map[string]interface{}{"status": reply.Status, "apps": reply.Apps}, reply.Apps, nil)
}

type HealthCommand struct {
}

func (c *HealthCommand) Execute(args []string) error {
	InitNoLogin()
	Log("Manager Health Check...")
	arg := ManagerHealthCheckArg{}
	var reply ManagerHealthCheckReply
	err := rpcClient.Call("HealthCheck", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> region: %s", reply.Region)
	Log("-> status: %s", reply.Status)
	return Output(map[string]interface{}{"status": reply.Status, "region": reply.Region}, reply.Region, nil)
}

type RegisterManagerCommand struct {
	Host   string `short:"H" long:"host" description:"the host to register"`
	Region string `short:"r" long:"region" description:"the region to unregister"`
}

func (c *RegisterManagerCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Register Manager...")
	args = ExtractArgs([]*string{&c.Host, &c.Region}, args)
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerRegisterManagerArg{authArg, c.Host, c.Region}
	var reply ManagerRegisterManagerReply
	err = rpcClient.Call("RegisterManager", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	return Output(map[string]interface{}{"status": reply.Status}, nil, nil)
}

type UnregisterManagerCommand struct {
	Host   string `short:"H" long:"host" description:"the manager host to unregister"`
	Region string `short:"r" long:"region" description:"the region to unregister"`
}

func (c *UnregisterManagerCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Unregister Manager...")
	args = ExtractArgs([]*string{&c.Host, &c.Region}, args)
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerRegisterManagerArg{authArg, c.Host, c.Region}
	var reply ManagerRegisterManagerReply
	err = rpcClient.Call("UnregisterManager", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	return Output(map[string]interface{}{"status": reply.Status}, nil, nil)
}

type ListManagersCommand struct {
}

func (c *ListManagersCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("List Managers..")
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerListManagersArg{authArg}
	var reply ManagerListManagersReply
	err = rpcClient.Call("ListManagers", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	for region, managers := range reply.Managers {
		Log("-> %s:", region)
		for _, manager := range managers {
			Log("->   %s", manager)
		}
	}
	return Output(map[string]interface{}{"status": reply.Status, "managers": reply.Managers}, reply.Managers, nil)
}

type RegisterSupervisorCommand struct {
	Host string `short:"H" long:"host" description:"the supervisor host to register"`
}

func (c *RegisterSupervisorCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Register Supervisor...")
	args = ExtractArgs([]*string{&c.Host}, args)
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerRegisterSupervisorArg{authArg, c.Host}
	var reply ManagerRegisterSupervisorReply
	err = rpcClient.Call("RegisterSupervisor", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	return Output(map[string]interface{}{"status": reply.Status}, nil, nil)
}

type UnregisterSupervisorCommand struct {
	Host string `short:"H" long:"host" description:"the supervisor host to register"`
}

func (c *UnregisterSupervisorCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Unregister Supervisor...")
	args = ExtractArgs([]*string{&c.Host}, args)
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerRegisterSupervisorArg{authArg, c.Host}
	var reply ManagerRegisterSupervisorReply
	err = rpcClient.Call("UnregisterSupervisor", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	return Output(map[string]interface{}{"status": reply.Status}, nil, nil)
}

type ListSupervisorsCommand struct {
}

func (c *ListSupervisorsCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("List Supervisors..")
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerListSupervisorsArg{authArg}
	var reply ManagerListSupervisorsReply
	err = rpcClient.Call("ListSupervisors", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	for _, supervisor := range reply.Supervisors {
		Log("->   %s", supervisor)
	}
	return Output(map[string]interface{}{"status": reply.Status, "supervisors": reply.Supervisors}, reply.Supervisors,
		nil)
}
