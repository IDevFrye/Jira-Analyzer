package repository

import (
	"fmt"
	"github.com/endpointhandler/model"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"strconv"
	"strings"
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

func GetFilteredProjects(limit, offset int, search string) ([]model.UIProject, int, error) {
	var projects []model.UIProject
	var total int

	search = strings.ToLower(search)
	searchQuery := ""
	args := []interface{}{}

	if search != "" {
		searchQuery = "WHERE LOWER(title) LIKE $1"
		args = append(args, "%"+search+"%")
	}

	paramLimit := len(args) + 1
	paramOffset := len(args) + 2

	query := fmt.Sprintf(`
        SELECT id, title AS key, title AS name, '' AS self, TRUE as existence
        FROM Projects
        %s
        ORDER BY title
        LIMIT $%d OFFSET $%d
    `, searchQuery, paramLimit, paramOffset)

	args = append(args, limit, offset)

	err := DB.Select(&projects, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("select error: %w", err)
	}

	countQuery := "SELECT COUNT(*) FROM Projects"
	if search != "" {
		countQuery += " " + searchQuery
	}

	err = DB.Get(&total, countQuery, args[:len(args)-2]...)
	if err != nil {
		return nil, 0, fmt.Errorf("count error: %w", err)
	}

	return projects, total, nil
}

func GetAllProjects() ([]model.Project, error) {
	var dbProjects []model.DBProject
	err := DB.Select(&dbProjects, "SELECT id, title FROM Projects")
	if err != nil {
		return nil, err
	}

	// Преобразование в []Project
	var result []model.Project
	for _, p := range dbProjects {
		result = append(result, model.Project{
			ID:   strconv.Itoa(p.ID), // Преобразование int → string
			Key:  p.Title,
			Name: p.Title,
			Self: fmt.Sprintf("https://issues.apache.org/jira/rest/api/2/project/%d", p.ID),
		})
	}
	return result, nil
}

func GetStats(projectID int) (model.ProjectStats, error) {
	var stats model.ProjectStats

	err := DB.Get(&stats.TotalIssues, "SELECT COUNT(*) FROM Issue WHERE projectId=$1", projectID)
	if err != nil {
		return stats, err
	}

	err = DB.Get(&stats.OpenIssues, "SELECT COUNT(*) FROM Issue WHERE projectId=$1 AND status NOT IN ('Closed','Resolved')", projectID)
	if err != nil {
		return stats, err
	}

	err = DB.Get(&stats.ClosedIssues, "SELECT COUNT(*) FROM Issue WHERE projectId=$1 AND status='Closed'", projectID)
	if err != nil {
		return stats, err
	}

	err = DB.Get(&stats.ReopenedIssues, `
		SELECT COUNT(*) FROM StatusChanges 
		WHERE issueId IN (SELECT id FROM Issue WHERE projectId=$1) AND toStatus='Reopened'`, projectID)
	if err != nil {
		return stats, err
	}

	err = DB.Get(&stats.ResolvedIssues, "SELECT COUNT(*) FROM Issue WHERE projectId=$1 AND status='Resolved'", projectID)
	if err != nil {
		return stats, err
	}

	err = DB.Get(&stats.InProgressIssues, "SELECT COUNT(*) FROM Issue WHERE projectId=$1 AND status='In progress'", projectID)
	if err != nil {
		return stats, err
	}

	err = DB.Get(&stats.AvgResolutionTimeH, `
		SELECT COALESCE(AVG(EXTRACT(EPOCH FROM closedTime - createdTime))/3600, 0) 
		FROM Issue WHERE projectId=$1 AND closedTime IS NOT NULL`, projectID)
	if err != nil {
		return stats, err
	}

	err = DB.Get(&stats.AvgCreatedPerDay7d, `
		SELECT COUNT(*) / 7.0 
		FROM Issue WHERE projectId=$1 AND createdTime > $2`, projectID, time.Now().AddDate(0, 0, -7))
	if err != nil {
		return stats, err
	}

	return stats, nil
}

func DeleteProject(projectID int) error {
	_, err := DB.Exec("DELETE FROM Projects WHERE id=$1", projectID)
	return err
}
