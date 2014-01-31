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
	"strings"
)

func StrikethroughString(str string) string {
	s := strings.Split(str, "")
	t := "̶"
	for _, v := range s {
		t = t + v + "̶"
	}
	return t
}

func GetColorAttr(path string) (string, string, string) {
	if IsRule(path) {
		return ruleColor, ruleFill, ruleFontColor
	} else if IsTrie(path) {
		return trieColor, trieFill, trieFontColor
	} else if IsHosts(path) {
		return hostColor, hostFill, hostFontColor
	} else if IsPool(path) {
		return poolColor, poolFill, poolFontColor
	}
	return "", "", ""
}

func IsPool(path string) bool {
	return strings.Contains(path, poolPath)
}

func IsHosts(path string) bool {
	a := strings.Split(path, poolPath)
	b := strings.Split(a[1], "/")
	if len(b) > 2 {
		return true
	}
	return false
}

func IsRule(path string) bool {
	return strings.Contains(path, rulePath)
}

func IsTrie(path string) bool {
	return strings.Contains(path, triePath)
}
