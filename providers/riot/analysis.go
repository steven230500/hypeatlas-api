package riot

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
)

// RegionMetaReport contiene el reporte del análisis del meta regional
type RegionMetaReport struct {
	Region             string             `json:"region"`
	AnalyzedMatches    int                `json:"analyzed_matches"`
	AnalyzedPlayers    int                `json:"analyzed_players"`
	TopChampions       []ChampionMetaStat `json:"top_champions"`
	TopBans            []ChampionMetaStat `json:"top_bans"`
	CommonCompositions []CompositionStat  `json:"common_compositions"`
}

// ChampionMetaStat estadísticas de un campeón en el meta
type ChampionMetaStat struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	PickRate float64 `json:"pick_rate"`
	WinRate  float64 `json:"win_rate"`
	BanRate  float64 `json:"ban_rate"`
	Games    int     `json:"games"`
	Wins     int     `json:"wins"`
	Bans     int     `json:"bans"`
}

// CompositionStat estadísticas de una composición
type CompositionStat struct {
	Champions []string `json:"champions"`
	Wins      int      `json:"wins"`
	Games     int      `json:"games"`
	WinRate   float64  `json:"win_rate"`
}

// AnalyzeRegionMeta analiza el meta de una región basado en partidas de Challenger
func (c *Client) AnalyzeRegionMeta(platform string, playerSampleSize, matchesPerPlayer int) (*RegionMetaReport, error) {
	// 0. Obtener datos estáticos de campeones (Data Dragon) para resolver nombres
	latestVersion, err := c.GetLatestVersion()
	if err != nil {
		return nil, fmt.Errorf("error getting latest version: %w", err)
	}
	championsData, err := c.GetChampions(latestVersion)
	if err != nil {
		return nil, fmt.Errorf("error getting champions data: %w", err)
	}

	// Mapa ID -> Nombre
	champIDToName := make(map[int]string)
	for _, champ := range championsData.Data {
		id, _ := strconv.Atoi(champ.Key)
		champIDToName[id] = champ.Name
	}

	// 1. Obtener Challenger League
	league, err := c.GetChallengerLeague(platform, "RANKED_SOLO_5x5")
	if err != nil {
		return nil, fmt.Errorf("error getting challenger league: %w", err)
	}

	// 2. Seleccionar top jugadores
	entries := league.Entries
	// Ordenar por LP descendente
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].LeaguePoints > entries[j].LeaguePoints
	})

	if len(entries) > playerSampleSize {
		entries = entries[:playerSampleSize]
	}

	// 3. Recolectar Match IDs únicos
	uniqueMatchIds := make(map[string]bool)
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Canal para limitar concurrencia de jugadores
	semPlayers := make(chan struct{}, 10)

	fmt.Printf("Analyzing %d players from %s...\n", len(entries), platform)

	for _, entry := range entries {
		wg.Add(1)
		go func(entry LeagueEntry) {
			defer wg.Done()
			semPlayers <- struct{}{} // Acquire token
			defer func() { <-semPlayers }() // Release token

			// Obtener Match IDs
			puuid := entry.Puuid
			if puuid == "" {
				summoner, err := c.GetSummonerBySummonerId(platform, entry.SummonerID)
				if err != nil {
					fmt.Printf("Error getting summoner %s: %v\n", entry.SummonerName, err)
					return
				}
				puuid = summoner.Puuid
			}

			ids, err := c.GetMatchIds(platform, puuid, matchesPerPlayer)
			if err != nil {
				fmt.Printf("Error getting matches for %s: %v\n", entry.SummonerName, err)
				return
			}

			mu.Lock()
			for _, id := range ids {
				uniqueMatchIds[id] = true
			}
			mu.Unlock()
		}(entry)
	}

	wg.Wait()

	// 4. Analizar partidas
	matchIds := make([]string, 0, len(uniqueMatchIds))
	for id := range uniqueMatchIds {
		matchIds = append(matchIds, id)
	}

	fmt.Printf("Found %d unique matches. Analyzing details...\n", len(matchIds))

	champStats := make(map[int]*ChampionMetaStat)
	totalGames := 0

	// Worker pool para procesar partidas concurrentemente
	// Aumentamos concurrencia ya que el RateLimiter interno gestionará los límites
	concurrency := 20
	semMatches := make(chan struct{}, concurrency)
	var statsMu sync.Mutex

	for _, matchId := range matchIds {
		wg.Add(1)
		go func(mId string) {
			defer wg.Done()
			semMatches <- struct{}{} // Acquire token
			defer func() { <-semMatches }() // Release token

			match, err := c.GetMatch(platform, mId)
			if err != nil {
				fmt.Printf("Error getting match %s: %v\n", mId, err)
				return
			}

			statsMu.Lock()
			defer statsMu.Unlock()

			totalGames++

			// Procesar Bans
			for _, team := range match.Info.Teams {
				for _, ban := range team.Bans {
					if _, exists := champStats[ban.ChampionID]; !exists {
						name := champIDToName[ban.ChampionID]
						if name == "" {
							name = "Unknown"
						}
						champStats[ban.ChampionID] = &ChampionMetaStat{ID: ban.ChampionID, Name: name}
					}
					champStats[ban.ChampionID].Bans++
				}
			}

			// Procesar Picks y Wins
			for _, p := range match.Info.Participants {
				if _, exists := champStats[p.ChampionID]; !exists {
					name := champIDToName[p.ChampionID]
					if name == "" {
						// Fallback al nombre de la API si no está en Data Dragon (raro)
						name = p.ChampionName
					}
					champStats[p.ChampionID] = &ChampionMetaStat{ID: p.ChampionID, Name: name}
				}
				
				// Asegurar que tenemos el nombre correcto
				if champStats[p.ChampionID].Name == "Unknown" || champStats[p.ChampionID].Name == "" {
					name := champIDToName[p.ChampionID]
					if name != "" {
						champStats[p.ChampionID].Name = name
					} else {
						champStats[p.ChampionID].Name = p.ChampionName
					}
				}
				
				champStats[p.ChampionID].Games++
				if p.Win {
					champStats[p.ChampionID].Wins++
				}
			}
		}(matchId)
	}

	wg.Wait()

	// 5. Generar reporte
	report := &RegionMetaReport{
		Region:             platform,
		AnalyzedMatches:    totalGames,
		AnalyzedPlayers:    len(entries),
		TopChampions:       make([]ChampionMetaStat, 0),
		TopBans:            make([]ChampionMetaStat, 0),
		CommonCompositions: make([]CompositionStat, 0), // Inicializar vacío para evitar null
	}

	for _, stat := range champStats {
		if totalGames > 0 {
			stat.PickRate = (float64(stat.Games) / float64(totalGames*10)) * 100 // 10 jugadores por partida
			stat.BanRate = (float64(stat.Bans) / float64(totalGames)) * 100
			if stat.Games > 0 {
				stat.WinRate = (float64(stat.Wins) / float64(stat.Games)) * 100
			}
		}
		
		// Filtrar campeones con ID 0 o -1 (no ban)
		if stat.ID > 0 {
			report.TopChampions = append(report.TopChampions, *stat)
			report.TopBans = append(report.TopBans, *stat)
		}
	}

	// Ordenar Top Champions por Pick Rate
	sort.Slice(report.TopChampions, func(i, j int) bool {
		return report.TopChampions[i].PickRate > report.TopChampions[j].PickRate
	})

	// Ordenar Top Bans por Ban Rate
	sort.Slice(report.TopBans, func(i, j int) bool {
		return report.TopBans[i].BanRate > report.TopBans[j].BanRate
	})

	// Limitar a Top 20
	if len(report.TopChampions) > 20 {
		report.TopChampions = report.TopChampions[:20]
	}
	if len(report.TopBans) > 20 {
		report.TopBans = report.TopBans[:20]
	}

	return report, nil
}
