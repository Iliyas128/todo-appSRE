package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateTodo(t *testing.T) {
	// Initialize router
	router := setupRouter()

	// Test cases
	tests := []struct {
		name       string
		payload    map[string]interface{}
		wantStatus int
	}{
		{
			name: "Valid todo",
			payload: map[string]interface{}{
				"title":       "Test Todo",
				"description": "Test Description",
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "Invalid todo - missing title",
			payload: map[string]interface{}{
				"description": "Test Description",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			jsonData, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("POST", "/todos", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			rr := httptest.NewRecorder()

			// Serve request
			router.ServeHTTP(rr, req)

			// Check status code
			if rr.Code != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					rr.Code, tt.wantStatus)
			}
		})
	}
}

func TestGetTodos(t *testing.T) {
	// Initialize router
	router := setupRouter()

	// Create test request
	req, _ := http.NewRequest("GET", "/todos", nil)
	rr := httptest.NewRecorder()

	// Serve request
	router.ServeHTTP(rr, req)

	// Check status code
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}

	// Check response body
	var response []map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}
}

func TestMetricsEndpoint(t *testing.T) {
	// Initialize router
	router := setupRouter()

	// Create test request
	req, _ := http.NewRequest("GET", "/metrics", nil)
	rr := httptest.NewRecorder()

	// Serve request
	router.ServeHTTP(rr, req)

	// Check status code
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}

	// Check if response contains expected metrics
	body := rr.Body.String()
	expectedMetrics := []string{
		"http_requests_total",
		"http_request_duration_seconds",
		"http_requests_in_flight",
	}

	for _, metric := range expectedMetrics {
		if !bytes.Contains(rr.Body.Bytes(), []byte(metric)) {
			t.Errorf("Response body does not contain expected metric: %s", metric)
		}
	}
}
