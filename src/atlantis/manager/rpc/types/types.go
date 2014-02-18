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

package types

import (
	"atlantis/router/config"
	. "atlantis/supervisor/rpc/types"
)

type App struct {
	NonAtlantis     bool
	Internal        bool // atlantis apps only
	Name            string
	Email           string
	Repo            string // atlantis apps only
	Root            string // atlantis apps only
	DependerEnvData map[string]*DependerEnvData
	DependerAppData map[string]*DependerAppData
}

type DependerEnvData struct {
	Name          string
	SecurityGroup []string
	EncryptedData string
	DataMap       map[string]interface{} `json:",omitempty"`
}

type DependerAppData struct {
	Name            string
	DependerEnvData map[string]*DependerEnvData
}

type Router struct {
	Internal  bool
	Zone      string
	Host      string
	CName     string
	IP        string
	RecordIDs []string
}

type AppEnv struct {
	App string
	Env string
}

func (a AppEnv) String() string {
	return a.App + "." + a.Env
}

type RouterPorts struct {
	Internal  bool
	PortMap   map[string]AppEnv
	AppEnvMap map[string]string
}

type Manager struct {
	Region           string
	Host             string
	ManagerCName     string
	ManagerRecordID  string
	RegistryCName    string
	RegistryRecordID string
	Roles            map[string]map[string]bool
}

// Manager RPC Types

// ------------ Health Check ------------
// Used to check the health and stats of Manager
type ManagerHealthCheckArg struct {
}

type ManagerHealthCheckReply struct {
	Region string
	Zone   string
	Status string
}

// ------------ Register Supervisor ------------
// Used to register an Supervisor
type ManagerRegisterSupervisorArg struct {
	ManagerAuthArg
	Host string
}

type ManagerRegisterSupervisorReply struct {
	Status string
}

// ------------ Register Manager ------------
// Used to register an Manager
type ManagerRegisterManagerArg struct {
	ManagerAuthArg
	Region        string
	Host          string
	RegistryCName string
	ManagerCName  string
}

type ManagerRegisterManagerReply struct {
	Status  string
	Manager *Manager
}

// get a manager
type ManagerGetManagerArg struct {
	ManagerAuthArg
	Region string
	Host   string
}

type ManagerGetManagerReply struct {
	Status  string
	Manager *Manager
}

type ManagerGetSelfArg struct {
	ManagerAuthArg
}

// Used to manager Manager Roles
type ManagerRoleArg struct {
	ManagerAuthArg
	Region string
	Host   string
	Role   string
	Type   string
}

type ManagerRoleReply struct {
	Status  string
	Manager *Manager
}

type ManagerHasRoleReply struct {
	Status  string
	HasRole bool
}

// ------------ Register App ------------
// Used to register an App
type ManagerRegisterAppArg struct {
	ManagerAuthArg
	NonAtlantis bool
	Internal    bool
	Name        string
	Repo        string
	Root        string
	Email       string
}

type ManagerRegisterAppReply struct {
	Status string
}

// ------------ Get App ------------
// Used to get an App
type ManagerGetAppArg struct {
	ManagerAuthArg
	Name string
}

type ManagerGetAppReply struct {
	Status string
	App    *App
}

// ------------ Request App Dependency ------------
// User for RquestAppDependency
type ManagerRequestAppDependencyArg struct {
	ManagerAuthArg
	App        string
	Dependency string
	Envs       []string
}

type ManagerRequestAppDependencyReply struct {
	Status string
}

// ------------ Add Depender App Data ------------
// Used Add AppDependerData to an App
type ManagerAddDependerAppDataArg struct {
	ManagerAuthArg
	App             string
	DependerAppData *DependerAppData
}

type ManagerAddDependerAppDataReply struct {
	Status string
	App    *App
}

// ------------ Remove Depender App Data ------------
// Used Remove AppDependerData from an App
type ManagerRemoveDependerAppDataArg struct {
	ManagerAuthArg
	App      string
	Depender string
}

type ManagerRemoveDependerAppDataReply struct {
	Status string
	App    *App
}

// ------------ Get Depender App Data ------------
// Used Get AppDependerData from an App
type ManagerGetDependerAppDataArg struct {
	ManagerAuthArg
	App      string
	Depender string
}

type ManagerGetDependerAppDataReply struct {
	Status          string
	DependerAppData *DependerAppData
}

// ------------ Add Depender Env Data ------------
// Used Add EnvDependerData to an App
type ManagerAddDependerEnvDataArg struct {
	ManagerAuthArg
	App             string
	DependerEnvData *DependerEnvData
}

type ManagerAddDependerEnvDataReply struct {
	Status string
	App    *App
}

// ------------ Remove Depender Env Data ------------
// Used Remove EnvDependerData from an App
type ManagerRemoveDependerEnvDataArg struct {
	ManagerAuthArg
	App string
	Env string
}

type ManagerRemoveDependerEnvDataReply struct {
	Status string
	App    *App
}

// ------------ Get Depender Env Data ------------
// Used Get EnvDependerData from an App
type ManagerGetDependerEnvDataArg struct {
	ManagerAuthArg
	App string
	Env string
}

type ManagerGetDependerEnvDataReply struct {
	Status          string
	DependerEnvData *DependerEnvData
}

// ------------ Add Depender Env Data For a Depender App ------------
// Used Add EnvDependerData to a DependerApp in an App
type ManagerAddDependerEnvDataForDependerAppArg struct {
	ManagerAuthArg
	App             string
	Depender        string
	DependerEnvData *DependerEnvData
}

type ManagerAddDependerEnvDataForDependerAppReply struct {
	Status string
	App    *App
}

// ------------ Remove Depender Env Data For a Depender App ------------
// Used Remove EnvDependerData from a DependerApp in an App
type ManagerRemoveDependerEnvDataForDependerAppArg struct {
	ManagerAuthArg
	App      string
	Depender string
	Env      string
}

type ManagerRemoveDependerEnvDataForDependerAppReply struct {
	Status string
	App    *App
}

// ------------ Get Depender Env Data For a Depender App ------------
// Used Get EnvDependerData from a DependerApp in an App
type ManagerGetDependerEnvDataForDependerAppArg struct {
	ManagerAuthArg
	App      string
	Depender string
	Env      string
}

type ManagerGetDependerEnvDataForDependerAppReply struct {
	Status          string
	DependerEnvData *DependerEnvData
}

// ------------ ListRegisteredApps ------------
// List all apps
type ManagerListRegisteredAppsArg struct {
	ManagerAuthArg
}

type ManagerListRegisteredAppsReply struct {
	Apps   []string
	Status string
}

// ------------ Register Router ------------
// Used to register an Router
type ManagerRegisterRouterArg struct {
	ManagerAuthArg
	Internal bool
	Zone     string
	Host     string
	IP       string
}

type ManagerRegisterRouterReply struct {
	Router *Router
	Status string
}

// ------------ Get Router ------------
// Used to get an Router
type ManagerGetRouterArg struct {
	ManagerAuthArg
	Internal bool
	Zone     string
	Host     string
}

type ManagerGetRouterReply struct {
	Status string
	Router *Router
}

// ------------ ListRouters ------------
// List all apps
type ManagerListRoutersArg struct {
	ManagerAuthArg
	Internal bool
}

type ManagerListRoutersReply struct {
	Routers map[string][]string
	Status  string
}

// ------------ List Supervisors ------------
// Used to list available Supervisors
type ManagerListSupervisorsArg struct {
	ManagerAuthArg
}

type ManagerListSupervisorsReply struct {
	Supervisors []string
	Status      string
}

// ------------ List Managers ------------
// Used to list available Managers
type ManagerListManagersArg struct {
	ManagerAuthArg
}

type ManagerListManagersReply struct {
	Managers map[string][]string
	Status   string
}

// ------------ Manager Dependency Management ------------
// Used to update, get, or delete a dependency
type ManagerDepArg struct {
	ManagerAuthArg
	Env   string
	Name  string
	Value string
}

type ManagerDepReply struct {
	Value  string
	Status string
}

// ------------ Manager Environment Management ------------
// Used to update, get, or delete a environment
type ManagerEnvArg struct {
	ManagerAuthArg
	Name string
}

type ManagerEnvReply struct {
	Status string
}

// ------------ Deploy ------------
// Used to deploy an app+sha+env
type ManagerDeployArg struct {
	ManagerAuthArg
	App         string
	Sha         string
	Env         string
	Instances   uint
	CPUShares   uint // relative shares
	MemoryLimit uint // MBytes
	Dev         bool // if true, only install 1 instance in 1 zone
}

type ManagerDeployReply struct {
	Status     string
	Containers []*Container
}

// ------------ CopyContainer ------------
// Used to deploy by copying a container
type ManagerCopyContainerArg struct {
	ManagerAuthArg
	ContainerID string
	Instances   uint
}

// CopyContainer uses ManagerDeployReply

// ------------ MoveContainer ------------
// Used to deploy by copying a container
type ManagerMoveContainerArg struct {
	ManagerAuthArg
	ContainerID string
}

// MoveContainer uses ManagerDeployReply

// ------------ CopyOrphaned ------------
// Used to deploy by copying a container
type ManagerCopyOrphanedArg struct {
	ManagerAuthArg
	ContainerID string
	Host        string
	CleanupZk   bool
}

// CopyOrphaned uses ManagerDeployReply

// ------------ ResolveDeps ------------
// used to resolve deps in an environment to see what the deploy will contain
type ManagerResolveDepsArg struct {
	ManagerAuthArg
	App      string
	Env      string
	DepNames []string
}

type ManagerResolveDepsReply struct {
	Status string
	Deps   map[string]DepsType
}

// ------------ Teardown ------------
// Teardown containers by app, app+sha, app+sha+container, or just simply all
type ManagerTeardownArg struct {
	ManagerAuthArg
	App         string
	Sha         string
	Env         string
	ContainerID string
	All         bool
}

type ManagerTeardownReply struct {
	ContainerIDs []string
	Status       string
}

// ------------ GetContainer ------------
// Used to get a container
type ManagerGetContainerArg struct {
	ManagerAuthArg
	ContainerID string
}

type ManagerGetContainerReply struct {
	Container *Container
	Status    string
}

// ------------ ListContainers ------------
// List all containers that are part of the app+sha+env combo
type ManagerListContainersArg struct {
	ManagerAuthArg
	App string
	Sha string
	Env string
}

type ManagerListContainersReply struct {
	ContainerIDs []string
	Status       string
}

// ------------ ListEnvs ------------
// List all envs that are part of the app+sha combo
type ManagerListEnvsArg struct {
	ManagerAuthArg
	App string
	Sha string
}

type ManagerListEnvsReply struct {
	Envs   []string
	Status string
}

// ------------ ListShas ------------
// List all shas that are part of the app
type ManagerListShasArg struct {
	ManagerAuthArg
	App string
}

type ManagerListShasReply struct {
	Shas   []string
	Status string
}

// ------------ ListApps ------------
// List all apps
type ManagerListAppsArg struct {
	ManagerAuthArg
}

type ManagerListAppsReply struct {
	Apps   []string
	Status string
}

// ------------ UpdatePool ------------
// used to update pools
type ManagerUpdatePoolArg struct {
	ManagerAuthArg
	Pool config.Pool
}

type ManagerUpdatePoolReply struct {
	Status string
}

// ------------ DeletePool ------------
// used to delete a pool
type ManagerDeletePoolArg struct {
	ManagerAuthArg
	Name     string
	Internal bool
}

type ManagerDeletePoolReply struct {
	Status string
}

// ------------ GetPool ------------
// used to get a pool
type ManagerGetPoolArg struct {
	ManagerAuthArg
	Name     string
	Internal bool
}

type ManagerGetPoolReply struct {
	Pool   config.Pool
	Status string
}

// ------------ ListPools ------------
// used to list all pools
type ManagerListPoolsArg struct {
	ManagerAuthArg
	Internal bool
}

type ManagerListPoolsReply struct {
	Pools  []string
	Status string
}

// ------------ UpdateRule ------------
// used to update rules
type ManagerUpdateRuleArg struct {
	ManagerAuthArg
	Rule config.Rule
}

type ManagerUpdateRuleReply struct {
	Status string
}

// ------------ DeleteRule ------------
// used to delete a rule
type ManagerDeleteRuleArg struct {
	ManagerAuthArg
	Name     string
	Internal bool
}

type ManagerDeleteRuleReply struct {
	Status string
}

// ------------ GetRule ------------
// used to get a pool
type ManagerGetRuleArg struct {
	ManagerAuthArg
	Name     string
	Internal bool
}

type ManagerGetRuleReply struct {
	Rule   config.Rule
	Status string
}

// ------------ ListRules ------------
// used to list all pools
type ManagerListRulesArg struct {
	ManagerAuthArg
	Internal bool
}

type ManagerListRulesReply struct {
	Rules  []string
	Status string
}

// ------------ UpdateTrie ------------
// used to update tries
type ManagerUpdateTrieArg struct {
	ManagerAuthArg
	Trie config.Trie
}

type ManagerUpdateTrieReply struct {
	Status string
}

// ------------ DeleteTrie ------------
// used to delete a trie
type ManagerDeleteTrieArg struct {
	ManagerAuthArg
	Name     string
	Internal bool
}

type ManagerDeleteTrieReply struct {
	Status string
}

// ------------ GetTrie ------------
// used to get a pool
type ManagerGetTrieArg struct {
	ManagerAuthArg
	Name     string
	Internal bool
}

type ManagerGetTrieReply struct {
	Trie   config.Trie
	Status string
}

// ------------ ListTries ------------
// used to list all pools
type ManagerListTriesArg struct {
	ManagerAuthArg
	Internal bool
}

type ManagerListTriesReply struct {
	Tries  []string
	Status string
}

// ------------ GetAppEnvPort ------------
// used to get a port
type ManagerGetAppEnvPortArg struct {
	ManagerAuthArg
	App string
	Env string
}

type ManagerGetAppEnvPortReply struct {
	Port   config.Port
	Status string
}

// ------------ ListAppEnvsWithPort ------------
// used to list all app envs with ports
type ManagerListAppEnvsWithPortArg struct {
	ManagerAuthArg
	Internal bool
}

type ManagerListAppEnvsWithPortReply struct {
	AppEnvs []AppEnv
	Status  string
}

// ------------ UpdatePort ------------
// used to update ports
type ManagerUpdatePortArg struct {
	ManagerAuthArg
	Port config.Port
}

type ManagerUpdatePortReply struct {
	Status string
}

// ------------ DeletePort ------------
// used to delete a port
type ManagerDeletePortArg struct {
	ManagerAuthArg
	Port     uint16
	Internal bool
}

type ManagerDeletePortReply struct {
	Status string
}

// ------------ GetPort ------------
// used to get a port
type ManagerGetPortArg struct {
	ManagerAuthArg
	Port     uint16
	Internal bool
}

type ManagerGetPortReply struct {
	Port   config.Port
	Status string
}

// ------------ ListPorts ------------
// used to list all ports
type ManagerListPortsArg struct {
	ManagerAuthArg
	Internal bool
}

type ManagerListPortsReply struct {
	Ports  []uint16
	Status string
}

// ------------ Login -----------
// used for LDAP Logins
type ManagerLoginArg struct {
	User   string
	Pass   string
	Secret string
}

type ManagerLoginReply struct {
	LoggedIn bool
	Secret   string
}

// ------------ Group ----------
// used for group-user mappings
type ManagerUserMapArg struct {
	User  string
	Group string
}

// ------------ Team ----------
// used for add/removing teams
type ManagerTeamArg struct {
	ManagerAuthArg
	Team string
}

type ManagerTeamReply struct {
}

// ------------ ListTeams ----------
// used for listing teams
type ManagerListTeamsArg struct {
	ManagerAuthArg
}

type ManagerListTeamsReply struct {
	Teams []string
}

// ------------ ListTeamEmails ----------
// used for listing team Emails
type ManagerListTeamEmailsArg struct {
	ManagerAuthArg
	Team string
}

type ManagerListTeamEmailsReply struct {
	TeamEmails []string
}

// ------------ ListTeamAdmins ----------
// used for listing team Admins
type ManagerListTeamAdminsArg struct {
	ManagerAuthArg
	Team string
}

type ManagerListTeamAdminsReply struct {
	TeamAdmins []string
}

// ------------ ListTeamMembers ----------
// used for listing team Members
type ManagerListTeamMembersArg struct {
	ManagerAuthArg
	Team string
}

type ManagerListTeamMembersReply struct {
	TeamMembers []string
}

// ------------ ListTeamApps ----------
// used for listing team Apps
type ManagerListTeamAppsArg struct {
	ManagerAuthArg
	Team string
}

type ManagerListTeamAppsReply struct {
	TeamApps []string
}

// ------------ Team Member----------
// used for add/removing team members
type ManagerTeamMemberArg struct {
	ManagerAuthArg
	Team string
	User string
}

type ManagerTeamMemberReply struct {
}

// -- App (owned by Teams) --
// used for adding/removing apps to teams
type ManagerAppArg struct {
	ManagerAuthArg
	App  string
	Team string
}

type ManagerAppReply struct {
}

// ------------ Authorize SSH ------------
// Authorize SSH
type ManagerAuthorizeSSHArg struct {
	ContainerID string
	User        string
	PublicKey   string
}

type ManagerAuthorizeSSHReply struct {
	Host   string
	Port   uint16
	Status string
}

// ------------ Deauthorize SSH ------------
// Deauthorize SSH
type ManagerDeauthorizeSSHArg struct {
	ContainerID string
	User        string
}

type ManagerDeauthorizeSSHReply struct {
	Status string
}

// ------------ Is App Allowed ---------
// Check whether an user has permissions to an app
type ManagerIsAppAllowedArg struct {
	ManagerAuthArg
	App  string
	User string // optional. superusers can check for another user
}

type ManagerIsAppAllowedReply struct {
	IsAllowed bool
}

// ------------ List Allowed Apps ---------
// List all allowed apps for this user
type ManagerListAllowedAppsArg struct {
	ManagerAuthArg
	User string // optional. superusers can check for another user
}

type ManagerListAllowedAppsReply struct {
	Apps []string
}

// --------- Team Admin Permissions -------
// Check whether one is the team admin
type ManagerTeamAdminArg struct {
	ManagerAuthArg
	Team string
}

type ManagerTeamAdminReply struct {
	IsAdmin bool
}

type ManagerModifyTeamAdminArg struct {
	ManagerAuthArg
	Team string
	User string
}

type ManagerModifyTeamAdminReply struct {
}

// -------- Team Email -----------
// Adding/Deleting Team Emails
type ManagerEmailArg struct {
	ManagerAuthArg
	Team  string
	Email string
}

type ManagerEmailReply struct {
}

// -------- Super User Permissions ---------
// Check whether one is in the super user group
type ManagerSuperUserArg struct {
	ManagerAuthArg
}

type ManagerSuperUserReply struct {
	IsSuperUser bool
}

// ------------ ContainerMaintenance ------------
// Set Container Maintenance Mode
type ManagerContainerMaintenanceArg struct {
	ManagerAuthArg
	ContainerID string
	Maintenance bool
}

type ManagerContainerMaintenanceReply struct {
	Status string
}

// ------------ Idle ------------
// Check if Idle
type ManagerIdleArg struct {
}

type ManagerIdleReply struct {
	Idle   bool
	Status string
}

// ------- Authentication --------
// Used for authenticating and accessing current user's sessions
type ManagerAuthArg struct {
	User     string
	Password string
	Secret   string
}

func (o *ManagerAuthArg) Credentials() (user, password, secret string) {
	return o.User, o.Password, o.Secret
}
