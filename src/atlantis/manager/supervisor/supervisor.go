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
