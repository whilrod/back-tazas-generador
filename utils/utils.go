package utils

import (
	"encoding/json"
	"net/http"
	"strings"
)

// ParsePgArray convierte el string de PostgreSQL {a,b,c} a []string
func ParsePgArray(input string) []string {
	input = strings.Trim(input, "{}")
	if input == "" {
		return []string{}
	}
	parts := strings.Split(input, ",")
	for i := range parts {
		parts[i] = strings.Trim(parts[i], `"`)
	}
	return parts
}

// RespondJSON envía datos en formato JSON
func RespondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
