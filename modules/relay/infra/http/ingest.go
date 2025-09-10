package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	relaypg "github.com/steven230500/hypeatlas-api/modules/relay/infra/repository/postgres"
)

type IngestHandler struct {
	pool *pgxpool.Pool
}

func NewIngest(pool *pgxpool.Pool) *IngestHandler { return &IngestHandler{pool: pool} }

func (h *IngestHandler) Register(r chi.Router) {
	r.Post("/costreams:upsert", h.upsertCoStream)
}

type upsertCoStreamReq struct {
	EventSlug  string  `json:"event_slug"`
	EventTitle string  `json:"event_title"`
	Game       string  `json:"game"`   // "val"|"lol"
	League     string  `json:"league"` // opcional
	StartsAt   *string `json:"starts_at"`

	Platform string `json:"platform"` // "twitch"|"youtube"
	Handle   string `json:"handle"`
	URL      string `json:"url"`
	Lang     string `json:"lang"`
	Country  string `json:"country"`
	Verified bool   `json:"verified"`
	Viewers  int    `json:"viewers"`
	IsLive   bool   `json:"is_live"`
}

func (h *IngestHandler) upsertCoStream(w http.ResponseWriter, r *http.Request) {
	var req upsertCoStreamReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	repo := relaypg.NewRaw(h.pool)
	if err := repo.UpsertCoStream(
		r.Context(),
		req.EventSlug, req.EventTitle, req.Game, req.League, req.StartsAt,
		req.Platform, req.Handle, req.URL, req.Lang, req.Country, req.Verified, req.Viewers, req.IsLive,
	); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
