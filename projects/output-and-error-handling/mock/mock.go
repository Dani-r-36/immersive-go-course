package mock

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"github.com/stretchr/testify/require"
)

type responseOrError struct {
	response *http.Response
	err      error
}
//mocked tester with response of error, so can mock http response
type MockRoundTripper struct {
	t            *testing.T
	responses    []responseOrError
	requestCount int
}

func NewMockRoundTripper(t *testing.T) *MockRoundTripper {
	return &MockRoundTripper{
		t: t,
	}
}
//mocked http response to send back from given statuscode, body and header
func (m *MockRoundTripper) StubResponse(statusCode int, body string, header http.Header) {
	// We need to stub out a fair bit of the HTTP response in for the Go HTTP client to accept our response.
	response := &http.Response{
		Body:          io.NopCloser(strings.NewReader(body)),
		ContentLength: int64(len(body)),
		Header:        header,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Status:        fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode)),
		StatusCode:    statusCode,
	}
	m.responses = append(m.responses, responseOrError{response: response})
}


//check if mocked is right with requested data
func (m *MockRoundTripper) AssertGotRightCalls() {
	m.t.Helper()

	require.Equalf(m.t, len(m.responses), m.requestCount, "Expected %d requests, got %d", len(m.responses), m.requestCount)
}


/*this is what is used to mock the round trip method in http and returns the reponses slice and increments the requests count, 
to check that the number of mock request (m.responses) match the stubbed repsoness (m.request count)*/
func (m *MockRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	m.t.Helper()

	if m.requestCount >= len(m.responses) {
		m.t.Fatalf("MockRoundTripper expected %d requests but got more", len(m.responses))
	}
	resp := m.responses[m.requestCount]
	m.requestCount += 1
	return resp.response, resp.err
}