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
	supervisors, err := datamodel.ListSupervisors()
	if err != nil {
		return err
	}
	for _, super := range supervisors {
		t.Log("Updating %s", super)
		_, err := supervisor.UpdateProxy(super, e.arg.Sha)
		if err != nil {
			e.reply.Status = StatusError
			return err
		}
	}
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

func (m *ManagerRPC) UpdateProxy(arg ManagerUpdateProxyArg, reply *ManagerUpdateProxyReply) error {
	return NewTask("UpdateProxy", &UpdateProxyExecutor{arg, reply}).Run()
}
