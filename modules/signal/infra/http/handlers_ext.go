package http

import (
	"encoding/json"
	"net/http"
)

// @Summary Listar campeones
// @Tags signal
// @Param version query string false "Versión del juego (ej. 14.14)"
// @Produce json
// @Success 200 {object} map[string]any
// @Failure 400 {string} string "error"
// @Router /v1/signal/champions [get]
func (h *Handler) champions(w http.ResponseWriter, r *http.Request) {
	version := r.URL.Query().Get("version")
	items, err := h.svc.ListChampions(r.Context(), version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"items": items})
}

// @Summary Listar regiones disponibles
// @Tags signal
// @Produce json
// @Success 200 {object} map[string]any
// @Failure 500 {string} string "error"
// @Router /v1/signal/regions [get]
func (h *Handler) regions(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.ListRegions(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"items": items})
}
