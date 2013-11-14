package rpc

import (
	. "atlantis/common"
	"atlantis/manager/ldap"
	. "atlantis/manager/rpc/types"
	"errors"
)

type Authorizer struct {
	User     string
	Password string
	Secret   string
}

func (a *Authorizer) Authenticate() (err error) {
	a.Secret, err = ldap.Login(a.User, a.Password, a.Secret)
	return err
}

func SimpleAuthorize(AuthArg *ManagerAuthArg) error {
	user, password, secret := AuthArg.Credentials()
	auther := Authorizer{user, password, secret}
	if err := auther.Authenticate(); err != nil {
		return err
	}
	return nil
}

func AuthorizeTeamAdmin(AuthArg *ManagerAuthArg, team string) error {
	if err := SimpleAuthorize(AuthArg); err != nil {
		return err
	}
	req := ManagerTeamAdminArg{*AuthArg, team}
	var res ManagerTeamAdminReply
	err := NewTask("Authorizer-IsTeamAdmin", &IsTeamAdminExecutor{req, &res}).Run()
	if err != nil || !res.IsAdmin {
		return errors.New("Not a Team Admin")
	}
	return nil
}

func AuthorizeApp(AuthArg *ManagerAuthArg, app string) error {
	if err := SimpleAuthorize(AuthArg); err != nil {
		return err
	}
	var reply ManagerAppPermissionsReply
	arg := ManagerAppPermissionsArg{*AuthArg, app}
	err := NewTask("Authorizer-HasAppPermissions", &HasAppPermissionsExecutor{arg, &reply}).Run()
	if err != nil || !reply.Permission {
		return errors.New("Not Authorized to Deploy App")
	}
	return nil
}

func AuthorizeSuperUser(AuthArg *ManagerAuthArg) error {
	if err := SimpleAuthorize(AuthArg); err != nil {
		return err
	}
	var reply ManagerSuperUserReply
	arg := ManagerSuperUserArg{*AuthArg}
	err := NewTask("Authorizer-IsSuperUser", &IsSuperUserExecutor{arg, &reply}).Run()
	if err != nil || !reply.IsSuperUser {
		return errors.New("Not a Super User")
	}
	return nil
}
