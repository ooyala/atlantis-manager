package client

import (
	atlantis "atlantis/common"
	. "atlantis/manager/constant"
)

func NewManagerRPCClient(hostAndPort string) *atlantis.RPCClient {
	return atlantis.NewRPCClient(hostAndPort, "ManagerRPC", ManagerRPCVersion, true)
}

func NewManagerRPCClientWithConfig(cfg atlantis.RPCServerOpts) *atlantis.RPCClient {
	return atlantis.NewRPCClientWithConfig(cfg, "ManagerRPC", ManagerRPCVersion, true)
}
