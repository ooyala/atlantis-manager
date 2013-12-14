package client

import (
	. "atlantis/manager/rpc/types"
)

func OutputDependerAppReply(reply *ManagerDependerAppReply) error {
	Log("-> Status: %s", reply.Status)
	if reply.Dependee != nil {
		Log("-> Dependee:")
		Log("->   Name:  %s", reply.Dependee.Name)
		Log("->   Repo:  %s", reply.Dependee.Repo)
		Log("->   Root:  %s", reply.Dependee.Root)
		Log("->   Email: %s", reply.Dependee.Email)
		Log("->   Dependers:")
		for app, depends := range reply.Dependee.AllowedDependerApps {
			Log("->     %s : %t", app, depends)
		}
	}
	return Output(map[string]interface{}{"status": reply.Status, "dependee": reply.Dependee}, nil, nil)
}

type AddDependerAppCommand struct {
	Dependee string `short:"e" long:"dependee" description:"the app to add a depender for"`
	Depender string `short:"r" long:"depender" description:"the depender app"`
}

func (c *AddDependerAppCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Add Depender App...")
	args = ExtractArgs([]*string{&c.Dependee, &c.Depender}, args)
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerDependerAppArg{ManagerAuthArg: authArg, Dependee: c.Dependee, Depender: c.Depender}
	var reply ManagerDependerAppReply
	err = rpcClient.Call("AddDependerApp", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	return OutputDependerAppReply(&reply)
}

type RemoveDependerAppCommand struct {
	Dependee string `short:"e" long:"dependee" description:"the app to remove a depender for"`
	Depender string `short:"r" long:"depender" description:"the depender app"`
}

func (c *RemoveDependerAppCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Remove Depender App...")
	args = ExtractArgs([]*string{&c.Dependee, &c.Depender}, args)
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerDependerAppArg{ManagerAuthArg: authArg, Dependee: c.Dependee, Depender: c.Depender}
	var reply ManagerDependerAppReply
	err = rpcClient.Call("RemoveDependerApp", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	return OutputDependerAppReply(&reply)
}
