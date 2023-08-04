package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetThresholdRecommendationFromChatGPT(t *testing.T) {
	expectedRecommendation := "Try increasing the threshold to 85% to reduce false alarms."
	mockResponse := ChatGPTResponse{
		Choices: []struct {
			Text string `json:"text"`
		}{
			{Text: expectedRecommendation},
		},
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(mockResponseJSON)
	}))
	defer server.Close()

	chatGPTAPIURL = server.URL

	prompt := "Based on historical data..."
	actualRecommendation, err := GetThresholdRecommendationFromChatGPT(prompt)
	if err != nil {
		t.Errorf("Error getting threshold recommendation: %v", err)
	}

	if actualRecommendation != expectedRecommendation {
		t.Errorf("Expected recommendation: %s\nGot: %s", expectedRecommendation, actualRecommendation)
	}
}

func TestOptimizeThresholdHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", optimizeRoute, nil)
	w := httptest.NewRecorder()

	OptimizeThresholdHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var responseData map[string]string
	json.Unmarshal(w.Body.Bytes(), &responseData)

	if responseData["recommendation"] == "" {
		t.Errorf("Expected non-empty recommendation in response")
	}
}

func TestSendResponse(t *testing.T) {
	expectedStatus := http.StatusOK
	expectedData := map[string]string{"recommendation": "Adjust the threshold to 85%"}

	w := httptest.NewRecorder()
	sendResponse(w, expectedStatus, expectedData)

	if w.Code != expectedStatus {
		t.Errorf("Expected status code %d, got %d", expectedStatus, w.Code)
	}

	var responseData map[string]string
	json.Unmarshal(w.Body.Bytes(), &responseData)

	if responseData["recommendation"] != expectedData["recommendation"] {
		t.Errorf("Expected recommendation: %s\nGot: %s", expectedData["recommendation"], responseData["recommendation"])
	}
}

func TestGetRequestData(t *testing.T) {
	reqBody := map[string]string{
		"metric_key": "cpu_usage",
		"value":      "80.5",
	}
	reqBodyJSON, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", ingestMetricRoute, strings.NewReader(string(reqBodyJSON)))

	data, err := getRequestData(req)
	if err != nil {
		t.Errorf("Error getting request data: %v", err)
	}

	if data["metric_key"] != "cpu_usage" {
		t.Errorf("Expected metric_key: %s, got: %s", "cpu_usage", data["metric_key"])
	}

	if data["value"] != "80.5" {
		t.Errorf("Expected value: %s, got: %s", "80.5", data["value"])
	}
}

func TestGetRequestDataInvalidJSON(t *testing.T) {
	req, _ := http.NewRequest("POST", optimizeRoute, strings.NewReader("invalid json"))

	_, err := getRequestData(req)
	if err == nil {
		t.Error("Expected error when parsing invalid JSON")
	}
}

func TestSendErrorResponse(t *testing.T) {
	expectedStatus := http.StatusBadRequest
	expectedMessage := "Bad request"

	w := httptest.NewRecorder()
	sendErrorResponse(w, expectedStatus, expectedMessage)

	if w.Code != expectedStatus {
		t.Errorf("Expected status code %d, got %d", expectedStatus, w.Code)
	}

	var responseData map[string]string
	json.Unmarshal(w.Body.Bytes(), &responseData)

	if responseData["error"] != expectedMessage {
		t.Errorf("Expected error message: %s\nGot: %s", expectedMessage, responseData["error"])
	}
}

func TestSendSuccessResponse(t *testing.T) {
	expectedData := map[string]string{"message": "Success"}

	w := httptest.NewRecorder()
	sendSuccessResponse(w, expectedData)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var responseData map[string]string
	json.Unmarshal(w.Body.Bytes(), &responseData)

	if responseData["message"] != expectedData["message"] {
		t.Errorf("Expected message: %s\nGot: %s", expectedData["message"], responseData["message"])
	}
}
