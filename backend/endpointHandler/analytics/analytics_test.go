package analytics

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"github.com/endpointhandler/repository"
)

func setupMockDB(t *testing.T) sqlmock.Sqlmock {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %s", err)
	}
	// конвертируем *sql.DB в *sqlx.DB
	repository.DB = sqlx.NewDb(db, "sqlmock")
	return mock
}

func performRequest(method, path string, handlerFunc gin.HandlerFunc) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(method, path, nil)
	handlerFunc(c)
	return w
}

func TestTimeOpenAnalytics(t *testing.T) {
	mock := setupMockDB(t)

	mock.ExpectQuery("SELECT.*FROM.*Projects p").
		WithArgs("test-project").
		WillReturnRows(sqlmock.NewRows([]string{"range", "count"}).
			AddRow("0-1", 5).
			AddRow("1-2", 3),
		)

	w := performRequest(http.MethodGet, "/analytics/time-open?key=test-project", TimeOpenAnalytics)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestStatusDistribution(t *testing.T) {
	mock := setupMockDB(t)

	mock.ExpectQuery("SELECT i.status, COUNT").
		WithArgs("test-project").
		WillReturnRows(sqlmock.NewRows([]string{"status", "count"}).
			AddRow("Open", 10).
			AddRow("In Progress", 4),
		)

	w := performRequest(http.MethodGet, "/analytics/status-distribution?key=test-project", StatusDistribution)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestTimeSpentAnalytics(t *testing.T) {
	mock := setupMockDB(t)

	mock.ExpectQuery("SELECT a.name AS author").
		WithArgs("test-project").
		WillReturnRows(sqlmock.NewRows([]string{"author", "total_time_spent"}).
			AddRow("Alice", 120).
			AddRow("Bob", 90),
		)

	w := performRequest(http.MethodGet, "/analytics/time-spent?key=test-project", TimeSpentAnalytics)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestPriorityAnalytics(t *testing.T) {
	mock := setupMockDB(t)

	mock.ExpectQuery("SELECT i.priority, COUNT").
		WithArgs("test-project").
		WillReturnRows(sqlmock.NewRows([]string{"priority", "count"}).
			AddRow("High", 7).
			AddRow("Low", 3),
		)

	w := performRequest(http.MethodGet, "/analytics/priority?key=test-project", PriorityAnalytics)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestTimeOpenAnalytics_DBError(t *testing.T) {
	mock := setupMockDB(t)

	mock.ExpectQuery("SELECT.*FROM.*Projects p").
		WithArgs("test-project").
		WillReturnError(fmt.Errorf("db error"))

	w := performRequest(http.MethodGet, "/analytics/time-open?key=test-project", TimeOpenAnalytics)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}
}

func TestStatusDistribution_DBError(t *testing.T) {
	mock := setupMockDB(t)

	mock.ExpectQuery("SELECT i.status, COUNT").
		WithArgs("test-project").
		WillReturnError(fmt.Errorf("db error"))

	w := performRequest(http.MethodGet, "/analytics/status-distribution?key=test-project", StatusDistribution)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}
}

func TestTimeSpentAnalytics_DBError(t *testing.T) {
	mock := setupMockDB(t)

	mock.ExpectQuery("SELECT a.name AS author").
		WithArgs("test-project").
		WillReturnError(fmt.Errorf("db error"))

	w := performRequest(http.MethodGet, "/analytics/time-spent?key=test-project", TimeSpentAnalytics)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}
}

func TestPriorityAnalytics_DBError(t *testing.T) {
	mock := setupMockDB(t)

	mock.ExpectQuery("SELECT i.priority, COUNT").
		WithArgs("test-project").
		WillReturnError(fmt.Errorf("db error"))

	w := performRequest(http.MethodGet, "/analytics/priority?key=test-project", PriorityAnalytics)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}
}
