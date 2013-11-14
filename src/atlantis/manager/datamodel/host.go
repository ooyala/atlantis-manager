package datamodel

import (
	. "atlantis/common"
	. "atlantis/manager/constant"
	"atlantis/manager/helper"
	"atlantis/manager/supervisor"
	"errors"
	"fmt"
	"log"
	"sort"
)

type ZkHost string

type HostData struct {
	PortMap map[string]uint16
}

func (h *HostData) HasAppShaEnv(app, sha, env string) bool {
	for container, _ := range h.PortMap {
		zi, err := GetInstance(container)
		if err != nil {
			continue
		}
		if zi.App == app && zi.Sha == sha && zi.Env == env {
			return true
		}
	}
	return false
}

func (h *HostData) CountAppShaEnv(app, sha, env string) int {
	count := 0
	for container, _ := range h.PortMap {
		zi, err := GetInstance(container)
		if err != nil {
			continue
		}
		if zi.App == app && zi.Sha == sha && zi.Env == env {
			count++
		}
	}
	return count
}

type ContainerData struct {
	port uint16
}

func Host(name string) ZkHost {
	return ZkHost(name)
}

func (h ZkHost) Touch() error {
	_, err := Zk.Touch(h.path())
	return err
}

// Delete the host node and all child container nodes of that host
func (h ZkHost) Delete() error {
	return Zk.RecursiveDelete(h.path())
}

// Supervisor will tell us the port -> container mapping for a given host, and we will store this back in zk in
// the /host/[host] node
func (h ZkHost) SetContainerAndPort(container string, port uint16) (err error) {
	err = h.createOrUpdateContainer(container, &ContainerData{port})
	if err != nil {
		log.Printf("Error setting mapping in container node. Error: %s.", err.Error())
		return
	}
	err = h.addRelation(container, port)
	if err != nil {
		log.Printf("Error setting mapping in host node. Error: %s.", err.Error())
	}
	return
}

func (h ZkHost) RemoveContainer(container string) (retErr error) {
	err := h.deleteContainer(container)
	if err != nil {
		log.Printf("Error deleting container %s. Error: %s.", container, err.Error())
		retErr = err
	}
	err = h.removeRelation(container)
	if err != nil {
		log.Printf("Error removing relationship from host node %s. Error: %s.", h.path(), err.Error())
		retErr = err
	}
	return
}

func ListHostsForApp(app string) (hosts []string, err error) {
	hosts, err = ListHosts()
	// TODO(jbhat): Filter out hosts that already are running the app or cannot run the app (whitelist/blacklist).
	// For now, return all hosts.
	return
}

func ListHosts() (hosts []string, err error) {
	hosts, _, err = Zk.Children(helper.GetBaseHostPath())
	if err != nil {
		log.Printf("Error getting list of hosts. Error: %s.", err.Error())
	}
	if hosts == nil {
		log.Println("No hosts found. Returning empty list.")
		hosts = []string{}
	}
	return
}

func (h ZkHost) Info() (*HostData, error) {
	data := &HostData{}
	err := getJson(h.path(), data)
	if err != nil {
		log.Printf("Error retrieving host data. Error: %s.", err.Error())
		return nil, err
	}
	return data, nil
}

// We will create private functions for use within this package

func (h ZkHost) deleteContainer(container string) (err error) {
	nodePath := h.containerPath(container)
	err = Zk.Delete(nodePath, -1)
	if err != nil {
		log.Printf("Error deleting node from zookeeper. Error: %s", err.Error())
	}
	return
}

func (h ZkHost) createOrUpdateContainer(container string, data *ContainerData) error {
	err := h.Touch()
	if err != nil {
		return err
	}
	return setJson(h.containerPath(container), data)
}

func (h ZkHost) containerPath(container string) string {
	return helper.GetBaseHostPath(string(h), container)
}

func (h ZkHost) path() string {
	return helper.GetBaseHostPath(string(h))
}

func (h ZkHost) addRelation(container string, port uint16) (err error) {
	data := HostData{}
	err = getJson(h.path(), &data)
	if err != nil {
		log.Printf("Error getting json from host node. Error: %s.", err.Error())
		return
	}
	if data.PortMap == nil {
		data.PortMap = map[string]uint16{}
	}
	data.PortMap[container] = port
	err = setJson(h.path(), &data)
	if err != nil {
		log.Printf("Error setting json to host node. Error: %s.", err.Error())
	}
	return
}

func (h ZkHost) removeRelation(container string) (retErr error) {
	data := HostData{}
	err := getJson(h.path(), &data)
	if err != nil {
		log.Printf("Error getting json from host node %s. Error: %s.", h.path(), err.Error())
		retErr = err
	}
	if data.PortMap == nil {
		data.PortMap = map[string]uint16{}
	}
	_, ok := data.PortMap[container]
	if ok {
		delete(data.PortMap, container)
		err = setJson(h.path(), &data)
		if err != nil {
			log.Printf("Error setting json for host node %s. Error: %s.", h.path(), err.Error())
			retErr = err
		}
	} else {
		retErr = errors.New(fmt.Sprintf("No port mapping exists on host %s for container %s\n", h.path(),
			container))
		log.Println(retErr.Error())
	}
	return
}

// used to sort host+weight
type HostAndWeight struct {
	Host   string
	Zone   string
	Free   uint
	Weight float64
}

type HostAndWeightList []HostAndWeight

func (h HostAndWeightList) Len() int {
	return len(h)
}

func (h HostAndWeightList) Less(i, j int) bool {
	return h[i].Weight < h[j].Weight
}

func (h HostAndWeightList) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

// Choses hosts and sorts them based on how "free" they are. returns a map of zone -> host slice.
func ChooseHosts(app, sha, env string, instances, cpu, memory uint, zones []string,
	excludeHosts map[string]bool) (map[string][]string, error) {
	hosts, err := ListHostsForApp(app)
	if err != nil {
		log.Println("Error listing hosts for app "+app+":", err)
		return nil, err
	}
	if len(hosts) == 0 {
		return nil, errors.New("No hosts available for app " + app)
	}
	list := HostAndWeightList{}
	for _, host := range hosts {
		if excludeHosts != nil && excludeHosts[host] {
			continue
		}
		// check if this host already has this app-sha
		hostInfo, err := Host(host).Info()
		if err != nil {
			continue // bad host, skip.
		}
		health, err := supervisor.HealthCheck(host)
		if err != nil || health.Status != StatusOk {
			continue // health check fail
		}
		if health.Containers.Free == 0 || health.Memory.Free < memory || health.CPUShares.Free < cpu {
			continue
		}
		// figure out how many we can stack on
		free := health.Containers.Free
		if health.Memory.Free/memory < free {
			free = health.Memory.Free / memory
		}
		if health.CPUShares.Free/cpu < free {
			free = health.CPUShares.Free / cpu
		}
		// we're chillin. add the weight to the host map
		// +2 weight for every one of this app/sha/env we see
		weight := float64(2*hostInfo.CountAppShaEnv(app, sha, env)) +
			(float64(health.Memory.Used+memory) / float64(health.Memory.Total)) +
			(float64(health.CPUShares.Used+cpu) / float64(health.CPUShares.Total))
		list = append(list, HostAndWeight{Host: host, Zone: health.Zone, Free: free, Weight: weight})
	}
	sort.Sort(list) // sort in weight order, lowest to highest
	chosenHosts := map[string][]string{}
	freeZones := map[string]uint{}
	for _, host := range list {
		if hosts, ok := chosenHosts[host.Zone]; !ok || hosts == nil {
			chosenHosts[host.Zone] = []string{host.Host}
		} else {
			chosenHosts[host.Zone] = append(hosts, host.Host)
		}
		freeZones[host.Zone] = freeZones[host.Zone] + host.Free
	}
	// ensure all zones are represented and have enough free
	for _, zone := range zones {
		if hosts, ok := chosenHosts[zone]; !ok || hosts == nil {
			msg := fmt.Sprintf("No host for app %s available in zone %s", app, zone)
			log.Println(msg)
			return nil, errors.New(msg)
		}
		if freeZones[zone] < instances {
			msg := fmt.Sprintf("Not enough instances for app %s available in zone %s (%d reqd, %d free)", app, zone,
				instances, freeZones[zone])
			log.Println(msg)
			return nil, errors.New(msg)
		}
	}
	return chosenHosts, nil
}
