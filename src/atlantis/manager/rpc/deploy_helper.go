package rpc

import (
	. "atlantis/common"
	. "atlantis/manager/constant"
	"atlantis/manager/datamodel"
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

func validateDeploy(auth *ManagerAuthArg, manifest *Manifest, sha, env string, t *Task) (err error) {
	// authorize that we're allowed to use the app
	if err = AuthorizeApp(auth, manifest.Name); err != nil {
		return errors.New("Permission Denied: " + err.Error())
	}
	// fetch the environment
	t.LogStatus("Fetching Environment")
	zkEnv, err := datamodel.GetEnv(env)
	if err != nil {
		return errors.New("Environment Error: " + err.Error())
	}
	// lock the deploy
	dl := datamodel.NewDeployLock(t.Id, manifest.Name, sha, env)
	if err := dl.Lock(); err != nil {
		return err
	}
	defer dl.Unlock()
	if manifest.Instances <= 0 {
		return errors.New(fmt.Sprintf("Invalid Number of Instances: %d", manifest.Instances))
	}
	if manifest.CPUShares < 0 ||
		(manifest.CPUShares > 0 && manifest.CPUShares != 1 && manifest.CPUShares%CPUSharesIncrement != 0) {
		return errors.New(fmt.Sprintf("CPU Shares should be 1 or a multiple of %d", CPUSharesIncrement))
	}
	if manifest.MemoryLimit < 0 ||
		(manifest.MemoryLimit > 0 && manifest.MemoryLimit%MemoryLimitIncrement != 0) {
		return errors.New(fmt.Sprintf("Memory Limit should be a multiple of %d", MemoryLimitIncrement))
	}
	t.LogStatus("Resolving Dependencies")
	manifest.Deps, err = zkEnv.ResolveDepValues(manifest.DepNames())
	if err != nil {
		return errors.New("Dependency Error: " + err.Error())
	}
	return nil
}

func deployToHostsInZones(manifest *Manifest, sha, env string, hosts map[string][]string, zones []string,
	t *Task) ([]*Container, error) {
	// deploy to hosts
	deployedContainers := []*Container{}
	deployedIds := []string{}
	for _, zone := range zones {
		// fail if zone has no hosts
		if hosts[zone] == nil || len(hosts[zone]) == 0 {
			cleanup(deployedContainers, t)
			return nil, errors.New(fmt.Sprintf("No hosts available for app %s in zone %s", manifest.Name, zone))
		}
		failures := 0
		maxFailures := len(hosts[zone]) // allow as many failures as there are hosts
		hostNum := 0
		deployed := uint(0)
		// keep deploying to hosts until we've deployed the right number in this zone OR we've failed too much
		for deployed < manifest.Instances && failures < maxFailures {
			host := hosts[zone][hostNum]
			instance, err := datamodel.CreateInstance(manifest.Internal, manifest.Name, sha, env, host)
			t.LogStatus("Deploying %s to %s", instance.Id, host)
			ihReply, err := supervisor.Deploy(host, manifest.Name, sha, env, instance.Id, manifest)
			if err != nil {
				instance.Delete()
				t.LogStatus("Supervisor " + host + " Deploy " + instance.Id + " Failed: " + err.Error())
				failures++
				continue // try another host
			}
			if ihReply.Status != StatusOk {
				instance.Delete()
				t.LogStatus("Supervisor " + host + " Deploy " + instance.Id + " Status Not OK: " + ihReply.Status)
				failures++
				continue // try another host
			}
			ihReply.Container.Host = host
			instance.SetPort(ihReply.Container.PrimaryPort)
			datamodel.Host(host).SetContainerAndPort(instance.Id, ihReply.Container.PrimaryPort)
			deployedContainers = append(deployedContainers, ihReply.Container)
			deployedIds = append(deployedIds, ihReply.Container.Id)
			AddAppShaToEnv(manifest.Name, sha, env)
			deployed++
			hostNum++ // increment through hosts to cycle through them. If we've gone too far, reset
			if hostNum >= len(hosts[zone]) {
				hostNum = 0
			}
		}
		if failures >= maxFailures { // if we've failed out for this zone, clean up and fail
			cleanup(deployedContainers, t)
			return nil, errors.New(fmt.Sprintf("Failed to deploy %d instances to zone %s", manifest.Instances, zone))
		}
	}
	t.LogStatus("Updating Router")
	err := datamodel.AddToPool(deployedIds)
	if err != nil { // if we can't add the pool, clean up and fail
		cleanup(deployedContainers, t)
		return nil, errors.New("Update Pool Error: " + err.Error())
	}
	return deployedContainers, nil
}

func deploy(auth *ManagerAuthArg, manifest *Manifest, sha, env string, t *Task) ([]*Container, error) {
	if err := validateDeploy(auth, manifest, sha, env, t); err != nil {
		return nil, err
	}
	// choose hosts
	t.LogStatus("Choosing Hosts")
	hosts, err := datamodel.ChooseHosts(manifest.Name, sha, env, manifest.Instances, manifest.CPUShares,
		manifest.MemoryLimit, AvailableZones, map[string]bool{})
	if err != nil {
		return nil, errors.New("Choose Hosts Error: " + err.Error())
	}
	return deployToHostsInZones(manifest, sha, env, hosts, AvailableZones, t)
}

func moveContainer(auth *ManagerAuthArg, cont *Container, t *Task) (*Container, error) {
	manifest := cont.Manifest
	manifest.Instances = 1 // we only want 1 instance
	if err := validateDeploy(auth, manifest, cont.Sha, cont.Env, t); err != nil {
		return nil, err
	}
	// choose host
	t.LogStatus("Choosing Host")
	zone, err := supervisor.GetZone(cont.Host)
	if err != nil {
		return nil, err
	}
	hosts, err := datamodel.ChooseHosts(manifest.Name, cont.Sha, cont.Env, manifest.Instances, manifest.CPUShares,
		manifest.MemoryLimit, []string{zone}, map[string]bool{cont.Host: true})
	if err != nil {
		return nil, err
	}
	deployed, err := deployToHostsInZones(manifest, cont.Sha, cont.Env, hosts, []string{zone}, t)
	if err != nil {
		return nil, err
	}
	// should only deploy 1 since we're only moving 1
	if len(deployed) != 1 {
		cleanup(deployed, t)
		return nil, errors.New(fmt.Sprintf("Didn't deploy 1 container. Deployed %d", len(deployed)))
	}
	cleanup([]*Container{cont}, t) // cleanup the old container
	return deployed[0], nil
}

func cleanup(deployedContainers []*Container, t *Task) {
	// kill all references to deployed containers as well as the container itself
	for _, container := range deployedContainers {
		supervisor.Teardown(container.Host, []string{container.Id}, false)
		if instance, err := datamodel.GetInstance(container.Id); err == nil {
			instance.Delete()
		} else {
			t.Log(fmt.Sprintf("Failed to clean up instance %s: %s", container.Id, err.Error()))
		}
		DeleteAppShaFromEnv(container.App, container.Sha, container.Env)
		datamodel.Host(container.Host).RemoveContainer(container.Id)
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
		hosts, err = datamodel.ListHosts()
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
