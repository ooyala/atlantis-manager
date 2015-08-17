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

package datamodel

import (
	"atlantis/manager/helper"
	"fmt"
	zookeeper "github.com/ghao-ooyala/gozk-recipes"
	"strings"
)

const allPath = "/"

type LockConflictError string

func (e LockConflictError) Error() string {
	return "Lock Conflict with: " + string(e)
}

func getLockedPaths(path string) (map[string]string, error) {
	var lockedPaths map[string]string
	if err := getJson(path, &lockedPaths); err != nil {
		return nil, err
	}
	if lockedPaths == nil {
		lockedPaths = map[string]string{}
	}
	return lockedPaths, nil
}

type DeployLock struct {
	id     string
	path   string
	locked bool
}

func NewDeployLock(id, app, sha, env string) *DeployLock {
	return &DeployLock{id: id, path: fmt.Sprintf("/%s/%s/%s", app, sha, env)}
}

func (l *DeployLock) Lock() error {
	if l.locked {
		return nil
	}
	// check if we can lock by checking if any thing in the lock file is a prefix to us
	path := helper.GetBaseLockPath("deploy")
	mutex := zookeeper.NewMutex(Zk.Conn, path)
	if err := mutex.Lock(); err != nil {
		return err
	}
	defer mutex.Unlock()
	lockedPaths, err := getLockedPaths(path)
	if err != nil {
		return err
	}
	if allID, ok := lockedPaths[allPath]; ok && allID != "" {
		return LockConflictError(allID)
	}
	for p, conflictID := range lockedPaths {
		if strings.HasPrefix(l.path, p) {
			return LockConflictError(conflictID)
		}
	}
	// if no conflicts, register our lock
	lockedPaths[l.path] = l.id
	if err := setJson(path, lockedPaths); err != nil {
		return err
	}
	l.locked = true
	return nil
}

func (l *DeployLock) Unlock() error {
	if !l.locked {
		return nil
	}
	// remove ourselves from the lock file
	path := helper.GetBaseLockPath("deploy")
	mutex := zookeeper.NewMutex(Zk.Conn, path)
	if err := mutex.Lock(); err != nil {
		return err
	}
	defer mutex.Unlock()
	lockedPaths, err := getLockedPaths(path)
	if err != nil {
		return err
	}
	delete(lockedPaths, l.path)
	if err := setJson(path, lockedPaths); err != nil {
		return err
	}
	l.locked = false
	return nil
}

type TeardownLock struct {
	id     string
	path   string
	locked bool
}

func NewTeardownLock(id string, args ...string) *TeardownLock {
	path := allPath
	if len(args) > 0 {
		path = helper.JoinWithBase("/", args...)
	}
	return &TeardownLock{id: id, path: path}
}

func (l *TeardownLock) Lock() error {
	if l.locked {
		return nil
	}
	// check if we can lock by checking if we are a prefix to anything in the lock file
	path := helper.GetBaseLockPath("deploy")
	mutex := zookeeper.NewMutex(Zk.Conn, path)
	if err := mutex.Lock(); err != nil {
		return err
	}
	defer mutex.Unlock()
	lockedPaths, err := getLockedPaths(path)
	if err != nil {
		return err
	}
	if allID, ok := lockedPaths[allPath]; ok && allID != "" {
		return LockConflictError(allID)
	}
	for p, conflictID := range lockedPaths {
		if strings.HasPrefix(p, l.path) {
			return LockConflictError(conflictID)
		}
	}
	// if no conflicts, register our lock
	lockedPaths[l.path] = l.id
	if err := setJson(path, lockedPaths); err != nil {
		return err
	}
	l.locked = true
	return nil
}

func (l *TeardownLock) Unlock() error {
	if !l.locked {
		return nil
	}
	// remove ourselves from the lock file
	path := helper.GetBaseLockPath("deploy")
	mutex := zookeeper.NewMutex(Zk.Conn, path)
	if err := mutex.Lock(); err != nil {
		return err
	}
	defer mutex.Unlock()
	lockedPaths, err := getLockedPaths(path)
	if err != nil {
		return err
	}
	delete(lockedPaths, l.path)
	if err := setJson(path, lockedPaths); err != nil {
		return err
	}
	l.locked = false
	return nil
}

func NewRouterPortsLock(internal bool) *RouterPortsLock {
	return &RouterPortsLock{internal: internal}
}

type RouterPortsLock struct {
	internal bool
	locked   bool
	mutex    *zookeeper.Mutex
}

func (l *RouterPortsLock) Lock() error {
	var path string
	if l.locked {
		return nil
	}
	if l.internal {
		path = helper.GetBaseLockPath("router_ports_internal")
	} else {
		path = helper.GetBaseLockPath("router_ports_external")
	}
	l.mutex = zookeeper.NewMutex(Zk.Conn, path)
	if err := l.mutex.Lock(); err != nil {
		return err
	}
	l.locked = true
	return nil
}

func (l *RouterPortsLock) Unlock() error {
	if !l.locked {
		return nil
	}
	if err := l.mutex.Unlock(); err != nil {
		return err
	}
	l.locked = false
	return nil
}
