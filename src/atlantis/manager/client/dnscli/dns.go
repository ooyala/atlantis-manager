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

package dnscli

import (
	"atlantis/manager/dns"
	"errors"
)

var suffix = ""

func InitDNSProvider(provider, zone string, ttl uint) error {
	var err error
	switch provider {
	case "route53":
		dns.Provider, err = dns.NewRoute53Provider(map[string]string{"cli": zone}, ttl)
		if err != nil {
			return err
		}
		suffix, err = dns.Provider.Suffix("cli")
		if err != nil {
			return err
		}
		return nil
	default:
		dns.Provider = nil
	}
	return errors.New("Invalid DNS Provider")
}

type DNSCreateARecordCommand struct {
	Provider      string `long:"provider" default:"route53" description:"the dns provider"`
	ZoneID        string `short:"z" long:"zone" description:"the dns zone to use"`
	TTL           uint   `long:"ttl" default:"10" description:"the ttl to use"`
	Prefix        string `short:"p" long:"prefix" description:"the name prefix to use"`
	IP            string `short:"i" long:"ip" description:"the ip to use"`
	HealthCheckID string `short:"H" long:"health-check-id" description:"the health check id to use"`
	Failover      string `short:"f" long:"failover" description:"the failover policy to use"`
	Weight        uint8  `short:"w" long:"weight" description:"the record's weight"`
	Comment       string `long:"comment" default:"CLIENT" description:"the comment to use"`
}

func (c *DNSCreateARecordCommand) Execute(args []string) error {
	if err := InitDNSProvider(c.Provider, c.ZoneID, c.TTL); err != nil {
		return Output(nil, nil, err)
	}
	arecord := dns.ARecord{
		Name:          c.Prefix + "." + suffix,
		IP:            c.IP,
		HealthCheckID: c.HealthCheckID,
		Failover:      c.Failover,
		Weight:        c.Weight,
	}
	err := dns.Provider.CreateRecords("cli", c.Comment, []dns.Record{&arecord})
	if err != nil {
		return Output(nil, nil, err)
	}
	Log("-> created %s", arecord.ID())
	return Output(map[string]interface{}{"id": arecord.ID()}, arecord.ID(), err)
}

type DNSCreateAliasCommand struct {
	Provider string `long:"provider" default:"route53" description:"the dns provider"`
	ZoneID   string `short:"z" long:"zone" description:"the dns zone to use"`
	TTL      uint   `long:"ttl" default:"10" description:"the ttl to use"`
	Prefix   string `short:"p" long:"prefix" description:"the name prefix to use for the alias"`
	Original string `short:"o" long:"original" description:"the target of the alias"`
	Failover string `short:"f" long:"failover" description:"the failover policy to use"`
	Comment  string `long:"comment" default:"CLIENT" description:"the comment to use"`
}

func (c *DNSCreateAliasCommand) Execute(args []string) error {
	if err := InitDNSProvider(c.Provider, c.ZoneID, c.TTL); err != nil {
		return Output(nil, nil, err)
	}
	alias := dns.Alias{
		Alias:    c.Prefix + "." + suffix,
		Original: c.Original,
		Failover: c.Failover,
	}
	err := dns.Provider.CreateRecords("cli", c.Comment, []dns.Record{&alias})
	if err != nil {
		return Output(nil, nil, err)
	}
	Log("-> created %s", alias.ID())
	return Output(map[string]interface{}{"id": alias.ID()}, alias.ID(), err)
}

type DNSCreateCNameCommand struct {
	Provider string `long:"provider" default:"route53" description:"the dns provider"`
	ZoneID   string `short:"z" long:"zone" description:"the dns zone to use"`
	TTL      uint   `long:"ttl" default:"10" description:"the ttl to use"`
	Prefix   string `short:"p" long:"prefix" description:"the name prefix to use for the alias"`
	Original string `short:"o" long:"original" description:"the target of the alias"`
	Failover string `short:"f" long:"failover" description:"the failover policy to use"`
	Weight   uint8  `short:"w" long:"weight" description:"the record's weight"`
	Comment  string `long:"comment" default:"CLIENT" description:"the comment to use"`
}

func (c *DNSCreateCNameCommand) Execute(args []string) error {
	if err := InitDNSProvider(c.Provider, c.ZoneID, c.TTL); err != nil {
		return Output(nil, nil, err)
	}
	cname := dns.CName{
		Name:     c.Prefix + "." + suffix,
		Original: c.Original,
		Failover: c.Failover,
		Weight:   c.Weight,
	}
	err := dns.Provider.CreateRecords("cli", c.Comment, []dns.Record{&cname})
	if err != nil {
		return Output(nil, nil, err)
	}
	Log("-> created %s", cname.ID())
	return Output(map[string]interface{}{"id": cname.ID()}, cname.ID(), err)
}

type DNSDeleteRecordsCommand struct {
	Provider  string   `long:"provider" default:"route53" description:"the dns provider"`
	ZoneID    string   `short:"z" long:"zone" description:"the dns zone to use"`
	TTL       uint     `long:"ttl" default:"10" description:"the ttl to use"`
	RecordIDs []string `short:"i" long:"id" description:"the record ids to delete"`
	Comment   string   `long:"comment" default:"CLIENT" description:"the comment to use"`
}

func (c *DNSDeleteRecordsCommand) Execute(args []string) error {
	if err := InitDNSProvider(c.Provider, c.ZoneID, c.TTL); err != nil {
		return Output(nil, nil, err)
	}
	err, errChan := dns.Provider.DeleteRecords("cli", c.Comment, c.RecordIDs...)
	if err != nil {
		return Output(nil, nil, err)
	}
	err = <-errChan
	if err == nil {
		Log("-> deleted:")
		for _, id := range c.RecordIDs {
			Log("->   %s", id)
		}
	}
	return Output(map[string]interface{}{"ids": c.RecordIDs}, c.RecordIDs, err)
}

type DNSDeleteCNameCommand struct {
	Provider  string `long:"provider" default:"route53" description:"the dns provider"`
	ZoneID    string `short:"z" long:"zone" description:"the dns zone to use"`
	TTL       uint   `long:"ttl" default:"10" description:"the ttl to use"`
	Prefix    string `short:"p" long:"prefix" description:"the name prefix to use for the alias"`
	Original  string `short:"o" long:"original" description:"the target of the alias"`
	Comment   string `long:"comment" default:"CLIENT" description:"the comment to use"`
}

func (c *DNSDeleteCNameCommand) Execute(args []string) error {
	if err := InitDNSProvider(c.Provider, c.ZoneID, c.TTL); err != nil {
		return Output(nil, nil, err)
	}
	cname := dns.CName{
		Name:     c.Prefix + "." + suffix,
		Original: c.Original,
	}
	recordID := cname.ID()
	err, errChan := dns.Provider.DeleteRecords("cli", c.Comment, recordID)
	if err != nil {
		return Output(nil, nil, err)
	}
	err = <-errChan
	if err == nil {
		Log("-> deleted: ", recordID)
	}
	return Output(map[string]interface{}{"id": recordID}, recordID, err)
}

type DNSGetRecordsForValueCommand struct {
	Provider string `long:"provider" default:"route53" description:"the dns provider"`
	ZoneID   string `short:"z" long:"zone" description:"the dns zone to use"`
	TTL      uint   `long:"ttl" default:"10" description:"the ttl to use"`
	Value    string `short:"v" long:"value" description:"the value to use"`
}

func (c *DNSGetRecordsForValueCommand) Execute(args []string) error {
	if err := InitDNSProvider(c.Provider, c.ZoneID, c.TTL); err != nil {
		return Output(nil, nil, err)
	}
	recordIDs, err := dns.Provider.GetRecordsForValue("cli", c.Value)
	if err == nil {
		Log("-> records:")
		for _, id := range recordIDs {
			Log("->   %s", id)
		}
	}
	return Output(map[string]interface{}{"ids": recordIDs}, recordIDs, err)
}

type DNSDeleteRecordsForValueCommand struct {
	Provider string `long:"provider" default:"route53" description:"the dns provider"`
	ZoneID   string `short:"z" long:"zone" description:"the dns zone to use"`
	TTL      uint   `long:"ttl" default:"10" description:"the ttl to use"`
	Value    string `short:"v" long:"value" description:"the value to use"`
}

func (c *DNSDeleteRecordsForValueCommand) Execute(args []string) error {
	if err := InitDNSProvider(c.Provider, c.ZoneID, c.TTL); err != nil {
		return Output(nil, nil, err)
	}
	err := dns.DeleteRecordsForValue("cli", c.Value)
	if err == nil {
		Log("-> deleted")
	}
	return Output(nil, nil, err)
}
