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
	hosts := make(map[string]config.Host, len(c.Hosts))
	for _, host := range c.Hosts {
		hosts[host] = config.Host{Address: host}
	}
	arg := ManagerUpdatePoolArg{dummyAuthArg, config.Pool{Name: c.Name, Hosts: hosts, Internal: c.Internal,
		Config: config.PoolConfig{HealthzEvery: c.HealthCheckEvery, HealthzTimeout: c.HealthzTimeout,
			RequestTimeout: c.RequestTimeout, Status: c.Status}}}
	var reply ManagerUpdatePoolReply
	err = rpcClient.CallAuthed("UpdatePool", &arg, &reply)
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
	arg := ManagerDeletePoolArg{dummyAuthArg, c.Name, c.Internal}
	var reply ManagerDeletePoolReply
	err = rpcClient.CallAuthed("DeletePool", &arg, &reply)
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
	arg := ManagerGetPoolArg{dummyAuthArg, c.Name, c.Internal}
	var reply ManagerGetPoolReply
	err = rpcClient.CallAuthed("GetPool", &arg, &reply)
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
	arg := ManagerListPoolsArg{dummyAuthArg, c.Internal}
	var reply ManagerListPoolsReply
	err = rpcClient.CallAuthed("ListPools", &arg, &reply)
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
	arg := ManagerUpdateRuleArg{dummyAuthArg, config.Rule{Name: c.Name, Type: c.Type, Value: c.Value, Next: c.Next,
		Pool: c.Pool, Internal: c.Internal}}
	var reply ManagerUpdateRuleReply
	err = rpcClient.Call("UpdateRule", &arg, &reply)
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
	arg := ManagerDeleteRuleArg{dummyAuthArg, c.Name, c.Internal}
	var reply ManagerDeleteRuleReply
	err = rpcClient.CallAuthed("DeleteRule", &arg, &reply)
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
	arg := ManagerGetRuleArg{dummyAuthArg, c.Name, c.Internal}
	var reply ManagerGetRuleReply
	err = rpcClient.Call("GetRule", &arg, &reply)
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
	arg := ManagerListRulesArg{dummyAuthArg, c.Internal}
	var reply ManagerListRulesReply
	err = rpcClient.CallAuthed("ListRules", &arg, &reply)
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
	arg := ManagerUpdateTrieArg{dummyAuthArg, config.Trie{Name: c.Name, Rules: c.Rules, Internal: c.Internal}}
	var reply ManagerUpdateTrieReply
	err = rpcClient.CallAuthed("UpdateTrie", &arg, &reply)
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
	arg := ManagerDeleteTrieArg{dummyAuthArg, c.Name, c.Internal}
	var reply ManagerDeleteTrieReply
	err = rpcClient.CallAuthed("DeleteTrie", &arg, &reply)
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
	arg := ManagerGetTrieArg{dummyAuthArg, c.Name, c.Internal}
	var reply ManagerGetTrieReply
	err = rpcClient.CallAuthed("GetTrie", &arg, &reply)
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
	arg := ManagerListTriesArg{dummyAuthArg, c.Internal}
	var reply ManagerListTriesReply
	err = rpcClient.CallAuthed("ListTries", &arg, &reply)
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
	arg := ManagerUpdatePortArg{
		ManagerAuthArg: dummyAuthArg,
		Port: config.Port{
			Port:     c.Port,
			Trie:     c.Trie,
			Internal: c.Internal,
		},
	}
	var reply ManagerUpdatePortReply
	err = rpcClient.CallAuthed("UpdatePort", &arg, &reply)
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
	arg := ManagerDeletePortArg{dummyAuthArg, c.Port, c.Internal}
	var reply ManagerDeletePortReply
	err = rpcClient.CallAuthed("DeletePort", &arg, &reply)
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
	arg := ManagerGetPortArg{dummyAuthArg, c.Port, c.Internal}
	var reply ManagerGetPortReply
	err = rpcClient.CallAuthed("GetPort", &arg, &reply)
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
	arg := ManagerListPortsArg{dummyAuthArg, c.Internal}
	var reply ManagerListPortsReply
	err = rpcClient.CallAuthed("ListPorts", &arg, &reply)
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
	arg := ManagerGetAppEnvPortArg{
		ManagerAuthArg: dummyAuthArg,
		App:            c.App,
		Env:            c.Env,
	}
	var reply ManagerGetAppEnvPortReply
	err = rpcClient.CallAuthed("GetAppEnvPort", &arg, &reply)
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
	arg := ManagerListAppEnvsWithPortArg{dummyAuthArg, c.Internal}
	var reply ManagerListAppEnvsWithPortReply
	err = rpcClient.CallAuthed("ListAppEnvsWithPort", &arg, &reply)
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
