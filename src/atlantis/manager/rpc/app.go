package rpc

import (
	. "atlantis/common"
	"atlantis/manager/datamodel"
	. "atlantis/manager/rpc/types"
	"errors"
	"fmt"
)

// ----------------------------------------------------------------------------------------------------------
// Depender App Data Methods
// ----------------------------------------------------------------------------------------------------------

type AddDependerAppDataExecutor struct {
	arg   ManagerAddDependerAppDataArg
	reply *ManagerAddDependerAppDataReply
}

func (e *AddDependerAppDataExecutor) Request() interface{} {
	return e.arg
}

func (e *AddDependerAppDataExecutor) Result() interface{} {
	return e.reply
}

func (e *AddDependerAppDataExecutor) Description() string {
	return fmt.Sprintf("[+] %s depender %+v", e.arg.App, e.arg.DependerAppData)
}

func (e *AddDependerAppDataExecutor) Authorize() error {
	return AuthorizeApp(&e.arg.ManagerAuthArg, e.arg.App)
}

func (e *AddDependerAppDataExecutor) Execute(t *Task) error {
	if e.arg.App == "" {
		return errors.New("Please specify an app")
	}
	if e.arg.DependerAppData == nil {
		return errors.New("Please specify data for the depender app")
	} else if e.arg.DependerAppData.Name == "" {
		return errors.New("Please specify name for the depender app")
	}
	zkApp, err := datamodel.GetApp(e.arg.App)
	if err != nil {
		e.reply.Status = StatusError
		return err
	}
	// go through all the env data and encrypt stuff
	err = zkApp.AddDependerAppData(e.arg.DependerAppData)
	if err != nil {
		e.reply.Status = StatusError
		return err
	}
	e.reply.Status = StatusOk
	castedApp := App(*zkApp)
	e.reply.App = &castedApp
	return err
}

func (m *ManagerRPC) AddDependerAppData(arg ManagerAddDependerAppDataArg, reply *ManagerAddDependerAppDataReply) error {
	return NewTask("AddDependerAppData", &AddDependerAppDataExecutor{arg, reply}).Run()
}

type RemoveDependerAppDataExecutor struct {
	arg   ManagerRemoveDependerAppDataArg
	reply *ManagerRemoveDependerAppDataReply
}

func (e *RemoveDependerAppDataExecutor) Request() interface{} {
	return e.arg
}

func (e *RemoveDependerAppDataExecutor) Result() interface{} {
	return e.reply
}

func (e *RemoveDependerAppDataExecutor) Description() string {
	return fmt.Sprintf("[-] %s depender %s", e.arg.App, e.arg.Depender)
}

func (e *RemoveDependerAppDataExecutor) Authorize() error {
	return AuthorizeApp(&e.arg.ManagerAuthArg, e.arg.App)
}

func (e *RemoveDependerAppDataExecutor) Execute(t *Task) error {
	if e.arg.App == "" {
		return errors.New("Please specify an app")
	}
	if e.arg.Depender == "" {
		return errors.New("Please specify a depender app")
	}
	zkApp, err := datamodel.GetApp(e.arg.App)
	if err != nil {
		e.reply.Status = StatusError
		return err
	}
	err = zkApp.RemoveDependerAppData(e.arg.Depender)
	if err != nil {
		e.reply.Status = StatusError
		return err
	}
	e.reply.Status = StatusOk
	castedApp := App(*zkApp)
	e.reply.App = &castedApp
	return err
}

func (m *ManagerRPC) RemoveDependerAppData(arg ManagerRemoveDependerAppDataArg, reply *ManagerRemoveDependerAppDataReply) error {
	return NewTask("RemoveDependerAppData", &RemoveDependerAppDataExecutor{arg, reply}).Run()
}

type GetDependerAppDataExecutor struct {
	arg   ManagerGetDependerAppDataArg
	reply *ManagerGetDependerAppDataReply
}

func (e *GetDependerAppDataExecutor) Request() interface{} {
	return e.arg
}

func (e *GetDependerAppDataExecutor) Result() interface{} {
	return e.reply
}

func (e *GetDependerAppDataExecutor) Description() string {
	return fmt.Sprintf("GET %s depender %s", e.arg.App, e.arg.Depender)
}

func (e *GetDependerAppDataExecutor) Authorize() error {
	return AuthorizeApp(&e.arg.ManagerAuthArg, e.arg.App)
}

func (e *GetDependerAppDataExecutor) Execute(t *Task) error {
	if e.arg.App == "" {
		return errors.New("Please specify an app")
	}
	if e.arg.Depender == "" {
		return errors.New("Please specify a depender app")
	}
	zkApp, err := datamodel.GetApp(e.arg.App)
	if err != nil {
		e.reply.Status = StatusError
		return err
	}
	data := zkApp.GetDependerAppData(e.arg.Depender, true)
	if err != nil {
		e.reply.Status = StatusError
		return err
	}
	e.reply.Status = StatusOk
	e.reply.DependerAppData = data
	return err
}

func (m *ManagerRPC) GetDependerAppData(arg ManagerGetDependerAppDataArg, reply *ManagerGetDependerAppDataReply) error {
	return NewTask("GetDependerAppData", &GetDependerAppDataExecutor{arg, reply}).Run()
}

// ----------------------------------------------------------------------------------------------------------
// Depender Env Data Methods
// ----------------------------------------------------------------------------------------------------------

type AddDependerEnvDataExecutor struct {
	arg   ManagerAddDependerEnvDataArg
	reply *ManagerAddDependerEnvDataReply
}

func (e *AddDependerEnvDataExecutor) Request() interface{} {
	return e.arg
}

func (e *AddDependerEnvDataExecutor) Result() interface{} {
	return e.reply
}

func (e *AddDependerEnvDataExecutor) Description() string {
	return fmt.Sprintf("[+] %s env %+v", e.arg.App, e.arg.DependerEnvData)
}

func (e *AddDependerEnvDataExecutor) Authorize() error {
	return AuthorizeApp(&e.arg.ManagerAuthArg, e.arg.App)
}

func (e *AddDependerEnvDataExecutor) Execute(t *Task) error {
	if e.arg.App == "" {
		return errors.New("Please specify an app")
	}
	if e.arg.DependerEnvData == nil {
		return errors.New("Please specify data for the env")
	} else if e.arg.DependerEnvData.Name == "" {
		return errors.New("Please specify name for the env")
	}
	zkApp, err := datamodel.GetApp(e.arg.App)
	if err != nil {
		e.reply.Status = StatusError
		return err
	}
	err = zkApp.AddDependerEnvData(e.arg.DependerEnvData)
	if err != nil {
		e.reply.Status = StatusError
		return err
	}
	e.reply.Status = StatusOk
	castedApp := App(*zkApp)
	e.reply.App = &castedApp
	return err
}

func (m *ManagerRPC) AddDependerEnvData(arg ManagerAddDependerEnvDataArg, reply *ManagerAddDependerEnvDataReply) error {
	return NewTask("AddDependerEnvData", &AddDependerEnvDataExecutor{arg, reply}).Run()
}

type RemoveDependerEnvDataExecutor struct {
	arg   ManagerRemoveDependerEnvDataArg
	reply *ManagerRemoveDependerEnvDataReply
}

func (e *RemoveDependerEnvDataExecutor) Request() interface{} {
	return e.arg
}

func (e *RemoveDependerEnvDataExecutor) Result() interface{} {
	return e.reply
}

func (e *RemoveDependerEnvDataExecutor) Description() string {
	return fmt.Sprintf("[-] %s env %s", e.arg.App, e.arg.Env)
}

func (e *RemoveDependerEnvDataExecutor) Authorize() error {
	return AuthorizeApp(&e.arg.ManagerAuthArg, e.arg.App)
}

func (e *RemoveDependerEnvDataExecutor) Execute(t *Task) error {
	if e.arg.App == "" {
		return errors.New("Please specify an app")
	}
	if e.arg.Env == "" {
		return errors.New("Please specify an env")
	}
	zkApp, err := datamodel.GetApp(e.arg.App)
	if err != nil {
		e.reply.Status = StatusError
		return err
	}
	err = zkApp.RemoveDependerEnvData(e.arg.Env)
	if err != nil {
		e.reply.Status = StatusError
		return err
	}
	e.reply.Status = StatusOk
	castedApp := App(*zkApp)
	e.reply.App = &castedApp
	return err
}

func (m *ManagerRPC) RemoveDependerEnvData(arg ManagerRemoveDependerEnvDataArg, reply *ManagerRemoveDependerEnvDataReply) error {
	return NewTask("RemoveDependerEnvData", &RemoveDependerEnvDataExecutor{arg, reply}).Run()
}

type GetDependerEnvDataExecutor struct {
	arg   ManagerGetDependerEnvDataArg
	reply *ManagerGetDependerEnvDataReply
}

func (e *GetDependerEnvDataExecutor) Request() interface{} {
	return e.arg
}

func (e *GetDependerEnvDataExecutor) Result() interface{} {
	return e.reply
}

func (e *GetDependerEnvDataExecutor) Description() string {
	return fmt.Sprintf("GET %s env %s", e.arg.App, e.arg.Env)
}

func (e *GetDependerEnvDataExecutor) Authorize() error {
	return AuthorizeApp(&e.arg.ManagerAuthArg, e.arg.App)
}

func (e *GetDependerEnvDataExecutor) Execute(t *Task) error {
	if e.arg.App == "" {
		return errors.New("Please specify an app")
	}
	if e.arg.Env == "" {
		return errors.New("Please specify an env")
	}
	zkApp, err := datamodel.GetApp(e.arg.App)
	if err != nil {
		e.reply.Status = StatusError
		return err
	}
	data := zkApp.GetDependerEnvData(e.arg.Env, true)
	if err != nil {
		e.reply.Status = StatusError
		return err
	}
	e.reply.Status = StatusOk
	e.reply.DependerEnvData = data
	return err
}

func (m *ManagerRPC) GetDependerEnvData(arg ManagerGetDependerEnvDataArg, reply *ManagerGetDependerEnvDataReply) error {
	return NewTask("GetDependerEnvData", &GetDependerEnvDataExecutor{arg, reply}).Run()
}

// ----------------------------------------------------------------------------------------------------------
// Depender Env Data for Depender App Methods
// ----------------------------------------------------------------------------------------------------------

type AddDependerEnvDataForDependerAppExecutor struct {
	arg   ManagerAddDependerEnvDataForDependerAppArg
	reply *ManagerAddDependerEnvDataForDependerAppReply
}

func (e *AddDependerEnvDataForDependerAppExecutor) Request() interface{} {
	return e.arg
}

func (e *AddDependerEnvDataForDependerAppExecutor) Result() interface{} {
	return e.reply
}

func (e *AddDependerEnvDataForDependerAppExecutor) Description() string {
	return fmt.Sprintf("[+] %s depender %s env %+v", e.arg.App, e.arg.Depender, e.arg.DependerEnvData)
}

func (e *AddDependerEnvDataForDependerAppExecutor) Authorize() error {
	return AuthorizeApp(&e.arg.ManagerAuthArg, e.arg.App)
}

func (e *AddDependerEnvDataForDependerAppExecutor) Execute(t *Task) error {
	if e.arg.App == "" {
		return errors.New("Please specify an app")
	}
	if e.arg.Depender == "" {
		return errors.New("Please specify a depender app")
	}
	if e.arg.DependerEnvData == nil {
		return errors.New("Please specify data for the env")
	} else if e.arg.DependerEnvData.Name == "" {
		return errors.New("Please specify name for the env")
	}
	zkApp, err := datamodel.GetApp(e.arg.App)
	if err != nil {
		e.reply.Status = StatusError
		return err
	}
	err = zkApp.AddDependerEnvDataForDependerApp(e.arg.Depender, e.arg.DependerEnvData)
	if err != nil {
		e.reply.Status = StatusError
		return err
	}
	e.reply.Status = StatusOk
	castedApp := App(*zkApp)
	e.reply.App = &castedApp
	return err
}

func (m *ManagerRPC) AddDependerEnvDataForDependerApp(arg ManagerAddDependerEnvDataForDependerAppArg,
	reply *ManagerAddDependerEnvDataForDependerAppReply) error {
	return NewTask("AddDependerEnvDataForDependerApp", &AddDependerEnvDataForDependerAppExecutor{arg, reply}).Run()
}

type RemoveDependerEnvDataForDependerAppExecutor struct {
	arg   ManagerRemoveDependerEnvDataForDependerAppArg
	reply *ManagerRemoveDependerEnvDataForDependerAppReply
}

func (e *RemoveDependerEnvDataForDependerAppExecutor) Request() interface{} {
	return e.arg
}

func (e *RemoveDependerEnvDataForDependerAppExecutor) Result() interface{} {
	return e.reply
}

func (e *RemoveDependerEnvDataForDependerAppExecutor) Description() string {
	return fmt.Sprintf("[-] %s depender %s env %s", e.arg.App, e.arg.Depender, e.arg.Env)
}

func (e *RemoveDependerEnvDataForDependerAppExecutor) Authorize() error {
	return AuthorizeApp(&e.arg.ManagerAuthArg, e.arg.App)
}

func (e *RemoveDependerEnvDataForDependerAppExecutor) Execute(t *Task) error {
	if e.arg.App == "" {
		return errors.New("Please specify an app")
	}
	if e.arg.Depender == "" {
		return errors.New("Please specify a depender app")
	}
	if e.arg.Env == "" {
		return errors.New("Please specify an env")
	}
	zkApp, err := datamodel.GetApp(e.arg.App)
	if err != nil {
		e.reply.Status = StatusError
		return err
	}
	err = zkApp.RemoveDependerEnvDataForDependerApp(e.arg.Depender, e.arg.Env)
	if err != nil {
		e.reply.Status = StatusError
		return err
	}
	e.reply.Status = StatusOk
	castedApp := App(*zkApp)
	e.reply.App = &castedApp
	return err
}

func (m *ManagerRPC) RemoveDependerEnvDataForDependerApp(arg ManagerRemoveDependerEnvDataForDependerAppArg,
	reply *ManagerRemoveDependerEnvDataForDependerAppReply) error {
	return NewTask("RemoveDependerEnvDataForDependerApp", &RemoveDependerEnvDataForDependerAppExecutor{arg, reply}).Run()
}

type GetDependerEnvDataForDependerAppExecutor struct {
	arg   ManagerGetDependerEnvDataForDependerAppArg
	reply *ManagerGetDependerEnvDataForDependerAppReply
}

func (e *GetDependerEnvDataForDependerAppExecutor) Request() interface{} {
	return e.arg
}

func (e *GetDependerEnvDataForDependerAppExecutor) Result() interface{} {
	return e.reply
}

func (e *GetDependerEnvDataForDependerAppExecutor) Description() string {
	return fmt.Sprintf("GET %s depender %s env %s", e.arg.App, e.arg.Depender, e.arg.Env)
}

func (e *GetDependerEnvDataForDependerAppExecutor) Authorize() error {
	return AuthorizeApp(&e.arg.ManagerAuthArg, e.arg.App)
}

func (e *GetDependerEnvDataForDependerAppExecutor) Execute(t *Task) error {
	if e.arg.App == "" {
		return errors.New("Please specify an app")
	}
	if e.arg.Depender == "" {
		return errors.New("Please specify a depender app")
	}
	if e.arg.Env == "" {
		return errors.New("Please specify an env")
	}
	zkApp, err := datamodel.GetApp(e.arg.App)
	if err != nil {
		e.reply.Status = StatusError
		return err
	}
	data := zkApp.GetDependerEnvDataForDependerApp(e.arg.Depender, e.arg.Env, true)
	if err != nil {
		e.reply.Status = StatusError
		return err
	}
	e.reply.Status = StatusOk
	e.reply.DependerEnvData = data
	return err
}

func (m *ManagerRPC) GetDependerEnvDataForDependerApp(arg ManagerGetDependerEnvDataForDependerAppArg,
	reply *ManagerGetDependerEnvDataForDependerAppReply) error {
	return NewTask("GetDependerEnvDataForDependerApp", &GetDependerEnvDataForDependerAppExecutor{arg, reply}).Run()
}
