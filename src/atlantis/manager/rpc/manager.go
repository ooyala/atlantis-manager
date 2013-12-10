package rpc

import (
	. "atlantis/common"
	"atlantis/manager/datamodel"
	. "atlantis/manager/rpc/types"
	"errors"
	"fmt"
)

type AddRoleExecutor struct {
	arg   ManagerRoleArg
	reply *ManagerRoleReply
}

func (e *AddRoleExecutor) Request() interface{} {
	return e.arg
}

func (e *AddRoleExecutor) Result() interface{} {
	return e.reply
}

func (e *AddRoleExecutor) Description() string {
	return fmt.Sprintf("%s in %s -> %s + %s", e.arg.Host, e.arg.Region, e.arg.Role, e.arg.Type)
}

func (e *AddRoleExecutor) Authorize() error {
	return AuthorizeSuperUser(&e.arg.ManagerAuthArg)
}

func (e *AddRoleExecutor) Execute(t *Task) error {
	if e.arg.Host == "" {
		return errors.New("Please specify a host")
	}
	if e.arg.Region == "" {
		return errors.New("Please specify a region")
	}
	if e.arg.Role == "" {
		return errors.New("Please specify a role")
	}
	zkManager, err := datamodel.GetManager(e.arg.Region, e.arg.Host)
	if err != nil {
		e.reply.Status = StatusError
		return err
	}
	err = zkManager.AddRole(e.arg.Role, e.arg.Type)
	if err != nil {
		e.reply.Status = StatusError
		return err
	}
	e.reply.Status = StatusOk
	castedManager := Manager(*zkManager)
	e.reply.Manager = &castedManager
	return err
}

func (m *ManagerRPC) AddRole(arg ManagerRoleArg, reply *ManagerRoleReply) error {
	return NewTask("AddRole", &AddRoleExecutor{arg, reply}).Run()
}

type RemoveRoleExecutor struct {
	arg   ManagerRoleArg
	reply *ManagerRoleReply
}

func (e *RemoveRoleExecutor) Request() interface{} {
	return e.arg
}

func (e *RemoveRoleExecutor) Result() interface{} {
	return e.reply
}

func (e *RemoveRoleExecutor) Description() string {
	return fmt.Sprintf("%s in %s -> %s - %s", e.arg.Host, e.arg.Region, e.arg.Role, e.arg.Type)
}

func (e *RemoveRoleExecutor) Authorize() error {
	return AuthorizeSuperUser(&e.arg.ManagerAuthArg)
}

func (e *RemoveRoleExecutor) Execute(t *Task) error {
	if e.arg.Host == "" {
		return errors.New("Please specify a host")
	}
	if e.arg.Region == "" {
		return errors.New("Please specify a region")
	}
	if e.arg.Role == "" {
		return errors.New("Please specify a role")
	}
	zkManager, err := datamodel.GetManager(e.arg.Region, e.arg.Host)
	if err != nil {
		e.reply.Status = StatusError
		return err
	}
	err = zkManager.RemoveRole(e.arg.Role, e.arg.Type)
	if err != nil {
		e.reply.Status = StatusError
		return err
	}
	e.reply.Status = StatusOk
	castedManager := Manager(*zkManager)
	e.reply.Manager = &castedManager
	return err
}

func (m *ManagerRPC) RemoveRole(arg ManagerRoleArg, reply *ManagerRoleReply) error {
	return NewTask("RemoveRole", &RemoveRoleExecutor{arg, reply}).Run()
}

type HasRoleExecutor struct {
	arg   ManagerRoleArg
	reply *ManagerHasRoleReply
}

func (e *HasRoleExecutor) Request() interface{} {
	return e.arg
}

func (e *HasRoleExecutor) Result() interface{} {
	return e.reply
}

func (e *HasRoleExecutor) Description() string {
	return fmt.Sprintf("%s in %s -> %s ? %s", e.arg.Host, e.arg.Region, e.arg.Role, e.arg.Type)
}

func (e *HasRoleExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

func (e *HasRoleExecutor) Execute(t *Task) error {
	if e.arg.Host == "" {
		return errors.New("Please specify a host")
	}
	if e.arg.Region == "" {
		return errors.New("Please specify a region")
	}
	if e.arg.Role == "" {
		return errors.New("Please specify a role")
	}
	if e.arg.Type == "" {
		return errors.New("Please specify a type")
	}
	has, err := datamodel.ManagerHasRole(e.arg.Region, e.arg.Host, e.arg.Role, e.arg.Type)
	if err != nil {
		e.reply.Status = StatusError
		return err
	}
	e.reply.HasRole = has
	e.reply.Status = StatusOk
	return err
}

func (m *ManagerRPC) HasRole(arg ManagerRoleArg, reply *ManagerHasRoleReply) error {
	return NewTask("HasRole", &HasRoleExecutor{arg, reply}).Run()
}
