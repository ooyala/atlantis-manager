/* Copyright 2014 Ooyala, Inc. All rights reserved.
 *
 * This file is licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
 * except in compliance with the License. You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is
 * distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

package rpc

import (
	bman "atlantis/builder/manifest"
	. "atlantis/common"
	"atlantis/manager/builder"
	"atlantis/manager/datamodel"
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
	return fmt.Sprintf("["+e.arg.ManagerAuthArg.User+"] %s @ %s in %s", e.arg.App, e.arg.Sha, e.arg.Env)
}

func (e *DeployExecutor) Authorize() error {
	if err := checkRole("deploys", "write"); err != nil {
		return err
	}
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
	data, err := bman.Read(manifestReader)
	if err != nil {
		return err
	}
	manifest, err := CreateManifest(data)
	if err != nil {
		return err
	}
	if e.arg.App != manifest.Name {
		// NOTE(edanaher): If we kick off two jobs simultaneously, they will assume they have the same job id, so
		// one of them will get the manifest from the other one's Jenkins job.  Unfortunately, Jenkins doesn't
		// give back any sort of useful information when you create the job, so we can't just use an ID easily.
		// Moreover, the API may not in any way acknowledge the job for several seconds, meaning that we can't
		// even scan the jobs created during a brief time interval for the job in question.  Rather, we have to
		// poll Jenkins until we find the job we're looking for.  This is sufficiently terrible that I'm just
		// erroring out and blaming Jenkins rather than adding a giant pile of code to handle that case.  We could
		// retry the job ourself, but after a day trying to beat Jenkins into submission, I have no interest in
		// applying further hacks.
		return errors.New("The app name you specified does not match the manifest.  This is probably due to an unavoidable race condition in Jenkin's RESTless API.  Please try again.")
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

func (m *ManagerRPC) Deploy(arg ManagerDeployArg, reply *AsyncReply) error {
	return NewTask("Deploy", &DeployExecutor{arg, &ManagerDeployReply{}}).RunAsync(reply)
}

type DeployContainerExecutor struct {
	arg   ManagerDeployContainerArg
	reply *ManagerDeployReply
}

func (e *DeployContainerExecutor) Request() interface{} {
	return e.arg
}

func (e *DeployContainerExecutor) Result() interface{} {
	return e.reply
}

func (e *DeployContainerExecutor) Description() string {
	return fmt.Sprintf("["+e.arg.ManagerAuthArg.User+"] %s x%d", e.arg.ContainerID, e.arg.Instances)
}

func (e *DeployContainerExecutor) Authorize() error {
	if err := checkRole("deploys", "write"); err != nil {
		return err
	}
	// app is authorized in deploy()
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

func (e *DeployContainerExecutor) Execute(t *Task) error {
	if e.arg.ContainerID == "" {
		return errors.New("Container ID is empty")
	}
	if e.arg.Instances <= 0 {
		return errors.New("Instances should be > 0")
	}
	instance, err := datamodel.GetInstance(e.arg.ContainerID)
	if err != nil {
		return err
	}
	ihReply, err := supervisor.Get(instance.Host, instance.ID)
	if err != nil {
		return err
	}
	e.reply.Containers, err = deployContainer(&e.arg.ManagerAuthArg, ihReply.Container, e.arg.Instances, t)
	return err
}

func (m *ManagerRPC) DeployContainer(arg ManagerDeployContainerArg, reply *AsyncReply) error {
	return NewTask("DeployContainer", &DeployContainerExecutor{arg, &ManagerDeployReply{}}).RunAsync(reply)
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
	return fmt.Sprintf("["+e.arg.ManagerAuthArg.User+"] %s -> %s", e.arg.ContainerID, e.arg.ToHost)
}

func (e *CopyContainerExecutor) Authorize() error {
	if err := checkRole("deploys", "write"); err != nil {
		return err
	}
	// app is authorized in deploy()
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

func (e *CopyContainerExecutor) Execute(t *Task) error {
	if e.arg.ContainerID == "" {
		return errors.New("Container ID is empty")
	}
	if e.arg.ToHost == "" {
		return errors.New("To Host is empty")
	}
	cont, err := copyContainer(&e.arg.ManagerAuthArg, e.arg.ContainerID, e.arg.ToHost, t)
	if err != nil {
		return err
	}
	e.reply.Containers = []*Container{cont}
	switch e.arg.PostCopy {
	case PostCopyCleanup:
		// we want to only cleanup ZK
		// get old instance
		inst, err := datamodel.GetInstance(e.arg.ContainerID)
		if err != nil {
			return err
		}
		cleanupZk(inst, t)
	case PostCopyTeardown:
		// now time to remove the old instance
		// get old instance
		inst, err := datamodel.GetInstance(e.arg.ContainerID)
		if err != nil {
			return err
		}
		// get old container from supervisor
		ihReply, err := supervisor.Get(inst.Host, inst.ID)
		if err != nil {
			return err
		}
		cleanup(true, []*Container{ihReply.Container}, t)
	}
	return nil
}

func (m *ManagerRPC) CopyContainer(arg ManagerCopyContainerArg, reply *AsyncReply) error {
	return NewTask("CopyContainer", &CopyContainerExecutor{arg, &ManagerDeployReply{}}).RunAsync(reply)
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
	return fmt.Sprintf("["+e.arg.ManagerAuthArg.User+"] %s in %s -> %v", e.arg.App, e.arg.Env, e.arg.DepNames)
}

func (e *ResolveDepsExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

func (e *ResolveDepsExecutor) Execute(t *Task) error {
	zkEnv, err := datamodel.GetEnv(e.arg.Env)
	if err != nil {
		return errors.New("Environment Error: " + err.Error())
	}
	e.reply.Deps, err = ResolveDepValues(e.arg.App, zkEnv, e.arg.DepNames, false, t)
	if err != nil {
		return err
	}
	e.reply.Status = StatusOk
	return nil
}

func (m *ManagerRPC) ResolveDeps(arg ManagerResolveDepsArg, reply *ManagerResolveDepsReply) error {
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
	return fmt.Sprintf("["+e.arg.ManagerAuthArg.User+"] app: %s, sha: %s, env: %s, container: %s (all:%t)",
		e.arg.App, e.arg.Sha, e.arg.Env, e.arg.ContainerID, e.arg.All)
}

func (e *TeardownExecutor) Authorize() error {
	if err := checkRole("deploys", "write"); err != nil {
		return err
	}
	if e.arg.All {
		return AuthorizeSuperUser(&e.arg.ManagerAuthArg)
	}
	if e.arg.App == "" {
		return SimpleAuthorize(&e.arg.ManagerAuthArg)
	}
	return AuthorizeApp(&e.arg.ManagerAuthArg, e.arg.App)
}

func (e *TeardownExecutor) Execute(t *Task) error {
	hostMap, err := getContainerIDsToTeardown(t, e.arg)
	if err != nil {
		return err
	}
	if e.arg.All {
		tl := datamodel.NewTeardownLock(t.ID)
		if err := tl.Lock(); err != nil {
			return err
		}
		defer tl.Unlock()
	} else if e.arg.Env != "" {
		tl := datamodel.NewTeardownLock(t.ID, e.arg.App, e.arg.Sha, e.arg.Env)
		if err := tl.Lock(); err != nil {
			return err
		}
		defer tl.Unlock()
	} else if e.arg.Sha != "" {
		tl := datamodel.NewTeardownLock(t.ID, e.arg.App, e.arg.Sha)
		if err := tl.Lock(); err != nil {
			return err
		}
		defer tl.Unlock()
	} else if e.arg.App != "" {
		tl := datamodel.NewTeardownLock(t.ID, e.arg.App)
		if err := tl.Lock(); err != nil {
			return err
		}
		defer tl.Unlock()
	}
	tornContainers := []string{}
	for host, containerIDs := range hostMap {
		if e.arg.All {
			t.LogStatus("Tearing Down * from %s", host)
		} else {
			t.LogStatus("Tearing Down %v from %s", containerIDs, host)
		}

		ihReply, err := supervisor.Teardown(host, containerIDs, e.arg.All)
		if err != nil {
			return errors.New(fmt.Sprintf("Error Tearing Down %v from %s : %s", containerIDs, host,
				err.Error()))
		}
		tornContainers = append(tornContainers, ihReply.ContainerIDs...)
		for _, tornContainerID := range tornContainers {
			err := datamodel.DeleteFromPool([]string{tornContainerID})
			if err != nil {
				t.Log("Error removing %s from pool: %v", tornContainerID, err)
			}
			datamodel.Supervisor(host).RemoveContainer(tornContainerID)
			instance, err := datamodel.GetInstance(tornContainerID)
			if err != nil {
				continue
			}
			last, _ := instance.Delete()
			if last {
				DeleteAppShaFromEnv(instance.App, instance.Sha, instance.Env)
			}
		}
	}
	e.reply.ContainerIDs = tornContainers
	return nil
}

func (m *ManagerRPC) Teardown(arg ManagerTeardownArg, reply *AsyncReply) error {
	return NewTask("Teardown", &TeardownExecutor{arg, &ManagerTeardownReply{}}).RunAsync(reply)
}

func (m *ManagerRPC) DeployResult(id string, result *ManagerDeployReply) error {
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

func (m *ManagerRPC) TeardownResult(id string, result *ManagerTeardownReply) error {
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
