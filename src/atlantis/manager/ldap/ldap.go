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
	"errors"
	"github.com/mavricknz/ldap"
	"log"
	"regexp"
	"strconv"
	"strings"
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
	SearchUserDn         string
	SearchUserPwd	     string
	TeamBlackList        []string
)

type Request struct {
	User     string
	Secret   string
	LoggedIn bool
	RespChan chan bool
}

type Session struct {

	//TODO remove LDAPConn from session
	//since we can no longer re-use ldap conn with jump cloud
 
	LDAPConn *ldap.LDAPConnection
	Team     []string
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

func CreateSession(req *Request, lc *ldap.LDAPConnection, team []string) {
	if SessionMap[req.User] == nil {
		SessionMap[req.User] = map[string]*Session{req.Secret: &Session{LDAPConn: lc, Team: team,
			Timer: time.AfterFunc(30*time.Minute, func() {
				SessionDestroyChan <- req
			})},
		}
	} else if SessionMap[req.User][req.Secret] == nil {
		SessionMap[req.User][req.Secret] = &Session{LDAPConn: lc, Team: team, Timer: time.AfterFunc(30*time.Minute, func() {
			SessionDestroyChan <- req
		})}
	} else {
		SessionMap[req.User][req.Secret].Timer.Reset(30*time.Minute)
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

func LookupConnection(user, secret string) *ldap.LDAPConnection {
	if SessionMap[user] != nil && SessionMap[user][secret] != nil {
		return SessionMap[user][secret].LDAPConn
	}
	return nil
}

func LookupTeam(user, secret string) []string {
	if SessionMap[user] != nil && SessionMap[user][secret] != nil {
		return SessionMap[user][secret].Team
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
	var LDAPConn *ldap.LDAPConnection
	var err error
	req := &Request{user, secret, false, make(chan bool)}
	LookupSession(req)
	teamList := []string{}
	if !req.LoggedIn {
		LDAPConn, err = CreateLdapConn(LdapServer, LdapPort, TlsConfig)
		
		if err != nil {
			return "", err
		}
		err = LoginBind(user, pass, LDAPConn)
		if err != nil {
			return "", err
		}

		//bind with account with dn search enabled
		err = LoginBind(SearchUserDn, SearchUserPwd, LDAPConn)
		if err != nil {
			log.Println("Warning: ldap binding with dn search account failed; ", err)
		}

		now := strconv.FormatInt(time.Now().Unix(), 10)
		sec := string(crypto.Encrypt([]byte(pass + now)))
		re := regexp.MustCompile("[^a-zA-Z0-9]")
		sec = re.ReplaceAllString(sec, "")
		req.Secret = sec
		teamList, err = GetTeamList(LDAPConn, user)
		if err != nil {
			log.Println("Warning: get team list failed for user; ", err)
		}
	}
	CreateSession(req, LDAPConn, teamList)
	return req.Secret, nil
}

func CreateLdapConn(server string, port uint16, tlsConf *tls.Config) (*ldap.LDAPConnection, error) {
	LDAPConn := ldap.NewLDAPSSLConnection(server, port, tlsConf)
	err := LDAPConn.Connect()
	if err != nil {
		return nil, err
	}
	return LDAPConn, nil
}

func GetTeamList(LDAPConn *ldap.LDAPConnection , user string) ([]string, error){
	//should be something like (&(objectClass=posixAccount)(uid=xxxx))
        filterStr := "(&(objectClass=" + UserClass + ")(" + UserCommonName + "=" + user + "))" 

	searchReq := ldap.NewSimpleSearchRequest(BaseDomain, 2, filterStr, []string{"memberOf"})
	sr, err := LDAPConn.Search(searchReq)

	ret := []string{}
	if err != nil || sr == nil {
		return ret, err
	}

	for _, entry := range sr.Entries {
		vals := entry.GetAttributeValues("memberOf")
		r, _ := regexp.Compile("^cn=([^,]+)")
		if len(vals) > 0 {
			for _, teamDn := range vals {
				substrings := strings.Split(r.FindString(teamDn), "=")
				
				if len(substrings) == 2 && !contains(TeamBlackList, substrings[1]) {
					ret = append(ret, substrings[1])
				}
			}
		}
	}

	return ret, nil
}


func LoginBind(user, pass string, lc *ldap.LDAPConnection) error {
	var dnInfo string
	dnInfo = UserClassAttr + "=" + user + "," + UserOu
	err := lc.Bind(dnInfo, pass)
	if err != nil {
		log.Println("ERROR : Login not Successful")
		return errors.New("Session Expired/Invalid Credentials")
	}
	return nil
}

func contains(arr []string, str string) bool {
   for _, a := range arr {
      if a == str {
         return true
      }
   }
   return false
}
