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

package ldap

import (
	"atlantis/crypto"
	"crypto/tls"
	"fmt"
	goLdap "github.com/go-ldap/ldap"
	"log"
	"regexp"
	"strconv"
	"time"
)

var (
	SessionDestroyChan   chan *Request
	BaseDomain           string
	LdapServer           string
	LdapPort             uint16
	TlsConfig            *tls.Config
	skipLogin            bool
	SessionMap           map[string]map[string]*Session // map user -> (map secret -> *Session)
	AppClass             string
	UsernameAttr         string
	TeamAdminAttr        string
	AllowedAppAttr       string
	AllowedAppCommonName string
	TeamClass            string
	UserOu               string
	TeamOu               string
	UserCommonName       string
	TeamCommonName       string
	SuperUserGroup       string
	UserClass            string
	UserClassAttr        string
	SkipAuthorization    bool
)

type Request struct {
	User     string
	Secret   string
	LoggedIn bool
	RespChan chan bool
}

type Session struct {
	LDAPConn *goLdap.Conn
	Timer    *time.Timer
}

func Init(lserver string, lport uint16, baseDomain string) {
	if lserver == "" {
		// if we're being initialized empty, then don't try to log people in
		skipLogin = true
		return
	}
	TlsConfig = &tls.Config{}
	// TODO : Remove InsecureSkipVerify
	TlsConfig.InsecureSkipVerify = true
	LdapServer = lserver
	LdapPort = lport
	BaseDomain = baseDomain
	SessionMap = make(map[string]map[string]*Session)
	SessionDestroyChan = make(chan *Request)
	go SessionExpiryRoutine()
}

func CreateSession(req *Request, lc *goLdap.Conn) {
	if SessionMap[req.User] == nil {
		SessionMap[req.User] = map[string]*Session{req.Secret: &Session{LDAPConn: lc,
			Timer: time.AfterFunc(30*time.Minute, func() {
				SessionDestroyChan <- req
			})},
		}
	} else if SessionMap[req.User][req.Secret] == nil {
		SessionMap[req.User][req.Secret] = &Session{LDAPConn: lc, Timer: time.AfterFunc(30*time.Minute, func() {
			SessionDestroyChan <- req
		})}
	} else {
		SessionMap[req.User][req.Secret].Timer.Reset(30 * time.Minute)
	}
}

func LookupSession(req *Request) {
	req.LoggedIn = false
	if SessionMap[req.User] != nil && req.Secret != "" {
		if SessionMap[req.User][req.Secret] != nil {
			req.LoggedIn = true
		}
	}
}

func LookupConnection(user, secret string) *goLdap.Conn {
	if SessionMap[user] != nil && SessionMap[user][secret] != nil {
		return SessionMap[user][secret].LDAPConn
	}
	return nil
}

func SessionExpiryRoutine() {
	for {
		select {
		case req := <-SessionDestroyChan:
			SessionMap[req.User][req.Secret].LDAPConn.Close()
			delete(SessionMap[req.User], req.Secret)
			break
		}
	}
}

func Login(user, pass, secret string) (string, error) {
	if skipLogin {
		return "dummysecret", nil // just let everything pass
	}

	// Checking if we are already logged in
	var Conn *goLdap.Conn
	req := &Request{user, secret, false, make(chan bool)}
	LookupSession(req)
	if !req.LoggedIn {
		var err error
		log.Printf("Dialing TLS connection on port %d", LdapPort)
		Conn, err = goLdap.Dial("tcp", fmt.Sprintf("%s:%d", LdapServer, LdapPort))
		if err != nil {
			return "", err
		}
		log.Printf("Starting TLS connection")
		err = Conn.StartTLS(TlsConfig)
		if err != nil {
			return "", err
		}
		username := fmt.Sprintf("uid=%s,ou=humans,ou=users,dc=ooyala,dc=com", user)
		log.Printf("Binding TLS connection")
		err = Conn.Bind(username, pass)
		if err != nil {
			return "", err
		}

		now := strconv.FormatInt(time.Now().Unix(), 10)
		sec := string(crypto.Encrypt([]byte(pass + now)))
		re := regexp.MustCompile("[^a-zA-Z0-9]")
		sec = re.ReplaceAllString(sec, "")
		req.Secret = sec

		log.Printf("TLS connection successful")
	}

	CreateSession(req, Conn)
	return req.Secret, nil
}
