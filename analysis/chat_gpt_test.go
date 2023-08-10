package analysis

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAskGPT(t *testing.T) {
	tests := []struct {
		name         string
		mockResponse string
		statusCode   int
		expectedErr  bool
		expectedResp string
	}{
		{
			"Valid GPT Response",
			`{"choices": [{"text": "The metric might be frequently alerted due to misconfigured thresholds."}]}`,
			http.StatusOK,
			false,
			"The metric might be frequently alerted due to misconfigured thresholds.",
		},
		//{
		//	"Bad GPT response format",
		//	`{"invalid": "response"}`,
		//	http.StatusOK,
		//	true,
		//	"",
		//},
		//{
		//	"API error response",
		//	`{"error": {"message": "API error"}}`,
		//	http.StatusInternalServerError,
		//	true,
		//	"",
		//},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(test.statusCode)
				w.Write([]byte(test.mockResponse))
			}))

			// Create a mock GPT client with a mock server URL
			gptClient := NewOpenAIClient("testToken")
			openAIEndpoint = mockServer.URL // Override the URL with our mock server address
			defer mockServer.Close()

			response, err := gptClient.AskGPT("Test question")

			if test.expectedErr && err == nil {
				t.Errorf("Expected error, got nil")
			} else if !test.expectedErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			} else if response != test.expectedResp {
				t.Errorf("Expected response: %s, got %s", test.expectedResp, response)
			}
		})
	}
}
