package client

import (
	. "atlantis/manager/rpc/types"
)

type UpdateDepCommand struct {
	Name  string `short:"n" long:"name" description:"the name of the dependency"`
	Value string `short:"v" long:"value" description:"the value of the dependency"`
	Env   string `short:"e" long:"env" description:"the environment of the dependency"`
}

func (c *UpdateDepCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	Log("Update Dep...")
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerDepArg{authArg, c.Env, c.Name, c.Value}
	var reply ManagerDepReply
	if err := rpcClient.Call("UpdateDep", arg, &reply); err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	Log("-> value: %s", reply.Value)
	return Output(map[string]interface{}{"status": reply.Status, "value": reply.Value}, reply.Value, nil)
}

type ResolveDepsCommand struct {
	App      string   `short:"a" long:"app" description:"the app the resolve dependencies for"`
	Env      string   `short:"e" long:"env" description:"the environment of the dependencies to resolve"`
	DepNames []string `short:"d" long:"dep" description:"the dep names to resolve"`
}

func (c *ResolveDepsCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	Log("Resolve Deps...")
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerResolveDepsArg{ManagerAuthArg: authArg, App: c.App, Env: c.Env, DepNames: c.DepNames}
	var reply ManagerResolveDepsReply
	if err := rpcClient.Call("ResolveDeps", arg, &reply); err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	Log("-> deps:")
	for zone, deps := range reply.Deps {
		Log("->   %s", zone)
		for name, value := range deps {
			Log("->     %s : %s", name, value)
		}
	}
	return Output(map[string]interface{}{"status": reply.Status, "Deps": reply.Deps}, reply.Deps, nil)
}

type GetDepCommand struct {
	Name string `short:"n" long:"name" description:"the name of the dependency"`
	Env  string `short:"e" long:"env" description:"the environment of the dependency"`
}

func (c *GetDepCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	Log("Get Dep...")
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerDepArg{authArg, c.Env, c.Name, ""}
	var reply ManagerDepReply
	if err := rpcClient.Call("GetDep", arg, &reply); err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	Log("-> value: %s", reply.Value)
	return Output(map[string]interface{}{"status": reply.Status, "value": reply.Value}, reply.Value, nil)
}

type DeleteDepCommand struct {
	Name string `short:"n" long:"name" description:"the name of the dependency"`
	Env  string `short:"e" long:"env" description:"the environment of the dependency"`
}

func (c *DeleteDepCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	Log("Delete Dep...")
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerDepArg{authArg, c.Env, c.Name, ""}
	var reply ManagerDepReply
	if err := rpcClient.Call("DeleteDep", arg, &reply); err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	Log("-> value: %s", reply.Value)
	return Output(map[string]interface{}{"status": reply.Status, "value": reply.Value}, reply.Value, nil)
}

type UpdateEnvCommand struct {
	Name   string `short:"n" long:"name" description:"the name of the environment"`
	Parent string `short:"p" long:"parent" description:"the parent of the environment"`
}

func (c *UpdateEnvCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	Log("Update Env...")
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerEnvArg{authArg, c.Name, c.Parent}
	var reply ManagerEnvReply
	if err := rpcClient.Call("UpdateEnv", arg, &reply); err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	return Output(map[string]interface{}{"status": reply.Status}, nil, nil)
}

type GetEnvCommand struct {
	Name string `short:"n" long:"name" description:"the name of the environment"`
}

func (c *GetEnvCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	Log("Get Env...")
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerEnvArg{authArg, c.Name, ""}
	var reply ManagerEnvReply
	if err := rpcClient.Call("GetEnv", arg, &reply); err != nil {
		return OutputError(err)
	}
	quietMap := map[string]string{}
	Log("-> status: %s", reply.Status)
	Log("-> parent: %s", reply.Parent)
	quietMap["parent"] = reply.Parent
	Log("-> deps:")
	for name, value := range reply.Deps {
		quietMap["dep "+name] = value
		Log("->   %s = %s", name, value)
	}
	Log("-> resolved deps:")
	for name, value := range reply.ResolvedDeps {
		quietMap["resolved "+name] = value
		Log("->   %s = %s", name, value)
	}
	return Output(map[string]interface{}{"status": reply.Status, "parent": reply.Parent, "deps": reply.Deps,
		"resolvedDeps": reply.ResolvedDeps}, quietMap, nil)
}

type DeleteEnvCommand struct {
	Name string `short:"n" long:"name" description:"the name of the environment"`
}

func (c *DeleteEnvCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	Log("Delete Env...")
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerEnvArg{authArg, c.Name, ""}
	var reply ManagerEnvReply
	if err := rpcClient.Call("DeleteEnv", arg, &reply); err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	return Output(map[string]interface{}{"status": reply.Status}, nil, nil)
}
