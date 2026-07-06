package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aakashloyar/elevate/assessment_runner/config"
	"github.com/aakashloyar/elevate/assessment_runner/internal/application/ports/out"
)

type ClientFactory struct {
	services config.ServiceConfig
}

func NewClientFactory(services config.ServiceConfig) out.ClientFactory {
	return &ClientFactory{services: services}
}

func (f *ClientFactory) AssessmentClient() out.ServiceClient {
	return &ServiceClient{baseURL: f.services.AssessmentServiceURL, client: &http.Client{Timeout: 5 * time.Second}}
}

func (f *ClientFactory) ProblemClient() out.ServiceClient {
	return &ServiceClient{baseURL: f.services.ProblemServiceURL, client: &http.Client{Timeout: 5 * time.Second}}
}

func (f *ClientFactory) SubmissionClient() out.ServiceClient {
	return &ServiceClient{baseURL: f.services.SubmissionServiceURL, client: &http.Client{Timeout: 5 * time.Second}}
}

type ServiceClient struct {
	baseURL string
	client  *http.Client
}

func (c *ServiceClient) Get(ctx context.Context, path string) (*http.Response, error) {
	return c.do(ctx, http.MethodGet, path, nil)
}

func (c *ServiceClient) Post(ctx context.Context, path string, body any) (*http.Response, error) {
	return c.do(ctx, http.MethodPost, path, body)
}

func (c *ServiceClient) do(ctx context.Context, method, path string, body any) (*http.Response, error) {
	var payload *bytes.Buffer
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		payload = bytes.NewBuffer(data)
	}
	request, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, payload)
	if err != nil {
		return nil, err
	}
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= http.StatusBadRequest {
		resp.Body.Close()
		return nil, fmt.Errorf("request failed with status %d", resp.StatusCode)
	}
	return resp, nil
}
