package rpc

import (
	. "atlantis/common"
	"atlantis/manager/datamodel"
	. "atlantis/manager/rpc/types"
	"atlantis/manager/supervisor"
	"errors"
	"fmt"
)

// NOTE[jigish]: these are simple pass-throughs to supervisor.

type AuthorizeSSHExecutor struct {
	arg   ManagerAuthorizeSSHArg
	reply *ManagerAuthorizeSSHReply
}

func (e *AuthorizeSSHExecutor) Request() interface{} {
	return e.arg
}

func (e *AuthorizeSSHExecutor) Result() interface{} {
	return e.reply
}

func (e *AuthorizeSSHExecutor) Description() string {
	return fmt.Sprintf("%s @ %s : \n%s", e.arg.User, e.arg.ContainerId, e.arg.PublicKey)
}

func (e *AuthorizeSSHExecutor) Execute(t *Task) error {
	if e.arg.PublicKey == "" {
		return errors.New("Please specify an SSH public key.")
	}
	if e.arg.ContainerId == "" {
		return errors.New("Please specify a container id.")
	}
	if e.arg.User == "" {
		return errors.New("Please specify a user.")
	}
	instance, err := datamodel.GetInstance(e.arg.ContainerId)
	if err != nil {
		return err
	}
	ihReply, err := supervisor.AuthorizeSSH(instance.Host, e.arg.ContainerId, e.arg.User, e.arg.PublicKey)
	if err != nil {
		e.reply.Status = StatusError
		return err
	}
	e.reply.Host = instance.Host
	e.reply.Port = ihReply.Port
	e.reply.Status = ihReply.Status
	return err
}

func (e *AuthorizeSSHExecutor) Authorize() error {
	return nil
}

type DeauthorizeSSHExecutor struct {
	arg   ManagerAuthorizeSSHArg
	reply *ManagerAuthorizeSSHReply
}

func (e *DeauthorizeSSHExecutor) Request() interface{} {
	return e.arg
}

func (e *DeauthorizeSSHExecutor) Result() interface{} {
	return e.reply
}

func (e *DeauthorizeSSHExecutor) Description() string {
	return fmt.Sprintf("%s @ %s", e.arg.User, e.arg.ContainerId)
}

func (e *DeauthorizeSSHExecutor) Execute(t *Task) error {
	if e.arg.ContainerId == "" {
		return errors.New("Please specify a container id.")
	}
	if e.arg.User == "" {
		return errors.New("Please specify a user.")
	}
	instance, err := datamodel.GetInstance(e.arg.ContainerId)
	if err != nil {
		return err
	}
	ihReply, err := supervisor.DeauthorizeSSH(instance.Host, e.arg.ContainerId, e.arg.User)
	e.reply.Status = ihReply.Status
	return err
}

func (e *DeauthorizeSSHExecutor) Authorize() error {
	return nil
}

func (m *ManagerRPC) AuthorizeSSH(arg ManagerAuthorizeSSHArg, reply *ManagerAuthorizeSSHReply) error {
	return NewTask("AuthorizeSSH", &AuthorizeSSHExecutor{arg, reply}).Run()
}

func (m *ManagerRPC) DeauthorizeSSH(arg ManagerAuthorizeSSHArg, reply *ManagerAuthorizeSSHReply) error {
	return NewTask("DeauthorizeSSH", &DeauthorizeSSHExecutor{arg, reply}).Run()
}
