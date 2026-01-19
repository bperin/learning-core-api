package http

import (
	"net/http"
	"strconv"
)

// PaginationParams holds pagination parameters from request
type PaginationParams struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Offset   int `json:"offset"`
	Limit    int `json:"limit"`
}

// PaginatedResponse wraps paginated data with metadata
type PaginatedResponse[T any] struct {
	Data       []T              `json:"data"`
	Pagination PaginationMeta   `json:"pagination"`
}

// PaginationMeta contains pagination metadata
type PaginationMeta struct {
	Page         int   `json:"page"`
	PageSize     int   `json:"page_size"`
	Total        int64 `json:"total"`
	TotalPages   int   `json:"total_pages"`
	HasNext      bool  `json:"has_next"`
	HasPrevious  bool  `json:"has_previous"`
}

// GetPaginationParams extracts pagination parameters from request
func GetPaginationParams(r *http.Request) PaginationParams {
	page := getIntParam(r, "page", 1)
	pageSize := getIntParam(r, "page_size", 20)
	
	// Ensure minimum values
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100 // Maximum page size
	}
	
	offset := (page - 1) * pageSize
	
	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
		Offset:   offset,
		Limit:    pageSize,
	}
}

// NewPaginatedResponse creates a paginated response
func NewPaginatedResponse[T any](data []T, params PaginationParams, total int64) PaginatedResponse[T] {
	totalPages := int((total + int64(params.PageSize) - 1) / int64(params.PageSize))
	
	return PaginatedResponse[T]{
		Data: data,
		Pagination: PaginationMeta{
			Page:         params.Page,
			PageSize:     params.PageSize,
			Total:        total,
			TotalPages:   totalPages,
			HasNext:      params.Page < totalPages,
			HasPrevious:  params.Page > 1,
		},
	}
}

func getIntParam(r *http.Request, key string, defaultValue int) int {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}
	
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	
	return intValue
}
