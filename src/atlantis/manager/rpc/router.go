package rpc

import (
	. "atlantis/common"
	. "atlantis/manager/constant"
	"atlantis/manager/datamodel"
	"atlantis/manager/dns"
	"atlantis/manager/helper"
	. "atlantis/manager/rpc/types"
	routerzk "atlantis/router/zk"
	"errors"
	"fmt"
	"strings"
)

// ----------------------------------------------------------------------------------------------------------
// Port Related
// ----------------------------------------------------------------------------------------------------------

type GetAppEnvPortExecutor struct {
	arg   ManagerGetAppEnvPortArg
	reply *ManagerGetAppEnvPortReply
}

func (e *GetAppEnvPortExecutor) Request() interface{} {
	return e.arg
}

func (e *GetAppEnvPortExecutor) Result() interface{} {
	return e.reply
}

func (e *GetAppEnvPortExecutor) Description() string {
	return fmt.Sprintf("%s in %s", e.arg.App, e.arg.Env)
}

func (e *GetAppEnvPortExecutor) Execute(t *Task) (err error) {
	if e.arg.App == "" {
		return errors.New("Please specify an app")
	} else if e.arg.Env == "" {
		return errors.New("Please specify an environment")
	}
	zkApp, err := datamodel.GetApp(e.arg.App)
	if err != nil {
		return err
	}
	helper.SetRouterRoot(zkApp.Internal)
	e.reply.Port, err = routerzk.GetPort(datamodel.Zk.Conn, helper.GetAppEnvTrieName(e.arg.App, e.arg.Env))
	return err
}

func (e *GetAppEnvPortExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

type ListAppEnvsWithPortExecutor struct {
	arg   ManagerListAppEnvsWithPortArg
	reply *ManagerListAppEnvsWithPortReply
}

func (e *ListAppEnvsWithPortExecutor) Request() interface{} {
	return e.arg
}

func (e *ListAppEnvsWithPortExecutor) Result() interface{} {
	return e.reply
}

func (e *ListAppEnvsWithPortExecutor) Description() string {
	return fmt.Sprintf("internal: %t", e.arg.Internal)
}

func (e *ListAppEnvsWithPortExecutor) Execute(t *Task) (err error) {
	zrp := datamodel.GetRouterPorts(e.arg.Internal)
	e.reply.AppEnvs = []AppEnv{}
	for _, appEnv := range zrp.PortMap {
		e.reply.AppEnvs = append(e.reply.AppEnvs, appEnv)
	}
	return err
}

func (e *ListAppEnvsWithPortExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

type UpdatePortExecutor struct {
	arg   ManagerUpdatePortArg
	reply *ManagerUpdatePortReply
}

func (e *UpdatePortExecutor) Request() interface{} {
	return e.arg
}

func (e *UpdatePortExecutor) Result() interface{} {
	return e.reply
}

func (e *UpdatePortExecutor) Description() string {
	return fmt.Sprintf("%+v", e.arg.Port)
}

func (e *UpdatePortExecutor) Execute(t *Task) (err error) {
	if e.arg.Port.Name == "" {
		return errors.New("Please specify a name")
	}
	if e.arg.Port.Trie == "" {
		return errors.New("Please specify a trie")
	}
	if e.arg.Port.Port == uint16(0) {
		return errors.New("Please specify a port")
	}
	helper.SetRouterRoot(e.arg.Port.Internal)
	return routerzk.SetPort(datamodel.Zk.Conn, e.arg.Port)
}

func (e *UpdatePortExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

type DeletePortExecutor struct {
	arg   ManagerDeletePortArg
	reply *ManagerDeletePortReply
}

func (e *DeletePortExecutor) Request() interface{} {
	return e.arg
}

func (e *DeletePortExecutor) Result() interface{} {
	return e.reply
}

func (e *DeletePortExecutor) Description() string {
	return e.arg.Name
}

func (e *DeletePortExecutor) Execute(t *Task) (err error) {
	if e.arg.Name == "" {
		return errors.New("Please specify a name")
	}
	helper.SetRouterRoot(e.arg.Internal)
	err = routerzk.DelPort(datamodel.Zk.Conn, e.arg.Name)
	return err
}

func (e *DeletePortExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

type GetPortExecutor struct {
	arg   ManagerGetPortArg
	reply *ManagerGetPortReply
}

func (e *GetPortExecutor) Request() interface{} {
	return e.arg
}

func (e *GetPortExecutor) Result() interface{} {
	return e.reply
}

func (e *GetPortExecutor) Description() string {
	return e.arg.Name
}

func (e *GetPortExecutor) Execute(t *Task) (err error) {
	if e.arg.Name == "" {
		return errors.New("Please specify a name")
	}
	helper.SetRouterRoot(e.arg.Internal)
	e.reply.Port, err = routerzk.GetPort(datamodel.Zk.Conn, e.arg.Name)
	return err
}

func (e *GetPortExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

type ListPortsExecutor struct {
	arg   ManagerListPortsArg
	reply *ManagerListPortsReply
}

func (e *ListPortsExecutor) Request() interface{} {
	return e.arg
}

func (e *ListPortsExecutor) Result() interface{} {
	return e.reply
}

func (e *ListPortsExecutor) Description() string {
	return "ListPorts"
}

func (e *ListPortsExecutor) Execute(t *Task) (err error) {
	helper.SetRouterRoot(e.arg.Internal)
	e.reply.Ports, err = routerzk.ListPorts(datamodel.Zk.Conn)
	return err
}

func (e *ListPortsExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

func (m *ManagerRPC) GetAppEnvPort(arg ManagerGetAppEnvPortArg, reply *ManagerGetAppEnvPortReply) error {
	return NewTask("GetAppEnvPort", &GetAppEnvPortExecutor{arg, reply}).Run()
}

func (m *ManagerRPC) ListAppEnvsWithPort(arg ManagerListAppEnvsWithPortArg, reply *ManagerListAppEnvsWithPortReply) error {
	return NewTask("ListAppEnvsWithPort", &ListAppEnvsWithPortExecutor{arg, reply}).Run()
}

func (m *ManagerRPC) UpdatePort(arg ManagerUpdatePortArg, reply *ManagerUpdatePortReply) error {
	return NewTask("UpdatePort", &UpdatePortExecutor{arg, reply}).Run()
}

func (m *ManagerRPC) DeletePort(arg ManagerDeletePortArg, reply *ManagerDeletePortReply) error {
	return NewTask("DeletePort", &DeletePortExecutor{arg, reply}).Run()
}

func (m *ManagerRPC) GetPort(arg ManagerGetPortArg, reply *ManagerGetPortReply) error {
	return NewTask("GetPort", &GetPortExecutor{arg, reply}).Run()
}

func (m *ManagerRPC) ListPorts(arg ManagerListPortsArg, reply *ManagerListPortsReply) error {
	return NewTask("ListPorts", &ListPortsExecutor{arg, reply}).Run()
}

// ----------------------------------------------------------------------------------------------------------
// Pool Related
// ----------------------------------------------------------------------------------------------------------

type UpdatePoolExecutor struct {
	arg   ManagerUpdatePoolArg
	reply *ManagerUpdatePoolReply
}

func (e *UpdatePoolExecutor) Request() interface{} {
	return e.arg
}

func (e *UpdatePoolExecutor) Result() interface{} {
	return e.reply
}

func (e *UpdatePoolExecutor) Description() string {
	return fmt.Sprintf("%+v", e.arg.Pool)
}

func (e *UpdatePoolExecutor) Execute(t *Task) error {
	if e.arg.Pool.Name == "" {
		return errors.New("Please specify a name")
	} else if e.arg.Pool.Config.HealthzEvery == "" {
		return errors.New("Please specify a healthz check frequency")
	} else if e.arg.Pool.Config.HealthzTimeout == "" {
		return errors.New("Please specify a healthz timeout")
	} else if e.arg.Pool.Config.RequestTimeout == "" {
		return errors.New("Please specify a request timeout")
	} // no need to check hosts. an empty pool is still a valid pool
	helper.SetRouterRoot(e.arg.Pool.Internal)
	err := routerzk.SetPool(datamodel.Zk.Conn, e.arg.Pool)
	if err != nil {
		e.reply.Status = StatusError
	} else {
		e.reply.Status = StatusOk
	}
	return err
}

func (e *UpdatePoolExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

type DeletePoolExecutor struct {
	arg   ManagerDeletePoolArg
	reply *ManagerDeletePoolReply
}

func (e *DeletePoolExecutor) Request() interface{} {
	return e.arg
}

func (e *DeletePoolExecutor) Result() interface{} {
	return e.reply
}

func (e *DeletePoolExecutor) Description() string {
	return fmt.Sprintf("%s", e.arg.Name)
}

func (e *DeletePoolExecutor) Execute(t *Task) (err error) {
	if e.arg.Name == "" {
		return errors.New("Please specify a name")
	}
	helper.SetRouterRoot(e.arg.Internal)
	err = routerzk.DelPool(datamodel.Zk.Conn, e.arg.Name)
	if err != nil {
		e.reply.Status = StatusError
	} else {
		e.reply.Status = StatusOk
	}
	return err
}

func (e *DeletePoolExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

type GetPoolExecutor struct {
	arg   ManagerGetPoolArg
	reply *ManagerGetPoolReply
}

func (e *GetPoolExecutor) Request() interface{} {
	return e.arg
}

func (e *GetPoolExecutor) Result() interface{} {
	return e.reply
}

func (e *GetPoolExecutor) Description() string {
	return fmt.Sprintf("%s", e.arg.Name)
}

func (e *GetPoolExecutor) Execute(t *Task) (err error) {
	if e.arg.Name == "" {
		return errors.New("Please specify a name")
	}
	helper.SetRouterRoot(e.arg.Internal)
	e.reply.Pool, err = routerzk.GetPool(datamodel.Zk.Conn, e.arg.Name)
	if err != nil {
		e.reply.Status = StatusError
	} else {
		e.reply.Status = StatusOk
	}
	return err
}

func (e *GetPoolExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

type ListPoolsExecutor struct {
	arg   ManagerListPoolsArg
	reply *ManagerListPoolsReply
}

func (e *ListPoolsExecutor) Request() interface{} {
	return e.arg
}

func (e *ListPoolsExecutor) Result() interface{} {
	return e.reply
}

func (e *ListPoolsExecutor) Description() string {
	return "ListPools"
}

func (e *ListPoolsExecutor) Execute(t *Task) (err error) {
	helper.SetRouterRoot(e.arg.Internal)
	e.reply.Pools, err = routerzk.ListPools(datamodel.Zk.Conn)
	if err != nil {
		e.reply.Status = StatusError
	} else {
		e.reply.Status = StatusOk
	}
	return err
}

func (e *ListPoolsExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

func (m *ManagerRPC) UpdatePool(arg ManagerUpdatePoolArg, reply *ManagerUpdatePoolReply) error {
	return NewTask("UpdatePool", &UpdatePoolExecutor{arg, reply}).Run()
}

func (m *ManagerRPC) DeletePool(arg ManagerDeletePoolArg, reply *ManagerDeletePoolReply) error {
	return NewTask("DeletePool", &DeletePoolExecutor{arg, reply}).Run()
}

func (m *ManagerRPC) GetPool(arg ManagerGetPoolArg, reply *ManagerGetPoolReply) error {
	return NewTask("GetPool", &GetPoolExecutor{arg, reply}).Run()
}

func (m *ManagerRPC) ListPools(arg ManagerListPoolsArg, reply *ManagerListPoolsReply) error {
	return NewTask("ListPools", &ListPoolsExecutor{arg, reply}).Run()
}

// ----------------------------------------------------------------------------------------------------------
// Rule Related
// ----------------------------------------------------------------------------------------------------------

type UpdateRuleExecutor struct {
	arg   ManagerUpdateRuleArg
	reply *ManagerUpdateRuleReply
}

func (e *UpdateRuleExecutor) Request() interface{} {
	return e.arg
}

func (e *UpdateRuleExecutor) Result() interface{} {
	return e.reply
}

func (e *UpdateRuleExecutor) Description() string {
	return fmt.Sprintf("%+v", e.arg.Rule)
}

func (e *UpdateRuleExecutor) Execute(t *Task) (err error) {
	if e.arg.Rule.Name == "" {
		return errors.New("Please specify a name")
	} else if e.arg.Rule.Type == "" {
		return errors.New("Please specify a type")
	} else if e.arg.Rule.Value == "" {
		return errors.New("Please specify a value")
	} else if e.arg.Rule.Next == "" && e.arg.Rule.Pool == "" {
		return errors.New("Please specify either a next trie or a pool")
	}
	// fill in current cname suffixes in multi-host rules
	if e.arg.Rule.Type == "multi-host" && dns.Provider != nil {
		suffix, err := dns.Provider.Suffix(Region)
		if err != nil {
			return err
		}
		list := strings.Join(helper.GetAppCNameSuffixes(suffix), ",")
		e.arg.Rule.Value = fmt.Sprintf("%s:%s", e.arg.Rule.Value, list)
	}

	helper.SetRouterRoot(e.arg.Rule.Internal)
	err = routerzk.SetRule(datamodel.Zk.Conn, e.arg.Rule)
	if err != nil {
		e.reply.Status = StatusError
	} else {
		e.reply.Status = StatusOk
	}
	return err
}

func (e *UpdateRuleExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

type DeleteRuleExecutor struct {
	arg   ManagerDeleteRuleArg
	reply *ManagerDeleteRuleReply
}

func (e *DeleteRuleExecutor) Request() interface{} {
	return e.arg
}

func (e *DeleteRuleExecutor) Result() interface{} {
	return e.reply
}

func (e *DeleteRuleExecutor) Description() string {
	return fmt.Sprintf("%s", e.arg.Name)
}

func (e *DeleteRuleExecutor) Execute(t *Task) (err error) {
	if e.arg.Name == "" {
		return errors.New("Please specify a name")
	}
	helper.SetRouterRoot(e.arg.Internal)
	err = routerzk.DelRule(datamodel.Zk.Conn, e.arg.Name)
	if err != nil {
		e.reply.Status = StatusError
	} else {
		e.reply.Status = StatusOk
	}
	return err
}

func (e *DeleteRuleExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

type GetRuleExecutor struct {
	arg   ManagerGetRuleArg
	reply *ManagerGetRuleReply
}

func (e *GetRuleExecutor) Request() interface{} {
	return e.arg
}

func (e *GetRuleExecutor) Result() interface{} {
	return e.reply
}

func (e *GetRuleExecutor) Description() string {
	return fmt.Sprintf("%s", e.arg.Name)
}

func (e *GetRuleExecutor) Execute(t *Task) (err error) {
	if e.arg.Name == "" {
		return errors.New("Please specify a name")
	}
	helper.SetRouterRoot(e.arg.Internal)
	e.reply.Rule, err = routerzk.GetRule(datamodel.Zk.Conn, e.arg.Name)
	if err != nil {
		e.reply.Status = StatusError
	} else {
		e.reply.Status = StatusOk
	}
	return err
}

func (e *GetRuleExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

type ListRulesExecutor struct {
	arg   ManagerListRulesArg
	reply *ManagerListRulesReply
}

func (e *ListRulesExecutor) Request() interface{} {
	return e.arg
}

func (e *ListRulesExecutor) Result() interface{} {
	return e.reply
}

func (e *ListRulesExecutor) Description() string {
	return "ListRules"
}

func (e *ListRulesExecutor) Execute(t *Task) (err error) {
	helper.SetRouterRoot(e.arg.Internal)
	e.reply.Rules, err = routerzk.ListRules(datamodel.Zk.Conn)
	if err != nil {
		e.reply.Status = StatusError
	} else {
		e.reply.Status = StatusOk
	}
	return err
}

func (e *ListRulesExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

func (m *ManagerRPC) UpdateRule(arg ManagerUpdateRuleArg, reply *ManagerUpdateRuleReply) error {
	return NewTask("UpdateRule", &UpdateRuleExecutor{arg, reply}).Run()
}

func (m *ManagerRPC) DeleteRule(arg ManagerDeleteRuleArg, reply *ManagerDeleteRuleReply) error {
	return NewTask("DeleteRule", &DeleteRuleExecutor{arg, reply}).Run()
}

func (m *ManagerRPC) GetRule(arg ManagerGetRuleArg, reply *ManagerGetRuleReply) error {
	return NewTask("GetRule", &GetRuleExecutor{arg, reply}).Run()
}

func (m *ManagerRPC) ListRules(arg ManagerListRulesArg, reply *ManagerListRulesReply) error {
	return NewTask("ListRules", &ListRulesExecutor{arg, reply}).Run()
}

// ----------------------------------------------------------------------------------------------------------
// Trie Related
// ----------------------------------------------------------------------------------------------------------

type UpdateTrieExecutor struct {
	arg   ManagerUpdateTrieArg
	reply *ManagerUpdateTrieReply
}

func (e *UpdateTrieExecutor) Request() interface{} {
	return e.arg
}

func (e *UpdateTrieExecutor) Result() interface{} {
	return e.reply
}

func (e *UpdateTrieExecutor) Description() string {
	return fmt.Sprintf("%+v", e.arg.Trie)
}

func (e *UpdateTrieExecutor) Execute(t *Task) (err error) {
	if e.arg.Trie.Name == "" {
		return errors.New("Please specify a name")
	}
	helper.SetRouterRoot(e.arg.Trie.Internal)
	err = routerzk.SetTrie(datamodel.Zk.Conn, e.arg.Trie)
	if err != nil {
		e.reply.Status = StatusError
	} else {
		e.reply.Status = StatusOk
	}
	return err
}

func (e *UpdateTrieExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

type DeleteTrieExecutor struct {
	arg   ManagerDeleteTrieArg
	reply *ManagerDeleteTrieReply
}

func (e *DeleteTrieExecutor) Request() interface{} {
	return e.arg
}

func (e *DeleteTrieExecutor) Result() interface{} {
	return e.reply
}

func (e *DeleteTrieExecutor) Description() string {
	return fmt.Sprintf("%s", e.arg.Name)
}

func (e *DeleteTrieExecutor) Execute(t *Task) (err error) {
	if e.arg.Name == "" {
		return errors.New("Please specify a name")
	}
	helper.SetRouterRoot(e.arg.Internal)
	err = routerzk.DelTrie(datamodel.Zk.Conn, e.arg.Name)
	if err != nil {
		e.reply.Status = StatusError
	} else {
		e.reply.Status = StatusOk
	}
	return err
}

func (e *DeleteTrieExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

type GetTrieExecutor struct {
	arg   ManagerGetTrieArg
	reply *ManagerGetTrieReply
}

func (e *GetTrieExecutor) Request() interface{} {
	return e.arg
}

func (e *GetTrieExecutor) Result() interface{} {
	return e.reply
}

func (e *GetTrieExecutor) Description() string {
	return fmt.Sprintf("%s", e.arg.Name)
}

func (e *GetTrieExecutor) Execute(t *Task) (err error) {
	if e.arg.Name == "" {
		return errors.New("Please specify a name")
	}
	helper.SetRouterRoot(e.arg.Internal)
	e.reply.Trie, err = routerzk.GetTrie(datamodel.Zk.Conn, e.arg.Name)
	if err != nil {
		e.reply.Status = StatusError
	} else {
		e.reply.Status = StatusOk
	}
	return err
}

func (e *GetTrieExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

type ListTriesExecutor struct {
	arg   ManagerListTriesArg
	reply *ManagerListTriesReply
}

func (e *ListTriesExecutor) Request() interface{} {
	return e.arg
}

func (e *ListTriesExecutor) Result() interface{} {
	return e.reply
}

func (e *ListTriesExecutor) Description() string {
	return "ListTries"
}

func (e *ListTriesExecutor) Execute(t *Task) (err error) {
	helper.SetRouterRoot(e.arg.Internal)
	e.reply.Tries, err = routerzk.ListTries(datamodel.Zk.Conn)
	if err != nil {
		e.reply.Status = StatusError
	} else {
		e.reply.Status = StatusOk
	}
	return err
}

func (e *ListTriesExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

func (m *ManagerRPC) UpdateTrie(arg ManagerUpdateTrieArg, reply *ManagerUpdateTrieReply) error {
	return NewTask("UpdateTrie", &UpdateTrieExecutor{arg, reply}).Run()
}

func (m *ManagerRPC) DeleteTrie(arg ManagerDeleteTrieArg, reply *ManagerDeleteTrieReply) error {
	return NewTask("DeleteTrie", &DeleteTrieExecutor{arg, reply}).Run()
}

func (m *ManagerRPC) GetTrie(arg ManagerGetTrieArg, reply *ManagerGetTrieReply) error {
	return NewTask("GetTrie", &GetTrieExecutor{arg, reply}).Run()
}

func (m *ManagerRPC) ListTries(arg ManagerListTriesArg, reply *ManagerListTriesReply) error {
	return NewTask("ListTries", &ListTriesExecutor{arg, reply}).Run()
}
