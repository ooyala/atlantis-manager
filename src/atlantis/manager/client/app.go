package client

import (
	. "atlantis/manager/rpc/types"
	"encoding/json"
	"os"
)

type RequestAppDependencyCommand struct {
	App        string   `short:"a" long:"app" description:"the app to request a dependency for"`
	Dependency string   `short:"a" long:"dependency" description:"the dependency to request"`
	Envs       []string `short:"e" long:"env" description:"the envs to request the dependency in"`
}

func (c *RequestAppDependencyCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Request App Dependency...")
	args = ExtractArgs([]*string{&c.App, &c.Dependency}, args)
	user, secret, err := GetSecret()
	if err != nil {
		return OutputError(err)
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerRequestAppDependencyArg{
		ManagerAuthArg: authArg,
		App:            c.App,
		Dependency:     c.Dependency,
		Envs:           c.Envs,
	}
	var reply ManagerRequestAppDependencyReply
	err = rpcClient.Call("RequestAppDependency", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> Status: %s", reply.Status)
	return Output(map[string]interface{}{"status": reply.Status}, nil, nil)
}

// ----------------------------------------------------------------------------------------------------------
// Depender App Data Methods
// ----------------------------------------------------------------------------------------------------------

type AddDependerAppDataCommand struct {
	App      string `short:"a" long:"app" description:"the app to add a depender for"`
	FromFile string `short:"f" long:"file" description:"the file to pull the data from"`
}

func (c *AddDependerAppDataCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Add Depender App...")
	args = ExtractArgs([]*string{&c.App, &c.FromFile}, args)
	user, secret, err := GetSecret()
	if err != nil {
		return OutputError(err)
	}
	data := &DependerAppData{}
	file, err := os.Open(c.FromFile)
	if err != nil {
		return OutputError(err)
	}
	jsonDec := json.NewDecoder(file)
	if err := jsonDec.Decode(data); err != nil {
		return OutputError(err)
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerAddDependerAppDataArg{ManagerAuthArg: authArg, App: c.App, DependerAppData: data}
	var reply ManagerAddDependerAppDataReply
	err = rpcClient.Call("AddDependerAppData", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> Status: %s", reply.Status)
	LogApp(reply.App)
	return Output(map[string]interface{}{"status": reply.Status, "app": reply.App}, nil, nil)
}

type RemoveDependerAppDataCommand struct {
	App      string `short:"a" long:"app" description:"the app to remove a depender from"`
	Depender string `short:"r" long:"depender" description:"the depender app to remove"`
}

func (c *RemoveDependerAppDataCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Remove Depender App...")
	args = ExtractArgs([]*string{&c.App, &c.Depender}, args)
	user, secret, err := GetSecret()
	if err != nil {
		return OutputError(err)
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerRemoveDependerAppDataArg{ManagerAuthArg: authArg, App: c.App, Depender: c.Depender}
	var reply ManagerRemoveDependerAppDataReply
	err = rpcClient.Call("RemoveDependerAppData", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> Status: %s", reply.Status)
	LogApp(reply.App)
	return Output(map[string]interface{}{"status": reply.Status, "app": reply.App}, nil, nil)
}

type GetDependerAppDataCommand struct {
	App      string `short:"a" long:"app" description:"the app to get a depender from"`
	Depender string `short:"r" long:"depender" description:"the depender app to get"`
}

func (c *GetDependerAppDataCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Get Depender App...")
	args = ExtractArgs([]*string{&c.App, &c.Depender}, args)
	user, secret, err := GetSecret()
	if err != nil {
		return OutputError(err)
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerGetDependerAppDataArg{ManagerAuthArg: authArg, App: c.App, Depender: c.Depender}
	var reply ManagerGetDependerAppDataReply
	err = rpcClient.Call("GetDependerAppData", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> Status: %s", reply.Status)
	Log("-> DependerAppData:")
	LogDependerAppData("  ", reply.DependerAppData)
	return Output(map[string]interface{}{"status": reply.Status, "dependerAppData": reply.DependerAppData}, nil, nil)
}

// ----------------------------------------------------------------------------------------------------------
// Depender Env Data Methods
// ----------------------------------------------------------------------------------------------------------

type AddDependerEnvDataCommand struct {
	App      string `short:"a" long:"app" description:"the app to add an env for"`
	FromFile string `short:"f" long:"file" description:"the file to pull the data from"`
}

func (c *AddDependerEnvDataCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Add Depender Env...")
	args = ExtractArgs([]*string{&c.App, &c.FromFile}, args)
	user, secret, err := GetSecret()
	if err != nil {
		return OutputError(err)
	}
	data := &DependerEnvData{}
	file, err := os.Open(c.FromFile)
	if err != nil {
		return OutputError(err)
	}
	jsonDec := json.NewDecoder(file)
	if err := jsonDec.Decode(data); err != nil {
		return OutputError(err)
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerAddDependerEnvDataArg{ManagerAuthArg: authArg, App: c.App, DependerEnvData: data}
	var reply ManagerAddDependerEnvDataReply
	err = rpcClient.Call("AddDependerEnvData", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> Status: %s", reply.Status)
	LogApp(reply.App)
	return Output(map[string]interface{}{"status": reply.Status, "app": reply.App}, nil, nil)
}

type RemoveDependerEnvDataCommand struct {
	App string `short:"a" long:"app" description:"the app to remove an env from"`
	Env string `short:"e" long:"env" description:"the env to remove"`
}

func (c *RemoveDependerEnvDataCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Remove Depender Env...")
	args = ExtractArgs([]*string{&c.App, &c.Env}, args)
	user, secret, err := GetSecret()
	if err != nil {
		return OutputError(err)
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerRemoveDependerEnvDataArg{ManagerAuthArg: authArg, Env: c.Env, App: c.App}
	var reply ManagerRemoveDependerEnvDataReply
	err = rpcClient.Call("RemoveDependerEnvData", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> Status: %s", reply.Status)
	LogApp(reply.App)
	return Output(map[string]interface{}{"status": reply.Status, "app": reply.App}, nil, nil)
}

type GetDependerEnvDataCommand struct {
	App string `short:"a" long:"app" description:"the app to get an env from"`
	Env string `short:"e" long:"depender" description:"the env to get"`
}

func (c *GetDependerEnvDataCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Get Depender Env...")
	args = ExtractArgs([]*string{&c.App, &c.Env}, args)
	user, secret, err := GetSecret()
	if err != nil {
		return OutputError(err)
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerGetDependerEnvDataArg{ManagerAuthArg: authArg, Env: c.Env, App: c.App}
	var reply ManagerGetDependerEnvDataReply
	err = rpcClient.Call("GetDependerEnvData", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> Status: %s", reply.Status)
	Log("-> DependerEnvData:")
	LogDependerEnvData("  ", reply.DependerEnvData)
	return Output(map[string]interface{}{"status": reply.Status, "dependerEnvData": reply.DependerEnvData}, nil, nil)
}

// ----------------------------------------------------------------------------------------------------------
// Depender Env Data For Depender App Methods
// ----------------------------------------------------------------------------------------------------------

type AddDependerEnvDataForDependerAppCommand struct {
	App      string `short:"a" long:"app" description:"the app to add an env for"`
	Depender string `short:"r" long:"depender" description:"the depender to add an env for"`
	FromFile string `short:"f" long:"file" description:"the file to pull the data from"`
}

func (c *AddDependerEnvDataForDependerAppCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Add Depender Env For Depender App...")
	args = ExtractArgs([]*string{&c.App, &c.FromFile}, args)
	user, secret, err := GetSecret()
	if err != nil {
		return OutputError(err)
	}
	data := &DependerEnvData{}
	file, err := os.Open(c.FromFile)
	if err != nil {
		return OutputError(err)
	}
	jsonDec := json.NewDecoder(file)
	if err := jsonDec.Decode(data); err != nil {
		return OutputError(err)
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerAddDependerEnvDataForDependerAppArg{
		ManagerAuthArg:  authArg,
		App:             c.App,
		Depender:        c.Depender,
		DependerEnvData: data,
	}
	var reply ManagerAddDependerEnvDataForDependerAppReply
	err = rpcClient.Call("AddDependerEnvDataForDependerApp", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> Status: %s", reply.Status)
	LogApp(reply.App)
	return Output(map[string]interface{}{"status": reply.Status, "app": reply.App}, nil, nil)
}

type RemoveDependerEnvDataForDependerAppCommand struct {
	App      string `short:"a" long:"app" description:"the app to remove an env from"`
	Depender string `short:"r" long:"depender" description:"the depender to add an env for"`
	Env      string `short:"e" long:"env" description:"the env to remove"`
}

func (c *RemoveDependerEnvDataForDependerAppCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Remove Depender Env For Depender App...")
	args = ExtractArgs([]*string{&c.App, &c.Env}, args)
	user, secret, err := GetSecret()
	if err != nil {
		return OutputError(err)
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerRemoveDependerEnvDataForDependerAppArg{
		ManagerAuthArg: authArg,
		Env:            c.Env,
		Depender:       c.Depender,
		App:            c.App,
	}
	var reply ManagerRemoveDependerEnvDataForDependerAppReply
	err = rpcClient.Call("RemoveDependerEnvDataForDependerApp", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> Status: %s", reply.Status)
	LogApp(reply.App)
	return Output(map[string]interface{}{"status": reply.Status, "app": reply.App}, nil, nil)
}

type GetDependerEnvDataForDependerAppCommand struct {
	App      string `short:"a" long:"app" description:"the app to get an env from"`
	Depender string `short:"r" long:"depender" description:"the depender to add an env for"`
	Env      string `short:"e" long:"env" description:"the env to get"`
}

func (c *GetDependerEnvDataForDependerAppCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Get Depender Env For Depender App...")
	args = ExtractArgs([]*string{&c.App, &c.Env}, args)
	user, secret, err := GetSecret()
	if err != nil {
		return OutputError(err)
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerGetDependerEnvDataForDependerAppArg{
		ManagerAuthArg: authArg,
		Env:            c.Env,
		Depender:       c.Depender,
		App:            c.App,
	}
	var reply ManagerGetDependerEnvDataForDependerAppReply
	err = rpcClient.Call("GetDependerEnvDataForDependerApp", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> Status: %s", reply.Status)
	Log("-> DependerEnvData:")
	LogDependerEnvData("  ", reply.DependerEnvData)
	return Output(map[string]interface{}{"status": reply.Status, "dependerEnvData": reply.DependerEnvData}, nil, nil)
}
