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

func Register(region, value, registryCName, managerCName string) (*datamodel.ZkManager, error) {
	if tmpManager, err := datamodel.GetManager(region, value); err == nil {
		return tmpManager, errors.New("Already registered.")
	}

	// NOTE[jigish]: health check removed because we can't actually do it security-group wise.

	// set up datamodel
	zkManager := datamodel.Manager(region, value)
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

	suffix, err := dns.Provider.Suffix(region)
	if err != nil {
		return nil, err
	}


	// set up unspecified cnames
	// first delete all entries we may already have for this Value in DNS
	err = dns.DeleteRecordsForValue(region, value)
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

	cnames := []dns.Record{}
	if zkManager.ManagerCName == "" {
		managerNum := 1
		zkManager.ManagerCName = helper.GetManagerCName(managerNum, suffix)
		for ; managerMap[zkManager.ManagerCName]; managerNum++ {
			zkManager.ManagerCName = helper.GetManagerCName(managerNum, suffix)
		}
		// managerX.<region>.<suffix>
		cname := dns.NewRecord(zkManager.ManagerCName, zkManager.Host, 0)
		zkManager.ManagerRecordId = cname.Id()
		cnames = append(cnames, cname)
	}
	if zkManager.RegistryCName == "" {
		registryNum := 1
		zkManager.RegistryCName = helper.GetRegistryCName(registryNum, suffix)
		for ; registryMap[zkManager.RegistryCName]; registryNum++ {
			zkManager.RegistryCName = helper.GetRegistryCName(registryNum, suffix)
		}
		// registryX.<region>.<suffix>
		cname := dns.NewRecord(zkManager.RegistryCName, zkManager.Host, 0)
		zkManager.RegistryRecordId = cname.Id()
		cnames = append(cnames, cname)
	}

	if len(cnames) == 0 {
		return zkManager, nil
	}
	err = dns.Provider.CreateRecords(region, "CREATE_MANAGER "+zkManager.Host+" in "+region, cnames)
	if err != nil {
		return zkManager, err
	}
	return zkManager, zkManager.Save()
}

func Unregister(region, value string) error {
	zkManager, err := datamodel.GetManager(region, value)
	if err != nil {
		return err
	}
	if dns.Provider == nil {
		// if we have no dns provider then just save here
		return zkManager.Delete()
	}
	records := []string{}
	if zkManager.ManagerRecordId != "" {
		records = append(records, zkManager.ManagerRecordId)
	}
	if zkManager.RegistryRecordId != "" {
		records = append(records, zkManager.RegistryRecordId)
	}
	if len(records) > 0 {
		err, errChan := dns.Provider.DeleteRecords(region, "DELETE_MANAGER "+value+" in "+region, records...)
		if err != nil {
			return err
		}
		err = <-errChan // wait for it to propagate
		if err != nil {
			return err
		}
	}
	return zkManager.Delete()
}

func HealthCheck(host string) (*ManagerHealthCheckReply, error) {
	args := ManagerHealthCheckArg{}
	var reply ManagerHealthCheckReply
	return &reply, NewManagerRPCClient(host+":"+Port).Call("HealthCheck", args, &reply)
}
