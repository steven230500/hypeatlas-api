package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	in "github.com/steven230500/hypeatlas-api/modules/signal/domain/ports/in"
)

type Handler struct{ svc in.Service }

func New(s in.Service) *Handler { return &Handler{svc: s} }

func (h *Handler) Register(r chi.Router) {
	r.Route("/signal", func(r chi.Router) {
		r.Get("/patches", h.patches)
		r.Get("/changes", h.changes)
	})
}

func (h *Handler) patches(w http.ResponseWriter, r *http.Request) {
	game := r.URL.Query().Get("game")
	items, err := h.svc.ListPatches(r.Context(), game)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"items": items})
}

func (h *Handler) changes(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	items, err := h.svc.ListChanges(r.Context(), q.Get("game"), q.Get("version"), q.Get("type"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"items": items})
}
