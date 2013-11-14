package rpc

import (
	. "atlantis/common"
	. "atlantis/manager/constant"
)

type VersionExecutor struct {
	arg   VersionArg
	reply *VersionReply
}

func (e *VersionExecutor) Request() interface{} {
	return e.arg
}

func (e *VersionExecutor) Result() interface{} {
	return e.reply
}

func (e *VersionExecutor) Description() string {
	return "Version"
}

func (e *VersionExecutor) Execute(t *Task) error {
	e.reply.RPCVersion = ManagerRPCVersion
	e.reply.APIVersion = ManagerAPIVersion
	t.Log("-> RPC: %s", ManagerRPCVersion)
	t.Log("-> API: %s", ManagerAPIVersion)
	return nil
}

func (e *VersionExecutor) Authorize() error {
	return nil
}

func (e *VersionExecutor) AllowDuringMaintenance() bool {
	return true
}

func (o *Manager) Version(arg VersionArg, reply *VersionReply) error {
	return NewTask("Version", &VersionExecutor{arg, reply}).Run()
}
