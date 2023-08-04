package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	chatGPTAPIURL  = "https://api.openai.com/v1/engines/davinci-codex/completions"
	openAIAPIToken = "YOUR_OPENAI_API_TOKEN"
	maxTokens      = 50
	serverPort     = ":8080"
	optimizeRoute  = "/optimize-threshold"
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

func main() {
	http.HandleFunc(optimizeRoute, OptimizeThresholdHandler)
	http.ListenAndServe(serverPort, nil)
}

func OptimizeThresholdHandler(w http.ResponseWriter, r *http.Request) {
	requestData, err := getRequestData(r)
	if err != nil {
		sendResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		return
	}

	incidentDetails := requestData["incident_details"]
	historicalData := requestData["historical_data"]

	prompt := generatePrompt(incidentDetails, historicalData)
	recommendation, err := GetThresholdRecommendationFromChatGPT(prompt)
	if err != nil {
		sendResponse(w, http.StatusInternalServerError, map[string]string{"error": "ChatGPT API request failed"})
		return
	}

	sendResponse(w, http.StatusOK, map[string]string{"recommendation": recommendation})
}

func getRequestData(r *http.Request) (map[string]string, error) {
	body, err := ioutil.ReadAll(r.Body)
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
	return fmt.Sprintf("Based on historical data, it seems that this incident is similar to past false alarms. To make the alert less sensitive, you can consider adjusting the alert threshold to a higher value like 85%%. This should help in reducing unnecessary alerts while still capturing genuine incidents. Please review and test this change before applying it in the production environment.\n\nIncident Details: %s\nHistorical Data: %s", incidentDetails, historicalData)
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

func sendResponse(w http.ResponseWriter, status int, responseData interface{}) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}
