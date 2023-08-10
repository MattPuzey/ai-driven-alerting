package prometheus_ingestion

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	defaultPrometheusURL = "http://YOUR_PROMETHEUS_SERVER:9090/api/v1/"
)

// MetricFetcher fetches metrics from a Prometheus server
type MetricFetcher struct {
	client *http.Client
	url    string
}

// NewMetricFetcher creates a new MetricFetcher
func NewMetricFetcher(client *http.Client, url string) *MetricFetcher {
	if url == "" {
		url = defaultPrometheusURL
	}
	return &MetricFetcher{client: client, url: url}
}

// MetricData represents the structure of Prometheus response for a query
type MetricData struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Value  []interface{}     `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

func main() {
	// Sample: Collect data for a specific query. Adjust this as per your needs.
	mf := NewMetricFetcher(http.DefaultClient, "")
	query := "up"
	data, err := mf.getMetrics(query)
	if err != nil {
		log.Fatalf("Error fetching metrics: %v", err)
	}

	for _, result := range data.Data.Result {
		fmt.Printf("Metric: %v, Value: %v\n", result.Metric, result.Value[1])
	}
}

// getMetrics retrieves metric data from Prometheus for the given query
func (mf *MetricFetcher) getMetrics(query string) (*MetricData, error) {
	resp, err := mf.client.Get(mf.url + "?query=" + query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch metrics: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var data MetricData
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return &data, nil
}
