package dbpusher

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	configreader "github.com/jiraconnector/internal/configReader"
	datatransformer "github.com/jiraconnector/internal/dataTransformer"
	myerr "github.com/jiraconnector/internal/dbPusher/errors"
	"github.com/jiraconnector/internal/structures"
	_ "github.com/lib/pq"
)

type DbPusher struct {
	db *sql.DB
}

func NewDbPusher(cfg configreader.Config) (*DbPusher, error) {
	connStr := buildConnectionstring(&cfg.DBCfg)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		ansErr := fmt.Errorf("%w :: %w", myerr.ErrOpenDb, err)
		log.Println(ansErr)
		return nil, ansErr
	}

	return &DbPusher{
		db: db,
	}, nil
}

func (dbp *DbPusher) Close() {
	dbp.db.Close()
}

func (dbp *DbPusher) PushProject(project structures.DBProject) (int, error) {
	var projectId int
	query := "INSERT INTO project (title) VALUES ($1) ON CONFLICT (title) DO NOTHING RETURNING id"

	if err := dbp.db.QueryRow(query, project.Title).Scan(&projectId); err != nil {
		ansErr := fmt.Errorf("%w - %s :: %w", myerr.ErrInsertProject, project.Title, err)
		log.Println(ansErr)
		return 0, ansErr
	}

	return projectId, nil
}

func (dbp *DbPusher) PushProjects(projects []structures.DBProject) error {
	tx, err := dbp.db.Begin()
	if err != nil {
		ansErr := fmt.Errorf("%w :: %w", myerr.ErrTranBegin, err)
		log.Println(ansErr)
		return ansErr
	}

	for _, project := range projects {
		_, err = dbp.PushProject(project)
		if err != nil {
			ansErr := fmt.Errorf("%w - %s :: %w", myerr.ErrPushProject, project.Title, err)
			log.Println(ansErr)
			tx.Rollback()
			return ansErr
		}
	}

	if err := tx.Commit(); err != nil {
		ansErr := fmt.Errorf("%w :: %w", myerr.ErrTranClose, err)
		log.Println(ansErr)
		return ansErr
	}

	log.Println("All projects were saved")
	return nil

}

func (dbp *DbPusher) PushAuthor(author structures.DBAuthor) (int, error) {
	var authorId int
	query := "INSERT INTO author (name) VALUES ($1) ON CONFLICT (name) DO NOTHING RETURNING id"

	if err := dbp.db.QueryRow(query, author.Name).Scan(&authorId); err != nil {
		ansErr := fmt.Errorf("%w - %s :: %w", myerr.ErrInsertAuthor, author.Name, err)
		log.Println(ansErr)
		return 0, ansErr
	}

	return authorId, nil

}

func (dbp *DbPusher) PushStatusChanges(issue int, changes datatransformer.DataTransformer) error {
	query := "INSERT INTO statuschanges (issueId, authorId, changeTime, fromStatus, toStatus) VALUES ($1, $2, $3, $4, $5)"
	for author, statusChange := range changes.StatusChanges {
		if dbp.hasStatusChange(issue, statusChange.ChangeTime) {
			log.Println("already has such status change")
			return nil
		}
		authorId, err := dbp.getAuthorId(structures.DBAuthor{Name: author})
		if err != nil {
			return err
		}
		if _, err := dbp.db.Exec(query, issue, authorId, statusChange.ChangeTime, statusChange.FromStatus, statusChange.ToStatus); err != nil {
			return err
		}
	}
	return nil
}

func (dbp *DbPusher) PushIssue(project string, issue datatransformer.DataTransformer) (int, error) {
	projectId, err := dbp.getProjectId(project)
	if err != nil {
		return 0, err
	}

	authorId, err := dbp.getAuthorId(issue.Author)
	if err != nil {
		return 0, err
	}

	assegneeId, err := dbp.getAuthorId(issue.Assignee)
	if err != nil {
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

		ansErr := fmt.Errorf("%w - %s :: %w", myerr.ErrInsertIssue, project, err)
		log.Println(ansErr)
		return 0, ansErr
	}

	return issueId, nil
}

func (dbp *DbPusher) PushIssues(project string, issues []datatransformer.DataTransformer) error {
	tx, err := dbp.db.Begin()
	if err != nil {
		ansErr := fmt.Errorf("%w :: %w", myerr.ErrTranBegin, err)
		log.Println(ansErr)
		return ansErr
	}

	for _, issue := range issues {
		issueId, err := dbp.PushIssue(project, issue)
		if err != nil {
			ansErr := fmt.Errorf("%w - %s :: %w", myerr.ErrPushIssue, project, err)
			log.Println(ansErr)
			tx.Rollback()
			return ansErr
		}

		if err := dbp.PushStatusChanges(issueId, issue); err != nil {
			ansErr := fmt.Errorf("%w :: %w", myerr.ErrInsertStatusChange, err)
			log.Println(ansErr)
			tx.Rollback()
			return ansErr
		}
	}

	if err := tx.Commit(); err != nil {
		ansErr := fmt.Errorf("%w :: %w", myerr.ErrTranClose, err)
		log.Println(ansErr)
		return ansErr
	}

	log.Println("All issues were saved")
	return nil
}

func (dbp *DbPusher) getAuthorId(author structures.DBAuthor) (int, error) {
	var authorId int
	var err error
	query := "SELECT id FROM author WHERE name=$1"

	_ = dbp.db.QueryRow(query, author.Name).Scan(&authorId)
	if authorId == 0 {
		authorId, err = dbp.PushAuthor(author)
		if err != nil {
			ansErr := fmt.Errorf("%w - %s :: %w", myerr.ErrSelectAuthor, author.Name, err)
			log.Println(ansErr)
			return 0, ansErr
		}
	}

	return authorId, nil
}

func (dbp *DbPusher) getProjectId(project string) (int, error) {
	var projectId int
	var err error
	query := "SELECT id FROM project WHERE title=$1"

	_ = dbp.db.QueryRow(query, project).Scan(&projectId)
	if projectId == 0 {
		projectId, err = dbp.PushProject(structures.DBProject{Title: project})
		if err != nil {
			ansErr := fmt.Errorf("%w - %s :: %w", myerr.ErrSelectProject, project, err)
			log.Println(ansErr)
			return 0, ansErr
		}
	}

	return projectId, nil
}

func (dbp *DbPusher) hasStatusChange(issue int, time time.Time) bool {
	var count int
	query := "SELECT COUNT(*) FROM statuschanges WHERE issueId=$1 AND changeTime=$2"
	if err := dbp.db.QueryRow(query, issue, time).Scan(&count); err != nil {
		log.Printf("err select status change %v", err)
		return false
	}
	return count != 0
}

func buildConnectionstring(cfg *configreader.DBConfig) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
	)
}
