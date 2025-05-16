package handler

import (
	"github.com/endpointhandler/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

func TimeOpenAnalytics(c *gin.Context) {
	projectKey := c.Query("key")
	if projectKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing project key"})
		return
	}

	var result []struct {
		Range string `db:"range" json:"range"`
		Count int    `db:"count" json:"count"`
	}

	err := repository.DB.Select(&result, `
		SELECT CASE
			WHEN now() - i.created <= interval '1 day' THEN '0-1'
			WHEN now() - i.created <= interval '2 days' THEN '1-2'
			WHEN now() - i.created <= interval '3 days' THEN '2-3'
			WHEN now() - i.created <= interval '5 days' THEN '3-5'
			WHEN now() - i.created <= interval '7 days' THEN '5-7'
			WHEN now() - i.created <= interval '10 days' THEN '7-10'
			WHEN now() - i.created <= interval '14 days' THEN '10-14'
			WHEN now() - i.created <= interval '21 days' THEN '14-21'
			WHEN now() - i.created <= interval '30 days' THEN '21-30'
			ELSE '30+'
		END AS range,
		COUNT(*) AS count
		FROM Issue i
		JOIN Projects p ON i.projectId = p.id
		WHERE p.key = $1 AND i.status NOT IN ('Closed', 'Resolved')
		GROUP BY range
		ORDER BY MIN(i.created)
	`, projectKey)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func StatusDistribution(c *gin.Context) {
	projectKey := c.Query("key")
	if projectKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing project key"})
		return
	}

	var result []struct {
		Status string `db:"status" json:"status"`
		Count  int    `db:"count" json:"count"`
	}

	err := repository.DB.Select(&result, `
		SELECT i.status, COUNT(*) AS count
		FROM Issue i
		JOIN Projects p ON i.projectId = p.id
		WHERE p.key = $1
		GROUP BY i.status
		ORDER BY i.status
	`, projectKey)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func TimeSpentAnalytics(c *gin.Context) {
	projectKey := c.Query("key")
	if projectKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing project key"})
		return
	}

	var result struct {
		TotalTime float64 `db:"total_time" json:"total_time"`
	}

	err := repository.DB.Get(&result, `
		SELECT SUM(i.timeSpent) AS total_time
		FROM Issue i
		JOIN Projects p ON i.projectId = p.id
		WHERE p.key = $1 AND i.timeSpent IS NOT NULL
	`, projectKey)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func PriorityAnalytics(c *gin.Context) {
	projectKey := c.Query("key")
	if projectKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing project key"})
		return
	}

	var result []struct {
		Priority string `db:"priority" json:"priority"`
		Count    int    `db:"count" json:"count"`
	}

	err := repository.DB.Select(&result, `
		SELECT i.priority, COUNT(*) AS count
		FROM Issue i
		JOIN Projects p ON i.projectId = p.id
		WHERE p.key = $1
		GROUP BY i.priority
		ORDER BY i.priority
	`, projectKey)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
