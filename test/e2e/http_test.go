package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

type HTTPClient struct {
	baseURL string
	client  *http.Client
	t       *testing.T
}

func NewHTTPClient(t *testing.T, baseURL string) *HTTPClient {
	return &HTTPClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		t: t,
	}
}

type Request struct {
	Method  string
	Path    string
	Body    interface{}
	Headers map[string]string
}

type Response struct {
	StatusCode int
	Body       []byte
	Headers    http.Header
}

func (hc *HTTPClient) Do(req *Request) (*Response, error) {
	var body io.Reader
	if req.Body != nil {
		jsonBody, err := json.Marshal(req.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		body = bytes.NewReader(jsonBody)
	}

	httpReq, err := http.NewRequest(req.Method, fmt.Sprintf("%s%s", hc.baseURL, req.Path), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	if req.Headers != nil {
		for key, value := range req.Headers {
			httpReq.Header.Set(key, value)
		}
	}

	resp, err := hc.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Body:       respBody,
		Headers:    resp.Header,
	}, nil
}

func (hc *HTTPClient) Post(path string, body interface{}) (*Response, error) {
	return hc.Do(&Request{
		Method: http.MethodPost,
		Path:   path,
		Body:   body,
	})
}

func (hc *HTTPClient) Get(path string) (*Response, error) {
	return hc.Do(&Request{
		Method: http.MethodGet,
		Path:   path,
	})
}

func (hc *HTTPClient) Put(path string, body interface{}) (*Response, error) {
	return hc.Do(&Request{
		Method: http.MethodPut,
		Path:   path,
		Body:   body,
	})
}

func (hc *HTTPClient) Delete(path string) (*Response, error) {
	return hc.Do(&Request{
		Method: http.MethodDelete,
		Path:   path,
	})
}

func (r *Response) UnmarshalJSON(target interface{}) error {
	if err := json.Unmarshal(r.Body, target); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return nil
}

func (r *Response) AssertStatusCode(t *testing.T, expected int) {
	if r.StatusCode != expected {
		t.Errorf("Expected status code %d, got %d. Body: %s", expected, r.StatusCode, string(r.Body))
	}
}

type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

type SuccessResponse struct {
	Data interface{} `json:"data"`
}

type UserResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UsersListResponse struct {
	Data []UserResponse `json:"data"`
}
