package jiraservice

import (
	"fmt"
	"log/slog"

	datatransformer "github.com/jiraconnector/internal/dataTransformer"
	"github.com/jiraconnector/internal/structures"
	"github.com/jiraconnector/pkg/config"
	"github.com/jiraconnector/pkg/logger"
)

//go:generate mockery

type JiraConnectorInterface interface {
	GetAllProjects() ([]structures.JiraProject, error)
	GetProjectsPage(search string, limit, page int) (*structures.ResponseProject, error)
	GetProjectIssues(project string) ([]structures.JiraIssue, error)
	GetProjectByKey(projectKey string) (*structures.JiraProject, error)
}

type DataTransformerInterface interface {
	TransformStatusDB(jiraChanges *structures.Changelog) map[string]structures.DBStatusChanges
	TransformAuthorDB(jiraAuthor *structures.User) *structures.DBAuthor
	TransformProjectDB(jiraProject *structures.JiraProject) *structures.DBProject
	TransformIssueDB(jiraIssue *structures.JiraIssue) *structures.DBIssue
	TransformToDbIssueSet(project *structures.JiraProject, jiraIssue *structures.JiraIssue) *datatransformer.DataTransformer
}

type DbPusherInterface interface {
	PushProject(project *structures.DBProject) (int, error)
	PushProjects(projects []structures.DBProject) error
	PushStatusChanges(issue int, changes *datatransformer.DataTransformer) error
	PushIssue(project *structures.DBProject, issue *datatransformer.DataTransformer) (int, error)
	PushIssues(project *structures.DBProject, issues []datatransformer.DataTransformer) error
	Close()
}

type JiraService struct {
	jiraConnector   JiraConnectorInterface
	dataTransformer DataTransformerInterface
	dbPusher        DbPusherInterface
	log             *slog.Logger
}

func NewJiraService(
	config *config.Config,
	jiraConnector JiraConnectorInterface,
	dataTransformer DataTransformerInterface,
	dbPusher DbPusherInterface,
	log *slog.Logger) (*JiraService, error) {
	return &JiraService{
		jiraConnector:   jiraConnector,
		dataTransformer: dataTransformer,
		dbPusher:        dbPusher,
		log:             log,
	}, nil
}

func (js *JiraService) GetProjectsPage(search string, limit, page int) (*structures.ResponseProject, error) {
	js.log.Info("get project page", "page", page, "search", search, "limit", limit)
	return js.jiraConnector.GetProjectsPage(search, limit, page)
}
func (js *JiraService) UpdateProjects(projectId string) ([]structures.JiraIssue, error) {
	js.log.Info("upd project page", "projectId", projectId)
	return js.jiraConnector.GetProjectIssues(projectId)
}

func (js *JiraService) PushDataToDb(project string, issues []structures.JiraIssue) error {
	prj, err := js.jiraConnector.GetProjectByKey(project)
	if err != nil {
		js.log.Error("error Get Project By Key", logger.Err(err))
		return fmt.Errorf("%w", err)
	}
	data := js.TransformDataToDb(prj, issues)
	prjDB := js.dataTransformer.TransformProjectDB(prj)
	if err := js.dbPusher.PushIssues(prjDB, data); err != nil {
		js.log.Error("error push issues", logger.Err(err))
		return fmt.Errorf("%w", err)
	}

	js.log.Info("push data to db", "project", project)

	return nil

}

func (js *JiraService) TransformDataToDb(project *structures.JiraProject, issues []structures.JiraIssue) []datatransformer.DataTransformer {
	var issuesDb []datatransformer.DataTransformer

	for _, issue := range issues {
		issuesDb = append(issuesDb, *js.dataTransformer.TransformToDbIssueSet(project, &issue))
	}

	js.log.Info("transform data for db", "project", project)

	return issuesDb
}
