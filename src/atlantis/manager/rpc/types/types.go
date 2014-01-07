package types

import (
	"atlantis/router/config"
	. "atlantis/supervisor/rpc/types"
)

type App struct {
	NonAtlantis         bool
	Type                string // proxy type. http for atlantis apps, tcp or http for non-atlantis apps
	Name                string
	Repo                string            // atlantis apps only
	Root                string            // atlantis apps only
	Addrs               map[string]string // non-atlantis apps only. env -> addr
	Email               string
	AllowedDependerApps map[string]bool
}

type Router struct {
	Internal  bool
	Zone      string
	Host      string
	CName     string
	RecordIDs []string
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
	Type        string
	Name        string
	Repo        string
	Root        string
	Addrs       map[string]string
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

// ------------ Depender App ------------
// Used to get an App
type ManagerDependerAppArg struct {
	ManagerAuthArg
	Depender string
	Dependee string
}

type ManagerDependerAppReply struct {
	Status   string
	Dependee *App
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
	Name   string
	Parent string
}

type ManagerEnvReply struct {
	Parent       string
	Deps         map[string]string
	ResolvedDeps map[string]string
	Status       string
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
	Deps   map[string]map[string]string
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

// ------------ App Permissions ---------
// Check whether an user has permissions to an app
type ManagerAppPermissionsArg struct {
	ManagerAuthArg
	App string
}

type ManagerAppPermissionsReply struct {
	Permission bool
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

// ------------ UpdateProxy ------------
// Update Proxy Version
type ManagerUpdateProxyArg struct {
	ManagerAuthArg
	Sha string
}

type ManagerUpdateProxyReply struct {
	Status string
}

// ------------ ConfigureProxy ------------
// Configure Proxy Version
type ManagerConfigureProxyArg struct {
	ManagerAuthArg
}

type ManagerConfigureProxyReply struct {
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
