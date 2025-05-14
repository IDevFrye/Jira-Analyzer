package jirahandlers

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	myErr "github.com/jiraconnector/internal/apiJiraConnector/jiraHandlers/errors"
	"github.com/jiraconnector/internal/apiJiraConnector/jiraHandlers/responseutils"
	datatransformer "github.com/jiraconnector/internal/dataTransformer"
	"github.com/jiraconnector/internal/structures"
	"github.com/jiraconnector/pkg/middleware"
)

//go:generate mockery

type JiraServiceInterface interface {
	GetProjectsPage(search string, limit, page int) (*structures.ResponseProject, error)
	UpdateProjects(projectId string) ([]structures.JiraIssue, error)

	PushDataToDb(project string, issues []structures.JiraIssue) error
	TransformDataToDb(project string, issues []structures.JiraIssue) []datatransformer.DataTransformer
}

type handler struct {
	service JiraServiceInterface
	log     *slog.Logger
}

func NewHandler(service JiraServiceInterface, router *mux.Router, log *slog.Logger) *mux.Router {
	h := handler{service: service, log: log}

	router.Use(middleware.NewLoggerMiddleware(log))

	router.HandleFunc("/api/v1/connector/projects", h.projects).Methods(http.MethodOptions, http.MethodGet)
	router.HandleFunc("/api/v1/connector/updateProject", h.updateProject).Methods(http.MethodOptions, http.MethodPost)
	log.Info("create router")
	return router
}

// @Summary Get paginated list of Jira projects
// @Description Получение проектов с пагинацией
// @Tags projects
// @Accept  json
// @Produce  json
// @Param   limit  query  int     false  "Items per page"
// @Param   page   query  int     false  "Page number"
// @Param   search query  string  false  "Search filter"
// @Success 200 {object} structures.ResponseProject
// @Failure 400 {object} responseutils.ErrorResponse
// @Failure 500 {object} responseutils.ErrorResponse
// @Router /api/v1/connector/projects [get]
func (h *handler) projects(w http.ResponseWriter, r *http.Request) {
	limit, page, search, err := getProjectParams(r)
	if err != nil {
		responseutils.WriteError(w, h.log, myErr.GetStatusCode(myErr.ErrorsProject, myErr.ErrParamLimitPage), myErr.ErrParamLimitPage.Error(), err)
		return
	}

	projects, err := h.service.GetProjectsPage(search, limit, page)
	if err != nil {
		responseutils.WriteError(w, h.log, myErr.GetStatusCode(myErr.ErrorsProject, myErr.ErrGetProjectPage), myErr.ErrGetProjectPage.Error(), err)
		return
	}

	responseutils.WriteSuccess(w, h.log, http.StatusOK, projects)
	h.log.Info("Got project page", "page", page)
}

// @Summary Update Jira project and push issues to DB
// @Description Обновляет проект в Jira, загружает задачи и сохраняет их в базу данных
// @Tags projects
// @Accept  json
// @Produce  json
// @Param   project  query  string  true  "Project Key or ID (required)"
// @Success 200 {object} structures.ResponseUpdate
// @Failure 400 {object} responseutils.ErrorResponse
// @Failure 404 {object} responseutils.ErrorResponse
// @Failure 500 {object} responseutils.ErrorResponse
// @Router /api/v1/connector/updateProject [post]
func (h *handler) updateProject(w http.ResponseWriter, r *http.Request) {
	project := r.URL.Query().Get("project")
	if project == "" {
		responseutils.WriteError(w, h.log, myErr.GetStatusCode(myErr.ErrorsUpdate, myErr.ErrParamProject), myErr.ErrParamProject.Error(), nil)
		return
	}

	issues, err := h.service.UpdateProjects(project)
	if err != nil {
		if errors.Is(err, myErr.ErrNoProject) {
			responseutils.WriteError(w, h.log, myErr.GetStatusCode(myErr.ErrorsUpdate, myErr.ErrNoProject), myErr.ErrNoProject.Error(), err)
		} else {
			responseutils.WriteError(w, h.log, myErr.GetStatusCode(myErr.ErrorsUpdate, myErr.ErrUpdProject), myErr.ErrUpdProject.Error(), err)
		}
		return
	}

	if err := h.service.PushDataToDb(project, issues); err != nil {
		responseutils.WriteError(w, h.log, myErr.GetStatusCode(myErr.ErrorsUpdate, myErr.ErrPushProject), myErr.ErrPushProject.Error(), err)
		return
	}

	responseutils.WriteSuccess(w, h.log, http.StatusOK, structures.ResponseUpdate{Project: project, Status: "updated"})
	h.log.Info("Update issues", "project", project)
}

func getProjectParams(r *http.Request) (int, int, string, error) {
	var err error
	limit := 20
	page := 1
	search := ""

	if r.URL.Query().Get("limit") != "" {
		limit, err = strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil || limit <= 0 {
			return 0, 0, "", myErr.ErrParamLimitPage
		}
	}

	if r.URL.Query().Get("page") != "" {
		page, err = strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil || page <= 0 {
			return 0, 0, "", myErr.ErrParamLimitPage
		}
	}

	if r.URL.Query().Get("search") != "" {
		search = r.URL.Query().Get("search")
	}

	return limit, page, search, nil
}
