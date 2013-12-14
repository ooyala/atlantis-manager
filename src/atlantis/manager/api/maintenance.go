package api

import (
	. "atlantis/manager/rpc/types"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func ContainerMaintenance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := ManagerAuthArg{r.FormValue("User"), "", r.FormValue("Secret")}
	maintenance, err := strconv.ParseBool(r.FormValue("Maintenance"))
	if err != nil {
		fmt.Fprintf(w, "%s", Output(map[string]interface{}{}, err))
	}
	arg := ManagerContainerMaintenanceArg{auth, vars["ID"], maintenance}
	var reply ManagerContainerMaintenanceReply
	err = manager.ContainerMaintenance(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"Status": reply.Status}, err))
}
