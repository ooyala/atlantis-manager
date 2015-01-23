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

package main

import (
	"atlantis/manager/builder"
	"atlantis/manager/server"
	"github.com/BurntSushi/toml"
	"log"
	"time"
)

type ManagerServerConfig struct {
	JenkinsURI        string `toml:"jenkins_uri"`
	JenkinsJob        string `toml:"jenkins_job"`
	SimpleBuilderHost string `toml:"simple_builder_host"`
}

func main() {
	managerd := server.New()
	config := &ManagerServerConfig{}
	_, err := toml.DecodeFile(managerd.Opts.ConfigFile, config)
	if err != nil {
		log.Fatalln(err)
	}
	bldr := configureBuilder(config)
	managerd.Run(bldr)
}

func configureBuilder(config *ManagerServerConfig) builder.Builder {
	var bldr builder.Builder
	if config.JenkinsURI != "" {
		log.Printf("Initializing Jenkins Builder with URI: %s, Job: %s", config.JenkinsURI, config.JenkinsJob)
		bldr = builder.NewJenkinsBuilder(config.JenkinsURI, config.JenkinsJob)
	} else if config.SimpleBuilderHost != "" {
		uri := "http://" + config.SimpleBuilderHost
		log.Printf("Initializing Simple Builder with URI: %s, Job: %s", uri, config.JenkinsJob)
		bldr = builder.NewSimpleBuilder(uri, 120*time.Second)
	} else {
		log.Printf("No builder configured")
	}
	return bldr
}
