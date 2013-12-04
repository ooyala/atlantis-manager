package datamodel

import (
	"atlantis/manager/helper"
	"atlantis/manager/rpc/types"
	"atlantis/router/config"
	routerzk "atlantis/router/zk"
	"fmt"
	"log"
)

// Router Registration

type ZkRouter types.Router

func Router(internal bool, zone, privateIP, publicIP string) *ZkRouter {
	return &ZkRouter{Internal: internal, Zone: zone, PrivateIP: privateIP, PublicIP: publicIP}
}

func (r *ZkRouter) Save() error {
	return setJson(r.path(), r)
}

func (r *ZkRouter) Delete() error {
	return Zk.RecursiveDelete(r.path())
}

func ListRouterZones(internal bool) (zones []string, err error) {
	basePath := helper.GetBaseRouterPath(internal)
	zones, _, err = Zk.Children(basePath)
	if err != nil {
		log.Printf("Error getting list of zones. Error: %s.", err.Error())
	}
	if zones == nil {
		log.Println("No zones found. Returning empty list.")
		zones = []string{}
	}
	return
}

func ListRoutersInZone(internal bool, zone string) (routers []string, err error) {
	basePath := helper.GetBaseRouterPath(internal, zone)
	routers, _, err = Zk.Children(basePath)
	if err != nil {
		log.Printf("Error getting list of routers for zone %s. Error: %s.", zone, err.Error())
	}
	if routers == nil {
		log.Printf("No routers found in zone %s", zone)
		routers = []string{}
	}
	return
}

func ListRouters(internal bool) (routers map[string][]string, err error) {
	routers = map[string][]string{}
	zones, err := ListRouterZones(internal)
	if err != nil {
		return
	}
	for _, zone := range zones {
		routers[zone], err = ListRoutersInZone(internal, zone)
		if err != nil {
			return
		}
	}
	return
}

func GetRouter(internal bool, zone, ip string) (zr *ZkRouter, err error) {
	zr = &ZkRouter{}
	err = getJson(helper.GetBaseRouterPath(internal, zone, ip), zr)
	return
}

func (r *ZkRouter) path() string {
	return helper.GetBaseRouterPath(r.Internal, r.Zone, r.PublicIP)
}

// Routing Datamodel

func defaultPool(name string, internal bool) config.Pool {
	return config.Pool{
		Name:     name,
		Config:   config.PoolConfig{HealthzEvery: "5s", HealthzTimeout: "5s", RequestTimeout: "120s", Status: "OK"},
		Hosts:    map[string]config.Host{},
		Internal: internal,
	}
}

func AddToPool(containers []string) error {
	pools := map[bool]map[string][]*ZkInstance{}
	pools[true] = map[string][]*ZkInstance{}
	pools[false] = map[string][]*ZkInstance{}
	for _, cont := range containers {
		inst, err := GetInstance(cont)
		if err != nil {
			// instance doesn't exist
			continue
		}
		name := helper.CreatePoolName(inst.App, inst.Sha, inst.Env)
		currInsts := pools[inst.Internal][name]
		if currInsts == nil {
			currInsts = []*ZkInstance{}
		}
		pools[inst.Internal][name] = append(currInsts, inst)
	}
	for internal, allPools := range pools {
		helper.SetRouterRoot(internal)
		for name, insts := range allPools {
			// create pool if we need to
			if exists, err := routerzk.PoolExists(Zk.Conn, name); !exists || err != nil {
				if err = routerzk.SetPool(Zk.Conn, defaultPool(name, internal)); err != nil {
					return err
				}
			}
			// add hosts
			hosts := map[string]config.Host{}
			for _, inst := range insts {
				address := fmt.Sprintf("%s:%d", inst.Host, inst.Port)
				hosts[address] = config.Host{Address: address}
			}
			if err := routerzk.AddHosts(Zk.Conn, name, hosts); err != nil {
				return err
			}
		}
	}
	return nil
}

func DeleteFromPool(containers []string) error {
	pools := map[bool]map[string][]*ZkInstance{}
	pools[true] = map[string][]*ZkInstance{}
	pools[false] = map[string][]*ZkInstance{}
	for _, cont := range containers {
		inst, err := GetInstance(cont)
		if err != nil {
			// instance doesn't exist
			continue
		}
		name := helper.CreatePoolName(inst.App, inst.Sha, inst.Env)
		currInsts := pools[inst.Internal][name]
		if currInsts == nil {
			currInsts = []*ZkInstance{}
		}
		pools[inst.Internal][name] = append(currInsts, inst)
	}
	for internal, allPools := range pools {
		helper.SetRouterRoot(internal)
		for name, insts := range allPools {
			// remove hosts
			hosts := []string{}
			for _, inst := range insts {
				hosts = append(hosts, fmt.Sprintf("%s:%d", inst.Host, inst.Port))
			}
			routerzk.DelHosts(Zk.Conn, name, hosts)
			// delete pool if no hosts exist
			getHosts, err := routerzk.GetHosts(Zk.Conn, name)
			if err != nil || len(getHosts) == 0 {
				err = routerzk.DelPool(Zk.Conn, name)
				if err != nil {
					log.Println("Error trying to clean up pool:", err)
				}
			}
		}
	}
	return nil
}
