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

package dnscli

import (
	"encoding/json"
	"fmt"
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

type ClientOpts struct {
	// Only use capital letters here. Also, "H" is off limits. kthxbye.
	Json       bool `long:"json" description:"print the output as JSON. useful for scripting."`
	PrettyJson bool `long:"pretty-json" description:"print the output as pretty JSON. useful for scripting."`
	Quiet      bool `long:"quiet" description:"no logs, only print relevant output. useful for scripting."`
}

var clientOpts = &ClientOpts{}

type DNSClient struct {
	*flags.Parser
}

func New() *DNSClient {
	o := &DNSClient{flags.NewParser(clientOpts, flags.Default)}

	// DNS Management
	o.AddCommand("create-arecord", "create an a record", "", &DNSCreateARecordCommand{})
	o.AddCommand("create-alias", "create an alias record", "", &DNSCreateAliasCommand{})
	o.AddCommand("create-cname", "create a cname record", "", &DNSCreateCNameCommand{})
	o.AddCommand("delete-records", "delete records", "", &DNSDeleteRecordsCommand{})
	o.AddCommand("delete-cname", "delete cname", "", &DNSDeleteCNameCommand{})
	o.AddCommand("delete-records-value", "delete records for value", "", &DNSDeleteRecordsForValueCommand{})
	o.AddCommand("get-records-value", "get records for value", "", &DNSGetRecordsForValueCommand{})

	return o
}

// Runs the  Passing -d as the first flag will run the server, otherwise the client is run.
func (o *DNSClient) Run() {
	o.Parse()
}
