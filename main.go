package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	openAIAPIToken      = "YOUR_OPENAI_API_TOKEN"
	maxTokens           = 50
	serverPort          = ":8080"
	optimizeRoute       = "/optimize-threshold"
	ingestMetricRoute   = "/ingest-metric"
	ingestIncidentRoute = "/ingest-incident"
)

var (
	chatGPTAPIURL = "https://api.openai.com/v1/engines/davinci-codex/completions"
)

type ChatGPTRequest struct {
	Prompt    string `json:"prompt"`
	MaxTokens int    `json:"max_tokens"`
}

type ChatGPTResponse struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

type Metric struct {
	ID        int       `json:"id"`
	MetricKey string    `json:"metric_key"`
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

type Incident struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
}

var metrics []Metric
var incidents []Incident

func main() {
	// manual note: why is there endpoints for ingestion when scheduling could be done in the application in the first
	// instance? a separate component for ingestion to a store would be better
	http.HandleFunc(optimizeRoute, OptimizeThresholdHandler)
	http.HandleFunc(ingestMetricRoute, IngestMetricHandler)
	http.HandleFunc(ingestIncidentRoute, IngestIncidentHandler)
	go mockAsyncMetricIngestion() // Start asynchronous metric ingestion
	http.ListenAndServe(serverPort, nil)
}

func OptimizeThresholdHandler(w http.ResponseWriter, r *http.Request) {
	// Your threshold optimization logic here
	sendSuccessResponse(w, map[string]string{"recommendation": "Adjust the threshold to 85%"})
}

func IngestMetricHandler(w http.ResponseWriter, r *http.Request) {
	requestData, err := getRequestData(r)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request data")
		return
	}

	metricKey := requestData["metric_key"]
	value := requestData["value"]
	valueFloat, err := strconv.ParseFloat(value, 64)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid value format")
		return
	}

	metric := Metric{ID: len(metrics) + 1, MetricKey: metricKey, Value: valueFloat, Timestamp: time.Now()}
	metrics = append(metrics, metric)

	sendSuccessResponse(w, map[string]string{"message": "Metric ingested successfully"})
}

func IngestIncidentHandler(w http.ResponseWriter, r *http.Request) {
	requestData, err := getRequestData(r)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request data")
		return
	}

	description := requestData["description"]
	incident := Incident{ID: len(incidents) + 1, Description: description, Timestamp: time.Now()}
	incidents = append(incidents, incident)

	sendSuccessResponse(w, map[string]string{"message": "Incident ingested successfully"})
}

func mockAsyncMetricIngestion() {
	// Simulate asynchronous metric ingestion from Prometheus
	for {
		time.Sleep(30 * time.Second) // Simulate every 30 seconds
		metricKey := "cpu_usage"
		value := fmt.Sprintf("%.2f", rand.Float64()*100)
		resp, _ := http.Post(fmt.Sprintf("http://localhost%s%s", serverPort, ingestMetricRoute),
			"application/json", strings.NewReader(fmt.Sprintf(`{"metric_key": "%s", "value": "%s"}`, metricKey, value)))
		resp.Body.Close()
	}
}

func getRequestData(r *http.Request) (map[string]string, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	var requestData map[string]string
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		return nil, err
	}

	return requestData, nil
}

func generatePrompt(incidentDetails, historicalData string) string {
	return fmt.Sprintf("Based on historical data, it seems that this incident is similar to past false alarms. "+
		"To make the alert less sensitive, you can consider adjusting the alert threshold to a higher value like 85%%. "+
		"This should help in reducing unnecessary alerts while still capturing genuine incidents. Please review and test "+
		"this change before applying it in the production environment.\n\nIncident Details: %s\nHistorical Data: %s", incidentDetails, historicalData)
}

func GetThresholdRecommendationFromChatGPT(prompt string) (string, error) {
	reqBody, _ := json.Marshal(ChatGPTRequest{
		Prompt:    prompt,
		MaxTokens: maxTokens,
	})

	req, _ := http.NewRequest("POST", chatGPTAPIURL, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", openAIAPIToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var chatGPTResponse ChatGPTResponse
	err = json.NewDecoder(resp.Body).Decode(&chatGPTResponse)
	if err != nil {
		return "", err
	}

	return chatGPTResponse.Choices[0].Text, nil
}

func sendSuccessResponse(w http.ResponseWriter, responseData interface{}) {
	sendResponse(w, http.StatusOK, responseData)
}

func sendErrorResponse(w http.ResponseWriter, statusCode int, errorMessage string) {
	sendResponse(w, statusCode, map[string]string{"error": errorMessage})
}

func sendResponse(w http.ResponseWriter, statusCode int, responseData interface{}) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}
