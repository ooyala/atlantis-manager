package rpc

import (
	. "atlantis/common"
	"atlantis/manager/datamodel"
	"atlantis/manager/helper"
	. "atlantis/manager/rpc/types"
	routerzk "atlantis/router/zk"
	"errors"
	"fmt"
)

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
