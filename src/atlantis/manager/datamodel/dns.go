package datamodel

import (
	"atlantis/manager/helper"
)

type ZkDNS struct {
	App       string
	Env       string
	Shas      map[string]bool
	RecordIDs []string
}

func DNS(app, env string) *ZkDNS {
	return &ZkDNS{App: app, Env: env}
}

func GetDNS(app, env string) (zd *ZkDNS, err error) {
	zd = &ZkDNS{}
	err = getJson(helper.GetBaseDNSPath(app, env), zd)
	return
}

func (d *ZkDNS) Save() error {
	return setJson(d.path(), d)
}

func (d *ZkDNS) Delete() error {
	return Zk.RecursiveDelete(d.path())
}

func (r *ZkDNS) path() string {
	return helper.GetBaseDNSPath(r.App, r.Env)
}
