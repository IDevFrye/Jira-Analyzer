package jirahandlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	myErr "github.com/jiraconnector/internal/apiJiraConnector/jiraHandlers/errors"
	datatransformer "github.com/jiraconnector/internal/dataTransformer"
	"github.com/jiraconnector/internal/structures"
	"github.com/jiraconnector/pkg/logger"
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

	router.HandleFunc("/projects", h.projects).Methods(http.MethodOptions, http.MethodGet)
	router.HandleFunc("/updateProject", h.updateProject).Methods(http.MethodOptions, http.MethodPost)
	log.Info("create router")
	return router
}

func (h *handler) projects(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	limit, page, search, err := getProjectParams(r)
	if err != nil {
		h.log.Error(myErr.ErrParamLimitPage.Error(), logger.Err(err))
		http.Error(w, myErr.ErrParamLimitPage.Error(), myErr.GetStatusCode(myErr.ErrorsProject, myErr.ErrParamLimitPage))
		return
	}

	projects, err := h.service.GetProjectsPage(search, limit, page)
	if err != nil {
		h.log.Error(myErr.ErrGetProjectPage.Error(), logger.Err(err))
		http.Error(w, myErr.ErrGetProjectPage.Error(), myErr.GetStatusCode(myErr.ErrorsProject, myErr.ErrGetProjectPage))
		return
	}

	if err = json.NewEncoder(w).Encode(projects); err != nil {
		h.log.Error(myErr.ErrEncodeAns.Error(), logger.Err(err))
		http.Error(w, myErr.ErrEncodeAns.Error(), myErr.GetStatusCode(myErr.ErrorsProject, myErr.ErrEncodeAns))
		return
	}

	h.log.Info("Got project page: %d", "page", page)
	w.WriteHeader(http.StatusOK)

}

func (h *handler) updateProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	project := r.URL.Query().Get("project")
	if project == "" {
		h.log.Error(myErr.ErrParamProject.Error())
		http.Error(w, myErr.ErrParamProject.Error(), myErr.GetStatusCode(myErr.ErrorsUpdate, myErr.ErrParamProject))
		return
	}

	issues, err := h.service.UpdateProjects(project)
	if err != nil {
		if errors.Is(err, myErr.ErrNoProject) {
			h.log.Error(myErr.ErrNoProject.Error(), "project", project, logger.Err(err))
			http.Error(w, myErr.ErrNoProject.Error(), myErr.GetStatusCode(myErr.ErrorsUpdate, myErr.ErrNoProject))
		} else {
			h.log.Error(myErr.ErrUpdProject.Error(), "project", project, logger.Err(err))
			http.Error(w, myErr.ErrUpdProject.Error(), myErr.GetStatusCode(myErr.ErrorsUpdate, myErr.ErrUpdProject))
		}
		return
	}

	if err := h.service.PushDataToDb(project, issues); err != nil {
		h.log.Error(myErr.ErrPushProject.Error(), "project", project, logger.Err(err))
		http.Error(w, myErr.ErrPushProject.Error(), myErr.GetStatusCode(myErr.ErrorsUpdate, myErr.ErrPushProject))
		return
	}

	if err = json.NewEncoder(w).Encode(map[string]string{project: "updated"}); err != nil {
		h.log.Error(myErr.ErrEncodeAns.Error(), logger.Err(err))
		http.Error(w, myErr.ErrEncodeAns.Error(), myErr.GetStatusCode(myErr.ErrorsProject, myErr.ErrEncodeAns))
		return
	}

	h.log.Info("Update issues")
	w.WriteHeader(http.StatusOK)
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
