package dns

import (
	. "atlantis/manager/constant"
	"atlantis/manager/datamodel"
	"atlantis/manager/helper"
	"crypto/sha256"
	"fmt"
	"regexp"
)

var (
	IPRegexp = regexp.MustCompile("^[0-9]+\\.[0-9]+\\.[0-9]+\\.[0-9]+$")
)

var Provider DNSProvider

type DNSProvider interface {
	CreateRecords(string, string, []Record) error
	// CreateCNames(string, string, []*CName) (error, chan error) // used for CreateRecords
	// CreateARecords(string, string, []*ARecord) (error, chan error) // used for CreateRecords
	// CreateAliases(string, string, []*Alias) (error, chan error) // unused
	GetRecordsForValue(string, string) ([]string, error)
	DeleteRecords(string, string, ...string) (error, chan error)
	// CreateHealthCheck(string, uint16) (string, error) // unused
	// DeleteHealthCheck(string) error // unused
	Suffix(string) (string, error)
}

func NewRecord(name, original string, weight uint8) Record {
	if IPRegexp.MatchString(original) {
		return &ARecord{
			Name: name,
			IP:   original,
		}
	}
	return &CName{
		Name:     name,
		Original: original,
	}
}

type Record interface {
	ID() string
}

type Alias struct {
	Alias    string
	Original string
	Failover string
}

func (a *Alias) ID() string {
	checksumArr := sha256.Sum256([]byte(fmt.Sprintf("%s %s", a.Original, a.Alias)))
	return fmt.Sprintf("%x", checksumArr[:sha256.Size])
}

type ARecord struct {
	Name          string
	IP            string
	HealthCheckID string
	Failover      string
	Weight        uint8
}

func (a *ARecord) ID() string {
	checksumArr := sha256.Sum256([]byte(fmt.Sprintf("%s %s", a.IP, a.Name)))
	return fmt.Sprintf("%x", checksumArr[:sha256.Size])
}

type CName struct {
	Name          string
	Original      string
	HealthCheckID string
	Failover      string
	Weight        uint8
}

func (c *CName) ID() string {
	checksumArr := sha256.Sum256([]byte(fmt.Sprintf("%s %s", c.Original, c.Name)))
	return fmt.Sprintf("%x", checksumArr[:sha256.Size])
}

func CreateAppCNames(internal bool, app, sha, env string) error {
	// cnames are created for only this manager's region.
	suffix, err := Provider.Suffix(Region)
	if err != nil {
		return err
	}
	// check if records were created already, if so add sha to list
	zkDNS, err := datamodel.GetDNS(app, env)
	if zkDNS != nil && err == nil {
		if zkDNS.Shas == nil {
			zkDNS.Shas = map[string]bool{}
		}
		zkDNS.Shas[sha] = true
		return zkDNS.Save()
	}
	// for each zone
	zkDNS = datamodel.DNS(app, env)
	if zkDNS.Shas == nil {
		zkDNS.Shas = map[string]bool{}
	}
	zkDNS.Shas[sha] = true
	if Provider == nil {
		return zkDNS.Save()
	}

	zkDNS.RecordIDs = []string{}
	cnames := []Record{}
	// set up zone cname
	for _, zone := range AvailableZones {
		newCName := NewRecord(helper.GetZoneAppCName(app, env, zone, suffix),
			helper.GetZoneRouterCName(internal, zone, suffix), 1)
		cnames = append(cnames, newCName)
		zkDNS.RecordIDs = append(zkDNS.RecordIDs, newCName.ID())
	}
	// region-wide entry (for referencing outside of atlantis)
	regionCName := NewRecord(helper.GetRegionAppCName(app, env, suffix),
		helper.GetRegionRouterCName(internal, suffix), 1)
	cnames = append(cnames, regionCName)
	zkDNS.RecordIDs = append(zkDNS.RecordIDs, regionCName.ID())

	err = Provider.CreateRecords(Region, "CREATE_APP "+app+" in "+env, cnames)
	if err != nil {
		return err
	}
	// save records made in router zone path
	return zkDNS.Save()
}

func DeleteAppCNames(app, sha, env string) error {
	// find ids for app+env
	zkDNS, err := datamodel.GetDNS(app, env)
	if err != nil {
		return err
	}
	if zkDNS.Shas == nil {
		zkDNS.Shas = map[string]bool{}
	}
	// remove sha from sha references
	delete(zkDNS.Shas, sha)
	err = zkDNS.Save()
	if err != nil {
		return err
	}
	// if this was *not* the last sha, don't delete anything
	if len(zkDNS.Shas) > 0 {
		return nil
	}
	if Provider == nil {
		return zkDNS.Delete()
	}
	// delete all the record ids
	err, errChan := Provider.DeleteRecords(Region, "DELETE_APP "+app+" in "+env, zkDNS.RecordIDs...)
	if err != nil {
		return err
	}
	err = <-errChan // wait for change to propagate
	if err != nil {
		return err
	}
	// remove dns datamodel
	return zkDNS.Delete()
}

func DeleteRecordsForValue(region, value string) error {
	if Provider == nil {
		return nil
	}
	ids, err := Provider.GetRecordsForValue(region, value)
	if err != nil {
		return err
	}
	err, errChan := Provider.DeleteRecords(region, "DELETE_ALL "+value, ids...)
	if err != nil {
		return err
	}
	return <-errChan
}
