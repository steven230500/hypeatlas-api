package memory

import (
	"context"

	"github.com/steven230500/hypeatlas-api/modules/signal/domain/entities"
	out "github.com/steven230500/hypeatlas-api/modules/signal/domain/ports/out"
)

type memRepo struct {
	patches []entities.Patch
	changes []entities.Change
}

func New() out.Repository {
	return &memRepo{
		patches: []entities.Patch{
			{ID: "p1", Game: "val", Version: "9.15", ReleasedAt: "2025-09-01"},
			{ID: "p2", Game: "lol", Version: "14.14", ReleasedAt: "2025-08-27"},
		},
		changes: []entities.Change{
			{PatchID: "p1", EntityType: "agent", EntityID: "sova", Field: "recon bolt cd", Old: "40s", New: "45s", Impact: 0.6},
			{PatchID: "p2", EntityType: "champion", EntityID: "azir", Field: "W mana", Old: "40", New: "30", Impact: 0.4},
		},
	}
}

func (m *memRepo) PatchesByGame(_ context.Context, game string) ([]entities.Patch, error) {
	var outv []entities.Patch
	for _, p := range m.patches {
		if p.Game == game {
			outv = append(outv, p)
		}
	}
	return outv, nil
}

func (m *memRepo) Changes(_ context.Context, game, version, entityType string) ([]entities.Change, error) {
	// Mock simple: devuelve todos; luego filtramos por game/version/type
	return m.changes, nil
}
