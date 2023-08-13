package main

import (
	"flag"
	"fmt"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/cmd/config"
	"log"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/sender_config.toml", "Path to configuration file")
}

func main() {
	cfg, err := config.NewSenderRabbitConfig(configFile)
	if err != nil {
		log.Fatalf("rabbit sender config error: %s", err.Error())
	}

	//logg := logger.NewLogger(cfg.Logger.Level, cfg.Logger.Representation)

	fmt.Printf("%+v\n", cfg)
}
