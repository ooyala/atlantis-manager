package graph

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

var (
	staticDir     = "static"
	visualizeBase = staticDir + "/router/config"
)

func VisualizeIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, visualizeBase+"/index.html")
}

func TrieJson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, err := DotTrie(vars["name"], true)
	if err != nil {
		str = "{ \"Error\" : \"" + err.Error() + "\"}"
	}
	fmt.Fprintf(w, "%s", str)
}

func TrieDot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, err := DotTrie(vars["name"], false)
	if err != nil {
		str = err.Error()
	}
	fmt.Fprintf(w, "%s", str)
}

func TrieSvg(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, err := DotTrie(vars["name"], false)
	if err != nil {
		fmt.Fprintf(w, "%s", "")
	} else {
		g := "<script type=\"text/vnd.graphviz\" id=\"output\">" + str + "</script>"
		h := "<script>window.onload = function () { var d = document.getElementById(\"output\").innerHTML; d = d + new Array(d.length).join(\" \");  document.body.innerHTML += Viz(d, \"svg\"); }</script>"
		html := "<html><body>" + g + "<script src=\"/static/router/config/viz.js\"></script>" + h + "</body></html>"
		fmt.Fprintf(w, "%s", html)
	}
}
