package api

import (
	. "atlantis/manager/rpc/types"
	"fmt"
	"net/http"
)

func Health(w http.ResponseWriter, r *http.Request) {
	var reply ManagerHealthCheckReply
	arg := ManagerHealthCheckArg{}
	err := manager.HealthCheck(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Region": reply.Region, "Status": reply.Status}, err))
}
