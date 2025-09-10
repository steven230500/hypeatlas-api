package postgres

import (
	"context"
)

type CreatorLite struct {
	Handle   string
	Lang     string
	Country  string
	Verified bool
}

// ListCreatorHandles devuelve handles (lowercase) por plataforma.
// Si verifiedOnly=true, filtra por verified=true.
func (r *Repo) ListCreatorHandles(ctx context.Context, platform string, verifiedOnly bool) ([]CreatorLite, error) {
	q := `
SELECT lower(handle) AS handle, COALESCE(lang,''), COALESCE(country,''), verified
FROM app.creators
WHERE platform = $1
`
	args := []any{platform}
	if verifiedOnly {
		q += " AND verified = true"
	}
	q += " ORDER BY handle"
	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []CreatorLite
	for rows.Next() {
		var c CreatorLite
		if err := rows.Scan(&c.Handle, &c.Lang, &c.Country, &c.Verified); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}
