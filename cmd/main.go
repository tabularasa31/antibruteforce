package main

import (
	"flag"
	"github.com/tabularasa31/antibruteforce/config"
	"github.com/tabularasa31/antibruteforce/internal/app"
	"log"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./config/config", "Path to configuration file")
}

func main() {
	flag.Parse()

	// Configuration
	//cfg, err := config.NewConfig(configFile)
	//if err != nil {
	//	log.Fatalf("Config error: %s", err)
	//}

	//configPath := config.GetConfigPath(os.Getenv("config"))
	cfg, err := config.GetConfig(configFile)
	if err != nil {
		log.Fatalf("Loading config: %v", err)
	}

	// Run
	app.Run(cfg)
}
