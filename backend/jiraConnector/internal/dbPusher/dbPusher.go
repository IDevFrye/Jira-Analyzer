package dbpusher

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	datatransformer "github.com/jiraconnector/internal/dataTransformer"
	myerr "github.com/jiraconnector/internal/dbPusher/errors"
	"github.com/jiraconnector/internal/structures"
	"github.com/jiraconnector/pkg/config"
	"github.com/jiraconnector/pkg/logger"
	_ "github.com/lib/pq"
)

type DbPusher struct {
	db  *sql.DB
	log *slog.Logger
}

func NewDbPusher(cfg *config.Config, log *slog.Logger) (*DbPusher, error) {
	connStr := buildConnectionstring(&cfg.DBCfg)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		ansErr := fmt.Errorf("%w: %w", myerr.ErrOpenDb, err)
		log.Error(ansErr.Error())
		return nil, ansErr
	}

	return &DbPusher{
		db:  db,
		log: log,
	}, nil
}

func (dbp *DbPusher) Db() *sql.DB {
	return dbp.db
}

func (dbp *DbPusher) Close() {
	dbp.log.Info("close db connection")
	dbp.db.Close()
}

func (dbp *DbPusher) PushProject(project *structures.DBProject) (int, error) {
	var projectId int
	query := `INSERT INTO projects (title, key, url) VALUES ($1, $2, $3) ON CONFLICT (title) DO NOTHING RETURNING id`
	if err := dbp.db.QueryRow(query, project.Title, project.Key, project.Url).Scan(&projectId); err != nil {
		return 0, fmt.Errorf("%w - %s: %w", myerr.ErrInsertProject, project.Title, err)
	}
	dbp.log.Info("success push project", "project", project.Title)
	return projectId, nil
}

func (dbp *DbPusher) PushProjects(projects []structures.DBProject) error {
	tx, err := dbp.db.Begin()
	if err != nil {
		ansErr := fmt.Errorf("%w: %w", myerr.ErrTranBegin, err)
		dbp.log.Error(ansErr.Error())
		return ansErr
	}

	for _, project := range projects {
		_, err = dbp.PushProject(&project)
		if err != nil {
			ansErr := fmt.Errorf("%w - %s: %w", myerr.ErrPushProject, project.Title, err)
			dbp.log.Error(ansErr.Error())
			tx.Rollback()
			return ansErr
		}
	}

	if err := tx.Commit(); err != nil {
		ansErr := fmt.Errorf("%w: %w", myerr.ErrTranClose, err)
		dbp.log.Error(ansErr.Error())
		return ansErr
	}

	dbp.log.Error("success save all projects")
	return nil

}

func (dbp *DbPusher) PushAuthor(author *structures.DBAuthor) (int, error) {
	var authorId int
	query := "INSERT INTO author (name) VALUES ($1) ON CONFLICT (name) DO NOTHING RETURNING id"

	if err := dbp.db.QueryRow(query, author.Name).Scan(&authorId); err != nil {
		ansErr := fmt.Errorf("%w - %s: %w", myerr.ErrInsertAuthor, author.Name, err)
		dbp.log.Error(ansErr.Error())
		return 0, ansErr
	}

	dbp.log.Info("success push author", "author", author.Name)
	return authorId, nil

}

func (dbp *DbPusher) PushStatusChanges(issue int, changes *datatransformer.DataTransformer) error {
	query := "INSERT INTO statuschanges (issueId, authorId, changeTime, fromStatus, toStatus) VALUES ($1, $2, $3, $4, $5)"
	for author, statusChange := range changes.StatusChanges {
		if dbp.hasStatusChange(issue, statusChange.ChangeTime) {
			dbp.log.Warn("already has such status change")
			return nil
		}
		authorId, err := dbp.getAuthorId(&structures.DBAuthor{Name: author})
		if err != nil {
			dbp.log.Error("err get author Id", "author", author)
			return err
		}
		if _, err := dbp.db.Exec(query, issue, authorId, statusChange.ChangeTime, statusChange.FromStatus, statusChange.ToStatus); err != nil {
			dbp.log.Error("err insert status change", "author", author)
			return err
		}
	}

	dbp.log.Info("success push status changes")
	return nil
}

func (dbp *DbPusher) PushIssue(project *structures.DBProject, issue *datatransformer.DataTransformer) (int, error) {
	projectId, err := dbp.getProjectId(project)
	if err != nil {
		dbp.log.Error("err get project", "project", project)
		return 0, err
	}

	authorId, err := dbp.getAuthorId(&issue.Author)
	if err != nil {
		dbp.log.Error("err get author Id", "author", issue.Author.Name)
		return 0, err
	}

	assegneeId, err := dbp.getAuthorId(&issue.Assignee)
	if err != nil {
		dbp.log.Error("err get assignee Id", "author", issue.Assignee.Name)
		return 0, err
	}

	query := `
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
   `

	var issueId int
	iss := issue.Issue
	iss.ProjectId = projectId
	iss.AuthorId = authorId
	iss.AssigneeId = assegneeId

	if err := dbp.db.QueryRow(
		query, iss.ProjectId, iss.AuthorId, iss.AssigneeId,
		iss.Key, iss.Summary, iss.Description, iss.Type,
		iss.Priority, iss.Status, iss.CreatedTime,
		iss.ClosedTime, iss.UpdatedTime, iss.TimeSpent).Scan(&issueId); err != nil {

		ansErr := fmt.Errorf("%w - %s: %w", myerr.ErrInsertIssue, project.Title, err)
		dbp.log.Error(ansErr.Error())
		return 0, ansErr
	}

	dbp.log.Info("success push issue", "project", project)
	return issueId, nil
}

func (dbp *DbPusher) PushIssues(project *structures.DBProject, issues []datatransformer.DataTransformer) error {
	tx, err := dbp.db.Begin()
	if err != nil {
		ansErr := fmt.Errorf("%w: %w", myerr.ErrTranBegin, err)
		dbp.log.Error(ansErr.Error(), "project", project)
		return ansErr
	}

	for _, issue := range issues {
		issueId, err := dbp.PushIssue(project, &issue)
		if err != nil {
			ansErr := fmt.Errorf("%w - %s: %w", myerr.ErrPushIssue, project.Title, err)
			dbp.log.Error(ansErr.Error())
			tx.Rollback()
			return ansErr
		}

		if err := dbp.PushStatusChanges(issueId, &issue); err != nil {
			ansErr := fmt.Errorf("%w: %w", myerr.ErrInsertStatusChange, err)
			dbp.log.Error(ansErr.Error(), "project", project)
			tx.Rollback()
			return ansErr
		}
	}

	if err := tx.Commit(); err != nil {
		ansErr := fmt.Errorf("%w: %w", myerr.ErrTranClose, err)
		dbp.log.Error(ansErr.Error(), "project", project)
		return ansErr
	}

	dbp.log.Info("success save all issues", "project", project)
	return nil
}

func (dbp *DbPusher) getAuthorId(author *structures.DBAuthor) (int, error) {
	var authorId int
	var err error
	query := "SELECT id FROM author WHERE name=$1"

	_ = dbp.db.QueryRow(query, author.Name).Scan(&authorId)
	if authorId == 0 {
		authorId, err = dbp.PushAuthor(author)
		if err != nil {
			ansErr := fmt.Errorf("%w - %s: %w", myerr.ErrSelectAuthor, author.Name, err)
			dbp.log.Error(ansErr.Error())
			return 0, ansErr
		}
	}

	dbp.log.Info("success get author\assignee id", "author\assignee", author.Name)
	return authorId, nil
}

func (dbp *DbPusher) getProjectId(project *structures.DBProject) (int, error) {
	var projectId int
	var err error
	query := "SELECT id FROM projects WHERE title=$1"

	_ = dbp.db.QueryRow(query, project.Title).Scan(&projectId)
	if projectId == 0 {
		projectId, err = dbp.PushProject(project)
		if err != nil {
			ansErr := fmt.Errorf("%w - %s: %w", myerr.ErrSelectProject, project.Title, err)
			dbp.log.Error(ansErr.Error())
			return 0, ansErr
		}
	}

	dbp.log.Info("success ge priject ID", "project", project)
	return projectId, nil
}

func (dbp *DbPusher) hasStatusChange(issue int, time time.Time) bool {
	var count int
	query := "SELECT COUNT(*) FROM statuschanges WHERE issueId=$1 AND changeTime=$2"
	if err := dbp.db.QueryRow(query, issue, time).Scan(&count); err != nil {
		dbp.log.Error("err select status change", logger.Err(err))
		return false
	}

	dbp.log.Info("success check has status change")
	return count != 0
}

func buildConnectionstring(cfg *config.DBConfig) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
	)
}
