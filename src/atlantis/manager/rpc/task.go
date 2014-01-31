package rpc

import (
	. "atlantis/common"
	. "atlantis/manager/rpc/types"
	"errors"
)

func (m *ManagerRPC) Status(id string, status *TaskStatus) error {
	if id == "" {
		return errors.New("ID empty")
	}
	getStatus, getError := Tracker.Status(id)
	if getStatus == nil {
		*status = *TaskStatusUnknown
	} else {
		*status = *getStatus
	}
	return getError
}

func (m *ManagerRPC) ListTaskIDs(arg ManagerAuthArg, ids *[]string) error {
	if err := SimpleAuthorize(&arg); err != nil {
		return err
	}
	types := []string{"Deploy", "Teardown"}
	if AuthorizeSuperUser(&arg) == nil {
		// superuser, return all types
		types = append(types, []string{
			"RegisterRouter",
			"UnregisterRouter",
			"RegisterManager",
			"UnregisterManager",
			"RegisterSupervisor",
			"UnregisterSupervisor",
		}...)
	}
	*ids = Tracker.ListIDs(types)
	return nil
}
