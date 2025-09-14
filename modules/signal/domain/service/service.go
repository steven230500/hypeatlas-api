package service

import (
	"context"
	"errors"

	"github.com/steven230500/hypeatlas-api/domain/entities"
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

func (s *svc) ListChanges(ctx context.Context, game, version, entityType string) ([]entities.PatchChange, error) {
	if game == "" || version == "" {
		return nil, errors.New("game and version required")
	}
	return s.repo.PatchChanges(ctx, game, version, entityType)
}

func (s *svc) ListLeagues(ctx context.Context, game, region string) ([]entities.League, error) {
	if game == "" {
		return nil, errors.New("game required")
	}
	return s.repo.Leagues(ctx, game, region)
}

func (s *svc) ListComps(ctx context.Context, game, region, league, patch, mapp, side string, limit int) ([]entities.Comp, error) {
	if game == "" || region == "" || patch == "" {
		return nil, errors.New("game, region and patch required")
	}
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	return s.repo.Comps(ctx, game, region, league, patch, mapp, side, limit)
}
