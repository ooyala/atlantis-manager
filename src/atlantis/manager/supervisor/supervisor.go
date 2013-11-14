package supervisor

import (
	. "atlantis/supervisor/rpc/client"
	. "atlantis/supervisor/rpc/types"
	"fmt"
)

var Port string

func Init(port uint16) {
	Port = fmt.Sprintf("%d", port)
}

func Deploy(host, app, sha, env, container string, man *Manifest) (*SupervisorDeployReply, error) {
	args := SupervisorDeployArg{Host: host, App: app, Sha: sha, Env: env, ContainerId: container, Manifest: man}
	var reply SupervisorDeployReply
	return &reply, NewSupervisorRPCClient(host+":"+Port).Call("Deploy", args, &reply)
}

func Teardown(host string, containerIds []string, all bool) (*SupervisorTeardownReply, error) {
	args := SupervisorTeardownArg{containerIds, all}
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

func Get(host, containerId string) (*SupervisorGetReply, error) {
	args := SupervisorGetArg{containerId}
	var reply SupervisorGetReply
	return &reply, NewSupervisorRPCClient(host+":"+Port).Call("Get", args, &reply)
}

func AuthorizeSSH(host, containerId, user, publicKey string) (*SupervisorAuthorizeSSHReply, error) {
	args := SupervisorAuthorizeSSHArg{containerId, user, publicKey}
	var reply SupervisorAuthorizeSSHReply
	return &reply, NewSupervisorRPCClient(host+":"+Port).Call("AuthorizeSSH", args, &reply)
}

func DeauthorizeSSH(host, containerId, user string) (*SupervisorDeauthorizeSSHReply, error) {
	args := SupervisorDeauthorizeSSHArg{containerId, user}
	var reply SupervisorDeauthorizeSSHReply
	return &reply, NewSupervisorRPCClient(host+":"+Port).Call("DeauthorizeSSH", args, &reply)
}

func ContainerMaintenance(host, containerId string, maint bool) (*SupervisorContainerMaintenanceReply, error) {
	args := SupervisorContainerMaintenanceArg{containerId, maint}
	var reply SupervisorContainerMaintenanceReply
	return &reply, NewSupervisorRPCClient(host+":"+Port).Call("ContainerMaintenance", args, &reply)
}
