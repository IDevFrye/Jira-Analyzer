package handler

import (
	"fmt"
	"github.com/endpointhandler/config"
	"net/http"
	"net/url"
	"strconv"

	"github.com/endpointhandler/service"
	"github.com/gin-gonic/gin"
)

func GetProjects(c *gin.Context, cfg *config.Config) {
	projectsResp, err := service.FetchAndStoreProjects(cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, projectsResp)
}

func GetProjectStats(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project ID"})
		return
	}

	stats, err := service.GetProjectStats(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get project stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func DeleteProject(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := service.DeleteProject(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func GetJiraProjects(c *gin.Context, cfg *config.Config) {
	limit := c.DefaultQuery("limit", "20")
	page := c.DefaultQuery("page", "1")
	search := c.DefaultQuery("search", "")

	reqURL := fmt.Sprintf("%s/projects?limit=%s&page=%s&search=%s",
		cfg.Connector.BaseURL,
		url.QueryEscape(limit),
		url.QueryEscape(page),
		url.QueryEscape(search),
	)

	resp, err := http.Get(reqURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to contact connector"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{"error": "connector returned error"})
		return
	}

	c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
}

func UpdateJiraProject(c *gin.Context, cfg *config.Config) {
	key := c.Query("project")
	result, err := service.UpdateJiraProject(cfg, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
