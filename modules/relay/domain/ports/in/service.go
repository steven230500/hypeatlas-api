package in

import (
	"context"

	"github.com/steven230500/hypeatlas-api/domain/entities"
)

type Service interface {
	// Relay básicos
	ListLiveCoStreams(ctx context.Context, eventID, lang string) ([]entities.CoStream, error)

	// Events
	ListEvents(ctx context.Context, game, league, status string, limit, offset int) ([]entities.Event, int, error)
	GetEvent(ctx context.Context, slug string) (*entities.Event, error)
	CreateEvent(ctx context.Context, event *entities.Event) error
	UpdateEvent(ctx context.Context, slug string, updates map[string]interface{}) error
	DeleteEvent(ctx context.Context, slug string) error

	// HypeMap
	HypeMapLive(ctx context.Context, game, lang string, limit, offset int) ([]entities.HypeMapItem, error)
	HypeMapSummary(ctx context.Context, game, lang string, limit, offset int) ([]entities.HypeMapSummaryItem, error)
}
