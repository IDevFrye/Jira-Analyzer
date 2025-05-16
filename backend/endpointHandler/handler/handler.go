package handler

import (
	"net/http"
	"strconv"

	"github.com/endpointhandler/service"
	"github.com/gin-gonic/gin"
)

func GetProjects(c *gin.Context) {
	projects, err := service.GetAllProjects()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, projects)
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

func GetJiraProjects(c *gin.Context) {
	projects, err := service.FetchJiraProjects()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, projects)
}

func UpdateJiraProject(c *gin.Context) {
	key := c.Query("project")
	result, err := service.UpdateJiraProject(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func CompareTask(c *gin.Context) {
	// mock example
	task := c.Param("taskNumber")
	c.JSON(http.StatusOK, gin.H{"task": task, "result": "comparison placeholder"})
}

func CompareAllProjects(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"result": "compare all projects placeholder"})
}
