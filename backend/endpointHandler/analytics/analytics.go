package analytics

import (
	"github.com/endpointhandler/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

func TimeOpenAnalytics(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project key is required"})
		return
	}

	var result []struct {
		Range string `json:"range"`
		Count int    `json:"count"`
	}

	err := repository.DB.Select(&result, `
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
			WHERE i.status NOT IN ('Closed', 'Resolved') AND p.key = $1
		) sub
		GROUP BY range
		ORDER BY MIN(age)
	`, key)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func StatusDistribution(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project key is required"})
		return
	}

	var result []struct {
		Status string `json:"status"`
		Count  int    `json:"count"`
	}

	err := repository.DB.Select(&result, `
		SELECT i.status, COUNT(*) AS count
		FROM Projects p
		JOIN Issue i ON p.id = i.projectId
		WHERE p.key = $1
		GROUP BY i.status
		ORDER BY i.status
	`, key)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func TimeSpentAnalytics(c *gin.Context) {
	projectKey := c.Query("key")
	if projectKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing project key"})
		return
	}

	var result []struct {
		Author         string `db:"author" json:"author"`
		TotalTimeSpent int    `db:"total_time_spent" json:"total_time_spent"`
	}

	err := repository.DB.Select(&result, `
		SELECT 
			a.name AS author,
			SUM(i.timeSpent) AS total_time_spent
		FROM Projects p
		JOIN Issue i ON p.id = i.projectId
		JOIN Author a ON a.id = i.authorId
		WHERE p.key = $1
		  AND i.timeSpent IS NOT NULL
		GROUP BY a.name
		ORDER BY total_time_spent DESC;
	`, projectKey)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func PriorityAnalytics(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project key is required"})
		return
	}

	var result []struct {
		Priority string `json:"priority"`
		Count    int    `json:"count"`
	}

	err := repository.DB.Select(&result, `
		SELECT i.priority, COUNT(*) AS count
		FROM Projects p
		JOIN Issue i ON p.id = i.projectId
		WHERE p.key = $1
		GROUP BY i.priority
		ORDER BY i.priority
	`, key)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
