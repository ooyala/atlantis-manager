package api

import (
	. "atlantis/manager/rpc/types"
	"fmt"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	arg := ManagerLoginArg{r.FormValue("User"), r.FormValue("Password"), r.FormValue("Secret")}
	var reply ManagerLoginReply
	err := manager.Login(arg, &reply)
	fmt.Fprintf(w, "%s", Output(map[string]interface{}{"User": r.FormValue("User"), "Secret": reply.Secret}, err))
}
