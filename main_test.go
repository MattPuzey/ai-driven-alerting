package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGeneratePrompt(t *testing.T) {
	incidentDetails := "Incident: CPU spike"
	historicalData := "Past incidents: ..."
	expectedPrompt := "Based on historical data, it seems that this incident is similar to past false alarms. To make the alert less sensitive, you can consider adjusting the alert threshold to a higher value like 85%%. This should help in reducing unnecessary alerts while still capturing genuine incidents. Please review and test this change before applying it in the production environment.\n\nIncident Details: Incident: CPU spike\nHistorical Data: Past incidents: ..."
	actualPrompt := generatePrompt(incidentDetails, historicalData)

	if actualPrompt != expectedPrompt {
		t.Errorf("Expected prompt:\n%s\n\nGot:\n%s", expectedPrompt, actualPrompt)
	}
}

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
	actualRecommendation := GetThresholdRecommendationFromChatGPT(prompt)

	if actualRecommendation != expectedRecommendation {
		t.Errorf("Expected recommendation: %s\nGot: %s", expectedRecommendation, actualRecommendation)
	}
}

func TestOptimizeThresholdHandler(t *testing.T) {
	reqBody := map[string]string{
		"incident_details": "CPU spike incident",
		"historical_data":  "Past incidents: ...",
	}
	reqBodyJSON, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", optimizeRoute, bytes.NewBuffer(reqBodyJSON))
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
