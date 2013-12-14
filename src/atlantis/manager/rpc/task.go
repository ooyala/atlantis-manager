package rpc

import (
	. "atlantis/common"
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
