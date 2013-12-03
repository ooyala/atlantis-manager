package rpc

import (
	. "atlantis/common"
	"atlantis/manager/datamodel"
	. "atlantis/manager/rpc/types"
	"atlantis/manager/supervisor"
	"errors"
	"fmt"
)

// Container Maintenance Management
// NOTE[jigish]: this is a simple pass-through to supervisor with auth

type ContainerMaintenanceExecutor struct {
	arg   ManagerContainerMaintenanceArg
	reply *ManagerContainerMaintenanceReply
}

func (e *ContainerMaintenanceExecutor) Request() interface{} {
	return e.arg
}

func (e *ContainerMaintenanceExecutor) Result() interface{} {
	return e.reply
}

func (e *ContainerMaintenanceExecutor) Description() string {
	return fmt.Sprintf("%s : %t", e.arg.ContainerId, e.arg.Maintenance)
}

func (e *ContainerMaintenanceExecutor) Execute(t *Task) error {
	if e.arg.ContainerId == "" {
		return errors.New("Please specify a container id.")
	}
	instance, err := datamodel.GetInstance(e.arg.ContainerId)
	if err != nil {
		return err
	}
	ihReply, err := supervisor.ContainerMaintenance(instance.Host, e.arg.ContainerId, e.arg.Maintenance)
	e.reply.Status = ihReply.Status
	return err
}

func (e *ContainerMaintenanceExecutor) Authorize() error {
	instance, err := datamodel.GetInstance(e.arg.ContainerId)
	if err != nil {
		return err
	}
	return AuthorizeApp(&e.arg.ManagerAuthArg, instance.App)
}

func (m *ManagerRPC) ContainerMaintenance(arg ManagerContainerMaintenanceArg,
	reply *ManagerContainerMaintenanceReply) error {
	return NewTask("ContainerMaintenance", &ContainerMaintenanceExecutor{arg, reply}).Run()
}

// Manager Idle Check

type IdleExecutor struct {
	arg   ManagerIdleArg
	reply *ManagerIdleReply
}

func (e *IdleExecutor) Request() interface{} {
	return e.arg
}

func (e *IdleExecutor) Result() interface{} {
	return e.reply
}

func (e *IdleExecutor) Description() string {
	return "Idle?"
}

func (e *IdleExecutor) Execute(t *Task) error {
	e.reply.Idle = Tracker.Idle(t)
	e.reply.Status = StatusOk
	return nil
}

func (e *IdleExecutor) Authorize() error {
	return nil // let anybody ask if we're idle. i dont care.
}

func (e *IdleExecutor) AllowDuringMaintenance() bool {
	return true // allow running thus during maintenance
}

func (m *ManagerRPC) Idle(arg ManagerIdleArg, reply *ManagerIdleReply) error {
	return NewTask("Idle", &IdleExecutor{arg, reply}).Run()
}
