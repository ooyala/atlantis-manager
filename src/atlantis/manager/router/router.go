package router

import (
	. "atlantis/manager/constant"
	"atlantis/manager/datamodel"
	"atlantis/manager/dns"
	"atlantis/manager/helper"
)

// Registering the 4th router in us-east-1a, with host (say) one.two.com needs to create the following
// entries in DNS (route53 in our case), assuming that atlantis.com is the DNS zone delegated to
// the deployment. Also assume that the deployment spans zones a, b and c in region us-east-1.
//
// - A weight 1 CNAME for value router4.us-east-1a.atlantis.com pointing at one.two.com
// - A weight 1 CNAME for value router.us-east-1a.atlantis.com pointing at one.two.com
// - A weight 0 CNAME for value router.us-east-1b.atlantis.com pointing to one.two.com
// - A weight 0 CNAME for value router.us-east-1c.atlantis.com pointing to one.two.com
// - A weight 1 CNAME for value router.us-east-1.atlantis.com pointing at one.two.com
//
// The weight 0 A records will be activated if health checks fail on the router in that zone (not implemented)
//
// Deleting the router, simply deletes all the records created when adding it.

func preCreateRecordSets(internal bool, zone, host string, zkRouter *datamodel.ZkRouter) (string, error) {
	// we're only creating records for this region
	suffix, err := dns.Provider.Suffix(Region)
	if err != nil {
		return "", err
	}
	// first delete all entries we may already have for this Value in DNS
	err = dns.DeleteRecordsForValue(Region, host)
	if err != nil {
		return "", err
	}
	// choose cname
	routers, err := datamodel.ListRoutersInZone(internal, zone)
	if err != nil {
		return "", err
	}
	routerMap := map[string]bool{}
	for _, router := range routers {
		tmpRouter, err := datamodel.GetRouter(internal, zone, router)
		if err != nil {
			return "", err
		}
		routerMap[tmpRouter.CName] = true
	}
	routerNum := 1
	myCName := helper.GetRouterCName(internal, routerNum, zone, suffix)
	for ; routerMap[myCName]; routerNum++ {
		myCName = helper.GetRouterCName(internal, routerNum, zone, suffix)
	}
	zkRouter.CName = myCName
	return suffix, nil
}

func createRecordSets(internal bool, zone, value string, zkRouter *datamodel.ZkRouter) ([]dns.Record, error) {
	var suffix string
	var err error
	if suffix, err = preCreateRecordSets(internal, zone, value, zkRouter); err != nil {
		return nil, err
	}

	records := make([]dns.Record, 3)

	// WEIGHT=1 router.<suffix>
	records[0] = dns.NewRecord(helper.GetRegionRouterCName(internal, suffix), value, 1)

	// WEIGHT=1 routerX.<zone>.<suffix>
	records[1] = dns.NewRecord(zkRouter.CName, value, 1)

	// WEIGHT=1 router.<zone>.<suffix>
	records[2] = dns.NewRecord(helper.GetZoneRouterCName(internal, zkRouter.Zone, suffix), value, 1)

	/*// WEIGHT=0 router.<zone>.<suffix> -> will be activated when needed
	for _, azone := range AvailableZones {
		if azone == zone {
			continue
		}
		record := dns.NewRecord(helper.GetZoneRouterCName(internal, azone, suffix), value, 0)
		records = append(records, record)
	}*/ // we don't need this yet
	return records, nil
}

func Register(internal bool, zone, host, ip string) (*datamodel.ZkRouter, error) {
	// create ZkRouter
	zkRouter := datamodel.Router(internal, zone, host, ip)
	zkRouter.RecordIDs = []string{}
	if dns.Provider == nil {
		// if we have no dns provider then just save here
		return zkRouter, zkRouter.Save()
	}

	// get record sets
	cnames, err := createRecordSets(internal, zone, host, zkRouter)
	if err != nil {
		return zkRouter, err
	}
	if len(cnames) == 0 { // no need to do anything, there are no cnames to save
		return zkRouter, nil
	}

	err = dns.Provider.CreateRecords(Region, "CREATE_ROUTER "+host+" with ip "+ip+" in "+zone, cnames)
	if err != nil {
		return zkRouter, err
	}

	// add RecordIDs
	for _, cname := range cnames {
		zkRouter.RecordIDs = append(zkRouter.RecordIDs, cname.ID())
	}

	return zkRouter, zkRouter.Save()
}

func Unregister(internal bool, zone, value string) error {
	zkRouter, err := datamodel.GetRouter(internal, zone, value)
	if err != nil {
		return err
	}
	if dns.Provider == nil || len(zkRouter.RecordIDs) == 0 {
		// if we have no dns provider or there aren't any record IDs then just save here
		return zkRouter.Delete()
	}
	err, errChan := dns.Provider.DeleteRecords(Region, "DELETE_ROUTER "+value+" in "+zone, zkRouter.RecordIDs...)
	if err != nil {
		return err
	}
	err = <-errChan // wait for it to propagate
	if err != nil {
		return err
	}
	return zkRouter.Delete()
}
