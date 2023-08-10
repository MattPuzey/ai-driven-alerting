package prometheus_ingestion

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test the getMetrics function
func TestGetMetrics(t *testing.T) {
	tests := []struct {
		mockResponse string
		statusCode   int
		expectedErr  bool
	}{
		{`{"status": "success", "data": {"resultType": "vector", "result": [{"metric": {}, "value": [1628919140, "1"]}]}}`, http.StatusOK, false},
		//{`{"status": "failure", "data": {"error": "Bad request"}}`, http.StatusBadRequest, true},
	}

	for _, test := range tests {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(test.statusCode)
			w.Write([]byte(test.mockResponse))
		}))

		// Override the Prometheus URL with our mock server address
		prometheusURL = mockServer.URL
		defer mockServer.Close()

		data, err := getMetrics("up")

		if test.expectedErr && err == nil {
			t.Errorf("Expected error, got nil")
		} else if !test.expectedErr && err != nil {
			t.Errorf("Unexpected error: %v", err)
		} else if !test.expectedErr && len(data.Data.Result) == 0 {
			t.Errorf("Expected data, got none")
		}
	}
}
