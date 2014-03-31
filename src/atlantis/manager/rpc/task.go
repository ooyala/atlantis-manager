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
	. "atlantis/manager/rpc/types"
	"errors"
	"sort"
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
	sort.Strings(*ids)
	return nil
}
