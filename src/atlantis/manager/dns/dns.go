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
	CreateAliases(string, string, []Alias) (error, chan error)
	CreateARecords(string, string, []ARecord) (error, chan error)
	GetRecordsForIP(string, string) ([]string, error)
	DeleteRecords(string, string, ...string) (error, chan error)
	CreateHealthCheck(string, uint16) (string, error)
	DeleteHealthCheck(string) error
	Suffix(string) (string, error)
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
	// aliases are created for only this manager's region.
	suffix, err := Provider.Suffix(Region)
	if err != nil {
		return err
	}
	// check if records were created already, if so add sha to list
	zkDNS, err := datamodel.GetDNS(app, env)
	if zkDNS != nil && err == nil {
		if zkDNS.Shas == nil {
			zkDNS.Shas = map[string]bool{}
		}
		zkDNS.Shas[sha] = true
		return zkDNS.Save()
	}
	// for each zone
	zkDNS = datamodel.DNS(app, env)
	if zkDNS.Shas == nil {
		zkDNS.Shas = map[string]bool{}
	}
	zkDNS.Shas[sha] = true
	if Provider == nil {
		return zkDNS.Save()
	}

	zkDNS.RecordIds = []string{}
	aliases := []Alias{}
	// set up private zone aliases (no publics, just change the host header you lazy bum!)
	for _, zone := range AvailableZones {
		newAlias := Alias{
			Alias:    helper.GetZoneAppAlias(true, app, env, zone, suffix),
			Original: helper.GetZoneRouterCName(true, internal, zone, suffix),
		}
		aliases = append(aliases, newAlias)
		zkDNS.RecordIds = append(zkDNS.RecordIds, newAlias.Id())
	}
	// region-wide entry (for referencing outside of atlantis, use private for ec2, public for others)
	privateAlias := Alias{
		Alias:    helper.GetRegionAppAlias(true, app, env, suffix),
		Original: helper.GetRegionRouterCName(true, internal, suffix),
	}
	aliases = append(aliases, privateAlias)
	zkDNS.RecordIds = append(zkDNS.RecordIds, privateAlias.Id())
	publicAlias := Alias{
		Alias:    helper.GetRegionAppAlias(false, app, env, suffix),
		Original: helper.GetRegionRouterCName(false, internal, suffix),
	}
	aliases = append(aliases, publicAlias)
	zkDNS.RecordIds = append(zkDNS.RecordIds, publicAlias.Id())

	err, errChan := Provider.CreateAliases(Region, "CREATE_APP "+app+" in "+env, aliases)
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
	if zkDNS.Shas == nil {
		zkDNS.Shas = map[string]bool{}
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
	err, errChan := Provider.DeleteRecords(Region, "DELETE_APP "+app+" in "+env, zkDNS.RecordIds...)
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

func DeleteRecordsForIP(region, ip string) error {
	if Provider == nil {
		return nil
	}
	ids, err := Provider.GetRecordsForIP(region, ip)
	if err != nil {
		return err
	}
	err, errChan := Provider.DeleteRecords(region, "DELETE_ALL_IP "+ip, ids...)
	if err != nil {
		return err
	}
	return <-errChan
}
