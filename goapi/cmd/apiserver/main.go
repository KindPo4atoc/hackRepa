package main

import (
	"flag"
	apiserver "goapi/internal/app"
	"log"

	"github.com/BurntSushi/toml"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "config/apiserver.toml", "path to config file")
}

func main() {
	flag.Parse()

	config := apiserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)

	if err != nil {
		log.Fatal(err)
	}

	s := apiserver.New(config)

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}

}
