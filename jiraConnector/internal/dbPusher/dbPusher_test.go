package dbpusher

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	datatransformer "github.com/jiraconnector/internal/dataTransformer"
	"github.com/jiraconnector/internal/structures"
	"github.com/stretchr/testify/assert"
)

func TestPushProject(t *testing.T) {
	// Create mock for database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbp := &DbPusher{db: db}

	// Set requests
	mock.ExpectQuery("INSERT INTO projects").
		WithArgs("Test Project").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// CheckFunc
	projectID, err := dbp.PushProject(structures.DBProject{Title: "Test Project"})

	assert.NoError(t, err)
	assert.Equal(t, 1, projectID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPushAuthor(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbp := &DbPusher{db: db}

	mock.ExpectQuery("INSERT INTO author").
		WithArgs("Test Author").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	authorId, err := dbp.PushAuthor(structures.DBAuthor{Name: "Test Author"})

	assert.NoError(t, err)
	assert.Equal(t, 1, authorId)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPushIssue(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbp := &DbPusher{db: db}

	// Mock get projectId
	mock.ExpectQuery(`INSERT INTO projects \(title\) VALUES \(\$1\) ON CONFLICT \(title\) DO NOTHING RETURNING id`).
		WithArgs("Test Project").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Mock get authorId
	mock.ExpectQuery(`INSERT INTO author \(name\) VALUES \(\$1\) ON CONFLICT \(name\) DO NOTHING RETURNING id`).
		WithArgs("Test Author").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

	// Mock get assignId
	mock.ExpectQuery(`INSERT INTO author \(name\) VALUES \(\$1\) ON CONFLICT \(name\) DO NOTHING RETURNING id`).
		WithArgs("Test Assignee").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(3))

	// Set request
	mock.ExpectQuery("INSERT INTO issue").
		WithArgs(
			1, 2, 3, "Test Issue", "Test Summary", "Test Description", "Bug", "High", "Open",
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 3600).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(10))

	// Test issue
	issue := datatransformer.DataTransformer{
		Issue: structures.DBIssue{
			ProjectId:   1,
			AuthorId:    2,
			AssigneeId:  3,
			Key:         "Test Issue",
			Summary:     "Test Summary",
			Description: "Test Description",
			Type:        "Bug",
			Priority:    "High",
			Status:      "Open",
			CreatedTime: time.Now(),
			ClosedTime:  sql.NullTime{}.Time,
			UpdatedTime: time.Now(),
			TimeSpent:   3600,
		},
		Author:   structures.DBAuthor{Name: "Test Author"},
		Assignee: structures.DBAuthor{Name: "Test Assignee"},
	}

	issueID, err := dbp.PushIssue("Test Project", issue)

	assert.NoError(t, err)
	assert.Equal(t, 10, issueID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPushStatusChanges(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbp := &DbPusher{db: db}

	issueID := 123
	authorName := "John Doe"
	changeTime := time.Now()
	fromStatus := "Open"
	toStatus := "In Progress"

	changes := datatransformer.DataTransformer{
		StatusChanges: map[string]structures.DBStatusChanges{
			authorName: {
				ChangeTime: changeTime,
				FromStatus: fromStatus,
				ToStatus:   toStatus,
			},
		},
	}

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM statuschanges WHERE issueId=\$1 AND changeTime=\$2`).
		WithArgs(issueID, changeTime).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	mock.ExpectQuery(`SELECT id FROM author WHERE name=\$1`).
		WithArgs(authorName).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(4))

	mock.ExpectExec(`INSERT INTO statuschanges \(issueId, authorId, changeTime, fromStatus, toStatus\) VALUES \(\$1, \$2, \$3, \$4, \$5\)`).
		WithArgs(issueID, 4, changeTime, fromStatus, toStatus).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = dbp.PushStatusChanges(issueID, changes)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
