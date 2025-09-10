package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	in "github.com/steven230500/hypeatlas-api/modules/signal/domain/ports/in"
)

type Handler struct{ svc in.Service }

func New(s in.Service) *Handler { return &Handler{svc: s} }

func (h *Handler) Register(r chi.Router) {
	r.Route("/signal", func(r chi.Router) {
		r.Get("/patches", h.patches)
		r.Get("/changes", h.changes)
		r.Get("/leagues", h.leagues)
		r.Get("/comps", h.comps)
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
	_ = json.NewEncoder(w).Encode(map[string]any{"items": items})
}

func (h *Handler) changes(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	items, err := h.svc.ListChanges(r.Context(), q.Get("game"), q.Get("version"), q.Get("type"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"items": items})
}

func (h *Handler) leagues(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	items, err := h.svc.ListLeagues(r.Context(), q.Get("game"), q.Get("region"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"items": items})
}

func (h *Handler) comps(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit, _ := strconv.Atoi(q.Get("limit"))
	items, err := h.svc.ListComps(
		r.Context(),
		q.Get("game"),
		q.Get("region"),
		q.Get("league"),
		q.Get("patch"),
		q.Get("map"),
		q.Get("side"),
		limit,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"items": items})
}
