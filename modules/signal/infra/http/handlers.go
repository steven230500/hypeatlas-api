package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	in "github.com/steven230500/hypeatlas-api/modules/signal/domain/ports/in"

	"github.com/steven230500/hypeatlas-api/modules/signal/domain/entities"
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

// ====== Wrappers de respuesta para Swagger ======
type PatchesResp struct {
	Items []entities.Patch `json:"items"`
}
type ChangesResp struct {
	Items []entities.Change `json:"items"`
}
type LeaguesResp struct {
	Items []entities.League `json:"items"`
}
type CompsResp struct {
	Items []entities.Comp `json:"items"`
}

// patches godoc
// @Summary      Listar parches por juego
// @Tags         signal
// @Param        game   query string true  "lol | val" enums(lol,val)
// @Produce      json
// @Success      200 {object} PatchesResp
// @Failure      400 {string} string "game required"
// @Router       /v1/signal/patches [get]
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

// changes godoc
// @Summary      Listar cambios de un parche
// @Tags         signal
// @Param        game     query string true  "lol | val" enums(lol,val)
// @Param        version  query string true  "Ej: 14.14 | 9.15"
// @Param        type     query string false "agent|champion|item|weapon|map" enums(agent,champion,item,weapon,map)
// @Produce      json
// @Success      200 {object} ChangesResp
// @Failure      400 {string} string "game and version required"
// @Router       /v1/signal/changes [get]
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

// leagues godoc
// @Summary      Listar ligas por juego y región
// @Tags         signal
// @Param        game   query string true  "lol | val" enums(lol,val)
// @Param        region query string false "EMEA | APAC | AMERICAS" enums(EMEA,APAC,AMERICAS)
// @Produce      json
// @Success      200 {object} LeaguesResp
// @Failure      400 {string} string "game required"
// @Router       /v1/signal/leagues [get]
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

// comps godoc
// @Summary      Listar composiciones (LoL/VAL) filtradas
// @Tags         signal
// @Param        game   query string true  "lol | val" enums(lol,val)
// @Param        region query string true  "EMEA"
// @Param        league query string false "LEC"
// @Param        patch  query string true  "14.14 | 9.15"
// @Param        map    query string false "Ascent (solo VAL)"
// @Param        side   query string false "LoL: blue/red | VAL: attack/defense" enums(blue,red,attack,defense)
// @Param        limit  query int    false "1-100" minimum(1) maximum(200) default(50)
// @Produce      json
// @Success      200 {object} CompsResp
// @Failure      400 {string} string "game, region and patch required"
// @Router       /v1/signal/comps [get]
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
