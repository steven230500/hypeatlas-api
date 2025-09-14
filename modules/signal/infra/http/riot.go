package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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
	imageSvc              *riot.ImageService
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
		imageSvc:              riot.NewImageService(tempClient),
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

		// Image endpoints
		r.Get("/images/champions/{version}/{championID}", h.getChampionImages)
		r.Get("/images/champions/{version}/{championID}/{skinNum}", h.getChampionSkinImages)
		r.Get("/images/items/{version}/{itemID}", h.getItemImage)
		r.Get("/images/spells/{version}/{spellName}", h.getSpellImage)
		r.Get("/images/runes/{runeIcon}", h.getRuneImage)
		r.Get("/images/profile-icons/{version}/{iconID}", h.getProfileIconImage)
		r.Get("/images/maps/{version}/{mapID}", h.getMapImage)
		r.Get("/images/abilities/{version}/{abilityName}", h.getAbilityImage)
		r.Get("/images/passives/{version}/{passiveFile}", h.getPassiveImage)
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

// Image Handlers

type ChampionImagesResponse struct {
	Success  bool                   `json:"success"`
	Version  string                 `json:"version"`
	Champion string                 `json:"champion"`
	Data     map[string]interface{} `json:"data"`
}

// @Summary Get champion image URLs
// @Description Retrieve all available image URLs for a specific champion
// @Tags riot
// @Accept json
// @Produce json
// @Param version path string true "Game version (e.g., 13.24.1)"
// @Param championID path string true "Champion ID (e.g., Ahri, Jinx)"
// @Success 200 {object} ChampionImagesResponse "Champion image URLs"
// @Failure 400 {string} string "Version and championID parameters are required"
// @Failure 500 {string} string "Internal server error"
// @Router /v1/signal/riot/images/champions/{version}/{championID} [get]
func (h *RiotHandler) getChampionImages(w http.ResponseWriter, r *http.Request) {
	version := chi.URLParam(r, "version")
	championID := chi.URLParam(r, "championID")

	if version == "" || championID == "" {
		http.Error(w, "Version and championID parameters are required", http.StatusBadRequest)
		return
	}

	images, err := h.imageSvc.GetChampionImageURLs(r.Context(), version, championID, 0)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting champion images: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(ChampionImagesResponse{
		Success:  true,
		Version:  version,
		Champion: championID,
		Data:     images,
	})
}

type ChampionSkinImagesResponse struct {
	Success  bool                   `json:"success"`
	Version  string                 `json:"version"`
	Champion string                 `json:"champion"`
	SkinNum  int                    `json:"skin_num"`
	Data     map[string]interface{} `json:"data"`
}

// @Summary Get champion skin image URLs
// @Description Retrieve all available image URLs for a specific champion skin
// @Tags riot
// @Accept json
// @Produce json
// @Param version path string true "Game version (e.g., 13.24.1)"
// @Param championID path string true "Champion ID (e.g., Ahri, Jinx)"
// @Param skinNum path int true "Skin number (0 for default skin)"
// @Success 200 {object} ChampionSkinImagesResponse "Champion skin image URLs"
// @Failure 400 {string} string "Parameters are required"
// @Failure 500 {string} string "Internal server error"
// @Router /v1/signal/riot/images/champions/{version}/{championID}/{skinNum} [get]
func (h *RiotHandler) getChampionSkinImages(w http.ResponseWriter, r *http.Request) {
	version := chi.URLParam(r, "version")
	championID := chi.URLParam(r, "championID")
	skinNumStr := chi.URLParam(r, "skinNum")

	if version == "" || championID == "" || skinNumStr == "" {
		http.Error(w, "Version, championID, and skinNum parameters are required", http.StatusBadRequest)
		return
	}

	skinNum, err := strconv.Atoi(skinNumStr)
	if err != nil {
		http.Error(w, "Invalid skin number", http.StatusBadRequest)
		return
	}

	images, err := h.imageSvc.GetChampionImageURLs(r.Context(), version, championID, skinNum)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting champion skin images: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(ChampionSkinImagesResponse{
		Success:  true,
		Version:  version,
		Champion: championID,
		SkinNum:  skinNum,
		Data:     images,
	})
}

type ItemImageResponse struct {
	Success bool                   `json:"success"`
	Version string                 `json:"version"`
	ItemID  string                 `json:"item_id"`
	Data    map[string]interface{} `json:"data"`
}

// @Summary Get item image URL
// @Description Retrieve the image URL for a specific item
// @Tags riot
// @Accept json
// @Produce json
// @Param version path string true "Game version (e.g., 13.24.1)"
// @Param itemID path string true "Item ID (e.g., 1001, 3153)"
// @Success 200 {object} ItemImageResponse "Item image URL"
// @Failure 400 {string} string "Parameters are required"
// @Failure 500 {string} string "Internal server error"
// @Router /v1/signal/riot/images/items/{version}/{itemID} [get]
func (h *RiotHandler) getItemImage(w http.ResponseWriter, r *http.Request) {
	version := chi.URLParam(r, "version")
	itemID := chi.URLParam(r, "itemID")

	if version == "" || itemID == "" {
		http.Error(w, "Version and itemID parameters are required", http.StatusBadRequest)
		return
	}

	image, err := h.imageSvc.GetItemImageURL(r.Context(), version, itemID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting item image: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(ItemImageResponse{
		Success: true,
		Version: version,
		ItemID:  itemID,
		Data:    image,
	})
}

type SpellImageResponse struct {
	Success   bool                   `json:"success"`
	Version   string                 `json:"version"`
	SpellName string                 `json:"spell_name"`
	Data      map[string]interface{} `json:"data"`
}

// @Summary Get summoner spell image URL
// @Description Retrieve the image URL for a specific summoner spell
// @Tags riot
// @Accept json
// @Produce json
// @Param version path string true "Game version (e.g., 13.24.1)"
// @Param spellName path string true "Spell name (e.g., SummonerFlash, SummonerHeal)"
// @Success 200 {object} SpellImageResponse "Spell image URL"
// @Failure 400 {string} string "Parameters are required"
// @Failure 500 {string} string "Internal server error"
// @Router /v1/signal/riot/images/spells/{version}/{spellName} [get]
func (h *RiotHandler) getSpellImage(w http.ResponseWriter, r *http.Request) {
	version := chi.URLParam(r, "version")
	spellName := chi.URLParam(r, "spellName")

	if version == "" || spellName == "" {
		http.Error(w, "Version and spellName parameters are required", http.StatusBadRequest)
		return
	}

	image, err := h.imageSvc.GetSpellImageURL(r.Context(), version, spellName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting spell image: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(SpellImageResponse{
		Success:   true,
		Version:   version,
		SpellName: spellName,
		Data:      image,
	})
}

type RuneImageResponse struct {
	Success  bool                   `json:"success"`
	RuneIcon string                 `json:"rune_icon"`
	Data     map[string]interface{} `json:"data"`
}

// @Summary Get rune image URL
// @Description Retrieve the image URL for a specific rune
// @Tags riot
// @Accept json
// @Produce json
// @Param runeIcon path string true "Rune icon path (e.g., perk-images/Styles/Domination/Electrocute/Electrocute.png)"
// @Success 200 {object} RuneImageResponse "Rune image URL"
// @Failure 400 {string} string "Rune icon parameter is required"
// @Failure 500 {string} string "Internal server error"
// @Router /v1/signal/riot/images/runes/{runeIcon} [get]
func (h *RiotHandler) getRuneImage(w http.ResponseWriter, r *http.Request) {
	runeIcon := chi.URLParam(r, "runeIcon")

	if runeIcon == "" {
		http.Error(w, "Rune icon parameter is required", http.StatusBadRequest)
		return
	}

	image, err := h.imageSvc.GetRuneImageURL(r.Context(), runeIcon)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting rune image: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(RuneImageResponse{
		Success:  true,
		RuneIcon: runeIcon,
		Data:     image,
	})
}

type ProfileIconImageResponse struct {
	Success bool                   `json:"success"`
	Version string                 `json:"version"`
	IconID  int                    `json:"icon_id"`
	Data    map[string]interface{} `json:"data"`
}

// @Summary Get profile icon image URL
// @Description Retrieve the image URL for a specific profile icon
// @Tags riot
// @Accept json
// @Produce json
// @Param version path string true "Game version (e.g., 13.24.1)"
// @Param iconID path int true "Profile icon ID"
// @Success 200 {object} ProfileIconImageResponse "Profile icon image URL"
// @Failure 400 {string} string "Parameters are required"
// @Failure 500 {string} string "Internal server error"
// @Router /v1/signal/riot/images/profile-icons/{version}/{iconID} [get]
func (h *RiotHandler) getProfileIconImage(w http.ResponseWriter, r *http.Request) {
	version := chi.URLParam(r, "version")
	iconIDStr := chi.URLParam(r, "iconID")

	if version == "" || iconIDStr == "" {
		http.Error(w, "Version and iconID parameters are required", http.StatusBadRequest)
		return
	}

	iconID, err := strconv.Atoi(iconIDStr)
	if err != nil {
		http.Error(w, "Invalid icon ID", http.StatusBadRequest)
		return
	}

	image, err := h.imageSvc.GetProfileIconImageURL(r.Context(), version, iconID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting profile icon image: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(ProfileIconImageResponse{
		Success: true,
		Version: version,
		IconID:  iconID,
		Data:    image,
	})
}

type MapImageResponse struct {
	Success bool                   `json:"success"`
	Version string                 `json:"version"`
	MapID   int                    `json:"map_id"`
	Data    map[string]interface{} `json:"data"`
}

// @Summary Get map image URL
// @Description Retrieve the image URL for a specific map
// @Tags riot
// @Accept json
// @Produce json
// @Param version path string true "Game version (e.g., 13.24.1)"
// @Param mapID path int true "Map ID (e.g., 11 for Summoner's Rift)"
// @Success 200 {object} MapImageResponse "Map image URL"
// @Failure 400 {string} string "Parameters are required"
// @Failure 500 {string} string "Internal server error"
// @Router /v1/signal/riot/images/maps/{version}/{mapID} [get]
func (h *RiotHandler) getMapImage(w http.ResponseWriter, r *http.Request) {
	version := chi.URLParam(r, "version")
	mapIDStr := chi.URLParam(r, "mapID")

	if version == "" || mapIDStr == "" {
		http.Error(w, "Version and mapID parameters are required", http.StatusBadRequest)
		return
	}

	mapID, err := strconv.Atoi(mapIDStr)
	if err != nil {
		http.Error(w, "Invalid map ID", http.StatusBadRequest)
		return
	}

	image, err := h.imageSvc.GetMapImageURL(r.Context(), version, mapID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting map image: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(MapImageResponse{
		Success: true,
		Version: version,
		MapID:   mapID,
		Data:    image,
	})
}

type AbilityImageResponse struct {
	Success     bool                   `json:"success"`
	Version     string                 `json:"version"`
	AbilityName string                 `json:"ability_name"`
	Data        map[string]interface{} `json:"data"`
}

// @Summary Get champion ability image URL
// @Description Retrieve the image URL for a specific champion ability
// @Tags riot
// @Accept json
// @Produce json
// @Param version path string true "Game version (e.g., 13.24.1)"
// @Param abilityName path string true "Ability name (e.g., AhriQ, AhriW)"
// @Success 200 {object} AbilityImageResponse "Ability image URL"
// @Failure 400 {string} string "Parameters are required"
// @Failure 500 {string} string "Internal server error"
// @Router /v1/signal/riot/images/abilities/{version}/{abilityName} [get]
func (h *RiotHandler) getAbilityImage(w http.ResponseWriter, r *http.Request) {
	version := chi.URLParam(r, "version")
	abilityName := chi.URLParam(r, "abilityName")

	if version == "" || abilityName == "" {
		http.Error(w, "Version and abilityName parameters are required", http.StatusBadRequest)
		return
	}

	image, err := h.imageSvc.GetAbilityImageURL(r.Context(), version, abilityName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting ability image: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(AbilityImageResponse{
		Success:     true,
		Version:     version,
		AbilityName: abilityName,
		Data:        image,
	})
}

type PassiveImageResponse struct {
	Success     bool                   `json:"success"`
	Version     string                 `json:"version"`
	PassiveFile string                 `json:"passive_file"`
	Data        map[string]interface{} `json:"data"`
}

// @Summary Get champion passive image URL
// @Description Retrieve the image URL for a specific champion passive
// @Tags riot
// @Accept json
// @Produce json
// @Param version path string true "Game version (e.g., 13.24.1)"
// @Param passiveFile path string true "Passive file name"
// @Success 200 {object} PassiveImageResponse "Passive image URL"
// @Failure 400 {string} string "Parameters are required"
// @Failure 500 {string} string "Internal server error"
// @Router /v1/signal/riot/images/passives/{version}/{passiveFile} [get]
func (h *RiotHandler) getPassiveImage(w http.ResponseWriter, r *http.Request) {
	version := chi.URLParam(r, "version")
	passiveFile := chi.URLParam(r, "passiveFile")

	if version == "" || passiveFile == "" {
		http.Error(w, "Version and passiveFile parameters are required", http.StatusBadRequest)
		return
	}

	image, err := h.imageSvc.GetPassiveImageURL(r.Context(), version, passiveFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting passive image: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(PassiveImageResponse{
		Success:     true,
		Version:     version,
		PassiveFile: passiveFile,
		Data:        image,
	})
}
