package main

import (
	"flag"
	"log"

	_ "github.com/lib/pq"
	"github.com/tabularasa31/antibruteforce/config"
	"github.com/tabularasa31/antibruteforce/internal/app"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./config/config.yml", "Path to configuration file")
}

func main() {
	flag.Parse()

	// Configuration
	cfg, err := config.GetConfig(configFile)
	if err != nil {
		log.Fatalf("Loading config: %v", err)
	}

	// Run
	app.Run(cfg)
}
