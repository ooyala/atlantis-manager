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
	if conn, err = net.Dial("tcp", "localhost:25"); err != nil {
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
	if w, err := client.Data(); err != nil {
		return logError(err, to, subject, body)
	} else {
		fmt.Fprintf(w, msg)
		w.Close()
		log.Printf("[SMTP] SENT: %v '%s'", to, subject)
	}
	return nil
}
