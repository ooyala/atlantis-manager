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
	. "atlantis/manager/constant"
	"atlantis/manager/datamodel"
	ggv "code.google.com/p/gographviz"
	gozk "github.com/scalingdata/gozk"
)

var (
	Zk            *gozk.Conn
	graphName     = "router"
	rootName      = "/atlantis/router/" + Region + "/external"
	rulePath      = rootName + "/rules"
	poolPath      = rootName + "/pools"
	triePath      = rootName + "/tries"
	font          = "Tahoma"
	ruleColor     = "navy"
	ruleFill      = "skyblue"
	ruleFontColor = "black"
	trieColor     = "red3"
	trieFill      = "mistyrose1"
	trieFontColor = "black"
	poolColor     = "darkgoldenrod4"
	poolFill      = "wheat1"
	poolFontColor = "black"
	hostColor     = "purple"
	hostFill      = "lavender"
	hostFontColor = "black"
	EdgeAttr      = map[string]string{"color": trieColor}
	PoolEdgeAttr  = map[string]string{"color": poolColor}
	NextEdgeAttr  = map[string]string{"color": ruleColor}
	RuleEdgeAttr  = map[string]string{"color": ruleColor}
	HostEdgeAttr  = map[string]string{"color": hostColor}
)

func DotTrie(name string, json bool) (string, error) {
	var err error
	Zk = datamodel.Zk.Conn
	ZC := zCfg{map[string]Pool{}, map[string]Rule{}, map[string]Trie{}}
	Edges := Edges{map[string]bool{}, map[string]bool{}, map[string]bool{}}
	graph := ggv.NewGraph()
	graph.SetName(graphName)
	graph.SetDir(true)
	result := zkGraph{graph, &ZC, &Edges}
	if name == "all" {
		err = ZC.ParseAllTrie()
	} else {
		err = ZC.ParseWholeTrie(name)
	}
	if err != nil {
		return "{}", err
	}
	str, err := ZC.ToJSON()
	if err != nil {
		return "{}", err
	}
	if json {
		return str, nil
	}
	str = result.JSONtoDot(str)
	return str, nil
}
