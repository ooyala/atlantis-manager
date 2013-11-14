package datamodel

import (
	"atlantis/manager/helper"
	"fmt"
	zookeeper "github.com/jigish/gozk-recipes"
	"strings"
)

const allPath = "/"

type LockConflictError string

func (e LockConflictError) Error() string {
	return fmt.Sprintf("Lock Conflict with: %s", e)
}

func (e LockConflictError) String() string {
	return fmt.Sprintf("Lock Conflict with: %s", e)
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
	if allId, ok := lockedPaths[allPath]; ok && allId != "" {
		return LockConflictError(allId)
	}
	for p, conflictId := range lockedPaths {
		if strings.HasPrefix(l.path, p) {
			return LockConflictError(conflictId)
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
	if allId, ok := lockedPaths[allPath]; ok && allId != "" {
		return LockConflictError(allId)
	}
	for p, conflictId := range lockedPaths {
		if strings.HasPrefix(p, l.path) {
			return LockConflictError(conflictId)
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
