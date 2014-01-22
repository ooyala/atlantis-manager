package client

import (
	. "atlantis/common"
	"errors"
)

type StatusCommand struct {
	ID string `short:"i" long:"id" description:"the task ID to fetch the status for"`
}

func (c *StatusCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	args = ExtractArgs([]*string{&c.ID}, args)
	Log("Task Status...")
	arg := c.ID
	var reply TaskStatus
	err = rpcClient.Call("Status", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> %s", reply.String())
	return Output(reply.Map(), reply.Status, nil)
}

type ResultCommand struct {
	ID string `short:"i" long:"id" description:"the task ID to fetch the result for"`
}

func (c *ResultCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	args = ExtractArgs([]*string{&c.ID}, args)
	arg := c.ID
	var reply TaskStatus
	err = rpcClient.Call("Status", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	switch reply.Name {
	case "Deploy":
		return (&DeployResultCommand{c.ID}).Execute(args)
	case "Teardown":
		return (&TeardownResultCommand{c.ID}).Execute(args)
	case "RegisterManager":
		return (&RegisterManagerResultCommand{c.ID}).Execute(args)
	case "UnregisterManager":
		return (&UnregisterManagerResultCommand{c.ID}).Execute(args)
	case "RegisterRouter":
		return (&RegisterRouterResultCommand{c.ID}).Execute(args)
	case "UnregisterRouter":
		return (&UnregisterRouterResultCommand{c.ID}).Execute(args)
	case "RegisterSupervisor":
		return (&RegisterSupervisorResultCommand{c.ID}).Execute(args)
	case "UnregisterSupervisor":
		return (&UnregisterSupervisorResultCommand{c.ID}).Execute(args)
	default:
		return OutputError(errors.New("Invalid Task Name: " + reply.Name))
	}
}

type WaitCommand struct {
	ID string `short:"i" long:"id" description:"the task ID to wait on"`
}

func (c *WaitCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	args = ExtractArgs([]*string{&c.ID}, args)
	Log("Waiting...")
	arg := c.ID
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
		if err := rpcClient.Call("Status", c.ID, &statusReply); err != nil {
			return OutputError(err)
		}
	}
	return (&ResultCommand{c.ID}).Execute(args)
}

type ListTaskIDsCommand struct {
	Type string `short:"t" long:"type" description:"the type of the task to list ids for"`
}

func (c *ListTaskIDsCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	args = ExtractArgs([]*string{&c.Type}, args)
	var ids []string
	if err := rpcClient.Call("ListTaskIDs", c.Type, &ids); err != nil {
		return OutputError(err)
	}
	return Output(map[string]interface{}{"ids": ids}, ids, nil)
}
