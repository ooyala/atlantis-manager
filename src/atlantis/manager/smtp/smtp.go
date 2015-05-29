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

package smtp

import (
	"fmt"
	"log"
	"net"
	"net/smtp"
	"time"
)

var (
	addr string
	from string
	cc   []string
)

func Init(myAddr, myFrom string, myCC []string) {
	addr = myAddr
	from = myFrom
	cc = myCC
}

func logError(err error, to []string, subject, body string) error {
	log.Printf("[SMTP] ERROR: Failed to send: %s", err.Error())
	log.Printf("[SMTP]        addr: %s", addr)
	log.Printf("[SMTP]        from: %s", from)
	log.Printf("[SMTP]        to:   %v", to)
	log.Printf("[SMTP]        body:\n%s", body)
	return err
}

func SendMail(to []string, subject, body string) (err error) {
	if addr == "" {
		// skip if not configured
		return nil
	}
	// create client
	var conn net.Conn
	var client *smtp.Client
	if conn, err = net.Dial("tcp", addr+":25"); err != nil {
		return logError(err, to, subject, body)
	}
	if client, err = smtp.NewClient(conn, ""); err != nil {
		return logError(err, to, subject, body)
	}
	defer client.Quit()
	msg := fmt.Sprintf("From: %s\n", from)
	if err = client.Mail(from); err != nil {
		return logError(err, to, subject, body)
	}
	for _, addr := range to {
		msg += fmt.Sprintf("To: %s\n", addr)
		if err = client.Rcpt(addr); err != nil {
			return logError(err, to, subject, body)
		}
	}
	for _, addr := range cc {
		msg += fmt.Sprintf("Cc: %s\n", addr)
		if err = client.Rcpt(addr); err != nil {
			return logError(err, to, subject, body)
		}
	}
	msg += fmt.Sprintf("Date: %s\n", time.Now().String())
	msg += fmt.Sprintf("Subject: %s\n\n", subject)
	msg += body

	w, err := client.Data()
	if err != nil {
		return logError(err, to, subject, body)
	}

	fmt.Fprintf(w, msg)
	w.Close()
	log.Printf("[SMTP] SENT: %v '%s'", to, subject)
	return nil
}
