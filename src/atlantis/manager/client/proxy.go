package client

import (
	atlantis "atlantis/common"
	. "atlantis/manager/rpc/types"
)

type UpdateProxyCommand struct {
	Wait bool   `long:"wait" description:"wait until done before exiting"`
	Sha  string `short:"s" long:"sha" description:"the sha to update to"`
}

func (c *UpdateProxyCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	args = ExtractArgs([]*string{&c.Sha}, args)
	Log("UpdateProxy...")
	arg := ManagerUpdateProxyArg{Sha: c.Sha}
	var reply atlantis.AsyncReply
	if err := rpcClient.Call("UpdateProxy", arg, &reply); err != nil {
		return OutputError(err)
	}
	Log("-> ID: %s", reply.ID)
	if !c.Wait {
		return Output(map[string]interface{}{"id": reply.ID}, reply.ID, nil)
	}
	return (&WaitCommand{reply.ID}).Execute(args)
}

func OutputUpdateProxyReply(reply *ManagerUpdateProxyReply) error {
	Log("-> Status: %s", reply.Status)
	return Output(map[string]interface{}{"status": reply.Status}, nil, nil)
}

type UpdateProxyResultCommand struct {
	ID string `short:"i" long:"id" description:"the task ID to fetch the result for"`
}

func (c *UpdateProxyResultCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	args = ExtractArgs([]*string{&c.ID}, args)
	Log("UpdateProxy Result...")
	arg := c.ID
	var reply ManagerUpdateProxyReply
	if err := rpcClient.Call("UpdateProxyResult", arg, &reply); err != nil {
		return OutputError(err)
	}
	return OutputUpdateProxyReply(&reply)
}

type ConfigureProxyCommand struct {
	Wait bool `long:"wait" description:"wait until done before exiting"`
}

func (c *ConfigureProxyCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	Log("ConfigureProxy...")
	arg := ManagerConfigureProxyArg{}
	var reply atlantis.AsyncReply
	if err := rpcClient.Call("ConfigureProxy", arg, &reply); err != nil {
		return OutputError(err)
	}
	Log("-> ID: %s", reply.ID)
	if !c.Wait {
		return Output(map[string]interface{}{"id": reply.ID}, reply.ID, nil)
	}
	return (&WaitCommand{reply.ID}).Execute(args)
}

func OutputConfigureProxyReply(reply *ManagerConfigureProxyReply) error {
	Log("-> Status: %s", reply.Status)
	return Output(map[string]interface{}{"status": reply.Status}, nil, nil)
}

type ConfigureProxyResultCommand struct {
	ID string `short:"i" long:"id" description:"the task ID to fetch the result for"`
}

func (c *ConfigureProxyResultCommand) Execute(args []string) error {
	if err := Init(); err != nil {
		return OutputError(err)
	}
	args = ExtractArgs([]*string{&c.ID}, args)
	Log("ConfigureProxy Result...")
	arg := c.ID
	var reply ManagerConfigureProxyReply
	if err := rpcClient.Call("ConfigureProxyResult", arg, &reply); err != nil {
		return OutputError(err)
	}
	return OutputConfigureProxyReply(&reply)
}
