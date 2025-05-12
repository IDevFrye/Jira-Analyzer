package main

import (
	"log/slog"

	"github.com/jiraconnector/cmd/app"
	"github.com/jiraconnector/pkg/config"
	"github.com/jiraconnector/pkg/logger"
)

func main() {
	//read config
	cfg := config.LoadConfig()

	//setting logger
	log := logger.SetupLogger(cfg.Env, cfg.LogFile)
	log.Info("starting url-shortener", slog.String("env", cfg.Env))

	//create connector app
	a, err := app.NewApp(cfg, log)
	if err != nil {
		log.Error("error create app")
		panic(err)
	}
	log.Info("created app")

	//start app
	if err := a.Run(); err != nil {
		log.Error("error run app")
		panic(err)
	}
	defer a.Close()
}
