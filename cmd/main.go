package main

import (
	"flag"
	"log"

	"github.com/tabularasa31/antibruteforce/config"
	"github.com/tabularasa31/antibruteforce/internal/app"

	_ "github.com/lib/pq"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./config/config", "Path to configuration file")
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
