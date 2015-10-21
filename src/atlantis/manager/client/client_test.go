package client

import (
	atlantis "atlantis/common"
	. "atlantis/manager/rpc/types"
	"atlantis/manager/util"
	"fmt"
	"github.com/jigish/go-flags"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

var CONFIG map[string]string

func init() {
	dir, _ := os.Getwd()
	CONFIG = util.GetTestConfig(filepath.Base(dir))
}

func setHostConfig() {
	cfg = []atlantis.RPCServerOpts{}
	rpcClient.Opts = cfg
	clientOpts.Regions = []string{CONFIG["region"]}
	for _, r := range cfg {
		log.Printf("Setting region: %v\n", r)
	}
}

func assertNoErr(t *testing.T, err error, msg string) bool {
	if err != nil {
		t.Errorf("%v:\n  %v\n", msg, err)
	}
	return err == nil
}

func assertEqual(t *testing.T, expected, actual, msg string) {
	if expected != actual {
		t.Errorf("%v\n  Expected: %v\nActual: %v\n", msg, expected, actual)
	}
}

func compare(t *testing.T, field string, expected, actual reflect.Value) {
	if !actual.IsValid() {
		t.Errorf("Invalid value on field %v: %v", field, actual)
		return
	}
	if !expected.IsValid() {
		t.Errorf("Invalid expected value on field %v: %v", field, actual)
		return
	}
	if expected.Kind() == reflect.Ptr && expected.IsNil() {
		return
	}
	if expected.Type() != actual.Type() {
		t.Errorf("Wrong type on %v\n  Expected: %v\n  Actual:   %v\n", field, expected.Type(), actual.Type())
		return
	}

	switch expected.Kind() {
	case reflect.Ptr:
		compare(t, field, expected.Elem(), actual.Elem())
	case reflect.Struct:
		for f := 0; f < expected.NumField(); f++ {
			name := expected.Type().Field(f).Name
			actualField := actual.FieldByName(name)
			if !actualField.IsValid() {
				t.Errorf("Expected field %v on %v not found", name, field)
				continue
			}
			compare(t, field+"."+name, expected.Field(f), actualField)
		}
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		for i := 0; i < expected.Len(); i++ {
			found := false
			for j := 0; j < actual.Len(); j++ {
				if expected.Index(i).Interface() == actual.Index(j).Interface() {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("In field %v, expected value %v not found", field, expected.Index(i))
			}
		}
	default:
		if expected.Interface() != actual.Interface() {
			t.Errorf("Error in field %v:\n  Expected: %v\n  Actual: %v\n", field, expected.Interface(), actual.Interface())
		}
	}
}

func checkResult(t *testing.T, expected, actual interface{}) {
	compare(t, reflect.ValueOf(actual).Type().String(), reflect.ValueOf(expected), reflect.ValueOf(actual))
}

func checkCommand(t *testing.T, command interface{}, line string, expected interface{}) interface{} {
	setHostConfig()
	args := strings.Fields(line)
	// If it has a built-in execute, just run it for now and don't worry about the result.
	if command, ok := command.(flags.Command); ok {
		err := command.Execute(args)
		assertNoErr(t, err, "Error executing opaque command")
		return nil
	}
	status, reply, name, data, err := genericResult(command, args)
	//t.Logf("err is %v\n", err)
	_ = status
	_ = name
	_ = data
	_ = reply
	if assertNoErr(t, err, "Error executing command") && expected != nil {
		checkResult(t, expected, reply[CONFIG["region"]])
	}
	return reply[CONFIG["region"]]
}

func TestRegisterApp(t *testing.T) {
	testName := "e2e-test-" + time.Now().Format("2006-01-02T15-04-05")

	log.Print("== Registering dummy app ==")
	registerCommand := &RegisterAppCommand{Name: testName, Repo: "fake-git-repo", Root: "/path/to/app",
		Email: "owner@example.com"}
	checkCommand(t, registerCommand, "", &ManagerRegisterAppReply{"OK"})

	log.Print("== Listing apps and checking existence ==")
	checkCommand(t, &ListRegisteredAppsCommand{}, "", &ManagerListRegisteredAppsReply{[]string{testName}, "OK"})

	log.Print("== Unregistering dummy app ==")
	checkCommand(t, &UnregisterAppCommand{}, testName, &ManagerRegisterAppReply{"OK"})
}

func writeDepFile(t *testing.T, env string, data string) string {
	filename := "/tmp/atlantis-test-" + env
	fileData := fmt.Sprintf(`{ "Name": "%s", "DataMap": %s }`, env, data)
	err := ioutil.WriteFile(filename, []byte(fileData), 0644)
	if err != nil {
		t.Errorf("Error writing dependency file: %v", err)
	}
	return filename
}

func checkUrl(t *testing.T, url, name, expected string) bool {
	resp, err := http.Get(url)
	if err != nil {
		t.Errorf("Error fetching from %s: %v\n", name, err)
		return false
	} else if resp.StatusCode != 200 {
		t.Errorf("Non-200 status code from %s: %v\n", name, resp.StatusCode)
		return false
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			t.Errorf("Error reading response body from %s: %v\n", name, err)
			return false
		} else {
			if !strings.Contains(string(body), expected) {
				t.Errorf("Response from %s is missing '%s': %v\n", name, expected, string(body))
				return false
			}
		}
	}
	return true
}

func TestFullDeploy(t *testing.T) {
	testName := "e2e-test-" + time.Now().Format("2006-01-02T15-04-05")

	checkCommand(t, &ListAppsCommand{}, "", &ManagerListAppsReply{[]string{CONFIG["app"]}, "OK"})

	log.Print("== Creating environment ==")
	checkCommand(t, &UpdateEnvCommand{}, testName, &ManagerEnvReply{"OK"})
	checkCommand(t, &ListEnvsCommand{}, "", &ManagerListEnvsReply{[]string{testName}, "OK"})

	log.Print("== Setting cmk dependency ==")
	envFile := writeDepFile(t, testName, `{ "contact_group": CONFIG["contact_group"] }`)
	cmkDepCommand := &AddDependerEnvDataForDependerAppCommand{App: "cmk", Depender: CONFIG["app"], FromFile: envFile}
	checkCommand(t, cmkDepCommand, "", &ManagerAddDependerEnvDataForDependerAppReply{"OK", nil})

	log.Print("== Deploying hello-go ==")
	deployCommand := &DeployCommand{App: CONFIG["app"], Sha: CONFIG["sha"],
		Env: testName, Dev: true, Wait: true}
	// TODO(edanaher): Status should be OK...
	deployi := checkCommand(t, deployCommand, "", &ManagerDeployReply{"", nil})
	if deploy, ok := deployi.(*ManagerDeployReply); ok {
		container := deploy.Containers[0]

		log.Print("== Checking container is accessible ==")
		url := fmt.Sprintf("http://%s:%d/", container.Host, container.PrimaryPort)
		checkUrl(t, url, "deployed container", "Hello from Go")

		log.Print("== Looking up port information ==")
		getPortCommand := &GetAppEnvPortCommand{App: CONFIG["app"], Env: testName, Internal: true}
		portReplyi := checkCommand(t, getPortCommand, "", nil)
		if portReply, ok := portReplyi.(*ManagerGetAppEnvPortReply); ok {
			log.Print("== Checking router port ==")
			routerPort := portReply.Port.Port
			url := fmt.Sprintf("http://%s:%d/", CONFIG["host"], routerPort)
			checkUrl(t, url, "router port", "Hello from Go")
		}

		log.Print("== Tearing down hello-go ==")
		teardownCommand := &TeardownCommand{}
		teardownCommand.ContainerID = container.ID
		teardownCommand.Wait = true
		// TODO(edanaher): Status should be OK...
		checkCommand(t, teardownCommand, "", &ManagerTeardownReply{nil, ""})
	} else {
		t.Errorf("Bad response from manager for deploy")
	}

	// Clean up the environment
	log.Print("== Deleting environment ==")
	log.Printf("%s\n", testName)
	checkCommand(t, &DeleteEnvCommand{}, testName, &ManagerEnvReply{"OK"})
}
