package service

import (
	"context"
	"fmt"
	"time"

	out "github.com/steven230500/hypeatlas-api/modules/signal/domain/ports/out"
	riot "github.com/steven230500/hypeatlas-api/providers/riot"
)

// MetaGameService implementa análisis de meta-game
type MetaGameService struct {
	repo       out.Repository
	riotSvc    *riot.Service
	dataDragon *riot.Client
}

// NewMetaGameService crea un nuevo servicio de análisis de meta-game
func NewMetaGameService(repo out.Repository, riotSvc *riot.Service) *MetaGameService {
	return &MetaGameService{
		repo:       repo,
		riotSvc:    riotSvc,
		dataDragon: riot.NewClient(""), // Cliente sin API key para Data Dragon
	}
}

// ChampionRotationAnalysis representa el análisis de una rotación de campeones
type ChampionRotationAnalysis struct {
	FreeChampions        []ChampionInfo `json:"free_champions"`
	NewPlayerChampions   []ChampionInfo `json:"new_player_champions"`
	ImpactScore          float64        `json:"impact_score"`
	MetaShiftProbability float64        `json:"meta_shift_probability"`
	RecommendedChampions []string       `json:"recommended_champions"`
}

// ChampionInfo información básica de un campeón
type ChampionInfo struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Role     string  `json:"role"`
	Tier     string  `json:"tier"`
	PickRate float64 `json:"pick_rate"`
	WinRate  float64 `json:"win_rate"`
	BanRate  float64 `json:"ban_rate"`
}

// AnalyzeChampionRotation analiza el impacto de la rotación semanal de campeones
func (s *MetaGameService) AnalyzeChampionRotation(ctx context.Context, platform string) (*ChampionRotationAnalysis, error) {
	// Obtener rotación actual desde Riot API
	rotation, err := s.riotSvc.GetChampionRotation(ctx, platform)
	if err != nil {
		return nil, fmt.Errorf("error getting champion rotation: %w", err)
	}

	// Obtener la versión más reciente desde Data Dragon
	latestVersion, err := s.dataDragon.GetLatestVersion()
	if err != nil {
		return nil, fmt.Errorf("error getting latest version: %w", err)
	}

	// Obtener datos de campeones desde Data Dragon (sin autenticación)
	champions, err := s.dataDragon.GetChampions(latestVersion)
	if err != nil {
		return nil, fmt.Errorf("error getting champions data: %w", err)
	}

	// Analizar campeones gratuitos
	freeChampions := s.analyzeChampions(rotation.FreeChampionIDs, champions)
	newPlayerChampions := s.analyzeChampions(rotation.FreeChampionIDsForNewPlayers, champions)

	// Calcular impacto en el meta
	impactScore := s.calculateRotationImpact(freeChampions, newPlayerChampions)
	metaShiftProb := s.calculateMetaShiftProbability(freeChampions)

	// Recomendaciones basadas en el análisis
	recommendations := s.generateRecommendations(freeChampions, newPlayerChampions)

	return &ChampionRotationAnalysis{
		FreeChampions:        freeChampions,
		NewPlayerChampions:   newPlayerChampions,
		ImpactScore:          impactScore,
		MetaShiftProbability: metaShiftProb,
		RecommendedChampions: recommendations,
	}, nil
}

// analyzeChampions analiza una lista de IDs de campeones
func (s *MetaGameService) analyzeChampions(championIDs []int, allChampions *riot.ChampionsResponse) []ChampionInfo {
	var result []ChampionInfo

	// Crear un mapa de key -> champion data para búsqueda rápida
	keyToChampion := make(map[string]riot.ChampionData)
	for _, champData := range allChampions.Data {
		keyToChampion[champData.Key] = champData
	}

	for _, id := range championIDs {
		// Buscar campeón por key (el ID numérico convertido a string)
		idStr := fmt.Sprintf("%d", id)

		if champData, exists := keyToChampion[idStr]; exists {
			// Estimar tier basado en el nombre del campeón
			tier := s.estimateChampionTier(champData.Name)
			pickRate := s.estimatePickRate(tier)
			winRate := s.estimateWinRate(tier)
			banRate := s.estimateBanRate(tier)

			result = append(result, ChampionInfo{
				ID:       id,
				Name:     champData.Name,
				Role:     s.estimateChampionRole(champData.Name),
				Tier:     tier,
				PickRate: pickRate,
				WinRate:  winRate,
				BanRate:  banRate,
			})
		} else {
			// Si no encontramos el campeón, crear una entrada básica
			result = append(result, ChampionInfo{
				ID:       id,
				Name:     fmt.Sprintf("Champion_%d", id),
				Role:     "Unknown",
				Tier:     "B",
				PickRate: 5.0,
				WinRate:  47.5,
				BanRate:  3.0,
			})
		}
	}

	return result
}

// estimateChampionTier estima el tier de un campeón basado en su nombre
func (s *MetaGameService) estimateChampionTier(championName string) string {
	// Campeones S-tier actuales (simplificado)
	sTierChampions := []string{
		"Yuumi", "Jinx", "Zeri", "Kai'Sa", "Jhin", "Vayne", "Caitlyn",
		"Lee Sin", "Akali", "Yone", "Fiora", "Jax", "Irelia",
	}

	// Campeones A-tier
	aTierChampions := []string{
		"Lux", "Ahri", "Miss Fortune", "Ezreal", "Thresh", "Leona",
		"Blitzcrank", "Zed", "Katarina", "Master Yi", "Darius",
	}

	for _, champ := range sTierChampions {
		if champ == championName {
			return "S"
		}
	}

	for _, champ := range aTierChampions {
		if champ == championName {
			return "A"
		}
	}

	return "B" // Default tier
}

// estimateChampionRole estima el rol principal de un campeón basado en su nombre
func (s *MetaGameService) estimateChampionRole(championName string) string {
	// Campeones ADC/Marksman
	adcChampions := []string{
		"Jinx", "Kai'Sa", "Jhin", "Vayne", "Caitlyn", "Miss Fortune",
		"Ezreal", "Lucian", "Ashe", "Draven", "Sivir", "Twitch",
	}

	// Campeones Assassin
	assassinChampions := []string{
		"Akali", "Katarina", "Zed", "Fiora", "Irelia", "Yone",
		"Qiyana", "Nocturne", "Rengar", "Master Yi", "Jax",
	}

	// Campeones Mage
	mageChampions := []string{
		"Ahri", "Lux", "Ziggs", "Veigar", "Karthus", "Anivia",
		"Brand", "Cassiopeia", "Kennen", "Orianna", "Swain",
	}

	// Campeones Tank
	tankChampions := []string{
		"Leona", "Thresh", "Nautilus", "Shen", "Malphite", "Maokai",
		"Poppy", "Rell", "Taric", "Yuumi", "Zac",
	}

	// Campeones Support
	supportChampions := []string{
		"Blitzcrank", "Janna", "Lulu", "Morgana", "Nami", "Soraka",
		"Zilean", "Bard", "Rakan", "Senna", "Yuumi",
	}

	for _, champ := range adcChampions {
		if champ == championName {
			return "ADC"
		}
	}

	for _, champ := range assassinChampions {
		if champ == championName {
			return "Assassin"
		}
	}

	for _, champ := range mageChampions {
		if champ == championName {
			return "Mage"
		}
	}

	for _, champ := range tankChampions {
		if champ == championName {
			return "Tank"
		}
	}

	for _, champ := range supportChampions {
		if champ == championName {
			return "Support"
		}
	}

	return "Fighter" // Default role
}

// estimatePickRate estima el pick rate basado en el tier
func (s *MetaGameService) estimatePickRate(tier string) float64 {
	switch tier {
	case "S":
		return 15.5
	case "A":
		return 12.3
	case "B":
		return 8.7
	default:
		return 5.2
	}
}

// estimateWinRate estima el win rate basado en el tier
func (s *MetaGameService) estimateWinRate(tier string) float64 {
	switch tier {
	case "S":
		return 52.1
	case "A":
		return 50.8
	case "B":
		return 49.2
	default:
		return 47.5
	}
}

// estimateBanRate estima el ban rate basado en el tier
func (s *MetaGameService) estimateBanRate(tier string) float64 {
	switch tier {
	case "S":
		return 25.3
	case "A":
		return 18.7
	case "B":
		return 12.1
	default:
		return 6.8
	}
}

// getPrimaryRole obtiene el rol primario de un campeón (deprecated - usar estimateChampionRole)
func (s *MetaGameService) getPrimaryRole(tags []string) string {
	return "Unknown"
}

// calculateRotationImpact calcula el impacto de la rotación en el meta
func (s *MetaGameService) calculateRotationImpact(free, newPlayer []ChampionInfo) float64 {
	totalImpact := 0.0

	// Campeones gratuitos tienen mayor impacto
	for _, champ := range free {
		switch champ.Tier {
		case "S":
			totalImpact += 10.0
		case "A":
			totalImpact += 7.5
		case "B":
			totalImpact += 5.0
		default:
			totalImpact += 2.5
		}
	}

	// Campeones para nuevos jugadores tienen impacto moderado
	for _, champ := range newPlayer {
		switch champ.Tier {
		case "S":
			totalImpact += 5.0
		case "A":
			totalImpact += 3.5
		case "B":
			totalImpact += 2.0
		default:
			totalImpact += 1.0
		}
	}

	return totalImpact
}

// calculateMetaShiftProbability calcula la probabilidad de cambio en el meta
func (s *MetaGameService) calculateMetaShiftProbability(champions []ChampionInfo) float64 {
	sTierCount := 0
	for _, champ := range champions {
		if champ.Tier == "S" {
			sTierCount++
		}
	}

	// Más campeones S-tier = mayor probabilidad de cambio
	baseProb := float64(sTierCount) * 15.0

	// Factor adicional por diversidad de roles
	roles := make(map[string]bool)
	for _, champ := range champions {
		roles[champ.Role] = true
	}

	diversityBonus := float64(len(roles)) * 5.0

	return min(baseProb+diversityBonus, 95.0)
}

// generateRecommendations genera recomendaciones basadas en el análisis
func (s *MetaGameService) generateRecommendations(free, newPlayer []ChampionInfo) []string {
	var recommendations []string

	// Recomendar campeones S-tier gratuitos
	for _, champ := range free {
		if champ.Tier == "S" {
			recommendations = append(recommendations,
				fmt.Sprintf("Practice %s (%s) - High impact champion now free",
					champ.Name, champ.Role))
		}
	}

	// Recomendar contra picks
	roleCount := make(map[string]int)
	for _, champ := range free {
		roleCount[champ.Role]++
	}

	for role, count := range roleCount {
		if count >= 2 {
			recommendations = append(recommendations,
				fmt.Sprintf("Watch %s counters - Multiple %s champions free",
					role, role))
		}
	}

	// Limitar a top 5 recomendaciones
	if len(recommendations) > 5 {
		recommendations = recommendations[:5]
	}

	return recommendations
}

// AnalyzeLeagueRankings analiza las estadísticas de una liga
func (s *MetaGameService) AnalyzeLeagueRankings(ctx context.Context, platform, queue string) (*LeagueAnalysis, error) {
	// Obtener datos de Challenger
	challenger, err := s.riotSvc.GetChallengerLeague(ctx, platform, queue)
	if err != nil {
		return nil, fmt.Errorf("error getting challenger league: %w", err)
	}

	// Calcular estadísticas
	totalPlayers := len(challenger.Entries)
	avgLP := s.calculateAverageLP(challenger.Entries)
	avgWinRate := s.calculateAverageWinRate(challenger.Entries)
	topLP := s.findTopLP(challenger.Entries)

	return &LeagueAnalysis{
		Queue:        queue,
		Tier:         challenger.Tier,
		TotalPlayers: totalPlayers,
		AvgLP:        avgLP,
		AvgWinRate:   avgWinRate,
		TopLP:        topLP,
		LeagueName:   challenger.Name,
	}, nil
}

// LeagueAnalysis representa el análisis de una liga
type LeagueAnalysis struct {
	Queue        string  `json:"queue"`
	Tier         string  `json:"tier"`
	TotalPlayers int     `json:"total_players"`
	AvgLP        float64 `json:"avg_lp"`
	AvgWinRate   float64 `json:"avg_win_rate"`
	TopLP        int     `json:"top_lp"`
	LeagueName   string  `json:"league_name"`
}

// calculateAverageLP calcula el LP promedio
func (s *MetaGameService) calculateAverageLP(entries []riot.LeagueEntry) float64 {
	if len(entries) == 0 {
		return 0
	}

	total := 0
	for _, entry := range entries {
		total += entry.LeaguePoints
	}

	return float64(total) / float64(len(entries))
}

// calculateAverageWinRate calcula el win rate promedio
func (s *MetaGameService) calculateAverageWinRate(entries []riot.LeagueEntry) float64 {
	if len(entries) == 0 {
		return 0
	}

	totalGames := 0
	totalWins := 0

	for _, entry := range entries {
		totalGames += entry.Wins + entry.Losses
		totalWins += entry.Wins
	}

	if totalGames == 0 {
		return 0
	}

	return (float64(totalWins) / float64(totalGames)) * 100
}

// findTopLP encuentra el LP más alto
func (s *MetaGameService) findTopLP(entries []riot.LeagueEntry) int {
	if len(entries) == 0 {
		return 0
	}

	maxLP := 0
	for _, entry := range entries {
		if entry.LeaguePoints > maxLP {
			maxLP = entry.LeaguePoints
		}
	}

	return maxLP
}

// GenerateMetaReport genera un reporte completo de meta-game
func (s *MetaGameService) GenerateMetaReport(ctx context.Context, platform string) (*MetaReport, error) {
	report := &MetaReport{
		Platform:    platform,
		GeneratedAt: time.Now(),
	}

	// Análisis de rotación de campeones
	rotationAnalysis, err := s.AnalyzeChampionRotation(ctx, platform)
	if err != nil {
		return nil, fmt.Errorf("error analyzing champion rotation: %w", err)
	}
	report.ChampionRotation = rotationAnalysis

	// Análisis de ligas
	soloQAnalysis, err := s.AnalyzeLeagueRankings(ctx, platform, "RANKED_SOLO_5x5")
	if err != nil {
		return nil, fmt.Errorf("error analyzing solo queue: %w", err)
	}
	report.LeagueAnalysis = soloQAnalysis

	// Generar insights
	report.Insights = s.generateMetaInsights(rotationAnalysis, soloQAnalysis)

	return report, nil
}

// MetaReport representa un reporte completo de meta-game
type MetaReport struct {
	Platform         string                    `json:"platform"`
	GeneratedAt      time.Time                 `json:"generated_at"`
	ChampionRotation *ChampionRotationAnalysis `json:"champion_rotation"`
	LeagueAnalysis   *LeagueAnalysis           `json:"league_analysis"`
	Insights         []string                  `json:"insights"`
}

// generateMetaInsights genera insights basados en el análisis
func (s *MetaGameService) generateMetaInsights(rotation *ChampionRotationAnalysis, league *LeagueAnalysis) []string {
	var insights []string

	// Insight sobre impacto de rotación
	if rotation.ImpactScore > 50 {
		insights = append(insights, "High impact champion rotation - expect meta shifts")
	} else if rotation.ImpactScore > 25 {
		insights = append(insights, "Moderate impact rotation - some champions may see increased play")
	} else {
		insights = append(insights, "Low impact rotation - minimal meta changes expected")
	}

	// Insight sobre win rate promedio
	if league.AvgWinRate > 52 {
		insights = append(insights, "High skill ceiling in Challenger - focus on mechanical champions")
	} else if league.AvgWinRate > 50 {
		insights = append(insights, "Balanced meta - team composition is key")
	} else {
		insights = append(insights, "Snowball meta - early game decisions are crucial")
	}

	// Insight sobre campeones S-tier gratuitos
	sTierFree := 0
	for _, champ := range rotation.FreeChampions {
		if champ.Tier == "S" {
			sTierFree++
		}
	}

	if sTierFree > 0 {
		insights = append(insights, fmt.Sprintf("%d S-tier champions are free - practice opportunities abound", sTierFree))
	}

	return insights
}

// Helper function
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
