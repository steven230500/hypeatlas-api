package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	out "github.com/steven230500/hypeatlas-api/modules/relay/domain/ports/out"
)

type IngestHandler struct{ repo out.Repository }

func NewIngest(repo out.Repository) *IngestHandler { return &IngestHandler{repo: repo} }

func (h *IngestHandler) Register(r chi.Router) {
	r.Post("/costreams:upsert", h.upsertCoStream)
}

type upsertCoStreamReq struct {
	EventSlug  string  `json:"event_slug"`
	EventTitle string  `json:"event_title"`
	Game       string  `json:"game"`
	League     string  `json:"league"`
	StartsAt   *string `json:"starts_at"`

	Platform string `json:"platform"` // twitch|youtube
	Handle   string `json:"handle"`
	URL      string `json:"url"`
	Lang     string `json:"lang"`
	Country  string `json:"country"`
	Verified bool   `json:"verified"`
	Viewers  int    `json:"viewers"`
	IsLive   bool   `json:"is_live"`
}

// @Summary     Ingest: upsert co-stream
// @Tags        ingest
// @Security    ApiKeyAuth
// @Accept      json
// @Param       body body   upsertCoStreamReq true "payload"
// @Success     204 "no content"
// @Failure     400 {string} string "bad json"
// @Failure     500 {string} string "db error"
// @Router      /v1/ingest/relay/costreams:upsert [post]
func (h *IngestHandler) upsertCoStream(w http.ResponseWriter, r *http.Request) {
	var req upsertCoStreamReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	if err := h.repo.UpsertCoStream(
		r.Context(),
		req.EventSlug, req.EventTitle, req.Game, req.League, req.StartsAt,
		req.Platform, req.Handle, req.URL, req.Lang, req.Country, req.Verified, req.Viewers, req.IsLive,
	); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
