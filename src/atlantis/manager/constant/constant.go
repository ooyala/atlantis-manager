package constant

var (
	Region         = "dev"
	Zone           = "dev"
	AvailableZones = []string{"dev"}
)

const (
	ManagerRPCVersion               = "0.2.0"
	ManagerAPIVersion               = "0.2.0"
	DefaultManagerRPCPort           = uint16(1338)
	DefaultManagerAPIPort           = uint16(8443)
	DefaultManagerKeyPath           = "~/.manager"
	DefaultResultDuration           = "30m"
	DefaultMaintenanceFile          = "/etc/atlantis/manager/maint"
	DefaultMaintenanceCheckInterval = "5s"
)
