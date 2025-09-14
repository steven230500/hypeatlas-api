package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/steven230500/hypeatlas-api/domain/entities"
	in "github.com/steven230500/hypeatlas-api/modules/relay/domain/ports/in"
)

type HypeMapHandler struct{ svc in.Service }

func NewHypeMapHandler(s in.Service) *HypeMapHandler { return &HypeMapHandler{svc: s} }

func (h *HypeMapHandler) Register(r chi.Router) {
	r.Route("/hypemap", func(r chi.Router) {
		r.Get("/live", h.live)
		r.Get("/summary", h.summary)
	})
}

type HypeMapLiveResp struct {
	Items      []entities.HypeMapItem `json:"items"`
	NextOffset int                    `json:"next_offset"`
}
type HypeMapSummaryResp struct {
	Items      []entities.HypeMapSummaryItem `json:"items"`
	NextOffset int                           `json:"next_offset"`
}

// @Summary      HypeMap: co-streams en vivo (ranking)
// @Tags         relay
// @Param        game    query string false "val|lol"
// @Param        lang    query string false "es|en|fr|pt"
// @Param        limit   query int    false "1-100" minimum(1) maximum(100) default(50)
// @Param        offset  query int    false "paginación" default(0)
// @Produce      json
// @Success      200 {object} HypeMapLiveResp
// @Router       /v1/hypemap/live [get]
func (h *HypeMapHandler) live(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))

	items, err := h.svc.HypeMapLive(r.Context(), q.Get("game"), q.Get("lang"), limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"items":       items,
		"next_offset": offset + len(items),
	})
}

// @Summary      HypeMap: resumen por evento (agregado)
// @Tags         relay
// @Param        game    query string false "val|lol"
// @Param        lang    query string false "es|en|fr|pt"
// @Param        limit   query int    false "1-100" minimum(1) maximum(100) default(20)
// @Param        offset  query int    false "paginación" default(0)
// @Produce      json
// @Success      200 {object} HypeMapSummaryResp
// @Router       /v1/hypemap/summary [get]
func (h *HypeMapHandler) summary(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))

	items, err := h.svc.HypeMapSummary(r.Context(), q.Get("game"), q.Get("lang"), limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"items":       items,
		"next_offset": offset + len(items),
	})
}
