package connector

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	handlerErr "github.com/jiraconnector/internal/apiJiraConnector/jiraHandlers/errors"
	configreader "github.com/jiraconnector/internal/configReader"
	myErr "github.com/jiraconnector/internal/connector/errors"
	"github.com/jiraconnector/internal/structures"
	"github.com/stretchr/testify/assert"
)

func mockConnectorWithURL(url string) *JiraConnector {
	cfg := configreader.Config{
		JiraCfg: configreader.JiraConfig{
			Url:           url,
			MinSleep:      int(10 * time.Millisecond),
			MaxSleep:      int(100 * time.Millisecond),
			ThreadCount:   2,
			IssueInOneReq: 1,
		},
	}
	return NewJiraConnector(cfg)
}

func TestGetAllProjects(t *testing.T) {
	tests := []struct {
		name           string
		handler        http.HandlerFunc
		expectErr      bool
		expectedLength int
	}{
		{
			name: "successful request",
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/rest/api/2/project", r.URL.Path)
				projects := []structures.JiraProject{
					{Id: "1", Name: "Test1"},
					{Id: "2", Name: "Test2"},
				}
				json.NewEncoder(w).Encode(projects)
			},
			expectErr:      false,
			expectedLength: 2,
		},
		{
			name: "inval Id JSON",
			handler: func(w http.ResponseWriter, r *http.Request) {
				io.WriteString(w, `invalid json`)
			},
			expectErr:      true,
			expectedLength: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			conn := mockConnectorWithURL(server.URL)
			projects, err := conn.GetAllProjects()

			assert.Equal(t, tt.expectErr, err != nil)
			assert.Len(t, projects, tt.expectedLength)
		})
	}
}

func TestGetProjectsPage(t *testing.T) {
	allProjects := []structures.JiraProject{
		{Id: "1", Name: "Alpha"},
		{Id: "2", Name: "Beta"},
		{Id: "3", Name: "Gamma"},
	}

	tests := []struct {
		name       string
		search     string
		page       int
		limit      int
		expectSize int
	}{
		{"no filter, page 1", "", 1, 2, 2},
		{"search Beta", "Beta", 1, 2, 1},
		{"out of range page", "", 10, 2, 0},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(allProjects)
	}))
	defer server.Close()

	conn := mockConnectorWithURL(server.URL)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := conn.GetProjectsPage(tt.search, tt.limit, tt.page)
			assert.NoError(t, err)
			assert.Len(t, result.Projects, tt.expectSize)
		})
	}
}

func TestRetryRequest(t *testing.T) {
	var attempt int32 = 0

	tests := []struct {
		name        string
		handler     http.HandlerFunc
		expectErr   bool
		expectRetry bool
	}{
		{
			name: "successful on first try",
			handler: func(w http.ResponseWriter, r *http.Request) {
				io.WriteString(w, "OK")
			},
			expectErr:   false,
			expectRetry: false,
		},
		{
			name: "too many retries",
			handler: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			},
			expectErr:   true,
			expectRetry: true,
		},
		{
			name: "404 returns immediately",
			handler: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Not Found", http.StatusNotFound)
			},
			expectErr:   true,
			expectRetry: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			atomic.StoreInt32(&attempt, 0)

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				atomic.AddInt32(&attempt, 1)
				tt.handler(w, r)
			}))
			defer server.Close()

			conn := mockConnectorWithURL(server.URL)
			resp, err := conn.retryRequest("GET", server.URL)

			if resp != nil {
				resp.Body.Close()
			}
			assert.Equal(t, tt.expectErr, err != nil)
		})
	}
}

func TestGetProjectIssues(t *testing.T) {
	type mockResponse struct {
		pathContains string
		response     string
	}

	// simulate issue retrievals
	responses := []mockResponse{
		{pathContains: "maxResults=0", response: `{"total": 2}`},
		{pathContains: "startAt=1", response: `{"issues":[{"id":"1","key":"ISSUE-1"}]}`},
		{pathContains: "startAt=2", response: `{"issues":[{"id":"2","key":"ISSUE-2"}]}`},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, resp := range responses {
			if strings.Contains(r.URL.String(), resp.pathContains) {
				io.WriteString(w, resp.response)
				return
			}
		}
		http.Error(w, "unexpected path", http.StatusNotFound)
	}))
	defer server.Close()

	conn := mockConnectorWithURL(server.URL)
	issues, err := conn.GetProjectIssues("TEST")
	assert.NoError(t, err)
	assert.Len(t, issues, 2)
	assert.Equal(t, "ISSUE-1", issues[0].Key)
}

func TestGetAllProjects_ErrorCases(t *testing.T) {
	tests := []struct {
		name      string
		handler   http.HandlerFunc
		expectErr error
	}{
		{
			name: "request error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Not Found", http.StatusNotFound)
			},
			expectErr: handlerErr.ErrNoProject,
		},
		{
			name: "read body error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Length", "1") // Force read error
			},
			expectErr: myErr.ErrReadResponseBody,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			conn := mockConnectorWithURL(server.URL)
			_, err := conn.GetAllProjects()
			assert.Error(t, err)
			if tt.expectErr != nil {
				assert.True(t, errors.Is(err, tt.expectErr), "expected error %v, got %v", tt.expectErr, err)
			}
		})
	}
}

func TestGetProjectsPage_ErrorCases(t *testing.T) {
	tests := []struct {
		name      string
		handler   http.HandlerFunc
		expectErr error
	}{
		{
			name: "error in GetAllProjects",
			handler: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Internal Error", http.StatusInternalServerError)
			},
			expectErr: myErr.ErrGetProjects,
		},
		{
			name: "invalid pagination",
			handler: func(w http.ResponseWriter, r *http.Request) {
				projects := []structures.JiraProject{{Id: "1", Name: "Test"}}
				json.NewEncoder(w).Encode(projects)
			},
			expectErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			conn := mockConnectorWithURL(server.URL)
			_, err := conn.GetProjectsPage("", -1, -1) // Invalid page/limit
			if tt.expectErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.expectErr))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetProjectIssues_ErrorCases(t *testing.T) {
	tests := []struct {
		name      string
		handler   http.HandlerFunc
		expectErr error
	}{
		{
			name: "error getting total issues",
			handler: func(w http.ResponseWriter, r *http.Request) {
				if strings.Contains(r.URL.String(), "maxResults=0") {
					http.Error(w, "Error", http.StatusInternalServerError)
					return
				}
				w.Write([]byte(`{"issues":[]}`))
			},
			expectErr: myErr.ErrGetIssues,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			conn := mockConnectorWithURL(server.URL)
			_, err := conn.GetProjectIssues("TEST")
			assert.Error(t, err)
			if tt.expectErr != nil {
				assert.True(t, errors.Is(err, tt.expectErr))
			}
		})
	}
}

func TestGetIssuesForOneThread_ErrorCases(t *testing.T) {
	tests := []struct {
		name      string
		handler   http.HandlerFunc
		expectErr error
	}{
		{
			name: "request error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Error", http.StatusInternalServerError)
			},
			expectErr: myErr.ErrGetIssues,
		},
		{
			name: "invalid JSON",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`invalid json`))
			},
			expectErr: myErr.ErrUnmarshalAns,
		},
		{
			name: "read body error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Length", "1000")
				w.Write([]byte("short body")) // Отправляем меньше данных чем объявлено
			},
			expectErr: myErr.ErrReadResponseBody,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			conn := mockConnectorWithURL(server.URL)
			_, err := conn.getIssuesForOneThread(0, "TEST")
			assert.Error(t, err)
			if tt.expectErr != nil {
				assert.True(t, errors.Is(err, tt.expectErr))
			}
		})
	}
}

func TestGetTotalIssues_ErrorCases(t *testing.T) {
	tests := []struct {
		name      string
		handler   http.HandlerFunc
		expectErr error
	}{
		{
			name: "request error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Error", http.StatusInternalServerError)
			},
			expectErr: myErr.ErrGetIssues,
		},
		{
			name: "invalid JSON",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`invalid json`))
			},
			expectErr: myErr.ErrUnmarshalAns,
		},
		{
			name: "read body error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Length", "1000")
				w.Write([]byte("short body")) // Отправляем меньше данных чем объявлено
			},
			expectErr: myErr.ErrReadResponseBody,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			conn := mockConnectorWithURL(server.URL)
			_, err := conn.getTotalIssues("TEST")
			assert.Error(t, err)
			if tt.expectErr != nil {
				assert.True(t, errors.Is(err, tt.expectErr))
			}
		})
	}
}

func TestRetryRequest_ErrorCases(t *testing.T) {
	tests := []struct {
		name      string
		handler   http.HandlerFunc
		expectErr error
	}{
		{
			name: "max retries exceeded",
			handler: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Error", http.StatusInternalServerError)
			},
			expectErr: myErr.ErrMaxTimeRequest,
		},
		{
			name: "bad request",
			handler: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Bad Request", http.StatusBadRequest)
			},
			expectErr: handlerErr.ErrNoProject,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			conn := mockConnectorWithURL(server.URL)
			_, err := conn.retryRequest("GET", server.URL)
			assert.Error(t, err)
			if tt.expectErr != nil {
				assert.True(t, errors.Is(err, tt.expectErr))
			}
		})
	}
}

func TestContainsSearchProject(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		substr   string
		expected bool
	}{
		{"exact match", "TestProject", "TestProject", true},
		{"case insensitive", "TestProject", "testproject", true},
		{"partial match", "TestProject", "Test", true},
		{"no match", "TestProject", "Other", false},
		{"empty substring", "TestProject", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsSearchProject(tt.str, tt.substr)
			assert.Equal(t, tt.expected, result)
		})
	}
}
