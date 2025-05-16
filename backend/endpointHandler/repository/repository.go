package repository

import (
	"github.com/endpointhandler/model"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"time"
)

var DB *sqlx.DB

func InitDB() {
	var err error
	DB, err = sqlx.Connect("postgres", "user=postgres password=00000 dbname=jira port=5432 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
}

func GetAllProjects() ([]model.Project, error) {
	var projects []model.Project
	err := DB.Select(&projects, "SELECT * FROM Projects")
	return projects, err
}

func GetStats(projectID int) (model.ProjectStats, error) {
	var stats model.ProjectStats
	err := DB.Get(&stats.TotalIssues, "SELECT COUNT(*) FROM Issue WHERE projectId=$1", projectID)
	if err != nil {
		return stats, err
	}
	DB.Get(&stats.OpenIssues, "SELECT COUNT(*) FROM Issue WHERE projectId=$1 AND status NOT IN ('Closed','Resolved')", projectID)
	DB.Get(&stats.ClosedIssues, "SELECT COUNT(*) FROM Issue WHERE projectId=$1 AND status='Closed'", projectID)
	DB.Get(&stats.ReopenedIssues, "SELECT COUNT(*) FROM StatusChanges WHERE issueId IN (SELECT id FROM Issue WHERE projectId=$1) AND toStatus='Reopened'", projectID)
	DB.Get(&stats.ResolvedIssues, "SELECT COUNT(*) FROM Issue WHERE projectId=$1 AND status='Resolved'", projectID)
	DB.Get(&stats.InProgressIssues, "SELECT COUNT(*) FROM Issue WHERE projectId=$1 AND status='In progress'", projectID)
	DB.Get(&stats.AvgResolutionTimeH, `SELECT COALESCE(AVG(EXTRACT(EPOCH FROM closedTime - createdTime))/3600, 0) FROM Issue WHERE projectId=$1 AND closedTime IS NOT NULL`, projectID)
	DB.Get(&stats.AvgCreatedPerDay7d, `SELECT COUNT(*)/7.0 FROM Issue WHERE projectId=$1 AND createdTime > $2`, projectID, time.Now().AddDate(0, 0, -7))
	return stats, nil
}

func DeleteProject(projectID int) error {
	_, err := DB.Exec("DELETE FROM Projects WHERE id=$1", projectID)
	return err
}
