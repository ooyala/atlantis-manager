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
	"io"
	"os"
)

const fileHeaderString = "## The following records are managed by Atlantis ##\n"
const fileFooterString = "## The preceding records are managed by Atlantis ##\n"

type DnsmasqProvider struct {
	file  string
	hosts map[string]int
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
	oldFile, err := os.Open(d.file)
	if err != nil {
		return err
	}
	defer oldFile.Close()
	oldReader := bufio.NewReader(oldFile)

	newFile, err := os.Create(d.file + ".new")
	if err != nil {
		return err
	}
	defer newFile.Close()

	readUntil(fileHeaderString, oldReader, newFile)
	readUntil(fileFooterString, oldReader, nil)
	readUntil("eof", oldReader, newFile)

	io.WriteString(newFile, fileHeaderString)
	io.WriteString(newFile, "127.0.0.1 magic\n")
	io.WriteString(newFile, fileFooterString)

	if err := os.Rename(d.file+".new", d.file); err != nil {
		return err
	}

	return nil
}

func (d *DnsmasqProvider) CreateRecords(region, comment string, records []Record) error {

	return nil
}

func (d *DnsmasqProvider) GetRecordsForValue(region, value string) ([]string, error) {
	return nil, nil
}

func (d *DnsmasqProvider) DeleteRecords(region, comment string, ids ...string) (error, chan error) {
	return nil, nil
}

func (d *DnsmasqProvider) Suffix(region string) (string, error) {
	return "", nil
}

func NewDnsmasqProvider(file string) (*DnsmasqProvider, error) {
	return &DnsmasqProvider{file: file}, nil
}
