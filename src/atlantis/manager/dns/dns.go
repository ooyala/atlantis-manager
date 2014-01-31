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
