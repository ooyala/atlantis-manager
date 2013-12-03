package manager

import (
	"atlantis/manager/datamodel"
	"atlantis/manager/dns"
	"atlantis/manager/helper"
	. "atlantis/manager/rpc/client"
	. "atlantis/manager/rpc/types"
	"errors"
)

var Port string

func Init(port string) {
	Port = port
}

func Register(region, ip, managerCName, registryCName string) (*datamodel.ZkManager, error) {
	if tmpManager, err := datamodel.GetManager(region, ip); err == nil {
		return tmpManager, errors.New("Already registered.")
	}

	// NOTE[jigish]: health check removed because we can't actually do it security-group wise.

	// set up datamodel
	zkManager := datamodel.Manager(region, ip)
	err := zkManager.Save()
	if err != nil {
		return zkManager, err
	}

	// pass through specified cnames
	if managerCName != "" {
		zkManager.ManagerCName = managerCName
	}
	if registryCName != "" {
		zkManager.RegistryCName = registryCName
	}
	err = zkManager.Save()
	if err != nil {
		return zkManager, err
	}

	if dns.Provider == nil {
		return zkManager, nil
	}

	// set up unspecified cnames
	// first delete all entries we may already have for this IP in DNS
	err = dns.DeleteRecordsForIP(ip)
	if err != nil {
		return zkManager, err
	}
	// choose cnames
	managers, err := datamodel.ListManagersInRegion(region)
	if err != nil {
		return zkManager, err
	}
	managerMap := map[string]bool{}
	registryMap := map[string]bool{}
	for _, manager := range managers {
		tmpManager, err := datamodel.GetManager(region, manager)
		if err != nil {
			return zkManager, err
		}
		managerMap[tmpManager.ManagerCName] = true
		registryMap[tmpManager.RegistryCName] = true
	}

	cnames := []dns.CName{}
	if zkManager.ManagerCName == "" {
		managerNum := 1
		zkManager.ManagerCName = helper.GetManagerCName(managerNum, region, dns.Provider.Suffix())
		for ; managerMap[zkManager.ManagerCName]; managerNum++ {
			zkManager.ManagerCName = helper.GetManagerCName(managerNum, region, dns.Provider.Suffix())
		}
		// create basic health check for manager
		zkManager.ManagerHealthCheckId, err = dns.Provider.CreateHealthCheck(zkManager.IP, uint16(443))
		if err != nil {
			return zkManager, err
		}
		// managerX.<region>.<suffix>
		cname := dns.CName{
			CName:         zkManager.ManagerCName,
			IP:            zkManager.IP,
			HealthCheckId: zkManager.ManagerHealthCheckId,
		}
		zkManager.ManagerRecordId = cname.Id()
		cnames = append(cnames, cname)
	}
	if zkManager.RegistryCName == "" {
		registryNum := 1
		zkManager.RegistryCName = helper.GetRegistryCName(registryNum, region, dns.Provider.Suffix())
		for ; registryMap[zkManager.RegistryCName]; registryNum++ {
			zkManager.RegistryCName = helper.GetRegistryCName(registryNum, region, dns.Provider.Suffix())
		}
		// create basic health check for registry
		zkManager.RegistryHealthCheckId, err = dns.Provider.CreateHealthCheck(zkManager.IP, uint16(80))
		if err != nil {
			return zkManager, err
		}
		// registryX.<region>.<suffix>
		cname := dns.CName{
			CName:         zkManager.RegistryCName,
			IP:            zkManager.IP,
			HealthCheckId: zkManager.RegistryHealthCheckId,
		}
		zkManager.RegistryRecordId = cname.Id()
		cnames = append(cnames, cname)
	}

	if len(cnames) == 0 {
		return zkManager, nil
	}
	err, errChan := dns.Provider.CreateCNames("CREATE_MANAGER "+ip+" in "+region, cnames)
	if err != nil {
		return zkManager, err
	}
	err = <-errChan // wait for change to propagate
	if err != nil {
		return zkManager, err
	}
	return zkManager, zkManager.Save()
}

func Unregister(region, ip string) error {
	zkManager, err := datamodel.GetManager(region, ip)
	if err != nil {
		return err
	}
	if dns.Provider == nil {
		// if we have no dns provider then just save here
		return zkManager.Delete()
	}
	err, errChan := dns.Provider.DeleteRecords("DELETE_MANAGER "+ip+" in "+region, zkManager.ManagerRecordId,
		zkManager.RegistryRecordId)
	if err != nil {
		return err
	}
	err = <-errChan // wait for it to propagate
	if err != nil {
		return err
	}
	// delete basic health check for manager
	err = dns.Provider.DeleteHealthCheck(zkManager.ManagerHealthCheckId)
	if err != nil {
		return err
	}
	// delete basic health check for registry
	err = dns.Provider.DeleteHealthCheck(zkManager.RegistryHealthCheckId)
	if err != nil {
		return err
	}
	return zkManager.Delete()
}

func HealthCheck(ip string) (*ManagerHealthCheckReply, error) {
	args := ManagerHealthCheckArg{}
	var reply ManagerHealthCheckReply
	return &reply, NewManagerRPCClient(ip+":"+Port).Call("HealthCheck", args, &reply)
}
