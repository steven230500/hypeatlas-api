package in

import (
	"context"

	"github.com/steven230500/hypeatlas-api/domain/entities"
)

type Service interface {
	// Patches & Changes
	ListPatches(ctx context.Context, game string) ([]entities.Patch, error)
	ListChanges(ctx context.Context, game, version, entityType string) ([]entities.PatchChange, error)

	// Leagues & Comps
	ListLeagues(ctx context.Context, game, region string) ([]entities.League, error)
	ListComps(ctx context.Context, game, region, league, patch, mapp, side string, limit int) ([]entities.Comp, error)
}
