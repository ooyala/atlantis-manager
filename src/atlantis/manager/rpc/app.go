package rpc

import (
	. "atlantis/common"
	"atlantis/manager/datamodel"
	. "atlantis/manager/rpc/types"
	"errors"
	"fmt"
)

type AddDependerAppExecutor struct {
	arg   ManagerDependerAppArg
	reply *ManagerDependerAppReply
}

func (e *AddDependerAppExecutor) Request() interface{} {
	return e.arg
}

func (e *AddDependerAppExecutor) Result() interface{} {
	return e.reply
}

func (e *AddDependerAppExecutor) Description() string {
	return fmt.Sprintf("[+] %s depends on %s", e.arg.Depender, e.arg.Dependee)
}

func (e *AddDependerAppExecutor) Authorize() error {
	return AuthorizeApp(&e.arg.ManagerAuthArg, e.arg.Dependee)
}

func (e *AddDependerAppExecutor) Execute(t *Task) error {
	if e.arg.Depender == "" {
		return errors.New("Please specify a depender app")
	}
	if e.arg.Dependee == "" {
		return errors.New("Please specify a dependee app")
	}
	zkApp, err := datamodel.GetApp(e.arg.Dependee)
	if err != nil {
		e.reply.Status = StatusError
		return err
	}
	err = zkApp.AddDepender(e.arg.Depender)
	if err != nil {
		e.reply.Status = StatusError
		return err
	}
	e.reply.Status = StatusOk
	castedApp := App(*zkApp)
	e.reply.Dependee = &castedApp
	return err
}

func (m *ManagerRPC) AddDependerApp(arg ManagerDependerAppArg, reply *ManagerDependerAppReply) error {
	return NewTask("AddDependerApp", &AddDependerAppExecutor{arg, reply}).Run()
}

type RemoveDependerAppExecutor struct {
	arg   ManagerDependerAppArg
	reply *ManagerDependerAppReply
}

func (e *RemoveDependerAppExecutor) Request() interface{} {
	return e.arg
}

func (e *RemoveDependerAppExecutor) Result() interface{} {
	return e.reply
}

func (e *RemoveDependerAppExecutor) Description() string {
	return fmt.Sprintf("[-] %s depends on %s", e.arg.Depender, e.arg.Dependee)
}

func (e *RemoveDependerAppExecutor) Authorize() error {
	return AuthorizeApp(&e.arg.ManagerAuthArg, e.arg.Dependee)
}

func (e *RemoveDependerAppExecutor) Execute(t *Task) error {
	if e.arg.Depender == "" {
		return errors.New("Please specify a depender app")
	}
	if e.arg.Dependee == "" {
		return errors.New("Please specify a dependee app")
	}
	zkApp, err := datamodel.GetApp(e.arg.Dependee)
	if err != nil {
		e.reply.Status = StatusError
		return err
	}
	err = zkApp.RemoveDepender(e.arg.Depender)
	if err != nil {
		e.reply.Status = StatusError
		return err
	}
	e.reply.Status = StatusOk
	castedApp := App(*zkApp)
	e.reply.Dependee = &castedApp
	return err
}

func (m *ManagerRPC) RemoveDependerApp(arg ManagerDependerAppArg, reply *ManagerDependerAppReply) error {
	return NewTask("RemoveDependerApp", &RemoveDependerAppExecutor{arg, reply}).Run()
}
