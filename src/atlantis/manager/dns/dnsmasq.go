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
	"bufio"
	"errors"
	"io"
	"os"
)

const fileHeaderString = "## This file managed by Atlantis manager.  DO NOT EDIT ##\n"

// TODO(edanaher): The records should be stored (as comments) in the hosts file to persist across master
// restarts.
type DnsmasqProvider struct {
	file    string
	records []Record
}

func readUntil(marker string, from *bufio.Reader, to *os.File) error {
	line, err := from.ReadString('\n')
	for ; err == nil; line, err = from.ReadString('\n') {
		if line == marker {
			break
		}
		if to != nil {
			io.WriteString(to, line)
		}
	}
	return err
}

func (d *DnsmasqProvider) rewriteHosts() error {
  // TODO(edanaher): This will be the code for reading the existing dns state.
	/*oldFile, err := os.Open(d.file)
	if err != nil {
		return err
	}
	defer oldFile.Close()
	oldReader := bufio.NewReader(oldFile) */

	newFile, err := os.Create(d.file + ".new")
	if err != nil {
		return err
	}
	defer newFile.Close()

	io.WriteString(newFile, fileHeaderString)
	hosts, err := d.getHosts()
	if err != nil {
		return err
	}
	for host, ip := range hosts {
		io.WriteString(newFile, ip+" "+host+"\n")
	}
	if err := os.Rename(d.file+".new", d.file); err != nil {
		return err
	}

	return nil
}

func (d *DnsmasqProvider) getHosts() (map[string]string, error) {
	// TODO(edanaher): This code is copied from route53.  Is it worth pulling out?
	aliases := []*Alias{}
	cnames := []*CName{}
	arecords := []*ARecord{}
	for _, record := range d.records {
		switch typedRecord := record.(type) {
		case *Alias:
			aliases = append(aliases, typedRecord)
		case *CName:
			cnames = append(cnames, typedRecord)
		case *ARecord:
			arecords = append(arecords, typedRecord)
		default:
			return nil, errors.New("Unsupported record type")
		}
	}
	hosts := map[string]string{}
	for _, arecord := range arecords {
		hosts[arecord.Name] = arecord.IP
	}
	// Loop until fixed point to handle recursive cnames and aliases.
	for len(cnames) > 0 {
		unknownNames := []*CName{}
		for _, cname := range cnames {
			if target, ok := hosts[cname.Original]; ok {
				hosts[cname.Name] = target
			} else {
				unknownNames = append(unknownNames, cname)
			}
		}
		if len(cnames) == len(unknownNames) {
			break // Should this be an error?  Unclear.
		}
		cnames = unknownNames
	}
	for len(aliases) > 0 {
		unknownAliases := []*Alias{}
		for _, alias := range aliases {
			if target, ok := hosts[alias.Original]; ok {
				hosts[alias.Alias] = target
			} else {
				unknownAliases = append(unknownAliases, alias)
			}
		}
		if len(aliases) == len(unknownAliases) {
			break // Should this be an error?  Unclear.
		}
		aliases = unknownAliases
	}
	return hosts, nil
}

func (d *DnsmasqProvider) CreateRecords(region, comment string, records []Record) error {
	d.records = append(d.records, records...)
	d.rewriteHosts()
	return nil
}

func (d *DnsmasqProvider) GetRecordsForValue(region, value string) ([]string, error) {
	records := []string{}
	for _, r := range d.records {
		switch typedRecord := r.(type) {
		case *ARecord:
			if typedRecord.IP == value {
				records = append(records, typedRecord.IP)
			}
		case *CName:
			if typedRecord.Original == value {
				records = append(records, typedRecord.Name)
			}
		case *Alias:
			if typedRecord.Original == value {
				records = append(records, typedRecord.Alias)
			}
		default:
			return nil, errors.New("Unsupported record type")
		}
	}
	return records, nil
}

func (d *DnsmasqProvider) DeleteRecords(region, comment string, ids ...string) (error, chan error) {
	idsMap := map[string]bool{}
	for _, id := range ids {
		idsMap[id] = true
	}
	for i, record := range d.records {
		if _, exists := idsMap[record.ID()]; exists {
			d.records = append(d.records[:i], d.records[i+1:]...)
		}
	}

	d.rewriteHosts()
	errChan := make(chan error)
	go func() { errChan <- nil }()
	return nil, errChan
}

func (d *DnsmasqProvider) Suffix(region string) (string, error) {
	return "aquarium", nil
}

func NewDnsmasqProvider(file string) (*DnsmasqProvider, error) {
	return &DnsmasqProvider{file: file}, nil
}
