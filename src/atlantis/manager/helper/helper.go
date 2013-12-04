package helper

import (
	atlantis "atlantis/common"
	. "atlantis/manager/constant"
	routerzk "atlantis/router/zk"
	"fmt"
	"path"
	"regexp"
	"strings"
)

const contRandIdSize = 6

func CreateContainerId(app, sha, env string) string {
	trimmedSha := sha
	if len(trimmedSha) > 6 {
		trimmedSha = trimmedSha[0:6]
	}
	return fmt.Sprintf("%s-%s-%s-%s", app, trimmedSha, env, atlantis.CreateRandomId(contRandIdSize))
}

func JoinWithBase(base string, args ...string) string {
	if len(args) == 0 {
		return base
	}
	return path.Join(base, path.Join(args...))
}

// get path to node in /apps
func GetBaseAppPath(args ...string) string {
	base := fmt.Sprintf("/atlantis/apps/%s", Region)
	return JoinWithBase(base, args...)
}

// Helper function to get path to node in /instances at any level
func GetBaseInstancePath(args ...string) string {
	base := fmt.Sprintf("/atlantis/instances/%s", Region)
	return JoinWithBase(base, args...)
}

func GetBaseInstanceDataPath(args ...string) string {
	base := fmt.Sprintf("/atlantis/instance_data/%s", Region)
	return JoinWithBase(base, args...)
}

func GetBaseHostPath(args ...string) string {
	base := fmt.Sprintf("/atlantis/hosts/%s", Region)
	return JoinWithBase(base, args...)
}

func CreatePoolName(app, sha, env string) string {
	return fmt.Sprintf("%s-%s-%s", app, sha, env)
}

func SetRouterRoot(internal bool) {
	internalStr := "external"
	if internal {
		internalStr = "internal"
	}
	root := fmt.Sprintf("/atlantis/router/%s/%s", Region, internalStr)
	routerzk.SetZkRoot(root)
}

func GetBaseDNSPath(args ...string) string {
	base := fmt.Sprintf("/atlantis/dns/%s", Region)
	return JoinWithBase(base, args...)
}

func GetBaseRouterPath(internal bool, args ...string) string {
	internalStr := "external"
	if internal {
		internalStr = "internal"
	}
	base := fmt.Sprintf("/atlantis/routers/%s/%s", Region, internalStr)
	return JoinWithBase(base, args...)
}

func GetBaseManagerPath(args ...string) string {
	base := "/atlantis/managers"
	return JoinWithBase(base, args...)
}

func GetManagerCName(num int, region, suffix string) string {
	return fmt.Sprintf("manager%d.%s.%s", num, region, suffix)
}

func GetRegistryCName(num int, region, suffix string) string {
	return fmt.Sprintf("registry%d.%s.%s", num, region, suffix)
}

func GetBaseDepPath(args ...string) string {
	switch len(args) {
	case 0:
		return GetBaseEnvPath()
	case 1:
		return GetBaseEnvPath(args...)
	default:
		return GetBaseEnvPath(args[0]) + fmt.Sprintf("/%s", args[1])
	}
}

func GetBaseEnvPath(args ...string) string {
	base := fmt.Sprintf("/atlantis/environments/%s", Region)
	return JoinWithBase(base, args...)
}

func GetBaseLockPath(args ...string) string {
	base := fmt.Sprintf("/atlantis/lock/%s", Region)
	return JoinWithBase(base, args...)
}

func GetRegionRouterCName(internal bool, suffix string) string {
	internalStr := ""
	if internal {
		internalStr = "internal-"
	}
	return fmt.Sprintf("%srouter.%s.%s", internalStr, Region, suffix)
}

func GetZoneRouterCName(internal bool, zone, suffix string) string {
	internalStr := ""
	if internal {
		internalStr = "internal-"
	}
	return fmt.Sprintf("%srouter.%s.%s", internalStr, RegionAndZone(zone), suffix)
}

func GetRouterCName(internal bool, num int, zone, suffix string) string {
	internalStr := ""
	if internal {
		internalStr = "internal-"
	}
	return fmt.Sprintf("%srouter%d.%s.%s", internalStr, num, RegionAndZone(zone), suffix)
}

func GetRegionAppAlias(app, env, suffix string) string {
	return fmt.Sprintf("%s%s.%s.%s", app, EmptyIfProdPrefix(env), Region, suffix)
}

func GetZoneAppAlias(app, env, zone, suffix string) string {
	return fmt.Sprintf("%s%s.%s.%s", app, EmptyIfProdPrefix(env), RegionAndZone(zone), suffix)
}

func RegionAndZone(zone string) string {
	if strings.Contains(zone, Region) {
		return zone
	}
	return Region + zone
}

var envSuffixRegexp = regexp.MustCompile("^(prod|production)([_-]|$)")

func EmptyIfProdPrefix(env string) string {
	if envSuffixRegexp.MatchString(env) {
		return ""
	}
	return "." + env
}
