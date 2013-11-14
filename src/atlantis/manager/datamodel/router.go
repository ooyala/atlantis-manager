package datamodel

import (
	"atlantis/manager/helper"
	"atlantis/router/config"
	routerzk "atlantis/router/zk"
	"fmt"
	"log"
)

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
