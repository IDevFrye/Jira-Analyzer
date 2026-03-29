//go:build integration
// +build integration

package integration

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNegativeScenarios(t *testing.T) {
	client := http.Client{Timeout: 20 * time.Second}
	baseURL := getBackendBaseURL()

	t.Run("GET /projects/{id} with invalid id returns bad request", func(t *testing.T) {
		resp, err := client.Get(fmt.Sprintf("%s/projects/not-an-int", baseURL))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("GET /analytics/status-distribution without key returns bad request", func(t *testing.T) {
		resp, err := client.Get(fmt.Sprintf("%s/analytics/status-distribution", baseURL))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("GET /compare/priority without key returns bad request", func(t *testing.T) {
		resp, err := client.Get(fmt.Sprintf("%s/compare/priority", baseURL))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("GET unknown endpoint returns not found", func(t *testing.T) {
		resp, err := client.Get(fmt.Sprintf("%s/connector/unknown", baseURL))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
