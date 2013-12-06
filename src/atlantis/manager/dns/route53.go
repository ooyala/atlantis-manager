package dns

import (
	"errors"
	"github.com/crowdmob/goamz/aws"
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
	Id   string
	Zone route53.HostedZone
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
	info, err := r.r53.ChangeRRSet(hostedZone.Id, changes, comment)
	if err != nil {
		return err, nil
	}
	return nil, info.PollForSync(5*time.Second, 10*time.Minute)
}

func (r *Route53Provider) baseRRSet(id, name, failover string) route53.RRSet {
	rrset := route53.RRSet{
		Name:          name,
		Type:          "A",
		SetIdentifier: id,
	}
	if failover == "PRIMARY" || failover == "SECONDARY" {
		rrset.Failover = failover
	}
	return rrset
}

func (r *Route53Provider) CreateAliases(region, comment string, aliases []Alias) (error, chan error) {
	hostedZone := r.Zones[region]
	if hostedZone == nil {
		return errors.New("Can't Find Route53 Hosted Zone for Region: " + region), nil
	}
	rrsets := make([]route53.RRSet, len(aliases))
	count := 0
	for _, alias := range aliases {
		rrsets[count] = r.baseRRSet(alias.Id(), alias.Alias, alias.Failover)
		rrsets[count].AliasTarget = &route53.AliasTarget{
			HostedZoneId:         hostedZone.Id,
			DNSName:              alias.Original,
			EvaluateTargetHealth: false,
		}
		count++
	}
	return r.createRecords(region, comment, rrsets...)
}

func (r *Route53Provider) CreateARecords(region, comment string, arecords []ARecord) (error, chan error) {
	rrsets := make([]route53.RRSet, len(arecords))
	count := 0
	for _, arecord := range arecords {
		rrsets[count] = r.baseRRSet(arecord.Id(), arecord.Name, arecord.Failover)
		rrsets[count].TTL = r.TTL
		rrsets[count].ResourceRecords = &route53.ResourceRecords{
			ResourceRecord: []route53.ResourceRecord{route53.ResourceRecord{Value: arecord.IP}},
		}
		rrsets[count].Weight = arecord.Weight
		if arecord.HealthCheckId != "" {
			rrsets[count].HealthCheckId = arecord.HealthCheckId
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
	rrsets, err := r.r53.ListRRSets(hostedZone.Id)
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
	info, err := r.r53.ChangeRRSet(hostedZone.Id, changes, comment)
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

func (r *Route53Provider) GetRecordsForIP(region, ip string) ([]string, error) {
	hostedZone := r.Zones[region]
	if hostedZone == nil {
		return nil, errors.New("Can't Find Route53 Hosted Zone for Region: " + region)
	}
	rrsets, err := r.r53.ListRRSets(hostedZone.Id)
	if err != nil {
		return nil, err
	}
	ids := []string{}
	for _, rrset := range rrsets {
		if rrset.ResourceRecords == nil {
			continue
		}
		for _, record := range rrset.ResourceRecords.ResourceRecord {
			if record.Value == ip {
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

func NewRoute53Provider(zoneIds map[string]string, ttl uint) (*Route53Provider, error) {
	route53.DebugOn()
	auth, err := aws.GetAuth("", "", "", time.Time{})
	if err != nil {
		return nil, err
	}
	r53 := route53.New(auth)
	r53.IncludeWeight = true
	zones := map[string]*Route53HostedZone{}
	for region, zoneId := range zoneIds {
		zone, err := r53.GetHostedZone(zoneId)
		if err != nil {
			return nil, err
		}
		zones[region] = &Route53HostedZone{Id: zoneId, Zone: zone}
	}
	return &Route53Provider{r53: r53, Zones: zones, TTL: ttl}, nil
}
