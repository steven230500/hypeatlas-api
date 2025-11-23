package http

import (
	"net/http"
	"strconv"
)

// PaginationParams contiene los parámetros de paginación parseados
type PaginationParams struct {
	Limit  int
	Offset int
}

// PaginationMeta contiene metadata sobre la paginación
type PaginationMeta struct {
	Limit   int  `json:"limit"`
	Offset  int  `json:"offset"`
	Total   int  `json:"total"`
	HasMore bool `json:"hasMore"`
}

// PaginatedResponse es una respuesta genérica con paginación
type PaginatedResponse struct {
	Items      interface{}    `json:"items"`
	Pagination PaginationMeta `json:"pagination"`
}

// ParsePaginationParams extrae y valida los parámetros de paginación de la request
func ParsePaginationParams(r *http.Request) PaginationParams {
	const (
		defaultLimit = 20
		maxLimit     = 100
	)

	limit := defaultLimit
	offset := 0

	// Parse limit
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			if parsedLimit > 0 && parsedLimit <= maxLimit {
				limit = parsedLimit
			} else if parsedLimit > maxLimit {
				limit = maxLimit
			}
		}
	}

	// Parse offset
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	return PaginationParams{
		Limit:  limit,
		Offset: offset,
	}
}

// ApplyPagination aplica paginación a un slice y retorna el subset
func ApplyPagination[T any](items []T, params PaginationParams) ([]T, PaginationMeta) {
	total := len(items)
	start := params.Offset
	end := start + params.Limit

	// Validar bounds
	if start >= total {
		return []T{}, PaginationMeta{
			Limit:   params.Limit,
			Offset:  params.Offset,
			Total:   total,
			HasMore: false,
		}
	}

	if end > total {
		end = total
	}

	return items[start:end], PaginationMeta{
		Limit:   params.Limit,
		Offset:  params.Offset,
		Total:   total,
		HasMore: end < total,
	}
}
