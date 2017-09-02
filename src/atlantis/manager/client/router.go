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
	"errors"
	"strings"
)

type CreatePoolCommand struct {
	Name             string   `short:"n" long:"name" description:"the name of the pool"`
	HealthCheckEvery string   `short:"e" long:"check-every" default:"5s" description:"how often to check healthz"`
	HealthzTimeout   string   `short:"z" long:"healthz-timeout" default:"5s" description:"timeout for healthz checks"`
	RequestTimeout   string   `short:"r" long:"request-timeout" default:"120s" description:"timeout for requests"`
	Status           string   `short:"s" long:"status" default:"OK" description:"the pool's status"`
	Hosts            []string `short:"H" long:"host" description:"the pool's hosts"`
	Internal         bool     `short:"i" long:"internal" description:"true if internal"`
	Headers          []string `short:"R" long:"response-headers" description:"response headers to set in case apps are not reachable"`
}

func (c *CreatePoolCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Create Pool...")
	if c.Name == "" {
		return OutputError(errors.New("Please specify pool name"))
	}
	hosts := make(map[string]config.Host, len(c.Hosts))
	for _, host := range c.Hosts {
		hosts[host] = config.Host{Address: host}
	}
	var headers []config.HttpHeader
	for _, h := range c.Headers {
		tokens := strings.Split(h, ":")
		if len(tokens) != 2 || tokens[0] == "" || tokens[1] == "" {
			return OutputError(errors.New("Invalid Header " + h))
		}
		hdr := config.HttpHeader{Key: tokens[0], Value: strings.TrimLeft(tokens[1], " ")}
		headers = append(headers, hdr)
	}
	arg := ManagerUpdatePoolArg{dummyAuthArg, config.Pool{Name: c.Name, Hosts: hosts, Internal: c.Internal, Headers: headers,
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

type UpdatePoolCommand struct {
	Name             string   `short:"n" long:"name" description:"the name of the pool"`
	HealthCheckEvery string   `short:"e" long:"check-every"  description:"how often to check healthz"`
	HealthzTimeout   string   `short:"z" long:"healthz-timeout" description:"timeout for healthz checks"`
	RequestTimeout   string   `short:"r" long:"request-timeout"  description:"timeout for requests"`
	Status           string   `short:"s" long:"status"  description:"the pool's status"`
	Ahosts           []string `short:"A" long:"add-host" description:"hosts to add to the pool"`
	Dhosts           []string `short:"D" long:"delete-host" description:"hosts to delete from the pool"`
	Internal         bool     `short:"i" long:"internal" description:"true if internal"`
	AHeaders         []string `short:"a" long:"add-response-headers" description:"response headers to add in case apps are not reachable"`
	DHeaders         []string `short:"d" long:"delete-response-headers" description:"response headers to delete in case apps are not reachable"`
}

func (c *UpdatePoolCommand) Execute(args []string) error {
	err := Init()
	if err != nil {
		return OutputError(err)
	}
	Log("Update Pool...")

	if c.Name == "" {
		return OutputError(errors.New("Please specify pool name"))
	}
	hosts := make(map[string]config.Host)
	var httpheaders []config.HttpHeader
	getPoolArg := ManagerGetPoolArg{dummyAuthArg, c.Name, c.Internal}
	var getPoolReply ManagerGetPoolReply
	err = rpcClient.CallAuthed("GetPool", &getPoolArg, &getPoolReply)
	if err != nil {
		return OutputError(err)
	}

	if getPoolReply.Status != "OK" {
		return Output(map[string]interface{}{"status": getPoolReply.Status}, nil, nil)
	}
	if c.HealthCheckEvery == "" {
		c.HealthCheckEvery = getPoolReply.Pool.Config.HealthzEvery
	}
	if c.HealthzTimeout == "" {
		c.HealthzTimeout = getPoolReply.Pool.Config.HealthzTimeout
	}
	if c.RequestTimeout == "" {
		c.RequestTimeout = getPoolReply.Pool.Config.RequestTimeout
	}
	if c.Status == "" {
		c.Status = getPoolReply.Pool.Config.Status
	}

	for _, element := range getPoolReply.Pool.Hosts {
		hosts[element.Address] = config.Host{Address: element.Address}
	}
	for _, element := range c.Dhosts {
		if _, ok := hosts[element]; ok {
			delete(hosts, element)
		}
	}
	for _, element := range c.Ahosts {
		hosts[element] = config.Host{Address: element}
	}

	httpheaders = make([]config.HttpHeader, len(getPoolReply.Pool.Headers))
	copy(httpheaders, getPoolReply.Pool.Headers)
	for _, h := range c.DHeaders {
		tokens := strings.Split(h, ":")
		if len(tokens) != 2 || tokens[0] == "" || tokens[1] == "" {
			return OutputError(errors.New("Invalid Header " + h))
		}
		hdr := config.HttpHeader{Key: tokens[0], Value: strings.TrimLeft(tokens[1], " ")}
		for i, v := range httpheaders {
			if v.Equals(hdr) {
				httpheaders[i] = httpheaders[len(httpheaders)-1]
				httpheaders = httpheaders[:len(httpheaders)-1]
				break
			}
		}
	}
	for _, h := range c.AHeaders {
		tokens := strings.Split(h, ":")
		if len(tokens) != 2 || tokens[0] == "" || tokens[1] == "" {
			return OutputError(errors.New("Invalid Header " + h))
		}
		hdr := config.HttpHeader{Key: tokens[0], Value: strings.TrimLeft(tokens[1], " ")}
		for i, v := range httpheaders {
			if v.Equals(hdr) {
				httpheaders[i] = httpheaders[len(httpheaders)-1]
				httpheaders = httpheaders[:len(httpheaders)-1]
				break
			}
		}
		httpheaders = append(httpheaders, hdr)
	}

	arg := ManagerUpdatePoolArg{dummyAuthArg, config.Pool{Name: c.Name, Hosts: hosts, Internal: c.Internal, Headers: httpheaders,
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

type GetPoolCommand struct {
	Name     string `short:"n" long:"name" description:"the name of the pool"`
	Internal bool   `short:"i" long:"internal" description:"true if internal"`
	Arg      ManagerGetPoolArg
	Reply    ManagerGetPoolReply
}

type ListPoolsCommand struct {
	Internal bool `short:"i" long:"internal" description:"true if internal"`
	Arg      ManagerListPoolsArg
	Reply    ManagerListPoolsReply
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
        err := Init()
        if err != nil {
                return OutputError(err)
        }
        Log("Update Rule...")
        arg := ManagerUpdateRuleArg{dummyAuthArg, config.Rule{Name: c.Name, Type: c.Type, Value: c.Value, 
					Next: c.Next, Pool: c.Pool, Internal: c.Internal}}
        var reply ManagerUpdateRuleReply
        err = rpcClient.CallAuthed("UpdateRule", &arg, &reply)
        if err != nil {
                return OutputError(err)
        }
        Log("-> status: %s", reply.Status)
        return Output(map[string]interface{}{"status": reply.Status}, nil, nil)
}

type DeleteRuleCommand struct {
	Name     string `short:"n" long:"name" description:"the name of the rule"`
	Internal bool   `short:"i" long:"internal" description:"true if internal"`
	Arg      ManagerDeleteRuleArg
	Reply    ManagerDeleteRuleReply
}

type GetRuleCommand struct {
	Name     string `short:"n" long:"name" description:"the name of the rule"`
	Internal bool   `short:"i" long:"internal" description:"true if internal"`
	Arg      ManagerGetRuleArg
	Reply    ManagerGetRuleReply
}

type ListRulesCommand struct {
	Internal bool `short:"i" long:"internal" description:"true if internal"`
	Arg      ManagerListRulesArg
	Reply    ManagerListRulesReply
}

type UpdateTrieCommand struct {
	Name     string   `short:"n" long:"name" description:"the name of the rule"`
	Rules    []string `short:"r" long:"rule" description:"the rules that make up the ruleset"`
	Internal bool     `short:"i" long:"internal" description:"true if internal"`
	Arg      ManagerUpdateTrieArg
	Reply    ManagerUpdateTrieReply
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
	Arg      ManagerDeleteTrieArg
	Reply    ManagerDeleteTrieReply
}

type GetTrieCommand struct {
	Name     string `short:"n" long:"name" description:"the name of the trie"`
	Internal bool   `short:"i" long:"internal" description:"true if internal"`
	Arg      ManagerGetTrieArg
	Reply    ManagerGetTrieReply
}

type ListTriesCommand struct {
	Internal bool `short:"i" long:"internal" description:"true if internal"`
	Arg      ManagerListTriesArg
	Reply    ManagerListTriesReply
}

type UpdatePortCommand struct {
	Port     uint16 `short:"p" long:"port" description:"the actual port to listen on"`
	Trie     string `short:"t" long:"trie" description:"the trie to use as root for this port"`
	Internal bool   `short:"i" long:"internal" description:"true if internal"`
	Arg      ManagerUpdatePortArg
	Reply    ManagerUpdatePortReply
}

type DeletePortCommand struct {
	Port     uint16 `short:"p" long:"port" description:"the port number"`
	Internal bool   `short:"i" long:"internal" description:"true if internal"`
	Arg      ManagerDeletePortArg
	Reply    ManagerDeletePortReply
}

type GetPortCommand struct {
	Port     uint16 `short:"p" long:"port" description:"the port number"`
	Internal bool   `short:"i" long:"internal" description:"true if internal"`
	Arg      ManagerGetPortArg
	Reply    ManagerGetPortReply
}

type ListPortsCommand struct {
	Internal bool `short:"i" long:"internal" description:"true if internal"`
	Arg      ManagerListPortsArg
	Reply    ManagerListPortsReply
}

type GetAppEnvPortCommand struct {
	App        string `short:"a" long:"app" description:"the app of the port"`
	Env        string `short:"e" long:"env" description:"the env of the port"`
	Internal   bool   `short:"i" long:"internal" description:"true if internal"`
	Properties string `message:"Get AppEnv Port"`
	Arg        ManagerGetAppEnvPortArg
	Reply      ManagerGetAppEnvPortReply
}

type ListAppEnvsWithPortCommand struct {
	Internal   bool   `short:"i" long:"internal" description:"true if internal"`
	Properties string `message:"List AppEnvs With Ports" field:"AppEnvs" name:"app+envs"`
	Arg        ManagerListAppEnvsWithPortArg
	Reply      ManagerListAppEnvsWithPortReply
}
