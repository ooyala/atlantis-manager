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

package status

import (
	"atlantis/manager/datamodel"
	. "atlantis/manager/rpc/types"
	"atlantis/manager/supervisor"
)

func GetUsage() (map[string]*SupervisorUsage, error) {
	// for each supervisor
	//   get health check to figure out total CPUShares, Memory, Price
	//   get list of containers
	//   fill in data in SupervisorUsage
	supers, err := datamodel.ListSupervisors()
	if err != nil {
		return nil, err
	}
	usageMap := map[string]*SupervisorUsage{}
	for _, super := range supers {
		usage := &SupervisorUsage{Containers: map[string]*ContainerUsage{}}
		hreply, err := supervisor.HealthCheck(super)
		if err != nil {
			return nil, err
		}
		usage.Price = hreply.Price
		usage.TotalContainers = hreply.Containers.Total
		usage.TotalCPUShares = hreply.CPUShares.Total
		usage.TotalMemory = hreply.Memory.Total
		lreply, err := supervisor.List(super)
		if err != nil {
			return nil, err
		}
		for id, cont := range lreply.Containers {
			usage.Containers[id] = &ContainerUsage{
				ID:        id,
				App:       cont.App,
				Sha:       cont.Sha,
				Env:       cont.Env,
				CPUShares: cont.Manifest.CPUShares,
				Memory:    cont.Manifest.MemoryLimit,
			}
		}
		usageMap[super] = usage
	}
	return usageMap, nil
}
