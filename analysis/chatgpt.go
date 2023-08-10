package analysis

import (
	"bytes"
	"encoding/json"
	"net/http"
)

var (
	openAIEndpoint = "https://api.openai.com/v1/engines/davinci/completions"
)

const (
	openAIApiKey = "YOUR_OPENAI_API_KEY" // Never hard-code this in a real application.
)

type OpenAIClient struct {
	client *http.Client
	apiKey string
}

func NewOpenAIClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{
		client: &http.Client{},
		apiKey: apiKey,
	}
}

func (oai *OpenAIClient) AskGPT(question string) (string, error) {
	requestBody := map[string]string{
		"prompt":     question,
		"max_tokens": "150",
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", openAIEndpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", "Bearer "+oai.apiKey)
	req.Header.Add("Content-Type", "application/json")
	resp, err := oai.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result["choices"].([]interface{})[0].(map[string]interface{})["text"].(string), nil
}
