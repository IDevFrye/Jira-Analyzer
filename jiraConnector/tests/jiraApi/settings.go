package jiraapiintegrations

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/jiraconnector/internal/structures"
)

func setupTestServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/rest/api/2/project":
			handleProjectsRequest(w, r)
		case strings.Contains(r.URL.Path, "/rest/api/2/search"):
			handleSearchRequest(w, r)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func handleProjectsRequest(w http.ResponseWriter, r *http.Request) {
	projects := []structures.JiraProject{
		{Key: "TEST1", Name: "Test Project 1"},
		{Key: "TEST2", Name: "Test Project 2"},
		{Key: "DEMO", Name: "Demo Project"},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

func handleSearchRequest(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	project := query.Get("jql")[8:] // extract project from "project=XXX"
	maxResults, _ := strconv.Atoi(query.Get("maxResults"))
	startAt, _ := strconv.Atoi(query.Get("startAt"))

	// Simulate pagination
	var issues []structures.JiraIssue
	if maxResults == 0 {
		// Total count request
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(structures.JiraIssues{Total: 150})
		return
	}

	for i := startAt; i < startAt+maxResults && i < 150; i++ {
		issues = append(issues, structures.JiraIssue{
			Key: fmt.Sprintf("%s-%d", project, i+1),
			Fields: structures.Field{
				Summary: fmt.Sprintf("Issue %d", i+1),
			},
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(structures.JiraIssues{
		Issues: issues,
		Total:  150,
	})
}
