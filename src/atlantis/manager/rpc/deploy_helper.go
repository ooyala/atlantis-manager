/* Copyright 2014 Ooyala, Inc. All rights reserved.
 *
 * This file is licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
 * except in compliance with the License. You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is
 * distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

package rpc

import (
	. "atlantis/common"
	. "atlantis/manager/constant"
	"atlantis/manager/datamodel"
	"atlantis/manager/dns"
	"atlantis/manager/helper"
	"atlantis/manager/netsec"
	. "atlantis/manager/rpc/types"
	"atlantis/manager/supervisor"
	"atlantis/supervisor/crypto"
	. "atlantis/supervisor/rpc/types"
	"errors"
	"fmt"
	"strconv"
)

//
// Deploy Stuff
//

func deployContainer(auth *ManagerAuthArg, cont *Container, instances uint, t *Task) ([]*Container, error) {
	manifest := cont.Manifest
	manifest.Instances = instances
	return deploy(auth, manifest, cont.Sha, cont.Env, t)
}

func MergeDependerEnvData(dst *DependerEnvData, src *DependerEnvData) *DependerEnvData {
	data := &DependerEnvData{
		Name:          dst.Name,
		SecurityGroup: dst.SecurityGroup,
		DataMap:       map[string]interface{}{},
	}
	if dst != nil {
		for key, val := range dst.DataMap {
			data.DataMap[key] = val
		}
	}
	if src != nil {
		for key, val := range src.DataMap {
			data.DataMap[key] = val
		}
		if src.SecurityGroup != nil {
			data.SecurityGroup = src.SecurityGroup
		}
	}
	return data
}

func ResolveDepValuesForZone(app string, zkEnv *datamodel.ZkEnv, zone string, names []string, encrypt bool, t *Task) (DepsType, error) {
	var (
		err    error
	)
	deps := DepsType{}
	// if we're using DNS and the app is registered, try to get the app cname (if deployed)
	if dns.Provider != nil {
		_, err = dns.Provider.Suffix(Region)
		if err != nil {
			return deps, err
		}
	}
	for _, name := range names {
		// if app is registered for this dependency name
		zkApp, err := datamodel.GetApp(name)
		if err != nil {
			continue
		}
		appEnvData := zkApp.GetDependerEnvDataForDependerApp(app, zkEnv.Name, true)
		if appEnvData == nil {
			continue
		}
		envData := zkApp.GetDependerEnvData(zkEnv.Name, true)
		if envData == nil {
			envData = &DependerEnvData{Name: appEnvData.Name}
		}
		// merge the data
		mergedEnvData := MergeDependerEnvData(envData, appEnvData)
		appDep := &AppDep{
			SecurityGroup: mergedEnvData.SecurityGroup,
			DataMap:       mergedEnvData.DataMap,
		}
		if dns.Provider != nil && !zkApp.NonAtlantis && zkApp.Internal {
			// auto-populate Address
			port, created, err := datamodel.ReserveRouterPortAndUpdateTrie(zkApp.Internal, name, "", zkEnv.Name)
			if err != nil {
				return deps, err
			}
			if created {
				// add warning since this means that the app has not been deployed in this env yet
				t.AddWarning("App dependency " + name + " has not yet been deployed in environment " + zkEnv.Name)
			}
			if appDep.DataMap == nil {
				appDep.DataMap = map[string]interface{}{}
			}
			appDep.DataMap["address"] = helper.GetZoneRouterConsulCName(true, zone) + ":" + port

			// auto-populate SecurityGroup
			portUint, err := strconv.ParseUint(port, 10, 16)
			if err != nil {
				return deps, err
			}
			appDep.SecurityGroup = map[string][]uint16{netsec.InternalRouterIPGroup: []uint16{uint16(portUint)}}
		}
		deps[name] = appDep
	}
	if encrypt {
		for _, value := range deps {
			crypto.EncryptAppDep(value)
		}
	}
	for _, name := range names {
		if _, ok := deps[name]; !ok {
			return deps, errors.New("Could not resolve dep " + name)
		}
	}
	return deps, nil
}

func ResolveDepValues(app string, zkEnv *datamodel.ZkEnv, names []string, encrypt bool, t *Task) (deps map[string]DepsType, err error) {
	deps = map[string]DepsType{}
	for _, zone := range AvailableZones {
		deps[zone], err = ResolveDepValuesForZone(app, zkEnv, zone, names, encrypt, t)
		if err != nil {
			return nil, errors.New("Dependency Error: " + err.Error())
		}
	}
	return deps, nil
}

func validateDeploy(auth *ManagerAuthArg, manifest *Manifest, sha, env string, t *Task) (deps map[string]DepsType, err error) {
	t.LogStatus("Validate Deploy")
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
	dl := datamodel.NewDeployLock(t.ID, manifest.Name, sha, env)
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
	return ResolveDepValues(manifest.Name, zkEnv, manifest.DepNames(), true, t)
}

type DeployHostResult struct {
	Host      string
	Container *Container
	Error     error
}

func deployToHost(respCh chan *DeployHostResult, manifest *Manifest, sha, env, host string) {
	instance, err := datamodel.CreateInstance(manifest.Name, sha, env, host)
	if err != nil {
		respCh <- &DeployHostResult{Host: host, Container: nil, Error: err}
		return
	}
	ihReply, err := supervisor.Deploy(host, manifest.Name, sha, env, instance.ID, manifest)
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
	instance.SetManifest(ihReply.Container.Manifest)
	AddAppShaToEnv(manifest.Name, sha, env)
	respCh <- &DeployHostResult{Host: host, Container: ihReply.Container, Error: nil}
}

type DeployZoneResult struct {
	Zone       string
	Containers []*Container
	Error      error
}

func deployToZone(respCh chan *DeployZoneResult, deps map[string]DepsType, rawManifest *Manifest, sha,
	env string, hosts []string, zone string) {
	hostNum := 0
	failures := 0
	deployed := uint(0)
	maxFailures := len(hosts)
	deployedContainers := []*Container{}
	for deployed < rawManifest.Instances && failures < maxFailures {
		numToDeploy := rawManifest.Instances - deployed
		respCh := make(chan *DeployHostResult, numToDeploy)
		for i := uint(0); i < numToDeploy; i++ {
			host := hosts[hostNum]
			// check health on host to figure out its zone to get the deps
			ihReply, err := supervisor.HealthCheck(host)
			if err == nil && ihReply.Status == StatusOk {
				// only try to deploy if health was fine
				// duplicate manifest and get deps
				manifest := rawManifest.Dup()
				manifest.Deps = deps[ihReply.Zone]
				go deployToHost(respCh, manifest, sha, env, host)
			}
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
			Error: errors.New(fmt.Sprintf("Failed to deploy %d instances in zone %s.", rawManifest.Instances,
				zone)),
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

func deployToHostsInZones(deps map[string]DepsType, manifest *Manifest, sha, env string,
	hosts map[string][]string, zones []string, t *Task) ([]*Container, error) {
	deployedContainers := []*Container{}
	// fetch the app
	zkApp, err := datamodel.GetApp(manifest.Name)
	if err != nil {
		return nil, err
	}
	// first check if zones have enough hosts
	for _, zone := range zones {
		// fail if zone has no hosts
		if hosts[zone] == nil || len(hosts[zone]) == 0 {
			return nil, errors.New(fmt.Sprintf("No hosts available for app %s in zone %s", manifest.Name, zone))
		}
	}
	// now that we know that enough hosts are available
	t.LogStatus("Deploying to zones: %v", zones)
	respCh := make(chan *DeployZoneResult, len(zones))
	for _, zone := range zones {
		go deployToZone(respCh, deps, manifest, sha, env, hosts[zone], zone)
	}
	numResults := 0
	status := "Deployed to zones: "
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
	if err != nil {
		cleanup(false, deployedContainers, t)
		return nil, err
	}

	// set ports on zk supervisor - can't do this in parallel. we may deploy to the same host at the same time
	// and since we don't lock zookeeper (maybe we should), this would result in a race condition.
	t.LogStatus("Updating Zookeeper")
	for _, cont := range deployedContainers {
		datamodel.Supervisor(cont.Host).SetContainerAndPort(cont.ID, cont.PrimaryPort)
	}

	// we're good now, so lets move on
	t.LogStatus("Updating Router")
	deployedIDs := make([]string, len(deployedContainers))
	count := 0
	for _, cont := range deployedContainers {
		deployedIDs[count] = cont.ID
		count++
	}
	err = datamodel.AddToPool(deployedIDs)
	if err != nil { // if we can't add the pool, clean up and fail
		cleanup(true, deployedContainers, t)
		return nil, errors.New("Update Pool Error: " + err.Error())
	}
	if zkApp.Internal {
		// reserve router port if needed and add app+env
		_, _, err = datamodel.ReserveRouterPortAndUpdateTrie(zkApp.Internal, manifest.Name, sha, env)
		if err != nil {
			datamodel.DeleteFromPool(deployedIDs)
			cleanup(true, deployedContainers, t)
			return nil, errors.New("Reserve Router Port Error: " + err.Error())
		}
	} else {
		// only update trie
		_, err = datamodel.UpdateAppEnvTrie(zkApp.Internal, manifest.Name, sha, env)
		if err != nil {
			datamodel.DeleteFromPool(deployedIDs)
			cleanup(true, deployedContainers, t)
			return nil, errors.New("Reserve Router Port Error: " + err.Error())
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
	return deployToHostsInZones(deps, manifest, sha, env, map[string][]string{"[any]": hosts}, []string{"[any]"}, t)
}

func copyContainer(auth *ManagerAuthArg, cid, toHost string, t *Task) (*Container, error) {
	// get old instance
	inst, err := datamodel.GetInstance(cid)
	if err != nil {
		return nil, err
	}

	// get manifest
	manifest := inst.Manifest
	if manifest == nil {
		// if we don't have the manifest in zk, try to get it from the supervisor
		ihReply, err := supervisor.Get(inst.Host, inst.ID)
		if err != nil {
			return nil, err
		}
		manifest = ihReply.Container.Manifest
	}
	manifest.Instances = 1

	// validate and get deps
	deps, err := validateDeploy(auth, manifest, inst.Sha, inst.Env, t)
	if err != nil {
		return nil, err
	}

	// get zone of toHost
	zone, err := supervisor.GetZone(toHost)
	if err != nil {
		return nil, err
	}

	deployed, err := deployToHostsInZones(deps, manifest, inst.Sha, inst.Env,
		map[string][]string{zone: []string{toHost}}, []string{zone}, t)
	if err != nil {
		return nil, err
	}
	// should only deploy 1 since we're only moving 1
	if len(deployed) != 1 {
		cleanup(true, deployed, t)
		return nil, errors.New(fmt.Sprintf("Didn't deploy 1 container. Deployed %d", len(deployed)))
	}
	return deployed[0], nil
}

func cleanup(removeContainerFromHost bool, deployedContainers []*Container, t *Task) {
	// kill all references to deployed containers as well as the container itself
	for _, container := range deployedContainers {
		supervisor.Teardown(container.Host, []string{container.ID}, false)
		if instance, err := datamodel.GetInstance(container.ID); err == nil {
			instance.Delete()
		} else {
			t.Log(fmt.Sprintf("Failed to clean up instance %s: %s", container.ID, err.Error()))
		}
		DeleteAppShaFromEnv(container.App, container.Sha, container.Env)
		if removeContainerFromHost {
			datamodel.Supervisor(container.Host).RemoveContainer(container.ID)
		}
	}
}

func cleanupZk(inst *datamodel.ZkInstance, t *Task) {
	// don't teardown from supervisor, this is meant as a pure zk cleanup
	inst.Delete()
	DeleteAppShaFromEnv(inst.App, inst.Sha, inst.Env)
}

//
// Teardown Stuff
//

func getContainerIDsOfShaEnv(t *Task, app, sha, env string) ([]string, error) {
	containerIDs, err := datamodel.ListInstances(app, sha, env)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error listing instances of %s @ %s in %s: %s", app, sha, env,
			err.Error()))
	}
	return containerIDs, nil
}

func getContainerIDsOfEnv(t *Task, app, env string) ([]string, error) {
	containerIDs := []string{}
	shas, err := datamodel.ListShas(app)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error listing shas of %s : %s", app, err.Error()))
	}
	for _, sha := range shas {
		containerIDs, err = datamodel.ListInstances(app, sha, env)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error listing instances of %s @ %s in %s: %s", app, sha, env,
				err.Error()))
		}
	}
	return containerIDs, nil
}

func getContainerIDsOfSha(t *Task, app, sha string) ([]string, error) {
	containerIDs := []string{}
	envs, err := datamodel.ListAppEnvs(app, sha)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error listing environments of %s @ %s: %s", app, sha, err.Error()))
	}
	for _, env := range envs {
		tmpContainerIDs, err := getContainerIDsOfShaEnv(t, app, sha, env)
		if err != nil {
			return nil, err
		}
		containerIDs = append(containerIDs, tmpContainerIDs...)
	}
	return containerIDs, nil
}

func getContainerIDsOfApp(t *Task, app string) ([]string, error) {
	containerIDs := []string{}
	shas, err := datamodel.ListShas(app)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error listing shas of %s: %s", app, err.Error()))
	}
	for _, sha := range shas {
		tmpContainerIDs, err := getContainerIDsOfSha(t, app, sha)
		if err != nil {
			return nil, err
		}
		containerIDs = append(containerIDs, tmpContainerIDs...)
	}
	return containerIDs, nil
}

func getContainerIDsToTeardown(t *Task, arg ManagerTeardownArg) (hostMap map[string][]string, err error) {
	hostMap = map[string][]string{} // map of host -> []string container ids
	// TODO(edanaher,2014-07-02): This pile of conditionals is braindead and caused us to ignore an environment
	// with no sha, tearing down all of of an app instead of just one environment.
	// We really need to fix this to be reasonable, but for the moment, to fix it, I'm just adding another case.
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
	} else if arg.ContainerID != "" {
		var instance *datamodel.ZkInstance
		instance, err = datamodel.GetInstance(arg.ContainerID)
		if err != nil {
			return
		}
		hostMap[instance.Host] = []string{arg.ContainerID}
		return
	} else if arg.App != "" {
		containerIDs := []string{}
		if arg.Sha != "" {
			if arg.Env != "" {
				if containerIDs, err = getContainerIDsOfShaEnv(t, arg.App, arg.Sha, arg.Env); err != nil {
					return nil, err
				}
			} else {
				if containerIDs, err = getContainerIDsOfSha(t, arg.App, arg.Sha); err != nil {
					return nil, err
				}
			}
		} else {
			if arg.Env != "" {
				if containerIDs, err = getContainerIDsOfEnv(t, arg.App, arg.Env); err != nil {
					return nil, err
				}
			} else {
				if containerIDs, err = getContainerIDsOfApp(t, arg.App); err != nil {
					return nil, err
				}
			}
		}
		for _, containerID := range containerIDs {
			instance, err := datamodel.GetInstance(containerID)
			if err != nil {
				continue
			}
			currentIDs := hostMap[instance.Host]
			if currentIDs == nil {
				hostMap[instance.Host] = []string{containerID}
			} else {
				hostMap[instance.Host] = append(currentIDs, containerID)
			}
		}
		return
	}
	return nil, errors.New("Invalid Arguments")
}
