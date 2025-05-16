package handler

import (
	"fmt"
	"github.com/endpointhandler/repository"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strconv"
	"strings"
)

func parseProjectIDs(c *gin.Context) ([]int, error) {
	raw := c.Query("ids")
	if raw == "" {
		return nil, fmt.Errorf("missing ?ids=1,2,...")
	}
	parts := strings.Split(raw, ",")
	var ids []int
	for _, p := range parts {
		id, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil {
			return nil, fmt.Errorf("invalid id: %s", p)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func CompareTimeOpen(c *gin.Context) {
	ids, err := parseProjectIDs(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query, args, _ := sqlx.In(`
		SELECT p.id AS project_id, p.title, COUNT(*) AS open_issues
		FROM Projects p
		JOIN Issue i ON p.id = i.projectId
		WHERE p.id IN (?) AND i.status NOT IN ('Closed', 'Resolved')
		GROUP BY p.id, p.title`, ids)
	query = repository.DB.Rebind(query)

	var result []struct {
		ProjectID  int    `db:"project_id" json:"project_id"`
		Title      string `db:"title" json:"title"`
		OpenIssues int    `db:"open_issues" json:"open_issues"`
	}
	err = repository.DB.Select(&result, query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func CompareStatusDistribution(c *gin.Context) {
	ids, err := parseProjectIDs(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query, args, _ := sqlx.In(`
		SELECT p.id AS project_id, p.title, i.status, COUNT(*) AS count
		FROM Projects p
		JOIN Issue i ON p.id = i.projectId
		WHERE p.id IN (?)
		GROUP BY p.id, p.title, i.status
		ORDER BY p.id, i.status`, ids)
	query = repository.DB.Rebind(query)

	var result []struct {
		ProjectID int    `db:"project_id" json:"project_id"`
		Title     string `db:"title" json:"title"`
		Status    string `db:"status" json:"status"`
		Count     int    `db:"count" json:"count"`
	}
	err = repository.DB.Select(&result, query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func CompareTimeSpent(c *gin.Context) {
	ids, err := parseProjectIDs(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query, args, _ := sqlx.In(`
		SELECT p.id AS project_id, p.title, AVG(i.timeSpent) AS avg_time_spent
		FROM Projects p
		JOIN Issue i ON p.id = i.projectId
		WHERE p.id IN (?) AND i.timeSpent IS NOT NULL
		GROUP BY p.id, p.title`, ids)
	query = repository.DB.Rebind(query)

	var result []struct {
		ProjectID    int     `db:"project_id" json:"project_id"`
		Title        string  `db:"title" json:"title"`
		AvgTimeSpent float64 `db:"avg_time_spent" json:"avg_time_spent"`
	}
	err = repository.DB.Select(&result, query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func ComparePriority(c *gin.Context) {
	ids, err := parseProjectIDs(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query, args, _ := sqlx.In(`
		SELECT p.id AS project_id, p.title, i.priority, COUNT(*) AS count
		FROM Projects p
		JOIN Issue i ON p.id = i.projectId
		WHERE p.id IN (?)
		GROUP BY p.id, p.title, i.priority
		ORDER BY p.id, i.priority`, ids)
	query = repository.DB.Rebind(query)

	var result []struct {
		ProjectID int    `db:"project_id" json:"project_id"`
		Title     string `db:"title" json:"title"`
		Priority  string `db:"priority" json:"priority"`
		Count     int    `db:"count" json:"count"`
	}
	err = repository.DB.Select(&result, query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
