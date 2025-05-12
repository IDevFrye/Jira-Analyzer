package jirahandlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	myErr "github.com/jiraconnector/internal/apiJiraConnector/jiraHandlers/errors"
	"github.com/jiraconnector/internal/structures"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Projects(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    map[string]string
		mockReturn     *structures.ResponseProject
		mockError      error
		expectedStatus int
		expectedError  error
	}{
		{
			name: "successful request with default params",
			queryParams: map[string]string{
				"limit":  "",
				"page":   "",
				"search": "",
			},
			mockReturn:     &structures.ResponseProject{},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedError:  nil,
		},
		{
			name: "successful request with custom params",
			queryParams: map[string]string{
				"limit":  "10",
				"page":   "2",
				"search": "test",
			},
			mockReturn:     &structures.ResponseProject{},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedError:  nil,
		},
		{
			name: "invalid limit parameter",
			queryParams: map[string]string{
				"limit": "invalid",
			},
			mockReturn:     nil,
			mockError:      nil,
			expectedStatus: myErr.GetStatusCode(myErr.ErrorsProject, myErr.ErrParamLimitPage),
			expectedError:  myErr.ErrParamLimitPage,
		},
		{
			name: "invalid page parameter",
			queryParams: map[string]string{
				"page": "invalid",
			},
			mockReturn:     nil,
			mockError:      nil,
			expectedStatus: myErr.GetStatusCode(myErr.ErrorsProject, myErr.ErrParamLimitPage),
			expectedError:  myErr.ErrParamLimitPage,
		},
		{
			name: "service error",
			queryParams: map[string]string{
				"search": "test",
				"limit":  "10",
				"page":   "1",
			},
			mockReturn:     nil,
			mockError:      errors.New("service error"),
			expectedStatus: myErr.GetStatusCode(myErr.ErrorsProject, myErr.ErrGetProjectPage),
			expectedError:  myErr.ErrGetProjectPage,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockJiraServiceInterface)

			limit, _ := strconv.Atoi(tt.queryParams["limit"])
			if limit == 0 {
				limit = 20
			}
			page, _ := strconv.Atoi(tt.queryParams["page"])
			if page == 0 {
				page = 1
			}
			search := tt.queryParams["search"]

			mockService.On("GetProjectsPage", search, limit, page).Return(tt.mockReturn, tt.mockError)

			router := mux.NewRouter()
			_ = NewHandler(mockService, router)

			req, err := http.NewRequest("GET", "/projects", nil)
			assert.NoError(t, err)

			q := req.URL.Query()
			for key, value := range tt.queryParams {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()

			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedError == nil && tt.mockError == nil {
				mockService.AssertExpectations(t)
			}
		})
	}
}

func TestHandler_UpdateProject(t *testing.T) {
	tests := []struct {
		name           string
		queryParam     string
		mockIssues     []structures.JiraIssue
		mockError      error
		pushError      error
		expectedStatus int
		expectedError  error
	}{
		{
			name:           "successful update",
			queryParam:     "TESTPROJ",
			mockIssues:     []structures.JiraIssue{},
			mockError:      nil,
			pushError:      nil,
			expectedStatus: http.StatusOK,
			expectedError:  nil,
		},
		{
			name:           "missing project parameter",
			queryParam:     "",
			mockIssues:     nil,
			mockError:      nil,
			pushError:      nil,
			expectedStatus: myErr.GetStatusCode(myErr.ErrorsUpdate, myErr.ErrParamProject),
			expectedError:  myErr.ErrParamProject,
		},
		{
			name:           "project not found",
			queryParam:     "UNKNOWN",
			mockIssues:     nil,
			mockError:      myErr.ErrNoProject,
			pushError:      nil,
			expectedStatus: myErr.GetStatusCode(myErr.ErrorsUpdate, myErr.ErrNoProject),
			expectedError:  myErr.ErrNoProject,
		},
		{
			name:           "update project error",
			queryParam:     "TESTPROJ",
			mockIssues:     nil,
			mockError:      errors.New("update error"),
			pushError:      nil,
			expectedStatus: myErr.GetStatusCode(myErr.ErrorsUpdate, myErr.ErrUpdProject),
			expectedError:  myErr.ErrUpdProject,
		},
		{
			name:           "push to db error",
			queryParam:     "TESTPROJ",
			mockIssues:     []structures.JiraIssue{},
			mockError:      nil,
			pushError:      errors.New("push error"),
			expectedStatus: myErr.GetStatusCode(myErr.ErrorsUpdate, myErr.ErrPushProject),
			expectedError:  myErr.ErrPushProject,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем mock сервиса
			mockService := new(MockJiraServiceInterface)

			// Настраиваем mock только если есть project параметр
			if tt.queryParam != "" {
				mockService.On("UpdateProjects", tt.queryParam).Return(tt.mockIssues, tt.mockError)

				// Если нет ошибки обновления, ожидаем вызов PushDataToDb
				if tt.mockError == nil {
					mockService.On("PushDataToDb", tt.queryParam, tt.mockIssues).Return(tt.pushError)
				}
			}

			// Создаем router и хендлер
			router := mux.NewRouter()
			_ = NewHandler(mockService, router)

			// Создаем запрос с query параметром
			req, err := http.NewRequest("POST", "/updateProject", nil)
			assert.NoError(t, err)

			q := req.URL.Query()
			if tt.queryParam != "" {
				q.Add("project", tt.queryParam)
			}
			req.URL.RawQuery = q.Encode()

			// Создаем ResponseRecorder для записи ответа
			rr := httptest.NewRecorder()

			// Выполняем запрос
			router.ServeHTTP(rr, req)

			// Проверяем статус код
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Проверяем вызовы mock
			mockService.AssertExpectations(t)
		})
	}
}

func TestGetProjectParams(t *testing.T) {
	tests := []struct {
		name        string
		queryParams map[string][]string // Изменено на []string
		expected    struct {
			limit  int
			page   int
			search string
			err    error
		}
	}{
		{
			name:        "default values",
			queryParams: map[string][]string{},
			expected: struct {
				limit  int
				page   int
				search string
				err    error
			}{
				limit:  20,
				page:   1,
				search: "",
				err:    nil,
			},
		},
		{
			name: "custom values",
			queryParams: map[string][]string{
				"limit":  {"50"},
				"page":   {"3"},
				"search": {"test"},
			},
			expected: struct {
				limit  int
				page   int
				search string
				err    error
			}{
				limit:  50,
				page:   3,
				search: "test",
				err:    nil,
			},
		},
		{
			name: "invalid limit",
			queryParams: map[string][]string{
				"limit": {"invalid"},
			},
			expected: struct {
				limit  int
				page   int
				search string
				err    error
			}{
				limit:  0,
				page:   0,
				search: "",
				err:    myErr.ErrParamLimitPage,
			},
		},
		{
			name: "invalid page",
			queryParams: map[string][]string{
				"page": {"invalid"},
			},
			expected: struct {
				limit  int
				page   int
				search string
				err    error
			}{
				limit:  0,
				page:   0,
				search: "",
				err:    myErr.ErrParamLimitPage,
			},
		},
		{
			name: "zero limit",
			queryParams: map[string][]string{
				"limit": {"0"},
			},
			expected: struct {
				limit  int
				page   int
				search string
				err    error
			}{
				limit:  0,
				page:   0,
				search: "",
				err:    myErr.ErrParamLimitPage,
			},
		},
		{
			name: "zero page",
			queryParams: map[string][]string{
				"page": {"0"},
			},
			expected: struct {
				limit  int
				page   int
				search string
				err    error
			}{
				limit:  0,
				page:   0,
				search: "",
				err:    myErr.ErrParamLimitPage,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем запрос с query параметрами
			req := &http.Request{
				URL: &url.URL{
					RawQuery: url.Values(tt.queryParams).Encode(),
				},
			}

			limit, page, search, err := getProjectParams(req)

			assert.Equal(t, tt.expected.limit, limit)
			assert.Equal(t, tt.expected.page, page)
			assert.Equal(t, tt.expected.search, search)
			assert.Equal(t, tt.expected.err, err)
		})
	}
}
