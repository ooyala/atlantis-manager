package rpc

import (
	. "atlantis/common"
	"atlantis/manager/builder"
	"atlantis/manager/datamodel"
	"atlantis/manager/dns"
	. "atlantis/manager/rpc/types"
	"atlantis/manager/supervisor"
	. "atlantis/supervisor/rpc/types"
	"errors"
	"fmt"
)

type DeployExecutor struct {
	arg   ManagerDeployArg
	reply *ManagerDeployReply
}

func (e *DeployExecutor) Request() interface{} {
	return e.arg
}

func (e *DeployExecutor) Result() interface{} {
	return e.reply
}

func (e *DeployExecutor) Description() string {
	return fmt.Sprintf("%s @ %s in %s", e.arg.App, e.arg.Sha, e.arg.Env)
}

func (e *DeployExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

func (e *DeployExecutor) Execute(t *Task) error {
	// error checking
	if e.arg.App == "" {
		return errors.New("Please specify an app")
	}
	if e.arg.Sha == "" {
		return errors.New("Please specify a sha")
	}
	if e.arg.Env == "" {
		return errors.New("Please specify an environment")
	}
	if e.arg.CPUShares < 0 ||
		(e.arg.CPUShares > 0 && e.arg.CPUShares != 1 && e.arg.CPUShares%CPUSharesIncrement != 0) {
		return errors.New(fmt.Sprintf("CPU Shares should be 1 or a multiple of %d", CPUSharesIncrement))
	}
	if e.arg.MemoryLimit < 0 ||
		(e.arg.MemoryLimit > 0 && e.arg.MemoryLimit%MemoryLimitIncrement != 0) {
		return errors.New(fmt.Sprintf("Memory should be a multiple of %d", MemoryLimitIncrement))
	}
	// fetch the repo and root
	app, err := datamodel.GetApp(e.arg.App)
	if err != nil {
		return errors.New("App " + e.arg.App + " is not registered: " + err.Error())
	}
	// fetch and parse manifest for app name
	manifestReader, err := builder.DefaultBuilder.Build(t, app.Repo, app.Root, e.arg.Sha)
	if err != nil {
		return errors.New("Build Error: " + err.Error())
	}
	defer manifestReader.Close()
	t.LogStatus("Reading Manifest")
	manifest, err := ReadManifest(manifestReader)
	if err != nil {
		return err
	}
	if e.arg.App != manifest.Name {
		return errors.New("The app name you specified does not match the manifest")
	}
	if e.arg.CPUShares > 0 {
		manifest.CPUShares = e.arg.CPUShares
	}
	if e.arg.MemoryLimit > 0 {
		manifest.MemoryLimit = e.arg.MemoryLimit
	}
	// figure out how many instances we need
	if e.arg.Instances > 0 {
		manifest.Instances = e.arg.Instances
	} else if manifest.Instances == 0 {
		manifest.Instances = uint(1) // default to 1 instance
	}
	if e.arg.Dev {
		e.reply.Containers, err = devDeploy(&e.arg.ManagerAuthArg, manifest, e.arg.Sha, e.arg.Env, t)
	} else {
		e.reply.Containers, err = deploy(&e.arg.ManagerAuthArg, manifest, e.arg.Sha, e.arg.Env, t)
	}
	return err
}

func (o *Manager) Deploy(arg ManagerDeployArg, reply *AsyncReply) error {
	return NewTask("Deploy", &DeployExecutor{arg, &ManagerDeployReply{}}).RunAsync(reply)
}

type CopyContainerExecutor struct {
	arg   ManagerCopyContainerArg
	reply *ManagerDeployReply
}

func (e *CopyContainerExecutor) Request() interface{} {
	return e.arg
}

func (e *CopyContainerExecutor) Result() interface{} {
	return e.reply
}

func (e *CopyContainerExecutor) Description() string {
	return fmt.Sprintf("%s x%d", e.arg.ContainerId, e.arg.Instances)
}

func (e *CopyContainerExecutor) Authorize() error {
	// app is authorized in deploy()
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

func (e *CopyContainerExecutor) Execute(t *Task) error {
	if e.arg.ContainerId == "" {
		return errors.New("Container ID is empty")
	}
	if e.arg.Instances <= 0 {
		return errors.New("Instances should be > 0")
	}
	instance, err := datamodel.GetInstance(e.arg.ContainerId)
	if err != nil {
		return err
	}
	ihReply, err := supervisor.Get(instance.Host, instance.Id)
	if err != nil {
		return err
	}
	e.reply.Containers, err = deployContainer(&e.arg.ManagerAuthArg, ihReply.Container, e.arg.Instances, t)
	return err
}

func (o *Manager) CopyContainer(arg ManagerCopyContainerArg, reply *AsyncReply) error {
	return NewTask("CopyContainer", &CopyContainerExecutor{arg, &ManagerDeployReply{}}).RunAsync(reply)
}

type MoveContainerExecutor struct {
	arg   ManagerMoveContainerArg
	reply *ManagerDeployReply
}

func (e *MoveContainerExecutor) Request() interface{} {
	return e.arg
}

func (e *MoveContainerExecutor) Result() interface{} {
	return e.reply
}

func (e *MoveContainerExecutor) Description() string {
	return fmt.Sprintf("%s", e.arg.ContainerId)
}

func (e *MoveContainerExecutor) Authorize() error {
	// app is authorized in deploy()
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

func (e *MoveContainerExecutor) Execute(t *Task) error {
	if e.arg.ContainerId == "" {
		return errors.New("Container ID is empty")
	}
	instance, err := datamodel.GetInstance(e.arg.ContainerId)
	if err != nil {
		return err
	}
	ihReply, err := supervisor.Get(instance.Host, instance.Id)
	if err != nil {
		return err
	}
	cont, err := moveContainer(&e.arg.ManagerAuthArg, ihReply.Container, t)
	e.reply.Containers = []*Container{cont}
	return err
}

func (o *Manager) MoveContainer(arg ManagerMoveContainerArg, reply *AsyncReply) error {
	return NewTask("MoveContainer", &MoveContainerExecutor{arg, &ManagerDeployReply{}}).RunAsync(reply)
}

type ResolveDepsExecutor struct {
	arg   ManagerResolveDepsArg
	reply *ManagerResolveDepsReply
}

func (e *ResolveDepsExecutor) Request() interface{} {
	return e.arg
}

func (e *ResolveDepsExecutor) Result() interface{} {
	return e.reply
}

func (e *ResolveDepsExecutor) Description() string {
	return fmt.Sprintf("%s -> %v", e.arg.Env, e.arg.DepNames)
}

func (e *ResolveDepsExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

func (e *ResolveDepsExecutor) Execute(t *Task) error {
	zkEnv, err := datamodel.GetEnv(e.arg.Env)
	if err != nil {
		return errors.New("Environment Error: " + err.Error())
	}
	e.reply.Deps, err = ResolveDepValues(zkEnv, e.arg.DepNames, false)
	if err != nil {
		return err
	}
	e.reply.Status = StatusOk
	return nil
}

func (o *Manager) ResolveDeps(arg ManagerResolveDepsArg, reply *ManagerResolveDepsReply) error {
	return NewTask("ResolveDeps", &ResolveDepsExecutor{arg, reply}).Run()
}

type TeardownExecutor struct {
	arg   ManagerTeardownArg
	reply *ManagerTeardownReply
}

func (e *TeardownExecutor) Request() interface{} {
	return e.arg
}

func (e *TeardownExecutor) Result() interface{} {
	return e.reply
}

func (e *TeardownExecutor) Description() string {
	return fmt.Sprintf("app: %s, sha: %s, env: %s, container: %s (all:%t)", e.arg.App, e.arg.Sha, e.arg.Env,
		e.arg.ContainerId, e.arg.All)
}

func (e *TeardownExecutor) Authorize() error {
	if e.arg.All {
		return AuthorizeSuperUser(&e.arg.ManagerAuthArg)
	}
	if e.arg.App == "" {
		return SimpleAuthorize(&e.arg.ManagerAuthArg)
	}
	return AuthorizeApp(&e.arg.ManagerAuthArg, e.arg.App)
}

func (e *TeardownExecutor) Execute(t *Task) error {
	hostMap, err := getContainerIdsToTeardown(t, e.arg)
	if err != nil {
		return err
	}
	if e.arg.All {
		tl := datamodel.NewTeardownLock(t.Id)
		if err := tl.Lock(); err != nil {
			return err
		}
		defer tl.Unlock()
	} else if e.arg.Env != "" {
		tl := datamodel.NewTeardownLock(t.Id, e.arg.App, e.arg.Sha, e.arg.Env)
		if err := tl.Lock(); err != nil {
			return err
		}
		defer tl.Unlock()
	} else if e.arg.Sha != "" {
		tl := datamodel.NewTeardownLock(t.Id, e.arg.App, e.arg.Sha)
		if err := tl.Lock(); err != nil {
			return err
		}
		defer tl.Unlock()
	} else if e.arg.App != "" {
		tl := datamodel.NewTeardownLock(t.Id, e.arg.App)
		if err := tl.Lock(); err != nil {
			return err
		}
		defer tl.Unlock()
	}
	tornContainers := []string{}
	for host, containerIds := range hostMap {
		if e.arg.All {
			t.LogStatus("Tearing Down * from %s", host)
		} else {
			t.LogStatus("Tearing Down %v from %s", containerIds, host)
		}

		ihReply, err := supervisor.Teardown(host, containerIds, e.arg.All)
		if err != nil {
			return errors.New(fmt.Sprintf("Error Tearing Down %v from %s : %s", containerIds, host,
				err.Error()))
		}
		tornContainers = append(tornContainers, ihReply.ContainerIds...)
		for _, tornContainerId := range tornContainers {
			err := datamodel.DeleteFromPool([]string{tornContainerId})
			if err != nil {
				t.Log("Error removing %s from pool: %v", tornContainerId, err)
			}
			datamodel.Host(host).RemoveContainer(tornContainerId)
			instance, err := datamodel.GetInstance(tornContainerId)
			if err != nil {
				continue
			}
			last, _ := instance.Delete()
			if last {
				if instance.Internal {
					dns.DeleteAppAliases(instance.App, instance.Sha, instance.Env)
				}
				DeleteAppShaFromEnv(instance.App, instance.Sha, instance.Env)
			}
		}
	}
	e.reply.ContainerIds = tornContainers
	return nil
}

func (o *Manager) Teardown(arg ManagerTeardownArg, reply *AsyncReply) error {
	return NewTask("Teardown", &TeardownExecutor{arg, &ManagerTeardownReply{}}).RunAsync(reply)
}

func (o *Manager) DeployResult(id string, result *ManagerDeployReply) error {
	if id == "" {
		return errors.New("ID empty")
	}
	status, err := Tracker.Status(id)
	if status.Status == StatusUnknown {
		return errors.New("Unknown ID.")
	}
	if status.Name != "Deploy" {
		return errors.New("ID is not a Deploy.")
	}
	if !status.Done {
		return errors.New("Deploy isn't done.")
	}
	if status.Status == StatusError || err != nil {
		return err
	}
	getResult := Tracker.Result(id)
	switch r := getResult.(type) {
	case *ManagerDeployReply:
		*result = *r
	default:
		// this should never happen
		return errors.New("Invalid Result Type.")
	}
	return nil
}

func (o *Manager) TeardownResult(id string, result *ManagerTeardownReply) error {
	if id == "" {
		return errors.New("ID empty")
	}
	status, err := Tracker.Status(id)
	if status.Status == StatusUnknown {
		return errors.New("Unknown ID.")
	}
	if status.Name != "Teardown" {
		return errors.New("ID is not a Teardown.")
	}
	if !status.Done {
		return errors.New("Teardown isn't done.")
	}
	if status.Status == StatusError || err != nil {
		return err
	}
	getResult := Tracker.Result(id)
	switch r := getResult.(type) {
	case *ManagerTeardownReply:
		*result = *r
	default:
		// this should never happen
		return errors.New("Invalid Result Type.")
	}
	return nil
}
