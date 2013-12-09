package client

import (
	. "atlantis/manager/constant"
	"atlantis/manager/rpc/client"
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/jigish/go-flags"
	"log"
	"os"
	"reflect"
)

func Log(format string, args ...interface{}) {
	if !IsJson() && !clientOpts.Quiet {
		log.Printf(format, args...)
	}
}

func Fatal(format string, args ...interface{}) {
	if !IsJson() && !clientOpts.Quiet {
		log.Fatalf(format, args...)
	}
}

func IsQuiet() bool {
	return clientOpts.Quiet
}

func IsJson() bool {
	return clientOpts.Json || clientOpts.PrettyJson
}

func quietOutput(prefix string, val interface{}) {
	quietValue := val
	indVal := reflect.Indirect(reflect.ValueOf(val))
	if kind := indVal.Kind(); kind == reflect.Struct {
		quietValue = destructify(val)
	}
	switch t := quietValue.(type) {
	case bool:
		fmt.Printf("%s%t\n", prefix, t)
	case int:
		fmt.Printf("%s%d\n", prefix, t)
	case uint:
		fmt.Printf("%s%d\n", prefix, t)
	case uint16:
		fmt.Printf("%s%d\n", prefix, t)
	case float64:
		fmt.Printf("%s%f\n", prefix, t)
	case string:
		fmt.Printf("%s%s\n", prefix, t)
	case []string:
		for _, value := range t {
			quietOutput(prefix, value)
		}
	case []interface{}:
		for _, value := range t {
			quietOutput(prefix, destructify(value))
		}
	case map[string]string:
		for key, value := range t {
			quietOutput(prefix+key+" ", value)
		}
	case map[string][]string:
		for key, value := range t {
			quietOutput(prefix+key+" ", value)
		}
	case map[string]interface{}:
		for key, value := range t {
			quietOutput(prefix+key+" ", destructify(value))
		}
	default:
		panic(fmt.Sprintf("invalid quiet type %T", t))
	}
}

func destructify(val interface{}) interface{} {
	indVal := reflect.Indirect(reflect.ValueOf(val))
	if kind := indVal.Kind(); kind == reflect.Struct {
		typ := indVal.Type()
		mapVal := map[string]interface{}{}
		for i := 0; i < typ.NumField(); i++ {
			field := indVal.Field(i)
			mapVal[typ.Field(i).Name] = destructify(field.Interface())
		}
		return mapVal
	} else if kind == reflect.Array || kind == reflect.Slice {
		if k := indVal.Type().Elem().Kind(); k != reflect.Array && k != reflect.Slice && k != reflect.Map &&
			k != reflect.Struct {
			return indVal.Interface()
		}
		arr := make([]interface{}, indVal.Len())
		for i := 0; i < indVal.Len(); i++ {
			field := indVal.Index(i)
			arr[i] = destructify(field.Interface())
		}
		return arr
	} else if kind == reflect.Map {
		if k := indVal.Type().Elem().Kind(); k != reflect.Array && k != reflect.Slice && k != reflect.Map &&
			k != reflect.Struct {
			return indVal.Interface()
		}
		keys := indVal.MapKeys()
		mapVal := make(map[string]interface{}, len(keys))
		for _, key := range keys {
			field := indVal.MapIndex(key)
			mapVal[fmt.Sprintf("%v", key.Interface())] = destructify(field.Interface())
		}
		return mapVal
	} else {
		return val
	}
}

func Output(obj map[string]interface{}, quiet interface{}, err error) error {
	if !IsJson() && !clientOpts.Quiet {
		return err
	}
	if IsJson() && obj != nil {
		obj["error"] = err
		var bytes []byte
		if clientOpts.PrettyJson {
			bytes, err = json.MarshalIndent(obj, "", "  ")
		} else {
			bytes, err = json.Marshal(obj)
		}
		if err != nil {
			fmt.Printf("{\"error\":\"%s\"}\n", err.Error())
		} else {
			fmt.Printf("%s\n", bytes)
		}
	} else if clientOpts.Quiet && quiet != nil {
		quietOutput("", quiet)
	}
	if err != nil { // denote failure with non-zero exit code
		os.Exit(1)
	}
	return nil
}

func OutputError(err error) error {
	return Output(map[string]interface{}{}, nil, err)
}

func OutputEmpty() error {
	return Output(map[string]interface{}{}, nil, nil)
}

// this will scan through checkArgs. If one of the elements is nil or empty it will pop off the first arg and
// use that as the value of the element. returns the resulting args slice
func ExtractArgs(checkArgs []*string, args []string) []string {
	for _, checkArg := range checkArgs {
		if len(args) == 0 {
			return args
		}
		if checkArg == nil || *checkArg == "" {
			*checkArg = args[0]
			args = args[1:]
		}
	}
	return args
}

type ClientConfig struct {
	Host    string `toml:"host"`
	Port    uint16 `toml:"port"`
	KeyPath string `toml:"key_path"`
}

func (c *ClientConfig) RPCHostAndPort() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

type ClientOpts struct {
	// Only use capital letters here. Also, "H" is off limits. kthxbye.
	Host       string `short:"M" long:"manager-host" description:"the manager host"`
	Port       uint16 `short:"P" long:"manager-port" description:"the manager port"`
	Config     string `short:"F" long:"config-file" default:"" description:"the config file to use"`
	Region     string `short:"R" long:"manager-region" default:"us-east-1" description:"the region to use"`
	KeyPath    string `short:"K" long:"key-path" description:"path to store the LDAP secret key"`
	Json       bool   `long:"json" description:"print the output as JSON. useful for scripting."`
	PrettyJson bool   `long:"pretty-json" description:"print the output as pretty JSON. useful for scripting."`
	Quiet      bool   `long:"quiet" description:"no logs, only print relevant output. useful for scripting."`
}

var clientOpts = &ClientOpts{}
var cfg = &ClientConfig{"localhost", DefaultManagerRPCPort, DefaultManagerKeyPath}
var rpcClient = client.NewManagerRPCClientWithConfig(cfg)

type ManagerClient struct {
	*flags.Parser
}

func New() *ManagerClient {
	o := &ManagerClient{flags.NewParser(clientOpts, flags.Default)}

	// Manager Management
	o.AddCommand("login", "login to the system", "", &LoginCommand{})
	o.AddCommand("version", "check manager' client and server versions", "", &VersionCommand{})
	o.AddCommand("health", "check manager' health", "", &HealthCommand{})
	o.AddCommand("idle", "check if manager is idle", "", &IdleCommand{})
	o.AddCommand("register-manager", "[async] register an manager", "", &RegisterManagerCommand{})
	o.AddCommand("unregister-manager", "[async] unregister an manager", "", &UnregisterManagerCommand{})
	o.AddCommand("list-managers", "list available managers", "", &ListManagersCommand{})

	// Supervisor Management
	o.AddCommand("register-supervisor", "register an supervisor", "", &RegisterSupervisorCommand{})
	o.AddCommand("unregister-supervisor", "unregister an supervisor", "", &UnregisterSupervisorCommand{})
	o.AddCommand("list-supervisors", "list available supervisors", "", &ListSupervisorsCommand{})

	// Router Management
	o.AddCommand("register-router", "[async] register an router", "", &RegisterRouterCommand{})
	o.AddCommand("unregister-router", "[async] unregister an router", "", &UnregisterRouterCommand{})
	o.AddCommand("list-routers", "list routers", "", &ListRoutersCommand{})
	o.AddCommand("get-router", "get an router", "", &GetRouterCommand{})

	// App Management
	o.AddCommand("register-app", "register an app", "", &RegisterAppCommand{})
	o.AddCommand("unregister-app", "unregister an app", "", &UnregisterAppCommand{})
	o.AddCommand("list-registered-apps", "list registered apps", "", &ListRegisteredAppsCommand{})
	o.AddCommand("get-app", "get a registered app", "", &GetAppCommand{})

	// Environment Management
	o.AddCommand("create-dep", "create a dependency", "", &UpdateDepCommand{}) // alias to update
	o.AddCommand("update-dep", "update a dependency", "", &UpdateDepCommand{})
	o.AddCommand("get-dep", "get a dependency", "", &GetDepCommand{})
	o.AddCommand("resolve-deps", "resolve dependencies in an environment", "", &ResolveDepsCommand{})
	o.AddCommand("delete-dep", "delete a dependency", "", &DeleteDepCommand{})
	o.AddCommand("create-env", "create a environment", "", &UpdateEnvCommand{}) // alias to update
	o.AddCommand("update-env", "update a environment", "", &UpdateEnvCommand{})
	o.AddCommand("get-env", "get a environment", "", &GetEnvCommand{})
	o.AddCommand("delete-env", "delete a environment", "", &DeleteEnvCommand{})
	o.AddCommand("list-envs", "list evironments (available or deployed)", "",
		&ListEnvsCommand{})

	// Container Management
	o.AddCommand("list-containers", "list deployed containers", "", &ListContainersCommand{})
	o.AddCommand("list-shas", "list deployed shas", "", &ListShasCommand{})
	o.AddCommand("list-apps", "list deployed apps", "", &ListAppsCommand{})
	o.AddCommand("deploy", "[async] deploy something", "", &DeployCommand{})
	o.AddCommand("copy-container", "[async] deploy by copying a container", "", &CopyContainerCommand{})
	o.AddCommand("move-container", "[async] deploy by moving a container", "", &MoveContainerCommand{})
	o.AddCommand("teardown", "[async] teardown something", "", &TeardownCommand{})
	o.AddCommand("get-container", "get a container", "", &GetContainerCommand{})

	// Router Config Management
	o.AddCommand("create-pool", "create a router pool", "", &UpdatePoolCommand{}) // alias to update
	o.AddCommand("update-pool", "update a router pool", "", &UpdatePoolCommand{})
	o.AddCommand("delete-pool", "delete a router pool", "", &DeletePoolCommand{})
	o.AddCommand("get-pool", "get a router pool", "", &GetPoolCommand{})
	o.AddCommand("list-pools", "list router pools", "", &ListPoolsCommand{})
	o.AddCommand("create-rule", "create a router rule", "", &UpdateRuleCommand{}) // alias to update
	o.AddCommand("update-rule", "update a router rule", "", &UpdateRuleCommand{})
	o.AddCommand("delete-rule", "delete a router rule", "", &DeleteRuleCommand{})
	o.AddCommand("get-rule", "get a router rule", "", &GetRuleCommand{})
	o.AddCommand("list-rules", "list router rules", "", &ListRulesCommand{})
	o.AddCommand("create-trie", "create a router trie", "", &UpdateTrieCommand{}) // alias to update
	o.AddCommand("update-trie", "update a router trie", "", &UpdateTrieCommand{})
	o.AddCommand("delete-trie", "delete a router trie", "", &DeleteTrieCommand{})
	o.AddCommand("get-trie", "get a router trie", "", &GetTrieCommand{})
	o.AddCommand("list-tries", "list router tries", "", &ListTriesCommand{})

	// LDAP Management
	o.AddCommand("create-team", "create a team", "", &CreateTeamCommand{})
	o.AddCommand("delete-team", "delete a team", "", &DeleteTeamCommand{})
	o.AddCommand("list-teams", "list all teams", "", &ListTeamsCommand{})
	o.AddCommand("list-team-emails", "list all team emails", "", &ListTeamEmailsCommand{})
	o.AddCommand("list-team-admins", "list all team admins", "", &ListTeamAdminsCommand{})
	o.AddCommand("list-team-members", "list all team members", "", &ListTeamMembersCommand{})
	o.AddCommand("list-team-apps", "list all team appss", "", &ListTeamAppsCommand{})
	o.AddCommand("add-team-member", "add member to a team", "", &AddTeamMemberCommand{})
	o.AddCommand("remove-team-member", "remove member from team", "", &RemoveTeamMemberCommand{})
	o.AddCommand("allow-app", "allow an app for deploy by a team", "", &AllowTeamAppCommand{})
	o.AddCommand("disallow-app", "disallow an app for deploy by a team", "", &DisallowTeamAppCommand{})
	o.AddCommand("add-team-admin", "add team admin to a team", "", &AddTeamAdminCommand{})
	o.AddCommand("remove-team-admin", "dellete team admin from a team", "", &RemoveTeamAdminCommand{})
	o.AddCommand("add-team-email", "add an email address to a team", "", &AddTeamEmailCommand{})
	o.AddCommand("remove-team-email", "delete an email address from a team", "", &RemoveTeamEmailCommand{})

	// Container Utilities
	o.AddCommand("ssh", "ssh into a container", "", &SSHCommand{})
	o.AddCommand("tail", "tail the current day's log of a container, won't advance to next day", "", &TailCommand{})
	o.AddCommand("container-maintenance", "set or unset maintenance mode for a container", "",
		&ContainerMaintenanceCommand{})

	// Task Management
	o.AddCommand("status", "get the status of an async command", "", &StatusCommand{})
	o.AddCommand("result", "get the result of an async command", "", &ResultCommand{})
	o.AddCommand("wait", "get the wait of an async command", "", &WaitCommand{})
	o.AddCommand("deploy-result", "get the result of an async deploy", "", &DeployResultCommand{})
	o.AddCommand("teardown-result", "get the result of an async teardown", "", &TeardownResultCommand{})

	return o
}

// Runs the  Passing -d as the first flag will run the server, otherwise the client is run.
func (o *ManagerClient) Run() {
	o.Parse()
}

// Used to initialize a before a command is run. This is the first thing that should happen in all Executes.
func Init() error {
	overlayConfig()
	return AutoLoginDefault()
}

func InitNoLogin() {
	overlayConfig()
}

func overlayConfig() {
	if clientOpts.Config != "" {
		_, err := toml.DecodeFile(clientOpts.Config, cfg)
		if err != nil {
			Log(err.Error())
			// no need to panic here. we have reasonable defaults.
		}
	} else if clientOpts.Region != "" {
		_, err := toml.DecodeFile("/usr/local/etc/atlantis/manager/client."+clientOpts.Region+".toml", cfg)
		if err != nil {
			Log(err.Error())
			// no need to panic here. we have reasonable defaults.
		}
	}
	if clientOpts.Host != "" {
		cfg.Host = clientOpts.Host
	}
	if clientOpts.Port != 0 {
		cfg.Port = clientOpts.Port
	}
	if clientOpts.KeyPath != "" {
		cfg.KeyPath = clientOpts.KeyPath
	}
}
