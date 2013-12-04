package router

import (
	. "atlantis/manager/constant"
	"atlantis/manager/datamodel"
	"atlantis/manager/dns"
	"atlantis/manager/helper"
)

// Registering the 4th router in us-east-1a, with IP (say) 10.0.0.4 needs to create the following
// entries in DNS (route53 in our case), assuming that atlantis.com is the DNS zone delegated to
// the deployment. Also assume that the deployment spans zones a, b and c in region us-east-1.
//
// - A weight 1 A record for value router4.us-east-1a.atlantis.com pointing at 10.0.0.4
// - A weight 1 A record for value router.us-east-1a.atlantis.com pointing at 10.0.0.4
// - A weight 0 A record for value router.us-east-1b.atlantis.com pointing to 10.0.0.4
// - A weight 0 A record for value router.us-east-1c.atlantis.com pointing to 10.0.0.4
// - A weight 1 A record for value router.us-east-1.atlantis.com pointing at 10.0.0.4
//
// The weight 0 A records will be activated if health checks fail on the router in that zone (not implemented)
//
// Deleting the router, simply deletes all the records created when adding it.

func Register(internal bool, zone, ip string) (*datamodel.ZkRouter, error) {
	// create ZkRouter
	zkRouter := datamodel.Router(internal, zone, ip)
	if !internal || dns.Provider == nil {
		// only internal routers need DNS stuff
		// if we have no dns provider then just save here
		return zkRouter, zkRouter.Save()
	}
	// first delete all entries we may already have for this IP in DNS
	err := dns.DeleteRecordsForIP(ip)
	if err != nil {
		return nil, err
	}
	// choose cname
	routers, err := datamodel.ListRoutersInZone(internal, zone)
	if err != nil {
		return zkRouter, err
	}
	routerMap := map[string]bool{}
	for _, router := range routers {
		tmpRouter, err := datamodel.GetRouter(internal, zone, router)
		if err != nil {
			return zkRouter, err
		}
		routerMap[tmpRouter.CName] = true
	}
	routerNum := 1
	zkRouter.CName = helper.GetRouterCName(internal, routerNum, zone, dns.Provider.Suffix())
	for ; routerMap[zkRouter.CName]; routerNum++ {
		zkRouter.CName = helper.GetRouterCName(internal, routerNum, zone, dns.Provider.Suffix())
	}

	zkRouter.RecordIds = make([]string, 3)
	cnames := make([]dns.ARecord, 3)

	// WEIGHT=1 router.<region>.<suffix>
	cnames[0] = dns.ARecord{
		Name:   helper.GetRegionRouterCName(internal, dns.Provider.Suffix()),
		IP:     zkRouter.IP,
		Weight: 1,
	}
	zkRouter.RecordIds[0] = cnames[0].Id()

	// WEIGHT=1 routerX.<region+zone>.<suffix>
	cnames[1] = dns.ARecord{
		Name:   zkRouter.CName,
		IP:     zkRouter.IP,
		Weight: 1,
	}
	zkRouter.RecordIds[1] = cnames[1].Id()

	// WEIGHT=1 router.<region+zone>.<suffix>
	cnames[2] = dns.ARecord{
		Name:   helper.GetZoneRouterCName(internal, zkRouter.Zone, dns.Provider.Suffix()),
		IP:     zkRouter.IP,
		Weight: 1,
	}
	zkRouter.RecordIds[2] = cnames[2].Id()

	// WEIGHT=0 router.<region+zone>.<suffix> -> will be activated when needed
	for _, azone := range AvailableZones {
		if azone == zone {
			continue
		}
		cname := dns.ARecord{
			Name:   helper.GetZoneRouterCName(internal, azone, dns.Provider.Suffix()),
			IP:     zkRouter.IP,
			Weight: 0,
		}
		zkRouter.RecordIds = append(zkRouter.RecordIds, cname.Id())
		cnames = append(cnames, cname)
	}

	err, errChan := dns.Provider.CreateARecords("CREATE_ROUTER "+ip+" in "+zone, cnames)
	if err != nil {
		return zkRouter, err
	}
	err = <-errChan // wait for change to propagate
	if err != nil {
		return zkRouter, err
	}
	return zkRouter, zkRouter.Save()
}

func Unregister(internal bool, zone, ip string) error {
	zkRouter, err := datamodel.GetRouter(internal, zone, ip)
	if err != nil {
		return err
	}
	if dns.Provider == nil {
		// if we have no dns provider then just save here
		return zkRouter.Delete()
	}
	err, errChan := dns.Provider.DeleteRecords("DELETE_ROUTER "+ip+" in "+zone, zkRouter.RecordIds...)
	if err != nil {
		return err
	}
	err = <-errChan // wait for it to propagate
	if err != nil {
		return err
	}
	return zkRouter.Delete()
}
