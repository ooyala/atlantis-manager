package rpc

import (
	. "atlantis/common"
	"atlantis/manager/datamodel"
	. "atlantis/manager/rpc/types"
	"atlantis/manager/supervisor"
	. "atlantis/supervisor/rpc/types"
	"errors"
	"fmt"
)

type GetContainerExecutor struct {
	arg   ManagerGetContainerArg
	reply *ManagerGetContainerReply
}

func (e *GetContainerExecutor) Request() interface{} {
	return e.arg
}

func (e *GetContainerExecutor) Result() interface{} {
	return e.reply
}

func (e *GetContainerExecutor) Description() string {
	return e.arg.ContainerID
}

func (e *GetContainerExecutor) Execute(t *Task) (err error) {
	if e.arg.ContainerID == "" {
		return errors.New("Container ID is empty")
	}
	instance, err := datamodel.GetInstance(e.arg.ContainerID)
	if err != nil {
		return err
	}
	var ihReply *SupervisorGetReply
	ihReply, err = supervisor.Get(instance.Host, e.arg.ContainerID)
	e.reply.Status = ihReply.Status
	ihReply.Container.Host = instance.Host
	e.reply.Container = ihReply.Container
	if err != nil {
		e.reply.Status = StatusError
	} else {
		e.reply.Status = StatusOk
	}
	return err
}

func (e *GetContainerExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

func (m *ManagerRPC) GetContainer(arg ManagerGetContainerArg, reply *ManagerGetContainerReply) error {
	return NewTask("GetContainer", &GetContainerExecutor{arg, reply}).Run()
}

type ListContainersExecutor struct {
	arg   ManagerListContainersArg
	reply *ManagerListContainersReply
}

func (e *ListContainersExecutor) Request() interface{} {
	return e.arg
}

func (e *ListContainersExecutor) Result() interface{} {
	return e.reply
}

func (e *ListContainersExecutor) Description() string {
	return fmt.Sprintf("%s @ %s in %s", e.arg.App, e.arg.Sha, e.arg.Env)
}

func (e *ListContainersExecutor) Execute(t *Task) error {
	var err error
	if e.arg.App == "" && e.arg.Sha == "" && e.arg.Env == "" {
		// try to list all instances
		e.reply.ContainerIDs, err = datamodel.ListAllInstances()
		if err != nil {
			e.reply.Status = StatusError
		} else {
			e.reply.Status = StatusOk
		}
		return err
	}
	if e.arg.App == "" {
		return errors.New("App is empty")
	}
	if e.arg.Sha == "" {
		return errors.New("Sha is empty")
	}
	if e.arg.Env == "" {
		return errors.New("Environment is empty")
	}
	e.reply.ContainerIDs, err = datamodel.ListInstances(e.arg.App, e.arg.Sha, e.arg.Env)
	if err != nil {
		e.reply.Status = StatusError
	} else {
		e.reply.Status = StatusOk
	}
	return err
}

func (e *ListContainersExecutor) Authorize() error {
	if e.arg.App == "" && e.arg.Sha == "" && e.arg.Env == "" {
		// list all containers
		return AuthorizeSuperUser(&e.arg.ManagerAuthArg)
	}
	return AuthorizeApp(&e.arg.ManagerAuthArg, e.arg.App)
}

func (m *ManagerRPC) ListContainers(arg ManagerListContainersArg, reply *ManagerListContainersReply) error {
	return NewTask("ListContainers", &ListContainersExecutor{arg, reply}).Run()
}

type ListEnvsExecutor struct {
	arg   ManagerListEnvsArg
	reply *ManagerListEnvsReply
}

func (e *ListEnvsExecutor) Request() interface{} {
	return e.arg
}

func (e *ListEnvsExecutor) Result() interface{} {
	return e.reply
}

func (e *ListEnvsExecutor) Description() string {
	if e.arg.App == "" || e.arg.Sha == "" {
		return fmt.Sprintf("%s @ %s", e.arg.App, e.arg.Sha)
	}
	return "All"
}

func (e *ListEnvsExecutor) Execute(t *Task) error {
	var err error
	if e.arg.App == "" || e.arg.Sha == "" {
		e.reply.Envs, err = datamodel.ListEnvs()
	} else {
		e.reply.Envs, err = datamodel.ListAppEnvs(e.arg.App, e.arg.Sha)
	}
	if err != nil {
		e.reply.Status = StatusError
	} else {
		e.reply.Status = StatusOk
	}
	return err
}

func (e *ListEnvsExecutor) Authorize() error {
	if e.arg.App == "" {
		return SimpleAuthorize(&e.arg.ManagerAuthArg)
	}
	return AuthorizeApp(&e.arg.ManagerAuthArg, e.arg.App)
}

func (m *ManagerRPC) ListEnvs(arg ManagerListEnvsArg, reply *ManagerListEnvsReply) error {
	return NewTask("ListEnvs", &ListEnvsExecutor{arg, reply}).Run()
}

type ListShasExecutor struct {
	arg   ManagerListShasArg
	reply *ManagerListShasReply
}

func (e *ListShasExecutor) Request() interface{} {
	return e.arg
}

func (e *ListShasExecutor) Result() interface{} {
	return e.reply
}

func (e *ListShasExecutor) Description() string {
	return e.arg.App
}

func (e *ListShasExecutor) Execute(t *Task) error {
	var err error
	if e.arg.App == "" {
		return errors.New("App is empty")
	}
	e.reply.Shas, err = datamodel.ListShas(e.arg.App)
	if err != nil {
		e.reply.Status = StatusError
	} else {
		e.reply.Status = StatusOk
	}
	return err
}
func (e *ListShasExecutor) Authorize() error {
	return AuthorizeApp(&e.arg.ManagerAuthArg, e.arg.App)
}

func (m *ManagerRPC) ListShas(arg ManagerListShasArg, reply *ManagerListShasReply) error {
	return NewTask("ListShas", &ListShasExecutor{arg, reply}).Run()
}

type ListAppsExecutor struct {
	arg   ManagerListAppsArg
	reply *ManagerListAppsReply
}

func (e *ListAppsExecutor) Request() interface{} {
	return e.arg
}

func (e *ListAppsExecutor) Result() interface{} {
	return e.reply
}

func (e *ListAppsExecutor) Description() string {
	return "ListApps"
}

func (e *ListAppsExecutor) Execute(t *Task) error {
	var err error
	apps, err := datamodel.ListApps()
	err = AuthorizeSuperUser(&e.arg.ManagerAuthArg)
	if err == nil {
		e.reply.Apps = apps
	} else {
		allowedApps := GetAllowedApps(&e.arg.ManagerAuthArg)
		appsCount := len(allowedApps)
		totalAppsCount := len(apps)
		e.reply.Apps = make([]string, 0, appsCount)
		for i := 0; i < totalAppsCount; i++ {
			if _, ok := allowedApps[apps[i]]; ok {
				e.reply.Apps = append(e.reply.Apps, apps[i])
			}
		}
		err = nil
	}
	if err != nil {
		e.reply.Status = StatusError
	} else {
		e.reply.Status = StatusOk
	}
	return err
}

func (e *ListAppsExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

func (m *ManagerRPC) ListApps(arg ManagerListAppsArg, reply *ManagerListAppsReply) error {
	return NewTask("ListApps", &ListAppsExecutor{arg, reply}).Run()
}
