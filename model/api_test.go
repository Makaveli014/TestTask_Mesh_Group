package model

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchSAPData(t *testing.T) {
	// Mock data to return
	mockData := []Segmentation{
		{AddressSapID: "sap123", AdrSegment: "seg1", SegmentID: 100},
		{AddressSapID: "sap456", AdrSegment: "seg2", SegmentID: 200},
	}
	mockJSON, _ := json.Marshal(mockData)

	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check User-Agent
		if r.Header.Get("User-Agent") != "test-agent" {
			t.Errorf("Expected User-Agent 'test-agent', got '%s'", r.Header.Get("User-Agent"))
		}
		// Check Basic Auth (Basic base64("user:pass") = Basic dXNlcjpwYXNz)
		if r.Header.Get("Authorization") != "Basic dXNlcjpwYXNz" {
			t.Errorf("Expected Basic Auth, got '%s'", r.Header.Get("Authorization"))
		}

		w.WriteHeader(http.StatusOK)
		w.Write(mockJSON)
	}))
	defer server.Close()

	cfg := APIConfig{
		URI:             server.URL,
		AuthLoginPwd:    "user:pass",
		UserAgent:       "test-agent",
		ImportBatchSize: 50,
	}

	data, err := FetchSAPData(server.Client(), cfg, 1)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(data) != 2 {
		t.Errorf("Expected 2 records, got %d", len(data))
	}

	if data[0].AddressSapID != "sap123" {
		t.Errorf("Expected AddressSapID 'sap123', got '%s'", data[0].AddressSapID)
	}
}

func TestFetchSAPData_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("")) // Empty body
	}))
	defer server.Close()

	cfg := APIConfig{
		URI:             server.URL,
		AuthLoginPwd:    "user:pass",
		UserAgent:       "test-agent",
		ImportBatchSize: 50,
	}

	data, err := FetchSAPData(server.Client(), cfg, 1)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(data) != 0 {
		t.Errorf("Expected 0 records for empty body, got %d", len(data))
	}
}
