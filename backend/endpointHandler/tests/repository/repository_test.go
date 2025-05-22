package repository_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/endpointhandler/model"
	"github.com/endpointhandler/repository"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setupMockDB(t *testing.T) (sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New() // теперь 3 значения
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}

	// Обернем *sql.DB в *sqlx.DB
	repository.DB = sqlx.NewDb(db, "sqlmock")

	return mock, func() {
		repository.DB.Close()
	}
}

func TestGetFilteredProjects(t *testing.T) {
	mock, closeDB := setupMockDB(t)
	defer closeDB()

	rows := sqlmock.NewRows([]string{"id", "key", "name", "self", "existence"}).
		AddRow(1, "PROJ1", "Project One", "", true).
		AddRow(2, "PROJ2", "Project Two", "", true)

	mock.ExpectQuery("SELECT id, title AS key, title AS name, '' AS self, TRUE as existence FROM Projects").
		WillReturnRows(rows)

	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM Projects").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	projects, total, err := repository.GetFilteredProjects(10, 0, "")
	assert.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, projects, 2)
}

func TestGetStats(t *testing.T) {
	mock, closeDB := setupMockDB(t)
	defer closeDB()

	projectID := 1

	// Настраиваем ожидания для каждого запроса подсчёта
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM Issue WHERE projectId=\\$1").
		WithArgs(projectID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(10))
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM Issue WHERE projectId=\\$1 AND status NOT IN").
		WithArgs(projectID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM Issue WHERE projectId=\\$1 AND status='Closed'").
		WithArgs(projectID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM StatusChanges WHERE issueId IN").
		WithArgs(projectID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM Issue WHERE projectId=\\$1 AND status='Resolved'").
		WithArgs(projectID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM Issue WHERE projectId=\\$1 AND status='In progress'").
		WithArgs(projectID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
	mock.ExpectQuery("SELECT COALESCE\\(AVG\\(EXTRACT\\(EPOCH FROM closedTime - createdTime\\)/3600\\), 0\\) FROM Issue WHERE projectId=\\$1 AND closedTime IS NOT NULL AND closedTime > createdTime").
		WithArgs(projectID).
		WillReturnRows(sqlmock.NewRows([]string{"avg"}).AddRow(24.5))
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) / 7.0 FROM Issue WHERE projectId=\\$1 AND createdTime > \\$2").
		WithArgs(projectID, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1.5))

	stats, err := repository.GetStats(projectID)
	assert.NoError(t, err)
	assert.Equal(t, 10, stats.TotalIssues)
	assert.Equal(t, 3, stats.OpenIssues)
	assert.Equal(t, 5, stats.ClosedIssues)
	assert.Equal(t, 1, stats.ReopenedIssues)
	assert.Equal(t, 1, stats.ResolvedIssues)
	assert.Equal(t, 2, stats.InProgressIssues)
	assert.InDelta(t, 24.5, stats.AvgResolutionTimeH, 0.001)
	assert.InDelta(t, 1.5, stats.AvgCreatedPerDay7d, 0.001)
}

func TestDeleteProject(t *testing.T) {
	mock, closeDB := setupMockDB(t)
	defer closeDB()

	mock.ExpectExec("DELETE FROM Projects WHERE id=\\$1").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repository.DeleteProject(1)
	assert.NoError(t, err)
}

func TestSaveProject(t *testing.T) {
	mock, closeDB := setupMockDB(t)
	defer closeDB()

	project := model.Project{Key: "KEY1", Name: "Project 1", Self: "http://url"}

	mock.ExpectExec("INSERT INTO Projects").
		WithArgs(project.Key, project.Name, project.Self).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repository.SaveProject(project)
	assert.NoError(t, err)
}
