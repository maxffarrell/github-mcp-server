package transport

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/github/github-mcp-server/pkg/http/headers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIVersionTransport(t *testing.T) {
	t.Parallel()

	var capturedVersion string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedVersion = r.Header.Get(headers.GitHubAPIVersionHeader)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	tr := &APIVersionTransport{
		Transport: http.DefaultTransport,
	}

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	require.NoError(t, err)

	resp, err := tr.RoundTrip(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, headers.GitHubAPIVersion, capturedVersion)
}

func TestAPIVersionTransport_OverridesExistingHeader(t *testing.T) {
	t.Parallel()

	var capturedVersion string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedVersion = r.Header.Get(headers.GitHubAPIVersionHeader)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	tr := &APIVersionTransport{
		Transport: http.DefaultTransport,
	}

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	require.NoError(t, err)
	// Set an old version, simulating what go-github does
	req.Header.Set(headers.GitHubAPIVersionHeader, "2022-11-28")

	resp, err := tr.RoundTrip(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, headers.GitHubAPIVersion, capturedVersion)
}

func TestAPIVersionTransport_NilTransport(t *testing.T) {
	t.Parallel()

	var capturedVersion string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedVersion = r.Header.Get(headers.GitHubAPIVersionHeader)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// nil Transport should fall back to http.DefaultTransport
	tr := &APIVersionTransport{}

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	require.NoError(t, err)

	resp, err := tr.RoundTrip(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, headers.GitHubAPIVersion, capturedVersion)
}

func TestAPIVersionTransport_DoesNotMutateOriginalRequest(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	tr := &APIVersionTransport{
		Transport: http.DefaultTransport,
	}

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	require.NoError(t, err)
	req.Header.Set(headers.GitHubAPIVersionHeader, "2022-11-28")

	resp, err := tr.RoundTrip(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Original request should still have the old version
	assert.Equal(t, "2022-11-28", req.Header.Get(headers.GitHubAPIVersionHeader))
}
