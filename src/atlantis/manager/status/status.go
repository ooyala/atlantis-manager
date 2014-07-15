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
	"strconv"
)

func toPrice(price float64) float64 {
	p, _ := strconv.ParseFloat(strconv.FormatFloat(price, 'f', 2, 64), 64)
	return p
}

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
		usage.Host = super
		price := hreply.Price
		total_cpu := hreply.CPUShares.Total
		total_mem := hreply.Memory.Total
		usage.TotalPrice = price
		usage.TotalContainers = hreply.Containers.Total
		usage.TotalCPUShares = total_cpu
		usage.TotalMemory = total_mem
		lreply, err := supervisor.List(super)
		if err != nil {
			return nil, err
		}
		var conts uint = 0
		var cpu uint = 0
		var mem uint = 0
		var cpu_price float64 = 0.0
		var mem_price float64 = 0.0
		for id, cont := range lreply.Containers {
			conts += 1
			cpu += cont.Manifest.CPUShares
			mem += cont.Manifest.MemoryLimit
			c := price * (float64(cont.Manifest.CPUShares) / float64(total_cpu))
			m := price * (float64(cont.Manifest.MemoryLimit) / float64(total_mem))
			cpu_price += c
			mem_price += m
			usage.Containers[id] = &ContainerUsage{
				ID:        id,
				App:       cont.App,
				Sha:       cont.Sha,
				Env:       cont.Env,
				CPUShares: cont.Manifest.CPUShares,
				Memory:    cont.Manifest.MemoryLimit,
				CPUPrice:  toPrice(c),
				MemPrice:  toPrice(m),
			}
		}
		usage.UsedContainers = conts
		usage.UsedCPUShares = cpu
		usage.UsedMemory = mem
		usage.UsedCPUPrice = toPrice(cpu_price)
		usage.UsedMemPrice = toPrice(mem_price)
		usageMap[super] = usage
	}
	return usageMap, nil
}
