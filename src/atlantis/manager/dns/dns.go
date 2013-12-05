package dns

import (
	. "atlantis/manager/constant"
	"atlantis/manager/datamodel"
	"atlantis/manager/helper"
	"crypto/sha256"
	"fmt"
)

var Provider DNSProvider

type DNSProvider interface {
	CreateAliases(string, []Alias) (error, chan error)
	CreateARecords(string, []ARecord) (error, chan error)
	GetRecordsForIP(string) ([]string, error)
	DeleteRecords(string, ...string) (error, chan error)
	CreateHealthCheck(string, uint16) (string, error)
	DeleteHealthCheck(string) error
	Suffix() string
}

type Alias struct {
	Alias    string
	Original string
	Failover string
}

func (a *Alias) Id() string {
	checksumArr := sha256.Sum256([]byte(fmt.Sprintf("%s %s", a.Original, a.Alias)))
	return fmt.Sprintf("%x", checksumArr[:sha256.Size])
}

type ARecord struct {
	Name          string
	IP            string
	HealthCheckId string
	Failover      string
	Weight        uint8
}

func (a *ARecord) Id() string {
	checksumArr := sha256.Sum256([]byte(fmt.Sprintf("%s %s", a.IP, a.Name)))
	return fmt.Sprintf("%x", checksumArr[:sha256.Size])
}

func CreateAppAliases(internal bool, app, sha, env string) error {
	// check if records were created already, if so add sha to list
	zkDNS, err := datamodel.GetDNS(app, env)
	if zkDNS != nil && err == nil {
		zkDNS.Shas[sha] = true
		return zkDNS.Save()
	}
	// for each zone
	zkDNS = datamodel.DNS(app, env)
	zkDNS.Shas[sha] = true
	if Provider == nil {
		return zkDNS.Save()
	}
	zkDNS.RecordIds = make([]string, 2*(len(AvailableZones)+1))
	aliases := make([]Alias, 2*(len(AvailableZones)+1))
	for i, zone := range AvailableZones {
		idx := i * 2
		aliases[idx] = Alias{
			Alias:    helper.GetZoneAppAlias(true, app, env, zone, Provider.Suffix()),
			Original: helper.GetZoneRouterCName(true, internal, zone, Provider.Suffix()),
		}
		zkDNS.RecordIds[idx] = aliases[idx].Id()
		idx++
		aliases[idx] = Alias{
			Alias:    helper.GetZoneAppAlias(false, app, env, zone, Provider.Suffix()),
			Original: helper.GetZoneRouterCName(false, internal, zone, Provider.Suffix()),
		}
		zkDNS.RecordIds[idx] = aliases[idx].Id()
	}
	// region-wide entry (for referencing outside of atlantis)
	aliases[len(aliases)-2] = Alias{
		Alias:    helper.GetRegionAppAlias(true, app, env, Provider.Suffix()),
		Original: helper.GetRegionRouterCName(true, internal, Provider.Suffix()),
	}
	zkDNS.RecordIds[len(aliases)-2] = aliases[len(aliases)-2].Id()
	aliases[len(aliases)-1] = Alias{
		Alias:    helper.GetRegionAppAlias(false, app, env, Provider.Suffix()),
		Original: helper.GetRegionRouterCName(false, internal, Provider.Suffix()),
	}
	zkDNS.RecordIds[len(aliases)-1] = aliases[len(aliases)-1].Id()
	err, errChan := Provider.CreateAliases("CREATE_APP "+app+" in "+env, aliases)
	if err != nil {
		return err
	}
	err = <-errChan // wait for change to propagate
	if err != nil {
		return err
	}
	// save records made in router zone path
	return zkDNS.Save()
}

func DeleteAppAliases(app, sha, env string) error {
	// find ids for app+env
	zkDNS, err := datamodel.GetDNS(app, env)
	if err != nil {
		return err
	}
	// remove sha from sha references
	delete(zkDNS.Shas, sha)
	err = zkDNS.Save()
	if err != nil {
		return err
	}
	// if this was *not* the last sha, don't delete anything
	if len(zkDNS.Shas) > 0 {
		return nil
	}
	if Provider == nil {
		return zkDNS.Delete()
	}
	// delete all the record ids
	err, errChan := Provider.DeleteRecords("DELETE_APP "+app+" in "+env, zkDNS.RecordIds...)
	if err != nil {
		return err
	}
	err = <-errChan // wait for change to propagate
	if err != nil {
		return err
	}
	// remove dns datamodel
	return zkDNS.Delete()
}

func DeleteRecordsForIP(ip string) error {
	if Provider == nil {
		return nil
	}
	ids, err := Provider.GetRecordsForIP(ip)
	if err != nil {
		return err
	}
	err, errChan := Provider.DeleteRecords("DELETE_ALL_IP "+ip, ids...)
	if err != nil {
		return err
	}
	return <-errChan
}
