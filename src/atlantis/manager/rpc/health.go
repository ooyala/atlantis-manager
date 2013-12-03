package rpc

import (
	. "atlantis/common"
	. "atlantis/manager/constant"
	. "atlantis/manager/rpc/types"
)

type HealthCheckExecutor struct {
	arg   ManagerHealthCheckArg
	reply *ManagerHealthCheckReply
}

func (e *HealthCheckExecutor) Request() interface{} {
	return e.arg
}

func (e *HealthCheckExecutor) Result() interface{} {
	return e.reply
}

func (e *HealthCheckExecutor) Description() string {
	return "HealthCheck"
}

func (e *HealthCheckExecutor) Execute(t *Task) error {
	e.reply.Region = Region
	e.reply.Zone = Zone
	if Tracker.UnderMaintenance() {
		e.reply.Status = StatusMaintenance
	} else {
		e.reply.Status = StatusOk
	}
	t.Log("[RPC][HealthCheck] -> region: %s", e.reply.Region)
	t.Log("[RPC][HealthCheck] -> zone: %s", e.reply.Zone)
	t.Log("[RPC][HealthCheck] -> status: %s", e.reply.Status)
	return nil
}

func (e *HealthCheckExecutor) Authorize() error {
	return nil // allow anyone to check health
}

func (e *HealthCheckExecutor) AllowDuringMaintenance() bool {
	return true // allow checking health during maintenance.
}

func (m *ManagerRPC) HealthCheck(arg ManagerHealthCheckArg, reply *ManagerHealthCheckReply) error {
	return NewTask("HealthCheck", &HealthCheckExecutor{arg, reply}).Run()
}
