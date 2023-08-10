package pagerduty_ingestion

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetIncidents(t *testing.T) {
	tests := []struct {
		name          string
		mockResponse  string
		statusCode    int
		expectedErr   bool
		expectedCount int
	}{
		{
			"Valid response",
			`{"incidents": [{"id": "P1234", "summary": "Test Incident", "status": "triggered", "description": "Sample incident for testing."}]}`,
			http.StatusOK,
			false,
			1,
		},
		{
			"Bad response format",
			`{"invalid": "response"}`,
			http.StatusOK,
			true,
			0,
		},
		{
			"API error response",
			`{"error": {"message": "API error", "code": 5000}}`,
			http.StatusInternalServerError,
			true,
			0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(test.statusCode)
				w.Write([]byte(test.mockResponse))
			}))

			// Create a PagerDuty client with a mock server URL
			pdClient := NewPagerDutyClient("testToken")
			pagerDutyAPIURL = mockServer.URL // Override the URL with our mock server address
			defer mockServer.Close()

			incidents, err := pdClient.GetIncidents()

			if test.expectedErr && err == nil {
				t.Errorf("Expected error, got nil")
			} else if !test.expectedErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			} else if len(incidents) != test.expectedCount {
				t.Errorf("Expected %d incidents, got %d", test.expectedCount, len(incidents))
			}
		})
	}
}
