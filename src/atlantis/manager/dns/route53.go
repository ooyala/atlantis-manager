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

package dns

import (
	"errors"
	"github.com/jigish/route53/src/route53"
	"strings"
	"time"
)

type Route53Provider struct {
	r53   *route53.Route53
	Zones map[string]*Route53HostedZone // region -> hosted zone
	TTL   uint
}

type Route53HostedZone struct {
	ID   string
	Zone route53.HostedZone
}

func (r *Route53Provider) CreateRecords(region, comment string, records []Record) error {
	aliases := []*Alias{}
	cnames := []*CName{}
	arecords := []*ARecord{}
	for _, record := range records {
		switch typedRecord := record.(type) {
		case *Alias:
			aliases = append(aliases, typedRecord)
		case *CName:
			cnames = append(cnames, typedRecord)
		case *ARecord:
			arecords = append(arecords, typedRecord)
		default:
			return errors.New("Unsupported record type")
		}
	}
	if len(aliases) > 0 {
		err, errChan := r.CreateAliases(region, comment, aliases)
		if err != nil {
			return err
		}
		err = <-errChan
		if err != nil {
			return err
		}
	}
	if len(cnames) > 0 {
		err, errChan := r.CreateCNames(region, comment, cnames)
		if err != nil {
			return err
		}
		err = <-errChan
		if err != nil {
			return err
		}
	}
	if len(arecords) > 0 {
		err, errChan := r.CreateARecords(region, comment, arecords)
		if err != nil {
			return err
		}
		err = <-errChan
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Route53Provider) createRecords(region, comment string, rrsets ...route53.RRSet) (error, chan error) {
	hostedZone := r.Zones[region]
	if hostedZone == nil {
		return errors.New("Can't Find Route53 Hosted Zone for Region: " + region), nil
	}
	changes := make([]route53.RRSetChange, len(rrsets))
	for i, rrset := range rrsets {
		changes[i] = route53.RRSetChange{
			Action: "CREATE",
			RRSet:  rrset,
		}
	}
	info, err := r.r53.ChangeRRSet(hostedZone.ID, changes, comment)
	if err != nil {
		return err, nil
	}
	return nil, info.PollForSync(5*time.Second, 10*time.Minute)
}

func (r *Route53Provider) baseRRSet(typ, id, name, failover string) route53.RRSet {
	rrset := route53.RRSet{
		Name:          name,
		Type:          typ,
		SetIdentifier: id,
	}
	if failover == "PRIMARY" || failover == "SECONDARY" {
		rrset.Failover = failover
	}
	return rrset
}

func (r *Route53Provider) CreateCNames(region, comment string, cnames []*CName) (error, chan error) {
	rrsets := make([]route53.RRSet, len(cnames))
	count := 0
	for _, cname := range cnames {
		rrsets[count] = r.baseRRSet("CNAME", cname.ID(), cname.Name, cname.Failover)
		rrsets[count].TTL = r.TTL
		rrsets[count].ResourceRecords = &route53.ResourceRecords{
			ResourceRecord: []route53.ResourceRecord{route53.ResourceRecord{Value: cname.Original}},
		}
		rrsets[count].Weight = cname.Weight
		if cname.HealthCheckID != "" {
			rrsets[count].HealthCheckID = cname.HealthCheckID
		}
		count++
	}
	return r.createRecords(region, comment, rrsets...)
}

func (r *Route53Provider) CreateAliases(region, comment string, aliases []*Alias) (error, chan error) {
	hostedZone := r.Zones[region]
	if hostedZone == nil {
		return errors.New("Can't Find Route53 Hosted Zone for Region: " + region), nil
	}
	rrsets := make([]route53.RRSet, len(aliases))
	count := 0
	for _, alias := range aliases {
		rrsets[count] = r.baseRRSet("A", alias.ID(), alias.Alias, alias.Failover)
		rrsets[count].AliasTarget = &route53.AliasTarget{
			HostedZoneID:         hostedZone.ID,
			DNSName:              alias.Original,
			EvaluateTargetHealth: false,
		}
		count++
	}
	return r.createRecords(region, comment, rrsets...)
}

func (r *Route53Provider) CreateARecords(region, comment string, arecords []*ARecord) (error, chan error) {
	rrsets := make([]route53.RRSet, len(arecords))
	count := 0
	for _, arecord := range arecords {
		rrsets[count] = r.baseRRSet("A", arecord.ID(), arecord.Name, arecord.Failover)
		rrsets[count].TTL = r.TTL
		rrsets[count].ResourceRecords = &route53.ResourceRecords{
			ResourceRecord: []route53.ResourceRecord{route53.ResourceRecord{Value: arecord.IP}},
		}
		rrsets[count].Weight = arecord.Weight
		if arecord.HealthCheckID != "" {
			rrsets[count].HealthCheckID = arecord.HealthCheckID
		}
		count++
	}
	return r.createRecords(region, comment, rrsets...)
}

func (r *Route53Provider) DeleteRecords(region, comment string, ids ...string) (error, chan error) {
	hostedZone := r.Zones[region]
	if hostedZone == nil {
		return errors.New("Can't Find Route53 Hosted Zone for Region: " + region), nil
	}
	if len(ids) == 0 {
		errChan := make(chan error)
		go func(ch chan error) { // fake channel with nil error
			ch <- nil
		}(errChan)
		return nil, errChan
	}
	// fetch all records
	rrsets, err := r.r53.ListRRSets(hostedZone.ID)
	if err != nil {
		return err, nil
	}
	// create record map to make things easier
	rrsetMap := map[string]route53.RRSet{}
	for _, rrset := range rrsets {
		rrsetMap[rrset.SetIdentifier] = rrset
	}
	// filter by id and delete
	toDelete := []route53.RRSet{}
	for _, id := range ids {
		if rrset, exists := rrsetMap[id]; exists {
			toDelete = append(toDelete, rrset)
		}
	}
	changes := make([]route53.RRSetChange, len(toDelete))
	for i, rrset := range toDelete {
		changes[i] = route53.RRSetChange{
			Action: "DELETE",
			RRSet:  rrset,
		}
	}
	info, err := r.r53.ChangeRRSet(hostedZone.ID, changes, comment)
	if err != nil {
		return err, nil
	}
	return nil, info.PollForSync(5*time.Second, 10*time.Minute)
}

func (r *Route53Provider) CreateHealthCheck(ip string, port uint16) (string, error) {
	// health check to make sure TCP 80 is reachable
	config := route53.HealthCheckConfig{
		IPAddress: ip,
		Port:      port,
		Type:      "TCP",
	}
	// add health check for ip, return health check id
	return r.r53.CreateHealthCheck(config, "")
}

func (r *Route53Provider) DeleteHealthCheck(id string) error {
	return r.r53.DeleteHealthCheck(id)
}

func (r *Route53Provider) GetRecordsForValue(region, value string) ([]string, error) {
	hostedZone := r.Zones[region]
	if hostedZone == nil {
		return nil, errors.New("Can't Find Route53 Hosted Zone for Region: " + region)
	}
	rrsets, err := r.r53.ListRRSets(hostedZone.ID)
	if err != nil {
		return nil, err
	}
	ids := []string{}
	for _, rrset := range rrsets {
		if rrset.ResourceRecords != nil {
			for _, record := range rrset.ResourceRecords.ResourceRecord {
				if record.Value == value {
					ids = append(ids, rrset.SetIdentifier)
				}
			}
		} else if rrset.AliasTarget != nil {
			if rrset.AliasTarget.DNSName == value {
				ids = append(ids, rrset.SetIdentifier)
			}
		}
	}
	return ids, nil
}

func (r *Route53Provider) Suffix(region string) (string, error) {
	hostedZone := r.Zones[region]
	if hostedZone == nil {
		return "", errors.New("Can't Find Route53 Hosted Zone for Region: " + region)
	}
	return strings.TrimRight(hostedZone.Zone.Name, "."), nil
}

func NewRoute53Provider(zoneIDs map[string]string, ttl uint) (*Route53Provider, error) {
	route53.DebugOn()
	r53, err := route53.New()
	if err != nil {
		return nil, err
	}
	r53.IncludeWeight = true
	zones := map[string]*Route53HostedZone{}
	for region, zoneID := range zoneIDs {
		zone, err := r53.GetHostedZone(zoneID)
		if err != nil {
			return nil, err
		}
		zones[region] = &Route53HostedZone{ID: zoneID, Zone: zone}
	}
	return &Route53Provider{r53: r53, Zones: zones, TTL: ttl}, nil
}
