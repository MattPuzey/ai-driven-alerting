package pagerduty_ingestion

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

var (
	pagerDutyAPIURL = "https://api.pagerduty.com/incidents"
)

const (
	pagerDutyAPIToken = "YOUR_PAGERDUTY_API_TOKEN" // Don't hard-code this in a real application
)

// PagerDutyClient communicates with the PagerDuty API
type PagerDutyClient struct {
	client   *http.Client
	apiToken string
}

// Incident represents the structure of an incident in PagerDuty's API response
type Incident struct {
	ID          string `json:"id"`
	Summary     string `json:"summary"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

// IncidentListResponse is the expected structure of the list incidents API call response
type IncidentListResponse struct {
	Incidents []Incident `json:"incidents"`
}

// NewPagerDutyClient creates a new PagerDutyClient
func NewPagerDutyClient(apiToken string) *PagerDutyClient {
	return &PagerDutyClient{
		client:   &http.Client{},
		apiToken: apiToken,
	}
}

// GetIncidents retrieves incidents from PagerDuty
func (pd *PagerDutyClient) GetIncidents() ([]Incident, error) {
	req, err := http.NewRequest("GET", pagerDutyAPIURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Token token="+pd.apiToken)
	req.Header.Add("Accept", "application/vnd.pagerduty+json;version=2")

	resp, err := pd.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response IncidentListResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Incidents, nil
}

func main() {
	client := NewPagerDutyClient(pagerDutyAPIToken)
	incidents, err := client.GetIncidents()
	if err != nil {
		log.Fatalf("Error fetching incidents: %v", err)
	}

	for _, incident := range incidents {
		fmt.Printf("ID: %s, Summary: %s, Status: %s\n", incident.ID, incident.Summary, incident.Status)
	}
}
