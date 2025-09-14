package repository

import (
	"context"

	"github.com/steven230500/hypeatlas-api/domain/entities"
	out "github.com/steven230500/hypeatlas-api/modules/signal/domain/ports/out"
	"github.com/steven230500/hypeatlas-api/shared/db"
	"gorm.io/gorm"
)

type Repo struct{ db *gorm.DB }

func New(db *gorm.DB) out.Repository { return &Repo{db: db} }

func (r *Repo) PatchesByGame(ctx context.Context, game string) ([]entities.Patch, error) {
	var patches []entities.Patch
	result := db.Call(
		r.db.WithContext(ctx).
			Where("game = ?", game).
			Order("released_at DESC NULLS LAST, version DESC").
			Find(&patches),
	)
	return patches, result.Error
}

func (r *Repo) PatchChanges(ctx context.Context, game, version, entityType string) ([]entities.PatchChange, error) {
	var changes []entities.PatchChange
	query := r.db.WithContext(ctx).Joins("JOIN app.patches p ON p.uuid = app.patch_changes.patch_uuid").Where("p.game = ? AND p.version = ?", game, version)
	if entityType != "" {
		query = query.Where("app.patch_changes.entity_type = ?", entityType)
	}
	result := db.Call(query.Order("app.patch_changes.impact_score DESC, app.patch_changes.id").Find(&changes))
	return changes, result.Error
}

func (r *Repo) Leagues(ctx context.Context, game, region string) ([]entities.League, error) {
	var leagues []entities.League
	query := r.db.WithContext(ctx).Where("game = ?", game)
	if region != "" {
		query = query.Where("region = ?", region)
	}
	result := db.Call(query.Order("region, name").Find(&leagues))
	return leagues, result.Error
}

func (r *Repo) Comps(ctx context.Context, game, region, league, patch, mapp, side string, limit int) ([]entities.Comp, error) {
	var comps []entities.Comp
	query := r.db.WithContext(ctx).Where("game = ? AND region = ? AND patch = ?", game, region, patch)
	if league != "" {
		query = query.Where("league = ?", league)
	}
	if mapp != "" {
		query = query.Where("map = ?", mapp)
	}
	if side != "" {
		query = query.Where("side = ?", side)
	}
	result := db.Call(query.Order("win_rate DESC NULLS LAST, pick_rate DESC NULLS LAST, uuid").Limit(limit).Find(&comps))
	return comps, result.Error
}

// UpsertComp inserta/actualiza una composición.
// Usa la UNIQUE (game,region,league,patch,map,side,slots_fp).
func (r *Repo) UpsertComp(
	ctx context.Context,
	game, region, league, patch, mapp, side string,
	slotsJSON string, // JSON en texto
	pickRate, winRate, deltaWin *float64, // pueden ser nil
) error {
	// language=SQL
	const q = `
INSERT INTO app.comps
  (game, region, league, patch, map, side, slots, pick_rate, win_rate, delta_win)
VALUES (?, ?, ?, ?, ?, ?, ?::jsonb, ?, ?, ?)
ON CONFLICT (game, region, league, patch, map, side, slots_fp)
DO UPDATE SET
  pick_rate  = EXCLUDED.pick_rate,
  win_rate   = EXCLUDED.win_rate,
  delta_win  = EXCLUDED.delta_win,
  updated_at = now();
`
	result := db.Call(r.db.WithContext(ctx).Exec(q,
		game, region, league, patch, mapp, side, slotsJSON, pickRate, winRate, deltaWin,
	))
	return result.Error
}
