package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/steven230500/hypeatlas-api/modules/relay/domain/entities"
	out "github.com/steven230500/hypeatlas-api/modules/relay/domain/ports/out"
)

type Repo struct{ db *pgxpool.Pool }

// API (devuelve interfaz)
func New(db *pgxpool.Pool) out.Repository { return &Repo{db: db} }

// Worker (tipo concreto)
func NewRaw(db *pgxpool.Pool) *Repo { return &Repo{db: db} }

// Lee co-streams en vivo por evento y (opcional) idioma.
func (r *Repo) FindLiveByEvent(ctx context.Context, eventID, lang string) ([]entities.CoStream, error) {
	const q = `
SELECT cs.id, cs.event_id, cs.platform, cs.url, cs.lang, cs.country,
       cs.viewers, cs.verified, cs.is_live
FROM app.co_streams cs
WHERE cs.event_id = (
  SELECT id FROM app.events WHERE slug = $1 OR id::text = $1
)
AND cs.is_live = true
AND ($2 = '' OR cs.lang = $2)
ORDER BY cs.viewers DESC NULLS LAST, cs.id;
`
	rows, err := r.db.Query(ctx, q, eventID, lang)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []entities.CoStream
	for rows.Next() {
		var cs entities.CoStream
		if err := rows.Scan(
			&cs.ID, &cs.EventID, &cs.Platform, &cs.URL, &cs.Lang, &cs.Country,
			&cs.Viewers, &cs.Verified, &cs.Live,
		); err != nil {
			return nil, err
		}
		list = append(list, cs)
	}
	return list, rows.Err()
}
