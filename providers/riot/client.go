package riot

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// RateLimiter maneja los límites de rate de Riot Games
type RateLimiter struct {
	mu           sync.Mutex
	requests     []time.Time
	maxPerSecond int
	maxPerMinute int
}

// NewRateLimiter crea un nuevo rate limiter
func NewRateLimiter(maxPerSecond, maxPerMinute int) *RateLimiter {
	return &RateLimiter{
		requests:     make([]time.Time, 0),
		maxPerSecond: maxPerSecond,
		maxPerMinute: maxPerMinute,
	}
}

// Wait bloquea hasta que sea seguro hacer una petición
func (rl *RateLimiter) Wait() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// Limpiar requests antiguos (más de 1 minuto)
	cutoff := now.Add(-time.Minute)
	validRequests := make([]time.Time, 0)
	for _, req := range rl.requests {
		if req.After(cutoff) {
			validRequests = append(validRequests, req)
		}
	}
	rl.requests = validRequests

	// Verificar límites
	if len(rl.requests) >= rl.maxPerMinute {
		// Esperar hasta que pase 1 minuto desde la primera request
		waitTime := time.Minute - now.Sub(rl.requests[0])
		if waitTime > 0 {
			time.Sleep(waitTime)
			now = time.Now()
		}
	}

	// Verificar límite por segundo
	if len(rl.requests) > 0 {
		lastRequest := rl.requests[len(rl.requests)-1]
		if now.Sub(lastRequest) < time.Second {
			time.Sleep(time.Second - now.Sub(lastRequest))
		}
	}

	rl.requests = append(rl.requests, now)
}

// Client para Riot Games API
type Client struct {
	apiKey      string
	client      *http.Client
	rateLimiter *RateLimiter
}

// NewClient crea un nuevo cliente de Riot Games API
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		rateLimiter: NewRateLimiter(18, 95), // 18 req/s, 95 req/min (margen de seguridad)
	}
}

// makeRequest hace una petición HTTP con el token de Riot y rate limiting
func (c *Client) makeRequest(method, url string, body io.Reader) (*http.Response, error) {
	return c.makeRequestWithAuth(method, url, body, true)
}

// makeRequestWithoutAuth hace una petición HTTP sin autenticación (para Data Dragon)
func (c *Client) makeRequestWithoutAuth(method, url string, body io.Reader) (*http.Response, error) {
	return c.makeRequestWithAuth(method, url, body, false)
}

// makeRequestWithAuth hace una petición HTTP con control de autenticación
func (c *Client) makeRequestWithAuth(method, url string, body io.Reader, withAuth bool) (*http.Response, error) {
	// Aplicar rate limiting solo para requests autenticados
	if withAuth {
		c.rateLimiter.Wait()
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Agregar header de autenticación solo si se requiere
	if withAuth {
		req.Header.Set("X-Riot-Token", c.apiKey)
		fmt.Printf("[RIOT-DEBUG] Making authenticated request to: %s\n", url)
		fmt.Printf("[RIOT-DEBUG] X-Riot-Token: %s\n", c.apiKey)
	} else {
		fmt.Printf("[RIOT-DEBUG] Making unauthenticated request to: %s\n", url)
	}
	req.Header.Set("User-Agent", "HypeAtlas-API/1.0")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	fmt.Printf("[RIOT-DEBUG] Response status: %d for URL: %s\n", resp.StatusCode, url)

	// Manejar rate limiting (429) solo para requests autenticados
	if withAuth && resp.StatusCode == 429 {
		resp.Body.Close()

		// Obtener tiempo de espera del header Retry-After
		retryAfter := resp.Header.Get("Retry-After")
		if retryAfter != "" {
			if seconds, err := strconv.Atoi(retryAfter); err == nil {
				time.Sleep(time.Duration(seconds) * time.Second)
				// Reintentar la petición
				return c.makeRequestWithAuth(method, url, body, withAuth)
			}
		}

		// Si no hay Retry-After, esperar 1 segundo por defecto
		time.Sleep(time.Second)
		return c.makeRequestWithAuth(method, url, body, withAuth)
	}

	return resp, nil
}

// VersionResponse respuesta de la API de versiones
type VersionResponse []string

// ChampionData datos de un campeón
type ChampionData struct {
	ID    string `json:"id"`
	Key   string `json:"key"`
	Name  string `json:"name"`
	Title string `json:"title"`
}

// ChampionsResponse respuesta de la API de campeones
type ChampionsResponse struct {
	Data map[string]ChampionData `json:"data"`
}

// GetLatestVersion obtiene la versión más reciente del juego desde Data Dragon
func (c *Client) GetLatestVersion() (string, error) {
	url := "https://ddragon.leagueoflegends.com/api/versions.json"

	resp, err := c.makeRequestWithoutAuth("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	var versions VersionResponse
	if err := json.Unmarshal(body, &versions); err != nil {
		return "", fmt.Errorf("error parsing response: %w", err)
	}

	if len(versions) == 0 {
		return "", fmt.Errorf("no versions found")
	}

	return versions[0], nil
}

// GetChampions obtiene la lista de campeones para una versión específica desde Data Dragon
func (c *Client) GetChampions(version string) (*ChampionsResponse, error) {
	url := fmt.Sprintf("https://ddragon.leagueoflegends.com/cdn/%s/data/en_US/champion.json", version)

	resp, err := c.makeRequestWithoutAuth("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	var champions ChampionsResponse
	if err := json.Unmarshal(body, &champions); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return &champions, nil
}

// GetChampion obtiene datos específicos de un campeón desde Data Dragon
func (c *Client) GetChampion(version, championID string) (*ChampionData, error) {
	url := fmt.Sprintf("https://ddragon.leagueoflegends.com/cdn/%s/data/en_US/champion/%s.json", version, championID)

	resp, err := c.makeRequestWithoutAuth("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	var result struct {
		Data map[string]ChampionData `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	if champion, exists := result.Data[championID]; exists {
		return &champion, nil
	}

	return nil, fmt.Errorf("champion %s not found", championID)
}

// SummonerData respuesta de la API de summoner
type SummonerData struct {
	ID            string `json:"id"`
	AccountID     string `json:"accountId"`
	PUUID         string `json:"puuid"`
	Name          string `json:"name"`
	ProfileIconID int    `json:"profileIconId"`
	RevisionDate  int64  `json:"revisionDate"`
	SummonerLevel int    `json:"summonerLevel"`
}

// GetSummonerByPUUID obtiene información de un summoner por PUUID
func (c *Client) GetSummonerByPUUID(platform, puuid string) (*SummonerData, error) {
	url := fmt.Sprintf("https://%s.api.riotgames.com/lol/summoner/v4/summoners/by-puuid/%s", platform, puuid)

	resp, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	var summoner SummonerData
	if err := json.Unmarshal(body, &summoner); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return &summoner, nil
}

// MatchIDsResponse lista de IDs de partidas
type MatchIDsResponse []string

// GetMatchIDsByPUUID obtiene lista de IDs de partidas recientes de un jugador
func (c *Client) GetMatchIDsByPUUID(region, puuid string, count int) (MatchIDsResponse, error) {
	if count <= 0 || count > 100 {
		count = 20
	}

	url := fmt.Sprintf("https://%s.api.riotgames.com/lol/match/v5/matches/by-puuid/%s/ids?count=%d", region, puuid, count)

	resp, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	var matchIDs MatchIDsResponse
	if err := json.Unmarshal(body, &matchIDs); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return matchIDs, nil
}

// VALContentResponse respuesta de la API de contenido de Valorant
type VALContentResponse struct {
	Version string `json:"version"`
	// Aquí irían más campos según la documentación
}

// GetVALContent obtiene el contenido de Valorant (acts, seasons, etc.)
func (c *Client) GetVALContent() (*VALContentResponse, error) {
	url := "https://na.api.riotgames.com/val/content/v1/contents"

	resp, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	var content VALContentResponse
	if err := json.Unmarshal(body, &content); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return &content, nil
}

// ChampionRotationResponse respuesta de la API de rotación de campeones
type ChampionRotationResponse struct {
	FreeChampionIDs              []int `json:"freeChampionIds"`
	FreeChampionIDsForNewPlayers []int `json:"freeChampionIdsForNewPlayers"`
	MaxNewPlayerLevel            int   `json:"maxNewPlayerLevel"`
}

// GetChampionRotation obtiene la rotación semanal de campeones gratuitos
func (c *Client) GetChampionRotation(platform string) (*ChampionRotationResponse, error) {
	url := fmt.Sprintf("https://%s.api.riotgames.com/lol/platform/v3/champion-rotations", platform)

	resp, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	var rotation ChampionRotationResponse
	if err := json.Unmarshal(body, &rotation); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return &rotation, nil
}

// LeagueEntry respuesta para entradas de liga
type LeagueEntry struct {
	SummonerID   string `json:"summonerId"`
	SummonerName string `json:"summonerName"`
	LeaguePoints int    `json:"leaguePoints"`
	Rank         string `json:"rank"`
	Wins         int    `json:"wins"`
	Losses       int    `json:"losses"`
	Veteran      bool   `json:"veteran"`
	Inactive     bool   `json:"inactive"`
	FreshBlood   bool   `json:"freshBlood"`
	HotStreak    bool   `json:"hotStreak"`
}

// LeagueEntriesResponse lista de entradas de liga
type LeagueEntriesResponse []LeagueEntry

// GetLeagueEntries obtiene entradas de liga por división y tier
func (c *Client) GetLeagueEntries(platform, queue, tier, division string, page int) (LeagueEntriesResponse, error) {
	url := fmt.Sprintf("https://%s.api.riotgames.com/lol/league/v4/entries/%s/%s/%s?page=%d",
		platform, queue, tier, division, page)

	resp, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	var entries LeagueEntriesResponse
	if err := json.Unmarshal(body, &entries); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return entries, nil
}

// GetChallengerLeague obtiene la liga Challenger para una cola específica
func (c *Client) GetChallengerLeague(platform, queue string) (*LeagueData, error) {
	url := fmt.Sprintf("https://%s.api.riotgames.com/lol/league/v4/challengerleagues/by-queue/%s", platform, queue)

	resp, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	var league LeagueData
	if err := json.Unmarshal(body, &league); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return &league, nil
}

// LeagueData respuesta para datos de liga
type LeagueData struct {
	Tier     string        `json:"tier"`
	LeagueID string        `json:"leagueId"`
	Queue    string        `json:"queue"`
	Name     string        `json:"name"`
	Entries  []LeagueEntry `json:"entries"`
}

// ChampionMastery respuesta para maestría de campeón
type ChampionMastery struct {
	ChampionID                   int    `json:"championId"`
	ChampionLevel                int    `json:"championLevel"`
	ChampionPoints               int    `json:"championPoints"`
	LastPlayTime                 int64  `json:"lastPlayTime"`
	ChampionPointsSinceLastLevel int    `json:"championPointsSinceLastLevel"`
	ChampionPointsUntilNextLevel int    `json:"championPointsUntilNextLevel"`
	ChestGranted                 bool   `json:"chestGranted"`
	TokensEarned                 int    `json:"tokensEarned"`
	SummonerID                   string `json:"summonerId"`
}

// ChampionMasteriesResponse lista de maestrías de campeón
type ChampionMasteriesResponse []ChampionMastery

// GetChampionMasteries obtiene todas las maestrías de campeón de un summoner
func (c *Client) GetChampionMasteries(platform, puuid string) (ChampionMasteriesResponse, error) {
	url := fmt.Sprintf("https://%s.api.riotgames.com/lol/champion-mastery/v4/champion-masteries/by-puuid/%s", platform, puuid)

	resp, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	var masteries ChampionMasteriesResponse
	if err := json.Unmarshal(body, &masteries); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return masteries, nil
}

// GetAllLeagues obtiene todas las ligas disponibles para una plataforma
func (c *Client) GetAllLeagues(platform string) ([]string, error) {
	// Riot Games no tiene un endpoint directo para obtener todas las ligas
	// Pero podemos obtener las ligas Challenger para las colas principales
	queues := []string{"RANKED_SOLO_5x5", "RANKED_FLEX_SR", "RANKED_FLEX_TT"}

	var leagues []string
	for _, queue := range queues {
		league, err := c.GetChallengerLeague(platform, queue)
		if err != nil {
			// Log error but continue with other queues
			fmt.Printf("[RIOT-DEBUG] Error getting league for queue %s: %v\n", queue, err)
			continue
		}

		// Add league name if not already in the list
		found := false
		for _, existing := range leagues {
			if existing == league.Name {
				found = true
				break
			}
		}
		if !found {
			leagues = append(leagues, league.Name)
		}
	}

	return leagues, nil
}

// GetGames obtiene la lista de juegos disponibles de Riot Games
func (c *Client) GetGames() ([]string, error) {
	// Lista de juegos principales de Riot Games
	// Esta información podría venir de una API o ser hardcodeada
	games := []string{
		"League of Legends",
		"Valorant",
		"Teamfight Tactics",
		"Legends of Runeterra",
		"Wild Rift",
	}

	return games, nil
}

// GetRegions obtiene la lista de regiones disponibles para League of Legends
func (c *Client) GetRegions() ([]string, error) {
	// Lista de regiones oficiales de League of Legends
	regions := []string{
		"BR1", "EUN1", "EUW1", "JP1", "KR", "LA1", "LA2",
		"NA1", "OC1", "PH2", "RU", "SG2", "TH2", "TR1", "TW2", "VN2",
	}

	return regions, nil
}

// GetChampionStats obtiene estadísticas de uso de campeones desde Data Dragon
func (c *Client) GetChampionStats(version string) (map[string]interface{}, error) {
	// Esta sería una implementación más avanzada que requeriría
	// análisis de datos de partidos o estadísticas agregadas
	// Por ahora retornamos una estructura básica
	stats := map[string]interface{}{
		"version": version,
		"note":    "Advanced champion statistics require match data analysis",
		"available_endpoints": []string{
			"champion-rotations",
			"champion-mastery",
			"league-entries",
		},
	}

	return stats, nil
}

// GetPatchChanges obtiene cambios de campeones entre parches desde Data Dragon
func (c *Client) GetPatchChanges(fromVersion, toVersion string) (map[string]interface{}, error) {
	// Comparación de cambios entre versiones
	// Esto requeriría análisis de datos de Data Dragon
	changes := map[string]interface{}{
		"from_version": fromVersion,
		"to_version":   toVersion,
		"note":         "Patch change analysis requires Data Dragon comparison",
		"buffs":        []string{},
		"nerfs":        []string{},
		"new_features": []string{},
	}

	return changes, nil
}
