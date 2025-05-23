//go:build integration
// +build integration

package dbintegrations

import (
	"context"
	"testing"
	"time"

	datatransformer "github.com/jiraconnector/internal/dataTransformer"
	"github.com/jiraconnector/internal/structures"
	"github.com/stretchr/testify/assert"
)

func TestLoadNewProject(t *testing.T) {
	resetTestDB(t)

	tx, err := DB.Db().Begin()
	if err != nil {
		t.Fatalf("failed to begin tx: %v", err)
	}
	defer tx.Rollback()

	testDB := *DB

	ctx := context.Background()

	testIssues, testProject := setupTestData()

	// 1. Проверяем, что проект отсутствует в БД перед загрузкой
	var count int
	err = tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM projects WHERE id = $1", testProject.Id).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 0, count, "Project should not exist before test")

	// 2. Загружаем проект
	_, err = testDB.PushProject(&testProject)
	assert.NoError(t, err)

	// 3. Загружаем задачи
	err = testDB.PushIssues(&testProject, testIssues)
	assert.NoError(t, err)

	// 4. Проверяем, что проект и задачи сохранились
	err = tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM projects WHERE id = $1", testProject.Id).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count, "Project should be saved")

	err = tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM issue WHERE projectId=$1", testProject.Id).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, len(testIssues), count, "All issues should be saved")
}

func TestUpdateExistingProject(t *testing.T) {
	resetTestDB(t)

	tx, err := DB.Db().Begin()
	if err != nil {
		t.Fatalf("failed to begin tx: %v", err)
	}
	defer tx.Rollback()

	testDB := *DB

	ctx := context.Background()

	testIssue, testProject := setupTestData()

	_, err = testDB.PushProject(&testProject)
	assert.NoError(t, err)

	err = testDB.PushIssues(&testProject, testIssue)
	assert.NoError(t, err)

	// 2. Подготавливаем обновленные данные
	updatedIssues := setUpdTestData()

	// 3. Обновляем проект (симулируем вызов /updateProject)
	// обновляем задачи
	err = testDB.PushIssues(&testProject, updatedIssues)
	assert.NoError(t, err)

	// 4. Проверяем результаты
	var count int

	// Проверяем количество задач
	err = tx.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM issue WHERE projectId=$1", testProject.Id).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, len(updatedIssues), count)

	// Проверяем добавление новой задачи
	err = tx.QueryRowContext(ctx,
		"SELECT 1 FROM issue WHERE id = $1", 5).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestProjectEdgeCases(t *testing.T) {
	resetTestDB(t)

	t.Run("DuplicateProject", func(t *testing.T) {
		project := structures.DBProject{Title: "Test"}

		// Первое сохранение должно пройти успешно
		_, err := DB.PushProject(&project)
		assert.NoError(t, err)

		_, err = DB.PushProject(&project)
		assert.Error(t, err, "Should reject duplicate project title")
	})
}

func setUpdTestData() []datatransformer.DataTransformer {
	now := time.Now()
	testStatusChange1 := map[string]structures.DBStatusChanges{
		"user1": {
			IssueId:    1,
			AuthorId:   1,
			ChangeTime: now.Add(60 * time.Minute),
			FromStatus: "process",
			ToStatus:   "approved",
		},
		"user2": {
			IssueId:    1,
			AuthorId:   2,
			ChangeTime: now.Add(120 * time.Minute),
			FromStatus: "process",
			ToStatus:   "process",
		},
	}
	testStatusChange2 := map[string]structures.DBStatusChanges{
		"user1": {
			IssueId:    2,
			AuthorId:   1,
			ChangeTime: now.Add(60 * time.Minute),
			FromStatus: "process",
			ToStatus:   "approved",
		},
		"user2": {
			IssueId:    2,
			AuthorId:   2,
			ChangeTime: now.Add(120 * time.Minute),
			FromStatus: "process",
			ToStatus:   "process",
		},
	}

	testStatusChange3 := map[string]structures.DBStatusChanges{
		"user1": {
			IssueId:    3,
			AuthorId:   1,
			ChangeTime: now.Add(160 * time.Minute),
			FromStatus: "process",
			ToStatus:   "approved",
		},
		"user2": {
			IssueId:    3,
			AuthorId:   2,
			ChangeTime: now.Add(170 * time.Minute),
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
				Status:      "Close",
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
		{
			Issue: structures.DBIssue{
				Key:         "PRJ-3",
				Summary:     "Test issue",
				Description: "Test description",
				Type:        "Task",
				Priority:    "High",
				Status:      "Open",
				CreatedTime: now,
			},
			Author:        structures.DBAuthor{Name: "user1"},
			Assignee:      structures.DBAuthor{Name: "user2"},
			StatusChanges: testStatusChange3,
		},
	}
	return testIssues
}

func setupTestData() ([]datatransformer.DataTransformer, structures.DBProject) {
	// Подготовка тестовых данных
	now := time.Now()
	testStatusChange1 := map[string]structures.DBStatusChanges{
		"user1": {
			IssueId:    1,
			AuthorId:   1,
			ChangeTime: now.Add(60 * time.Minute),
			FromStatus: "process",
			ToStatus:   "approved",
		},
		"user2": {
			IssueId:    1,
			AuthorId:   2,
			ChangeTime: now.Add(120 * time.Minute),
			FromStatus: "process",
			ToStatus:   "process",
		},
	}
	testStatusChange2 := map[string]structures.DBStatusChanges{
		"user1": {
			IssueId:    2,
			AuthorId:   1,
			ChangeTime: now.Add(60 * time.Minute),
			FromStatus: "process",
			ToStatus:   "approved",
		},
		"user2": {
			IssueId:    2,
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

	projectTestID := 1
	testProject := structures.DBProject{
		Id:    projectTestID,
		Title: "Test Project",
	}

	return testIssues, testProject
}
