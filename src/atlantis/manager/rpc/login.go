package rpc

import (
	. "atlantis/common"
	. "atlantis/manager/rpc/types"
)

type LoginExecutor struct {
	arg   ManagerLoginArg
	reply *ManagerLoginReply
}

func (e *LoginExecutor) Request() interface{} {
	return e.arg
}

func (e *LoginExecutor) Result() interface{} {
	return e.reply
}

func (e *LoginExecutor) Description() string {
	return "Login"
}

func (e *LoginExecutor) Execute(t *Task) error {
	e.reply.Secret = e.arg.Secret
	e.reply.LoggedIn = false
	auther := Authorizer{e.arg.User, e.arg.Pass, e.arg.Secret}
	err := auther.Authenticate()
	if err != nil {
		return err
	}
	e.reply.Secret = auther.Secret
	// reply.Secret should be blank if the username/pass is incorrect
	if e.reply.Secret != "" {
		e.reply.LoggedIn = true
	}
	return nil
}

func (e *LoginExecutor) Authorize() error {
	return nil
}

func (o *Manager) Login(arg ManagerLoginArg, reply *ManagerLoginReply) error {
	return NewTask("Login", &LoginExecutor{arg, reply}).Run()
}
