package main

import (
	"atlantis/manager/builder"
	"atlantis/manager/server"
	"github.com/BurntSushi/toml"
	"log"
)

type OoyalaServerConfig struct {
	JenkinsURI string `toml:"jenkins_uri"`
}

func main() {
	managerd := server.New()
	config := &OoyalaServerConfig{}
	_, err := toml.DecodeFile(managerd.Opts.ConfigFile, config)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Initializing Jenkins Builder with URI: %s", config.JenkinsURI)
	bldr := builder.NewJenkinsBuilder(config.JenkinsURI)
	managerd.Run(bldr)
}
