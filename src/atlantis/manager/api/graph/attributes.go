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
