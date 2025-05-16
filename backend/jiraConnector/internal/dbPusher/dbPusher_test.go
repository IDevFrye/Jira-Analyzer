package dbpusher

import (
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	datatransformer "github.com/jiraconnector/internal/dataTransformer"
	myerr "github.com/jiraconnector/internal/dbPusher/errors"
	"github.com/jiraconnector/internal/structures"
	"github.com/stretchr/testify/assert"
)

func TestPushProject(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbp := &DbPusher{db: db, log: slog.Default()}

	tests := []struct {
		name       string
		title      string
		mockQuery  func()
		wantErr    bool
		wantResult int
	}{
		{
			name:  "success insert",
			title: "Test Project",
			mockQuery: func() {
				mock.ExpectQuery("INSERT INTO projects").
					WithArgs("Test Project", "").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			wantErr:    false,
			wantResult: 1,
		},
		{
			name:  "insert error",
			title: "Bad Project",
			mockQuery: func() {
				mock.ExpectQuery("INSERT INTO projects").
					WithArgs("Bad Project", "").
					WillReturnError(fmt.Errorf("insert error"))
			},
			wantErr:    true,
			wantResult: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockQuery()
			id, err := dbp.PushProject(structures.DBProject{Title: tt.title})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantResult, id)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPushProjects(t *testing.T) {

	tests := []struct {
		name        string
		projects    []structures.DBProject
		mockSetup   func(*sqlmock.Sqlmock)
		wantErr     bool
		expectedErr error
	}{
		{
			name: "successful transaction with multiple projects",
			projects: []structures.DBProject{
				{Title: "Project A"},
				{Title: "Project B"},
			},
			mockSetup: func(m *sqlmock.Sqlmock) {
				(*m).ExpectBegin()
				(*m).ExpectQuery("INSERT INTO projects").
					WithArgs("Project A", "").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				(*m).ExpectQuery("INSERT INTO projects").
					WithArgs("Project B", "").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				(*m).ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "failed to begin transaction",
			projects: []structures.DBProject{
				{Title: "Project A"},
			},
			mockSetup: func(m *sqlmock.Sqlmock) {
				(*m).ExpectBegin().WillReturnError(errors.New("begin error"))
			},
			wantErr:     true,
			expectedErr: myerr.ErrTranBegin,
		},
		{
			name: "failed to insert project - rollback",
			projects: []structures.DBProject{
				{Title: "Project A"},
				{Title: "Project B"},
			},
			mockSetup: func(m *sqlmock.Sqlmock) {
				(*m).ExpectBegin()
				(*m).ExpectQuery("INSERT INTO projects").
					WithArgs("Project A", "").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				(*m).ExpectQuery("INSERT INTO projects").
					WithArgs("Project B", "").
					WillReturnError(errors.New("insert error"))
				(*m).ExpectRollback()
			},
			wantErr:     true,
			expectedErr: myerr.ErrPushProject,
		},
		{
			name: "failed to commit transaction",
			projects: []structures.DBProject{
				{Title: "Project A"},
			},
			mockSetup: func(m *sqlmock.Sqlmock) {
				(*m).ExpectBegin()
				(*m).ExpectQuery("INSERT INTO projects").
					WithArgs("Project A", "").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				(*m).ExpectCommit().WillReturnError(errors.New("commit error"))
			},
			wantErr:     true,
			expectedErr: myerr.ErrTranClose,
		},
		{
			name:     "empty projects list",
			projects: []structures.DBProject{},
			mockSetup: func(m *sqlmock.Sqlmock) {
				(*m).ExpectBegin()
				(*m).ExpectCommit()
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			if tt.mockSetup != nil {
				tt.mockSetup(&mock)
			}

			dbp := &DbPusher{db: db, log: slog.Default()}

			err = dbp.PushProjects(tt.projects)

			if (err != nil) != tt.wantErr {
				t.Errorf("PushProjects() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.expectedErr != nil {
				if !errors.Is(err, tt.expectedErr) {
					t.Errorf("PushProjects() error = %v, expectedErr %v", err, tt.expectedErr)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestPushAuthor(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbp := &DbPusher{db: db, log: slog.Default()}

	tests := []struct {
		name       string
		author     structures.DBAuthor
		mockQuery  func()
		wantErr    bool
		wantResult int
	}{
		{
			name:   "success insert",
			author: structures.DBAuthor{Name: "Author1"},
			mockQuery: func() {
				mock.ExpectQuery("INSERT INTO author").
					WithArgs("Author1").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			wantErr:    false,
			wantResult: 1,
		},
		{
			name:   "insert error",
			author: structures.DBAuthor{Name: "Author2"},
			mockQuery: func() {
				mock.ExpectQuery("INSERT INTO author").
					WithArgs("Author2").
					WillReturnError(errors.New("db error"))
			},
			wantErr:    true,
			wantResult: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockQuery()
			id, err := dbp.PushAuthor(tt.author)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantResult, id)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestHasStatusChange(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbp := &DbPusher{db: db, log: slog.Default()}

	timeVal := time.Now()

	tests := []struct {
		name      string
		issueID   int
		time      time.Time
		mockQuery func()
		want      bool
	}{
		{
			name:    "found",
			issueID: 1,
			time:    timeVal,
			mockQuery: func() {
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM statuschanges`).
					WithArgs(1, timeVal).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			want: true,
		},
		{
			name:    "not found",
			issueID: 2,
			time:    timeVal,
			mockQuery: func() {
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM statuschanges`).
					WithArgs(2, timeVal).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			want: false,
		},
		{
			name:    "query error",
			issueID: 3,
			time:    timeVal,
			mockQuery: func() {
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM statuschanges`).
					WithArgs(3, timeVal).
					WillReturnError(errors.New("db error"))
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockQuery()
			result := dbp.hasStatusChange(tt.issueID, tt.time)
			assert.Equal(t, tt.want, result)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPushStatusChanges(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbp := &DbPusher{db: db, log: slog.Default()}

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

func TestPushIssue(t *testing.T) {
	testIssue := datatransformer.DataTransformer{
		Issue: structures.DBIssue{
			Key:         "PRJ-1",
			Summary:     "Test issue",
			Description: "Test description",
			Type:        "Task",
			Priority:    "High",
			Status:      "Open",
			CreatedTime: time.Now(),
		},
		Author:   structures.DBAuthor{Name: "user1"},
		Assignee: structures.DBAuthor{Name: "user2"},
	}
	tests := []struct {
		name          string
		project       string
		issue         datatransformer.DataTransformer
		mockSetup     func(*sqlmock.Sqlmock)
		expectedId    int
		expectedError error
	}{
		{
			name:    "successful insert new issue",
			project: "Project1",
			issue:   testIssue,
			mockSetup: func(m *sqlmock.Sqlmock) {
				// Mock getProjectId - сначала SELECT возвращает 0, потом INSERT
				//(*m).ExpectQuery(`SELECT id FROM projects WHERE title=\$1`).WithArgs("Project1").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(0))
				(*m).ExpectQuery(regexp.QuoteMeta(`
					INSERT INTO projects (title, url) VALUES ($1, $2) 
					ON CONFLICT (title) DO NOTHING RETURNING id`)).
					WithArgs("Project1", "").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				// Mock getAuthorId (author) - сначала SELECT возвращает 0, потом INSERT
				(*m).ExpectQuery(`SELECT id FROM author WHERE name=\$1`).
					WithArgs("user1").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(0))
				(*m).ExpectQuery(regexp.QuoteMeta(`INSERT INTO author (name) VALUES ($1) ON CONFLICT (name) DO NOTHING RETURNING id`)).
					WithArgs("user1").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

				// Mock getAuthorId (assignee) - сначала SELECT возвращает 0, потом INSERT
				(*m).ExpectQuery(`SELECT id FROM author WHERE name=\$1`).
					WithArgs("user2").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(0))
				(*m).ExpectQuery(regexp.QuoteMeta(`INSERT INTO author (name) VALUES ($1) ON CONFLICT (name) DO NOTHING RETURNING id`)).
					WithArgs("user2").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(3))

				// Mock insert issue
				(*m).ExpectQuery(regexp.QuoteMeta(`
                   INSERT INTO issue
                       (projectId, authorId, assigneeId, key, summary, description, type, priority, status, createdTime, closedTime, updatedTime, timeSpent)
                   VALUES
                       ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
                   ON CONFLICT (key)
                   DO UPDATE SET
                       projectId = EXCLUDED.projectId,
                       authorId = EXCLUDED.authorId,
                       assigneeId = EXCLUDED.assigneeId,
                       summary = EXCLUDED.summary,
                       description = EXCLUDED.description,
                       type = EXCLUDED.type,
                       priority = EXCLUDED.priority,
                       status = EXCLUDED.status,
                       createdTime = EXCLUDED.createdTime,
                       closedTime = EXCLUDED.closedTime,
                       updatedTime = EXCLUDED.updatedTime,
                       timeSpent = EXCLUDED.timeSpent
                   RETURNING id
               `)).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(100))
			},
			expectedId: 100,
		},
		{
			name:    "project found in db",
			project: "Project1",
			issue:   testIssue,
			mockSetup: func(m *sqlmock.Sqlmock) {
				// Project уже существует
				(*m).ExpectQuery(regexp.QuoteMeta(`
					INSERT INTO projects (title, url) VALUES ($1, $2) 
					ON CONFLICT (title) DO NOTHING RETURNING id`)).
					WithArgs("Project1", "").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				// Авторы не существуют
				// Mock getAuthorId (author) - сначала SELECT возвращает 0, потом INSERT
				(*m).ExpectQuery(regexp.QuoteMeta(`SELECT id FROM author WHERE name=$1`)).WithArgs("user1").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(0))
				(*m).ExpectQuery(regexp.QuoteMeta(`INSERT INTO author (name) VALUES ($1) ON CONFLICT (name) DO NOTHING RETURNING id`)).WithArgs("user1").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

				// Mock getAuthorId (assignee) - сначала SELECT возвращает 0, потом INSERT
				(*m).ExpectQuery(regexp.QuoteMeta(`SELECT id FROM author WHERE name=$1`)).WithArgs("user2").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(0))
				(*m).ExpectQuery(regexp.QuoteMeta(`INSERT INTO author (name) VALUES ($1) ON CONFLICT (name) DO NOTHING RETURNING id`)).WithArgs("user2").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(3))

				// Insert issue
				(*m).ExpectQuery(regexp.QuoteMeta(`INSERT INTO issue`)).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(100))
			},
			expectedId: 100,
		},
		{
			name:    "failed to insert project",
			project: "Project1",
			issue: datatransformer.DataTransformer{
				Issue: structures.DBIssue{Key: "PRJ-1"},
			},
			mockSetup: func(m *sqlmock.Sqlmock) {
				// Project не найден
				//(*m).ExpectQuery(`SELECT id FROM projects WHERE title=\$1`).WithArgs("Project1").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(0))
				(*m).ExpectQuery(regexp.QuoteMeta(`
					INSERT INTO projects (title, url) VALUES ($1, $2) 
					ON CONFLICT (title) DO NOTHING RETURNING id`)).
					WithArgs("Project1", "").
					WillReturnError(myerr.ErrInsertProject)
			},
			expectedError: myerr.ErrSelectProject,
		},
		{
			name:    "failed to insert author",
			project: "Project1",
			issue: datatransformer.DataTransformer{
				Issue:  structures.DBIssue{Key: "PRJ-1"},
				Author: structures.DBAuthor{Name: "user1"},
			},
			mockSetup: func(m *sqlmock.Sqlmock) {
				// Project успешно находится
				(*m).ExpectQuery(regexp.QuoteMeta(`
					INSERT INTO projects (title, url) VALUES ($1, $2) 
					ON CONFLICT (title) DO NOTHING RETURNING id`)).
					WithArgs("Project1", "").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				// Author не найден и ошибка при вставке
				(*m).ExpectQuery(regexp.QuoteMeta(`SELECT id FROM author WHERE name=$1`)).
					WithArgs("user1").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(0))
				(*m).ExpectQuery(regexp.QuoteMeta(`INSERT INTO author (name) VALUES ($1) ON CONFLICT (name) DO NOTHING RETURNING id`)).
					WithArgs("user1").
					WillReturnError(myerr.ErrInsertAuthor)
			},
			expectedError: myerr.ErrSelectAuthor,
		},
		{
			name:    "failed to insert issue",
			project: "Project1",
			issue: datatransformer.DataTransformer{
				Issue: structures.DBIssue{
					Key:     "PRJ-1",
					Summary: "Test",
				},
				Author:   structures.DBAuthor{Name: "user1"},
				Assignee: structures.DBAuthor{Name: "user2"},
			},
			mockSetup: func(m *sqlmock.Sqlmock) {
				// Project успешно находится
				(*m).ExpectQuery(regexp.QuoteMeta(`
					INSERT INTO projects (title, url) VALUES ($1, $2) 
					ON CONFLICT (title) DO NOTHING RETURNING id`)).
					WithArgs("Project1", "").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				// Author успешно находится/вставляется
				(*m).ExpectQuery("SELECT id FROM author WHERE name=\\$1").
					WithArgs("user1").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

				// Assignee успешно находится/вставляется
				(*m).ExpectQuery("SELECT id FROM author WHERE name=\\$1").
					WithArgs("user2").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(3))

				// Ошибка при вставке issue
				(*m).ExpectQuery(regexp.QuoteMeta(`INSERT INTO issue`)).
					WillReturnError(errors.New("insert failed"))
			},
			expectedError: myerr.ErrInsertIssue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			if tt.mockSetup != nil {
				tt.mockSetup(&mock)
			}

			dbp := &DbPusher{db: db, log: slog.Default()}
			id, err := dbp.PushIssue(structures.DBProject{Title: tt.project}, tt.issue)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.expectedError), "expected error %v, got %v", tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedId, id)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
func TestPushIssues(t *testing.T) {
	now := time.Now()

	testStatusChange1 := map[string]structures.DBStatusChanges{
		"user1": {
			IssueId:    100,
			AuthorId:   1,
			ChangeTime: now.Add(60 * time.Minute),
			FromStatus: "process",
			ToStatus:   "approved",
		},
		"user2": {
			IssueId:    100,
			AuthorId:   2,
			ChangeTime: now.Add(120 * time.Minute),
			FromStatus: "process",
			ToStatus:   "process",
		},
	}
	testStatusChange2 := map[string]structures.DBStatusChanges{
		"user1": {
			IssueId:    101,
			AuthorId:   1,
			ChangeTime: now.Add(60 * time.Minute),
			FromStatus: "process",
			ToStatus:   "approved",
		},
		"user2": {
			IssueId:    101,
			AuthorId:   2,
			ChangeTime: now.Add(120 * time.Minute),
			FromStatus: "process",
			ToStatus:   "process",
		},
	}
	testIssues := []datatransformer.DataTransformer{
		{
			Issue: structures.DBIssue{
				Key:         "PRJ-1",
				Summary:     "Test issue",
				Description: "Test description",
				Type:        "Task",
				Priority:    "High",
				Status:      "Open",
				CreatedTime: now,
			},
			Author:        structures.DBAuthor{Name: "user1"},
			Assignee:      structures.DBAuthor{Name: "user2"},
			StatusChanges: testStatusChange1,
		},
		{
			Issue: structures.DBIssue{
				Key:         "PRJ-2",
				Summary:     "Test issue",
				Description: "Test description",
				Type:        "Task",
				Priority:    "High",
				Status:      "Open",
				CreatedTime: now,
			},
			Author:        structures.DBAuthor{Name: "user1"},
			Assignee:      structures.DBAuthor{Name: "user2"},
			StatusChanges: testStatusChange2,
		},
	}

	tests := []struct {
		name          string
		project       string
		issues        []datatransformer.DataTransformer
		mockSetup     func(*sqlmock.Sqlmock)
		expectedError error
	}{
		{
			name:    "successful insert multiple issues",
			project: "Project1",
			issues:  testIssues,
			mockSetup: func(m *sqlmock.Sqlmock) {
				// Begin transaction
				(*m).ExpectBegin()

				// First issue
				(*m).ExpectQuery(regexp.QuoteMeta(`
					INSERT INTO projects (title, url) VALUES ($1, $2) 
					ON CONFLICT (title) DO NOTHING RETURNING id`)).
					WithArgs("Project1", "").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				// Author queries for first issue
				(*m).ExpectQuery(regexp.QuoteMeta(`SELECT id FROM author WHERE name=$1`)).
					WithArgs("user1").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(0))
				(*m).ExpectQuery(regexp.QuoteMeta(`INSERT INTO author (name) VALUES ($1) ON CONFLICT (name) DO NOTHING RETURNING id`)).
					WithArgs("user1").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				(*m).ExpectQuery(regexp.QuoteMeta(`SELECT id FROM author WHERE name=$1`)).
					WithArgs("user2").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(0))
				(*m).ExpectQuery(regexp.QuoteMeta(`INSERT INTO author (name) VALUES ($1) ON CONFLICT (name) DO NOTHING RETURNING id`)).
					WithArgs("user2").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

				// Insert first issue
				(*m).ExpectQuery(regexp.QuoteMeta(`INSERT INTO issue`)).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(100))

				// Status changes for first issue
				for author, sc := range testStatusChange1 {
					(*m).ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(*) FROM statuschanges WHERE issueId=$1 AND changeTime=$2`)).
						WithArgs(sc.IssueId, sc.ChangeTime).
						WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

					// Author already exists for status changes
					(*m).ExpectQuery(regexp.QuoteMeta(`SELECT id FROM author WHERE name=$1`)).
						WithArgs(author).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(sc.AuthorId))

					(*m).ExpectExec(regexp.QuoteMeta(`INSERT INTO statuschanges`)).
						WithArgs(sc.IssueId, sc.AuthorId, sc.ChangeTime, sc.FromStatus, sc.ToStatus).
						WillReturnResult(sqlmock.NewResult(1, 1))
				}

				// Second issue - project already exists
				(*m).ExpectQuery(regexp.QuoteMeta(`
					INSERT INTO projects (title, url) VALUES ($1, $2) 
					ON CONFLICT (title) DO NOTHING RETURNING id`)).
					WithArgs("Project1", "").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				// Authors already exist for second issue
				(*m).ExpectQuery(regexp.QuoteMeta(`SELECT id FROM author WHERE name=$1`)).
					WithArgs("user1").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				(*m).ExpectQuery(regexp.QuoteMeta(`SELECT id FROM author WHERE name=$1`)).
					WithArgs("user2").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

				// Insert second issue
				(*m).ExpectQuery(regexp.QuoteMeta(`INSERT INTO issue`)).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(101))

				// Status changes for second issue
				for author, sc := range testStatusChange2 {
					(*m).ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(*) FROM statuschanges WHERE issueId=$1 AND changeTime=$2`)).
						WithArgs(sc.IssueId, sc.ChangeTime).
						WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

					// Author already exists for status changes
					(*m).ExpectQuery(regexp.QuoteMeta(`SELECT id FROM author WHERE name=$1`)).
						WithArgs(author).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(sc.AuthorId))

					(*m).ExpectExec(regexp.QuoteMeta(`INSERT INTO statuschanges`)).
						WithArgs(sc.IssueId, sc.AuthorId, sc.ChangeTime, sc.FromStatus, sc.ToStatus).
						WillReturnResult(sqlmock.NewResult(1, 1))
				}

				// Commit
				(*m).ExpectCommit()
			},
		},
		{
			name:    "failed to begin transaction",
			project: "Project1",
			issues:  testIssues,
			mockSetup: func(m *sqlmock.Sqlmock) {
				(*m).ExpectBegin().WillReturnError(errors.New("begin error"))
			},
			expectedError: myerr.ErrTranBegin,
		},
		{
			name:    "failed to push issue - rollback",
			project: "Project1",
			issues:  testIssues,
			mockSetup: func(m *sqlmock.Sqlmock) {
				(*m).ExpectBegin()
				(*m).ExpectQuery(regexp.QuoteMeta(`
					INSERT INTO projects (title, url) VALUES ($1, $2) 
					ON CONFLICT (title) DO NOTHING RETURNING id`)).
					WithArgs("Project1", "").
					WillReturnError(errors.New("project error"))
				(*m).ExpectRollback()
			},
			expectedError: myerr.ErrPushIssue,
		},
		{
			name:    "failed to push status changes - rollback",
			project: "Project1",
			issues:  testIssues,
			mockSetup: func(m *sqlmock.Sqlmock) {
				(*m).ExpectBegin()
				// First issue
				(*m).ExpectQuery(`SELECT id FROM projects WHERE title=\$1`).
					WithArgs("Project1").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(0))
				(*m).ExpectQuery(regexp.QuoteMeta(`INSERT INTO projects (title, url) VALUES ($1, $2) ON CONFLICT (title) DO NOTHING RETURNING id`)).
					WithArgs("Project1", "").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				// Author queries for first issue
				(*m).ExpectQuery(regexp.QuoteMeta(`SELECT id FROM author WHERE name=$1`)).
					WithArgs("user1").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(0))
				(*m).ExpectQuery(regexp.QuoteMeta(`INSERT INTO author (name) VALUES ($1) ON CONFLICT (name) DO NOTHING RETURNING id`)).
					WithArgs("user1").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				(*m).ExpectQuery(regexp.QuoteMeta(`SELECT id FROM author WHERE name=$1`)).
					WithArgs("user2").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(0))
				(*m).ExpectQuery(regexp.QuoteMeta(`INSERT INTO author (name) VALUES ($1) ON CONFLICT (name) DO NOTHING RETURNING id`)).
					WithArgs("user2").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

				// Insert first issue
				(*m).ExpectQuery(regexp.QuoteMeta(`INSERT INTO issue`)).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(100))

				// Status change fails
				(*m).ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(*) FROM statuschanges WHERE issueId=$1 AND changeTime=$2`)).
					WithArgs(testStatusChange1["user1"].IssueId, testStatusChange1["user1"].ChangeTime).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

				// Author already exists for status changes
				(*m).ExpectQuery(regexp.QuoteMeta(`SELECT id FROM author WHERE name=$1`)).
					WithArgs("user1").
					WillReturnError(myerr.ErrInsertAuthor)

				(*m).ExpectRollback()
			},
			expectedError: myerr.ErrInsertStatusChange,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			if tt.mockSetup != nil {
				tt.mockSetup(&mock)
			}

			dbp := &DbPusher{db: db, log: slog.Default()}
			err = dbp.PushIssues(structures.DBProject{Title: tt.project}, tt.issues)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.expectedError), "expected error %v, got %v", tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
