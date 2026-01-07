package render

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"
)

// DecodeJSON decodes JSON from the request body
func DecodeJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

// JSON sends a JSON response with the given status and data
func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Error sends an error response with the given status and message
func Error(w http.ResponseWriter, status int, message string) {
	if status >= 500 {
		log.Printf("ERROR: %d %s\n%s", status, message, debug.Stack())
	} else {
		log.Printf("INFO: %d %s", status, message)
	}
	JSON(w, status, map[string]string{"error": message})
}
