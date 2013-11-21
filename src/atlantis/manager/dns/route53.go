package dns

import (
	"github.com/alekar/route53/src/route53"
	"github.com/crowdmob/goamz/aws"
	"time"
)

type Route53Provider struct {
	r53  *route53.Route53
	Zone route53.HostedZone
	TTL  uint
}

func (r *Route53Provider) createRecords(comment string, rrsets ...route53.RRSet) (error, chan error) {
	changes := make([]route53.RRSetChange, len(rrsets))
	for i, rrset := range rrsets {
		changes[i] = route53.RRSetChange{
			Action: "CREATE",
			RRSet:  rrset,
		}
	}
	info, err := r.r53.ChangeRRSet(r.Zone.Id, changes, comment)
	if err != nil {
		return err, nil
	}
	return nil, r.r53.PollForSync(info.Id, time.Second, 60*time.Second)
}

func (r *Route53Provider) CreateAliases(comment string, aliases []Alias) (error, chan error) {
	rrsets := make([]route53.RRSet, len(aliases))
	count := 0
	for _, alias := range aliases {
		failover := "SECONDARY"
		if alias.Primary {
			failover = "PRIMARY"
		}
		rrsets[count] = route53.RRSet{
			Failover:             failover,
			Name:                 alias.Alias,
			Type:                 "A",
			TTL:                  r.TTL,
			SetIdentifier:        alias.Id(),
			Weight:               0,
			HostedZoneId:         r.Zone.Id,
			DNSName:              alias.Original,
			EvaluateTargetHealth: true,
		}
		count++
	}
	return r.createRecords(comment, rrsets...)
}

func (r *Route53Provider) CreateCNames(comment string, cnames []CName) (error, chan error) {
	rrsets := make([]route53.RRSet, len(cnames))
	count := 0
	for _, cname := range cnames {
		failover := "SECONDARY"
		if cname.Primary {
			failover = "PRIMARY"
		}
		rrsets[count] = route53.RRSet{
			Failover:      failover,
			Name:          cname.CName,
			Type:          "A",
			TTL:           r.TTL,
			Values:        []string{cname.IP},
			HealthCheckId: cname.HealthCheckId,
			SetIdentifier: cname.Id(),
			Weight:        0,
		}
		count++
	}
	return r.createRecords(comment, rrsets...)
}

func (r *Route53Provider) DeleteRecords(comment string, ids ...string) (error, chan error) {
	if len(ids) == 0 {
		errChan := make(chan error)
		go func(ch chan error) { // fake channel with nil error
			ch <- nil
		}(errChan)
		return nil, errChan
	}
	// fetch all records
	rrsets, err := r.r53.ListRRSets(r.Zone.Id)
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
	info, err := r.r53.ChangeRRSet(r.Zone.Id, changes, comment)
	if err != nil {
		return err, nil
	}
	return nil, r.r53.PollForSync(info.Id, time.Second, 60*time.Second)
}

func (r *Route53Provider) CreateHealthCheck(ip string) (string, error) {
	// health check to make sure TCP 80 is reachable
	config := route53.HealthCheckConfig{
		IPAddress: ip,
		Port:      80,
		Type:      "TCP",
	}
	// add health check for ip, return health check id
	return r.r53.CreateHealthCheck(config, "")
}

func (r *Route53Provider) DeleteHealthCheck(id string) error {
	return r.r53.DeleteHealthCheck(id)
}

func (r *Route53Provider) GetRecordsForIP(ip string) ([]string, error) {
	rrsets, err := r.r53.ListRRSets(r.Zone.Id)
	if err != nil {
		return nil, err
	}
	ids := []string{}
	for _, rrset := range rrsets {
		for _, value := range rrset.Values {
			if value == ip {
				ids = append(ids, rrset.SetIdentifier)
			}
		}
	}
	return ids, nil
}

func (r *Route53Provider) Suffix() string {
	return r.Zone.Name
}

func NewRoute53Provider(zoneId string, ttl uint) (*Route53Provider, error) {
	route53.DebugOn()
	auth, err := aws.GetAuth("", "", "", time.Time{})
	if err != nil {
		return nil, err
	}
	r53 := route53.New(auth)
	zone, err := r53.GetHostedZone(zoneId)
	if err != nil {
		return nil, err
	}
	return &Route53Provider{r53: r53, Zone: zone, TTL: ttl}, nil
}
