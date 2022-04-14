package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"
)

// Variable naming is matched to the API data field names, camel case is preferred.
type (

	// Service objects

	// Job Service object
	JobService struct {
		client *Client
	}

	// Job represents a job returned by the API
	Job struct {
		ID                       uint      `json:"id"`
		IsSample                 bool      `json:"is_sample"`
		ObservableName           string    `json:"observable_name"`
		ObservableClassification string    `json:"observable_classification"`
		FileName                 string    `json:"file_name"`
		FileMimetype             string    `json:"file_mimetype"`
		Status                   string    `json:"status"`
		Tags                     []TagData `json:"tags"`
		ProcessTime              uint      `json:"process_time"`
		NoOfAnalyzersExecuted    string    `json:"no_of_analyzers_executed"`
		NoOfConnectorsExecuted   string    `json:"no_of_connectors_executed"`
	}

	// Tags represents the tag metadata for a job
	TagData struct {
		ID    uint
		Label string
		Color string
	}

	// Client represents the IntelOwl client.
	Client struct {
		apiKey     string
		httpClient *http.Client
		baseURL    *url.URL

		Job *JobService
	}
)

// since IntelOwl is self-hosted, we might remove this and read from configs
const (
	baseURL     = "http://localhost:8081"
	contentType = "application/json"
)

// NewClient returns a go-intelowl client.
func NewClient(apiKey string) (*Client, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	c := &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: time.Minute}, // read timeout from config?
		baseURL:    u,
	}

	c.Job = &JobService{client: c}

	return c, nil
}

// buildReq returns a pre-configured HTTP request.
func (c *Client) buildReq(ctx context.Context, method, uri string, body interface{}) (*http.Request, error) {
	// prepare URI
	u, err := c.baseURL.Parse(uri)
	if err != nil {
		return nil, err
	}

	// prepare request body
	var respBody []byte
	if body != nil {
		respBody, err = json.MarshalIndent(body, "", "	")
		if err != nil {
			return nil, err
		}
	}

	// create request
	req, err := http.NewRequestWithContext(ctx, method, u.String(), bytes.NewBuffer(respBody))
	if err != nil {
		return nil, err
	}

	// set headers
	switch {
	case body == nil:
		req.Header.Set("Accept", contentType)
	default:
		req.Header.Set("Content-Type", contentType)
	}

	req.SetBasicAuth(c.apiKey, "")
	return req, nil
}

func (c *Client) performRequest(req *http.Request, expectedStatus int) ([]Job, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != expectedStatus {
		return nil, errors.New("status code doesn't match")
	}

	var jobs []Job
	err = json.NewDecoder(resp.Body).Decode(&jobs)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func (j *JobService) GetJobs(ctx context.Context) ([]Job, error) {
	jobsUri := "/jobs"

	req, err := j.client.buildReq(ctx, http.MethodGet, jobsUri, nil)
	if err != nil {
		return nil, err
	}

	resp, err := j.client.performRequest(req, 200)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
