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

package server

import (
	. "atlantis/common"
	"atlantis/crypto"
	"atlantis/manager/api"
	"atlantis/manager/builder"
	. "atlantis/manager/constant"
	"atlantis/manager/datamodel"
	"atlantis/manager/dns"
	"atlantis/manager/ldap"
	"atlantis/manager/rpc"
	"atlantis/manager/smtp"
	iconst "atlantis/supervisor/constant"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/jigish/go-flags"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type ServerConfig struct {
	RpcAddr                    string `toml:"rpc_addr"`
	ApiAddr                    string `toml:"api_addr"`
	SupervisorPort             uint16 `toml:"supervisor_port"`
	ZookeeperUri               string `toml:"zookeeper_uri"`
	LdapHost                   string `toml:"ldap_host"`
	LdapPort                   uint16 `toml:"ldap_port"`
	LdapBaseDomain             string `toml:"ldap_basedomain"`
	CPUSharesIncrement         uint   `toml:"cpu_shares_increment"`
	MemoryLimitIncrement       uint   `toml:"memory_limit_increment"`
	ResultDuration             string `toml:"result_duration"`
	LdapUserSearchBase         string `toml:"ldap_user_search_base"`
	LdapRoleUserSearchBase	   string `toml:"ldap_role_user_search_base"`
	LdapTeamSearchBase         string `toml:"ldap_team_search_base"`
	LdapUsernameAttr           string `toml:"ldap_username_attr"`
	LdapAppClass               string `toml:"ldap_app_class"`
	LdapTeamClass              string `toml:"ldap_team_class"`
	LdapTeamAdminAttr          string `toml:"ldap_team_admin_attr"`
	LdapAllowedAppAttr         string `toml:"ldap_allowed_app_attr"`
	LdapAllowedAppCommonName   string `toml:"ldap_allowed_app_common_name"`
	LdapUserCommonNamePrefix   string `toml:"ldap_user_common_name_prefix"`
	LdapTeamCommonNamePrefix   string `toml:"ldap_team_common_name_prefix"`
	LdapUserClass              string `toml:"ldap_user_class"`
	LdapUserClassAttr          string `toml:"ldap_user_class_attr"`
	SkipAuthorization          bool   `toml:"skip_authorization"`
	LdapSuperUserGroup         string `toml:"ldap_super_user_group"`
	Host                       string `toml:"host"`
	Region                     string `toml:"region"`
	Zone                       string `toml:"zone"`
	AvailableZones             string `toml:"available_zones"`
	MaintenanceFile            string `toml:"maintenance_file"`
	MaintenanceCheckInterval   string `toml:"maintenance_check_interval"`
	SuperUserOnlyFile          string `toml:"superuser_only_file"`
	SuperUserOnlyCheckInterval string `toml:"superuser_only_check_interval"`
	MinRouterPort              uint16 `toml:"min_router_port"`
	MaxRouterPort              uint16 `toml:"max_router_port"`
	SMTPAddr                   string `toml:"smtp_addr"`
	SMTPFrom                   string `toml:"smtp_from"`
	SMTPCC                     string `toml:"smtp_cc"`
}

type ServerOpts struct {
	RpcAddr                    string `long:"rpc" description:"the RPC listen addr"`
	SupervisorPort             uint16 `long:"supervisor" description:"the RPC port for supervisor"`
	ApiAddr                    string `long:"api" description:"the API listen addr"`
	ZookeeperUri               string `long:"zookeeper" description:"the uri of the zookeeper to connect to"`
	ConfigFile                 string `long:"config-file" default:"/etc/atlantis/manager/server.toml" description:"the config file to use"`
	LdapHost                   string `long:"ldap-host" description:"LDAP server to contact"`
	LdapPort                   uint16 `long:"ldap-port" description:"LDAP port to use"`
	LdapBaseDomain             string `long:"ldap-base-domain" description:"LDAP Base Domain Name to use"`
	CPUSharesIncrement         uint   `long:"cpu-shares-increment" description:"CPU shares increment"`
	MemoryLimitIncrement       uint   `long:"memory-limit-increment" description:"Memory Limit increment"`
	ResultDuration             string `long:"result-duration" description:"How long to keep the results of an Async Command"`
	SkipAuthorization          bool   `long:"skip-authorization" description:"Skip verification for LDAP UTA Details"`
	Host                       string `long:"host" description:"the host of this manager"`
	Region                     string `long:"region" description:"the region this manager is in"`
	Zone                       string `long:"zone" description:"the availability zone this manager is in"`
	AvailableZones             string `long:"available-zones" description:"the available availability zones"`
	MaintenanceFile            string `long:"maintenance-file" description:"the maintenance file to check"`
	MaintenanceCheckInterval   string `long:"maintenance-check-interval" description:"the interval to check the maintenance file"`
	SuperUserOnlyFile          string `long:"superuser-only-file" description:"if this file exists, only allow superusers to do things"`
	SuperUserOnlyCheckInterval string `long:"maintenance-check-interval" description:"the interval to check the superuser only file"`
	MinRouterPort              uint16 `long:"min-router-port"`
	MaxRouterPort              uint16 `long:"max-router-port"`
	SMTPAddr                   string `long:"smtp-addr"`
	SMTPFrom                   string `long:"smtp-from"`
	SMTPCC                     string `long:"smtp-cc"`
}

type ManagerServer struct {
	parser *flags.Parser
	Opts   *ServerOpts
	Config *ServerConfig
}

func New() *ManagerServer {
	opts := &ServerOpts{}
	manager := &ManagerServer{
		parser: flags.NewParser(opts, flags.Default),
		Opts:   opts,
		Config: &ServerConfig{
			RpcAddr:                    fmt.Sprintf(":%d", DefaultManagerRPCPort),
			SupervisorPort:             iconst.DefaultSupervisorRPCPort,
			ApiAddr:                    fmt.Sprintf(":%d", DefaultManagerAPIPort),
			LdapPort:                   DefaultLDAPPort,
			ZookeeperUri:               "localhost:2181",
			CPUSharesIncrement:         1,
			MemoryLimitIncrement:       1,
			ResultDuration:             DefaultResultDuration,
			SkipAuthorization:          false,
			Host:                       DefaultManagerHost,
			Region:                     DefaultRegion,
			Zone:                       DefaultZone,
			AvailableZones:             DefaultZone,
			MaintenanceFile:            DefaultMaintenanceFile,
			MaintenanceCheckInterval:   DefaultMaintenanceCheckInterval,
			SuperUserOnlyFile:          DefaultSuperUserOnlyFile,
			SuperUserOnlyCheckInterval: DefaultSuperUserOnlyCheckInterval,
			MinRouterPort:              DefaultMinRouterPort,
			MaxRouterPort:              DefaultMaxRouterPort,
			SMTPAddr:                   "",
			SMTPFrom:                   "",
			SMTPCC:                     "",
		},
	}
	manager.parser.Parse()
	manager.overlayConfig()
	return manager
}

func (m *ManagerServer) SetHandlerFunc(handlerFunc func(http.Handler) http.Handler) {
	api.HandlerFunc = handlerFunc
}

func (m *ManagerServer) SetDNSProvider(provider dns.DNSProvider) {
	dns.Provider = provider
}

func (m *ManagerServer) AddCommand(cmd, desc, longDesc string, data interface{}) {
	log.Fatalln("You may not add a command to the manager server")
}

func (m *ManagerServer) Run(bldr builder.Builder) {
	builder.DefaultBuilder = bldr
	smtp.Init(m.Config.SMTPAddr, m.Config.SMTPFrom, strings.Split(m.Config.SMTPCC, ","))
	if err := crypto.Init(); err != nil {
		panic(fmt.Sprintf("Error initializing crypto: %s", err.Error()))
	}
	log.Println("Fate rarely calls upon us at a moment of our choosing.")
	log.Println("                                                       -- Manager\n")
	ldap.Init(m.Config.LdapHost, m.Config.LdapPort, m.Config.LdapBaseDomain)
	Host = m.Config.Host
	Region = m.Config.Region
	Zone = m.Config.Zone
	AvailableZones = strings.Split(m.Config.AvailableZones, ",")
	log.Printf("Initializing Manager [%s] [%s] [%s]", Region, Zone, Host)
	datamodel.Init(m.Config.ZookeeperUri)
	datamodel.MinRouterPort = m.Config.MinRouterPort
	datamodel.MaxRouterPort = m.Config.MaxRouterPort
	resultDuration, err := time.ParseDuration(m.Config.ResultDuration)
	if err != nil {
		panic(fmt.Sprintf("Could not parse Result Duration: %s", err.Error()))
	}
	handleError(rpc.Init(m.Config.RpcAddr, m.Config.SupervisorPort, m.Config.CPUSharesIncrement,
		m.Config.MemoryLimitIncrement, resultDuration))
	handleError(api.Init(m.Config.ApiAddr))
	err = m.LDAPInit()
	if err != nil {
		log.Fatalln(err)
	}
	maintenanceCheckInterval, err := time.ParseDuration(m.Config.MaintenanceCheckInterval)
	if err != nil {
		log.Fatalln(err)
	}
	superUserCheckInterval, err := time.ParseDuration(m.Config.SuperUserOnlyCheckInterval)
	if err != nil {
		log.Fatalln(err)
	}
	MaintenanceChecker(m.Config.MaintenanceFile, maintenanceCheckInterval)
	rpc.SuperUserOnlyChecker(m.Config.SuperUserOnlyFile, superUserCheckInterval)
	go signalListener()
	go rpc.Listen()
	api.Listen()
}

func handleError(err error) {
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
}

func (m *ManagerServer) overlayConfig() {
	if m.Opts.ConfigFile != "" {
		_, err := toml.DecodeFile(m.Opts.ConfigFile, m.Config)
		if err != nil {
			log.Println(err)
			// no need to panic here. we have reasonable defaults.
		}
	}
	if m.Opts.RpcAddr != "" {
		m.Config.RpcAddr = m.Opts.RpcAddr
	}
	if m.Opts.ApiAddr != "" {
		m.Config.ApiAddr = m.Opts.ApiAddr
	}
	if m.Opts.SupervisorPort != 0 {
		m.Config.SupervisorPort = m.Opts.SupervisorPort
	}
	if m.Opts.ZookeeperUri != "" {
		m.Config.ZookeeperUri = m.Opts.ZookeeperUri
	}
	if m.Opts.LdapPort != 0 {
		m.Config.LdapPort = m.Opts.LdapPort
	}
	if m.Opts.LdapHost != "" {
		m.Config.LdapHost = m.Opts.LdapHost
	}
	if m.Opts.LdapBaseDomain != "" {
		m.Config.LdapBaseDomain = m.Opts.LdapBaseDomain
	}
	if m.Opts.CPUSharesIncrement != 0 {
		m.Config.CPUSharesIncrement = m.Opts.CPUSharesIncrement
	}
	if m.Opts.MemoryLimitIncrement != 0 {
		m.Config.MemoryLimitIncrement = m.Opts.MemoryLimitIncrement
	}
	if m.Opts.ResultDuration != "" {
		m.Config.ResultDuration = m.Opts.ResultDuration
	}
	if m.Opts.Host != "" {
		m.Config.Host = m.Opts.Host
	}
	if m.Opts.Region != "" {
		m.Config.Region = m.Opts.Region
	}
	if m.Opts.Zone != "" {
		m.Config.Zone = m.Opts.Zone
	}
	if m.Opts.AvailableZones != "" {
		m.Config.AvailableZones = m.Opts.AvailableZones
	}
	if m.Opts.MaintenanceFile != "" {
		m.Config.MaintenanceFile = m.Opts.MaintenanceFile
	}
	if m.Opts.MaintenanceCheckInterval != "" {
		m.Config.MaintenanceCheckInterval = m.Opts.MaintenanceCheckInterval
	}
	if m.Opts.SuperUserOnlyFile != "" {
		m.Config.SuperUserOnlyFile = m.Opts.SuperUserOnlyFile
	}
	if m.Opts.SuperUserOnlyCheckInterval != "" {
		m.Config.SuperUserOnlyCheckInterval = m.Opts.SuperUserOnlyCheckInterval
	}
	if m.Opts.MinRouterPort != 0 {
		m.Config.MinRouterPort = m.Opts.MinRouterPort
	}
	if m.Opts.MaxRouterPort != 0 {
		m.Config.MaxRouterPort = m.Opts.MaxRouterPort
	}
	if m.Opts.SMTPAddr != "" {
		m.Config.SMTPAddr = m.Opts.SMTPAddr
	}
	if m.Opts.SMTPFrom != "" {
		m.Config.SMTPFrom = m.Opts.SMTPFrom
	}
	if m.Opts.SMTPCC != "" {
		m.Config.SMTPCC = m.Opts.SMTPCC
	}
}

func (m *ManagerServer) LDAPInit() error {
	if m.Config.SkipAuthorization == false {
		if m.Config.LdapUserCommonNamePrefix == "" {
			return errors.New("Missing in server.toml: ldap_user_common_name_prefix")
		}
		if m.Config.LdapTeamCommonNamePrefix == "" {
			return errors.New("Missing in server.toml: ldap_team_common_name_prefix")
		}
		if m.Config.LdapUserSearchBase == "" || m.Config.LdapTeamSearchBase == "" ||
			m.Config.LdapAppClass == "" || m.Config.LdapUsernameAttr == "" ||
			m.Config.LdapTeamClass == "" || m.Config.LdapAllowedAppAttr == "" ||
			m.Config.LdapTeamAdminAttr == "" {
			return errors.New("LDAP User / Team / Applications Authorization Details Insufficient")
		}
	}
	ldap.UserOu = m.Config.LdapUserSearchBase
	ldap.RoleUserOu = m.Config.LdapRoleUserSearchBase
	ldap.TeamOu = m.Config.LdapTeamSearchBase
	ldap.UsernameAttr = m.Config.LdapUsernameAttr
	ldap.TeamClass = m.Config.LdapTeamClass
	ldap.TeamAdminAttr = m.Config.LdapTeamAdminAttr
	ldap.AllowedAppAttr = m.Config.LdapAllowedAppAttr
	ldap.AppClass = m.Config.LdapAppClass
	ldap.UserCommonName = m.Config.LdapUserCommonNamePrefix
	ldap.TeamCommonName = m.Config.LdapTeamCommonNamePrefix
	ldap.SkipAuthorization = m.Config.SkipAuthorization
	ldap.SuperUserGroup = m.Config.LdapSuperUserGroup
	ldap.UserClass = m.Config.LdapUserClass
	ldap.UserClassAttr = m.Config.LdapUserClassAttr
	return nil
}

func signalListener() {
	// wait for SIGTERM
	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGTERM)
	<-termChan
	signal.Stop(termChan)
	close(termChan)

	// wait indefinitely for idle before exit - we can always kill if we *really* want manager to die
	log.Println("[SIGTERM] Gracefully shutting down...")
	for !Tracker.Idle(nil) {
		log.Println("[SIGTERM] -> waiting for idle")
		time.Sleep(5 * time.Second)
	}
	os.Exit(0)
}
