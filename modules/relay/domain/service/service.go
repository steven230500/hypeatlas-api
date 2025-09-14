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

func (s *svc) HypeMapLive(ctx context.Context, game, lang string, limit, offset int) ([]entities.HypeMapItem, error) {
	return s.repo.HypeMapLive(ctx, game, lang, limit, offset)
}

func (s *svc) HypeMapSummary(ctx context.Context, game, lang string, limit, offset int) ([]entities.HypeMapSummaryItem, error) {
	return s.repo.HypeMapSummary(ctx, game, lang, limit, offset)
}
