package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestIngestIncidentHandler(t *testing.T) {
	reqBody := map[string]string{
		"description": "Network outage",
	}
	reqBodyJSON, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", ingestIncidentRoute, strings.NewReader(string(reqBodyJSON)))
	w := httptest.NewRecorder()

	IngestIncidentHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var responseData map[string]string
	json.Unmarshal(w.Body.Bytes(), &responseData)

	if responseData["message"] != "Incident ingested successfully" {
		t.Errorf("Expected message: %s\nGot: %s", "Incident ingested successfully", responseData["message"])
	}
}

func TestMockAsyncIncidentIngestion(t *testing.T) {
	go mockAsyncIncidentIngestion() // Start async incident ingestion
	time.Sleep(2 * time.Second)     // Wait for async ingestion to run

	if len(incidents) == 0 {
		t.Error("Expected incidents to be ingested asynchronously")
	}
}
