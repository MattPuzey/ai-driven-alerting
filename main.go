package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	openAIAPIToken        = "YOUR_OPENAI_API_TOKEN"
	maxTokens             = 50
	serverPort            = ":8080"
	optimizeRoute         = "/optimize-threshold"
	ingestIncidentRoute   = "/ingest-incident"
	pagerDutyAPIKey       = "YOUR_PAGERDUTY_API_KEY"
	pagerDutyIncidentsURL = "https://api.pagerduty.com/incidents"
	ingestMetricRoute     = "/ingest-metric"
	prometheusBaseURL     = "http://prometheus-server:9090" // Replace with your Prometheus server URL
	prometheusQueryPath   = "/api/v1/query"
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
	http.HandleFunc(optimizeRoute, OptimizeThresholdHandler)
	// TODO: manual note to remove these endpoints as they do not make sense, data ingestion should be done out of band
	// of request processing...
	http.HandleFunc(ingestMetricRoute, IngestMetricHandler)
	http.HandleFunc(ingestIncidentRoute, IngestIncidentHandler)
	go mockAsyncIncidentIngestion()
	go mockAsyncMetricIngestion()
	go http.ListenAndServe(serverPort, nil)
}

func getPrometheusMetricValue(metricKey string) string {
	// Replace with actual Prometheus query logic
	query := fmt.Sprintf(`sum(%s)`, metricKey)
	resp, err := http.Get(fmt.Sprintf("%s%s?query=%s", prometheusBaseURL, prometheusQueryPath, query))
	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		return ""
	}
	defer resp.Body.Close()

	//body, _ := io.ReadAll(resp.Body)
	// Parse and extract metric value from Prometheus response
	// This step depends on the Prometheus query response format

	return "123.45" // Replace with actual extracted metric value
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

func IngestPagerDutyIncident(incidentDescription string) error {
	payload := map[string]interface{}{
		"type":        "incident",
		"title":       incidentDescription,
		"description": "This is a generated incident description.",
		"service": map[string]string{
			"id":   "YOUR_PAGERDUTY_SERVICE_ID",
			"type": "service_reference",
		},
	}

	payloadBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", pagerDutyIncidentsURL, bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token token=%s", pagerDutyAPIKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

//func IngestMetricToChatGPT(metric Metric) {
//	// Generate a prompt using metric information and call Chat GPT
//	prompt := fmt.Sprintf("Based on historical data, it seems that the current metric value for '%s' is high. This could indicate a potential issue. To make the alert more sensitive, you can consider adjusting the alert threshold to a lower value. This should help in capturing such incidents early. Please review and test this change before applying it in the production environment.\n\nMetric Details:\nMetric Key: %s\nMetric Value: %.2f", metric.MetricKey, metric.MetricKey, metric.Value)
//
//	// Call Chat GPT for recommendations
//	recommendation, err := GetThresholdRecommendationFromChatGPT(prompt)
//	if err != nil {
//		fmt.Printf("Error generating recommendation: %v\n", err)
//		return
//	}
//
//	fmt.Printf("Recommendation: %s\n", recommendation)
//}

func OptimizeThresholdHandler(w http.ResponseWriter, r *http.Request) {
	requestData, err := getRequestData(r)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request data")
		return
	}
	prompt := requestData["prompt"]

	recommendation, err := GetThresholdRecommendationFromChatGPT(prompt)
	response := map[string]string{"recommendation": recommendation}

	sendSuccessResponse(w, response)
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

		// Replace with actual Prometheus query and metric extraction logic
		metricKey := "cpu_usage"
		value := getPrometheusMetricValue(metricKey)
		resp, _ := http.Post(fmt.Sprintf("http://localhost%s%s", serverPort, ingestMetricRoute), "application/json", strings.NewReader(fmt.Sprintf(`{"metric_key": "%s", "value": "%s"}`, metricKey, value)))
		resp.Body.Close()
	}
}

func mockAsyncIncidentIngestion() {
	// Simulate asynchronous incident ingestion from Prometheus
	for {
		time.Sleep(10 * time.Second) // Simulate every 10 seconds
		description := "CPU spike incident"
		incident := Incident{ID: len(incidents) + 1, Description: description, Timestamp: time.Now()}
		incidents = append(incidents, incident)
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
