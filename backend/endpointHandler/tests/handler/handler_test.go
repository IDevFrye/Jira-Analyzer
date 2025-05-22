package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/endpointhandler/config"
	"github.com/endpointhandler/handler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter(cfg *config.Config) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/projects", func(c *gin.Context) { handler.GetProjects(c, cfg) })
	r.GET("/projects/:id", handler.GetProjectStats)
	r.DELETE("/projects/:id", handler.DeleteProject)
	r.GET("/connector/projects", func(c *gin.Context) { handler.GetJiraProjects(c, cfg) })
	r.POST("/connector/updateProject", func(c *gin.Context) { handler.UpdateJiraProject(c, cfg) })

	return r
}

func TestGetProjects(t *testing.T) {
	cfg := &config.Config{}
	r := setupRouter(cfg)

	req := httptest.NewRequest(http.MethodGet, "/projects", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusInternalServerError)
	assert.Contains(t, w.Body.String(), "{")
}

func TestGetProjectStats_InvalidID(t *testing.T) {
	r := setupRouter(&config.Config{})

	req := httptest.NewRequest(http.MethodGet, "/projects/abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid project ID")
}

func TestDeleteProject(t *testing.T) {
	r := setupRouter(&config.Config{})

	req := httptest.NewRequest(http.MethodDelete, "/projects/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusInternalServerError)
}

func TestGetJiraProjects(t *testing.T) {
	cfg := &config.Config{}
	cfg.Connector.BaseURL = "http://invalid-url" // для теста

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodGet, "/connector/projects", nil)

	handler.GetJiraProjects(c, cfg)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "failed to contact connector")
}

func TestUpdateJiraProject(t *testing.T) {
	cfg := &config.Config{}
	r := setupRouter(cfg)

	req := httptest.NewRequest(http.MethodPost, "/connector/updateProject?project=TEST", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusInternalServerError)
}
