package datatransformer

import (
	"testing"
	"time"

	"github.com/jiraconnector/internal/structures"
	"github.com/stretchr/testify/assert"
)

func TestTransformStatusDB(t *testing.T) {

	tests := []struct {
		name     string
		input    structures.Changelog
		expected map[string]structures.DBStatusChanges
	}{
		{
			name: "single status change",
			input: structures.Changelog{
				Histories: []structures.History{
					{
						Created: "2023-01-01T10:00:00.000-0700",
						Author:  structures.User{Name: "user1"},
						Items: []structures.Item{
							{
								Field:      "status",
								FromString: "Open",
								ToString:   "In Progress",
							},
						},
					},
				},
			},
			expected: map[string]structures.DBStatusChanges{
				"user1": {
					ChangeTime: time.Date(2023, 1, 1, 10, 0, 0, 0, time.FixedZone("", -7*3600)),
					FromStatus: "Open",
					ToStatus:   "In Progress",
				},
			},
		},
		{
			name: "multiple changes, only status",
			input: structures.Changelog{
				Histories: []structures.History{
					{
						Created: "2023-01-01T10:00:00.000-0700",
						Author:  structures.User{Name: "user1"},
						Items: []structures.Item{
							{
								Field:      "status",
								FromString: "Open",
								ToString:   "In Progress",
							},
							{
								Field:      "priority",
								FromString: "Low",
								ToString:   "High",
							},
						},
					},
					{
						Created: "2023-01-02T11:00:00.000-0700",
						Author:  structures.User{Name: "user2"},
						Items: []structures.Item{
							{
								Field:      "status",
								FromString: "In Progress",
								ToString:   "Done",
							},
						},
					},
				},
			},
			expected: map[string]structures.DBStatusChanges{
				"user1": {
					ChangeTime: time.Date(2023, 1, 1, 10, 0, 0, 0, time.FixedZone("", -7*3600)),
					FromStatus: "Open",
					ToStatus:   "In Progress",
				},
				"user2": {
					ChangeTime: time.Date(2023, 1, 2, 11, 0, 0, 0, time.FixedZone("", -7*3600)),
					FromStatus: "In Progress",
					ToStatus:   "Done",
				},
			},
		},
		{
			name:     "empty changelog",
			input:    structures.Changelog{},
			expected: map[string]structures.DBStatusChanges{},
		},
	}

	dt := NewDataTransformer("base_url")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := dt.TransformStatusDB(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTransformAuthorDB(t *testing.T) {
	tests := []struct {
		name     string
		input    structures.User
		expected structures.DBAuthor
	}{
		{
			name:     "regular user",
			input:    structures.User{Name: "john.doe"},
			expected: structures.DBAuthor{Name: "john.doe"},
		},
		{
			name:     "empty user",
			input:    structures.User{},
			expected: structures.DBAuthor{},
		},
	}

	dt := NewDataTransformer("base_url")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := dt.TransformAuthorDB(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTransformProjectDB(t *testing.T) {
	tests := []struct {
		name     string
		input    structures.JiraProject
		expected structures.DBProject
	}{
		{
			name:     "regular project",
			input:    structures.JiraProject{Name: "Project X"},
			expected: structures.DBProject{Title: "Project X", Url: "/projects/Project_X"},
		},
		{
			name:     "empty project",
			input:    structures.JiraProject{},
			expected: structures.DBProject{Url: "/projects/"},
		},
	}

	dt := NewDataTransformer("")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := dt.TransformProjectDB(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTransformIssueDB(t *testing.T) {
	createdTime := "2023-01-01T10:00:00.000-0700"
	updatedTime := "2023-01-02T11:00:00.000-0700"
	closedTime := "2023-01-03T12:00:00.000-0700"

	parsedCreated, _ := time.Parse("2006-01-02T15:04:05.000-0700", createdTime)
	parsedUpdated, _ := time.Parse("2006-01-02T15:04:05.000-0700", updatedTime)
	parsedClosed, _ := time.Parse("2006-01-02T15:04:05.000-0700", closedTime)

	tests := []struct {
		name     string
		input    structures.JiraIssue
		expected structures.DBIssue
	}{
		{
			name: "full issue data",
			input: structures.JiraIssue{
				Key: "PRJ-123",
				Fields: structures.Field{
					Summary:     "Test issue",
					Description: "Test description",
					Type:        structures.IssueType{Description: "Task"},
					Project:     structures.JiraProject{Name: "Project X"},
					Status:      structures.IssueStatus{Name: "Done"},
					CreatedTime: createdTime,
					UpdatedTime: updatedTime,
					ClosedTime:  closedTime,
					TimeSpent:   3600,
					Author:      structures.User{Name: "author"},
					Assignee:    structures.User{Name: "assignee"},
				},
			},
			expected: structures.DBIssue{
				Key:         "PRJ-123",
				Summary:     "Test issue",
				Description: "Test description",
				Type:        "Task",
				Priority:    "Project X",
				Status:      "Done",
				CreatedTime: parsedCreated,
				UpdatedTime: parsedUpdated,
				ClosedTime:  parsedClosed,
				TimeSpent:   3600,
			},
		},
		{
			name: "minimum issue data",
			input: structures.JiraIssue{
				Key: "PRJ-124",
				Fields: structures.Field{
					Summary: "Minimal issue",
				},
			},
			expected: structures.DBIssue{
				Key:     "PRJ-124",
				Summary: "Minimal issue",
			},
		},
	}

	dt := NewDataTransformer("base_url")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := dt.TransformIssueDB(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTransformToDbIssueSet(t *testing.T) {
	createdTime := "2023-01-01T10:00:00.000-0700"
	parsedCreated, _ := time.Parse("2006-01-02T15:04:05.000-0700", createdTime)

	inputIssue := structures.JiraIssue{
		Key: "PRJ-123",
		Fields: structures.Field{
			Summary:     "Test issue",
			CreatedTime: createdTime,
			Author:      structures.User{Name: "author"},
			Assignee:    structures.User{Name: "assignee"},
		},
		Changelog: structures.Changelog{
			Histories: []structures.History{
				{
					Created: createdTime,
					Author:  structures.User{Name: "user1"},
					Items: []structures.Item{
						{
							Field:      "status",
							FromString: "Open",
							ToString:   "In Progress",
						},
					},
				},
			},
		},
	}

	expected := &DataTransformer{
		Project: structures.DBProject{Title: "TestProject", Url: "base_url/projects/TestProject"},
		Issue: structures.DBIssue{
			Key:         "PRJ-123",
			Summary:     "Test issue",
			CreatedTime: parsedCreated,
		},
		Author:   structures.DBAuthor{Name: "author"},
		Assignee: structures.DBAuthor{Name: "assignee"},
		StatusChanges: map[string]structures.DBStatusChanges{
			"user1": {
				ChangeTime: parsedCreated,
				FromStatus: "Open",
				ToStatus:   "In Progress",
			},
		},
	}

	dt := NewDataTransformer("base_url")
	result := dt.TransformToDbIssueSet(structures.JiraProject{Name: "TestProject"}, inputIssue)

	assert.Equal(t, expected.Project, result.Project)
	assert.Equal(t, expected.Issue.Key, result.Issue.Key)
	assert.Equal(t, expected.Issue.Summary, result.Issue.Summary)
	assert.Equal(t, expected.Author, result.Author)
	assert.Equal(t, expected.Assignee, result.Assignee)
	assert.Equal(t, expected.StatusChanges["user1"], result.StatusChanges["user1"])
}
