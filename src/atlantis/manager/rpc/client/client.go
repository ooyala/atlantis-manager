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
	atlantis "atlantis/common"
	. "atlantis/manager/constant"
	"reflect"
)

type ManagerRPCClient struct {
	atlantis.RPCClient
	User    string
	Secrets map[string]string
}

type authedArg interface {
	SetCredentials(string, string)
}

func (r *ManagerRPCClient) CallAuthed(name string, arg authedArg, reply interface{}) error {
	_, err := r.CallAuthedMulti(name, arg, reply)
	return err
}


/* This should be the only Call.  But in the interest of not having to change every request all at the same
* time... */
func (r *ManagerRPCClient) CallAuthedMulti(name string, arg authedArg, reply interface{}) (map[string]interface{}, error) {
	/* This is a terrible hack, but any other fix I can think of either requires changing a commonly used
	 * interface, and thus breaks a whole bunch of other code, or is a worse hack.  */
	replies := map[string]interface{}{}
	originalOpts := r.Opts
	for _, opt := range originalOpts {
		arg.SetCredentials(r.User, r.Secrets[opt.RPCHostAndPort()])
		err := r.RPCClient.Call(name, arg, reply)
		if err != nil {
			return replies, err
		}
		// NOTE(edanaher): reply is a pointer.  We need to copy it.  This apparently requires reflection.
		replies[opt.RPCHostAndPort()] = reflect.ValueOf(reply).Elem().Interface()
		r.Opts = r.Opts[1:]
	}
	r.Opts = originalOpts
	return replies, nil
}

func NewManagerRPCClient(hostAndPort string) *atlantis.RPCClient {
	return atlantis.NewRPCClient(hostAndPort, "ManagerRPC", ManagerRPCVersion, true)
}

func NewManagerRPCClientWithConfig(cfg []atlantis.RPCServerOpts) *atlantis.RPCClient {
	return atlantis.NewMultiRPCClientWithConfig(cfg, "ManagerRPC", ManagerRPCVersion, true)
}
