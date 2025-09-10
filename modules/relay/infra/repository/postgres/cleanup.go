package postgres

import (
	"context"
	"math"
	"time"
)

// MarkStaleCoStreamsOffline pone is_live = false a co_streams cuya última señal
// sea anterior a (NOW() - olderThan). Retorna filas afectadas.
func (r *Repo) MarkStaleCoStreamsOffline(ctx context.Context, olderThan time.Duration) (int64, error) {
	if olderThan <= 0 {
		olderThan = 10 * time.Minute
	}
	mins := int(math.Round(olderThan.Minutes()))
	q := `
UPDATE app.co_streams
SET is_live = FALSE
WHERE is_live = TRUE
  AND last_seen_at < NOW() - ($1::text || ' minutes')::interval
`
	cmd, err := r.db.Exec(ctx, q, mins)
	if err != nil {
		return 0, err
	}
	return cmd.RowsAffected(), nil
}
