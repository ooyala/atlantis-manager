package ldap

import (
	"atlantis/crypto"
	"crypto/tls"
	"errors"
	"github.com/mavricknz/ldap"
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
	LDAPConn *ldap.LDAPConnection
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

func CreateSession(req *Request, lc *ldap.LDAPConnection) {
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

func LookupConnection(user, secret string) *ldap.LDAPConnection {
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
		return "", nil // just let everything pass
	}
	// Checking if we are already logged in
	var LDAPConn *ldap.LDAPConnection
	req := &Request{user, secret, false, make(chan bool)}
	LookupSession(req)
	if !req.LoggedIn {
		LDAPConn = ldap.NewLDAPSSLConnection(LdapServer, LdapPort, TlsConfig)
		err := LDAPConn.Connect()
		if err != nil {
			return "", err
		}
		err = LoginBind(user, pass, LDAPConn)
		if err != nil {
			return "", err
		}
		now := strconv.FormatInt(time.Now().Unix(), 10)
		sec := string(crypto.Encrypt([]byte(pass + now)))
		re := regexp.MustCompile("[^a-zA-Z0-9]")
		sec = re.ReplaceAllString(sec, "")
		req.Secret = sec
	}
	CreateSession(req, LDAPConn)
	return req.Secret, nil
}

func LoginBind(user, pass string, lc *ldap.LDAPConnection) error {
	filterStr := "(&(objectClass=" + UserClass + ")(" + UserClassAttr + "=" + user + "))"
	searchReq := ldap.NewSimpleSearchRequest(BaseDomain, 2, filterStr, []string{UserClassAttr})
	sr, err := lc.Search(searchReq)
	if err != nil {
		return err
	}
	var dnInfo string
	if len(sr.Entries) == 1 {
		dnInfo = sr.Entries[0].DN
	} else {
		log.Println("No entries found")
		return errors.New("Invalid Credentials")
	}
	err = lc.Bind(dnInfo, pass)
	if err != nil {
		log.Println("ERROR : Login not Successful")
		return errors.New("Session Expired/Invalid Credentials")
	}
	return nil
}
