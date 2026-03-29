package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/endpointhandler/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	configPath := flag.String("config", "configs/local.yaml", "path to endpointHandler YAML config")
	sqlPath := flag.String("sql", "", "path to init.sql (required)")
	flag.Parse()
	if *sqlPath == "" {
		log.Fatal("missing --sql path to init.sql")
	}

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.User, cfg.Database.Password,
		cfg.Database.DBName, cfg.Database.Port, cfg.Database.SSLMode,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("connect postgres: %v", err)
	}
	defer db.Close()

	sqlBytes, err := os.ReadFile(*sqlPath)
	if err != nil {
		log.Fatalf("read sql file: %v", err)
	}

	if _, err := db.Exec(string(sqlBytes)); err != nil {
		log.Fatalf("apply sql: %v", err)
	}
	log.Println("init.sql applied OK")
}
