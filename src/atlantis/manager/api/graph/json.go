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

package graph

import (
	"atlantis/router/zk"
	"encoding/json"
	"fmt"
	ggv "github.com/awalterschulze/gographviz"
	"strconv"
)

type zkGraph struct {
	Graph  *ggv.Graph
	Config *zCfg
	Edges  *Edges
}

type zCfg struct {
	Pools map[string]Pool
	Rules map[string]Rule
	Tries map[string]Trie
}

type Edges struct {
	Pools map[string]bool
	Rules map[string]bool
	Tries map[string]bool
}

func (ZC *zCfg) ToJSON() (string, error) {
	res, err := json.Marshal(ZC)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

type Pool struct {
	HealthCheckEvery string
	Degraded         string
	Critical         string
	Healthz          string
	Request          string
	Hosts            []string
}

type Rule struct {
	Type  string
	Value string
	Next  string
	Pool  string
}

type Trie struct {
	Rules []string
}

func JSONToArray(data []byte) (map[string]interface{}, error) {
	var f interface{}
	if err := json.Unmarshal(data, &f); err != nil {
		// Assume that we just have a string that is not JSON
		return nil, err
	}
	m := f.(map[string]interface{})
	return m, nil
}

func (ZC *zCfg) ParseWholeTrie(name string) error {
	var err error
	ZC.Tries[name], err = ParseTrie(name)
	if err != nil {
		return err
	}
	for _, v := range ZC.Tries[name].Rules {
		if _, ok := ZC.Rules[v]; !ok {
			ZC.Rules[v], err = ParseRule(v)
			if err != nil {
				fmt.Println("Error Parsing Rule")
				return err
			}
			if ZC.Rules[v].Pool != "" {
				if _, ok := ZC.Pools[ZC.Rules[v].Pool]; !ok {
					ZC.Pools[ZC.Rules[v].Pool], err = ParsePool(ZC.Rules[v].Pool)
					if err != nil {
						fmt.Println("Error Getting Pool")
						return err
					}
				}
			}
			if ZC.Rules[v].Pool == "" && ZC.Rules[v].Next != "" {
				if _, ok := ZC.Tries[ZC.Rules[v].Next]; !ok {
					ZC.Tries[ZC.Rules[v].Next], err = ParseTrie(ZC.Rules[v].Next)
					if err != nil {
						fmt.Println("Error Getitng Next")
						return err
					}
				}
			}
		}
	}
	return nil
}

func (ZC *zCfg) ParseAllTrie() error {
	data, _, err := Zk.Children(triePath)
	if err != nil {
		return err
	}
	for _, v := range data {
		err := ZC.ParseWholeTrie(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func ParseTrie(name string) (Trie, error) {
	var res Trie
	data, err := zk.GetTrie(Zk, name)
	if err != nil {
		return res, err
	}
	return Trie{data.Rules}, nil
}

func ParseRule(name string) (Rule, error) {
	var res Rule
	data, err := zk.GetRule(Zk, name)
	if err != nil {
		return res, err
	}
	return Rule{data.Type, data.Value, data.Next, data.Pool}, nil
}

func ParsePool(name string) (Pool, error) {
	var res Pool
	data, err := zk.GetPool(Zk, name)
	if err != nil {
		return res, err
	}
	numHosts := len(data.Hosts)
	hosts := make([]string, numHosts)
	count := 0
	for _, v := range data.Hosts {
		hosts[count] = v.Address
		count++
	}
	return Pool{data.Config.HealthzEvery, "N/A", "N/A",
		data.Config.HealthzTimeout, data.Config.RequestTimeout, hosts}, nil

}

func (z *zkGraph) JSONtoDot(j string) string {
	poolm := map[string]bool{}
	rulem := map[string]bool{}
	triesm := map[string]bool{}
	var zc *zCfg
	json.Unmarshal([]byte(j), &zc)
	for k, t := range zc.Tries {
		triesm[k] = true
		DrawTrie(z.Graph, k, t)
		z.ProcessTrie(zc, k, t, poolm, rulem, triesm)
	}
	return z.Graph.String()
}

func DrawTrie(graph *ggv.Graph, name string, t Trie) {
	fpath := triePath + "/" + name
	attr := map[string]string{}
	lbl := " | <f1> Rules : "
	for _, v := range t.Rules {
		lbl = lbl + v + ", "
	}
	lbl = lbl[0 : len(lbl)-2]
	attr["shape"] = "record"
	attr["label"] = "\"{ <f0>" + name + lbl + "} \""
	attr["color"], attr["fillcolor"], attr["fontcolor"] = GetColorAttr(fpath)
	attr["fontname"] = font
	attr["style"] = "\"filled\""
	graph.AddNode(graphName, "\""+fpath+"\"", attr)
}

func (z *zkGraph) ProcessTrie(zc *zCfg, name string, t Trie, poolm map[string]bool, rulem map[string]bool, triesm map[string]bool) {
	for _, v := range t.Rules {
		if !rulem[v] {
			c := zc.Rules[v]
			rulem[v] = true
			DrawRule(z.Graph, v, c)
			if c.Pool != "" {
				if !poolm[c.Pool] {
					poolm[c.Pool] = true
					d := zc.Pools[c.Pool]
					DrawPool(z.Graph, c.Pool, d)
					if len(d.Hosts) > 0 {
						DrawHosts(z.Graph, c.Pool, d.Hosts)
					}
				}
				if !z.Edges.Rules["rules_"+v+"_pool_"+c.Pool] {
					z.Graph.AddEdge("\""+rulePath+"/"+v+"\":f4", "\""+poolPath+"/"+c.Pool+"\"", true, RuleEdgeAttr)
					z.Edges.Rules["rules_"+v+"_pool_"+c.Pool] = true
				}
			} else {
				if c.Next != "" {
					if !triesm[c.Next] {
						triesm[c.Next] = true
						DrawTrie(z.Graph, c.Next, zc.Tries[c.Next])
						z.ProcessTrie(zc, c.Next, zc.Tries[c.Next], poolm, rulem, triesm)
					}
					if !z.Edges.Rules["rules_"+v+"_trie_"+c.Next] {
						z.Graph.AddEdge("\""+rulePath+"/"+v+"\":f3", "\""+triePath+"/"+c.Next+"\"", true, RuleEdgeAttr)
						z.Edges.Rules["rules_"+v+"_trie_"+c.Next] = true
					}

				}
			}
		}
		if !z.Edges.Rules["trie_"+name+"_rule_"+v] {
			z.Edges.Rules["trie_"+name+"_rule_"+v] = true
			z.Graph.AddEdge("\""+triePath+"/"+name+"\":f1", "\""+rulePath+"/"+v+"\"", true, EdgeAttr)
		}
	}
}

func DrawRule(graph *ggv.Graph, name string, r Rule) {
	fpath := rulePath + "/" + name
	attr := map[string]string{}
	lbl := " | <f1> Type : " + r.Type + " | <f2> Value : " + r.Value
	if r.Pool != "" {
		lbl = lbl + "| <f3> " + StrikethroughString("Next : "+r.Next) + " | <f4> Pool : " + r.Pool
	} else {
		lbl = lbl + "| <f3> Next : " + r.Next + " | <f4> " + StrikethroughString("Pool : "+r.Pool)
	}
	attr["shape"] = "record"
	attr["label"] = "\"{ <f0>" + name + lbl + "} \""
	attr["color"], attr["fillcolor"], attr["fontcolor"] = GetColorAttr(fpath)
	attr["fontname"] = font
	attr["style"] = "\"filled\""
	graph.AddNode(graphName, "\""+fpath+"\"", attr)
}

func DrawPool(graph *ggv.Graph, name string, p Pool) {
	fpath := poolPath + "/" + name
	attr := map[string]string{}
	lbl := " | <f1> HealthCheckEvery : " + p.HealthCheckEvery + " | <f2> Degraded : " + p.Degraded +
		" | <f3> Critical : " + p.Critical + " | <f4> Healthz : " + p.Healthz + " | <f5> Request : " + p.Request

	attr["shape"] = "record"
	attr["label"] = "\"{ <f0>" + name + lbl + "} \""
	attr["color"], attr["fillcolor"], attr["fontcolor"] = GetColorAttr(fpath)
	attr["fontname"] = font
	attr["style"] = "\"filled\""
	graph.AddNode(graphName, "\""+fpath+"\"", attr)
}

func DrawHosts(graph *ggv.Graph, parentName string, hosts []string) {
	fpath := poolPath + "/" + parentName + "/hosts"
	attr := map[string]string{}
	lbl := ""
	i := 1
	for _, v := range hosts {
		lbl = lbl + "| <f" + strconv.FormatInt(int64(i), 10) + "> " + v + " "
	}
	attr["shape"] = "record"
	attr["label"] = "\"{ <f0> Hosts " + lbl + "} \""
	attr["color"], attr["fillcolor"], attr["fontcolor"] = GetColorAttr(fpath)
	attr["fontname"] = font
	attr["style"] = "\"filled\""
	graph.AddNode(graphName, "\""+fpath+"\"", attr)
	graph.AddEdge("\""+poolPath+"/"+parentName+"\"", "\""+fpath+"\"", true, PoolEdgeAttr)
}
