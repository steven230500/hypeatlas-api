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
		r.Get("/live", h.live)
		r.Get("/summary", h.summary)
	})
}

// ====== Swagger response wrappers ======
type HypeMapLiveResp struct {
	Items      []relaypg.HypeMapItem `json:"items"`
	NextOffset int                   `json:"next_offset"`
}
type HypeMapSummaryResp struct {
	Items      []relaypg.HypeMapSummaryItem `json:"items"`
	NextOffset int                          `json:"next_offset"`
}

// live godoc
// @Summary      HypeMap: co-streams en vivo (ranking)
// @Tags         relay
// @Param        game    query string false "val|lol"
// @Param        lang    query string false "es|en|fr|pt"
// @Param        limit   query int    false "1-100" minimum(1) maximum(100) default(50)
// @Param        offset  query int    false "paginación" default(0)
// @Produce      json
// @Success      200 {object} HypeMapLiveResp
// @Failure      500 {string} string "db error"
// @Router       /v1/hypemap/live [get]
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

// summary godoc
// @Summary      HypeMap: resumen por evento (agregado)
// @Tags         relay
// @Param        game    query string false "val|lol"
// @Param        lang    query string false "es|en|fr|pt"
// @Param        limit   query int    false "1-100" minimum(1) maximum(100) default(20)
// @Param        offset  query int    false "paginación" default(0)
// @Produce      json
// @Success      200 {object} HypeMapSummaryResp
// @Failure      500 {string} string "db error"
// @Router       /v1/hypemap/summary [get]
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
