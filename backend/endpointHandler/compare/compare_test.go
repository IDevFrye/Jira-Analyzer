package compare

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/endpointhandler/repository"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

import "database/sql/driver"

func toDriverValues(args []interface{}) []driver.Value {
	values := make([]driver.Value, len(args))
	for i, v := range args {
		values[i] = v
	}
	return values
}

func setupDB(t *testing.T) (sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	repository.DB = sqlx.NewDb(db, "postgres")

	return mock, func() {
		db.Close()
	}
}

func setupRouterWithHandler(path string, handlerFunc gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET(path, handlerFunc)
	return r
}

func TestCompareTimeOpen_Success(t *testing.T) {
	mock, closeDB := setupDB(t)
	defer closeDB()

	rows := sqlmock.NewRows([]string{"range", "count"}).
		AddRow("0-1", 10).
		AddRow("1-2", 5)

	mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT
				CASE
					WHEN age <= 1 THEN '0-1'
					WHEN age <= 2 THEN '1-2'
					WHEN age <= 3 THEN '2-3'
					WHEN age <= 5 THEN '3-5'
					WHEN age <= 7 THEN '5-7'
					WHEN age <= 10 THEN '7-10'
					WHEN age <= 14 THEN '10-14'
					WHEN age <= 21 THEN '14-21'
					WHEN age <= 30 THEN '21-30'
					ELSE '30+'
				END AS range,
				COUNT(*) AS count
			FROM (
				SELECT DATE_PART('day', NOW() - i.createdTime) AS age
				FROM Projects p
				JOIN Issue i ON p.id = i.projectId
				WHERE i.status NOT IN ('Closed', 'Resolved') AND p.key = $1
			) sub
			GROUP BY range
			ORDER BY MIN(age)
		`)).
		WithArgs("TESTKEY").
		WillReturnRows(rows)

	r := setupRouterWithHandler("/api/v1/compare/time-open", CompareTimeOpen)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/compare/time-open?key=TESTKEY", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string][]struct {
		Range string `json:"range"`
		Count int    `json:"count"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp["TESTKEY"], 2)
	assert.Equal(t, 10, resp["TESTKEY"][0].Count)
	assert.Equal(t, 5, resp["TESTKEY"][1].Count)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCompareTimeOpen_MissingKey(t *testing.T) {
	r := setupRouterWithHandler("/api/v1/compare/time-open", CompareTimeOpen)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/compare/time-open", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "missing ?key")
}

// --- CompareStatusDistribution ---

func TestCompareStatusDistribution(t *testing.T) {
	mock, closeDB := setupDB(t)
	defer closeDB()

	rows := sqlmock.NewRows([]string{"project", "status", "count"}).
		AddRow("PROJ1", "Open", 3).
		AddRow("PROJ1", "In Progress", 7).
		AddRow("PROJ2", "Open", 2)

	keys := []string{"PROJ1", "PROJ2"}
	query := `
		SELECT 
			p.key AS project,
			i.status,
			COUNT(*) AS count
		FROM Projects p
		JOIN Issue i ON p.id = i.projectId
		WHERE p.key IN (?)
		GROUP BY p.key, i.status
		ORDER BY p.key, i.status
	`
	rebQuery, args, err := sqlx.In(query, keys)
	assert.NoError(t, err)
	rebQuery = repository.DB.Rebind(rebQuery)

	mock.ExpectQuery(regexp.QuoteMeta(rebQuery)).
		WithArgs(toDriverValues(args)...).
		WillReturnRows(rows)

	r := setupRouterWithHandler("/api/v1/compare/status-distribution", CompareStatusDistribution)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/compare/status-distribution?key=PROJ1,PROJ2", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]map[string]int
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, 3, resp["PROJ1"]["Open"])
	assert.Equal(t, 7, resp["PROJ1"]["In Progress"])
	assert.Equal(t, 2, resp["PROJ2"]["Open"])

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCompareStatusDistribution_DBError(t *testing.T) {
	mock, closeDB := setupDB(t)
	defer closeDB()

	query := `
		SELECT 
			p.key AS project,
			i.status,
			COUNT(*) AS count
		FROM Projects p
		JOIN Issue i ON p.id = i.projectId
		WHERE p.key IN (?)
		GROUP BY p.key, i.status
		ORDER BY p.key, i.status
	`
	rebQuery, _, err := sqlx.In(query, []string{"PROJ1"})
	assert.NoError(t, err)
	rebQuery = repository.DB.Rebind(rebQuery)

	mock.ExpectQuery(regexp.QuoteMeta(rebQuery)).
		WillReturnError(errors.New("db failure"))

	r := setupRouterWithHandler("/api/v1/compare/status-distribution", CompareStatusDistribution)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/compare/status-distribution?key=PROJ1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "db failure")

	assert.NoError(t, mock.ExpectationsWereMet())
}

// --- CompareTimeSpent ---

func TestCompareTimeSpent(t *testing.T) {
	mock, closeDB := setupDB(t)
	defer closeDB()

	rows := sqlmock.NewRows([]string{"project", "author", "total_time_spent"}).
		AddRow("PROJ1", "Alice", 100).
		AddRow("PROJ1", "Bob", 50).
		AddRow("PROJ2", "Charlie", 75)

	keys := []string{"PROJ1", "PROJ2"}
	query := `
		SELECT 
			p.key AS project,
			a.name AS author,
			SUM(i.timeSpent) AS total_time_spent
		FROM Projects p
		JOIN Issue i ON p.id = i.projectId
		JOIN Author a ON a.id = i.authorId
		WHERE p.key IN (?) AND i.timeSpent IS NOT NULL
		GROUP BY p.key, a.name
		ORDER BY p.key, total_time_spent DESC
	`
	rebQuery, args, err := sqlx.In(query, keys)
	assert.NoError(t, err)
	rebQuery = repository.DB.Rebind(rebQuery)

	mock.ExpectQuery(regexp.QuoteMeta(rebQuery)).
		WithArgs(toDriverValues(args)...).
		WillReturnRows(rows)

	r := setupRouterWithHandler("/api/v1/compare/time-spent", CompareTimeSpent)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/compare/time-spent?key=PROJ1,PROJ2", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]struct {
		Authors []struct {
			Author         string `json:"author"`
			TotalTimeSpent int    `json:"total_time_spent"`
		} `json:"authors"`
	}

	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Len(t, resp["PROJ1"].Authors, 2)
	assert.Len(t, resp["PROJ2"].Authors, 1)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// --- ComparePriority ---

func TestComparePriority(t *testing.T) {
	mock, closeDB := setupDB(t)
	defer closeDB()

	rows := sqlmock.NewRows([]string{"project", "priority", "count"}).
		AddRow("PROJ1", "High", 5).
		AddRow("PROJ1", "Low", 3).
		AddRow("PROJ2", "Medium", 7)

	keys := []string{"PROJ1", "PROJ2"}
	query := `
		SELECT 
			p.key AS project,
			i.priority,
			COUNT(*) AS count
		FROM Projects p
		JOIN Issue i ON p.id = i.projectId
		WHERE p.key IN (?)
		GROUP BY p.key, i.priority
		ORDER BY p.key, i.priority
	`
	rebQuery, args, err := sqlx.In(query, keys)
	assert.NoError(t, err)
	rebQuery = repository.DB.Rebind(rebQuery)

	mock.ExpectQuery(regexp.QuoteMeta(rebQuery)).
		WithArgs(toDriverValues(args)...).
		WillReturnRows(rows)

	r := setupRouterWithHandler("/api/v1/compare/priority", ComparePriority)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/compare/priority?key=PROJ1,PROJ2", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]map[string]int
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, 5, resp["PROJ1"]["High"])
	assert.Equal(t, 3, resp["PROJ1"]["Low"])
	assert.Equal(t, 7, resp["PROJ2"]["Medium"])

	assert.NoError(t, mock.ExpectationsWereMet())
}
