package out

import (
	"context"

	"github.com/steven230500/hypeatlas-api/domain/entities"
)

type Repository interface {
	// Patches & Changes
	PatchesByGame(ctx context.Context, game string) ([]entities.Patch, error)
	PatchChanges(ctx context.Context, game, version, entityType string) ([]entities.PatchChange, error)

	// Leagues & Comps
	Leagues(ctx context.Context, game, region string) ([]entities.League, error)
	Comps(ctx context.Context, game, region, league, patch, mapp, side string, limit int) ([]entities.Comp, error)

	// Ingest
	UpsertComp(ctx context.Context, game, region, league, patch, mapp, side string, slotsJSON string, pickRate, winRate, deltaWin *float64) error
}
