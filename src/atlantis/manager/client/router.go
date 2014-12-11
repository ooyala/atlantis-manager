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
	Arg      ManagerDeletePoolArg
	Reply    ManagerDeletePoolReply
}

func (c *DeletePoolCommand) Execute(args []string) error {
	return genericExecuter(c, args)
}

type GetPoolCommand struct {
	Name     string `short:"n" long:"name" description:"the name of the pool"`
	Internal bool   `short:"i" long:"internal" description:"true if internal"`
	Arg      ManagerGetPoolArg
	Reply    ManagerGetPoolReply
}

func (c *GetPoolCommand) Execute(args []string) error {
	return genericExecuter(c, args)
}

type ListPoolsCommand struct {
	Internal bool `short:"i" long:"internal" description:"true if internal"`
	Arg      ManagerListPoolsArg
	Reply    ManagerListPoolsReply
}

func (c *ListPoolsCommand) Execute(args []string) error {
	return genericExecuter(c, args)
}

type UpdateRuleCommand struct {
	Name     string `short:"n" long:"name" description:"the name of the rule"`
	Type     string `short:"t" long:"type" description:"the type of the rule"`
	Value    string `short:"v" long:"value" description:"the rule's value"`
	Next     string `short:"x" long:"next" description:"the next ruleset"`
	Pool     string `short:"p" long:"pool" description:"the pool to point to if this rule succeeds"`
	Internal bool   `short:"i" long:"internal" description:"true if internal"`
	Arg      ManagerUpdateRuleArg
	Reply    ManagerUpdateRuleReply
}

func (c *UpdateRuleCommand) Execute(args []string) error {
	return genericExecuter(c, args)
}

type DeleteRuleCommand struct {
	Name     string `short:"n" long:"name" description:"the name of the rule"`
	Internal bool   `short:"i" long:"internal" description:"true if internal"`
	Arg      ManagerDeleteRuleArg
	Reply    ManagerDeleteRuleReply
}

func (c *DeleteRuleCommand) Execute(args []string) error {
	return genericExecuter(c, args)
}

type GetRuleCommand struct {
	Name     string `short:"n" long:"name" description:"the name of the rule"`
	Internal bool   `short:"i" long:"internal" description:"true if internal"`
	Arg      ManagerGetRuleArg
	Reply    ManagerGetRuleReply
}

func (c *GetRuleCommand) Execute(args []string) error {
	return genericExecuter(c, args)
}

type ListRulesCommand struct {
	Internal bool `short:"i" long:"internal" description:"true if internal"`
	Arg      ManagerListRulesArg
	Reply    ManagerListRulesReply
}

func (c *ListRulesCommand) Execute(args []string) error {
	return genericExecuter(c, args)
}

type UpdateTrieCommand struct {
	Name     string   `short:"n" long:"name" description:"the name of the rule"`
	Rules    []string `short:"r" long:"rule" description:"the rules that make up the ruleset"`
	Internal bool     `short:"i" long:"internal" description:"true if internal"`
	Arg      ManagerUpdateTrieArg
	Reply    ManagerUpdateTrieReply
}

func (c *UpdateTrieCommand) Execute(args []string) error {
	return genericExecuter(c, args)
}

type DeleteTrieCommand struct {
	Name     string `short:"n" long:"name" description:"the name of the trie"`
	Internal bool   `short:"i" long:"internal" description:"true if internal"`
	Arg      ManagerDeleteTrieArg
	Reply    ManagerDeleteTrieReply
}

func (c *DeleteTrieCommand) Execute(args []string) error {
	return genericExecuter(c, args)
}

type GetTrieCommand struct {
	Name     string `short:"n" long:"name" description:"the name of the trie"`
	Internal bool   `short:"i" long:"internal" description:"true if internal"`
	Arg      ManagerGetTrieArg
	Reply    ManagerGetTrieReply
}

func (c *GetTrieCommand) Execute(args []string) error {
	return genericExecuter(c, args)
}

type ListTriesCommand struct {
	Internal bool `short:"i" long:"internal" description:"true if internal"`
	Arg      ManagerListTriesArg
	Reply    ManagerListTriesReply
}

func (c *ListTriesCommand) Execute(args []string) error {
	return genericExecuter(c, args)
}

type UpdatePortCommand struct {
	Port     uint16 `short:"p" long:"port" description:"the actual port to listen on"`
	Trie     string `short:"t" long:"trie" description:"the trie to use as root for this port"`
	Internal bool   `short:"i" long:"internal" description:"true if internal"`
	Arg      ManagerUpdatePortArg
	Reply    ManagerUpdatePortReply
}

func (c *UpdatePortCommand) Execute(args []string) error {
	return genericExecuter(c, args)
}

type DeletePortCommand struct {
	Port     uint16 `short:"p" long:"port" description:"the port number"`
	Internal bool   `short:"i" long:"internal" description:"true if internal"`
	Arg      ManagerDeletePortArg
	Reply    ManagerDeletePortReply
}

func (c *DeletePortCommand) Execute(args []string) error {
	return genericExecuter(c, args)
}

type GetPortCommand struct {
	Port     uint16 `short:"p" long:"port" description:"the port number"`
	Internal bool   `short:"i" long:"internal" description:"true if internal"`
	Arg      ManagerGetPortArg
	Reply    ManagerGetPortReply
}

func (c *GetPortCommand) Execute(args []string) error {
	return genericExecuter(c, args)
}

type ListPortsCommand struct {
	Internal bool `short:"i" long:"internal" description:"true if internal"`
	Arg      ManagerListPortsArg
	Reply    ManagerListPortsReply
}

func (c *ListPortsCommand) Execute(args []string) error {
	return genericExecuter(c, args)
}

type GetAppEnvPortCommand struct {
	App        string `short:"a" long:"app" description:"the app of the port"`
	Env        string `short:"e" long:"env" description:"the env of the port"`
	Internal   bool   `short:"i" long:"internal" description:"true if internal"`
	Properties string `message:"Get AppEnv Port"`
	Arg        ManagerGetAppEnvPortArg
	Reply      ManagerGetAppEnvPortReply
}

func (c *GetAppEnvPortCommand) Execute(args []string) error {
	return genericExecuter(c, args)
}

type ListAppEnvsWithPortCommand struct {
	Internal   bool   `short:"i" long:"internal" description:"true if internal"`
	Properties string `message:"List AppEnvs With Ports" field:"AppEnvs" name:"app+envs"`
	Arg        ManagerListAppEnvsWithPortArg
	Reply      ManagerListAppEnvsWithPortReply
}

func (c *ListAppEnvsWithPortCommand) Execute(args []string) error {
	return genericExecuter(c, args)
}
