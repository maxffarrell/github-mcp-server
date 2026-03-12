package transport

import (
	"net/http"

	"github.com/github/github-mcp-server/pkg/http/headers"
)

// APIVersionTransport is an http.RoundTripper that sets the GitHub REST API
// version header on every request. This overrides the default version set by
// the go-github library, ensuring all requests use the configured API version.
type APIVersionTransport struct {
	// Transport is the underlying HTTP transport. If nil, http.DefaultTransport is used.
	Transport http.RoundTripper
}

// RoundTrip implements http.RoundTripper.
func (t *APIVersionTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	transport := t.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	req = req.Clone(req.Context())
	req.Header.Set(headers.GitHubAPIVersionHeader, headers.GitHubAPIVersion)

	return transport.RoundTrip(req)
}
