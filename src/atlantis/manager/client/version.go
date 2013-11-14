package client

import (
	. "atlantis/common"
	. "atlantis/manager/constant"
)

type VersionCommand struct {
}

func (c *VersionCommand) Execute(args []string) error {
	InitNoLogin()
	Log("Manager Version Check...")
	arg := VersionArg{}
	var reply VersionReply
	defer func() {
		if err := recover(); err != nil {
			reply.RPCVersion = "unknown"
			reply.APIVersion = "unknown"
		}
		Log("-> client rpc: %s", ManagerRPCVersion)
		Log("-> server rpc: %s", reply.RPCVersion)
		Log("-> server api: %s", reply.APIVersion)
		Output(map[string]interface{}{"client": map[string]string{"rpc": ManagerRPCVersion},
			"server": map[string]string{"rpc": reply.RPCVersion, "api": reply.APIVersion}},
			map[string]string{"client rpc": ManagerRPCVersion, "server rpc": reply.RPCVersion,
				"server api": reply.APIVersion}, nil)
	}()
	err := rpcClient.Call("Version", arg, &reply)
	return Output(nil, nil, err)
}
