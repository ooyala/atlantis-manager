package app

import (
	. "atlantis/manager/constant"
	"atlantis/manager/datamodel"
	"atlantis/manager/dns"
	"atlantis/manager/helper"
)

func CreateAppCNames(internal bool, app, sha, env string) error {
	// cnames are created for only this manager's region.
	suffix, err := dns.Provider.Suffix(Region)
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
	if dns.Provider == nil {
		return zkDNS.Save()
	}

	zkDNS.RecordIDs = []string{}
	cnames := []dns.Record{}
	// set up zone cname
	for _, zone := range AvailableZones {
		newCName := dns.NewRecord(helper.GetZoneAppCName(app, env, zone, suffix),
			helper.GetZoneRouterCName(internal, zone, suffix), 1)
		cnames = append(cnames, newCName)
		zkDNS.RecordIDs = append(zkDNS.RecordIDs, newCName.ID())
	}
	// region-wide entry (for referencing outside of atlantis)
	regionCName := dns.NewRecord(helper.GetRegionAppCName(app, env, suffix),
		helper.GetRegionRouterCName(internal, suffix), 1)
	cnames = append(cnames, regionCName)
	zkDNS.RecordIDs = append(zkDNS.RecordIDs, regionCName.ID())

	err = dns.Provider.CreateRecords(Region, "CREATE_APP "+app+" in "+env, cnames)
	if err != nil {
		return err
	}
	// save records made in router zone path
	return zkDNS.Save()
}

func DeleteAppCNames(app, sha, env string) error {
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
	if dns.Provider == nil {
		return zkDNS.Delete()
	}
	// delete all the record ids
	err, errChan := dns.Provider.DeleteRecords(Region, "DELETE_APP "+app+" in "+env, zkDNS.RecordIDs...)
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
