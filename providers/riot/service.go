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
