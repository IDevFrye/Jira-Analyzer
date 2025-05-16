package handler

import (
	"fmt"
	"net/http"
	"net/url"
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
	// Получаем query параметры
	limit := c.DefaultQuery("limit", "20")
	page := c.DefaultQuery("page", "1")
	search := c.DefaultQuery("search", "")

	// Формируем URL запроса к коннектору
	baseURL := "http://localhost:8080/api/v1/connector/projects"
	reqURL := fmt.Sprintf("%s?limit=%s&page=%s&search=%s",
		baseURL,
		url.QueryEscape(limit),
		url.QueryEscape(page),
		url.QueryEscape(search),
	)

	// Выполняем запрос к коннектору
	resp, err := http.Get(reqURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to contact connector"})
		return
	}
	defer resp.Body.Close()

	// Проверка статуса
	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{"error": "connector returned error"})
		return
	}

	// Возвращаем тело ответа клиенту как есть
	c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
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
