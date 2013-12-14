package client

import (
	atlantis "atlantis/common"
	. "atlantis/manager/rpc/types"
)

type RegisterRouterCommand struct {
	Wait     bool   `long:"wait" description:"wait until done before exiting"`
	Internal bool   `long:"internal" description:"true to list internal routers"`
	Zone     string `short:"z" long:"zone" description:"the zone to register in"`
	Host     string `short:"H" long:"host" description:"the host to register"`
}

func (c *RegisterRouterCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Register Router...")
	args = ExtractArgs([]*string{&c.Zone, &c.Host}, args)
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerRegisterRouterArg{
		ManagerAuthArg: authArg,
		Internal:       c.Internal,
		Zone:           c.Zone,
		Host:           c.Host,
	}
	var reply atlantis.AsyncReply
	err = rpcClient.Call("RegisterRouter", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> ID: %s", reply.ID)
	if !c.Wait {
		return Output(map[string]interface{}{"id": reply.ID}, reply.ID, nil)
	}
	return (&WaitCommand{reply.ID}).Execute(args)
}

type UnregisterRouterCommand struct {
	Wait     bool   `long:"wait" description:"wait until done before exiting"`
	Internal bool   `long:"internal" description:"true to list internal routers"`
	Zone     string `short:"z" long:"zone" description:"the zone to register in"`
	Host     string `short:"H" long:"host" description:"the host to unregister"`
}

func (c *UnregisterRouterCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Unregister Router...")
	args = ExtractArgs([]*string{&c.Zone, &c.Host}, args)
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerRegisterRouterArg{
		ManagerAuthArg: authArg,
		Internal:       c.Internal,
		Zone:           c.Zone,
		Host:           c.Host,
	}
	var reply atlantis.AsyncReply
	err = rpcClient.Call("UnregisterRouter", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> ID: %s", reply.ID)
	if !c.Wait {
		return Output(map[string]interface{}{"id": reply.ID}, reply.ID, nil)
	}
	return (&WaitCommand{reply.ID}).Execute(args)
}

func OutputRegisterRouterReply(reply *ManagerRegisterRouterReply) error {
	Log("-> Status: %s", reply.Status)
	if reply.Router != nil {
		Log("-> Router:")
		Log("->   Internal : %t", reply.Router.Internal)
		Log("->   Zone     : %s", reply.Router.Zone)
		Log("->   Host     : %s", reply.Router.Host)
		Log("->   CName    : %s", reply.Router.CName)
	}
	return Output(map[string]interface{}{"status": reply.Status, "router": reply.Router}, nil, nil)
}

type RegisterRouterResultCommand struct {
	ID string `short:"i" long:"id" description:"the task ID to fetch the result for"`
}

func (c *RegisterRouterResultCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	args = ExtractArgs([]*string{&c.ID}, args)
	Log("RegisterRouter Result...")
	arg := c.ID
	var reply ManagerRegisterRouterReply
	if err := rpcClient.Call("RegisterRouterResult", arg, &reply); err != nil {
		return OutputError(err)
	}
	return OutputRegisterRouterReply(&reply)
}

type UnregisterRouterResultCommand struct {
	ID string `short:"i" long:"id" description:"the task ID to fetch the result for"`
}

func (c *UnregisterRouterResultCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	args = ExtractArgs([]*string{&c.ID}, args)
	Log("UnregisterRouter Result...")
	arg := c.ID
	var reply ManagerRegisterRouterReply
	if err := rpcClient.Call("UnregisterRouterResult", arg, &reply); err != nil {
		return OutputError(err)
	}
	return OutputRegisterRouterReply(&reply)
}

type GetRouterCommand struct {
	Internal bool   `long:"internal" description:"true to list internal routers"`
	Zone     string `short:"z" long:"zone" description:"the zone to register in"`
	Host     string `short:"H" long:"host" description:"the host to get"`
}

func (c *GetRouterCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Get Router...")
	args = ExtractArgs([]*string{&c.Zone, &c.Host}, args)
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerGetRouterArg{ManagerAuthArg: authArg, Internal: c.Internal, Zone: c.Zone, Host: c.Host}
	var reply ManagerGetRouterReply
	err = rpcClient.Call("GetRouter", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> Status: %s", reply.Status)
	if reply.Router != nil {
		Log("-> Router:")
		Log("->   Internal : %t", reply.Router.Internal)
		Log("->   Zone     : %s", reply.Router.Zone)
		Log("->   Host     : %s", reply.Router.Host)
		Log("->   CName    : %s", reply.Router.CName)
	}
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
	App   string `short:"a" long:"app" description:"the app to register"`
	Repo  string `short:"g" long:"git" description:"the app's git repository"`
	Root  string `short:"r" long:"root" description:"the app's root within the repo"`
	Email string `short:"e" long:"email" description"the email of the app's owner"`
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
	arg := ManagerRegisterAppArg{ManagerAuthArg: authArg, Name: c.App, Repo: c.Repo, Root: c.Root, Email: c.Email}
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
	App string `short:"a" long:"app" description:"the app to get"`
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
	Log("-> name:  %s", reply.App.Name)
	Log("-> repo:  %s", reply.App.Repo)
	Log("-> root:  %s", reply.App.Root)
	Log("-> email: %s", reply.App.Email)
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
	Wait          bool   `long:"wait" description:"wait until done before exiting"`
	Host          string `short:"H" long:"host" description:"the host to register"`
	Region        string `short:"r" long:"region" description:"the region to register"`
	ManagerCName  string `long:"manager-cname" description:"the manager's cname if it already has one"`
	RegistryCName string `long:"registry-cname" description:"the registry's cname if it already has one"`
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
	arg := ManagerRegisterManagerArg{
		ManagerAuthArg: authArg,
		Host:           c.Host,
		Region:         c.Region,
		ManagerCName:   c.ManagerCName,
		RegistryCName:  c.RegistryCName,
	}
	var reply atlantis.AsyncReply
	err = rpcClient.Call("RegisterManager", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> ID: %s", reply.ID)
	if !c.Wait {
		return Output(map[string]interface{}{"id": reply.ID}, reply.ID, nil)
	}
	return (&WaitCommand{reply.ID}).Execute(args)
}

type UnregisterManagerCommand struct {
	Wait   bool   `long:"wait" description:"wait until done before exiting"`
	Host   string `short:"H" long:"host" description:"the host to register"`
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
	arg := ManagerRegisterManagerArg{ManagerAuthArg: authArg, Host: c.Host, Region: c.Region}
	var reply atlantis.AsyncReply
	err = rpcClient.Call("UnregisterManager", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> ID: %s", reply.ID)
	if !c.Wait {
		return Output(map[string]interface{}{"id": reply.ID}, reply.ID, nil)
	}
	return (&WaitCommand{reply.ID}).Execute(args)
}

func OutputRegisterManagerReply(reply *ManagerRegisterManagerReply) error {
	Log("-> Status: %s", reply.Status)
	if reply.Manager == nil {
		Log("-> Manager:")
		Log("->   Region:         %s", reply.Manager.Region)
		Log("->   Host:           %s", reply.Manager.Host)
		Log("->   Registry CName: %s", reply.Manager.RegistryCName)
		Log("->   Manager CName:  %s", reply.Manager.ManagerCName)
	}
	return Output(map[string]interface{}{"status": reply.Status, "manager": reply.Manager}, nil, nil)
}

type RegisterManagerResultCommand struct {
	ID string `short:"i" long:"id" description:"the task ID to fetch the result for"`
}

func (c *RegisterManagerResultCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	args = ExtractArgs([]*string{&c.ID}, args)
	Log("RegisterManager Result...")
	arg := c.ID
	var reply ManagerRegisterManagerReply
	if err := rpcClient.Call("RegisterManagerResult", arg, &reply); err != nil {
		return OutputError(err)
	}
	return OutputRegisterManagerReply(&reply)
}

type UnregisterManagerResultCommand struct {
	ID string `short:"i" long:"id" description:"the task ID to fetch the result for"`
}

func (c *UnregisterManagerResultCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	args = ExtractArgs([]*string{&c.ID}, args)
	Log("UnregisterManager Result...")
	arg := c.ID
	var reply ManagerRegisterManagerReply
	if err := rpcClient.Call("UnregisterManagerResult", arg, &reply); err != nil {
		return OutputError(err)
	}
	return OutputRegisterManagerReply(&reply)
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

func OutputGetManagerReply(reply *ManagerGetManagerReply) error {
	Log("-> Status: %s", reply.Status)
	if reply.Manager == nil {
		Log("-> Manager:")
		Log("->   Region:         %s", reply.Manager.Region)
		Log("->   Host:           %s", reply.Manager.Host)
		Log("->   Registry CName: %s", reply.Manager.RegistryCName)
		Log("->   Manager CName:  %s", reply.Manager.ManagerCName)
		Log("->   Roles:")
		for role, typeMap := range reply.Manager.Roles {
			Log("->     %s", role)
			for typeName, val := range typeMap {
				Log("->       %s : %t", typeName, val)
			}
		}
	}
	return Output(map[string]interface{}{"status": reply.Status, "manager": reply.Manager}, nil, nil)
}

type GetManagerCommand struct {
	Region string `short:"r" long:"region" description:"the region to get"`
	Host   string `short:"H" long:"host" description:"the host to get"`
}

func (c *GetManagerCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Get Manager...")
	args = ExtractArgs([]*string{&c.Region, &c.Host}, args)
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerGetManagerArg{ManagerAuthArg: authArg, Region: c.Region, Host: c.Host}
	var reply ManagerGetManagerReply
	err = rpcClient.Call("GetManager", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	return OutputGetManagerReply(&reply)
}

type GetSelfCommand struct {
}

func (c *GetSelfCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Get Self...")
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerGetSelfArg{ManagerAuthArg: authArg}
	var reply ManagerGetManagerReply
	err = rpcClient.Call("GetSelf", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	return OutputGetManagerReply(&reply)
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
