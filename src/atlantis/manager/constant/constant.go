package constant

import (
	"regexp"
)

var (
	Host           = "localhost"
	Region         = "dev"
	Zone           = "dev1"
	AvailableZones = []string{"dev1"}
	AppRegexp      = regexp.MustCompile("[A-Za-z0-9-]+") // apps can contain letters, numbers, and -
)

const (
	ManagerRPCVersion               = "1.0.0"
	ManagerAPIVersion               = "1.0.0"
	DefaultManagerHost              = "localhost"
	DefaultManagerRPCPort           = uint16(1338)
	DefaultManagerAPIPort           = uint16(443)
	DefaultManagerKeyPath           = "~/.manager"
	DefaultResultDuration           = "30m"
	DefaultMaintenanceFile          = "/etc/atlantis/manager/maint"
	DefaultMaintenanceCheckInterval = "5s"
)
