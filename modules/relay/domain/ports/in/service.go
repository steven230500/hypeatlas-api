package in

import (
	"context"

	"github.com/steven230500/hypeatlas-api/domain/entities"
)

type Service interface {
	// Relay básicos
	ListLiveCoStreams(ctx context.Context, eventID, lang string) ([]entities.CoStream, error)

	// HypeMap
	HypeMapLive(ctx context.Context, game, lang string, limit, offset int) ([]entities.HypeMapItem, error)
	HypeMapSummary(ctx context.Context, game, lang string, limit, offset int) ([]entities.HypeMapSummaryItem, error)
}
