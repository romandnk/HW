package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/romandnk/HW/hw12_13_14_15_calendar/cmd/config"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/sender_config.toml", "Path to configuration file")
}

func main() {
	cfg, err := config.NewSenderConfig(configFile)
	if err != nil {
		log.Fatalf("rabbit sender config error: %s", err.Error())
	}

	fmt.Printf("%+v\n", cfg)
}
