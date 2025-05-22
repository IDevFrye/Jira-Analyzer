package handler_test

import (
	"encoding/json"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/endpointhandler/handler"
	"github.com/endpointhandler/repository"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"net/http"
	"net/http/httptest"
	"testing"
)

// helper: запуск маршрута с нужным эндпоинтом и запросом с ключом
func setupRouterWithPath(path string, handlerFunc gin.HandlerFunc) (*gin.Engine, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET(path, handlerFunc)
	w := httptest.NewRecorder()
	return r, w
}

func TestTimeOpenAnalytics_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repository.DB = sqlx.NewDb(db, "postgres")

	mock.ExpectQuery("SELECT(.*)FROM \\(").
		WithArgs("PROJ").
		WillReturnRows(sqlmock.NewRows([]string{"range", "count"}).
			AddRow("0-1", 5).
			AddRow("1-2", 3),
		)

	r, w := setupRouterWithPath("/api/v1/analytics/time-open", handler.TimeOpenAnalytics)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/time-open?key=PROJ", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp []map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(resp) != 2 {
		t.Errorf("unexpected response length: %d", len(resp))
	}
}

func TestTimeOpenAnalytics_MissingKey(t *testing.T) {
	r, w := setupRouterWithPath("/api/v1/analytics/time-open", handler.TimeOpenAnalytics)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/time-open", nil) // без key
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestStatusDistribution_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repository.DB = sqlx.NewDb(db, "postgres")

	mock.ExpectQuery("SELECT i.status, COUNT").
		WithArgs("PROJ").
		WillReturnRows(sqlmock.NewRows([]string{"status", "count"}).
			AddRow("Open", 10).
			AddRow("In Progress", 5),
		)

	r, w := setupRouterWithPath("/api/v1/analytics/status-distribution", handler.StatusDistribution)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/status-distribution?key=PROJ", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestStatusDistribution_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repository.DB = sqlx.NewDb(db, "postgres")

	mock.ExpectQuery("SELECT i.status, COUNT").
		WithArgs("PROJ").
		WillReturnError(errors.New("db error"))

	r, w := setupRouterWithPath("/api/v1/analytics/status-distribution", handler.StatusDistribution)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/status-distribution?key=PROJ", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestTimeSpentAnalytics_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repository.DB = sqlx.NewDb(db, "postgres")

	mock.ExpectQuery("SELECT(.*)SUM\\(i.timeSpent\\)").
		WithArgs("PROJ").
		WillReturnRows(sqlmock.NewRows([]string{"author", "total_time_spent"}).
			AddRow("Alice", 100).
			AddRow("Bob", 50),
		)

	r, w := setupRouterWithPath("/api/v1/analytics/time-spent", handler.TimeSpentAnalytics)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/time-spent?key=PROJ", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestPriorityAnalytics_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repository.DB = sqlx.NewDb(db, "postgres")

	mock.ExpectQuery("SELECT i.priority, COUNT").
		WithArgs("PROJ").
		WillReturnRows(sqlmock.NewRows([]string{"priority", "count"}).
			AddRow("High", 3).
			AddRow("Low", 2),
		)

	r, w := setupRouterWithPath("/api/v1/analytics/priority", handler.PriorityAnalytics)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/priority?key=PROJ", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}
