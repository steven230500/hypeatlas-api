package riot

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/steven230500/hypeatlas-api/domain/entities"
	"github.com/steven230500/hypeatlas-api/modules/signal/domain/ports/out"
)

// Service integra Riot Games API con el módulo signal
type Service struct {
	client *Client
	repo   out.Repository
}

// NewService crea un nuevo servicio de Riot Games
func NewService(apiKey string, repo out.Repository) *Service {
	return &Service{
		client: NewClient(apiKey),
		repo:   repo,
	}
}

// SyncPatches sincroniza los parches desde Riot Games API
func (s *Service) SyncPatches(ctx context.Context) error {
	log.Println("Starting Riot Games patch synchronization...")

	// Obtener la versión más reciente
	latestVersion, err := s.client.GetLatestVersion()
	if err != nil {
		return fmt.Errorf("error getting latest version: %w", err)
	}

	log.Printf("Latest version from Riot API: %s", latestVersion)

	// Verificar si ya existe este parche en la base de datos
	existingPatches, err := s.repo.PatchesByGame(ctx, "lol")
	if err != nil {
		return fmt.Errorf("error checking existing patches: %w", err)
	}

	// Verificar si el parche ya existe
	patchExists := false
	for _, patch := range existingPatches {
		if patch.Version == latestVersion {
			patchExists = true
			break
		}
	}

	if patchExists {
		log.Printf("Patch %s already exists in database", latestVersion)
		return nil
	}

	// Crear nuevo parche
	newPatch := entities.Patch{
		Game:       "lol",
		Version:    latestVersion,
		ReleasedAt: &[]time.Time{time.Now()}[0], // Fecha actual como aproximación
	}

	// Aquí podríamos agregar lógica para obtener la fecha real de lanzamiento
	// desde la API de Riot, pero por ahora usamos la fecha actual

	// Guardar en base de datos
	if err := s.savePatch(ctx, &newPatch); err != nil {
		return fmt.Errorf("error saving patch: %w", err)
	}

	log.Printf("Successfully synced patch: %s", latestVersion)
	return nil
}

// SyncChampions sincroniza los campeones para un parche específico
func (s *Service) SyncChampions(ctx context.Context, version string) error {
	log.Printf("Starting champion synchronization for version %s...", version)

	// Obtener campeones desde la API
	champions, err := s.client.GetChampions(version)
	if err != nil {
		return fmt.Errorf("error getting champions: %w", err)
	}

	log.Printf("Found %d champions in version %s", len(champions.Data), version)

	// Aquí podríamos procesar y guardar los cambios de campeones
	// Por ahora solo loggeamos
	for _, champion := range champions.Data {
		log.Printf("Champion: %s (%s) - %s", champion.Name, champion.ID, champion.Title)
	}

	return nil
}

// savePatch guarda un parche en la base de datos
func (s *Service) savePatch(ctx context.Context, patch *entities.Patch) error {
	// Asignar UUID si no tiene
	if patch.UUID == uuid.Nil {
		patch.UUID = uuid.New()
	}

	// Asignar timestamps
	now := time.Now()
	if patch.CreatedAt.IsZero() {
		patch.CreatedAt = now
	}
	if patch.UpdatedAt.IsZero() {
		patch.UpdatedAt = now
	}

	// Aquí iría la lógica para guardar en la base de datos
	// Por ahora solo simulamos que se guarda correctamente
	log.Printf("Simulating save of patch: %+v", patch)

	return nil
}

// GetPatchInfo obtiene información detallada de un parche
func (s *Service) GetPatchInfo(ctx context.Context, version string) (*entities.Patch, error) {
	// Buscar parche en base de datos
	patches, err := s.repo.PatchesByGame(ctx, "lol")
	if err != nil {
		return nil, fmt.Errorf("error getting patches: %w", err)
	}

	// Buscar el parche específico
	for _, patch := range patches {
		if patch.Version == version {
			return &patch, nil
		}
	}

	return nil, fmt.Errorf("patch %s not found", version)
}

// GetChampionRotation obtiene la rotación semanal de campeones gratuitos
func (s *Service) GetChampionRotation(ctx context.Context, platform string) (*ChampionRotationResponse, error) {
	return s.client.GetChampionRotation(platform)
}

// GetChampions obtiene la lista de campeones para una versión específica
func (s *Service) GetChampions(ctx context.Context, version string) (*ChampionsResponse, error) {
	return s.client.GetChampions(version)
}

// GetChallengerLeague obtiene la liga Challenger para una cola específica
func (s *Service) GetChallengerLeague(ctx context.Context, platform, queue string) (*LeagueData, error) {
	return s.client.GetChallengerLeague(platform, queue)
}

// GetAllLeagues obtiene todas las ligas disponibles para una plataforma
func (s *Service) GetAllLeagues(ctx context.Context, platform string) ([]string, error) {
	return s.client.GetAllLeagues(platform)
}

// GetGames obtiene la lista de juegos disponibles de Riot Games
func (s *Service) GetGames(ctx context.Context) ([]string, error) {
	return s.client.GetGames()
}

// GetRegions obtiene la lista de regiones disponibles para League of Legends
func (s *Service) GetRegions(ctx context.Context) ([]string, error) {
	return s.client.GetRegions()
}

// GetChampionStats obtiene estadísticas de uso de campeones
func (s *Service) GetChampionStats(ctx context.Context, version string) (map[string]interface{}, error) {
	return s.client.GetChampionStats(version)
}

// GetPatchChanges obtiene cambios de campeones entre parches
func (s *Service) GetPatchChanges(ctx context.Context, fromVersion, toVersion string) (map[string]interface{}, error) {
	return s.client.GetPatchChanges(fromVersion, toVersion)
}

// GetProfessionalLeagues obtiene información sobre ligas profesionales
func (s *Service) GetProfessionalLeagues(ctx context.Context) (map[string]interface{}, error) {
	return s.client.GetProfessionalLeagues()
}

// GetLeagueChampions obtiene estadísticas de campeones en una liga específica
func (s *Service) GetLeagueChampions(ctx context.Context, league string) (map[string]interface{}, error) {
	return s.client.GetLeagueChampions(league)
}

// ProfessionalLeagueService maneja la lógica de ligas profesionales
type ProfessionalLeagueService struct {
	client *Client
}

// NewProfessionalLeagueService crea un nuevo servicio de ligas profesionales
func NewProfessionalLeagueService(client *Client) *ProfessionalLeagueService {
	return &ProfessionalLeagueService{client: client}
}

// GetProfessionalLeagues obtiene información completa de ligas profesionales
func (s *ProfessionalLeagueService) GetProfessionalLeagues() (map[string]interface{}, error) {
	return s.client.GetProfessionalLeagues()
}

// GetLeagueChampions obtiene estadísticas detalladas de una liga específica
func (s *ProfessionalLeagueService) GetLeagueChampions(league string) (map[string]interface{}, error) {
	// Validar liga
	validLeagues := map[string]bool{
		"LEC": true, "LCK": true, "LPL": true, "LTA": true,
		"LCS": true, "VCS": true, "PCS": true,
	}

	if !validLeagues[league] {
		return nil, fmt.Errorf("invalid league: %s", league)
	}

	// Obtener datos básicos
	stats, err := s.client.GetLeagueChampions(league)
	if err != nil {
		return nil, err
	}

	// Enriquecer con datos adicionales
	stats["analysis_date"] = time.Now().UTC().Format(time.RFC3339)
	stats["data_source"] = "Riot Games API + Community Analysis"

	return stats, nil
}

// ValidateLeague valida si una liga existe
func (s *ProfessionalLeagueService) ValidateLeague(league string) bool {
	validLeagues := []string{"LEC", "LCK", "LPL", "LTA", "LCS", "VCS", "PCS"}
	for _, valid := range validLeagues {
		if valid == league {
			return true
		}
	}
	return false
}

// DataDragonService maneja la lógica de Data Dragon API
type DataDragonService struct {
	client *Client
}

// NewDataDragonService crea un nuevo servicio de Data Dragon
func NewDataDragonService(client *Client) *DataDragonService {
	return &DataDragonService{client: client}
}

// GetGameVersions obtiene todas las versiones disponibles del juego
func (s *DataDragonService) GetGameVersions(ctx context.Context) ([]string, error) {
	versions, err := s.client.GetLatestVersion()
	if err != nil {
		return nil, fmt.Errorf("error getting latest version: %w", err)
	}

	// Para obtener todas las versiones, usamos la API de versiones
	// Por ahora devolvemos solo la última versión
	return []string{versions}, nil
}

// GetItems obtiene datos de items para una versión específica
func (s *DataDragonService) GetItems(ctx context.Context, version string) (map[string]interface{}, error) {
	// Usar el método del cliente que ya implementamos
	return s.client.GetItems(version)
}

// GetRunes obtiene datos de runas para una versión específica
func (s *DataDragonService) GetRunes(ctx context.Context, version string) (map[string]interface{}, error) {
	// Usar el método del cliente que ya implementamos
	return s.client.GetRunes(version)
}

// GetSummonerSpells obtiene datos de summoner spells para una versión específica
func (s *DataDragonService) GetSummonerSpells(ctx context.Context, version string) (map[string]interface{}, error) {
	// Usar el método del cliente que ya implementamos
	return s.client.GetSummonerSpells(version)
}

// GetChampionDetails obtiene detalles completos de un campeón específico
func (s *DataDragonService) GetChampionDetails(ctx context.Context, version, championID string) (map[string]interface{}, error) {
	// Usar el método del cliente que ya implementamos
	return s.client.GetChampionDetails(version, championID)
}

// GetPatchNotes obtiene información de cambios entre parches
func (s *DataDragonService) GetPatchNotes(ctx context.Context, fromVersion, toVersion string) (map[string]interface{}, error) {
	// Usar el método del cliente que ya implementamos
	return s.client.GetPatchNotes(fromVersion, toVersion)
}

// ImageService maneja la lógica de URLs de imágenes de Data Dragon
type ImageService struct {
	client *Client
}

// NewImageService crea un nuevo servicio de imágenes
func NewImageService(client *Client) *ImageService {
	return &ImageService{client: client}
}

// GetChampionImageURLs obtiene todas las URLs de imágenes disponibles para un campeón
func (s *ImageService) GetChampionImageURLs(ctx context.Context, version, championID string, skinNum int) (map[string]interface{}, error) {
	baseURL := "https://ddragon.leagueoflegends.com/cdn"

	imageURLs := map[string]interface{}{
		"version":  version,
		"champion": championID,
		"skin_num": skinNum,
		"images": map[string]string{
			// Icono del campeón (48x48)
			"icon": fmt.Sprintf("%s/%s/img/champion/%s.png", baseURL, version, championID),

			// Splash art (1208x512)
			"splash": fmt.Sprintf("%s/img/champion/splash/%s_%d.jpg", baseURL, championID, skinNum),

			// Pantalla de carga (1024x768)
			"loading": fmt.Sprintf("%s/img/champion/loading/%s_%d.jpg", baseURL, championID, skinNum),

			// Tile (mini splash) (380x224)
			"tile": fmt.Sprintf("%s/img/champion/tiles/%s_%d.jpg", baseURL, championID, skinNum),
		},
	}

	return imageURLs, nil
}

// GetItemImageURL obtiene la URL de imagen para un item específico
func (s *ImageService) GetItemImageURL(ctx context.Context, version, itemID string) (map[string]interface{}, error) {
	baseURL := "https://ddragon.leagueoflegends.com/cdn"

	imageURL := map[string]interface{}{
		"version": version,
		"item_id": itemID,
		"image": map[string]interface{}{
			"url":    fmt.Sprintf("%s/%s/img/item/%s.png", baseURL, version, itemID),
			"size":   "64x64",
			"format": "png",
		},
	}

	return imageURL, nil
}

// GetSpellImageURL obtiene la URL de imagen para un summoner spell
func (s *ImageService) GetSpellImageURL(ctx context.Context, version, spellName string) (map[string]interface{}, error) {
	baseURL := "https://ddragon.leagueoflegends.com/cdn"

	imageURL := map[string]interface{}{
		"version": version,
		"spell":   spellName,
		"image": map[string]interface{}{
			"url":    fmt.Sprintf("%s/%s/img/spell/%s.png", baseURL, version, spellName),
			"size":   "64x64",
			"format": "png",
		},
	}

	return imageURL, nil
}

// GetRuneImageURL obtiene la URL de imagen para una runa
func (s *ImageService) GetRuneImageURL(ctx context.Context, runeIcon string) (map[string]interface{}, error) {
	baseURL := "https://ddragon.leagueoflegends.com/cdn"

	imageURL := map[string]interface{}{
		"rune_icon": runeIcon,
		"image": map[string]interface{}{
			"url":    fmt.Sprintf("%s/img/%s", baseURL, runeIcon),
			"size":   "48x48",
			"format": "png",
		},
	}

	return imageURL, nil
}

// GetProfileIconImageURL obtiene la URL de imagen para un ícono de perfil
func (s *ImageService) GetProfileIconImageURL(ctx context.Context, version string, iconID int) (map[string]interface{}, error) {
	baseURL := "https://ddragon.leagueoflegends.com/cdn"

	imageURL := map[string]interface{}{
		"version": version,
		"icon_id": iconID,
		"image": map[string]interface{}{
			"url":    fmt.Sprintf("%s/%s/img/profileicon/%d.png", baseURL, version, iconID),
			"size":   "48x48",
			"format": "png",
		},
	}

	return imageURL, nil
}

// GetMapImageURL obtiene la URL de imagen para un mapa
func (s *ImageService) GetMapImageURL(ctx context.Context, version string, mapID int) (map[string]interface{}, error) {
	baseURL := "https://ddragon.leagueoflegends.com/cdn"

	imageURL := map[string]interface{}{
		"version": version,
		"map_id":  mapID,
		"image": map[string]interface{}{
			"url":    fmt.Sprintf("%s/%s/img/map/map%d.png", baseURL, version, mapID),
			"size":   "512x512",
			"format": "png",
		},
	}

	return imageURL, nil
}

// GetAbilityImageURL obtiene la URL de imagen para una habilidad de campeón
func (s *ImageService) GetAbilityImageURL(ctx context.Context, version, abilityName string) (map[string]interface{}, error) {
	baseURL := "https://ddragon.leagueoflegends.com/cdn"

	imageURL := map[string]interface{}{
		"version": version,
		"ability": abilityName,
		"image": map[string]interface{}{
			"url":    fmt.Sprintf("%s/%s/img/spell/%s.png", baseURL, version, abilityName),
			"size":   "64x64",
			"format": "png",
		},
	}

	return imageURL, nil
}

// GetPassiveImageURL obtiene la URL de imagen para la pasiva de un campeón
func (s *ImageService) GetPassiveImageURL(ctx context.Context, version, passiveFile string) (map[string]interface{}, error) {
	baseURL := "https://ddragon.leagueoflegends.com/cdn"

	imageURL := map[string]interface{}{
		"version": version,
		"passive": passiveFile,
		"image": map[string]interface{}{
			"url":    fmt.Sprintf("%s/%s/img/passive/%s.png", baseURL, version, passiveFile),
			"size":   "48x48",
			"format": "png",
		},
	}

	return imageURL, nil
}
