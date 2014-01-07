package rpc

import (
	. "atlantis/common"
	"atlantis/manager/datamodel"
	. "atlantis/manager/rpc/types"
	"atlantis/manager/supervisor"
	"errors"
	"fmt"
)

type UpdateProxyExecutor struct {
	arg   ManagerUpdateProxyArg
	reply *ManagerUpdateProxyReply
}

func (e *UpdateProxyExecutor) Request() interface{} {
	return e.arg
}

func (e *UpdateProxyExecutor) Result() interface{} {
	return e.reply
}

func (e *UpdateProxyExecutor) Description() string {
	return fmt.Sprintf("%s", e.arg.Sha)
}

func (e *UpdateProxyExecutor) Authorize() error {
	return AuthorizeSuperUser(&e.arg.ManagerAuthArg)
}

func (e *UpdateProxyExecutor) Execute(t *Task) error {
	if e.arg.Sha == "" {
		return errors.New("Please specify a sha")
	}
	t.LogStatus("Listing Supervisors")
	supervisors, err := datamodel.ListSupervisors()
	if err != nil {
		return err
	}
	for _, super := range supervisors {
		t.LogStatus("Updating %s", super)
		_, err := supervisor.UpdateProxy(super, e.arg.Sha)
		if err != nil {
			e.reply.Status = StatusError
			return err
		}
	}
	t.LogStatus("Saving New SHA")
	lock := datamodel.NewProxyLock()
	lock.Lock()
	defer lock.Unlock()
	zp := datamodel.GetProxy()
	zp.Sha = e.arg.Sha
	err = zp.Save()
	if err != nil {
		e.reply.Status = StatusError
	} else {
		e.reply.Status = StatusOk
	}
	return err
}

func (m *ManagerRPC) UpdateProxy(arg ManagerUpdateProxyArg, reply *AsyncReply) error {
	return NewTask("UpdateProxy", &UpdateProxyExecutor{arg, &ManagerUpdateProxyReply{}}).RunAsync(reply)
}

func (m *ManagerRPC) UpdateProxyResult(id string, result *ManagerUpdateProxyReply) error {
	if id == "" {
		return errors.New("ID empty")
	}
	status, err := Tracker.Status(id)
	if status.Status == StatusUnknown {
		return errors.New("Unknown ID.")
	}
	if status.Name != "UpdateProxy" {
		return errors.New("ID is not a UpdateProxy.")
	}
	if !status.Done {
		return errors.New("UpdateProxy isn't done.")
	}
	if status.Status == StatusError || err != nil {
		return err
	}
	getResult := Tracker.Result(id)
	switch r := getResult.(type) {
	case *ManagerUpdateProxyReply:
		*result = *r
	default:
		// this should never happen
		return errors.New("Invalid Result Type.")
	}
	return nil
}

type ConfigureProxyExecutor struct {
	arg   ManagerConfigureProxyArg
	reply *ManagerConfigureProxyReply
}

func (e *ConfigureProxyExecutor) Request() interface{} {
	return e.arg
}

func (e *ConfigureProxyExecutor) Result() interface{} {
	return e.reply
}

func (e *ConfigureProxyExecutor) Description() string {
	return "ConfigureProxy"
}

func (e *ConfigureProxyExecutor) Authorize() error {
	return AuthorizeSuperUser(&e.arg.ManagerAuthArg)
}

func (e *ConfigureProxyExecutor) Execute(t *Task) error {
	pLock := datamodel.NewProxyLock()
	pLock.Lock()
	defer pLock.Unlock()
	if err := datamodel.ConfigureProxy(); err != nil {
		e.reply.Status = StatusError
		return err
	}
	e.reply.Status = StatusOk
	return nil
}

func (m *ManagerRPC) ConfigureProxy(arg ManagerConfigureProxyArg, reply *AsyncReply) error {
	return NewTask("ConfigureProxy", &ConfigureProxyExecutor{arg, &ManagerConfigureProxyReply{}}).RunAsync(reply)
}

func (m *ManagerRPC) ConfigureProxyResult(id string, result *ManagerConfigureProxyReply) error {
	if id == "" {
		return errors.New("ID empty")
	}
	status, err := Tracker.Status(id)
	if status.Status == StatusUnknown {
		return errors.New("Unknown ID.")
	}
	if status.Name != "ConfigureProxy" {
		return errors.New("ID is not a ConfigureProxy.")
	}
	if !status.Done {
		return errors.New("ConfigureProxy isn't done.")
	}
	if status.Status == StatusError || err != nil {
		return err
	}
	getResult := Tracker.Result(id)
	switch r := getResult.(type) {
	case *ManagerConfigureProxyReply:
		*result = *r
	default:
		// this should never happen
		return errors.New("Invalid Result Type.")
	}
	return nil
}
