package router

import (
	. "atlantis/manager/constant"
	"atlantis/manager/datamodel"
	"atlantis/manager/dns"
	"atlantis/manager/helper"
)

func Register(zone, ip string) (*datamodel.ZkRouter, error) {
	// create ZkRouter
	zkRouter := datamodel.Router(zone, ip)
	if dns.Provider == nil {
		// if we have no dns provider then just save here
		return zkRouter, zkRouter.Save()
	}
	// choose cname
	routers, err := datamodel.ListRoutersInZone(zone)
	if err != nil {
		return zkRouter, err
	}
	routerMap := map[string]bool{}
	for _, router := range routers {
		tmpRouter, err := datamodel.GetRouter(zone, router)
		if err != nil {
			return zkRouter, err
		}
		routerMap[tmpRouter.CName] = true
	}
	routerNum := 1
	zkRouter.CName = helper.GetRouterCName(routerNum, zone, dns.Provider.Suffix())
	for ; routerMap[zkRouter.CName]; routerNum++ {
		zkRouter.CName = helper.GetRouterCName(routerNum, zone, dns.Provider.Suffix())
	}
	// create basic health check for router
	zkRouter.HealthCheckId, err = dns.Provider.CreateHealthCheck(zkRouter.IP)
	if err != nil {
		return zkRouter, err
	}
	// create the basic entries for the new router
	// e.g. region is us-east-1 && zone is us-east-1a && there are zones a, b, c, and d && this is the 2nd in a
	//   create router.us-east-1.<suffix>   primary
	//   create router2.us-east-1a.<suffix> primary
	//   create router.us-east-1a.<suffix>  primary
	//   create router.us-east-1b.<suffix>  secondary
	//   create router.us-east-1c.<suffix>  secondary
	//   create router.us-east-1d.<suffix>  secondary
	zkRouter.RecordIds = make([]string, 3)
	cnames := make([]dns.CName, 3)
	// PRIMARY router.<region>.<suffix>
	cnames[0] = dns.CName{
		Primary: true,
		CName:   helper.GetRegionRouterCName(dns.Provider.Suffix()),
		IP:      zkRouter.IP,
	}
	zkRouter.RecordIds[0] = cnames[0].Id()
	// PRIMARY routerX.<region+zone>.<suffix>
	cnames[1] = dns.CName{
		Primary: true,
		CName:   zkRouter.CName,
		IP:      zkRouter.IP,
	}
	zkRouter.RecordIds[1] = cnames[1].Id()
	// PRIMARY router.<region+zone>.<suffix>
	cnames[2] = dns.CName{
		Primary: true,
		CName:   helper.GetZoneRouterCName(zkRouter.Zone, dns.Provider.Suffix()),
		IP:      zkRouter.IP,
	}
	zkRouter.RecordIds[2] = cnames[2].Id()
	// SECONDARY router.<region+zone>.<suffix>
	for _, azone := range AvailableZones {
		if azone == zone {
			continue
		}
		cname := dns.CName{
			Primary: false,
			CName:   helper.GetZoneRouterCName(azone, dns.Provider.Suffix()),
			IP:      zkRouter.IP,
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

func Unregister(zone, ip string) error {
	zkRouter, err := datamodel.GetRouter(zone, ip)
	if err != nil {
		return err
	}
	if dns.Provider == nil {
		// if we have no dns provider then just save here
		return zkRouter.Delete()
	}
	// delete the basic entries for the new router
	// e.g. region is us-east-1 && zone is us-east-1a && there are zones a, b, c, and d && this is the 2nd in a
	//   delete router2.us-east-1a.<suffix> primary
	//   delete router.us-east-1a.<suffix>  primary
	//   delete router.us-east-1b.<suffix>  secondary
	//   delete router.us-east-1c.<suffix>  secondary
	//   delete router.us-east-1d.<suffix>  secondary
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
