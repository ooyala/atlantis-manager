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

package datamodel

import (
	"atlantis/manager/helper"
	"atlantis/manager/rpc/types"
	"errors"
	"fmt"
	"log"
)

type ZkIPGroup types.IPGroup

type IpGroup struct {
	Name string 	`db:"name"`
}

type IpGroupMember struct {
	Id int64	`db"id"`
	IpGroup string	`db:"ipgroup"`
	IP string	`db:"ip"`
}


func GetIPGroup(name string) (zig *ZkIPGroup, err error) {
	zig = &ZkIPGroup{}
	err = getJson(helper.GetBaseIPGroupPath(name), zig)

	//////////////////// SQL ///////////////////////////	
	var ips []string
	_, err = DbMap.Select(&ips, "select ip from ipgroupmember where ipgroup=?", name)
	if err != nil {
		fmt.Printf("\n Error selecting IPs, %v \n", err)
		zig = nil
		err := errors.New("No ipgroup found")
	}
	fmt.Printf("\n Found %d for ipgroup %s \n", len(ips), name)
	//////////////////////////////////////////////////
	
	return
}

func (zig *ZkIPGroup) Delete() error {

	/////////////////// SQL ////////////////////////
	_, err := DbMap.Exec("delete from ipgroupmember where ipgroup=?", zig.Name)
	if err != nil {

	}
	_, err = DbMap.Delete(&IpGroup{zig.Name})
	if err != nil {

	}
	//////////////////////////////////////////////

	return Zk.RecursiveDelete(zig.path())
}

func (zig *ZkIPGroup) path() string {
	return helper.GetBaseIPGroupPath(zig.Name)
}

func (zig *ZkIPGroup) Save() error {


	//////////////////// SQL ///////////////////////////
	if zig.Name != "" { 
		fmt.Printf("trying to save ipgroup: %v \n\n", zig)
		obj, err := DbMap.Get(IpGroup{}, zig.Name)
		if err != nil {
			fmt.Printf("Err: %v \n", err)
		}	
		if obj == nil {
			fmt.Printf("Ipgroup not exist yet, insert it")
			ipg := IpGroup{zig.Name}
			DbMap.Insert(&ipg)
		} else {
			fmt.Printf("IpGroup does exist, delete its ipgroupmembers so we can reload")
			ipg := obj.(*IpGroup)
			if ipg != nil {
				fmt.Printf("couldn't cast ipgroup")
			}
			_, err := DbMap.Exec("delete from ipgroupmember where ipgroup=?", zig.Name)
			if err != nil {
				fmt.Printf("error deleting its ipgroupmems")
			}
		}
		//populate ipgroupmem table 
		for _, ip := range zig.IPs {
			ipMem := IpGroupMember{IpGroup: zig.Name, IP: ip}
			DbMap.Insert(&ipMem)	
		}
	} else {
		fmt.Println("\n TRYING TO SAVE IPGROUP WITH EMPTY NAME \n")
	}
	///////////////////////////////////////////////////

	if err := setJson(zig.path(), zig); err != nil {
		return err
	}
	

	
	return nil
}

func ListIPGroups() (groups []string, err error) {
	groups, _, err = Zk.VisibleChildren(helper.GetBaseIPGroupPath())
	if err != nil {
		log.Printf("Error getting list of ip groups. Error: %s.", err.Error())
	}
	if groups == nil {
		log.Println("No ip groups found. Returning empty list.")
		groups = []string{}
	}
	/////////////////////// SQL //////////////////////////
	var igroups []string
	_, err = DbMap.Select(&igroups, "select name from ipgroup")
	if err != nil {

	}	
	/////////////////////////////////////////////////////

	return
}
