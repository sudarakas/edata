package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}

	return json.NewDecoder((r.Body)).Decode(payload)
}

func ValidatePayload(payload any) error {
	return Validate.Struct(payload)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) error {
	return WriteJSON(w, status, map[string]string{"error": err.Error()})
}

func WriteSuccess(w http.ResponseWriter, status int, data any) error {
	return WriteJSON(w, status, map[string]any{"data": data})
}

func WritePaginatedSuccess(w http.ResponseWriter, status int, data any, totalCount, limit, offset int) error {
	totalPages := 0
	if limit > 0 {
		totalPages = (totalCount + limit - 1) / limit
	}

	// Calculate the current page (1-based index)
	currentPage := offset/limit + 1

	// Determine if next page exists
	nextPage := currentPage + 1
	if nextPage > totalPages {
		nextPage = 0 // No next page
	}

	// Determine if previous page exists
	prevPage := currentPage - 1
	if prevPage < 1 {
		prevPage = 0 // No previous page
	}

	// Calculate last page (same as total_pages)
	lastPage := totalPages

	return WriteJSON(w, status, map[string]any{
		"data": map[string]any{
			"items":        data,
			"total_count":  totalCount,
			"limit":        limit,
			"offset":       offset,
			"current_page": currentPage,
			"total_pages":  totalPages,
			"next_page":    nextPage,
			"prev_page":    prevPage,
			"last_page":    lastPage,
		},
	})
}
