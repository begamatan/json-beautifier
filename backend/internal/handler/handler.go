package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/begamatan/json-beautifier/backend/internal/formatter"
)

const maxBodySize = 5 * 1024 * 1024 // 5 MB

// apiError is the consistent error response schema.
type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// beautifyRequest is the request body for /api/v1/beautify.
type beautifyRequest struct {
	JSON   string `json:"json"`
	Indent int    `json:"indent"`
}

// jsonRequest is the request body for /api/v1/minify and /api/v1/validate.
type jsonRequest struct {
	JSON string `json:"json"`
}

// outputResponse is the response body for beautify and minify.
type outputResponse struct {
	Result string `json:"result"`
}

// validateResponse is the response body for validate.
type validateResponse struct {
	Valid   bool   `json:"valid"`
	Message string `json:"message"`
}

// Health handles GET /api/v1/health.
func Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// Beautify handles POST /api/v1/beautify.
func Beautify(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	var req beautifyRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid request body", err.Error())
		return
	}

	indent := req.Indent
	if indent == 0 {
		indent = 2
	}

	result, err := formatter.Beautify([]byte(req.JSON), indent)
	if err != nil {
		if errors.Is(err, formatter.ErrInvalidJSON) {
			writeError(w, http.StatusUnprocessableEntity, "INVALID_JSON", "input is not valid JSON", "")
			return
		}
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error(), "")
		return
	}

	writeJSON(w, http.StatusOK, outputResponse{Result: string(result)})
}

// Minify handles POST /api/v1/minify.
func Minify(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	var req jsonRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid request body", err.Error())
		return
	}

	result, err := formatter.Minify([]byte(req.JSON))
	if err != nil {
		if errors.Is(err, formatter.ErrInvalidJSON) {
			writeError(w, http.StatusUnprocessableEntity, "INVALID_JSON", "input is not valid JSON", "")
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "unexpected error", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, outputResponse{Result: string(result)})
}

// Validate handles POST /api/v1/validate.
func Validate(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	var req jsonRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid request body", err.Error())
		return
	}

	if err := formatter.Validate([]byte(req.JSON)); err != nil {
		writeJSON(w, http.StatusOK, validateResponse{Valid: false, Message: "input is not valid JSON"})
		return
	}

	writeJSON(w, http.StatusOK, validateResponse{Valid: true, Message: "JSON is valid"})
}

// decodeJSON reads and decodes JSON from the request body.
func decodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}

// writeJSON encodes v as JSON and writes it to w with the given status code.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// writeError writes a consistent JSON error response.
func writeError(w http.ResponseWriter, status int, code, message, details string) {
	writeJSON(w, status, apiError{Code: code, Message: message, Details: details})
}
