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

package constant

import (
	"regexp"
)

var (
	Host                = "localhost"
	Region              = "dev"
	Zone                = "dev1"
	AvailableZones      = []string{"dev1"}
	AppRegexp           = regexp.MustCompile("[A-Za-z0-9-]+")                            // apps can contain letters, numbers, and -
	SecurityGroupRegexp = regexp.MustCompile("[0-9]+\\.[0-9]+\\.[0-9]+\\.[0-9]+:[0-9]+") // ip:port
)

const (
	ManagerRPCVersion                 = "3.0.0"
	ManagerAPIVersion                 = "3.0.0"
	DefaultManagerHost                = "localhost"
	DefaultManagerRPCPort             = uint16(1338)
	DefaultManagerAPIPort             = uint16(443)
	DefaultManagerKeyPath             = "~/.manager"
	DefaultResultDuration             = "30m"
	DefaultMaintenanceFile            = "/etc/atlantis/manager/maint"
	DefaultMaintenanceCheckInterval   = "5s"
	DefaultSuperUserOnlyFile          = "/etc/atlantis/manager/superuser"
	DefaultSuperUserOnlyCheckInterval = "5s"
	DefaultMinRouterPort              = uint16(49152)
	DefaultMaxRouterPort              = uint16(65535)
)
