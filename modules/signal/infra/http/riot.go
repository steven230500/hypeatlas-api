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
	riotSvc               *riot.Service
	sigSvc                in.Service
	metaGameSvc           *service.MetaGameService
	championStatsSvc      *riot.ChampionStatsService
	professionalLeagueSvc *riot.ProfessionalLeagueService
	dataDragonSvc         *riot.DataDragonService
}

func NewRiotHandler(riotSvc *riot.Service, sigSvc in.Service, metaGameSvc *service.MetaGameService) *RiotHandler {
	// Crear cliente temporal para los servicios especializados
	// En una implementación más robusta, el cliente debería ser inyectado
	tempClient := riot.NewClient("temp-key") // Esto será reemplazado por el cliente real del servicio

	return &RiotHandler{
		riotSvc:               riotSvc,
		sigSvc:                sigSvc,
		metaGameSvc:           metaGameSvc,
		championStatsSvc:      riot.NewChampionStatsService(tempClient),
		professionalLeagueSvc: riot.NewProfessionalLeagueService(tempClient),
		dataDragonSvc:         riot.NewDataDragonService(tempClient),
	}
}

func (h *RiotHandler) Register(r chi.Router) {
	r.Route("/riot", func(r chi.Router) {
		r.Post("/sync/patches", h.syncPatches)
		r.Get("/patches/{version}", h.getPatchInfo)
		r.Get("/metagame/rotation/{platform}", h.analyzeChampionRotation)
		r.Get("/metagame/league/{platform}/{queue}", h.analyzeLeagueRankings)
		r.Get("/metagame/report/{platform}", h.generateMetaReport)
		r.Get("/games", h.getGames)
		r.Get("/leagues/{platform}", h.getLeagues)
		r.Get("/regions", h.getRegions)
		r.Get("/champion-stats/{version}", h.getChampionStats)
		r.Get("/patch-changes/{fromVersion}/{toVersion}", h.getPatchChanges)
		r.Get("/pro-leagues", h.getProfessionalLeagues)
		r.Get("/pro-leagues/{league}/champions", h.getLeagueChampions)

		// Data Dragon endpoints
		r.Get("/versions", h.getGameVersions)
		r.Get("/items/{version}", h.getItems)
		r.Get("/runes/{version}", h.getRunes)
		r.Get("/summoner-spells/{version}", h.getSummonerSpells)
		r.Get("/champions/{version}/{championID}", h.getChampionDetails)
		r.Get("/patch-notes/{fromVersion}/{toVersion}", h.getPatchNotes)
	})
}

type SyncPatchesResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Version string `json:"version,omitempty"`
}

// @Summary Synchronize game patches from Riot Games
// @Description Sync latest patch information from Riot Games Data Dragon API
// @Tags riot
// @Accept json
// @Produce json
// @Success 200 {object} SyncPatchesResponse "Synchronization result"
// @Failure 500 {string} string "Internal server error"
// @Router /v1/signal/riot/sync/patches [post]
func (h *RiotHandler) syncPatches(w http.ResponseWriter, r *http.Request) {
	if err := h.riotSvc.SyncPatches(r.Context()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(SyncPatchesResponse{Success: true, Message: "Patches synchronized successfully"})
}

// @Summary Analyze weekly champion rotation
// @Description Get detailed analysis of free champion rotation including tier classification, pick rates, and strategic recommendations
// @Tags riot
// @Accept json
// @Produce json
// @Param platform path string true "Platform (e.g., na1, euw1, kr)"
// @Success 200 {object} map[string]interface{} "Analysis result with champion data and recommendations"
// @Failure 400 {string} string "Platform parameter is required"
// @Failure 500 {string} string "Internal server error"
// @Router /v1/signal/riot/metagame/rotation/{platform} [get]
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

// @Summary Analyze league rankings and statistics
// @Description Get detailed analysis of Challenger league including win rates, LP distribution, and player statistics
// @Tags riot
// @Accept json
// @Produce json
// @Param platform path string true "Platform (e.g., na1, euw1, kr)"
// @Param queue path string true "Queue type (e.g., RANKED_SOLO_5x5, RANKED_FLEX_SR)"
// @Success 200 {object} map[string]interface{} "League analysis with statistics and rankings"
// @Failure 400 {string} string "Platform and queue parameters are required"
// @Failure 500 {string} string "Internal server error"
// @Router /v1/signal/riot/metagame/league/{platform}/{queue} [get]
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

// @Summary Generate comprehensive meta-game report
// @Description Get complete meta-game analysis combining champion rotation and league statistics with insights and recommendations
// @Tags riot
// @Accept json
// @Produce json
// @Param platform path string true "Platform (e.g., na1, euw1, kr)"
// @Success 200 {object} map[string]interface{} "Complete meta-game report with analysis and insights"
// @Failure 400 {string} string "Platform parameter is required"
// @Failure 500 {string} string "Internal server error"
// @Router /v1/signal/riot/metagame/report/{platform} [get]
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

// @Summary Get detailed patch information
// @Description Retrieve detailed information about a specific game patch
// @Tags riot
// @Accept json
// @Produce json
// @Param version path string true "Patch version (e.g., 13.24.1)"
// @Success 200 {object} PatchInfoResponse "Patch information"
// @Failure 400 {string} string "Version parameter is required"
// @Failure 404 {string} string "Patch not found"
// @Router /v1/signal/riot/patches/{version} [get]
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

type GamesResponse struct {
	Success bool     `json:"success"`
	Games   []string `json:"games"`
}

// @Summary Get available Riot Games
// @Description Retrieve list of all available games from Riot Games
// @Tags riot
// @Accept json
// @Produce json
// @Success 200 {object} GamesResponse "List of available games"
// @Failure 500 {string} string "Internal server error"
// @Router /v1/signal/riot/games [get]
func (h *RiotHandler) getGames(w http.ResponseWriter, r *http.Request) {
	games, err := h.riotSvc.GetGames(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting games: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(GamesResponse{Success: true, Games: games})
}

type LeaguesResponse struct {
	Success  bool     `json:"success"`
	Platform string   `json:"platform"`
	Leagues  []string `json:"leagues"`
}

// @Summary Get available leagues for a platform
// @Description Retrieve list of all available leagues for a specific platform
// @Tags riot
// @Accept json
// @Produce json
// @Param platform path string true "Platform (e.g., na1, euw1, kr)"
// @Success 200 {object} LeaguesResponse "List of available leagues"
// @Failure 400 {string} string "Platform parameter is required"
// @Failure 500 {string} string "Internal server error"
// @Router /v1/signal/riot/leagues/{platform} [get]
func (h *RiotHandler) getLeagues(w http.ResponseWriter, r *http.Request) {
	platform := chi.URLParam(r, "platform")
	if platform == "" {
		http.Error(w, "Platform parameter is required", http.StatusBadRequest)
		return
	}

	leagues, err := h.riotSvc.GetAllLeagues(r.Context(), platform)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting leagues: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(LeaguesResponse{Success: true, Platform: platform, Leagues: leagues})
}

type RegionsResponse struct {
	Success bool     `json:"success"`
	Regions []string `json:"regions"`
}

// @Summary Get available League of Legends regions
// @Description Retrieve list of all available regions for League of Legends
// @Tags riot
// @Accept json
// @Produce json
// @Success 200 {object} RegionsResponse "List of available regions"
// @Failure 500 {string} string "Internal server error"
// @Router /v1/signal/riot/regions [get]
func (h *RiotHandler) getRegions(w http.ResponseWriter, r *http.Request) {
	regions, err := h.riotSvc.GetRegions(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting regions: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(RegionsResponse{Success: true, Regions: regions})
}

type ChampionStatsResponse struct {
	Success bool                   `json:"success"`
	Version string                 `json:"version"`
	Stats   map[string]interface{} `json:"stats"`
}

// @Summary Get champion statistics for a game version
// @Description Retrieve detailed statistics about champion usage, pick rates, and performance
// @Tags riot
// @Accept json
// @Produce json
// @Param version path string true "Game version (e.g., 13.24.1)"
// @Success 200 {object} ChampionStatsResponse "Champion statistics"
// @Failure 500 {string} string "Internal server error"
// @Router /v1/signal/riot/champion-stats/{version} [get]
func (h *RiotHandler) getChampionStats(w http.ResponseWriter, r *http.Request) {
	version := chi.URLParam(r, "version")
	if version == "" {
		http.Error(w, "Version parameter is required", http.StatusBadRequest)
		return
	}

	// Usar el servicio especializado para estadísticas de campeones
	stats, err := h.championStatsSvc.GetChampionStats(version)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting champion stats: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(ChampionStatsResponse{Success: true, Version: version, Stats: stats})
}

type PatchChangesResponse struct {
	Success     bool                   `json:"success"`
	FromVersion string                 `json:"from_version"`
	ToVersion   string                 `json:"to_version"`
	Changes     map[string]interface{} `json:"changes"`
}

// @Summary Get patch changes between versions
// @Description Compare changes between two game versions including buffs, nerfs, and new features
// @Tags riot
// @Accept json
// @Produce json
// @Param fromVersion path string true "From version (e.g., 13.23.1)"
// @Param toVersion path string true "To version (e.g., 13.24.1)"
// @Success 200 {object} PatchChangesResponse "Patch changes comparison"
// @Failure 400 {string} string "Version parameters are required"
// @Failure 500 {string} string "Internal server error"
// @Router /v1/signal/riot/patch-changes/{fromVersion}/{toVersion} [get]
func (h *RiotHandler) getPatchChanges(w http.ResponseWriter, r *http.Request) {
	fromVersion := chi.URLParam(r, "fromVersion")
	toVersion := chi.URLParam(r, "toVersion")

	if fromVersion == "" || toVersion == "" {
		http.Error(w, "Both fromVersion and toVersion parameters are required", http.StatusBadRequest)
		return
	}

	changes, err := h.riotSvc.GetPatchChanges(r.Context(), fromVersion, toVersion)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting patch changes: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(PatchChangesResponse{
		Success:     true,
		FromVersion: fromVersion,
		ToVersion:   toVersion,
		Changes:     changes,
	})
}

type ProfessionalLeaguesResponse struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data"`
}

// @Summary Get professional League of Legends leagues information
// @Description Retrieve detailed information about all major professional leagues (LEC, LCK, LPL, etc.)
// @Tags riot
// @Accept json
// @Produce json
// @Success 200 {object} ProfessionalLeaguesResponse "Professional leagues information"
// @Failure 500 {string} string "Internal server error"
// @Router /v1/signal/riot/pro-leagues [get]
func (h *RiotHandler) getProfessionalLeagues(w http.ResponseWriter, r *http.Request) {
	leagues, err := h.professionalLeagueSvc.GetProfessionalLeagues()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting professional leagues: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(ProfessionalLeaguesResponse{Success: true, Data: leagues})
}

type LeagueChampionsResponse struct {
	Success bool                   `json:"success"`
	League  string                 `json:"league"`
	Data    map[string]interface{} `json:"data"`
}

// @Summary Get champion statistics for a professional league
// @Description Retrieve detailed champion statistics including pick rates, win rates, and ban rates for a specific professional league
// @Tags riot
// @Accept json
// @Produce json
// @Param league path string true "League code (e.g., LEC, LCK, LPL, LTA)"
// @Success 200 {object} LeagueChampionsResponse "League champion statistics"
// @Failure 400 {string} string "Invalid league parameter"
// @Failure 500 {string} string "Internal server error"
// @Router /v1/signal/riot/pro-leagues/{league}/champions [get]
func (h *RiotHandler) getLeagueChampions(w http.ResponseWriter, r *http.Request) {
	league := chi.URLParam(r, "league")
	if league == "" {
		http.Error(w, "League parameter is required", http.StatusBadRequest)
		return
	}

	// Usar el servicio especializado para validar y obtener datos
	if !h.professionalLeagueSvc.ValidateLeague(league) {
		validLeagues := []string{"LEC", "LCK", "LPL", "LTA", "LCS", "VCS", "PCS"}
		http.Error(w, fmt.Sprintf("Invalid league: %s. Valid leagues: %v", league, validLeagues), http.StatusBadRequest)
		return
	}

	champions, err := h.professionalLeagueSvc.GetLeagueChampions(league)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting league champions: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(LeagueChampionsResponse{
		Success: true,
		League:  league,
		Data:    champions,
	})
}

type GameVersionsResponse struct {
	Success  bool     `json:"success"`
	Versions []string `json:"versions"`
}

// @Summary Get available game versions
// @Description Retrieve list of all available League of Legends game versions from Data Dragon
// @Tags riot
// @Accept json
// @Produce json
// @Success 200 {object} GameVersionsResponse "List of game versions"
// @Failure 500 {string} string "Internal server error"
// @Router /v1/signal/riot/versions [get]
func (h *RiotHandler) getGameVersions(w http.ResponseWriter, r *http.Request) {
	versions, err := h.dataDragonSvc.GetGameVersions(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting game versions: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(GameVersionsResponse{Success: true, Versions: versions})
}

type ItemsResponse struct {
	Success bool                   `json:"success"`
	Version string                 `json:"version"`
	Data    map[string]interface{} `json:"data"`
}

// @Summary Get items data for a specific version
// @Description Retrieve detailed information about all items for a specific game version from Data Dragon
// @Tags riot
// @Accept json
// @Produce json
// @Param version path string true "Game version (e.g., 13.24.1)"
// @Success 200 {object} ItemsResponse "Items data"
// @Failure 400 {string} string "Version parameter is required"
// @Failure 500 {string} string "Internal server error"
// @Router /v1/signal/riot/items/{version} [get]
func (h *RiotHandler) getItems(w http.ResponseWriter, r *http.Request) {
	version := chi.URLParam(r, "version")
	if version == "" {
		http.Error(w, "Version parameter is required", http.StatusBadRequest)
		return
	}

	items, err := h.dataDragonSvc.GetItems(r.Context(), version)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting items: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(ItemsResponse{Success: true, Version: version, Data: items})
}

type RunesResponse struct {
	Success bool                   `json:"success"`
	Version string                 `json:"version"`
	Data    map[string]interface{} `json:"data"`
}

// @Summary Get runes data for a specific version
// @Description Retrieve detailed information about all runes for a specific game version from Data Dragon
// @Tags riot
// @Accept json
// @Produce json
// @Param version path string true "Game version (e.g., 13.24.1)"
// @Success 200 {object} RunesResponse "Runes data"
// @Failure 400 {string} string "Version parameter is required"
// @Failure 500 {string} string "Internal server error"
// @Router /v1/signal/riot/runes/{version} [get]
func (h *RiotHandler) getRunes(w http.ResponseWriter, r *http.Request) {
	version := chi.URLParam(r, "version")
	if version == "" {
		http.Error(w, "Version parameter is required", http.StatusBadRequest)
		return
	}

	runes, err := h.dataDragonSvc.GetRunes(r.Context(), version)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting runes: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(RunesResponse{Success: true, Version: version, Data: runes})
}

type SummonerSpellsResponse struct {
	Success bool                   `json:"success"`
	Version string                 `json:"version"`
	Data    map[string]interface{} `json:"data"`
}

// @Summary Get summoner spells data for a specific version
// @Description Retrieve detailed information about all summoner spells for a specific game version from Data Dragon
// @Tags riot
// @Accept json
// @Produce json
// @Param version path string true "Game version (e.g., 13.24.1)"
// @Success 200 {object} SummonerSpellsResponse "Summoner spells data"
// @Failure 400 {string} string "Version parameter is required"
// @Failure 500 {string} string "Internal server error"
// @Router /v1/signal/riot/summoner-spells/{version} [get]
func (h *RiotHandler) getSummonerSpells(w http.ResponseWriter, r *http.Request) {
	version := chi.URLParam(r, "version")
	if version == "" {
		http.Error(w, "Version parameter is required", http.StatusBadRequest)
		return
	}

	spells, err := h.dataDragonSvc.GetSummonerSpells(r.Context(), version)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting summoner spells: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(SummonerSpellsResponse{Success: true, Version: version, Data: spells})
}

type ChampionDetailsResponse struct {
	Success    bool                   `json:"success"`
	Version    string                 `json:"version"`
	ChampionID string                 `json:"champion_id"`
	Data       map[string]interface{} `json:"data"`
}

// @Summary Get detailed champion information
// @Description Retrieve complete information about a specific champion for a game version from Data Dragon
// @Tags riot
// @Accept json
// @Produce json
// @Param version path string true "Game version (e.g., 13.24.1)"
// @Param championID path string true "Champion ID (e.g., Ahri, Jinx)"
// @Success 200 {object} ChampionDetailsResponse "Champion details"
// @Failure 400 {string} string "Version and championID parameters are required"
// @Failure 500 {string} string "Internal server error"
// @Router /v1/signal/riot/champions/{version}/{championID} [get]
func (h *RiotHandler) getChampionDetails(w http.ResponseWriter, r *http.Request) {
	version := chi.URLParam(r, "version")
	championID := chi.URLParam(r, "championID")

	if version == "" || championID == "" {
		http.Error(w, "Version and championID parameters are required", http.StatusBadRequest)
		return
	}

	details, err := h.dataDragonSvc.GetChampionDetails(r.Context(), version, championID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting champion details: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(ChampionDetailsResponse{
		Success:    true,
		Version:    version,
		ChampionID: championID,
		Data:       details,
	})
}

type PatchNotesResponse struct {
	Success     bool                   `json:"success"`
	FromVersion string                 `json:"from_version"`
	ToVersion   string                 `json:"to_version"`
	Data        map[string]interface{} `json:"data"`
}

// @Summary Get patch notes and changes between versions
// @Description Compare changes between two game versions including new/removed champions and modifications
// @Tags riot
// @Accept json
// @Produce json
// @Param fromVersion path string true "From version (e.g., 13.23.1)"
// @Param toVersion path string true "To version (e.g., 13.24.1)"
// @Success 200 {object} PatchNotesResponse "Patch notes comparison"
// @Failure 400 {string} string "Version parameters are required"
// @Failure 500 {string} string "Internal server error"
// @Router /v1/signal/riot/patch-notes/{fromVersion}/{toVersion} [get]
func (h *RiotHandler) getPatchNotes(w http.ResponseWriter, r *http.Request) {
	fromVersion := chi.URLParam(r, "fromVersion")
	toVersion := chi.URLParam(r, "toVersion")

	if fromVersion == "" || toVersion == "" {
		http.Error(w, "Both fromVersion and toVersion parameters are required", http.StatusBadRequest)
		return
	}

	notes, err := h.dataDragonSvc.GetPatchNotes(r.Context(), fromVersion, toVersion)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting patch notes: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(PatchNotesResponse{
		Success:     true,
		FromVersion: fromVersion,
		ToVersion:   toVersion,
		Data:        notes,
	})
}
