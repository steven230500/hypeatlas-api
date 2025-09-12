package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/steven230500/hypeatlas-api/modules/relay/domain/entities"
	in "github.com/steven230500/hypeatlas-api/modules/relay/domain/ports/in"
)

type Handler struct{ svc in.Service }

func New(s in.Service) *Handler { return &Handler{svc: s} }

func (h *Handler) Register(r chi.Router) {
	r.Route("/relay", func(r chi.Router) {
		r.Get("/costreams", h.list)
	})
}

// ====== Swagger response wrappers ======
type CoStreamsResp struct {
	Items []entities.CoStream `json:"items"`
}

// list godoc
// @Summary      Listar co-streams en vivo de un evento
// @Tags         relay
// @Param        event_id query string true  "Slug o ID del evento (ej. vct-emea-final)"
// @Param        lang     query string false "Filtro por idioma (ej. es)"
// @Produce      json
// @Success      200 {object} CoStreamsResp
// @Failure      400 {string} string "event_id required"
// @Router       /v1/relay/costreams [get]
func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	eventID := r.URL.Query().Get("event_id")
	lang := r.URL.Query().Get("lang")
	items, err := h.svc.ListLiveCoStreams(r.Context(), eventID, lang)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"items": items})
}
