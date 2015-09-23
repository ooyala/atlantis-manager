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

package rpc

import (
	. "atlantis/common"
	. "atlantis/manager/rpc/types"
	"atlantis/manager/datamodel"
	"fmt"
	"strings"
	"encoding/json"
)

type GetAutoScalerRuleExecutor struct {
	arg   ManagerGetAutoScalerRuleArg
	reply *ManagerGetAutoScalerRuleReply
}

func (e *GetAutoScalerRuleExecutor) Request() interface{} {
	return e.arg
}

func (e *GetAutoScalerRuleExecutor) Result() interface{} {
	return e.reply
}

func (e *GetAutoScalerRuleExecutor) Description() string {
	return fmt.Sprintf("["+e.arg.ManagerAuthArg.User+"] %s %s %s", e.arg.App, e.arg.Sha, e.arg.Env)
}

func (e *GetAutoScalerRuleExecutor) Execute(t *Task) (err error) {
        fmt.Println("/atlantis/autoscale-rules/"+ e.arg.App + "-" + e.arg.Sha + "-" + e.arg.Env)
	data, err := datamodel.GetAutoscaleRule("/atlantis/autoscale-rules/"+ e.arg.App + "-" + e.arg.Sha + "-" + e.arg.Env)
	if err != nil {
	      e.reply.Rule = ""
	      e.reply.Status = "NOT OK"
 	} else {
	      e.reply.Rule = data
	      e.reply.Status = "OK"
	}
	fmt.Println(err)
	return err
}

func (e *GetAutoScalerRuleExecutor) Authorize() error {
	return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

func (m *ManagerRPC) GetAutoScalerRule(arg ManagerGetAutoScalerRuleArg, reply *ManagerGetAutoScalerRuleReply) error {
	return NewTask("GetAutoScalerRule", &GetAutoScalerRuleExecutor{arg, reply}).Run()
}


type SetAutoScalerRuleExecutor struct {
        arg   ManagerSetAutoScalerRuleArg
        reply *ManagerSetAutoScalerRuleReply
}

func (e *SetAutoScalerRuleExecutor) Request() interface{} {
        return e.arg
}

func (e *SetAutoScalerRuleExecutor) Result() interface{} {
        return e.reply
}

func (e *SetAutoScalerRuleExecutor) Description() string {
        return fmt.Sprintf("["+e.arg.ManagerAuthArg.User+"] %s %s %s %s", e.arg.App, e.arg.Sha, e.arg.Env, e.arg.Data)
}

func (e *SetAutoScalerRuleExecutor) Execute(t *Task) (err error) {
        fmt.Println("/atlantis/autoscale-rules/"+ e.arg.App + "-" + e.arg.Sha + "-" + e.arg.Env+ "    " + e.arg.Data)


        if strings.Trim(e.arg.Data," ") == "" {
                // delete rules baby
                err = datamodel.DeleteAutoscaleRule("/atlantis/autoscale-rules/"+ e.arg.App + "-" + e.arg.Sha + "-" + e.arg.Env)
                if err != nil {
                        e.reply.Status = "NOT OK"
                } else {
                        e.reply.Status = "OK"
                }
                return err
        }

        var rule AutoscaleRule
        if err := json.Unmarshal([]byte(e.arg.Data), &rule); err != nil {
                e.reply.Status = "NOT OK"
                return err
        }

        err = datamodel.SetAutoscaleRule("/atlantis/autoscale-rules/"+ e.arg.App + "-" + e.arg.Sha + "-" + e.arg.Env, e.arg.Data)
        if err != nil {
              e.reply.Status = "NOT OK"
        } else {
              e.reply.Status = "OK"
        }
        fmt.Println(err)
        return err
}

func (e *SetAutoScalerRuleExecutor) Authorize() error {
        return SimpleAuthorize(&e.arg.ManagerAuthArg)
}

func (m *ManagerRPC) SetAutoScalerRule(arg ManagerSetAutoScalerRuleArg, reply *ManagerSetAutoScalerRuleReply) error {
        return NewTask("SetAutoScalerRule", &SetAutoScalerRuleExecutor{arg, reply}).Run()
}

//copied from autoscaler repo

type AutoscaleRule struct{
	IntervalInSec      string   
	MetricSourceType   string 

	MinInstances       string   
	MaxInstances	   string
	DatadogMetric      DatadogMetricType
	Pool               PoolInfo
	Policies           []ScalingPolicy  
}

type PoolInfo struct {
	AppName            string
	Sha                string
	Env                string
}

type ScalingPolicy struct {
	Min                string
	Max                string
	Change             string
}

type DatadogMetricType struct {
	MetricName         string
	IntervalInSec      string
	Tags               []string
	//StatsMethod         string
}
