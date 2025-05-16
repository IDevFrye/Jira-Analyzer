package app

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	jirahandlers "github.com/jiraconnector/internal/apiJiraConnector/jiraHandlers"
	jiraservice "github.com/jiraconnector/internal/apiJiraConnector/jiraService"
	"github.com/jiraconnector/internal/connector"
	datatransformer "github.com/jiraconnector/internal/dataTransformer"
	dbpusher "github.com/jiraconnector/internal/dbPusher"
	"github.com/jiraconnector/pkg/config"
	httpSwagger "github.com/swaggo/http-swagger"
)

type JiraApp struct {
	server        *http.Server
	jiraConnector *connector.JiraConnector
	db            *dbpusher.DbPusher
	log           *slog.Logger
}

func NewApp(cfg *config.Config, log *slog.Logger) (*JiraApp, error) {
	con := connector.NewJiraConnector(cfg, log)
	log.Info("created jira connection")

	dbPusher, err := dbpusher.NewDbPusher(cfg, log)
	if err != nil {
		return nil, err
	}

	datatransformer := datatransformer.NewDataTransformer(cfg.JiraCfg.Url)

	service, err := jiraservice.NewJiraService(cfg, con, datatransformer, dbPusher, log)
	if err != nil {
		ansErr := fmt.Errorf("error create service: %w", err)
		log.Error(ansErr.Error())
		return nil, ansErr
	}
	log.Info("created jira service")

	router := mux.NewRouter()
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	jiraHandler := jirahandlers.NewHandler(service, router, log)
	log.Info("created jira handlers")

	server := &http.Server{
		Addr:    cfg.ServerCfg.Port,
		Handler: jiraHandler,
	}
	log.Info("create jira server")

	return &JiraApp{
		server:        server,
		jiraConnector: con,
		db:            dbPusher,
		log:           log,
	}, nil
}

func (a *JiraApp) Run() error {
	a.log.Info("run app")
	return fmt.Errorf("run app err: %v", a.server.ListenAndServe())
}

func (a *JiraApp) Close() {
	a.log.Info("close app")
	a.db.Close()
	a.server.Close()
}

func (a *JiraApp) GetDB() *dbpusher.DbPusher {
	return a.db
}
