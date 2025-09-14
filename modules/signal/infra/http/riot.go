package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/steven230500/hypeatlas-api/domain/entities"
	in "github.com/steven230500/hypeatlas-api/modules/signal/domain/ports/in"
	"github.com/steven230500/hypeatlas-api/modules/signal/domain/service"
	"github.com/steven230500/hypeatlas-api/providers/riot"
)

type RiotHandler struct {
	riotSvc     *riot.Service
	sigSvc      in.Service
	metaGameSvc *service.MetaGameService
}

func NewRiotHandler(riotSvc *riot.Service, sigSvc in.Service, metaGameSvc *service.MetaGameService) *RiotHandler {
	return &RiotHandler{riotSvc: riotSvc, sigSvc: sigSvc, metaGameSvc: metaGameSvc}
}

func (h *RiotHandler) Register(r chi.Router) {
	r.Route("/riot", func(r chi.Router) {
		r.Post("/sync/patches", h.syncPatches)
		r.Get("/patches/{version}", h.getPatchInfo)
		r.Get("/metagame/rotation/{platform}", h.analyzeChampionRotation)
		r.Get("/metagame/league/{platform}/{queue}", h.analyzeLeagueRankings)
		r.Get("/metagame/report/{platform}", h.generateMetaReport)
	})
}

type SyncPatchesResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Version string `json:"version,omitempty"`
}

func (h *RiotHandler) syncPatches(w http.ResponseWriter, r *http.Request) {
	if err := h.riotSvc.SyncPatches(r.Context()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(SyncPatchesResponse{Success: true, Message: "Patches synchronized successfully"})
}

func (h *RiotHandler) analyzeChampionRotation(w http.ResponseWriter, r *http.Request) {
	platform := chi.URLParam(r, "platform")
	if platform == "" {
		http.Error(w, "Platform parameter is required", http.StatusBadRequest)
		return
	}
	analysis, err := h.metaGameSvc.AnalyzeChampionRotation(r.Context(), platform)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error analyzing champion rotation: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"success": true, "platform": platform, "analysis": analysis})
}

func (h *RiotHandler) analyzeLeagueRankings(w http.ResponseWriter, r *http.Request) {
	platform := chi.URLParam(r, "platform")
	queue := chi.URLParam(r, "queue")
	if platform == "" || queue == "" {
		http.Error(w, "Platform and queue parameters are required", http.StatusBadRequest)
		return
	}
	analysis, err := h.metaGameSvc.AnalyzeLeagueRankings(r.Context(), platform, queue)
	if err != nil {
		http.Error(w, fmt.Errorf("Error analyzing league rankings: %w", err).Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"success": true, "platform": platform, "queue": queue, "analysis": analysis})
}

func (h *RiotHandler) generateMetaReport(w http.ResponseWriter, r *http.Request) {
	platform := chi.URLParam(r, "platform")
	if platform == "" {
		http.Error(w, "Platform parameter is required", http.StatusBadRequest)
		return
	}
	report, err := h.metaGameSvc.GenerateMetaReport(r.Context(), platform)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating meta report: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"success": true, "platform": platform, "report": report})
}

type PatchInfoResponse struct {
	Success bool            `json:"success"`
	Patch   *entities.Patch `json:"patch,omitempty"`
}

func (h *RiotHandler) getPatchInfo(w http.ResponseWriter, r *http.Request) {
	version := chi.URLParam(r, "version")
	if version == "" {
		http.Error(w, "Version parameter is required", http.StatusBadRequest)
		return
	}
	patch, err := h.riotSvc.GetPatchInfo(r.Context(), version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(PatchInfoResponse{Success: true, Patch: patch})
}
