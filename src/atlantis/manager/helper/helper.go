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

package helper

import (
	atlantis "atlantis/common"
	. "atlantis/manager/constant"
	"atlantis/manager/rpc/types"
	routerzk "atlantis/router/zk"
	"fmt"
	"log"
	"path"
	"regexp"
	"strings"
)

const contRandIDSize = 6

func CreateContainerID(app, sha, env string) string {
	trimmedSha := sha
	if len(trimmedSha) > 6 {
		trimmedSha = trimmedSha[0:6]
	}
	return fmt.Sprintf("%s-%s-%s-%s", app, trimmedSha, env, atlantis.CreateRandomID(contRandIDSize))
}

func JoinWithBase(base string, args ...string) string {
	if len(args) == 0 {
		return base
	}
	return path.Join(base, path.Join(args...))
}

func GetBaseIPGroupPath(args ...string) string {
	base := fmt.Sprintf("/atlantis/ip_groups/%s", Region)
	return JoinWithBase(base, args...)
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

func GetBaseSupervisorPath(args ...string) string {
	base := fmt.Sprintf("/atlantis/supervisors/%s", Region)
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

func GetBaseRouterPortsPath(internal bool, args ...string) string {
	internalStr := "external"
	if internal {
		internalStr = "internal"
	}
	base := fmt.Sprintf("/atlantis/router_ports/%s/%s", Region, internalStr)
	return JoinWithBase(base, args...)
}

func GetAppEnvTrieName(app, env string) string {
	return types.AppEnv{App: app, Env: env}.String()
}

func GetAppShaEnvStaticRuleName(app, sha, env string) string {
	return fmt.Sprintf("static-%s-%s-%s", app, sha, env)
}

func GetBaseManagerPath(args ...string) string {
	base := "/atlantis/managers"
	return JoinWithBase(base, args...)
}

func GetManagerCName(num int, suffix string) string {
	return fmt.Sprintf("manager%d.%s", num, suffix)
}

func GetRegistryCName(num int, suffix string) string {
	return fmt.Sprintf("registry%d.%s", num, suffix)
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
	return fmt.Sprintf("%srouter.%s", internalStr, suffix)
}

func GetZoneRouterCName(internal bool, zone, suffix string) string {
	internalStr := ""
	if internal {
		internalStr = "internal-"
	}
	return fmt.Sprintf("%srouter.%s.%s", internalStr, ZoneMinusRegion(zone), suffix)
}

func GetRouterCName(internal bool, num int, zone, suffix string) string {
	internalStr := ""
	if internal {
		internalStr = "internal-"
	}
	return fmt.Sprintf("%srouter%d.%s.%s", internalStr, num, ZoneMinusRegion(zone), suffix)
}

func GetAppCNameSuffixes(suffix string) []string {
	suffixes := []string{suffix}
	for _, zone := range AvailableZones {
		suffixes = append(suffixes, fmt.Sprintf("%s.%s", ZoneMinusRegion(zone), suffix))
	}
	return suffixes
}

func GetRegionAppCName(app, env, suffix string) string {
	return fmt.Sprintf("%s%s.%s", app, EmptyIfProdPrefix(env), suffix)
}

func GetZoneAppCName(app, env, zone, suffix string) string {
	return fmt.Sprintf("%s%s.%s.%s", app, EmptyIfProdPrefix(env), ZoneMinusRegion(zone), suffix)
}

func RegionAndZone(zone string) string {
	if strings.Contains(zone, Region) {
		return zone
	}
	return Region + zone
}

func ZoneMinusRegion(zone string) string {
	zmr := strings.Replace(zone, Region, "", 1)
	if zmr == "" {
		results := strings.Split(zone, "-")
		if len(results) == 2 {
			zmr = results[1]
		}
	}
	log.Printf("ZoneMinusRegion: zone=%s Region=%s zmr='%s'", zone, Region, zmr)
	return zmr
}

var envSuffixRegexp = regexp.MustCompile("^(prod|production)([_-]|$)")

func EmptyIfProdPrefix(env string) string {
	if envSuffixRegexp.MatchString(env) {
		return ""
	}
	return "." + env
}
