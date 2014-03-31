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
	"atlantis/manager/rpc/types"
	"log"
)

type ZkIPGroup types.IPGroup

func GetIPGroup(name string) (zig *ZkIPGroup, err error) {
	zig = &ZkIPGroup{}
	err = getJson(helper.GetBaseIPGroupPath(name), zig)
	return
}

func (zig *ZkIPGroup) Delete() error {
	return Zk.RecursiveDelete(zig.path())
}

func (zig *ZkIPGroup) path() string {
	return helper.GetBaseIPGroupPath(zig.Name)
}

func (zig *ZkIPGroup) Save() error {
	if err := setJson(zig.path(), zig); err != nil {
		return err
	}
	return nil
}

func ListIPGroups() (groups []string, err error) {
	groups, _, err = Zk.VisibleChildren(helper.GetBaseIPGroupPath())
	if err != nil {
		log.Printf("Error getting list of ip groups. Error: %s.", err.Error())
	}
	if groups == nil {
		log.Println("No ip groups found. Returning empty list.")
		groups = []string{}
	}
	return
}
