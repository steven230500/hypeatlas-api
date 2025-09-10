package postgres

import (
	"context"
	"time"
)

type HypeMapSummaryItem struct {
	EventSlug    string    `json:"event_slug"`
	EventTitle   string    `json:"event_title"`
	Game         string    `json:"game"`
	League       string    `json:"league"`
	Streamers    int       `json:"streamers"`
	TotalViewers int       `json:"total_viewers"`
	LastSeenAt   time.Time `json:"last_seen_at"`
}

// HypeMapSummary agrupa co-streams "activos" por evento.
// Activo = is_live = true O visto en los últimos 5 minutos.
// Filtros opcionales por game/lang.
func (r *Repo) HypeMapSummary(ctx context.Context, game, lang string, limit, offset int) ([]HypeMapSummaryItem, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	q := `
WITH live AS (
  SELECT
    cs.*,
    c.handle,
    c.lang      AS creator_lang,
    e.slug      AS event_slug,
    e.title     AS event_title,
    e.league    AS league,
    e.game      AS game
  FROM app.co_streams cs
  JOIN app.creators c ON c.id = cs.creator_id
  LEFT JOIN app.events   e ON e.id = cs.event_id
  WHERE cs.is_live = TRUE
     OR cs.last_seen_at >= NOW() - INTERVAL '5 minutes'
)
SELECT
  COALESCE(event_slug, 'misc-live')                        AS event_slug,
  COALESCE(event_title, 'Community Live')                  AS event_title,
  COALESCE(game, 'val')                                    AS game,
  COALESCE(league, 'Community')                            AS league,
  COUNT(*)                                                 AS streamers,
  COALESCE(SUM(viewers),0)                                 AS total_viewers,
  MAX(last_seen_at)                                        AS last_seen_at
FROM live
WHERE ($1 = '' OR COALESCE(game,'') = $1)
  AND ($2 = '' OR COALESCE(lang,'') = $2)
GROUP BY 1,2,3,4
ORDER BY total_viewers DESC, last_seen_at DESC
LIMIT $3 OFFSET $4
`
	rows, err := r.db.Query(ctx, q, game, lang, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []HypeMapSummaryItem
	for rows.Next() {
		var it HypeMapSummaryItem
		if err := rows.Scan(&it.EventSlug, &it.EventTitle, &it.Game, &it.League,
			&it.Streamers, &it.TotalViewers, &it.LastSeenAt); err != nil {
			return nil, err
		}
		out = append(out, it)
	}
	return out, rows.Err()
}
