package postgres

import (
	"context"
)

// UpsertCoStream inserta/actualiza Event, Creator y CoStream atómicamente.
// eventSlug/title/game/league describen el evento (slug es único).
// creator: platform(twitch|youtube) + handle es único.
// lang: "es", country: "ES", viewers, verified, isLive son métricas actuales.
// startsAt puede venir vacío (usar NULL para now()).
func (r *Repo) UpsertCoStream(
	ctx context.Context,
	eventSlug, eventTitle, game, league string,
	startsAtNullable *string, // "YYYY-MM-DDTHH:MM:SSZ" o nil
	platform, handle, url, lang, country string,
	verified bool,
	viewers int,
	isLive bool,
) error {
	// language=SQL
	const q = `
WITH e AS (
  INSERT INTO app.events (slug, title, game, league, starts_at)
  VALUES ($1, $2, $3, $4, COALESCE($5::timestamptz, now()))
  ON CONFLICT (slug) DO UPDATE
    SET title = EXCLUDED.title,
        league = COALESCE(EXCLUDED.league, app.events.league)
  RETURNING id
),
c AS (
  INSERT INTO app.creators (platform, handle, url, lang, country, verified)
  VALUES ($6, $7, $8, $9, $10, $11)
  ON CONFLICT (platform, handle) DO UPDATE
    SET url = EXCLUDED.url,
        lang = EXCLUDED.lang,
        country = EXCLUDED.country,
        verified = EXCLUDED.verified
  RETURNING id
)
INSERT INTO app.co_streams
  (event_id, creator_id, platform, url, lang, country, viewers, verified, is_live, last_seen_at)
SELECT e.id, c.id, $6, $8, $9, $10, $12, $11, $13, now()
FROM e, c
ON CONFLICT (event_id, creator_id) DO UPDATE SET
  platform     = EXCLUDED.platform,
  url          = EXCLUDED.url,
  lang         = EXCLUDED.lang,
  country      = EXCLUDED.country,
  viewers      = EXCLUDED.viewers,
  verified     = EXCLUDED.verified,
  is_live      = EXCLUDED.is_live,
  last_seen_at = now();
`
	_, err := r.db.Exec(ctx, q,
		eventSlug, eventTitle, game, league, startsAtNullable,
		platform, handle, url, lang, country, verified,
		viewers, isLive,
	)
	return err
}
