package service

import (
	"context"
	"errors"

	"github.com/steven230500/hypeatlas-api/modules/signal/domain/entities"
	in "github.com/steven230500/hypeatlas-api/modules/signal/domain/ports/in"
	out "github.com/steven230500/hypeatlas-api/modules/signal/domain/ports/out"
)

type svc struct{ repo out.Repository }

func New(r out.Repository) in.Service { return &svc{repo: r} }

func (s *svc) ListPatches(ctx context.Context, game string) ([]entities.Patch, error) {
	if game == "" {
		return nil, errors.New("game required")
	}
	return s.repo.PatchesByGame(ctx, game)
}

func (s *svc) ListChanges(ctx context.Context, game, version, entityType string) ([]entities.Change, error) {
	if game == "" || version == "" {
		return nil, errors.New("game and version required")
	}
	return s.repo.Changes(ctx, game, version, entityType)
}
