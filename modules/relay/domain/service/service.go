package service

import (
	"context"
	"errors"

	"github.com/steven230500/hypeatlas-api/domain/entities"
	inport "github.com/steven230500/hypeatlas-api/modules/relay/domain/ports/in"
	outport "github.com/steven230500/hypeatlas-api/modules/relay/domain/ports/out"
)

type svc struct{ repo outport.Repository }

func New(r outport.Repository) inport.Service { return &svc{repo: r} }

func (s *svc) ListLiveCoStreams(ctx context.Context, eventID, lang string) ([]entities.CoStream, error) {
	if eventID == "" {
		return nil, errors.New("event_id required")
	}
	return s.repo.FindLiveByEvent(ctx, eventID, lang)
}

// ListEvents lista eventos con filtros opcionales
func (s *svc) ListEvents(ctx context.Context, game, league, status string, limit, offset int) ([]entities.Event, int, error) {
	filters := make(map[string]interface{})
	if game != "" {
		filters["game"] = game
	}
	if league != "" {
		filters["league"] = league
	}
	if status != "" {
		filters["status"] = status
	}
	return s.repo.FindEvents(ctx, filters, limit, offset)
}

// GetEvent obtiene un evento por su slug
func (s *svc) GetEvent(ctx context.Context, slug string) (*entities.Event, error) {
	if slug == "" {
		return nil, errors.New("slug required")
	}
	return s.repo.FindEventBySlug(ctx, slug)
}

// CreateEvent crea un nuevo evento
func (s *svc) CreateEvent(ctx context.Context, event *entities.Event) error {
	if event.Slug == "" {
		return errors.New("slug required")
	}
	if event.Title == "" {
		return errors.New("title required")
	}
	if event.Game == "" {
		return errors.New("game required")
	}
	return s.repo.CreateEvent(ctx, event)
}

// UpdateEvent actualiza un evento existente
func (s *svc) UpdateEvent(ctx context.Context, slug string, updates map[string]interface{}) error {
	if slug == "" {
		return errors.New("slug required")
	}
	if len(updates) == 0 {
		return errors.New("no updates provided")
	}
	return s.repo.UpdateEvent(ctx, slug, updates)
}

// DeleteEvent elimina un evento
func (s *svc) DeleteEvent(ctx context.Context, slug string) error {
	if slug == "" {
		return errors.New("slug required")
	}
	return s.repo.DeleteEvent(ctx, slug)
}

func (s *svc) HypeMapLive(ctx context.Context, game, lang string, limit, offset int) ([]entities.HypeMapItem, error) {
	return s.repo.HypeMapLive(ctx, game, lang, limit, offset)
}

func (s *svc) HypeMapSummary(ctx context.Context, game, lang string, limit, offset int) ([]entities.HypeMapSummaryItem, error) {
	return s.repo.HypeMapSummary(ctx, game, lang, limit, offset)
}
