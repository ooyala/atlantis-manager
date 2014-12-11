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
)

type UpdateDepCommand struct {
	Name       string `short:"n" long:"name" description:"the name of the dependency"`
	Value      string `short:"v" long:"value" description:"the value of the dependency"`
	Env        string `short:"e" long:"env" description:"the environment of the dependency"`
	Properties string `field:"Value" name:"value"`
	Arg        ManagerDepArg
	Reply      ManagerDepReply
}

func (c *UpdateDepCommand) Execute(args []string) error {
	return genericExecuter(c, args)
}

type ResolveDepsCommand struct {
	App      string   `short:"a" long:"app" description:"the app the resolve dependencies for"`
	Env      string   `short:"e" long:"env" description:"the environment of the dependencies to resolve"`
	DepNames []string `short:"d" long:"dep" description:"the dep names to resolve"`
	Arg      ManagerResolveDepsArg
	Reply    ManagerResolveDepsReply
}

func (c *ResolveDepsCommand) Execute(args []string) error {
	return genericExecuter(c, args)
}

type GetDepCommand struct {
	Name       string `short:"n" long:"name" description:"the name of the dependency"`
	Env        string `short:"e" long:"env" description:"the environment of the dependency"`
	Properties string `field:"Value" name:"value"`
	Arg        ManagerDepArg
	Reply      ManagerDepReply
}

func (c *GetDepCommand) Execute(args []string) error {
	return genericExecuter(c, args)
}

type DeleteDepCommand struct {
	Name       string `short:"n" long:"name" description:"the name of the dependency"`
	Env        string `short:"e" long:"env" description:"the environment of the dependency"`
	Properties string `field:"Value" name:"value"`
	Arg        ManagerDepArg
	Reply      ManagerDepReply
}

func (c *DeleteDepCommand) Execute(args []string) error {
	return genericExecuter(c, args)
}

type UpdateEnvCommand struct {
	Name  string `short:"n" long:"name" description:"the name of the environment"`
	Arg   ManagerEnvArg
	Reply ManagerEnvReply
}

func (c *UpdateEnvCommand) Execute(args []string) error {
	return genericExecuter(c, args)
}

type DeleteEnvCommand struct {
	Name  string `short:"n" long:"name" description:"the name of the environment"`
	Arg   ManagerEnvArg
	Reply ManagerEnvReply
}

func (c *DeleteEnvCommand) Execute(args []string) error {
	return genericExecuter(c, args)
}
