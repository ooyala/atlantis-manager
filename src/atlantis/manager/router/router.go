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
// - An A record for value router4.us-east-1a.atlantis.com pointing at 10.0.0.4
// - A primary failover A record for value router.us-east-1a.atlantis.com pointing at 10.0.0.4
// - A secondary failover A record for value router.us-east-1b.atlantis.com pointing to 10.0.0.4
// - A secondary failover A record for value router.us-east-1c.atlantis.com pointing to 10.0.0.4
// - A round robin A record for value router.us-east-1.atlantis.com pointing at 10.0.0.4
//
// Deleting the router, simply deletes all the records created when adding it.

func Register(internal bool, zone, ip string) (*datamodel.ZkRouter, error) {
	// create ZkRouter
	zkRouter := datamodel.Router(internal, zone, ip)
	if dns.Provider == nil {
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
	// create basic health check for router
	zkRouter.HealthCheckId, err = dns.Provider.CreateHealthCheck(zkRouter.IP, uint16(80))
	if err != nil {
		return zkRouter, err
	}
	zkRouter.RecordIds = make([]string, 3)
	cnames := make([]dns.CName, 3)
	// PRIMARY router.<region>.<suffix>
	cnames[0] = dns.CName{
		CName:         helper.GetRegionRouterCName(internal, dns.Provider.Suffix()),
		IP:            zkRouter.IP,
		HealthCheckId: zkRouter.HealthCheckId,
	}
	zkRouter.RecordIds[0] = cnames[0].Id()
	// PRIMARY routerX.<region+zone>.<suffix>
	cnames[1] = dns.CName{
		CName:         zkRouter.CName,
		IP:            zkRouter.IP,
		HealthCheckId: zkRouter.HealthCheckId,
	}
	zkRouter.RecordIds[1] = cnames[1].Id()
	// PRIMARY router.<region+zone>.<suffix>
	cnames[2] = dns.CName{
		Failover:      "PRIMARY",
		CName:         helper.GetZoneRouterCName(internal, zkRouter.Zone, dns.Provider.Suffix()),
		IP:            zkRouter.IP,
		HealthCheckId: zkRouter.HealthCheckId,
	}
	zkRouter.RecordIds[2] = cnames[2].Id()
	// SECONDARY router.<region+zone>.<suffix>
	for _, azone := range AvailableZones {
		if azone == zone {
			continue
		}
		cname := dns.CName{
			Failover:      "SECONDARY",
			CName:         helper.GetZoneRouterCName(internal, azone, dns.Provider.Suffix()),
			IP:            zkRouter.IP,
			HealthCheckId: zkRouter.HealthCheckId,
		}
		zkRouter.RecordIds = append(zkRouter.RecordIds, cname.Id())
		cnames = append(cnames, cname)
	}
	err, errChan := dns.Provider.CreateCNames("CREATE_ROUTER "+ip+" in "+zone, cnames)
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
	// delete basic health check for router
	err = dns.Provider.DeleteHealthCheck(zkRouter.HealthCheckId)
	if err != nil {
		return err
	}
	return zkRouter.Delete()
}
