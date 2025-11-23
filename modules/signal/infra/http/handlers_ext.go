package http

import (
	"encoding/json"
	"net/http"

	sharedhttp "github.com/steven230500/hypeatlas-api/shared/http"
)

// @Summary Listar campeones
// @Tags signal
// @Param version query string false "Versión del juego (ej. 14.14.1)"
// @Param limit query int false "Number of items per page (default: 20, max: 100)"
// @Param offset query int false "Number of items to skip (default: 0)"
// @Produce json
// @Success 200 {object} sharedhttp.PaginatedResponse
// @Failure 400 {string} string "error"
// @Router /v1/signal/champions [get]
func (h *Handler) champions(w http.ResponseWriter, r *http.Request) {
	version := r.URL.Query().Get("version")
	items, err := h.svc.ListChampions(r.Context(), version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Aplicar paginación
	params := sharedhttp.ParsePaginationParams(r)
	paginatedItems, meta := sharedhttp.ApplyPagination(items, params)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(sharedhttp.PaginatedResponse{
		Items:      paginatedItems,
		Pagination: meta,
	})
}

// @Summary Listar regiones disponibles
// @Tags signal
// @Param limit query int false "Number of items per page (default: 20, max: 100)"
// @Param offset query int false "Number of items to skip (default: 0)"
// @Produce json
// @Success 200 {object} sharedhttp.PaginatedResponse
// @Failure 500 {string} string "error"
// @Router /v1/signal/regions [get]
func (h *Handler) regions(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.ListRegions(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Aplicar paginación
	params := sharedhttp.ParsePaginationParams(r)
	paginatedItems, meta := sharedhttp.ApplyPagination(items, params)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(sharedhttp.PaginatedResponse{
		Items:      paginatedItems,
		Pagination: meta,
	})
}
