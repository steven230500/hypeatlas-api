package out

import (
	"context"
	"time"

	"github.com/steven230500/hypeatlas-api/domain/entities"
)

type Repository interface {
	// CoStreams
	FindLiveByEvent(ctx context.Context, eventID, lang string) ([]entities.CoStream, error)

	// HypeMap
	HypeMapLive(ctx context.Context, game, lang string, limit, offset int) ([]entities.HypeMapItem, error)
	HypeMapSummary(ctx context.Context, game, lang string, limit, offset int) ([]entities.HypeMapSummaryItem, error)

	// Ingest / mantenimiento
	UpsertCoStream(ctx context.Context,
		eventSlug, eventTitle, game, league string,
		startsAtNullable *string,
		platform, handle, url, lang, country string,
		verified bool, viewers int, isLive bool,
	) error

	MarkStaleCoStreamsOffline(ctx context.Context, olderThan time.Duration) (int64, error)

	// Worker helpers
	LoadStreamRules(ctx context.Context) (map[string]string, error)
	ActiveWindows(ctx context.Context, now time.Time) ([]entities.EventWindow, error)
	ListCreatorHandles(ctx context.Context, platform string, verified bool) ([]entities.Creator, error)
}
