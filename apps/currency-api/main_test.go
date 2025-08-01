// currency-api/main_test.go - Clean working version

package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	// Create test app
	app := createApp()

	// Create test request
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	// Test status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Test response body contains expected values
	body := w.Body.String()
	expectedStrings := []string{"healthy", "currency-exchange-api", "1.24"}

	for _, expected := range expectedStrings {
		if !strings.Contains(body, expected) {
			t.Errorf("Expected response to contain '%s', got: %s", expected, body)
		}
	}
}

func TestRatesEndpoint(t *testing.T) {
	app := createApp()

	t.Run("missing currencies parameter returns 400", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/rates", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}

		body := w.Body.String()
		if !strings.Contains(body, "currencies parameter is required") {
			t.Errorf("Expected error message about missing currencies parameter, got: %s", body)
		}
	})

	t.Run("with currencies parameter returns 200", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/rates?currencies=USD,EUR", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		body := w.Body.String()
		if !strings.Contains(body, "USD,EUR") {
			t.Errorf("Expected response to contain 'USD,EUR', got: %s", body)
		}
	})
}

func TestExchangeEndpoint(t *testing.T) {
	app := createApp()

	t.Run("missing parameters returns 400", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/exchange", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}

		body := w.Body.String()
		if !strings.Contains(body, "required") {
			t.Errorf("Expected error message about required parameters, got: %s", body)
		}
	})

	t.Run("with all parameters returns 200", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/exchange?from=WBTC&to=USDT&amount=1.0", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		body := w.Body.String()
		expectedParams := []string{"WBTC", "USDT", "1.0"}

		for _, param := range expectedParams {
			if !strings.Contains(body, param) {
				t.Errorf("Expected response to contain '%s', got: %s", param, body)
			}
		}
	})
}

// Benchmark test to show off
func BenchmarkHealthEndpoint(t *testing.B) {
	app := createApp()
	req, _ := http.NewRequest("GET", "/health", nil)

	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}

// Table-driven test example (idiomatic Go)
func TestCorsHeaders(t *testing.T) {
	app := createApp()

	tests := []struct {
		name           string
		method         string
		expectedHeader string
		expectedValue  string
	}{
		{
			name:           "OPTIONS request sets CORS headers",
			method:         "OPTIONS",
			expectedHeader: "Access-Control-Allow-Origin",
			expectedValue:  "*",
		},
		{
			name:           "GET request sets CORS headers",
			method:         "GET",
			expectedHeader: "Access-Control-Allow-Origin",
			expectedValue:  "*",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, "/health", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			w := httptest.NewRecorder()
			app.ServeHTTP(w, req)

			headerValue := w.Header().Get(tt.expectedHeader)
			if headerValue != tt.expectedValue {
				t.Errorf("Expected header %s to be %s, got %s",
					tt.expectedHeader, tt.expectedValue, headerValue)
			}
		})
	}
}