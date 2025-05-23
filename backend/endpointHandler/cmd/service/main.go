package main

import (
	"flag"
	"log"

	"github.com/endpointhandler/config"
	"github.com/endpointhandler/repository"
	"github.com/endpointhandler/router"
)

func main() {

	configPath := flag.String("config", "config.yaml", "Path to config file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config from %s: %v", *configPath, err)
	}

	err = repository.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to init DB: %v", err)
	}

	r := router.SetupRouter(cfg)
	log.Fatal(r.Run(":" + cfg.Server.Port))
}
