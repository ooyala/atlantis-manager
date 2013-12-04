package constant

import (
	"regexp"
)

var (
	Region         = "dev"
	Zone           = "dev"
	AvailableZones = []string{"dev"}
	AppRegexp      = regexp.MustCompile("[A-Za-z0-9-]+") // apps can contain letters, numbers, and -
)

const (
	ManagerRPCVersion               = "0.2.0"
	ManagerAPIVersion               = "0.2.0"
	DefaultManagerRPCPort           = uint16(1338)
	DefaultManagerAPIPort           = uint16(443)
	DefaultManagerKeyPath           = "~/.manager"
	DefaultResultDuration           = "30m"
	DefaultMaintenanceFile          = "/etc/atlantis/manager/maint"
	DefaultMaintenanceCheckInterval = "5s"
)
