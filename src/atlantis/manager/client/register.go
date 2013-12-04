package client

import (
	. "atlantis/manager/rpc/types"
)

type RegisterRouterCommand struct {
	Internal bool   `long:"internal" description:"true to list internal routers"`
	Zone     string `short:"z" long:"zone" description:"the zone to register in"`
	IP       string `short:"i" long:"ip" description:"the IP to register"`
}

func (c *RegisterRouterCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Register Router...")
	args = ExtractArgs([]*string{&c.Zone, &c.IP}, args)
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerRegisterRouterArg{ManagerAuthArg: authArg, Internal: c.Internal, Zone: c.Zone, IP: c.IP}
	var reply ManagerRegisterRouterReply
	err = rpcClient.Call("RegisterRouter", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	Log("-> zone            : %s", reply.Router.Zone)
	Log("-> ip              : %s", reply.Router.IP)
	Log("-> cname           : %s", reply.Router.CName)
	return Output(map[string]interface{}{"status": reply.Status}, nil, nil)
}

type UnregisterRouterCommand struct {
	Internal bool   `long:"internal" description:"true to list internal routers"`
	Zone     string `short:"z" long:"zone" description:"the zone to register in"`
	IP       string `short:"i" long:"ip" description:"the IP to register"`
}

func (c *UnregisterRouterCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Unregister Router...")
	args = ExtractArgs([]*string{&c.Zone, &c.IP}, args)
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerRegisterRouterArg{ManagerAuthArg: authArg, Internal: c.Internal, Zone: c.Zone, IP: c.IP}
	var reply ManagerRegisterRouterReply
	err = rpcClient.Call("UnregisterRouter", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	return Output(map[string]interface{}{"status": reply.Status}, nil, nil)
}

type GetRouterCommand struct {
	Internal bool   `long:"internal" description:"true to list internal routers"`
	Zone     string `short:"z" long:"zone" description:"the zone to register in"`
	IP       string `short:"i" long:"ip" description:"the IP to register"`
}

func (c *GetRouterCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Get Router...")
	args = ExtractArgs([]*string{&c.Zone, &c.IP}, args)
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerGetRouterArg{ManagerAuthArg: authArg, Internal: c.Internal, Zone: c.Zone, IP: c.IP}
	var reply ManagerGetRouterReply
	err = rpcClient.Call("GetRouter", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	Log("-> zone            : %s", reply.Router.Zone)
	Log("-> ip              : %s", reply.Router.IP)
	Log("-> cname           : %s", reply.Router.CName)
	return Output(map[string]interface{}{"status": reply.Status, "router": reply.Router}, nil, nil)
}

type ListRoutersCommand struct {
	Internal bool `long:"internal" description:"true to list internal routers"`
}

func (c *ListRoutersCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("List Routers..")
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerListRoutersArg{ManagerAuthArg: authArg, Internal: c.Internal}
	var reply ManagerListRoutersReply
	err = rpcClient.Call("ListRouters", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	for zone, routers := range reply.Routers {
		Log("->   %s", zone)
		for _, router := range routers {
			Log("->     %s", router)
		}
	}
	return Output(map[string]interface{}{"status": reply.Status, "routers": reply.Routers}, reply.Routers, nil)
}

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
	IP            string `short:"i" long:"ip" description:"the ip to register"`
	Region        string `short:"r" long:"region" description:"the region to unregister"`
	ManagerCName  string `long:"manager-cname" description:"the manager's cname if it already has one"`
	RegistryCName string `long:"registry-cname" description:"the registry's cname if it already has one"`
}

func (c *RegisterManagerCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Register Manager...")
	args = ExtractArgs([]*string{&c.IP, &c.Region}, args)
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerRegisterManagerArg{
		ManagerAuthArg: authArg,
		IP:             c.IP,
		Region:         c.Region,
		ManagerCName:   c.ManagerCName,
		RegistryCName:  c.RegistryCName,
	}
	var reply ManagerRegisterManagerReply
	err = rpcClient.Call("RegisterManager", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	Log("-> region:         %s", reply.Manager.Region)
	Log("-> ip:             %s", reply.Manager.IP)
	Log("-> manager cname:  %s", reply.Manager.ManagerCName)
	Log("-> registry cname: %s", reply.Manager.RegistryCName)
	return Output(map[string]interface{}{"status": reply.Status, "manager": reply.Manager}, nil, nil)
}

type UnregisterManagerCommand struct {
	IP     string `short:"i" long:"ip" description:"the ip to unregister"`
	Region string `short:"r" long:"region" description:"the region to ununregister"`
}

func (c *UnregisterManagerCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Unregister Manager...")
	args = ExtractArgs([]*string{&c.IP, &c.Region}, args)
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerRegisterManagerArg{ManagerAuthArg: authArg, IP: c.IP, Region: c.Region}
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
