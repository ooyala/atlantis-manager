package rpc

import (
	. "atlantis/common"
	"atlantis/crypto"
	. "atlantis/manager/constant"
	"atlantis/manager/datamodel"
	"atlantis/manager/dns"
	"atlantis/manager/helper"
	. "atlantis/manager/rpc/types"
	"atlantis/manager/supervisor"
	. "atlantis/supervisor/rpc/types"
	"errors"
	"fmt"
)

//
// Deploy Stuff
//

func deployContainer(auth *ManagerAuthArg, cont *Container, instances uint, t *Task) ([]*Container, error) {
	manifest := cont.Manifest
	manifest.Instances = instances
	return deploy(auth, manifest, cont.Sha, cont.Env, t)
}

func ResolveDepValuesForZone(zkEnv *datamodel.ZkEnv, zone string, names []string, encrypt bool) (map[string]string, error) {
	deps := map[string]string{}
	leftoverNames := []string{}
	// if we're using DNS and the app is registered, try to get the app cname (if deployed)
	if dns.Provider != nil {
		for _, name := range names {
			if datamodel.InternalAppExistsInEnv(name, zkEnv.Name) {
				// this is a registered, internal, deployed app, output the cname
				suffix, err := dns.Provider.Suffix(Region)
				if err != nil {
					leftoverNames = append(leftoverNames, name)
					continue
				}
				deps[name] = helper.GetZoneAppCName(name, zkEnv.Name, zone, suffix)
			} else {
				leftoverNames = append(leftoverNames, name)
			}
		}
	} else {
		leftoverNames = names
	}
	envDeps, err := zkEnv.ResolveDepValues(leftoverNames)
	if err != nil {
		return nil, err
	}
	var retDeps map[string]string
	if encrypt {
		for name, value := range deps {
			envDeps[name] = string(crypto.Encrypt([]byte(value)))
		}
		retDeps = envDeps
	} else {
		for name, value := range envDeps {
			deps[name] = string(crypto.Decrypt([]byte(value)))
		}
		retDeps = deps
	}
	for _, name := range names {
		if _, ok := retDeps[name]; !ok {
			return retDeps, errors.New("Could not resolve dep " + name)
		}
	}
	return retDeps, nil
}

func ResolveDepValues(zkEnv *datamodel.ZkEnv, names []string, encrypt bool) (deps map[string]map[string]string, err error) {
	deps = map[string]map[string]string{}
	for _, zone := range AvailableZones {
		deps[zone], err = ResolveDepValuesForZone(zkEnv, zone, names, encrypt)
		if err != nil {
			return nil, errors.New("Dependency Error: " + err.Error())
		}
	}
	return deps, nil
}

func validateDeploy(auth *ManagerAuthArg, manifest *Manifest, sha, env string, t *Task) (deps map[string]map[string]string, err error) {
	// authorize that we're allowed to use the app
	if err = AuthorizeApp(auth, manifest.Name); err != nil {
		return nil, errors.New("Permission Denied: " + err.Error())
	}
	// fetch the environment
	t.LogStatus("Fetching Environment")
	zkEnv, err := datamodel.GetEnv(env)
	if err != nil {
		return nil, errors.New("Environment Error: " + err.Error())
	}
	// lock the deploy
	dl := datamodel.NewDeployLock(t.Id, manifest.Name, sha, env)
	if err := dl.Lock(); err != nil {
		return nil, err
	}
	defer dl.Unlock()
	if manifest.Instances <= 0 {
		return nil, errors.New(fmt.Sprintf("Invalid Number of Instances: %d", manifest.Instances))
	}
	if manifest.CPUShares < 0 ||
		(manifest.CPUShares > 0 && manifest.CPUShares != 1 && manifest.CPUShares%CPUSharesIncrement != 0) {
		return nil, errors.New(fmt.Sprintf("CPU Shares should be 1 or a multiple of %d", CPUSharesIncrement))
	}
	if manifest.MemoryLimit < 0 ||
		(manifest.MemoryLimit > 0 && manifest.MemoryLimit%MemoryLimitIncrement != 0) {
		return nil, errors.New(fmt.Sprintf("Memory Limit should be a multiple of %d", MemoryLimitIncrement))
	}
	t.LogStatus("Resolving Dependencies")
	return ResolveDepValues(zkEnv, manifest.DepNames(), true)
}

type DeployHostResult struct {
	Host      string
	Container *Container
	Error     error
}

func deployToHost(respCh chan *DeployHostResult, manifest *Manifest, sha, env, host string) {
	instance, err := datamodel.CreateInstance(manifest.Internal, manifest.Name, sha, env, host)
	if err != nil {
		respCh <- &DeployHostResult{Host: host, Container: nil, Error: err}
		return
	}
	ihReply, err := supervisor.Deploy(host, manifest.Name, sha, env, instance.Id, manifest)
	if err != nil {
		instance.Delete()
		respCh <- &DeployHostResult{Host: host, Container: nil, Error: err}
		return
	}
	if ihReply.Status != StatusOk {
		instance.Delete()
		respCh <- &DeployHostResult{Host: host, Container: nil, Error: err}
		return
	}
	ihReply.Container.Host = host
	instance.SetPort(ihReply.Container.PrimaryPort)
	AddAppShaToEnv(manifest.Name, sha, env)
	respCh <- &DeployHostResult{Host: host, Container: ihReply.Container, Error: nil}
}

type DeployZoneResult struct {
	Zone       string
	Containers []*Container
	Error      error
}

func deployToZone(respCh chan *DeployZoneResult, manifest *Manifest, sha, env string, hosts []string, zone string) {
	hostNum := 0
	failures := 0
	deployed := uint(0)
	maxFailures := len(hosts)
	deployedContainers := []*Container{}
	for deployed < manifest.Instances && failures < maxFailures {
		numToDeploy := manifest.Instances - deployed
		respCh := make(chan *DeployHostResult, numToDeploy)
		for i := uint(0); i < numToDeploy; i++ {
			host := hosts[hostNum]
			go deployToHost(respCh, manifest, sha, env, host)
			hostNum++
			if hostNum >= len(hosts) {
				hostNum = 0
			}
		}
		numResult := uint(0)
		for result := range respCh {
			if result.Error != nil {
				failures++
			} else {
				deployed++
				deployedContainers = append(deployedContainers, result.Container)
			}
			numResult++
			if numResult >= numToDeploy { // we're done
				close(respCh)
			}
		}
	}
	if failures >= maxFailures {
		respCh <- &DeployZoneResult{
			Zone:       zone,
			Containers: deployedContainers,
			Error:      errors.New(fmt.Sprintf("Failed to deploy %d instances in zone %s.", manifest.Instances, zone)),
		}
		return
	}
	respCh <- &DeployZoneResult{
		Zone:       zone,
		Containers: deployedContainers,
		Error:      nil,
	}
	return
}

func deployToHostsInZones(deps map[string]map[string]string, manifest *Manifest, sha, env string,
	hosts map[string][]string, zones []string, t *Task) ([]*Container, error) {
	deployedContainers := []*Container{}
	// first check if zones have enough hosts
	for _, zone := range zones {
		// fail if zone has no hosts
		if hosts[zone] == nil || len(hosts[zone]) == 0 {
			return nil, errors.New(fmt.Sprintf("No hosts available for app %s in zone %s", manifest.Name, zone))
		}
	}
	// now that we know that enough hosts are available
	t.LogStatus("Deploying to: %v", zones)
	respCh := make(chan *DeployZoneResult, len(zones))
	for _, zone := range zones {
		zoneManifest := manifest.Dup()
		zoneManifest.Deps = deps[zone]
		go deployToZone(respCh, zoneManifest, sha, env, hosts[zone], zone)
	}
	var err error
	numResults := 0
	status := "Deployed to: "
	for result := range respCh {
		deployedContainers = append(deployedContainers, result.Containers...)
		if result.Error != nil {
			err = result.Error
			t.Log(err.Error())
			status += result.Zone + ":FAIL "
		} else {
			status += result.Zone + ":SUCCESS "
		}
		t.LogStatus(status)
		numResults++
		if numResults >= len(zones) { // we're done
			close(respCh)
		}
	}
	close(respCh)
	if err != nil {
		cleanup(false, deployedContainers, t)
		return nil, err
	}

	// set ports on zk supervisor - can't do this in parallel. we may deploy to the same host at the same time
	// and since we don't lock zookeeper (maybe we should), this would result in a race condition.
	t.LogStatus("Updating Zookeeper")
	for _, cont := range deployedContainers {
		datamodel.Supervisor(cont.Host).SetContainerAndPort(cont.Id, cont.PrimaryPort)
	}

	// we're good now, so lets move on
	t.LogStatus("Updating Router")
	deployedIds := make([]string, len(deployedContainers))
	count := 0
	for _, cont := range deployedContainers {
		deployedIds[count] = cont.Id
		count++
	}
	err = datamodel.AddToPool(deployedIds)
	if err != nil { // if we can't add the pool, clean up and fail
		cleanup(true, deployedContainers, t)
		return nil, errors.New("Update Pool Error: " + err.Error())
	}
	if manifest.Internal {
		t.LogStatus("Updating DNS")
		// if we're internal, handle the DNS stuff
		err := dns.CreateAppCNames(manifest.Internal, manifest.Name, sha, env)
		if err != nil { // if DNS fails, clean up and fail
			cleanup(true, deployedContainers, t)
			return nil, errors.New("Update DNS Error: " + err.Error())
		}
	}
	return deployedContainers, nil
}

func devDeployToHosts(deps map[string]map[string]string, manifest *Manifest, sha, env string, hosts []string, t *Task) ([]*Container, error) {
	// deploy to hosts
	deployedContainers := []*Container{}
	deployedIds := []string{}
	for _, host := range hosts {
		health, err := supervisor.HealthCheck(host)
		if err != nil {
			continue
		}
		manifest.Deps = deps[health.Zone]
		instance, err := datamodel.CreateInstance(manifest.Internal, manifest.Name, sha, env, host)
		if err != nil {
			continue
		}
		t.LogStatus("Deploying %s to %s", instance.Id, host)
		ihReply, err := supervisor.Deploy(host, manifest.Name, sha, env, instance.Id, manifest)
		if err != nil {
			instance.Delete()
			t.LogStatus("Supervisor " + host + " Deploy " + instance.Id + " Failed: " + err.Error())
			continue // try another host
		}
		if ihReply.Status != StatusOk {
			instance.Delete()
			t.LogStatus("Supervisor " + host + " Deploy " + instance.Id + " Status Not OK: " + ihReply.Status)
			continue // try another host
		}
		ihReply.Container.Host = host
		instance.SetPort(ihReply.Container.PrimaryPort)
		datamodel.Supervisor(host).SetContainerAndPort(instance.Id, ihReply.Container.PrimaryPort)
		deployedContainers = append(deployedContainers, ihReply.Container)
		deployedIds = append(deployedIds, ihReply.Container.Id)
		AddAppShaToEnv(manifest.Name, sha, env)
		break // only deploy 1
	}
	if len(deployedContainers) == 0 {
		return nil, errors.New("Could not deploy to any hosts")
	}
	t.LogStatus("Updating Router")
	err := datamodel.AddToPool(deployedIds)
	if err != nil { // if we can't add the pool, clean up and fail
		cleanup(true, deployedContainers, t)
		return nil, errors.New("Update Pool Error: " + err.Error())
	}
	if manifest.Internal {
		t.LogStatus("Updating DNS")
		// if we're internal, handle the DNS stuff
		err := dns.CreateAppCNames(manifest.Internal, manifest.Name, sha, env)
		if err != nil { // if DNS fails, clean up and fail
			cleanup(true, deployedContainers, t)
			return nil, errors.New("Update DNS Error: " + err.Error())
		}
	}
	return deployedContainers, nil
}

func deploy(auth *ManagerAuthArg, manifest *Manifest, sha, env string, t *Task) ([]*Container, error) {
	deps, err := validateDeploy(auth, manifest, sha, env, t)
	if err != nil {
		return nil, err
	}
	// choose hosts
	t.LogStatus("Choosing Supervisors")
	hosts, err := datamodel.ChooseSupervisors(manifest.Name, sha, env, manifest.Instances, manifest.CPUShares,
		manifest.MemoryLimit, AvailableZones, map[string]bool{})
	if err != nil {
		return nil, errors.New("Choose Supervisors Error: " + err.Error())
	}
	return deployToHostsInZones(deps, manifest, sha, env, hosts, AvailableZones, t)
}

func devDeploy(auth *ManagerAuthArg, manifest *Manifest, sha, env string, t *Task) ([]*Container, error) {
	manifest.Instances = 1 // set to 1 instance regardless of what came in
	deps, err := validateDeploy(auth, manifest, sha, env, t)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	// choose hosts
	t.LogStatus("Choosing Supervisors")
	list, err := datamodel.ChooseSupervisorsList(manifest.Name, sha, env, manifest.CPUShares, manifest.MemoryLimit,
		AvailableZones, map[string]bool{})
	if err != nil {
		return nil, errors.New("Choose Supervisors Error: " + err.Error())
	}
	hosts := make([]string, len(list))
	for i, elem := range list {
		hosts[i] = elem.Supervisor
	}
	return devDeployToHosts(deps, manifest, sha, env, hosts, t)
}

func moveContainer(auth *ManagerAuthArg, cont *Container, t *Task) (*Container, error) {
	manifest := cont.Manifest
	manifest.Instances = 1 // we only want 1 instance
	deps, err := validateDeploy(auth, manifest, cont.Sha, cont.Env, t)
	if err != nil {
		return nil, err
	}
	// choose host
	t.LogStatus("Choosing Supervisor")
	zone, err := supervisor.GetZone(cont.Host)
	if err != nil {
		return nil, err
	}
	hosts, err := datamodel.ChooseSupervisors(manifest.Name, cont.Sha, cont.Env, manifest.Instances, manifest.CPUShares,
		manifest.MemoryLimit, []string{zone}, map[string]bool{cont.Host: true})
	if err != nil {
		return nil, err
	}
	deployed, err := deployToHostsInZones(deps, manifest, cont.Sha, cont.Env, hosts, []string{zone}, t)
	if err != nil {
		return nil, err
	}
	// should only deploy 1 since we're only moving 1
	if len(deployed) != 1 {
		cleanup(true, deployed, t)
		return nil, errors.New(fmt.Sprintf("Didn't deploy 1 container. Deployed %d", len(deployed)))
	}
	cleanup(true, []*Container{cont}, t) // cleanup the old container
	return deployed[0], nil
}

func cleanup(removeContainerFromHost bool, deployedContainers []*Container, t *Task) {
	// kill all references to deployed containers as well as the container itself
	for _, container := range deployedContainers {
		supervisor.Teardown(container.Host, []string{container.Id}, false)
		if instance, err := datamodel.GetInstance(container.Id); err == nil {
			instance.Delete()
		} else {
			t.Log(fmt.Sprintf("Failed to clean up instance %s: %s", container.Id, err.Error()))
		}
		DeleteAppShaFromEnv(container.App, container.Sha, container.Env)
		if removeContainerFromHost {
			datamodel.Supervisor(container.Host).RemoveContainer(container.Id)
		}
	}
}

//
// Teardown Stuff
//

func getContainerIdsOfEnv(t *Task, app, sha, env string) ([]string, error) {
	containerIds, err := datamodel.ListInstances(app, sha, env)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error listing instances of %s @ %s in %s: %s", app, sha, env,
			err.Error()))
	}
	return containerIds, nil
}

func getContainerIdsOfSha(t *Task, app, sha string) ([]string, error) {
	containerIds := []string{}
	envs, err := datamodel.ListAppEnvs(app, sha)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error listing environments of %s @ %s: %s", app, sha, err.Error()))
	}
	for _, env := range envs {
		tmpContainerIds, err := getContainerIdsOfEnv(t, app, sha, env)
		if err != nil {
			return nil, err
		}
		containerIds = append(containerIds, tmpContainerIds...)
	}
	return containerIds, nil
}

func getContainerIdsOfApp(t *Task, app string) ([]string, error) {
	containerIds := []string{}
	shas, err := datamodel.ListShas(app)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error listing shas of %s: %s", app, err.Error()))
	}
	for _, sha := range shas {
		tmpContainerIds, err := getContainerIdsOfSha(t, app, sha)
		if err != nil {
			return nil, err
		}
		containerIds = append(containerIds, tmpContainerIds...)
	}
	return containerIds, nil
}

func getContainerIdsToTeardown(t *Task, arg ManagerTeardownArg) (hostMap map[string][]string, err error) {
	hostMap = map[string][]string{} // map of host -> []string container ids
	if arg.All {
		var hosts []string
		hosts, err = datamodel.ListSupervisors()
		if err != nil {
			return nil, errors.New("Error listing hosts: " + err.Error())
		}
		for _, host := range hosts {
			hostMap[host] = []string{}
		}
		return
	} else if arg.ContainerId != "" {
		var instance *datamodel.ZkInstance
		instance, err = datamodel.GetInstance(arg.ContainerId)
		if err != nil {
			return
		}
		hostMap[instance.Host] = []string{arg.ContainerId}
		return
	} else if arg.App != "" {
		containerIds := []string{}
		if arg.Sha != "" {
			if arg.Env != "" {
				if containerIds, err = getContainerIdsOfEnv(t, arg.App, arg.Sha, arg.Env); err != nil {
					return nil, err
				}
			} else {
				if containerIds, err = getContainerIdsOfSha(t, arg.App, arg.Sha); err != nil {
					return nil, err
				}
			}
		} else {
			if containerIds, err = getContainerIdsOfApp(t, arg.App); err != nil {
				return nil, err
			}
		}
		for _, containerId := range containerIds {
			instance, err := datamodel.GetInstance(containerId)
			if err != nil {
				continue
			}
			currentIds := hostMap[instance.Host]
			if currentIds == nil {
				hostMap[instance.Host] = []string{containerId}
			} else {
				hostMap[instance.Host] = append(currentIds, containerId)
			}
		}
		return
	}
	return nil, errors.New("Invalid Arguments")
}
