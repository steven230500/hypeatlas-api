package postgres

import (
	"context"
)

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
VALUES ($1, $2, $3, $4, $5, $6, $7::jsonb, $8, $9, $10)
ON CONFLICT (game, region, league, patch, map, side, slots_fp)
DO UPDATE SET
  pick_rate  = EXCLUDED.pick_rate,
  win_rate   = EXCLUDED.win_rate,
  delta_win  = EXCLUDED.delta_win,
  created_at = now();
`
	_, err := r.db.Exec(ctx, q,
		game, region, league, patch, mapp, side, slotsJSON, pickRate, winRate, deltaWin,
	)
	return err
}
