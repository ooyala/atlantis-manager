package rpc

import (
	. "atlantis/common"
	"atlantis/crypto"
	"atlantis/manager/datamodel"
	. "atlantis/manager/rpc/types"
	"errors"
	"fmt"
	"sync"
)

// NOTE[jigish]: We don't support links in dependencies. If we wanted to add it this is what you'd have to do:
// 1. create a new map that is parent -> child -> # of links that use this relation
// 2. create new maintenance methods that add/remove from the map
// 3. update LoadEnvs()
// 4. update isEnvInUseRecursive()
// 5. update Update/Delete Dep to keep track of the new links

var (
	children = map[string]map[string]bool{} // parent -> child -> true
	appShas  = map[string]map[string]int{}  // env -> app+sha -> # deployed
	envMutex = &sync.RWMutex{}
)

func LoadEnvs() error {
	envMutex.Lock()
	// get env -> children
	envs, err := datamodel.ListEnvs()
	if err != nil {
		envMutex.Unlock()
		return nil // if we can't get the envs, this means this is a fresh zookeeper. don't fail.
	}
	for _, name := range envs {
		env, err := datamodel.GetEnv(name)
		if err != nil {
			envMutex.Unlock()
			return err
		}
		if env.Parent != "" {
			if _, ok := children[env.Parent]; ok {
				children[env.Parent][env.Name] = true
			} else {
				children[env.Parent] = map[string]bool{env.Name: true}
			}
		}
		if _, ok := children[env.Name]; !ok {
			children[env.Name] = map[string]bool{}
		}
	}
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

func isEnvInUseRecursive(env string) bool {
	if envAppShas, ok := appShas[env]; ok {
		inUse := len(envAppShas) > 0
		if inUse {
			return true
		}
	}
	// check children (and their children, etc...)
	envChildren := children[env]
	if envChildren == nil || len(envChildren) == 0 {
		return false
	}
	for child, _ := range envChildren {
		if isEnvInUseRecursive(child) {
			return true
		}
	}
	return false
}

func IsEnvInUse(env string) bool {
	envMutex.RLock()
	inUse := isEnvInUseRecursive(env)
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

func HasChildren(parent string) bool {
	envMutex.RLock()
	envChildren, ok := children[parent]
	has := ok && len(envChildren) > 0
	envMutex.RUnlock()
	return has
}

func UpdateEnv(child, oldParent, newParent string) {
	envMutex.Lock()
	// first touch
	if _, ok := children[child]; !ok {
		children[child] = map[string]bool{}
	}
	// next delete child from old parent
	if oldParent != "" {
		delete(children[oldParent], child)
	}
	// finally add child to new parent
	if newParent != "" {
		children[newParent][child] = true
	}
	envMutex.Unlock()
}

func DeleteEnv(env string) {
	envMutex.Lock()
	// delete its children tracker
	delete(children, env)
	// delete wherever it is being tracked as a child
	for parent, envChildren := range children {
		if _, ok := envChildren[env]; ok {
			delete(children[parent], env)
		}
	}
	envMutex.Unlock()
}

func decryptAll(deps map[string]string) map[string]string {
	for key, encryptedValue := range deps {
		deps[key] = string(crypto.Decrypt([]byte(encryptedValue)))
	}
	return deps
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
	return fmt.Sprintf("%s, parent: %s", e.arg.Name, e.arg.Parent)
}

func (e *UpdateEnvExecutor) Execute(t *Task) error {
	// error checking
	if e.arg.Name == "" {
		return errors.New("Please specify a name")
	}
	if IsEnvInUse(e.arg.Name) {
		return errors.New(fmt.Sprintf("%s is in use and cannot be updated", e.arg.Name))
	}
	if e.arg.Parent != "" {
		_, err := datamodel.GetEnv(e.arg.Parent)
		if err != nil {
			return err
		}
	}
	gEnv, err := datamodel.GetEnv(e.arg.Name)
	if err != nil {
		gEnv = &datamodel.ZkEnv{}
	}
	env := datamodel.Env(e.arg.Name, e.arg.Parent)
	if err := env.Save(); err != nil {
		return err
	}
	UpdateEnv(e.arg.Name, gEnv.Parent, env.Parent)
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
	if HasChildren(e.arg.Name) {
		return errors.New(fmt.Sprintf("%s has children and cannot be deleted", e.arg.Name))
	}
	env := datamodel.Env(e.arg.Name, "")
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

type GetEnvExecutor struct {
	arg   ManagerEnvArg
	reply *ManagerEnvReply
}

func (e *GetEnvExecutor) Request() interface{} {
	return e.arg
}

func (e *GetEnvExecutor) Result() interface{} {
	return e.reply
}

func (e *GetEnvExecutor) Description() string {
	return fmt.Sprintf("%s", e.arg.Name)
}

func (e *GetEnvExecutor) Execute(t *Task) (err error) {
	// error checking
	if e.arg.Name == "" {
		return errors.New("Please specify a name")
	}
	env, err := datamodel.GetEnv(e.arg.Name)
	if err != nil {
		return err
	}
	e.reply.Parent = env.Parent
	if deps, err := env.AllDepValues(); err != nil {
		return err
	} else {
		e.reply.Deps = decryptAll(deps)
	}
	if resolvedDeps, err := env.ResolveAllDepValues(); err != nil {
		return err
	} else {
		e.reply.ResolvedDeps = decryptAll(resolvedDeps)
	}
	e.reply.Status = StatusOk
	return nil
}

func (e *GetEnvExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

func (o *Manager) UpdateEnv(arg ManagerEnvArg, reply *ManagerEnvReply) error {
	return NewTask("UpdateEnv", &UpdateEnvExecutor{arg, reply}).Run()
}

func (o *Manager) DeleteEnv(arg ManagerEnvArg, reply *ManagerEnvReply) error {
	return NewTask("DeleteEnv", &DeleteEnvExecutor{arg, reply}).Run()
}

func (o *Manager) GetEnv(arg ManagerEnvArg, reply *ManagerEnvReply) error {
	return NewTask("GetEnv", &GetEnvExecutor{arg, reply}).Run()
}

// ----------------------------------------------------------------------------------------------------------
// Update, Delete, Get Dependencies
// ----------------------------------------------------------------------------------------------------------

type UpdateDepExecutor struct {
	arg   ManagerDepArg
	reply *ManagerDepReply
}

func (e *UpdateDepExecutor) Request() interface{} {
	return e.arg
}

func (e *UpdateDepExecutor) Result() interface{} {
	return e.reply
}

func (e *UpdateDepExecutor) Description() string {
	return fmt.Sprintf("%s in %s", e.arg.Name, e.arg.Env)
}

func (e *UpdateDepExecutor) Execute(t *Task) (err error) {
	// error checking
	if e.arg.Name == "" {
		return errors.New("Please specify a name")
	}
	if e.arg.Value == "" {
		return errors.New("Please specify a value")
	}
	if e.arg.Env == "" {
		return errors.New("Please specify an environment")
	}
	if _, err := datamodel.GetApp(e.arg.Name); err == nil {
		return errors.New("Dep name conflicts with a registered app")
	}
	env, err := datamodel.GetEnv(e.arg.Env)
	if err != nil {
		return err
	}
	if _, err := env.GetDepValue(e.arg.Name); err == nil {
		// trying to update dep that already exists
		if IsEnvInUse(e.arg.Name) {
			return errors.New(fmt.Sprintf("%s is in use and cannot be updated", e.arg.Name))
		}
	}
	if err = env.UpdateDep(e.arg.Name, string(crypto.Encrypt([]byte(e.arg.Value)))); err != nil {
		return err
	}
	if value, err := env.GetDepValue(e.arg.Name); err == nil {
		e.reply.Value = string(crypto.Decrypt([]byte(value)))
	} else {
		return err
	}
	e.reply.Status = StatusOk
	return nil
}

func (e *UpdateDepExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

type DeleteDepExecutor struct {
	arg   ManagerDepArg
	reply *ManagerDepReply
}

func (e *DeleteDepExecutor) Request() interface{} {
	return e.arg
}

func (e *DeleteDepExecutor) Result() interface{} {
	return e.reply
}

func (e *DeleteDepExecutor) Description() string {
	return fmt.Sprintf("%s in %s", e.arg.Name, e.arg.Env)
}

func (e *DeleteDepExecutor) Execute(t *Task) (err error) {
	// error checking
	if e.arg.Name == "" {
		return errors.New("Please specify a name")
	}
	if e.arg.Env == "" {
		return errors.New("Please specify an environment")
	}
	if IsEnvInUse(e.arg.Name) {
		return errors.New(fmt.Sprintf("%s is in use and cannot be deleted", e.arg.Name))
	}
	env, err := datamodel.GetEnv(e.arg.Env)
	if err != nil {
		return err
	}
	if value, err := env.GetDepValue(e.arg.Name); err == nil {
		e.reply.Value = string(crypto.Decrypt([]byte(value)))
	} else {
		return err
	}
	if err := env.DeleteDep(e.arg.Name); err != nil {
		e.reply.Status = StatusError
	} else {
		e.reply.Status = StatusOk
	}
	return err
}

func (e *DeleteDepExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

type GetDepExecutor struct {
	arg   ManagerDepArg
	reply *ManagerDepReply
}

func (e *GetDepExecutor) Request() interface{} {
	return e.arg
}

func (e *GetDepExecutor) Result() interface{} {
	return e.reply
}

func (e *GetDepExecutor) Description() string {
	return fmt.Sprintf("%s in %s", e.arg.Name, e.arg.Env)
}

func (e *GetDepExecutor) Execute(t *Task) (err error) {
	// error checking
	if e.arg.Name == "" {
		return errors.New("Please specify a name")
	}
	if e.arg.Env == "" {
		return errors.New("Please specify an environment")
	}
	env, err := datamodel.GetEnv(e.arg.Env)
	if err != nil {
		return err
	}
	if value, err := env.GetDepValue(e.arg.Name); err != nil {
		return err
	} else {
		e.reply.Value = string(crypto.Decrypt([]byte(value)))
	}
	e.reply.Status = StatusOk
	return nil
}

func (e *GetDepExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

func (o *Manager) UpdateDep(arg ManagerDepArg, reply *ManagerDepReply) error {
	return NewTask("UpdateDep", &UpdateDepExecutor{arg, reply}).Run()
}

func (o *Manager) DeleteDep(arg ManagerDepArg, reply *ManagerDepReply) error {
	return NewTask("DeleteDep", &DeleteDepExecutor{arg, reply}).Run()
}

func (o *Manager) GetDep(arg ManagerDepArg, reply *ManagerDepReply) error {
	return NewTask("GetDep", &GetDepExecutor{arg, reply}).Run()
}
