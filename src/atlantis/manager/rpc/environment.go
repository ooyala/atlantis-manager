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
	. "atlantis/common"
	"atlantis/manager/datamodel"
	. "atlantis/manager/rpc/types"
	"errors"
	"fmt"
	"sync"
)

var (
	appShas  = map[string]map[string]int{} // env -> app+sha -> # deployed
	envMutex = &sync.RWMutex{}
)

func LoadEnvs() error {
	envMutex.Lock()
	// get app+sha -> env
	apps, err := datamodel.ListApps()
	if err != nil {
		envMutex.Unlock()
		return nil // if we can't get the apps, this means this is a fresh zookeeper. don't fail.
	}
	for _, app := range apps {
		shas, err := datamodel.ListShas(app)
		if err != nil {
			envMutex.Unlock()
			return err
		}
		for _, sha := range shas {
			appSha := app + sha
			envs, err := datamodel.ListAppEnvs(app, sha)
			if err != nil {
				envMutex.Unlock()
				return err
			}
			for _, env := range envs {
				if _, ok := appShas[env]; ok {
					appShas[env][appSha]++
				} else {
					appShas[env] = map[string]int{appSha: 1}
				}
			}
		}
	}
	envMutex.Unlock()
	return nil
}

func IsEnvInUse(env string) bool {
	envMutex.RLock()
	inUse := false
	if envAppShas, ok := appShas[env]; ok {
		inUse = len(envAppShas) > 0
	}
	envMutex.RUnlock()
	return inUse
}

func AddAppShaToEnv(app, sha, env string) {
	appSha := app + sha
	envMutex.Lock()
	if _, ok := appShas[env]; ok {
		appShas[env][appSha]++
	} else {
		appShas[env] = map[string]int{appSha: 1}
	}
	envMutex.Unlock()
}

func DeleteAppShaFromEnv(app, sha, env string) {
	appSha := app + sha
	envMutex.Lock()
	appShas[env][appSha]--
	if appShas[env][appSha] == 0 {
		delete(appShas[env], appSha)
	}
	envMutex.Unlock()
}

func DeleteEnv(env string) {
	envMutex.Lock()
	delete(appShas, env)
	envMutex.Unlock()
}

// ----------------------------------------------------------------------------------------------------------
// Update, Delete, Get Environments
// ----------------------------------------------------------------------------------------------------------

type UpdateEnvExecutor struct {
	arg   ManagerEnvArg
	reply *ManagerEnvReply
}

func (e *UpdateEnvExecutor) Request() interface{} {
	return e.arg
}

func (e *UpdateEnvExecutor) Result() interface{} {
	return e.reply
}

func (e *UpdateEnvExecutor) Description() string {
	return fmt.Sprintf("%s", e.arg.Name)
}

func (e *UpdateEnvExecutor) Execute(t *Task) error {
	// error checking
	if e.arg.Name == "" {
		return errors.New("Please specify a name")
	}
	if IsEnvInUse(e.arg.Name) {
		return errors.New(fmt.Sprintf("%s is in use and cannot be updated", e.arg.Name))
	}
	env := datamodel.Env(e.arg.Name)
	if err := env.Save(); err != nil {
		return err
	}
	e.reply.Status = StatusOk
	return nil
}

func (e *UpdateEnvExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

type DeleteEnvExecutor struct {
	arg   ManagerEnvArg
	reply *ManagerEnvReply
}

func (e *DeleteEnvExecutor) Request() interface{} {
	return e.arg
}

func (e *DeleteEnvExecutor) Result() interface{} {
	return e.reply
}

func (e *DeleteEnvExecutor) Description() string {
	return fmt.Sprintf("%s", e.arg.Name)
}

func (e *DeleteEnvExecutor) Execute(t *Task) (err error) {
	// error checking
	if e.arg.Name == "" {
		return errors.New("Please specify a name")
	}
	if IsEnvInUse(e.arg.Name) {
		return errors.New(fmt.Sprintf("%s is in use and cannot be deleted", e.arg.Name))
	}
	env := datamodel.Env(e.arg.Name)
	if err := env.Delete(); err != nil {
		e.reply.Status = StatusError
	} else {
		DeleteEnv(e.arg.Name)
		e.reply.Status = StatusOk
	}
	return err
}

func (e *DeleteEnvExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

func (m *ManagerRPC) UpdateEnv(arg ManagerEnvArg, reply *ManagerEnvReply) error {
	return NewTask("UpdateEnv", &UpdateEnvExecutor{arg, reply}).Run()
}

func (m *ManagerRPC) DeleteEnv(arg ManagerEnvArg, reply *ManagerEnvReply) error {
	return NewTask("DeleteEnv", &DeleteEnvExecutor{arg, reply}).Run()
}
