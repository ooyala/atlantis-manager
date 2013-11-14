package rpc

import (
	. "atlantis/common"
	"errors"
)

func (o *Manager) Status(id string, status *TaskStatus) error {
	if id == "" {
		return errors.New("Id empty")
	}
	getStatus, getError := Tracker.Status(id)
	if getStatus == nil {
		*status = *TaskStatusUnknown
	} else {
		*status = *getStatus
	}
	return getError
}
