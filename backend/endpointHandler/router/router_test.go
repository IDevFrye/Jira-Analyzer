package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/endpointhandler/config"
)

func TestSetupRouter_ProjectsRouteExists(t *testing.T) {
	cfg := &config.Config{}
	r := SetupRouter(cfg)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/projects", nil)
	r.ServeHTTP(w, req)

	if w.Code == http.StatusNotFound {
		t.Fatalf("expected /api/v1/projects route to be registered, got 404")
	}
}

func TestSetupRouter_AnalyticsRouteExists(t *testing.T) {
	cfg := &config.Config{}
	r := SetupRouter(cfg)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/analytics/time-open", nil)
	r.ServeHTTP(w, req)

	if w.Code == http.StatusNotFound {
		t.Fatalf("expected /api/v1/analytics/time-open route to be registered, got 404")
	}
}

func TestSetupRouter_ThroughputRouteExists(t *testing.T) {
	cfg := &config.Config{}
	r := SetupRouter(cfg)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/analytics/throughput?key=PRJ", nil)
	r.ServeHTTP(w, req)

	if w.Code == http.StatusNotFound {
		t.Fatalf("expected /api/v1/analytics/throughput route to be registered, got 404")
	}
}

