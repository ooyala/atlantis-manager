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

package supervisor

import (
	. "atlantis/supervisor/rpc/client"
	. "atlantis/supervisor/rpc/types"
)

var Port string

func Init(port string) {
	Port = port
}

func Deploy(host, app, sha, env, container string, man *Manifest) (*SupervisorDeployReply, error) {
	args := SupervisorDeployArg{Host: host, App: app, Sha: sha, Env: env, ContainerID: container, Manifest: man}
	var reply SupervisorDeployReply
	return &reply, NewSupervisorRPCClient(host+":"+Port).Call("Deploy", args, &reply)
}

func Teardown(host string, containerIDs []string, all bool) (*SupervisorTeardownReply, error) {
	args := SupervisorTeardownArg{containerIDs, all}
	var reply SupervisorTeardownReply
	return &reply, NewSupervisorRPCClient(host+":"+Port).Call("Teardown", args, &reply)
}

func HealthCheck(host string) (*SupervisorHealthCheckReply, error) {
	args := SupervisorHealthCheckArg{}
	var reply SupervisorHealthCheckReply
	return &reply, NewSupervisorRPCClient(host+":"+Port).Call("HealthCheck", args, &reply)
}

func GetZone(host string) (string, error) {
	hReply, err := HealthCheck(host)
	return hReply.Zone, err
}

func Get(host, containerID string) (*SupervisorGetReply, error) {
	args := SupervisorGetArg{containerID}
	var reply SupervisorGetReply
	return &reply, NewSupervisorRPCClient(host+":"+Port).Call("Get", args, &reply)
}

func AuthorizeSSH(host, containerID, user, publicKey string) (*SupervisorAuthorizeSSHReply, error) {
	args := SupervisorAuthorizeSSHArg{containerID, user, publicKey}
	var reply SupervisorAuthorizeSSHReply
	return &reply, NewSupervisorRPCClient(host+":"+Port).Call("AuthorizeSSH", args, &reply)
}

func DeauthorizeSSH(host, containerID, user string) (*SupervisorDeauthorizeSSHReply, error) {
	args := SupervisorDeauthorizeSSHArg{containerID, user}
	var reply SupervisorDeauthorizeSSHReply
	return &reply, NewSupervisorRPCClient(host+":"+Port).Call("DeauthorizeSSH", args, &reply)
}

func ContainerMaintenance(host, containerID string, maint bool) (*SupervisorContainerMaintenanceReply, error) {
	args := SupervisorContainerMaintenanceArg{containerID, maint}
	var reply SupervisorContainerMaintenanceReply
	return &reply, NewSupervisorRPCClient(host+":"+Port).Call("ContainerMaintenance", args, &reply)
}
