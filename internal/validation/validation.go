package validation

import (
	"encoding/json"
	"net/http"
)

// ErrorDetail memuat struktur error response terstandar sesuai SSOT §70.
type ErrorDetail struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

// ErrorResponse memuat wrapper error JSON.
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// WriteError mengirim response error JSON terstandar ke client.
func WriteError(w http.ResponseWriter, statusCode int, code, message string, fields map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
			Fields:  fields,
		},
	})
}
