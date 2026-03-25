package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/begamatan/json-beautifier/backend/internal/handler"
)

// --- helpers ---

func doRequest(t *testing.T, h http.HandlerFunc, method, target, body string) *httptest.ResponseRecorder {
	t.Helper()
	var reqBody *strings.Reader
	if body != "" {
		reqBody = strings.NewReader(body)
	} else {
		reqBody = strings.NewReader("")
	}
	req := httptest.NewRequest(method, target, reqBody)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h(rr, req)
	return rr
}

func decodeBody(t *testing.T, rr *httptest.ResponseRecorder, dst any) {
	t.Helper()
	if err := json.NewDecoder(rr.Body).Decode(dst); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
}

// --- health ---

func TestHealth(t *testing.T) {
	rr := doRequest(t, handler.Health, http.MethodGet, "/api/v1/health", "")
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
	var resp map[string]string
	decodeBody(t, rr, &resp)
	if resp["status"] != "ok" {
		t.Errorf("expected status=ok, got %q", resp["status"])
	}
}

// --- beautify ---

func TestBeautify(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantResult string
		wantCode   string
	}{
		{
			name:       "valid json default indent",
			body:       `{"json":"{\"a\":1}","indent":2}`,
			wantStatus: http.StatusOK,
			wantResult: "{\n  \"a\": 1\n}",
		},
		{
			name:       "valid json 4-space indent",
			body:       `{"json":"{\"a\":1}","indent":4}`,
			wantStatus: http.StatusOK,
			wantResult: "{\n    \"a\": 1\n}",
		},
		{
			name:       "invalid json input",
			body:       `{"json":"{bad}","indent":2}`,
			wantStatus: http.StatusUnprocessableEntity,
			wantCode:   "INVALID_JSON",
		},
		{
			name:       "invalid indent value",
			body:       `{"json":"{\"a\":1}","indent":3}`,
			wantStatus: http.StatusBadRequest,
			wantCode:   "BAD_REQUEST",
		},
		{
			name:       "malformed request body",
			body:       `not-json`,
			wantStatus: http.StatusBadRequest,
			wantCode:   "BAD_REQUEST",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rr := doRequest(t, handler.Beautify, http.MethodPost, "/api/v1/beautify", tc.body)
			if rr.Code != tc.wantStatus {
				t.Errorf("status: got %d, want %d", rr.Code, tc.wantStatus)
			}
			if tc.wantResult != "" {
				var resp map[string]string
				decodeBody(t, rr, &resp)
				if resp["result"] != tc.wantResult {
					t.Errorf("result: got %q, want %q", resp["result"], tc.wantResult)
				}
			}
			if tc.wantCode != "" {
				var resp map[string]string
				decodeBody(t, rr, &resp)
				if resp["code"] != tc.wantCode {
					t.Errorf("code: got %q, want %q", resp["code"], tc.wantCode)
				}
			}
		})
	}
}

// --- minify ---

func TestMinify(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantResult string
		wantCode   string
	}{
		{
			name:       "valid pretty json",
			body:       `{"json":"{\n  \"a\": 1\n}"}`,
			wantStatus: http.StatusOK,
		},
		{
			name:       "valid compact json",
			body:       `{"json":"{\"a\":1,\"b\":2}"}`,
			wantStatus: http.StatusOK,
			wantResult: `{"a":1,"b":2}`,
		},
		{
			name:       "invalid json",
			body:       `{"json":"{bad}"}`,
			wantStatus: http.StatusUnprocessableEntity,
			wantCode:   "INVALID_JSON",
		},
		{
			name:       "malformed body",
			body:       `oops`,
			wantStatus: http.StatusBadRequest,
			wantCode:   "BAD_REQUEST",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rr := doRequest(t, handler.Minify, http.MethodPost, "/api/v1/minify", tc.body)
			if rr.Code != tc.wantStatus {
				t.Errorf("status: got %d, want %d", rr.Code, tc.wantStatus)
			}
			if tc.wantResult != "" {
				var resp map[string]string
				decodeBody(t, rr, &resp)
				if resp["result"] != tc.wantResult {
					t.Errorf("result: got %q, want %q", resp["result"], tc.wantResult)
				}
			}
			if tc.wantCode != "" {
				var resp map[string]string
				decodeBody(t, rr, &resp)
				if resp["code"] != tc.wantCode {
					t.Errorf("code: got %q, want %q", resp["code"], tc.wantCode)
				}
			}
		})
	}
}

// --- validate ---

func TestValidate(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantValid  bool
		wantCode   string
	}{
		{
			name:       "valid json",
			body:       `{"json":"{\"a\":1}"}`,
			wantStatus: http.StatusOK,
			wantValid:  true,
		},
		{
			name:       "invalid json",
			body:       `{"json":"{bad}"}`,
			wantStatus: http.StatusOK,
			wantValid:  false,
		},
		{
			name:       "malformed body",
			body:       `oops`,
			wantStatus: http.StatusBadRequest,
			wantCode:   "BAD_REQUEST",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rr := doRequest(t, handler.Validate, http.MethodPost, "/api/v1/validate", tc.body)
			if rr.Code != tc.wantStatus {
				t.Errorf("status: got %d, want %d", rr.Code, tc.wantStatus)
			}
			if tc.wantCode == "" && tc.wantStatus == http.StatusOK {
				var resp map[string]any
				decodeBody(t, rr, &resp)
				if valid, ok := resp["valid"].(bool); !ok || valid != tc.wantValid {
					t.Errorf("valid: got %v, want %v", resp["valid"], tc.wantValid)
				}
			}
		})
	}
}

// --- size limit ---

func TestSizeLimit(t *testing.T) {
	bigPayload := bytes.Repeat([]byte("a"), 6*1024*1024)
	body := `{"json":"` + string(bigPayload) + `"}`
	rr := doRequest(t, handler.Beautify, http.MethodPost, "/api/v1/beautify", body)
	if rr.Code == http.StatusOK {
		t.Error("expected non-200 for oversized request")
	}
}
