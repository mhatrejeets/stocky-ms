//go:build integration

package tests

import (
	"context"
	"testing"

	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestRewardService_Integration(t *testing.T) {
	ctx := context.Background()
	// Start PostgreSQL container
	pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:15",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_DB":       "assignment",
				"POSTGRES_USER":     "stocky",
				"POSTGRES_PASSWORD": "password",
			},
			WaitingFor: wait.ForListeningPort("5432/tcp"),
		},
		Started: true,
	})
	assert.NoError(t, err)
	defer pgContainer.Terminate(ctx)

	// Start Redis container
	redisContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "redis:7",
			ExposedPorts: []string{"6379/tcp"},
			WaitingFor:   wait.ForListeningPort("6379/tcp"),
		},
		Started: true,
	})
	assert.NoError(t, err)
	defer redisContainer.Terminate(ctx)

	// TODO: Run migrations here (use os/exec or Go migration lib)

	// Start API server (use httptest or real server)
	server := httptest.NewServer(nil) // Replace nil with your Gin router
	defer server.Close()

	// POST /api/v1/reward
	body := bytes.NewBufferString(`{"stock_symbol":"RELIANCE","shares":"1.000000","rewarded_at":"2025-09-25T11:30:00Z"}`)
	resp, err := http.Post(server.URL+"/api/v1/reward", "application/json", body)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// TODO: Validate DB record exists

	// GET /api/v1/portfolio/{userId}
	resp, err = http.Get(server.URL + "/api/v1/portfolio/user-1")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	// TODO: Validate response body
}
