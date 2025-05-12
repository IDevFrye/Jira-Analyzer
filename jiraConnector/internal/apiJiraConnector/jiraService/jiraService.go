package jiraservice

import (
	"fmt"
	"log"

	configreader "github.com/jiraconnector/internal/configReader"
	datatransformer "github.com/jiraconnector/internal/dataTransformer"
	"github.com/jiraconnector/internal/structures"
)

//go:generate mockery

type JiraConnectorInterface interface {
	GetAllProjects() ([]structures.JiraProject, error)
	GetProjectsPage(search string, limit, page int) (*structures.ResponseProject, error)
	GetProjectIssues(project string) ([]structures.JiraIssue, error)
}

type DataTransformerInterface interface {
	TransformStatusDB(jiraChanges structures.Changelog) map[string]structures.DBStatusChanges
	TransformAuthorDB(jiraAuthor structures.User) structures.DBAuthor
	TransformProjectDB(jiraProject structures.JiraProject) structures.DBProject
	TransformIssueDB(jiraIssue structures.JiraIssue) structures.DBIssue
	TransformToDbIssueSet(projectName string, jiraIssue structures.JiraIssue) *datatransformer.DataTransformer
}

type DbPusherInterface interface {
	PushProject(project structures.DBProject) (int, error)
	PushProjects(projects []structures.DBProject) error
	PushStatusChanges(issue int, changes datatransformer.DataTransformer) error
	PushIssue(project string, issue datatransformer.DataTransformer) (int, error)
	PushIssues(project string, issues []datatransformer.DataTransformer) error
	Close()
}

type JiraService struct {
	jiraConnector   JiraConnectorInterface
	dataTransformer DataTransformerInterface
	dbPusher        DbPusherInterface
}

func NewJiraService(
	config *configreader.Config,
	jiraConnector JiraConnectorInterface,
	dataTransformer DataTransformerInterface,
	dbPusher DbPusherInterface) (*JiraService, error) {
	return &JiraService{
		jiraConnector:   jiraConnector,
		dataTransformer: dataTransformer,
		dbPusher:        dbPusher,
	}, nil
}

func (js JiraService) GetProjectsPage(search string, limit, page int) (*structures.ResponseProject, error) {
	return js.jiraConnector.GetProjectsPage(search, limit, page)
}
func (js JiraService) UpdateProjects(projectId string) ([]structures.JiraIssue, error) {
	return js.jiraConnector.GetProjectIssues(projectId)
}

func (js JiraService) PushDataToDb(project string, issues []structures.JiraIssue) error {
	data := js.TransformDataToDb(project, issues)

	if err := js.dbPusher.PushIssues(project, data); err != nil {
		log.Println(err)
		return fmt.Errorf("%w", err)
	}

	return nil

}

func (js JiraService) TransformDataToDb(project string, issues []structures.JiraIssue) []datatransformer.DataTransformer {
	var issuesDb []datatransformer.DataTransformer

	for _, issue := range issues {
		issuesDb = append(issuesDb, *js.dataTransformer.TransformToDbIssueSet(project, issue))
	}

	return issuesDb
}
