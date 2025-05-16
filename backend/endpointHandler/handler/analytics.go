package handler

import (
	"github.com/endpointhandler/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

func TimeOpenAnalytics(c *gin.Context) {
	var result []struct {
		ProjectID  int    `db:"project_id" json:"project_id"`
		Title      string `db:"title" json:"title"`
		OpenIssues int    `db:"open_issues" json:"open_issues"`
	}
	err := repository.DB.Select(&result, `
		SELECT p.id AS project_id, p.title, COUNT(*) AS open_issues
		FROM Projects p
		JOIN Issue i ON p.id = i.projectId
		WHERE i.status NOT IN ('Closed', 'Resolved')
		GROUP BY p.id, p.title`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func StatusDistribution(c *gin.Context) {
	var result []struct {
		ProjectID int    `db:"project_id" json:"project_id"`
		Title     string `db:"title" json:"title"`
		Status    string `db:"status" json:"status"`
		Count     int    `db:"count" json:"count"`
	}
	err := repository.DB.Select(&result, `
		SELECT p.id AS project_id, p.title, i.status, COUNT(*) AS count
		FROM Projects p
		JOIN Issue i ON p.id = i.projectId
		GROUP BY p.id, p.title, i.status
		ORDER BY p.id, i.status`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func TimeSpentAnalytics(c *gin.Context) {
	var result []struct {
		ProjectID    int     `db:"project_id" json:"project_id"`
		Title        string  `db:"title" json:"title"`
		AvgTimeSpent float64 `db:"avg_time_spent" json:"avg_time_spent"`
	}
	err := repository.DB.Select(&result, `
		SELECT p.id AS project_id, p.title, AVG(i.timeSpent) AS avg_time_spent
		FROM Projects p
		JOIN Issue i ON p.id = i.projectId
		WHERE i.timeSpent IS NOT NULL
		GROUP BY p.id, p.title`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func PriorityAnalytics(c *gin.Context) {
	var result []struct {
		ProjectID int    `db:"project_id" json:"project_id"`
		Title     string `db:"title" json:"title"`
		Priority  string `db:"priority" json:"priority"`
		Count     int    `db:"count" json:"count"`
	}
	err := repository.DB.Select(&result, `
		SELECT p.id AS project_id, p.title, i.priority, COUNT(*) AS count
		FROM Projects p
		JOIN Issue i ON p.id = i.projectId
		GROUP BY p.id, p.title, i.priority
		ORDER BY p.id, i.priority`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
