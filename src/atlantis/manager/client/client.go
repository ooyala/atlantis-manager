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
	atlantis "atlantis/common"
	. "atlantis/manager/constant"
	"atlantis/manager/rpc/client"
	rpcTypes "atlantis/manager/rpc/types"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/jigish/go-flags"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"
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

func quietOutput(prefix string, v reflect.Value) {
	// For pointers and interfaces, just grab the underlying value and try again
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		quietOutput(prefix, v.Elem())
		return
	}

	switch v.Kind() {
	case reflect.Struct:
		for f := 0; f < v.NumField(); f++ {
			name := v.Type().Field(f).Name
			quietOutput(prefix+name+" ", v.Field(f))
		}
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			quietOutput(prefix, v.Index(i))
		}
	case reflect.Map:
		for _, k := range v.MapKeys() {
			if name, ok := k.Interface().(string); ok {
				quietOutput(prefix+name+" ", v.MapIndex(k))
			} else {
				quietOutput(prefix+"UnknownField", v.MapIndex(k))
			}
		}
	default:
		if v.CanInterface() {
			fmt.Printf("%s%v\n", prefix, v.Interface())
		} else {
			fmt.Printf("%s <Invalid>", prefix)
		}
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
		quietOutput("", reflect.ValueOf(quiet))
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
	Host       string   `short:"M" long:"manager-host" description:"the manager host"`
	Port       uint16   `short:"P" long:"manager-port" description:"the manager port"`
	Config     string   `short:"F" long:"config-file" default:"" description:"the config file to use"`
	Regions    []string `short:"R" long:"manager-region" default:"us-east-1" description:"the regions to use"`
	KeyPath    string   `short:"K" long:"key-path" description:"path to store the LDAP secret key"`
	Json       bool     `long:"json" description:"print the output as JSON. useful for scripting."`
	PrettyJson bool     `long:"pretty-json" description:"print the output as pretty JSON. useful for scripting."`
	Quiet      bool     `long:"quiet" description:"no logs, only print relevant output. useful for scripting."`
}

var clientOpts = &ClientOpts{}
var cfg = []atlantis.RPCServerOpts{}
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
func genericLogData(name string, v reflect.Value, indent string, skip int) {
	prefix := "->" + indent
	if name != "" {
		prefix = "-> " + indent + name + ":"
	}

	// For pointers and interfaces, just grab the underlying value and try again
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		genericLogData(name, v.Elem(), indent, skip)
		return
	}

	// Only print if skip is 0.  Recursive calls should indent or decrement skip, depending on if anything was
	// printed.
	var skipped bool
	if skip == 0 {
		skipped = false
		indent = indent + "  "
	} else {
		skipped = true
		skip = skip - 1
	}

	switch v.Kind() {
	case reflect.Struct:
		if !skipped {
			Log("%s\n", prefix)
		}
		for f := 0; f < v.NumField(); f++ {
			name := v.Type().Field(f).Name
			genericLogData(name, v.Field(f), indent, skip)
		}
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		if !skipped {
			Log("%s\n", prefix)
		}
		for i := 0; i < v.Len(); i++ {
			genericLogData("", v.Index(i), indent, skip)
		}
	case reflect.Map:
		if !skipped {
			Log("%s\n", prefix)
		}
		for _, k := range v.MapKeys() {
			if name, ok := k.Interface().(string); ok {
				genericLogData(name, v.MapIndex(k), indent, skip)
			} else {
				genericLogData("Unknown field", v.MapIndex(k), indent, skip)
			}
		}
	default:
		if v.CanInterface() {
			Log("%s %v", prefix, v.Interface())
		} else {
			Log("%s <Invalid>", prefix)
		}
	}
}

func executeFlags(rv reflect.Value) (message, rpc, field, name, fileName, fileField string, fileData interface{},
	noauth, async, wait bool) {
	// Use flags from the Properties field, if it exists.
	if properties, ok := rv.Type().FieldByName("Properties"); ok {
		message = properties.Tag.Get("message")
		rpc = properties.Tag.Get("rpc")
		field = properties.Tag.Get("field")
		name = properties.Tag.Get("name")
		fileField = properties.Tag.Get("filefield")
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
		wait = waitField.Bool()
	}

	// If we need to read file data, get that set up
	if fileDataField := rv.FieldByName("FileData"); fileDataField.IsValid() {
		if fileNamev := rv.FieldByName("FromFile"); fileNamev.IsValid() {
			fileDatav := reflect.New(fileDataField.Type())
			fileData = fileDatav.Interface()
			fileName = fileNamev.String()
			if fileField == "" {
				// The default name is the unqualified type: *types.DependerEnvData => DependerEnvData
				components := strings.Split(fileDatav.Type().String(), ".")
				fileField = components[len(components)-1]
			}
		}
	}
	return
}

func copyType(rv reflect.Value, name string) (reflect.Value, error) {
	field := rv.FieldByName(name)
	if !field.IsValid() {
		return rv, errors.New("Internal error: " + rv.Type().String() + " missing field " + name)
	}
	return reflect.New(field.Type()), nil
}

func genericResult(command interface{}, args []string) (map[string]string, string, map[string]map[string]interface{}, error) {
	rv := reflect.ValueOf(command).Elem()
	// Extract all the configuration flags from the Command struct
	message, rpc, field, name, fileName, fileField, fileData, noauth, async, wait := executeFlags(rv)

	// Some command require auth, some don't
	if noauth {
		InitNoLogin()
	} else {
		if err := Init(); err != nil {
			return nil, "", nil, OutputError(err)
		}
	}

	// Set up the arg object based on types in the Command struct.
	argv, err := copyType(rv, "Arg")
	if err != nil {
		return nil, "", nil, OutputError(err)
	}

	// Async replies should use the ID field, not the final response field
	if async && !wait {
		field = "ID"
		name = "ID"
	}

	// Copy args from the CLI arguments to the RPC arguments, and handle any positional ones.
	genericExtractArgs(rv, args)
	genericCopyArgs(rv, argv.Elem())

	// Read in file data if necessary
	if fileData != nil {
		file, err := os.Open(fileName)
		if err != nil {
			return nil, "", nil, OutputError(err)
		}
		jsonDec := json.NewDecoder(file)
		if err := jsonDec.Decode(fileData); err != nil {
			return nil, "", nil, OutputError(err)
		}
		argv.Elem().FieldByName(fileField).Set(reflect.ValueOf(fileData))
	}

	Log(message + "...")

	// Set up some storage for the results...
	statuses := map[string]string{}
	replies := map[string]interface{}{}
	datas := map[string]map[string]interface{}{}

	// Now we're prepped; let's make the requests and store the results to return
	for region := range cfg {
		// Set up the reply object based on types in the Command struct.  But async replies always come back as
		// type AsyncReply.
		var reply interface{}
		var replyv reflect.Value
		if async {
			reply = &atlantis.AsyncReply{}
			replyv = reflect.ValueOf(reply)
		} else {
			replyv, err = copyType(rv, "Reply")
			if err != nil {
				return nil, "", nil, OutputError(err)
			}
			reply = replyv.Interface()
		}

		// Actually make the request, either with or without auth.
		if noauth {
			arg := argv.Interface()
			if err := rpcClient.CallMulti(rpc, arg, region, reply); err != nil {
				return nil, "", nil, OutputError(err)
			}
		} else {
			arg := argv.Interface().(client.AuthedArg)
			if err := rpcClient.CallAuthedMulti(rpc, arg, region, reply); err != nil {
				return nil, "", nil, OutputError(err)
			}
		}

		// If we're waiting on an async command, update the reply when it comes back.
		if wait {
			if idv := replyv.Elem().FieldByName("ID"); idv.IsValid() {
				Log("-> ID: %v", replyv.Elem().FieldByName("ID"))
				replyv = reflect.New(rv.FieldByName("Reply").Type())
				reply = replyv.Interface()
				if err := genericWait(command, rpc, idv.String(), reply); err != nil {
					return nil, "", nil, OutputError(err)
				}
			} else {
				Log("Error: No async ID found in response")
			}
		}

		// Extract that status and desired field from the RPC result.
		status := "Unknown"
		if v := replyv.Elem().FieldByName("Status"); v.IsValid() {
			status = v.Interface().(string)
		}
		data := map[string]interface{}{}
		if field != "" {
			if v := replyv.Elem().FieldByName(field); v.IsValid() {
				data[name] = v.Interface()
			}
		}

		// Async results don't give back status immediately.
		if async && !wait {
			status = ""
		}

		// Get the name of the region to return
		regionName := "Unknown region"
		if len(clientOpts.Regions) > region {
			regionName = clientOpts.Regions[region]
		}

		// And now we're done; save everything to return.
		statuses[regionName] = status
		replies[regionName] = reply
		datas[regionName] = data
	}

	return statuses, name, datas, nil
}

func genericWait(command interface{}, rpc, id string, reply interface{}) error {
	Log("Waiting...")
	var statusReply atlantis.TaskStatus
	var currentStatus string
	if err := rpcClient.Call("Status", id, &statusReply); err != nil {
		return OutputError(err)
	}
	for !statusReply.Done {
		if currentStatus != statusReply.Status {
			currentStatus = statusReply.Status
			Log(currentStatus)
		}
		time.Sleep(waitPollInterval)
		if err := rpcClient.Call("Status", id, &statusReply); err != nil {
			return OutputError(err)
		}
	}

	// And it's done!  Fetch the result
	if err := rpcClient.Call(rpc+"Result", id, reply); err != nil {
		return OutputError(err)
	}

	return nil
}

func genericExecuter(command interface{}, args []string) error {
	status, name, data, err := genericResult(command, args)
	if err != nil {
		return err
	}

	if len(status) == 1 {
		for _, s := range status {
			if s != "" {
				Log("-> status: %s", s)
			}
		}
	} else {
		genericLogData("status", reflect.ValueOf(status), "", 0)
	}
	if data != nil {
		// Skip printing the region if there's only one.
		skip := 0
		if len(data) == 1 {
			skip = 1
		}
		for r, d := range data {
			genericLogData(r, reflect.ValueOf(d), "", skip)
		}
	}

	// JSON output should be grouped by region
	jsonOutput := map[string]interface{}{}
	for region, s := range status {
		jsonOutput[region] = map[string]interface{}{"status": s, name: data[region][name]}
	}

	// Quiet data shouldn't include known field name
	quietOutputArray := map[string]interface{}{}
	for region, _ := range status {
		quietOutputArray[region] = data[region][name]
	}
	var quietOutput interface{}
	quietOutput = quietOutputArray

	// If there's only one region, omit it.
	if len(status) == 1 {
		for region, _ := range status {
			if regionData, ok := jsonOutput[region].(map[string]interface{}); ok {
				jsonOutput = regionData
			}
			quietOutput = quietOutputArray[region]
		}
	}

	return Output(jsonOutput, quietOutput, nil)
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

// TODO(edanaher): This config parsing is kind of horrendous, but hopefully does the right thing.
func overlayConfig() {
	if clientOpts.Config != "" {
		_, err := toml.DecodeFile(clientOpts.Config, cfg)
		if err != nil {
			fmt.Print("Error parsing config file " + clientOpts.Config + ":\n" + err.Error() + "\n")
			os.Exit(1)
		}
	} else {
		// NOTE(edanaher): The default doesn't get removed if more are passed in
		if len(clientOpts.Regions) > 1 {
			clientOpts.Regions = clientOpts.Regions[1:]
		}
		for _, region := range clientOpts.Regions {
			configFileFound := false
			for _, path := range configDirs {
				filename := path + "client." + region + ".toml"
				if ok, _ := exists(filename); ok {
					var curCfg ClientConfig
					_, err := toml.DecodeFile(filename, &curCfg)
					if err != nil {
						fmt.Print("Error parsing config file " + filename + ":\n" + err.Error() + "\n")
						os.Exit(1)
					}
					// Defaults need to be loaded for each file independently
					if curCfg.Port == 0 {
						curCfg.Port = DefaultManagerRPCPort
					}
					if curCfg.KeyPath == "" {
						curCfg.KeyPath = DefaultManagerKeyPath
					}
					cfg = append(cfg, &curCfg)
					configFileFound = true
					break
				}
			}
			if !configFileFound {
				fmt.Print("Error: could not find config file for " + region + "\n")
				os.Exit(1)
			}
		}
	}
	// If other options are passed in, assume there's only one region and we should use that one
	/* NOTE(edanaher): cfg has to be an array of interfaces, because arrays don't get auto-inferfaced properly.
	 * But then we have to cast it here.  *sigh* */
	if clientOpts.Host != "" {
		cfg[0].(*ClientConfig).Host = clientOpts.Host
	}
	if clientOpts.Port != 0 {
		cfg[0].(*ClientConfig).Port = clientOpts.Port
	}
	if clientOpts.KeyPath != "" {
		cfg[0].(*ClientConfig).KeyPath = clientOpts.KeyPath
	}
	// TODO(edanaher): This is aliased.  The appends above may have unaliased it.  Why do we do this?
	rpcClient.Opts = cfg
}
