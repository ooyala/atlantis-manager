package helper

import (
	atlantis "atlantis/common"
	. "atlantis/manager/constant"
	routerzk "atlantis/router/zk"
	"fmt"
	"path"
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

func GetBaseRouterPath(args ...string) string {
	base := fmt.Sprintf("/atlantis/routers/%s", Region)
	return JoinWithBase(base, args...)
}

func GetBaseManagerPath(args ...string) string {
	base := "/atlantis/managers"
	return JoinWithBase(base, args...)
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

func GetRegionRouterCName(suffix string) string {
	return fmt.Sprintf("router.%s.%s", Region, suffix)
}

func GetZoneRouterCName(zone, suffix string) string {
	return fmt.Sprintf("router.%s.%s", RegionAndZone(zone), suffix)
}

func GetRouterCName(num int, zone, suffix string) string {
	return fmt.Sprintf("router%d.%s.%s", num, RegionAndZone(zone), suffix)
}

func GetRegionAppAlias(app, env, suffix string) string {
	return fmt.Sprintf("%s.%s.%s.%s", app, env, Region, suffix)
}

func GetZoneAppAlias(app, env, zone, suffix string) string {
	return fmt.Sprintf("%s.%s.%s.%s", app, env, RegionAndZone(zone), suffix)
}

func RegionAndZone(zone string) string {
	if strings.Contains(zone, Region) {
		return zone
	}
	return Region + zone
}
