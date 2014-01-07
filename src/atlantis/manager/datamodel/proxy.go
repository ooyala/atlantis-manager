package datamodel

import (
	. "atlantis/common"
	. "atlantis/manager/constant"
	"atlantis/manager/dns"
	"atlantis/manager/helper"
	"atlantis/manager/supervisor"
	"atlantis/proxy/types"
	"errors"
	"fmt"
	"strings"
)

var (
	MinProxyPort uint16
	MaxProxyPort uint16
)

type ZkProxyAppEnv struct {
	App string
	Env string
}

type ZkProxy struct {
	Sha     string                   // current version
	PortMap map[string]ZkProxyAppEnv // port -> app+env
	AppMap  map[string]string        // app.env -> port
}

func GetProxy() *ZkProxy {
	zp := &ZkProxy{}
	err := getJson(helper.GetBaseProxyPath(), zp)
	if err != nil || zp == nil {
		zp = &ZkProxy{
			PortMap: map[string]ZkProxyAppEnv{},
			AppMap:  map[string]string{},
		}
		zp.Save()
	} else if zp.PortMap == nil || zp.AppMap == nil {
		zp.PortMap = map[string]ZkProxyAppEnv{}
		zp.AppMap = map[string]string{}
		zp.Save()
	}
	return zp
}

func (zp *ZkProxy) PortForAppEnv(app, env string) (string, error) {
	portStr, ok := zp.AppMap[app+"."+env]
	if !ok {
		return "", errors.New(app + " in " + env + " is not in the proxy.")
	}
	return portStr, nil
}

func ConfigureProxy() (err error) {
	zp := GetProxy()
	if zp.Sha == "" {
		return nil // no proxy deployed yet, no need to configure anything
	}
	// find all supervisors
	supers, err := ListSupervisors()
	if err != nil {
		return err
	}
	// map of zone -> config
	cfgs := map[string]map[string]*types.ProxyConfig{}
	for _, zone := range AvailableZones {
		cfgs[zone], err = GetProxyConfig(zone)
		if err != nil {
			return err
		}
	}
	// configure proxy on each supervisor one-by-one
	for _, super := range supers {
		zone, err := supervisor.GetZone(super)
		if err != nil {
			return err
		}
		if cfgs[zone] == nil {
			return errors.New("Invalid Zone ("+super+"): "+zone)
		}
		sReply, err := supervisor.ConfigureProxy(super, cfgs[zone])
		if err != nil {
			return err
		} else if sReply.Status != StatusOk {
			return errors.New("Configure Proxy Status ("+super+"): "+sReply.Status)
		}
	}
	return nil
}

func GetProxyConfig(zone string) (map[string]*types.ProxyConfig, error) {
	zp := GetProxy()
	cfg := map[string]*types.ProxyConfig{}
	for port, zkAppEnv := range zp.PortMap {
		zkApp, err := GetApp(zkAppEnv.App)
		if err != nil {
			return nil, err
		}
		lAddr := "0.0.0.0:"+port
		if zkApp.NonAtlantis {
			// use addrs from app
			rAddr := zkApp.Addrs[zkAppEnv.Env]
			if rAddr == "" {
				return nil, errors.New("app " + zkAppEnv.App + " does not have env " + zkAppEnv.Env)
			}
			cfg[lAddr] = &types.ProxyConfig{
				Type:       types.ProxyTypeTCP,
				LocalAddr:  lAddr,
				RemoteAddr: rAddr,
			}
			if strings.ToLower(zkApp.Type) == "http" {
				cfg[lAddr].Type = types.ProxyTypeHTTP
			}
		} else if dns.Provider != nil {
			suffix, err := dns.Provider.Suffix(Region)
			if err != nil {
				return nil, err
			}
			// use internal app dns stuff
			cfg[lAddr] = &types.ProxyConfig{
				Type:       types.ProxyTypeHTTP,
				LocalAddr:  lAddr,
				RemoteAddr: helper.GetZoneAppCName(zkAppEnv.App, zkAppEnv.Env, zone, suffix),
			}
		}
	}
	return cfg, nil
}

func (zp *ZkProxy) AddAll(app string, envs []string) error {
	i := MinProxyPort
	for envCount := 0; envCount < len(envs); envCount++ {
		if _, ok := zp.AppMap[app+"."+envs[envCount]]; ok {
			continue
		}
		// find a port
		for ; MinProxyPort <= i && i <= MaxProxyPort; i++ {
			iStr := fmt.Sprintf("%d", i)
			if _, taken := zp.PortMap[iStr]; taken {
				continue
			}
			zp.PortMap[iStr] = ZkProxyAppEnv{App: app, Env: envs[envCount]}
			zp.AppMap[app+"."+envs[envCount]] = iStr
			break
		}
		if i == MaxProxyPort+1 {
			// TODO[jigish] email appsplat. this is a problem lol.
			return errors.New("Not Enough Available Ports")
		}
	}
	ConfigureProxy()
	return zp.Save()
}

func (zp *ZkProxy) RemoveAll(app string, envs []string) error {
	for _, env := range envs {
		currentPort, ok := zp.AppMap[app+"."+env]
		if !ok {
			continue
		}
		delete(zp.AppMap, app+"."+env)
		delete(zp.PortMap, currentPort)
	}
	ConfigureProxy()
	return zp.Save()
}

func (zp *ZkProxy) AddAppEnv(app, env string) error {
	if _, ok := zp.AppMap[app+"."+env]; ok {
		return nil // we already have added this
	}
	// find a port
	var i uint16
	for i = MinProxyPort; MinProxyPort <= i && i <= MaxProxyPort; i++ {
		iStr := fmt.Sprintf("%d", i)
		if _, taken := zp.PortMap[iStr]; taken {
			continue
		}
		zp.PortMap[iStr] = ZkProxyAppEnv{App: app, Env: env}
		zp.AppMap[app+"."+env] = iStr
		break
	}
	if i == MaxProxyPort+1 {
		// TODO[jigish] email appsplat. this is a problem lol.
		return errors.New("No Available Ports")
	}
	ConfigureProxy()
	return zp.Save()
}

func (zp *ZkProxy) RemoveAppEnv(app, env string) error {
	currentPort, ok := zp.AppMap[app+"."+env]
	if !ok {
		return nil
	}
	delete(zp.AppMap, app+"."+env)
	delete(zp.PortMap, currentPort)
	ConfigureProxy()
	return zp.Save()
}

func (zp *ZkProxy) path() string {
	return helper.GetBaseProxyPath()
}

func (zp *ZkProxy) Save() error {
	if err := setJson(zp.path(), zp); err != nil {
		return err
	}
	return nil
}
