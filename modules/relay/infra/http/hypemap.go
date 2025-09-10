package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	relaypg "github.com/steven230500/hypeatlas-api/modules/relay/infra/repository/postgres"
)

type HypeMapHandler struct {
	Repo *relaypg.Repo
}

func NewHypeMapHandler(repo *relaypg.Repo) *HypeMapHandler { return &HypeMapHandler{Repo: repo} }

func (h *HypeMapHandler) Register(r chi.Router) {
	r.Route("/hypemap", func(r chi.Router) {
		r.Get("/live", h.live)       // si ya tienes el Live, queda igual
		r.Get("/summary", h.summary) // <-- NUEVO
	})
}

// GET /v1/hypemap/live?game=val&lang=es&limit=50&offset=0
func (h *HypeMapHandler) live(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	game := q.Get("game")
	lang := q.Get("lang")
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))

	items, err := h.Repo.HypeMapLive(r.Context(), game, lang, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := map[string]any{
		"items":       items,
		"next_offset": offset + len(items),
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// GET /v1/hypemap/summary?game=val&lang=es&limit=20&offset=0
func (h *HypeMapHandler) summary(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	game := q.Get("game")
	lang := q.Get("lang")
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))

	items, err := h.Repo.HypeMapSummary(r.Context(), game, lang, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := map[string]any{
		"items":       items,
		"next_offset": offset + len(items),
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
