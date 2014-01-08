package rpc

import (
	. "atlantis/common"
	. "atlantis/manager/constant"
	"atlantis/manager/datamodel"
	"atlantis/manager/manager"
	"atlantis/manager/router"
	. "atlantis/manager/rpc/types"
	"atlantis/manager/supervisor"
	"errors"
	"fmt"
)

// ----------------------------------------------------------------------------------------------------------
// Register Router
// ----------------------------------------------------------------------------------------------------------

type RegisterRouterExecutor struct {
	arg   ManagerRegisterRouterArg
	reply *ManagerRegisterRouterReply
}

func (e *RegisterRouterExecutor) Request() interface{} {
	return e.arg
}

func (e *RegisterRouterExecutor) Result() interface{} {
	return e.reply
}

func (e *RegisterRouterExecutor) Description() string {
	return fmt.Sprintf("%s in %s, internal: %t", e.arg.Host, e.arg.Zone, e.arg.Internal)
}

func (e *RegisterRouterExecutor) Authorize() error {
	return AuthorizeSuperUser(&e.arg.ManagerAuthArg)
}

func (e *RegisterRouterExecutor) Execute(t *Task) error {
	if e.arg.Host == "" {
		return errors.New("Please specify a host to register")
	}
	if e.arg.Zone == "" {
		return errors.New("Please specify a zone")
	}
	routerObj, err := router.Register(e.arg.Internal, e.arg.Zone, e.arg.Host)
	if err != nil {
		e.reply.Status = StatusError
	}
	castedRouter := Router(*routerObj)
	e.reply.Router = &castedRouter
	e.reply.Status = StatusOk
	return err
}

func (m *ManagerRPC) RegisterRouterResult(id string, result *ManagerRegisterRouterReply) error {
	if id == "" {
		return errors.New("ID empty")
	}
	status, err := Tracker.Status(id)
	if status.Status == StatusUnknown {
		return errors.New("Unknown ID.")
	}
	if status.Name != "RegisterRouter" {
		return errors.New("ID is not a RegisterRouter.")
	}
	if !status.Done {
		return errors.New("RegisterRouter isn't done.")
	}
	if status.Status == StatusError || err != nil {
		return err
	}
	getResult := Tracker.Result(id)
	switch r := getResult.(type) {
	case *ManagerRegisterRouterReply:
		*result = *r
	default:
		// this should never happen
		return errors.New("Invalid Result Type.")
	}
	return nil
}

type UnregisterRouterExecutor struct {
	arg   ManagerRegisterRouterArg
	reply *ManagerRegisterRouterReply
}

func (e *UnregisterRouterExecutor) Request() interface{} {
	return e.arg
}

func (e *UnregisterRouterExecutor) Result() interface{} {
	return e.reply
}

func (e *UnregisterRouterExecutor) Description() string {
	return fmt.Sprintf("%s in %s, internal: %t", e.arg.Host, e.arg.Zone, e.arg.Internal)
}

func (e *UnregisterRouterExecutor) Authorize() error {
	return AuthorizeSuperUser(&e.arg.ManagerAuthArg)
}

func (e *UnregisterRouterExecutor) Execute(t *Task) error {
	if e.arg.Host == "" {
		return errors.New("Please specify a host to uregister")
	}
	if e.arg.Zone == "" {
		return errors.New("Please specify a zone")
	}
	err := router.Unregister(e.arg.Internal, e.arg.Zone, e.arg.Host)
	if err != nil {
		e.reply.Status = StatusError
	}
	e.reply.Status = StatusOk
	return err
}

func (m *ManagerRPC) UnregisterRouterResult(id string, result *ManagerRegisterRouterReply) error {
	if id == "" {
		return errors.New("ID empty")
	}
	status, err := Tracker.Status(id)
	if status.Status == StatusUnknown {
		return errors.New("Unknown ID.")
	}
	if status.Name != "UnregisterRouter" {
		return errors.New("ID is not a UnregisterRouter.")
	}
	if !status.Done {
		return errors.New("UnregisterRouter isn't done.")
	}
	if status.Status == StatusError || err != nil {
		return err
	}
	getResult := Tracker.Result(id)
	switch r := getResult.(type) {
	case *ManagerRegisterRouterReply:
		*result = *r
	default:
		// this should never happen
		return errors.New("Invalid Result Type.")
	}
	return nil
}

type ListRoutersExecutor struct {
	arg   ManagerListRoutersArg
	reply *ManagerListRoutersReply
}

func (e *ListRoutersExecutor) Request() interface{} {
	return e.arg
}

func (e *ListRoutersExecutor) Result() interface{} {
	return e.reply
}

func (e *ListRoutersExecutor) Description() string {
	return "ListRouters"
}

func (e *ListRoutersExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

func (e *ListRoutersExecutor) Execute(t *Task) (err error) {
	e.reply.Routers, err = datamodel.ListRouters(e.arg.Internal)
	if err != nil {
		e.reply.Status = StatusError
	}
	e.reply.Status = StatusOk
	return err
}

type GetRouterExecutor struct {
	arg   ManagerGetRouterArg
	reply *ManagerGetRouterReply
}

func (e *GetRouterExecutor) Request() interface{} {
	return e.arg
}

func (e *GetRouterExecutor) Result() interface{} {
	return e.reply
}

func (e *GetRouterExecutor) Description() string {
	return fmt.Sprintf("%s in %s", e.arg.Host, e.arg.Zone)
}

func (e *GetRouterExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

func (e *GetRouterExecutor) Execute(t *Task) error {
	zkRouter, err := datamodel.GetRouter(e.arg.Internal, e.arg.Zone, e.arg.Host)
	castedRouter := Router(*zkRouter)
	e.reply.Router = &castedRouter
	if err != nil {
		e.reply.Status = StatusError
	}
	e.reply.Status = StatusOk
	return err
}

// ----------------------------------------------------------------------------------------------------------
// Register App
// ----------------------------------------------------------------------------------------------------------

type RegisterAppExecutor struct {
	arg   ManagerRegisterAppArg
	reply *ManagerRegisterAppReply
}

func (e *RegisterAppExecutor) Request() interface{} {
	return e.arg
}

func (e *RegisterAppExecutor) Result() interface{} {
	return e.reply
}

func (e *RegisterAppExecutor) Description() string {
	return fmt.Sprintf("%s -> %s:%s, non-atlantis: %t", e.arg.Name, e.arg.Repo, e.arg.Root, e.arg.NonAtlantis)
}

func (e *RegisterAppExecutor) Authorize() error {
	return AuthorizeApp(&e.arg.ManagerAuthArg, e.arg.Name)
}

func (e *RegisterAppExecutor) Execute(t *Task) error {
	if e.arg.Name == "" {
		return errors.New("Please specify an app name to register")
	}
	if !AppRegexp.MatchString(e.arg.Name) {
		return errors.New("App name must be [A-Za-z0-9-]+")
	}
	if !e.arg.NonAtlantis && e.arg.Repo == "" {
		return errors.New("Please specify a repo")
	}
	if !e.arg.NonAtlantis && e.arg.Root == "" {
		return errors.New("Please specify the repo's root")
	}
	if e.arg.Email == "" {
		return errors.New("Please specify the email of the app owner")
	}
	typ := "http"
	if e.arg.Type != "" {
		typ = e.arg.Type
	}
	if _, err := datamodel.GetApp(e.arg.Name); err == nil {
		return errors.New("Already Registered.")
	}
	_, err := datamodel.CreateOrUpdateApp(e.arg.NonAtlantis, typ, e.arg.Name, e.arg.Repo, e.arg.Root,
		e.arg.Email, e.arg.Addrs)
	if err != nil {
		e.reply.Status = StatusError
	}
	e.reply.Status = StatusOk
	return err
}

type UpdateAppExecutor struct {
	arg   ManagerRegisterAppArg
	reply *ManagerRegisterAppReply
}

func (e *UpdateAppExecutor) Request() interface{} {
	return e.arg
}

func (e *UpdateAppExecutor) Result() interface{} {
	return e.reply
}

func (e *UpdateAppExecutor) Description() string {
	return fmt.Sprintf("%s -> %s:%s, non-atlantis: %t", e.arg.Name, e.arg.Repo, e.arg.Root, e.arg.NonAtlantis)
}

func (e *UpdateAppExecutor) Authorize() error {
	return AuthorizeApp(&e.arg.ManagerAuthArg, e.arg.Name)
}

func (e *UpdateAppExecutor) Execute(t *Task) error {
	if e.arg.Name == "" {
		return errors.New("Please specify an app name to update")
	}
	if !AppRegexp.MatchString(e.arg.Name) {
		return errors.New("App name must be [A-Za-z0-9-]+")
	}
	if !e.arg.NonAtlantis && e.arg.Repo == "" {
		return errors.New("Please specify a repo")
	}
	if !e.arg.NonAtlantis && e.arg.Root == "" {
		return errors.New("Please specify the repo's root")
	}
	if e.arg.Email == "" {
		return errors.New("Please specify the email of the app owner")
	}
	typ := "http"
	if e.arg.Type != "" {
		typ = e.arg.Type
	}
	_, err := datamodel.CreateOrUpdateApp(e.arg.NonAtlantis, typ, e.arg.Name, e.arg.Repo, e.arg.Root,
		e.arg.Email, e.arg.Addrs)
	if err != nil {
		e.reply.Status = StatusError
	}
	e.reply.Status = StatusOk
	return err
}

type UnregisterAppExecutor struct {
	arg   ManagerRegisterAppArg
	reply *ManagerRegisterAppReply
}

func (e *UnregisterAppExecutor) Request() interface{} {
	return e.arg
}

func (e *UnregisterAppExecutor) Result() interface{} {
	return e.reply
}

func (e *UnregisterAppExecutor) Description() string {
	return fmt.Sprintf("%s", e.arg.Name)
}

func (e *UnregisterAppExecutor) Authorize() error {
	return AuthorizeApp(&e.arg.ManagerAuthArg, e.arg.Name)
}

func (e *UnregisterAppExecutor) Execute(t *Task) error {
	if e.arg.Name == "" {
		return errors.New("Please specify an app name to unregister")
	}
	app, err := datamodel.GetApp(e.arg.Name)
	if err != nil || app == nil {
		e.reply.Status = StatusError
		return errors.New("App " + e.arg.Name + " does not exist")
	}
	if err = app.Delete(); err != nil {
		e.reply.Status = StatusError
		return err
	}
	e.reply.Status = StatusOk
	return nil
}

type GetAppExecutor struct {
	arg   ManagerGetAppArg
	reply *ManagerGetAppReply
}

func (e *GetAppExecutor) Request() interface{} {
	return e.arg
}

func (e *GetAppExecutor) Result() interface{} {
	return e.reply
}

func (e *GetAppExecutor) Description() string {
	return fmt.Sprintf("%s", e.arg.Name)
}

func (e *GetAppExecutor) Authorize() error {
	return AuthorizeApp(&e.arg.ManagerAuthArg, e.arg.Name)
}

func (e *GetAppExecutor) Execute(t *Task) error {
	if e.arg.Name == "" {
		return errors.New("Please specify an app name to get")
	}
	app, err := datamodel.GetApp(e.arg.Name)
	if err != nil || app == nil {
		e.reply.Status = StatusError
		return errors.New("App " + e.arg.Name + " does not exist")
	}
	e.reply.Status = StatusOk
	castedApp := App(*app)
	e.reply.App = &castedApp
	return nil
}

type ListRegisteredAppsExecutor struct {
	arg   ManagerListRegisteredAppsArg
	reply *ManagerListRegisteredAppsReply
}

func (e *ListRegisteredAppsExecutor) Request() interface{} {
	return e.arg
}

func (e *ListRegisteredAppsExecutor) Result() interface{} {
	return e.reply
}

func (e *ListRegisteredAppsExecutor) Description() string {
	return "ListRegisteredApps"
}

func (e *ListRegisteredAppsExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

func (e *ListRegisteredAppsExecutor) Execute(t *Task) (err error) {
	e.reply.Apps, err = datamodel.ListRegisteredApps()
	if err != nil {
		e.reply.Status = StatusError
	} else {
		e.reply.Status = StatusOk
	}
	return err
}

type ListAuthorizedRegisteredAppsExecutor struct {
	arg   ManagerListRegisteredAppsArg
	reply *ManagerListRegisteredAppsReply
}

func (e *ListAuthorizedRegisteredAppsExecutor) Request() interface{} {
	return e.arg
}

func (e *ListAuthorizedRegisteredAppsExecutor) Result() interface{} {
	return e.reply
}

func (e *ListAuthorizedRegisteredAppsExecutor) Description() string {
	return "ListAuthorizedRegisteredApps"
}

func (e *ListAuthorizedRegisteredAppsExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

func (e *ListAuthorizedRegisteredAppsExecutor) Execute(t *Task) (err error) {
	apps, err := datamodel.ListRegisteredApps()
	if err != nil {
		e.reply.Status = StatusError
	} else {
		e.reply.Status = StatusOk
	}
	authdApps := []string{}
	for _, app := range apps {
		if err := AuthorizeApp(&e.arg.ManagerAuthArg, app); err == nil {
			authdApps = append(authdApps, app)
		}
	}
	e.reply.Apps = authdApps
	return err
}

// ----------------------------------------------------------------------------------------------------------
// Register Supervisor
// ----------------------------------------------------------------------------------------------------------

type RegisterSupervisorExecutor struct {
	arg   ManagerRegisterSupervisorArg
	reply *ManagerRegisterSupervisorReply
}

func (e *RegisterSupervisorExecutor) Request() interface{} {
	return e.arg
}

func (e *RegisterSupervisorExecutor) Result() interface{} {
	return e.reply
}

func (e *RegisterSupervisorExecutor) Description() string {
	return fmt.Sprintf("%s:%s", e.arg.Host, supervisor.Port)
}

func (e *RegisterSupervisorExecutor) Execute(t *Task) error {
	if e.arg.Host == "" {
		return errors.New("Please specify a host to register")
	}
	// check health of to be registered supervisor
	health, err := supervisor.HealthCheck(e.arg.Host)
	if err != nil {
		e.reply.Status = StatusError
		return err
	}
	if health.Status != StatusOk && health.Status != StatusFull {
		e.reply.Status = health.Status
		return errors.New("Status is " + health.Status)
	}
	if health.Region != Region {
		e.reply.Status = "Region Mismatch"
		return errors.New("Supervisor Region (" + health.Region + ") does not match Manager Region (" + Region + ")")
	}
	err = datamodel.Supervisor(e.arg.Host).Touch()
	if err != nil {
		e.reply.Status = StatusError
		return err
	}
	if err == nil {
		e.reply.Status = StatusOk
	} else {
		e.reply.Status = StatusError
	}
	return err
}

func (e *RegisterSupervisorExecutor) Authorize() error {
	return AuthorizeSuperUser(&e.arg.ManagerAuthArg)
}

func (m *ManagerRPC) RegisterSupervisorResult(id string, result *ManagerRegisterSupervisorReply) error {
	if id == "" {
		return errors.New("ID empty")
	}
	status, err := Tracker.Status(id)
	if status.Status == StatusUnknown {
		return errors.New("Unknown ID.")
	}
	if status.Name != "RegisterSupervisor" {
		return errors.New("ID is not a RegisterSupervisor.")
	}
	if !status.Done {
		return errors.New("RegisterSupervisor isn't done.")
	}
	if status.Status == StatusError || err != nil {
		return err
	}
	getResult := Tracker.Result(id)
	switch r := getResult.(type) {
	case *ManagerRegisterSupervisorReply:
		*result = *r
	default:
		// this should never happen
		return errors.New("Invalid Result Type.")
	}
	return nil
}

type UnregisterSupervisorExecutor struct {
	arg   ManagerRegisterSupervisorArg
	reply *ManagerRegisterSupervisorReply
}

func (e *UnregisterSupervisorExecutor) Request() interface{} {
	return e.arg
}

func (e *UnregisterSupervisorExecutor) Result() interface{} {
	return e.reply
}

func (e *UnregisterSupervisorExecutor) Description() string {
	return fmt.Sprintf("%s:%s", e.arg.Host, supervisor.Port)
}

func (e *UnregisterSupervisorExecutor) Execute(t *Task) error {
	if e.arg.Host == "" {
		return errors.New("Please specify a host to unregister")
	}
	// teardown containers
	supervisor.Teardown(e.arg.Host, []string{}, true)
	err := datamodel.Supervisor(e.arg.Host).Delete()
	if err != nil {
		return err
	}
	e.reply.Status = StatusOk
	return nil
}

func (e *UnregisterSupervisorExecutor) Authorize() error {
	return AuthorizeSuperUser(&e.arg.ManagerAuthArg)
}

func (m *ManagerRPC) UnregisterSupervisorResult(id string, result *ManagerRegisterSupervisorReply) error {
	if id == "" {
		return errors.New("ID empty")
	}
	status, err := Tracker.Status(id)
	if status.Status == StatusUnknown {
		return errors.New("Unknown ID.")
	}
	if status.Name != "UnregisterSupervisor" {
		return errors.New("ID is not a UnregisterSupervisor.")
	}
	if !status.Done {
		return errors.New("UnregisterSupervisor isn't done.")
	}
	if status.Status == StatusError || err != nil {
		return err
	}
	getResult := Tracker.Result(id)
	switch r := getResult.(type) {
	case *ManagerRegisterSupervisorReply:
		*result = *r
	default:
		// this should never happen
		return errors.New("Invalid Result Type.")
	}
	return nil
}

type ListSupervisorsExecutor struct {
	arg   ManagerListSupervisorsArg
	reply *ManagerListSupervisorsReply
}

func (e *ListSupervisorsExecutor) Request() interface{} {
	return e.arg
}

func (e *ListSupervisorsExecutor) Result() interface{} {
	return e.reply
}

func (e *ListSupervisorsExecutor) Description() string {
	return "ListSupervisors"
}

func (e *ListSupervisorsExecutor) Execute(t *Task) (err error) {
	e.reply.Supervisors, err = datamodel.ListSupervisors()
	if err != nil {
		e.reply.Status = StatusError
	} else {
		e.reply.Status = StatusOk
	}
	return err
}

func (e *ListSupervisorsExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

// ----------------------------------------------------------------------------------------------------------
// Register Manager
// ----------------------------------------------------------------------------------------------------------

type RegisterManagerExecutor struct {
	arg   ManagerRegisterManagerArg
	reply *ManagerRegisterManagerReply
}

func (e *RegisterManagerExecutor) Request() interface{} {
	return e.arg
}

func (e *RegisterManagerExecutor) Result() interface{} {
	return e.reply
}

func (e *RegisterManagerExecutor) Description() string {
	return fmt.Sprintf("%s:%s in %s", e.arg.Host, lPort, e.arg.Region)
}

func (e *RegisterManagerExecutor) Execute(t *Task) error {
	if e.arg.Host == "" {
		return errors.New("Please specify an Host to register")
	}
	if e.arg.Region == "" {
		return errors.New("Please specify a Region to register")
	}
	mgr, err := manager.Register(e.arg.Region, e.arg.Host, e.arg.RegistryCName, e.arg.ManagerCName)
	castedManager := Manager(*mgr)
	e.reply.Manager = &castedManager
	if err != nil {
		e.reply.Status = StatusError
	} else {
		e.reply.Status = StatusOk
	}
	return err
}

func (e *RegisterManagerExecutor) Authorize() error {
	return AuthorizeSuperUser(&e.arg.ManagerAuthArg)
}

func (m *ManagerRPC) RegisterManagerResult(id string, result *ManagerRegisterManagerReply) error {
	if id == "" {
		return errors.New("ID empty")
	}
	status, err := Tracker.Status(id)
	if status.Status == StatusUnknown {
		return errors.New("Unknown ID.")
	}
	if status.Name != "RegisterManager" {
		return errors.New("ID is not a RegisterManager.")
	}
	if !status.Done {
		return errors.New("RegisterManager isn't done.")
	}
	if status.Status == StatusError || err != nil {
		return err
	}
	getResult := Tracker.Result(id)
	switch r := getResult.(type) {
	case *ManagerRegisterManagerReply:
		*result = *r
	default:
		// this should never happen
		return errors.New("Invalid Result Type.")
	}
	return nil
}

type UnregisterManagerExecutor struct {
	arg   ManagerRegisterManagerArg
	reply *ManagerRegisterManagerReply
}

func (e *UnregisterManagerExecutor) Request() interface{} {
	return e.arg
}

func (e *UnregisterManagerExecutor) Result() interface{} {
	return e.reply
}

func (e *UnregisterManagerExecutor) Description() string {
	return fmt.Sprintf("%s:%s in %s", e.arg.Host, lPort, e.arg.Region)
}

func (e *UnregisterManagerExecutor) Execute(t *Task) error {
	if e.arg.Host == "" {
		return errors.New("Please specify an host to unregister")
	}
	if e.arg.Region == "" {
		return errors.New("Please specify a region to unregister")
	}
	err := manager.Unregister(e.arg.Region, e.arg.Host)
	if err != nil {
		e.reply.Status = StatusError
	} else {
		e.reply.Status = StatusOk
	}
	return err
}

func (e *UnregisterManagerExecutor) Authorize() error {
	return AuthorizeSuperUser(&e.arg.ManagerAuthArg)
}

func (m *ManagerRPC) UnregisterManagerResult(id string, result *ManagerRegisterManagerReply) error {
	if id == "" {
		return errors.New("ID empty")
	}
	status, err := Tracker.Status(id)
	if status.Status == StatusUnknown {
		return errors.New("Unknown ID.")
	}
	if status.Name != "UnregisterManager" {
		return errors.New("ID is not a UnregisterManager.")
	}
	if !status.Done {
		return errors.New("UnregisterManager isn't done.")
	}
	if status.Status == StatusError || err != nil {
		return err
	}
	getResult := Tracker.Result(id)
	switch r := getResult.(type) {
	case *ManagerRegisterManagerReply:
		*result = *r
	default:
		// this should never happen
		return errors.New("Invalid Result Type.")
	}
	return nil
}

type ListManagersExecutor struct {
	arg   ManagerListManagersArg
	reply *ManagerListManagersReply
}

func (e *ListManagersExecutor) Request() interface{} {
	return e.arg
}

func (e *ListManagersExecutor) Result() interface{} {
	return e.reply
}

func (e *ListManagersExecutor) Description() string {
	return "ListManagers"
}

func (e *ListManagersExecutor) Execute(t *Task) (err error) {
	e.reply.Managers, err = datamodel.ListManagers()
	if err != nil {
		e.reply.Status = StatusError
	} else {
		e.reply.Status = StatusOk
	}
	return
}

func (e *ListManagersExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

type GetManagerExecutor struct {
	arg   ManagerGetManagerArg
	reply *ManagerGetManagerReply
}

func (e *GetManagerExecutor) Request() interface{} {
	return e.arg
}

func (e *GetManagerExecutor) Result() interface{} {
	return e.reply
}

func (e *GetManagerExecutor) Description() string {
	return fmt.Sprintf("%s in %s", e.arg.Host, e.arg.Region)
}

func (e *GetManagerExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

func (e *GetManagerExecutor) Execute(t *Task) error {
	zkManager, err := datamodel.GetManager(e.arg.Region, e.arg.Host)
	castedManager := Manager(*zkManager)
	e.reply.Manager = &castedManager
	if err != nil {
		e.reply.Status = StatusError
	}
	e.reply.Status = StatusOk
	return err
}

type GetSelfExecutor struct {
	arg   ManagerGetSelfArg
	reply *ManagerGetManagerReply
}

func (e *GetSelfExecutor) Request() interface{} {
	return e.arg
}

func (e *GetSelfExecutor) Result() interface{} {
	return e.reply
}

func (e *GetSelfExecutor) Description() string {
	return fmt.Sprintf("%s in %s", Host, Region)
}

func (e *GetSelfExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

func (e *GetSelfExecutor) Execute(t *Task) error {
	zkManager, err := datamodel.GetManager(Host, Region)
	castedManager := Manager(*zkManager)
	e.reply.Manager = &castedManager
	if err != nil {
		e.reply.Status = StatusError
	}
	e.reply.Status = StatusOk
	return err
}

func (m *ManagerRPC) RegisterRouter(arg ManagerRegisterRouterArg, reply *AsyncReply) error {
	return NewTask("RegisterRouter", &RegisterRouterExecutor{arg, &ManagerRegisterRouterReply{}}).RunAsync(reply)
}

func (m *ManagerRPC) UnregisterRouter(arg ManagerRegisterRouterArg, reply *AsyncReply) error {
	return NewTask("UnregisterRouter", &UnregisterRouterExecutor{arg, &ManagerRegisterRouterReply{}}).RunAsync(reply)
}

func (m *ManagerRPC) GetRouter(arg ManagerGetRouterArg, reply *ManagerGetRouterReply) error {
	return NewTask("GetRouter", &GetRouterExecutor{arg, reply}).Run()
}

func (m *ManagerRPC) ListRouters(arg ManagerListRoutersArg, reply *ManagerListRoutersReply) error {
	return NewTask("ListRouters", &ListRoutersExecutor{arg, reply}).Run()
}

func (m *ManagerRPC) RegisterApp(arg ManagerRegisterAppArg, reply *ManagerRegisterAppReply) error {
	return NewTask("RegisterApp", &RegisterAppExecutor{arg, reply}).Run()
}

func (m *ManagerRPC) UpdateApp(arg ManagerRegisterAppArg, reply *ManagerRegisterAppReply) error {
	return NewTask("UpdateApp", &UpdateAppExecutor{arg, reply}).Run()
}

func (m *ManagerRPC) UnregisterApp(arg ManagerRegisterAppArg, reply *ManagerRegisterAppReply) error {
	return NewTask("UnregisterApp", &UnregisterAppExecutor{arg, reply}).Run()
}

func (m *ManagerRPC) GetApp(arg ManagerGetAppArg, reply *ManagerGetAppReply) error {
	return NewTask("GetApp", &GetAppExecutor{arg, reply}).Run()
}

func (m *ManagerRPC) ListRegisteredApps(arg ManagerListRegisteredAppsArg, reply *ManagerListRegisteredAppsReply) error {
	return NewTask("ListRegisteredApps", &ListRegisteredAppsExecutor{arg, reply}).Run()
}

func (m *ManagerRPC) ListAuthorizedRegisteredApps(arg ManagerListRegisteredAppsArg, reply *ManagerListRegisteredAppsReply) error {
	return NewTask("ListAuthorizedRegisteredApps", &ListRegisteredAppsExecutor{arg, reply}).Run()
}

func (m *ManagerRPC) RegisterSupervisor(arg ManagerRegisterSupervisorArg, reply *AsyncReply) error {
	return NewTask("RegisterSupervisor", &RegisterSupervisorExecutor{arg, &ManagerRegisterSupervisorReply{}}).RunAsync(reply)
}

func (m *ManagerRPC) UnregisterSupervisor(arg ManagerRegisterSupervisorArg, reply *AsyncReply) error {
	return NewTask("UnregisterSupervisor", &UnregisterSupervisorExecutor{arg, &ManagerRegisterSupervisorReply{}}).RunAsync(reply)
}

func (m *ManagerRPC) ListSupervisors(arg ManagerListSupervisorsArg, reply *ManagerListSupervisorsReply) error {
	return NewTask("ListSupervisors", &ListSupervisorsExecutor{arg, reply}).Run()
}

func (m *ManagerRPC) RegisterManager(arg ManagerRegisterManagerArg, reply *AsyncReply) error {
	return NewTask("RegisterManager", &RegisterManagerExecutor{arg, &ManagerRegisterManagerReply{}}).RunAsync(reply)
}

func (m *ManagerRPC) UnregisterManager(arg ManagerRegisterManagerArg, reply *AsyncReply) error {
	return NewTask("UnregisterManager", &UnregisterManagerExecutor{arg, &ManagerRegisterManagerReply{}}).RunAsync(reply)
}

func (m *ManagerRPC) ListManagers(arg ManagerListManagersArg, reply *ManagerListManagersReply) error {
	return NewTask("ListManagers", &ListManagersExecutor{arg, reply}).Run()
}

func (m *ManagerRPC) GetManager(arg ManagerGetManagerArg, reply *ManagerGetManagerReply) error {
	return NewTask("GetManager", &GetManagerExecutor{arg, reply}).Run()
}

func (m *ManagerRPC) GetSelf(arg ManagerGetSelfArg, reply *ManagerGetManagerReply) error {
	return NewTask("GetSelf", &GetSelfExecutor{arg, reply}).Run()
}
