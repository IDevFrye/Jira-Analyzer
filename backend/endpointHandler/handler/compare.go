package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/endpointhandler/repository"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func parseProjectKeys(c *gin.Context) ([]string, error) {
	raw := c.Query("key")
	if raw == "" {
		return nil, fmt.Errorf("missing ?key=KEY1,KEY2,...")
	}
	parts := strings.Split(raw, ",")
	var keys []string
	for _, p := range parts {
		key := strings.TrimSpace(p)
		if key == "" {
			continue
		}
		keys = append(keys, key)
	}
	if len(keys) == 0 {
		return nil, fmt.Errorf("no valid project keys provided")
	}
	return keys, nil
}

type AgeRangeCount struct {
	Range string `db:"range" json:"range"`
	Count int    `db:"count" json:"count"`
}

func CompareTimeOpen(c *gin.Context) {
	keys, err := parseProjectKeys(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := make(map[string][]AgeRangeCount)

	for _, key := range keys {
		var ranges []AgeRangeCount

		query := `
			SELECT
				CASE
					WHEN age <= 1 THEN '0-1'
					WHEN age <= 2 THEN '1-2'
					WHEN age <= 3 THEN '2-3'
					WHEN age <= 5 THEN '3-5'
					WHEN age <= 7 THEN '5-7'
					WHEN age <= 10 THEN '7-10'
					WHEN age <= 14 THEN '10-14'
					WHEN age <= 21 THEN '14-21'
					WHEN age <= 30 THEN '21-30'
					ELSE '30+'
				END AS range,
				COUNT(*) AS count
			FROM (
				SELECT DATE_PART('day', NOW() - i.createdTime) AS age
				FROM Projects p
				JOIN Issue i ON p.id = i.projectId
				WHERE i.status NOT IN ('Closed', 'Resolved') AND p.title = $1
			) sub
			GROUP BY range
			ORDER BY MIN(age)
		`

		if err := repository.DB.Select(&ranges, query, key); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response[key] = ranges
	}

	c.JSON(http.StatusOK, response)
}

func CompareStatusDistribution(c *gin.Context) {
	keys, err := parseProjectKeys(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query, args, _ := sqlx.In(`
		SELECT 
			p.title AS project,
			i.status,
			COUNT(*) AS count
		FROM Projects p
		JOIN Issue i ON p.id = i.projectId
		WHERE p.title IN (?)
		GROUP BY p.title, i.status
		ORDER BY p.title, i.status
	`, keys)
	query = repository.DB.Rebind(query)

	var rows []struct {
		Project string `db:"project" json:"project"`
		Status  string `db:"status" json:"status"`
		Count   int    `db:"count" json:"count"`
	}
	if err := repository.DB.Select(&rows, query, args...); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Группируем статусы по проекту
	response := make(map[string]map[string]int)
	for _, r := range rows {
		if _, exists := response[r.Project]; !exists {
			response[r.Project] = make(map[string]int)
		}
		response[r.Project][r.Status] = r.Count
	}

	c.JSON(http.StatusOK, response)
}

func CompareTimeSpent(c *gin.Context) {
	keys, err := parseProjectKeys(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query, args, _ := sqlx.In(`
		SELECT 
			p.title AS project,
			a.name AS author,
			SUM(i.timeSpent) AS total_time_spent
		FROM Projects p
		JOIN Issue i ON p.id = i.projectId
		JOIN Author a ON a.id = i.authorId
		WHERE p.title IN (?) AND i.timeSpent IS NOT NULL
		GROUP BY p.title, a.name
		ORDER BY p.title, total_time_spent DESC
	`, keys)
	query = repository.DB.Rebind(query)

	var rows []struct {
		Project        string `db:"project" json:"project"`
		Author         string `db:"author" json:"author"`
		TotalTimeSpent int    `db:"total_time_spent" json:"total_time_spent"`
	}
	if err := repository.DB.Select(&rows, query, args...); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Группируем по проекту, внутри список авторов с их временем
	type authorStat struct {
		Author         string `json:"author"`
		TotalTimeSpent int    `json:"total_time_spent"`
	}
	response := make(map[string]struct {
		Authors []authorStat `json:"authors"`
	})
	for _, r := range rows {
		projectBlock := response[r.Project]
		projectBlock.Authors = append(projectBlock.Authors, authorStat{
			Author:         r.Author,
			TotalTimeSpent: r.TotalTimeSpent,
		})
		response[r.Project] = projectBlock
	}

	c.JSON(http.StatusOK, response)
}

func ComparePriority(c *gin.Context) {
	keys, err := parseProjectKeys(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query, args, _ := sqlx.In(`
		SELECT 
			p.title AS project,
			i.priority,
			COUNT(*) AS count
		FROM Projects p
		JOIN Issue i ON p.id = i.projectId
		WHERE p.title IN (?)
		GROUP BY p.title, i.priority
		ORDER BY p.title, i.priority
	`, keys)
	query = repository.DB.Rebind(query)

	var rows []struct {
		Project  string `db:"project" json:"project"`
		Priority string `db:"priority" json:"priority"`
		Count    int    `db:"count" json:"count"`
	}
	if err := repository.DB.Select(&rows, query, args...); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Группируем по проекту, внутри - приоритеты и их количество
	response := make(map[string]map[string]int)
	for _, r := range rows {
		if _, exists := response[r.Project]; !exists {
			response[r.Project] = make(map[string]int)
		}
		response[r.Project][r.Priority] = r.Count
	}

	c.JSON(http.StatusOK, response)
}
