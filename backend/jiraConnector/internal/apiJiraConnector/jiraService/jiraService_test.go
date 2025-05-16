package jiraservice

import (
	"errors"
	"fmt"
	"log/slog"
	"testing"

	datatransformer "github.com/jiraconnector/internal/dataTransformer"
	"github.com/jiraconnector/internal/structures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewJiraService(t *testing.T) {
	mockJiraConn := new(MockJiraConnectorInterface)
	mockTransformer := new(MockDataTransformerInterface)
	mockDbPusher := new(MockDbPusherInterface)

	service, err := NewJiraService(
		nil,
		mockJiraConn,
		mockTransformer,
		mockDbPusher,
		slog.Default(),
	)

	assert.NoError(t, err)
	assert.NotNil(t, service)
}

func TestGetProjectsPage(t *testing.T) {
	tests := []struct {
		name          string
		search        string
		limit         int
		page          int
		mockReturn    *structures.ResponseProject
		mockError     error
		expectedError error
	}{
		{
			name:          "success",
			search:        "test",
			limit:         10,
			page:          1,
			mockReturn:    &structures.ResponseProject{},
			mockError:     nil,
			expectedError: nil,
		},
		{
			name:          "connector error",
			search:        "test",
			limit:         10,
			page:          1,
			mockReturn:    nil,
			mockError:     errors.New("connector error"),
			expectedError: errors.New("connector error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockJiraConn := new(MockJiraConnectorInterface)
			mockJiraConn.On("GetProjectsPage", tt.search, tt.limit, tt.page).Return(tt.mockReturn, tt.mockError)

			service := JiraService{
				jiraConnector: mockJiraConn,
				log:           slog.Default(),
			}

			result, err := service.GetProjectsPage(tt.search, tt.limit, tt.page)

			assert.Equal(t, tt.mockReturn, result)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			mockJiraConn.AssertExpectations(t)
		})
	}
}

func TestUpdateProjects(t *testing.T) {
	tests := []struct {
		name          string
		projectId     string
		mockReturn    []structures.JiraIssue
		mockError     error
		expectedError error
	}{
		{
			name:          "success",
			projectId:     "TEST",
			mockReturn:    []structures.JiraIssue{},
			mockError:     nil,
			expectedError: nil,
		},
		{
			name:          "connector error",
			projectId:     "TEST",
			mockReturn:    nil,
			mockError:     errors.New("connector error"),
			expectedError: errors.New("connector error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockJiraConn := new(MockJiraConnectorInterface)
			mockJiraConn.On("GetProjectIssues", tt.projectId).Return(tt.mockReturn, tt.mockError)

			service := JiraService{
				jiraConnector: mockJiraConn,
				log:           slog.Default(),
			}

			result, err := service.UpdateProjects(tt.projectId)

			assert.Equal(t, tt.mockReturn, result)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			mockJiraConn.AssertExpectations(t)
		})
	}
}

func TestPushDataToDb(t *testing.T) {
	tests := []struct {
		name          string
		project       structures.JiraProject
		issues        []structures.JiraIssue
		mockTransform []*datatransformer.DataTransformer
		mockError     error
		expectedError string
	}{
		{
			name: "success",
			project: structures.JiraProject{
				Name: "TEST",
			},
			issues: []structures.JiraIssue{
				{Id: "1"},
				{Id: "2"},
			},
			mockTransform: []*datatransformer.DataTransformer{
				{},
				{},
			},
			mockError:     nil,
			expectedError: "",
		},
		{
			name: "db push error",
			project: structures.JiraProject{
				Name: "TEST",
			},
			issues: []structures.JiraIssue{
				{Id: "1"},
			},
			mockTransform: []*datatransformer.DataTransformer{
				{},
			},
			mockError:     errors.New("db error"),
			expectedError: "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTransformer := new(MockDataTransformerInterface)
			mockDbPusher := new(MockDbPusherInterface)

			mockTransformer.On("TransformProjectDB", tt.project).
				Return(structures.DBProject{Title: tt.project.Name, Url: fmt.Sprintf("/projects/%s", tt.project.Name)})

			for i, issue := range tt.issues {
				mockTransformer.On("TransformToDbIssueSet", tt.project, issue).Return(tt.mockTransform[i])
			}

			mockDbPusher.On("PushIssues",
				structures.DBProject{Title: tt.project.Name, Url: fmt.Sprintf("/projects/%s", tt.project.Name)},
				mock.AnythingOfType("[]datatransformer.DataTransformer")).Return(tt.mockError)

			service := JiraService{
				dataTransformer: mockTransformer,
				dbPusher:        mockDbPusher,
				log:             slog.Default(),
			}

			err := service.PushDataToDb(tt.project, tt.issues)

			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			mockTransformer.AssertExpectations(t)
			mockDbPusher.AssertExpectations(t)
		})
	}
}

func TestTransformDataToDb(t *testing.T) {
	tests := []struct {
		name           string
		project        string
		issues         []structures.JiraIssue
		mockTransforms []*datatransformer.DataTransformer
		expectedResult []datatransformer.DataTransformer
	}{
		{
			name:    "single issue",
			project: "TEST",
			issues: []structures.JiraIssue{
				{Id: "1"},
			},
			mockTransforms: []*datatransformer.DataTransformer{
				{Issue: structures.DBIssue{Id: 1}},
			},
			expectedResult: []datatransformer.DataTransformer{
				{Issue: structures.DBIssue{Id: 1}},
			},
		},
		{
			name:    "multiple issues",
			project: "TEST",
			issues: []structures.JiraIssue{
				{Id: "1"},
				{Id: "2"},
			},
			mockTransforms: []*datatransformer.DataTransformer{
				{Issue: structures.DBIssue{Id: 1}},
				{Issue: structures.DBIssue{Id: 2}},
			},
			expectedResult: []datatransformer.DataTransformer{
				{Issue: structures.DBIssue{Id: 1}},
				{Issue: structures.DBIssue{Id: 2}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTransformer := new(MockDataTransformerInterface)

			for i, issue := range tt.issues {
				mockTransformer.On("TransformToDbIssueSet", structures.JiraProject{Name: tt.project}, issue).Return(tt.mockTransforms[i])
			}

			service := JiraService{
				dataTransformer: mockTransformer,
				log:             slog.Default(),
			}

			result := service.TransformDataToDb(structures.JiraProject{Name: tt.project}, tt.issues)

			assert.Equal(t, tt.expectedResult, result)
			mockTransformer.AssertExpectations(t)
		})
	}
}
