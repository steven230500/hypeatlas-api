package postgres

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/steven230500/hypeatlas-api/modules/signal/domain/entities"
	out "github.com/steven230500/hypeatlas-api/modules/signal/domain/ports/out"
)

type Repo struct{ db *pgxpool.Pool }

func New(db *pgxpool.Pool) out.Repository { return &Repo{db: db} }

func NewRaw(db *pgxpool.Pool) *Repo { return &Repo{db: db} }

func (r *Repo) PatchesByGame(ctx context.Context, game string) ([]entities.Patch, error) {
	const q = `
SELECT id, version,
       COALESCE(to_char(released_at, 'YYYY-MM-DD'), '')
FROM app.patches
WHERE game = $1
ORDER BY released_at DESC NULLS LAST, version DESC;`
	rows, err := r.db.Query(ctx, q, game)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var outv []entities.Patch
	for rows.Next() {
		var p entities.Patch
		p.Game = game
		if err := rows.Scan(&p.ID, &p.Version, &p.ReleasedAt); err != nil {
			return nil, err
		}
		outv = append(outv, p)
	}
	return outv, rows.Err()
}

func (r *Repo) Changes(ctx context.Context, game, version, entityType string) ([]entities.Change, error) {
	const q = `
SELECT pc.patch_id, pc.entity_type, pc.entity_id, pc.field, pc.old, pc.new, pc.impact_score
FROM app.patch_changes pc
JOIN app.patches p ON p.id = pc.patch_id
WHERE p.game = $1
  AND p.version = $2
  AND ($3 = '' OR pc.entity_type = $3)
ORDER BY pc.impact_score DESC, pc.id;`
	rows, err := r.db.Query(ctx, q, game, version, entityType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var outv []entities.Change
	for rows.Next() {
		var c entities.Change
		if err := rows.Scan(&c.PatchID, &c.EntityType, &c.EntityID, &c.Field, &c.Old, &c.New, &c.Impact); err != nil {
			return nil, err
		}
		outv = append(outv, c)
	}
	return outv, rows.Err()
}

func (r *Repo) Leagues(ctx context.Context, game, region string) ([]entities.League, error) {
	const q = `
SELECT id, game, region, name, slug
FROM app.leagues
WHERE game = $1 AND ($2 = '' OR region = $2)
ORDER BY region, name;`
	rows, err := r.db.Query(ctx, q, game, region)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var outv []entities.League
	for rows.Next() {
		var l entities.League
		if err := rows.Scan(&l.ID, &l.Game, &l.Region, &l.Name, &l.Slug); err != nil {
			return nil, err
		}
		outv = append(outv, l)
	}
	return outv, rows.Err()
}

func (r *Repo) Comps(ctx context.Context, game, region, league, patch, mapp, side string, limit int) ([]entities.Comp, error) {
	const q = `
SELECT id, game, region, COALESCE(league,''), patch, COALESCE(map,''), COALESCE(side,''), slots, pick_rate, win_rate, delta_win
FROM app.comps
WHERE game = $1
  AND region = $2
  AND patch  = $3
  AND ($4 = '' OR league = $4)
  AND ($5 = '' OR map    = $5)
  AND ($6 = '' OR side   = $6)
ORDER BY win_rate DESC NULLS LAST, pick_rate DESC NULLS LAST, id
LIMIT $7;`
	rows, err := r.db.Query(ctx, q, game, region, patch, league, mapp, side, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var outv []entities.Comp
	for rows.Next() {
		var c entities.Comp
		var raw []byte
		if err := rows.Scan(&c.ID, &c.Game, &c.Region, &c.League, &c.Patch, &c.Map, &c.Side, &raw, &c.PickRate, &c.WinRate, &c.DeltaWin); err != nil {
			return nil, err
		}
		if len(raw) > 0 {
			var tmp any
			_ = json.Unmarshal(raw, &tmp) // solo validar
			c.SlotsJSON = string(raw)
		}
		outv = append(outv, c)
	}
	return outv, rows.Err()
}
