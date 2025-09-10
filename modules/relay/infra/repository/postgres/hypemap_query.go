package postgres

import (
	"context"
)

type HypeMapItem struct {
	EventSlug  string  `json:"event_slug"`
	EventTitle string  `json:"event_title"`
	Game       string  `json:"game"`
	League     string  `json:"league"`
	Platform   string  `json:"platform"`
	Handle     string  `json:"handle"`
	Lang       string  `json:"lang"`
	Country    string  `json:"country"`
	Viewers    int     `json:"viewers"`
	IsLive     bool    `json:"is_live"`
	Score      float64 `json:"score"`
}

func (r *Repo) HypeMapLive(ctx context.Context, game, lang string, limit, offset int) ([]HypeMapItem, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	// Simple scoring: viewers + bonus verified + bonus idioma
	q := `
WITH live AS (
  SELECT cs.*, c.handle, c.lang AS creator_lang, c.country,
         e.title AS event_title, e.league, e.game,
         COALESCE(e.slug, 'misc-live') AS event_slug
  FROM app.co_streams cs
  JOIN app.creators c ON c.id = cs.creator_id
  LEFT JOIN app.events e ON e.id = cs.event_id
  WHERE cs.is_live = TRUE OR cs.last_seen_at >= NOW() - INTERVAL '5 minutes'
)
SELECT
  event_slug,
  COALESCE(event_title, 'Community Live') AS event_title,
  COALESCE(game, 'val') AS game,
  COALESCE(league, 'Community') AS league,
  platform,
  handle,
  COALESCE(lang, '') AS lang,
  COALESCE(country, '') AS country,
  viewers,
  is_live,
  (viewers
   + CASE WHEN verified THEN 300 ELSE 0 END
   + CASE WHEN lang = $1 THEN 100 ELSE 0 END
  )::float AS score
FROM live
WHERE ($2 = '' OR game = $2)
  AND ($1 = '' OR lang = $1)
ORDER BY score DESC
LIMIT $3 OFFSET $4
`
	rows, err := r.db.Query(ctx, q, lang, game, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []HypeMapItem
	for rows.Next() {
		var it HypeMapItem
		if err := rows.Scan(&it.EventSlug, &it.EventTitle, &it.Game, &it.League,
			&it.Platform, &it.Handle, &it.Lang, &it.Country, &it.Viewers, &it.IsLive, &it.Score); err != nil {
			return nil, err
		}
		out = append(out, it)
	}
	return out, rows.Err()
}
