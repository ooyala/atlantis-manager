package client
/*
import (
	"atlantis/manager/dns"
)

func InitDNSProvider(provider, zone string, ttl uint) {
	InitNoLogin()
	var err error
	switch provider {
	case "route53":
		dns.Provider, err = dns.NewRoute53Provider(zone, ttl)
		if err != nil {
			Output(nil, nil, err)
		}
	default:
		dns.Provider = nil
	}
}

type DNSCreateARecordCommand struct {
	Provider      string `long:"provider" default:"route53" description:"the dns provider"`
	ZoneId        string `short:"z" long:"zone" description:"the dns zone to use"`
	TTL           uint   `long:"ttl" default:"10" description:"the ttl to use"`
	Prefix        string `short:"p" long:"prefix" description:"the name prefix to use"`
	IP            string `short:"i" long:"ip" description:"the ip to use"`
	HealthCheckId string `short:"H" long:"health-check-id" description:"the health check id to use"`
	Failover      string `short:"f" long:"failover" description:"the failover policy to use"`
	Weight        uint8  `short:"w" long:"weight" description:"the record's weight"`
	Comment       string `long:"comment" default:"CLIENT" description:"the comment to use"`
}

func (c *DNSCreateARecordCommand) Execute(args []string) error {
	InitDNSProvider(c.Provider, c.ZoneId, c.TTL)
	arecord := dns.ARecord{
		Name:          c.Prefix + "." + dns.Provider.Suffix(),
		IP:            c.IP,
		HealthCheckId: c.HealthCheckId,
		Failover:      c.Failover,
		Weight:        c.Weight,
	}
	err, errChan := dns.Provider.CreateARecords(c.Comment, []dns.ARecord{arecord})
	if err != nil {
		return Output(nil, nil, err)
	}
	err = <-errChan
	if err == nil {
		Log("-> created %s", arecord.Id())
	}
	return Output(map[string]interface{}{"id": arecord.Id()}, arecord.Id(), err)
}

type DNSCreateAliasCommand struct {
	Provider string `long:"provider" default:"route53" description:"the dns provider"`
	ZoneId   string `short:"z" long:"zone" description:"the dns zone to use"`
	TTL      uint   `long:"ttl" default:"10" description:"the ttl to use"`
	Prefix   string `short:"p" long:"prefix" description:"the name prefix to use for the alias"`
	Original string `short:"o" long:"original" description:"the target of the alias"`
	Failover string `short:"f" long:"failover" description:"the failover policy to use"`
	Comment  string `long:"comment" default:"CLIENT" description:"the comment to use"`
}

func (c *DNSCreateAliasCommand) Execute(args []string) error {
	InitDNSProvider(c.Provider, c.ZoneId, c.TTL)
	alias := dns.Alias{
		Alias:    c.Prefix + "." + dns.Provider.Suffix(),
		Original: c.Original,
		Failover: c.Failover,
	}
	err, errChan := dns.Provider.CreateAliases(c.Comment, []dns.Alias{alias})
	if err != nil {
		return Output(nil, nil, err)
	}
	err = <-errChan
	if err == nil {
		Log("-> created %s", alias.Id())
	}
	return Output(map[string]interface{}{"id": alias.Id()}, alias.Id(), err)
}

type DNSCreateHealthCheckCommand struct {
	Provider string `long:"provider" default:"route53" description:"the dns provider"`
	ZoneId   string `short:"z" long:"zone" description:"the dns zone to use"`
	TTL      uint   `long:"ttl" default:"10" description:"the ttl to use"`
	IP       string `short:"i" long:"ip" description:"the ip to use"`
	Port     uint16 `short:"p" long:"port" description:"the port to check"`
}

func (c *DNSCreateHealthCheckCommand) Execute(args []string) error {
	InitDNSProvider(c.Provider, c.ZoneId, c.TTL)
	id, err := dns.Provider.CreateHealthCheck(c.IP, c.Port)
	if err == nil {
		Log("-> created %s", id)
	}
	return Output(map[string]interface{}{"id": id}, id, err)
}

type DNSDeleteHealthCheckCommand struct {
	Provider      string `long:"provider" default:"route53" description:"the dns provider"`
	ZoneId        string `short:"z" long:"zone" description:"the dns zone to use"`
	TTL           uint   `long:"ttl" default:"10" description:"the ttl to use"`
	HealthCheckId string `short:"i" long:"id" description:"the health check id to use"`
}

func (c *DNSDeleteHealthCheckCommand) Execute(args []string) error {
	InitDNSProvider(c.Provider, c.ZoneId, c.TTL)
	err := dns.Provider.DeleteHealthCheck(c.HealthCheckId)
	if err == nil {
		Log("-> deleted %s", c.HealthCheckId)
	}
	return Output(map[string]interface{}{"id": c.HealthCheckId}, c.HealthCheckId, err)
}

type DNSDeleteRecordsCommand struct {
	Provider  string   `long:"provider" default:"route53" description:"the dns provider"`
	ZoneId    string   `short:"z" long:"zone" description:"the dns zone to use"`
	TTL       uint     `long:"ttl" default:"10" description:"the ttl to use"`
	RecordIds []string `short:"i" long:"id" description:"the record ids to delete"`
	Comment   string   `long:"comment" default:"CLIENT" description:"the comment to use"`
}

func (c *DNSDeleteRecordsCommand) Execute(args []string) error {
	InitDNSProvider(c.Provider, c.ZoneId, c.TTL)
	err, errChan := dns.Provider.DeleteRecords(c.Comment, c.RecordIds...)
	if err != nil {
		return Output(nil, nil, err)
	}
	err = <-errChan
	if err == nil {
		Log("-> deleted:", c.RecordIds)
		for _, id := range c.RecordIds {
			Log("->   %s", id)
		}
	}
	return Output(map[string]interface{}{"ids": c.RecordIds}, c.RecordIds, err)
}

type DNSGetRecordsForIPCommand struct {
	Provider string `long:"provider" default:"route53" description:"the dns provider"`
	ZoneId   string `short:"z" long:"zone" description:"the dns zone to use"`
	TTL      uint   `long:"ttl" default:"10" description:"the ttl to use"`
	IP       string `short:"i" long:"ip" description:"the ip to use"`
}

func (c *DNSGetRecordsForIPCommand) Execute(args []string) error {
	InitDNSProvider(c.Provider, c.ZoneId, c.TTL)
	recordIds, err := dns.Provider.GetRecordsForIP(c.IP)
	if err == nil {
		Log("-> records:")
		for _, id := range recordIds {
			Log("->   %s", id)
		}
	}
	return Output(map[string]interface{}{"ids": recordIds}, recordIds, err)
}

type DNSDeleteRecordsForIPCommand struct {
	Provider string `long:"provider" default:"route53" description:"the dns provider"`
	ZoneId   string `short:"z" long:"zone" description:"the dns zone to use"`
	TTL      uint   `long:"ttl" default:"10" description:"the ttl to use"`
	IP       string `short:"i" long:"ip" description:"the ip to use"`
}

func (c *DNSDeleteRecordsForIPCommand) Execute(args []string) error {
	InitDNSProvider(c.Provider, c.ZoneId, c.TTL)
	err := dns.DeleteRecordsForIP(c.IP)
	if err == nil {
		Log("-> deleted")
	}
	return Output(nil, nil, err)
}
*/
