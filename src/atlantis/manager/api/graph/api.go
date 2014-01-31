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
