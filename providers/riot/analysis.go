package riot

import (
	"fmt"
	"sort"
	"sync"
	"time"
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

	// Canal para limitar concurrencia (evitar rate limits agresivos)
	sem := make(chan struct{}, 5) // 5 concurrent requests

	fmt.Printf("Analyzing %d players from %s...\n", len(entries), platform)

	for _, entry := range entries {
		wg.Add(1)
		go func(entry LeagueEntry) {
			defer wg.Done()
			sem <- struct{}{} // Acquire token
			defer func() { <-sem }() // Release token

			// Obtener Match IDs
			// Usamos el PUUID directamente de la entrada de liga (ahorra una llamada a API)
			puuid := entry.Puuid
			if puuid == "" {
				// Fallback si por alguna razón no viene el PUUID (aunque debería)
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

	// Procesar partidas (también con concurrencia limitada)
	for i, matchId := range matchIds {
		// Rate limiting manual simple para no explotar la API key de desarrollo
		if i > 0 && i%10 == 0 {
			time.Sleep(1 * time.Second)
		}

		match, err := c.GetMatch(platform, matchId)
		if err != nil {
			fmt.Printf("Error getting match %s: %v\n", matchId, err)
			continue
		}

		totalGames++

		// Procesar Bans
		for _, team := range match.Info.Teams {
			for _, ban := range team.Bans {
				if _, exists := champStats[ban.ChampionID]; !exists {
					champStats[ban.ChampionID] = &ChampionMetaStat{ID: ban.ChampionID}
				}
				champStats[ban.ChampionID].Bans++
			}
		}

		// Procesar Picks y Wins
		for _, p := range match.Info.Participants {
			if _, exists := champStats[p.ChampionID]; !exists {
				champStats[p.ChampionID] = &ChampionMetaStat{ID: p.ChampionID, Name: p.ChampionName}
			}
			// Actualizar nombre si falta (por bans)
			if champStats[p.ChampionID].Name == "" {
				champStats[p.ChampionID].Name = p.ChampionName
			}
			
			champStats[p.ChampionID].Games++
			if p.Win {
				champStats[p.ChampionID].Wins++
			}
		}
	}

	// 5. Generar reporte
	report := &RegionMetaReport{
		Region:          platform,
		AnalyzedMatches: totalGames,
		AnalyzedPlayers: len(entries),
		TopChampions:    make([]ChampionMetaStat, 0),
		TopBans:         make([]ChampionMetaStat, 0),
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
