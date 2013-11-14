package graph

import (
	. "atlantis/manager/constant"
	"atlantis/manager/datamodel"
	ggv "code.google.com/p/gographviz"
	gozk "launchpad.net/gozk"
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
