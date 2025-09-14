package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	signalrepo "github.com/steven230500/hypeatlas-api/modules/signal/infra/repository"
	"gorm.io/gorm"
)

type IngestHandler struct {
	db *gorm.DB
}

func NewIngest(db *gorm.DB) *IngestHandler { return &IngestHandler{db: db} }

func (h *IngestHandler) Register(r chi.Router) {
	r.Post("/comps:upsert", h.upsertComp)
}

type upsertCompReq struct {
	Game   string         `json:"game"`
	Region string         `json:"region"`
	League string         `json:"league"`
	Patch  string         `json:"patch"`
	Map    string         `json:"map"`
	Side   string         `json:"side"`
	Slots  map[string]any `json:"slots"`
	Pick   *float64       `json:"pick_rate"`
	Win    *float64       `json:"win_rate"`
	Delta  *float64       `json:"delta_win"`
}

// upsertComp godoc
// @Summary     Upsert de composición (ingesta)
// @Tags        ingest
// @Security    ApiKeyAuth
// @Accept      json
// @Param       body body   upsertCompReq true "payload"
// @Success     204 "no content"
// @Failure     400 {string} string "bad json"
// @Failure     500 {string} string "db error"
// @Router      /v1/ingest/signal/comps:upsert [post]
func (h *IngestHandler) upsertComp(w http.ResponseWriter, r *http.Request) {
	var req upsertCompReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	raw, _ := json.Marshal(req.Slots)
	repo := signalrepo.New(h.db)
	if pgRepo, ok := repo.(*signalrepo.Repo); ok {
		if err := pgRepo.UpsertComp(
			r.Context(),
			req.Game, req.Region, req.League, req.Patch, req.Map, req.Side,
			string(raw), req.Pick, req.Win, req.Delta,
		); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "invalid repo type", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
