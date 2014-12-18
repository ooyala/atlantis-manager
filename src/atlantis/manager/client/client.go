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

package client

import (
	. "atlantis/manager/constant"
	"atlantis/manager/rpc/client"
	rpcTypes "atlantis/manager/rpc/types"
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/jigish/go-flags"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"
)

// Directories to be searched for client configuration files.
var configDirs = []string{
	"/usr/local/etc/atlantis/manager/", // brew config path
	"/etc/atlantis/manager/",           // deb package config path
	"/opt/atlantis/manager/etc",        // just in case config path
}

func Log(format string, args ...interface{}) {
	if !IsJson() && !clientOpts.Quiet {
		//the standard logger which log.Printf uses
		//defaults to stderr, set out to stdout before
		//printing
		log.SetOutput(os.Stdout)
		log.Printf(format, args...)
	}
}

func Fatal(format string, args ...interface{}) {
	if !IsJson() && !clientOpts.Quiet {
		//the Log function above changes where the
		//standard logger prints to, thus reset it
		//to stderr before calling fatal
		log.SetOutput(os.Stderr)
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
	case []uint16:
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
var rpcClient = &client.ManagerRPCClient{*client.NewManagerRPCClientWithConfig(cfg), "", map[string]string{}}
var dummyAuthArg = rpcTypes.ManagerAuthArg{"", "", ""}

type commandWrapper struct {
	Command interface{}
}

func (c *commandWrapper) Execute(args []string) error {
	// If the command has an Execute method, honor it.  Otherwise, fall back to the generic Execute.
	if command, ok := c.Command.(flags.Command); ok {
		return command.Execute(args)
	} else {
		return genericExecuter(c.Command, args)
	}
}

type ManagerClient struct {
	*flags.Parser
}

func (m *ManagerClient) AddCommand(name, short, long string, data interface{}) {
	m.Parser.AddCommand(name, short, long, &commandWrapper{data})
}

func New() *ManagerClient {
	o := &ManagerClient{flags.NewParser(clientOpts, flags.Default)}

	// Manager Management
	o.AddCommand("login", "login to the system", "", &LoginCommand{})
	o.AddCommand("version", "check manager client and server versions", "", &VersionCommand{})
	o.AddCommand("health", "check manager health", "", &HealthCommand{})
	o.AddCommand("usage", "check manager usage stats", "", &UsageCommand{})
	o.AddCommand("idle", "check if manager is idle", "", &IdleCommand{})
	o.AddCommand("register-manager", "[async] register an manager", "", &RegisterManagerCommand{})
	o.AddCommand("unregister-manager", "[async] unregister an manager", "", &UnregisterManagerCommand{})
	o.AddCommand("list-managers", "list available managers", "", &ListManagersCommand{})
	o.AddCommand("get-manager", "get manager", "", &GetManagerCommand{})
	o.AddCommand("get-self", "get this manager", "", &GetSelfCommand{})
	o.AddCommand("add-role", "add role to manager", "", &AddRoleCommand{})
	o.AddCommand("remove-role", "remove role from manager", "", &RemoveRoleCommand{})
	o.AddCommand("has-role", "check role on manager", "", &HasRoleCommand{})

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
	o.AddCommand("list-authorized-registered-apps", "list authorized registered apps", "", &ListAuthorizedRegisteredAppsCommand{})
	o.AddCommand("get-app", "get a registered app", "", &GetAppCommand{})
	o.AddCommand("request-app-dependency", "request a dependency for an app", "", &RequestAppDependencyCommand{})
	o.AddCommand("add-depender-app-data", "add depender app data", "", &AddDependerAppDataCommand{})
	o.AddCommand("remove-depender-app-data", "remove depender app data", "", &RemoveDependerAppDataCommand{})
	o.AddCommand("get-depender-app-data", "get depender app data", "", &GetDependerAppDataCommand{})
	o.AddCommand("add-depender-env-data", "add depender env data", "", &AddDependerEnvDataCommand{})
	o.AddCommand("remove-depender-env-data", "remove depender env data", "", &RemoveDependerEnvDataCommand{})
	o.AddCommand("get-depender-env-data", "get depender env data", "", &GetDependerEnvDataCommand{})
	o.AddCommand("add-depender-env-data-for-depender-app", "add depender env data for depender app", "", &AddDependerEnvDataForDependerAppCommand{})
	o.AddCommand("remove-depender-env-data-for-depender-app", "remove depender env data for depender app", "", &RemoveDependerEnvDataForDependerAppCommand{})
	o.AddCommand("get-depender-env-data-for-depender-app", "get depender env data for depender app", "", &GetDependerEnvDataForDependerAppCommand{})

	// Environment Management
	o.AddCommand("create-dep", "create a dependency", "", &UpdateDepCommand{}) // alias to update
	o.AddCommand("update-dep", "update a dependency", "", &UpdateDepCommand{})
	o.AddCommand("get-dep", "get a dependency", "", &GetDepCommand{})
	o.AddCommand("resolve-deps", "resolve dependencies in an environment", "", &ResolveDepsCommand{})
	o.AddCommand("delete-dep", "delete a dependency", "", &DeleteDepCommand{})
	o.AddCommand("create-env", "create a environment", "", &UpdateEnvCommand{}) // alias to update
	o.AddCommand("update-env", "update a environment", "", &UpdateEnvCommand{})
	o.AddCommand("delete-env", "delete a environment", "", &DeleteEnvCommand{})
	o.AddCommand("list-envs", "list evironments (available or deployed)", "",
		&ListEnvsCommand{})

	// Container Management
	o.AddCommand("list-containers", "list deployed containers", "", &ListContainersCommand{})
	o.AddCommand("list-shas", "list deployed shas", "", &ListShasCommand{})
	o.AddCommand("list-apps", "list deployed apps", "", &ListAppsCommand{})
	o.AddCommand("deploy", "[async] deploy something", "", &DeployCommand{})
	o.AddCommand("deploy-container", "[async] deploy by replicating a container", "", &DeployContainerCommand{})
	o.AddCommand("copy-container", "[async] deploy by copying a single container to a specific host", "", &CopyContainerCommand{})
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
	o.AddCommand("create-port", "create a router port", "", &UpdatePortCommand{}) // alias to update
	o.AddCommand("update-port", "update a router port", "", &UpdatePortCommand{})
	o.AddCommand("delete-port", "delete a router port", "", &DeletePortCommand{})
	o.AddCommand("get-port", "get a router port", "", &GetPortCommand{})
	o.AddCommand("list-ports", "list router ports", "", &ListPortsCommand{})
	o.AddCommand("get-app-env-port", "get a router port for an appenv", "", &GetAppEnvPortCommand{})
	o.AddCommand("list-app-envs-with-port", "list app envs with a router port", "", &ListAppEnvsWithPortCommand{})

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
	o.AddCommand("is-app-allowed", "check if an app is allowed for deploy by a user", "", &IsAppAllowedCommand{})
	o.AddCommand("list-allowed-apps", "list all allowed apps for a user", "", &ListAllowedAppsCommand{})
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
	o.AddCommand("list-task-ids", "list all the task ids of async commands", "", &ListTaskIDsCommand{})
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

// Go through any positional arguments and add them to empty flag arguments
func genericExtractArgs(command reflect.Value, args []string) []string {
	for f := 0; f < command.NumField(); f++ {
		if len(args) == 0 {
			return args
		}
		if tp := command.Type().Field(f).Name; tp == "Properties" || tp == "Arg" || tp == "Reply" {
			continue
		}
		if str, ok := command.Field(f).Interface().(string); ok {
			if str == "" {
				command.Field(f).Set(reflect.ValueOf(args[0]))
				args = args[1:]
			}
		}
	}
	return args
}

// Copy arguments from the CLI Command struct to the RPCarg struct
func genericCopyArgs(command reflect.Value, arg reflect.Value) {
	for f := 0; f < arg.NumField(); f++ {
		name := arg.Type().Field(f).Name
		if v := command.FieldByName(name); v.IsValid() {
			arg.Field(f).Set(v)
		}
	}
}

// Pretty-print any random data we get back from the server
func genericLogData(name string, v reflect.Value, indent string) {
	prefix := "->" + indent
	if name != "" {
		prefix = "-> " + indent + name + ":"
	}
	switch v.Kind() {
	case reflect.Ptr:
		genericLogData(name, v.Elem(), indent)
	case reflect.Struct:
		Log("%s\n", prefix)
		for f := 0; f < v.NumField(); f++ {
			name := v.Type().Field(f).Name
			genericLogData(name, v.Field(f), indent+"  ")
		}
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		Log("%s\n", prefix)
		for i := 0; i < v.Len(); i++ {
			genericLogData("", v.Index(i), indent+"  ")
		}
	case reflect.Map:
		Log("%s\n", prefix)
		for _, k := range v.MapKeys() {
			if name, ok := k.Interface().(string); ok {
				genericLogData(name, v.MapIndex(k), indent+"  ")
			} else {
				genericLogData("Unknown field", v.MapIndex(k), indent+"  ")
			}
		}
	default:
    Log("%s %v", prefix, v.Interface())
	}
}

func executeFlags(rv reflect.Value) (message, rpc, field, name string, noauth, async, wait bool) {
	// Use flags from the Properties field, if it exists.
	if properties, ok := rv.Type().FieldByName("Properties"); ok {
		message = properties.Tag.Get("message")
		rpc = properties.Tag.Get("rpc")
		field = properties.Tag.Get("field")
		name = properties.Tag.Get("name")
		noauth = properties.Tag.Get("noauth") != ""
	}

	// The base name of the command is extracted from the Command type; e.g., client.ListAppsCommand is ListApps
	baseNameRegexp := regexp.MustCompile(".*\\.([A-Za-z]*)Command$")
	baseName := baseNameRegexp.ReplaceAll([]byte(rv.Type().String()), []byte("$1"))

	// The default RPC name is the base name.
	if rpc == "" {
		rpc = string(baseName)
	}

	// The default message is the base name with spaces; e.g. "List Apps"
	if message == "" {
		messageRegexp := regexp.MustCompile("(.)([A-Z])")
		message = string(messageRegexp.ReplaceAll(baseName, []byte("$1 $2")))
	}

	// The default field to print is the last component of the base name; e.g., Apps
	if field == "" {
		fieldRegexp := regexp.MustCompile("[A-Z][a-z]*$")
		field = string(fieldRegexp.Find(baseName))
	}

	// The default name for the field (for output) is the lowercase field name
	if name == "" {
		name = strings.ToLower(field)
	}

	// Async commands are weird; we get back an ID to wait on, and there's always an boolean --wait flag.
	if waitField := rv.FieldByName("Wait"); waitField.Kind() == reflect.Bool {
		async = true
		field = "ID"
		name = "ID"
		wait = waitField.Bool()
	}
	return
}

func genericResult(command interface{}, args []string) (string, interface{}, string, interface{}, *WaitCommand, error) {
	rv := reflect.ValueOf(command).Elem()
	// Extract all the configuration flags from the Command struct
	message, rpc, field, name, noauth, async, wait := executeFlags(rv)

	// Some command require auth, some don't
	if noauth {
		InitNoLogin()
	} else {
		if err := Init(); err != nil {
			return "", nil, "", nil, nil, OutputError(err)
		}
	}

	// Set up the arg and reply objects based on types in the Command struct.
	argv := reflect.New(rv.FieldByName("Arg").Type())
	replyv := reflect.New(rv.FieldByName("Reply").Type())
	reply := replyv.Interface()

	// Copy args from the CLI arguments to the RPC arguments, and handle any positional ones.
	genericExtractArgs(rv, args)
	genericCopyArgs(rv, argv.Elem())

	Log(message + "...")

	// Actually make the request, either with or without auth.
	if noauth {
		arg := argv.Interface()
		if err := rpcClient.Call(rpc, arg, reply); err != nil {
			return "", nil, "", nil, nil, OutputError(err)
		}
	} else {
		arg := argv.Interface().(client.AuthedArg)
		if err := rpcClient.CallAuthed(rpc, arg, reply); err != nil {
			return "", nil, "", nil, nil, OutputError(err)
		}
	}

	// Extract that status and desired field from the RPC result.
	status := "Unknown"
	if v := replyv.Elem().FieldByName("Status"); v.IsValid() {
		status = v.Interface().(string)
	}
	var data interface{}
	data = "Unknown"
	if field != "" {
		if v := replyv.Elem().FieldByName(field); v.IsValid() {
			data = v.Interface()
		}
	}

	// Async results don't give back status immediately.
	if async {
		status = ""
	}

	// If we're waiting on an async command, return the magic WaitCommand on the asyncID to poll until finished.
	var waitCommand *WaitCommand
	if wait {
		if idv := replyv.Elem().FieldByName("ID"); idv.IsValid() {
			waitCommand = &WaitCommand{idv.String()}
		} else {
			Log("Error: No async ID found in response")
		}
	}
	return status, reply, name, data, waitCommand, nil
}

func genericExecuter(command interface{}, args []string) error {
	status, reply, name, data, waitCommand, err := genericResult(command, args)
	if err != nil {
		return err
	}

	if status != "" {
		Log("-> status: %s", status)
	}
	if data != nil {
		genericLogData(name, reflect.ValueOf(data), "")
	}
	if waitCommand != nil {
		return waitCommand.Execute(args)
	}

	return Output(map[string]interface{}{"status": status, name: data}, reply, nil)
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func overlayConfig() {
	configFileFound := false
	if clientOpts.Config != "" {
		_, err := toml.DecodeFile(clientOpts.Config, cfg)
		if err != nil {
			fmt.Print("Error parsing config file " + clientOpts.Config + ":\n" + err.Error() + "\n")
			os.Exit(1)
		}
	} else if clientOpts.Region != "" {
		for _, path := range configDirs {
			filename := path + "client." + clientOpts.Region + ".toml"
			if ok, _ := exists(filename); ok {
				_, err := toml.DecodeFile(filename, cfg)
				if err != nil {
					fmt.Print("Error parsing config file " + filename + ":\n" + err.Error() + "\n")
					os.Exit(1)
				}
				configFileFound = true
				break
			}
		}
		if !configFileFound {
			fmt.Print("Error: could not find config file for " + clientOpts.Region + "\n")
			os.Exit(1)
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
