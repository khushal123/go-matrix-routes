package main

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	// Setup routes once for all tests
	setupRoutes()

	// Run all tests
	m.Run()
}

// Helper function to create multipart form request with file
func createTestRequest(t *testing.T, content string, filename string, path string) *http.Request {
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	if filename != "" {
		part, err := writer.CreateFormFile("file", filename)
		if err != nil {
			t.Fatal(err)
		}
		part.Write([]byte(content))
	}

	writer.Close()

	req := httptest.NewRequest("POST", path, &b)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

func TestMatrixEndpoints(t *testing.T) {
	// Initialize routes

	tests := []struct {
		name         string
		path         string
		method       string
		fileContent  string
		filename     string
		expectedCode int
		expectedBody string
	}{
		// Success cases
		{
			name:         "Echo Success",
			path:         "/echo",
			method:       "POST",
			fileContent:  "1,2\n3,4",
			filename:     "valid.csv",
			expectedCode: http.StatusOK,
			expectedBody: "1,2\n3,4",
		},
		{
			name:         "Invert Success",
			path:         "/invert",
			method:       "POST",
			fileContent:  "1,2\n3,4",
			filename:     "valid.csv",
			expectedCode: http.StatusOK,
			expectedBody: "1,3\n2,4",
		},
		{
			name:         "Flatten Success",
			path:         "/flatten",
			method:       "POST",
			fileContent:  "1,2\n3,4",
			filename:     "valid.csv",
			expectedCode: http.StatusOK,
			expectedBody: "1,2,3,4",
		},
		{
			name:         "Sum Success",
			path:         "/sum",
			method:       "POST",
			fileContent:  "1,2\n3,4",
			filename:     "valid.csv",
			expectedCode: http.StatusOK,
			expectedBody: "10",
		},
		{
			name:         "Multiply Success",
			path:         "/multiply",
			method:       "POST",
			fileContent:  "1,2\n3,4",
			filename:     "valid.csv",
			expectedCode: http.StatusOK,
			expectedBody: "24",
		},

		// Failure cases - Wrong HTTP Method
		{
			name:         "Echo Wrong Method",
			path:         "/echo",
			method:       "GET",
			fileContent:  "",
			filename:     "",
			expectedCode: http.StatusMethodNotAllowed,
		},

		// Failure cases - No File
		{
			name:         "Echo No File",
			path:         "/echo",
			method:       "POST",
			fileContent:  "",
			filename:     "",
			expectedCode: http.StatusBadRequest,
		},

		// Failure cases - Non-Square Matrix
		{
			name:         "Echo Non-Square Matrix",
			path:         "/echo",
			method:       "POST",
			fileContent:  "1,2\n3,4,5",
			filename:     "invalid.csv",
			expectedCode: http.StatusInternalServerError,
		},

		// Failure cases - Invalid Numbers
		{
			name:         "Echo Invalid Numbers",
			path:         "/echo",
			method:       "POST",
			fileContent:  "1,a\n3,4",
			filename:     "invalid.csv",
			expectedCode: http.StatusInternalServerError,
		},

		// Failure cases - Empty Matrix
		{
			name:         "Echo Empty Matrix",
			path:         "/echo",
			method:       "POST",
			fileContent:  "",
			filename:     "empty.csv",
			expectedCode: http.StatusInternalServerError,
		},

		// Failure cases - Wrong File Type
		{
			name:         "Echo Wrong File Type",
			path:         "/echo",
			method:       "POST",
			fileContent:  "1,2\n3,4",
			filename:     "file.txt",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.method == "GET" {
				req = httptest.NewRequest(tt.method, tt.path, nil)
			} else {
				req = createTestRequest(t, tt.fileContent, tt.filename, tt.path)
			}

			rr := httptest.NewRecorder()
			http.DefaultServeMux.ServeHTTP(rr, req)

			// Check status code
			if rr.Code != tt.expectedCode {
				t.Errorf("%s: handler returned wrong status code: got %v want %v",
					tt.name, rr.Code, tt.expectedCode)
			}

			// Check response body contains expected string
			if !strings.Contains(rr.Body.String(), tt.expectedBody) {
				t.Errorf("%s: handler returned unexpected body: got %v want %v",
					tt.name, rr.Body.String(), tt.expectedBody)
			}
		})
	}
}

// TestConcurrentRequests tests how the endpoints handle concurrent requests
func TestConcurrentRequests(t *testing.T) {
	paths := []string{"/echo", "/invert", "/flatten", "/sum", "/multiply"}
	for _, path := range paths {
		t.Run("Concurrent"+path, func(t *testing.T) {
			// Create a channel to wait for all requests to complete
			done := make(chan bool)

			// Launch 5 concurrent requests
			for i := 0; i < 5; i++ {
				go func() {
					req := createTestRequest(t, "1,2\n3,4", "valid.csv", path)
					rr := httptest.NewRecorder()
					http.DefaultServeMux.ServeHTTP(rr, req)

					if rr.Code != http.StatusOK {
						t.Errorf("concurrent request failed with status: %d", rr.Code)
					}
					done <- true
				}()
			}

			// Wait for all requests to complete
			for i := 0; i < 5; i++ {
				<-done
			}
		})
	}
}
