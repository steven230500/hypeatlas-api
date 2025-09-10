package service

import (
	"context"
	"errors"

	"github.com/steven230500/hypeatlas-api/modules/relay/domain/entities"
	in "github.com/steven230500/hypeatlas-api/modules/relay/domain/ports/in"
	out "github.com/steven230500/hypeatlas-api/modules/relay/domain/ports/out"
)

type svc struct{ repo out.Repository }

func New(r out.Repository) in.Service { return &svc{repo: r} }

func (s *svc) ListLiveCoStreams(ctx context.Context, eventID, lang string) ([]entities.CoStream, error) {
	if eventID == "" {
		return nil, errors.New("event_id required")
	}
	return s.repo.FindLiveByEvent(ctx, eventID, lang)
}
