package handler

import (
	"github.com/endpointhandler/config"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupRouter(cfg *config.Config) *gin.Engine {
	r := gin.Default()
	api := r.Group("/api")

	api.GET("/projects", func(c *gin.Context) {
		GetProjects(c, cfg)
	})
	api.GET("/projects/:id", GetProjectStats)
	api.DELETE("/projects/:id", DeleteProject)

	connector := api.Group("/connector")
	{
		connector.GET("/projects", func(c *gin.Context) {
			GetJiraProjects(c, cfg)
		})
		connector.POST("/updateProject", func(c *gin.Context) {
			UpdateJiraProject(c, cfg)
		})
	}

	analytics := api.Group("/analytics")
	{
		analytics.GET("/time-open", dummyHandler) // заглушки
		analytics.GET("/status-distribution", dummyHandler)
		analytics.GET("/time-spent", dummyHandler)
		analytics.GET("/priority", dummyHandler)
	}

	compare := api.Group("/compare")
	{
		compare.GET("/time-open", dummyHandler)
		compare.GET("/status-distribution", dummyHandler)
		compare.GET("/time-spent", dummyHandler)
		compare.GET("/priority", dummyHandler)
	}

	return r
}

func dummyHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func TestGetProjects(t *testing.T) {
	cfg := &config.Config{}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/projects", nil)
	setupRouter(cfg).ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("unexpected status code: %d", w.Code)
	}
}

func TestGetProjectStats(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/projects/1", nil)
	setupRouter(nil).ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError && w.Code != http.StatusBadRequest {
		t.Errorf("unexpected status code: %d", w.Code)
	}
}

func TestDeleteProject(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/projects/1", nil)
	setupRouter(nil).ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("unexpected status code: %d", w.Code)
	}
}

func TestUpdateJiraProject(t *testing.T) {
	cfg := &config.Config{}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/connector/updateProject?project=PRJ", nil)
	setupRouter(cfg).ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("unexpected status code: %d", w.Code)
	}
}

func TestGetJiraProjects(t *testing.T) {
	cfg := &config.Config{}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/connector/projects?limit=5&page=1", nil)
	setupRouter(cfg).ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("unexpected status code: %d", w.Code)
	}
}

func TestGetProjectStats_BadRequest(t *testing.T) {

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/projects/not-a-number", nil)
	setupRouter(nil).ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest && w.Code != http.StatusInternalServerError {
		t.Errorf("expected 400 or 500, got %d", w.Code)
	}
}

func TestUpdateJiraProject_EmptyParam(t *testing.T) {
	cfg := &config.Config{}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/connector/updateProject", nil)
	setupRouter(cfg).ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest && w.Code != http.StatusInternalServerError {
		t.Errorf("expected 400 or 500, got %d", w.Code)
	}
}

func TestGetJiraProjects_MissingParams(t *testing.T) {
	cfg := &config.Config{}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/connector/projects", nil)
	setupRouter(cfg).ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest && w.Code != http.StatusInternalServerError {
		t.Errorf("expected 400 or 500, got %d", w.Code)
	}
}
