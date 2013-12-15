package datamodel

import (
	"atlantis/manager/helper"
	"errors"
	"fmt"
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
	// TODO[jigish] add app+env to all proxies
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
	// TODO[jigish] remove app+env from all proxies
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
	// TODO[jigish] add app+env to all proxies
	return zp.Save()
}

func (zp *ZkProxy) RemoveAppEnv(app, env string) error {
	currentPort, ok := zp.AppMap[app+"."+env]
	if !ok {
		return nil
	}
	delete(zp.AppMap, app+"."+env)
	delete(zp.PortMap, currentPort)
	// TODO[jigish] remove app+env from all proxies
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
