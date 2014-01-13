package datamodel

import (
	. "atlantis/manager/constant"
	"atlantis/manager/helper"
	"atlantis/manager/rpc/types"
	routercfg "atlantis/router/config"
	routerzk "atlantis/router/zk"
	"errors"
	"fmt"
	"log"
	"strconv"
)

var (
	MinRouterPort = DefaultMinRouterPort
	MaxRouterPort = DefaultMaxRouterPort
)

// Router Port Reservation

type ZkRouterPorts types.RouterPorts

func GetRouterPorts(internal bool) (zr *ZkRouterPorts) {
	zr = &ZkRouterPorts{}
	err := getJson(helper.GetBaseRouterPortsPath(internal), zr)
	if err != nil || zr == nil {
		zr = &ZkRouterPorts{
			Internal:  internal,
			PortMap:   map[string]types.AppEnv{},
			AppEnvMap: map[string]string{},
		}
		zr.save()
	} else if zr.PortMap == nil || zr.AppEnvMap == nil {
		zr.Internal = internal
		zr.PortMap = map[string]types.AppEnv{}
		zr.AppEnvMap = map[string]string{}
		zr.save()
	}
	return zr
}

func HasRouterPortForAppEnv(internal bool, app, env string) bool {
	lock := NewRouterPortsLock(internal)
	lock.Lock()
	defer lock.Unlock()
	zrp := GetRouterPorts(internal)
	return zrp.hasPortForAppEnv(app, env)
}

func ReclaimRouterPortsForApp(internal bool, app string) error {
	helper.SetRouterRoot(internal)
	lock := NewRouterPortsLock(internal)
	lock.Lock()
	defer lock.Unlock()
	zrp := GetRouterPorts(internal)
	envs, err := zrp.reclaimApp(app)
	if err != nil {
		return err
	}
	for _, env := range envs {
		name := helper.GetAppEnvTrieName(app, env)
		if err := routerzk.DelPort(Zk.Conn, name); err != nil {
			log.Printf("Error reclaiming port %s", name)
			// don't fail here
			// TODO email appsplat
		}
	}
	return nil
}

func ReclaimRouterPortsForEnv(internal bool, env string) error {
	helper.SetRouterRoot(internal)
	lock := NewRouterPortsLock(internal)
	lock.Lock()
	defer lock.Unlock()
	zrp := GetRouterPorts(internal)
	apps, err := zrp.reclaimEnv(env)
	if err != nil {
		return err
	}
	for _, app := range apps {
		name := helper.GetAppEnvTrieName(app, env)
		if err := routerzk.DelPort(Zk.Conn, name); err != nil {
			log.Printf("Error reclaiming port %s", name)
			// don't fail here
			// TODO email appsplat
		}
	}
	return nil
}

func ReserveRouterPortAndUpdateTrie(app, sha, env string) (string, bool, error) {
	helper.SetRouterRoot(true)
	var (
		err     error
		created = false
		port    = ""
	)
	// reserve port for app env
	if !HasRouterPortForAppEnv(true, app, env) {
		created = true
	}
	if port, err = reserveRouterPort(true, app, env); err != nil {
		return port, created, err
	}
	// create trie (if it doesn't exist)
	trieName := helper.GetAppEnvTrieName(app, env)
	if exists, err := routerzk.TrieExists(Zk.Conn, trieName); !exists || err != nil {
		err = routerzk.SetTrie(Zk.Conn, routercfg.Trie{
			Name:     trieName,
			Rules:    []string{},
			Internal: true,
		})
		if err != nil {
			return port, created, err
		}
	}
	// if sha != "" attach pool as static rule (if trie is empty)
	if sha != "" {
		// if static rule does not exist, create it
		ruleName := helper.GetAppShaEnvStaticRuleName(app, sha, env)
		poolName := helper.CreatePoolName(app, sha, env)
		if exists, err := routerzk.RuleExists(Zk.Conn, ruleName); !exists || err != nil {
			err = routerzk.SetRule(Zk.Conn, routercfg.Rule{
				Name:     ruleName,
				Type:     "static",
				Value:    "true",
				Pool:     poolName,
				Internal: true,
			})
			if err != nil {
				return port, created, err
			}
		}
		trie, err := routerzk.GetTrie(Zk.Conn, trieName)
		if err != nil {
			return port, created, err
		}
		if len(trie.Rules) == 0 {
			trie.Rules = []string{ruleName}
		} else {
			trie.Rules = append(trie.Rules, ruleName)
		}
		if err = routerzk.SetTrie(Zk.Conn, trie); err != nil {
			return port, created, err
		}
	}
	// now port is reserved and trie is created so we can actually create port for router
	portUInt, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return port, created, err
	}
	err = routerzk.SetPort(Zk.Conn, routercfg.Port{
		Name: trieName,
		Port: uint16(portUInt),
		Trie: trieName,
	})
	if err != nil {
		return port, created, err
	}
	// return true if port was created
	return port, created, err
}

func reserveRouterPort(internal bool, app, env string) (string, error) {
	lock := NewRouterPortsLock(internal)
	lock.Lock()
	defer lock.Unlock()
	zrp := GetRouterPorts(internal)
	return zrp.getPortForAppEnv(app, env)
}

func (r *ZkRouterPorts) hasPortForAppEnv(app, env string) bool {
	appEnv := types.AppEnv{App: app, Env: env}
	port := r.AppEnvMap[appEnv.String()]
	return port != ""
}

func (r *ZkRouterPorts) getPortForAppEnv(app, env string) (string, error) {
	appEnv := types.AppEnv{App: app, Env: env}
	port := r.AppEnvMap[appEnv.String()]
	if port != "" {
		return port, nil
	}
	for i := MinRouterPort; MinRouterPort <= i && i <= MaxRouterPort; i++ {
		portStr := fmt.Sprintf("%d", i)
		if _, ok := r.PortMap[portStr]; ok {
			continue
		}
		r.PortMap[portStr] = appEnv
		r.AppEnvMap[appEnv.String()] = portStr
		return portStr, r.save()
	}
	// TODO email appsplat?
	return "", errors.New("No available ports")
}

func (r *ZkRouterPorts) reclaimApp(app string) ([]string, error) {
	reclaimedEnvs := []string{}
	for port, appEnv := range r.PortMap {
		if appEnv.App == app {
			delete(r.PortMap, port)
			delete(r.AppEnvMap, appEnv.String())
			reclaimedEnvs = append(reclaimedEnvs, appEnv.Env)
		}
	}
	return reclaimedEnvs, r.save()
}

func (r *ZkRouterPorts) reclaimEnv(env string) ([]string, error) {
	reclaimedApps := []string{}
	for port, appEnv := range r.PortMap {
		if appEnv.Env == env {
			delete(r.PortMap, port)
			delete(r.AppEnvMap, appEnv.String())
			reclaimedApps = append(reclaimedApps, appEnv.App)
		}
	}
	return reclaimedApps, r.save()
}

func (r *ZkRouterPorts) save() error {
	return setJson(r.path(), r)
}

func (r *ZkRouterPorts) path() string {
	return helper.GetBaseRouterPortsPath(r.Internal)
}

// Router Registration

type ZkRouter types.Router

func Router(internal bool, zone, value string) *ZkRouter {
	return &ZkRouter{Internal: internal, Zone: zone, Host: value}
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

func GetRouter(internal bool, zone, value string) (zr *ZkRouter, err error) {
	zr = &ZkRouter{}
	err = getJson(helper.GetBaseRouterPath(internal, zone, value), zr)
	return
}

func (r *ZkRouter) path() string {
	return helper.GetBaseRouterPath(r.Internal, r.Zone, r.Host)
}

// Routing Datamodel

func defaultPool(name string, internal bool) routercfg.Pool {
	return routercfg.Pool{
		Name:     name,
		Config:   routercfg.PoolConfig{HealthzEvery: "5s", HealthzTimeout: "5s", RequestTimeout: "120s", Status: "OK"},
		Hosts:    map[string]routercfg.Host{},
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
		zkApp, err := GetApp(inst.App)
		if err != nil {
			return err
		}
		currInsts := pools[zkApp.Internal][name]
		if currInsts == nil {
			currInsts = []*ZkInstance{}
		}
		pools[zkApp.Internal][name] = append(currInsts, inst)
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
			hosts := map[string]routercfg.Host{}
			for _, inst := range insts {
				address := fmt.Sprintf("%s:%d", inst.Host, inst.Port)
				hosts[address] = routercfg.Host{Address: address}
			}
			if err := routerzk.AddHosts(Zk.Conn, name, hosts); err != nil {
				return err
			}
		}
	}
	return nil
}

type poolDefinition struct {
	app   string
	sha   string
	env   string
	insts []*ZkInstance
}

func DeleteFromPool(containers []string) error {
	pools := map[bool]map[string]*poolDefinition{}
	pools[true] = map[string]*poolDefinition{}
	pools[false] = map[string]*poolDefinition{}
	for _, cont := range containers {
		inst, err := GetInstance(cont)
		if err != nil {
			// instance doesn't exist
			continue
		}
		name := helper.CreatePoolName(inst.App, inst.Sha, inst.Env)
		zkApp, err := GetApp(inst.App)
		if err != nil {
			return err
		}
		poolDef := pools[zkApp.Internal][name]
		if poolDef == nil {
			poolDef = &poolDefinition{
				app:   inst.App,
				sha:   inst.Sha,
				env:   inst.Env,
				insts: []*ZkInstance{},
			}
			pools[zkApp.Internal][name] = poolDef
		}
		pools[zkApp.Internal][name].insts = append(poolDef.insts, inst)
	}
	for internal, allPools := range pools {
		helper.SetRouterRoot(internal)
		for name, poolDef := range allPools {
			// remove hosts
			hosts := []string{}
			for _, inst := range poolDef.insts {
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
				err = CleanupCreatedPoolRefs(internal, poolDef.app, poolDef.sha, poolDef.env)
				if err != nil {
					log.Println("Error trying to clean up pool:", err)
				}
			}
		}
	}
	return nil
}

func CleanupCreatedPoolRefs(internal bool, app, sha, env string) error {
	helper.SetRouterRoot(internal)
	if !internal {
		return nil
	}
	// remove static rule, cleanup rule from trie if needed
	ruleName := helper.GetAppShaEnvStaticRuleName(app, sha, env)
	trieName := helper.GetAppEnvTrieName(app, env)
	// remove static rule from trie
	trie, err := routerzk.GetTrie(Zk.Conn, trieName)
	if err != nil {
		return err
	}
	newRules := []string{}
	for _, rule := range trie.Rules {
		if rule != ruleName {
			newRules = append(newRules, ruleName)
		}
	}
	if len(trie.Rules) != len(newRules) {
		trie.Rules = newRules
		err = routerzk.SetTrie(Zk.Conn, trie)
		if err != nil {
			return err
		}
	}
	// delete static rule
	err = routerzk.DelRule(Zk.Conn, ruleName)
	if err != nil {
		return err
	}
	return nil
}
