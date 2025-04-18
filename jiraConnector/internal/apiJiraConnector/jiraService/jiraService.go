package jiraservice

import (
	"fmt"
	"log"

	configreader "github.com/jiraconnector/internal/configReader"
	"github.com/jiraconnector/internal/connector"
	datatransformer "github.com/jiraconnector/internal/dataTransformer"
	dbpusher "github.com/jiraconnector/internal/dbPusher"
	"github.com/jiraconnector/internal/structures"
)

type JiraService struct {
	jiraConnector   connector.JiraConnector
	dataTransformer datatransformer.DataTransformer
	dbPusher        dbpusher.DbPusher
}

func NewJiraService(config configreader.Config, jiraConnector connector.JiraConnector, dbPusher dbpusher.DbPusher) (*JiraService, error) {
	return &JiraService{
		jiraConnector:   jiraConnector,
		dataTransformer: *datatransformer.NewDataTransformer(),
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
