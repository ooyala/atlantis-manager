package client

import (
	. "atlantis/common"
	"errors"
)

type StatusCommand struct {
	Id string `short:"i" long:"id" description:"the task ID to fetch the status for"`
}

func (c *StatusCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	args = ExtractArgs([]*string{&c.Id}, args)
	Log("Task Status...")
	arg := c.Id
	var reply TaskStatus
	err = rpcClient.Call("Status", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> %s", reply.String())
	return Output(reply.Map(), reply.Status, nil)
}

type ResultCommand struct {
	Id string `short:"i" long:"id" description:"the task ID to fetch the result for"`
}

func (c *ResultCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	args = ExtractArgs([]*string{&c.Id}, args)
	arg := c.Id
	var reply TaskStatus
	err = rpcClient.Call("Status", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	switch reply.Name {
	case "Deploy":
		return (&DeployResultCommand{c.Id}).Execute(args)
	case "Teardown":
		return (&TeardownResultCommand{c.Id}).Execute(args)
	case "RegisterManager":
		return (&RegisterManagerResultCommand{c.Id}).Execute(args)
	case "UnregisterManager":
		return (&UnregisterManagerResultCommand{c.Id}).Execute(args)
	case "RegisterRouter":
		return (&RegisterRouterResultCommand{c.Id}).Execute(args)
	case "UnregisterRouter":
		return (&UnregisterRouterResultCommand{c.Id}).Execute(args)
	default:
		return OutputError(errors.New("Invalid Task Name: " + reply.Name))
	}
}

type WaitCommand struct {
	Id string `short:"i" long:"id" description:"the task ID to wait on"`
}

func (c *WaitCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	args = ExtractArgs([]*string{&c.Id}, args)
	Log("Waiting...")
	arg := c.Id
	var statusReply TaskStatus
	var currentStatus string
	if err := rpcClient.Call("Status", arg, &statusReply); err != nil {
		return OutputError(err)
	}
	for !statusReply.Done {
		if currentStatus != statusReply.Status {
			currentStatus = statusReply.Status
			Log(currentStatus)
		}
		if err := rpcClient.Call("Status", c.Id, &statusReply); err != nil {
			return OutputError(err)
		}
	}
	return (&ResultCommand{c.Id}).Execute(args)
}
