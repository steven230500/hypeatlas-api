package postgres

import (
	"context"
	"time"
)

type StreamRule struct {
	Platform  string
	Handle    string
	EventSlug string
}

type EventWindow struct {
	EventSlug string
	StartsAt  time.Time
	EndsAt    time.Time
	Region    string
	Lang      string
}

// key: "<platform>:<handle>" (handle en minúsculas)
func (r *Repo) LoadStreamRules(ctx context.Context) (map[string]string, error) {
	q := `SELECT platform, lower(handle), event_slug FROM app.event_stream_rules`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make(map[string]string)
	for rows.Next() {
		var p, h, e string
		if err := rows.Scan(&p, &h, &e); err != nil {
			return nil, err
		}
		out[p+":"+h] = e
	}
	return out, rows.Err()
}

func (r *Repo) ActiveWindows(ctx context.Context, now time.Time) ([]EventWindow, error) {
	q := `
SELECT event_slug, starts_at, ends_at, COALESCE(region,''), COALESCE(lang,'')
FROM app.event_windows
WHERE starts_at <= $1 AND ends_at >= $1
`
	rows, err := r.db.Query(ctx, q, now)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []EventWindow
	for rows.Next() {
		var w EventWindow
		if err := rows.Scan(&w.EventSlug, &w.StartsAt, &w.EndsAt, &w.Region, &w.Lang); err != nil {
			return nil, err
		}
		res = append(res, w)
	}
	return res, rows.Err()
}
