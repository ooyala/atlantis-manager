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

package netsec

import (
	"atlantis/manager/datamodel"
	. "atlantis/manager/rpc/types"
	"atlantis/manager/supervisor"
	"sync"
)

const InternalRouterIPGroup = "internal-router"

var lock = sync.RWMutex{}

func updateSupervisors(name string, ips []string) error {
	// update all supervisors
	supers, err := datamodel.ListSupervisors()
	if err != nil {
		return err
	}
	for _, host := range supers {
		if _, err := supervisor.UpdateIPGroup(host, name, ips); err != nil {
			return err
		}
	}
	return nil
}

func UpdateIPGroup(name string, ips []string) error {
	lock.Lock()
	defer lock.Unlock()
	// save the IP group
	group := datamodel.ZkIPGroup{Name: name, IPs: ips}
	if err := group.Save(); err != nil {
		return err
	}
	return updateSupervisors(name, ips)
}

func AddIPToGroup(name, ip string) error {
	lock.Lock()
	defer lock.Unlock()
	group, err := datamodel.GetIPGroup(name)
	if err != nil {
		return err
	}
	// dedup
	ipsMap := map[string]bool{}
	for _, theIP := range group.IPs {
		ipsMap[theIP] = true
	}
	group.IPs = make([]string, len(ipsMap))
	i := 0
	for ip, _ := range ipsMap {
		group.IPs[i] = ip
		i++
	}

	if err := group.Save(); err != nil {
		return err
	}
	return updateSupervisors(name, group.IPs)
}

func RemoveIPFromGroup(name, ip string) error {
	lock.Lock()
	defer lock.Unlock()
	group, err := datamodel.GetIPGroup(name)
	if err != nil {
		return err
	}

	// dedup
	ipsMap := map[string]bool{}
	for _, theIP := range group.IPs {
		ipsMap[theIP] = true
	}
	delete(ipsMap, ip) // delete the ip we want to remove
	group.IPs = make([]string, len(ipsMap))
	i := 0
	for ip, _ := range ipsMap {
		group.IPs[i] = ip
		i++
	}

	if err := group.Save(); err != nil {
		return err
	}
	return updateSupervisors(name, group.IPs)
}

func DeleteIPGroup(name string) error {
	lock.Lock()
	defer lock.Unlock()
	// delete the IP group
	group, err := datamodel.GetIPGroup(name)
	if err != nil {
		return err
	}
	if err := group.Delete(); err != nil {
		return err
	}
	// update all supervisors
	supers, err := datamodel.ListSupervisors()
	if err != nil {
		return err
	}
	for _, host := range supers {
		if _, err := supervisor.DeleteIPGroup(host, name); err != nil {
			return err
		}
	}
	return nil
}

func GetIPGroup(name string) (*IPGroup, error) {
	lock.RLock()
	defer lock.RUnlock()
	// get the IP group
	group, err := datamodel.GetIPGroup(name)
	if err != nil {
		return nil, err
	}
	typedGroup := IPGroup(*group)
	return &typedGroup, nil
}