package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	in "github.com/steven230500/hypeatlas-api/modules/relay/domain/ports/in"
)

type Handler struct{ svc in.Service }

func New(s in.Service) *Handler { return &Handler{svc: s} }

func (h *Handler) Register(r chi.Router) {
	r.Route("/relay", func(r chi.Router) {
		r.Get("/costreams", h.list)
	})
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	eventID := r.URL.Query().Get("event_id")
	lang := r.URL.Query().Get("lang")
	items, err := h.svc.ListLiveCoStreams(r.Context(), eventID, lang)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"items": items})
}
