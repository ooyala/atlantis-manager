/* Copyright 2014 Ooyala, Inc. All rights reserved.
 *
 * This file is licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
 * except in compliance with the License. You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is
 * distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

package client

import (
	. "atlantis/manager/rpc/types"
	"atlantis/router/config"
)

type UpdatePoolCommand struct {
	Name             string   `short:"n" long:"name" description:"the name of the pool"`
	HealthCheckEvery string   `short:"e" long:"check-every" default:"5s" description:"how often to check healthz"`
	HealthzTimeout   string   `short:"z" long:"healthz-timeout" default:"5s" description:"timeout for healthz checks"`
	RequestTimeout   string   `short:"r" long:"request-timeout" default:"120s" description:"timeout for requests"`
	Status           string   `short:"s" long:"status" default:"OK" description:"the pool's status"`
	Hosts            []string `short:"H" long:"host" description:"the pool's hosts"`
	Internal         bool     `short:"i" long:"internal" description:"true if internal"`
}

func (c *UpdatePoolCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Update Pool...")
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	hosts := make(map[string]config.Host, len(c.Hosts))
	for _, host := range c.Hosts {
		hosts[host] = config.Host{Address: host}
	}
	arg := ManagerUpdatePoolArg{authArg, config.Pool{Name: c.Name, Hosts: hosts, Internal: c.Internal,
		Config: config.PoolConfig{HealthzEvery: c.HealthCheckEvery, HealthzTimeout: c.HealthzTimeout,
			RequestTimeout: c.RequestTimeout, Status: c.Status}}}
	var reply ManagerUpdatePoolReply
	err = rpcClient.Call("UpdatePool", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	return Output(map[string]interface{}{"status": reply.Status}, nil, nil)
}

type DeletePoolCommand struct {
	Name     string `short:"n" long:"name" description:"the name of the pool"`
	Internal bool   `short:"i" long:"internal" description:"true if internal"`
}

func (c *DeletePoolCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Delete Pool...")
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerDeletePoolArg{authArg, c.Name, c.Internal}
	var reply ManagerDeletePoolReply
	err = rpcClient.Call("DeletePool", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	return Output(map[string]interface{}{"status": reply.Status}, nil, nil)
}

type GetPoolCommand struct {
	Name     string `short:"n" long:"name" description:"the name of the pool"`
	Internal bool   `short:"i" long:"internal" description:"true if internal"`
}

func (c *GetPoolCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Get Pool...")
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerGetPoolArg{authArg, c.Name, c.Internal}
	var reply ManagerGetPoolReply
	err = rpcClient.Call("GetPool", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	Log("-> %v", reply.Pool.String())
	return Output(map[string]interface{}{"status": reply.Status, "pool": reply.Pool}, reply.Pool, nil)
}

type ListPoolsCommand struct {
	Internal bool `short:"i" long:"internal" description:"true if internal"`
}

func (c *ListPoolsCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("List Pools...")
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerListPoolsArg{authArg, c.Internal}
	var reply ManagerListPoolsReply
	err = rpcClient.Call("ListPools", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	Log("-> pools:")
	for _, pool := range reply.Pools {
		Log("->   %s", pool)
	}
	return Output(map[string]interface{}{"status": reply.Status, "pools": reply.Pools}, reply.Pools, nil)
}

type UpdateRuleCommand struct {
	Name     string `short:"n" long:"name" description:"the name of the rule"`
	Type     string `short:"t" long:"type" description:"the type of the rule"`
	Value    string `short:"v" long:"value" description:"the rule's value"`
	Next     string `short:"x" long:"next" description:"the next ruleset"`
	Pool     string `short:"p" long:"pool" description:"the pool to point to if this rule succeeds"`
	Internal bool   `short:"i" long:"internal" description:"true if internal"`
}

func (c *UpdateRuleCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Update Rule...")
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerUpdateRuleArg{authArg, config.Rule{Name: c.Name, Type: c.Type, Value: c.Value, Next: c.Next,
		Pool: c.Pool, Internal: c.Internal}}
	var reply ManagerUpdateRuleReply
	err = rpcClient.Call("UpdateRule", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	return Output(map[string]interface{}{"status": reply.Status}, nil, nil)
}

type DeleteRuleCommand struct {
	Name     string `short:"n" long:"name" description:"the name of the rule"`
	Internal bool   `short:"i" long:"internal" description:"true if internal"`
}

func (c *DeleteRuleCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Delete Rule...")
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerDeleteRuleArg{authArg, c.Name, c.Internal}
	var reply ManagerDeleteRuleReply
	err = rpcClient.Call("DeleteRule", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	return Output(map[string]interface{}{"status": reply.Status}, nil, nil)
}

type GetRuleCommand struct {
	Name     string `short:"n" long:"name" description:"the name of the rule"`
	Internal bool   `short:"i" long:"internal" description:"true if internal"`
}

func (c *GetRuleCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Get Rule...")
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerGetRuleArg{authArg, c.Name, c.Internal}
	var reply ManagerGetRuleReply
	err = rpcClient.Call("GetRule", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	Log("-> %v", reply.Rule.String())
	return Output(map[string]interface{}{"status": reply.Status, "rule": reply.Rule}, reply.Rule, nil)
}

type ListRulesCommand struct {
	Internal bool `short:"i" long:"internal" description:"true if internal"`
}

func (c *ListRulesCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("List Rules...")
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerListRulesArg{authArg, c.Internal}
	var reply ManagerListRulesReply
	err = rpcClient.Call("ListRules", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	Log("-> rules:")
	for _, rule := range reply.Rules {
		Log("->   %s", rule)
	}
	return Output(map[string]interface{}{"status": reply.Status, "rules": reply.Rules}, reply.Rules, nil)
}

type UpdateTrieCommand struct {
	Name     string   `short:"n" long:"name" description:"the name of the rule"`
	Rules    []string `short:"r" long:"rule" description:"the rules that make up the ruleset"`
	Internal bool     `short:"i" long:"internal" description:"true if internal"`
}

func (c *UpdateTrieCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Update Trie...")
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerUpdateTrieArg{authArg, config.Trie{Name: c.Name, Rules: c.Rules, Internal: c.Internal}}
	var reply ManagerUpdateTrieReply
	err = rpcClient.Call("UpdateTrie", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	return Output(map[string]interface{}{"status": reply.Status}, nil, nil)
}

type DeleteTrieCommand struct {
	Name     string `short:"n" long:"name" description:"the name of the trie"`
	Internal bool   `short:"i" long:"internal" description:"true if internal"`
}

func (c *DeleteTrieCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Delete Trie...")
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerDeleteTrieArg{authArg, c.Name, c.Internal}
	var reply ManagerDeleteTrieReply
	err = rpcClient.Call("DeleteTrie", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	return Output(map[string]interface{}{"status": reply.Status}, nil, nil)
}

type GetTrieCommand struct {
	Name     string `short:"n" long:"name" description:"the name of the trie"`
	Internal bool   `short:"i" long:"internal" description:"true if internal"`
}

func (c *GetTrieCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Get Trie...")
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerGetTrieArg{authArg, c.Name, c.Internal}
	var reply ManagerGetTrieReply
	err = rpcClient.Call("GetTrie", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	Log("-> %v", reply.Trie.String())
	return Output(map[string]interface{}{"status": reply.Status, "trie": reply.Trie}, reply.Trie, nil)
}

type ListTriesCommand struct {
	Internal bool `short:"i" long:"internal" description:"true if internal"`
}

func (c *ListTriesCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("List Tries...")
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerListTriesArg{authArg, c.Internal}
	var reply ManagerListTriesReply
	err = rpcClient.Call("ListTries", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	Log("-> tries:")
	for _, trie := range reply.Tries {
		Log("->   %s", trie)
	}
	return Output(map[string]interface{}{"status": reply.Status, "tries": reply.Tries}, reply.Tries, nil)
}

type UpdatePortCommand struct {
	Port     uint16 `short:"p" long:"port" description:"the actual port to listen on"`
	Trie     string `short:"t" long:"trie" description:"the trie to use as root for this port"`
	Internal bool   `short:"i" long:"internal" description:"true if internal"`
}

func (c *UpdatePortCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Update Port...")
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerUpdatePortArg{
		ManagerAuthArg: authArg,
		Port: config.Port{
			Port:     c.Port,
			Trie:     c.Trie,
			Internal: c.Internal,
		},
	}
	var reply ManagerUpdatePortReply
	err = rpcClient.Call("UpdatePort", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	return Output(map[string]interface{}{"status": reply.Status}, nil, nil)
}

type DeletePortCommand struct {
	Port     uint16 `short:"p" long:"port" description:"the port number"`
	Internal bool   `short:"i" long:"internal" description:"true if internal"`
}

func (c *DeletePortCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Delete Port...")
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerDeletePortArg{authArg, c.Port, c.Internal}
	var reply ManagerDeletePortReply
	err = rpcClient.Call("DeletePort", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	return Output(map[string]interface{}{"status": reply.Status}, nil, nil)
}

type GetPortCommand struct {
	Port     uint16 `short:"p" long:"port" description:"the port number"`
	Internal bool   `short:"i" long:"internal" description:"true if internal"`
}

func (c *GetPortCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Get Port...")
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerGetPortArg{authArg, c.Port, c.Internal}
	var reply ManagerGetPortReply
	err = rpcClient.Call("GetPort", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	Log("-> %v", reply.Port.String())
	return Output(map[string]interface{}{"status": reply.Status, "port": reply.Port}, reply.Port, nil)
}

type ListPortsCommand struct {
	Internal bool `short:"i" long:"internal" description:"true if internal"`
}

func (c *ListPortsCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("List Ports...")
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerListPortsArg{authArg, c.Internal}
	var reply ManagerListPortsReply
	err = rpcClient.Call("ListPorts", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	Log("-> ports:")
	for _, port := range reply.Ports {
		Log("->   %d", port)
	}
	return Output(map[string]interface{}{"status": reply.Status, "ports": reply.Ports}, reply.Ports, nil)
}

type GetAppEnvPortCommand struct {
	App      string `short:"a" long:"app" description:"the app of the port"`
	Env      string `short:"e" long:"env" description:"the env of the port"`
	Internal bool   `short:"i" long:"internal" description:"true if internal"`
}

func (c *GetAppEnvPortCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Get AppEnv Port...")
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerGetAppEnvPortArg{
		ManagerAuthArg: authArg,
		App:            c.App,
		Env:            c.Env,
	}
	var reply ManagerGetAppEnvPortReply
	err = rpcClient.Call("GetAppEnvPort", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	Log("-> %v", reply.Port.String())
	return Output(map[string]interface{}{"status": reply.Status, "port": reply.Port}, reply.Port, nil)
}

type ListAppEnvsWithPortCommand struct {
	Internal bool `short:"i" long:"internal" description:"true if internal"`
}

func (c *ListAppEnvsWithPortCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("List AppEnvs With Ports...")
	user, secret, err := GetSecret()
	if err != nil {
		return err
	}
	authArg := ManagerAuthArg{user, "", secret}
	arg := ManagerListAppEnvsWithPortArg{authArg, c.Internal}
	var reply ManagerListAppEnvsWithPortReply
	err = rpcClient.Call("ListAppEnvsWithPort", arg, &reply)
	if err != nil {
		return OutputError(err)
	}
	Log("-> status: %s", reply.Status)
	Log("-> app+envs:")
	for _, appEnv := range reply.AppEnvs {
		Log("->   %s in %s", appEnv.App, appEnv.Env)
	}
	return Output(map[string]interface{}{"status": reply.Status, "appEnvs": reply.AppEnvs}, reply.AppEnvs, nil)
}
