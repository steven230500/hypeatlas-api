package memory

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/steven230500/hypeatlas-api/modules/signal/domain/entities"
	out "github.com/steven230500/hypeatlas-api/modules/signal/domain/ports/out"
)

type memRepo struct {
	patches []entities.Patch
	changes []entities.Change
	leagues []entities.League
	comps   []entities.Comp
}

// New devuelve un repositorio en memoria que satisface out.Repository.
func New() out.Repository {
	// Datos de ejemplo/seed
	patchVAL := entities.Patch{ID: "p1", Game: "val", Version: "9.15", ReleasedAt: "2025-09-01"}
	patchLOL := entities.Patch{ID: "p2", Game: "lol", Version: "14.14", ReleasedAt: "2025-08-27"}

	// Comp de ejemplo (VAL / EMEA / Ascent / attack)
	slots := map[string]any{
		"roles": []string{"smokes", "initiator", "duelist", "sentinel", "flex"},
		"members": []map[string]string{
			{"agent": "omen"}, {"agent": "sova"}, {"agent": "jett"},
			{"agent": "killjoy"}, {"agent": "skye"},
		},
	}
	rawSlots, _ := json.Marshal(slots)
	pick := 24.300
	win := 52.100
	delta := 1.600

	return &memRepo{
		patches: []entities.Patch{patchVAL, patchLOL},
		changes: []entities.Change{
			{PatchID: patchVAL.ID, EntityType: "agent", EntityID: "sova", Field: "recon bolt cd", Old: "40s", New: "45s", Impact: 0.6},
			{PatchID: patchLOL.ID, EntityType: "champion", EntityID: "azir", Field: "W mana", Old: "40", New: "30", Impact: 0.4},
		},
		leagues: []entities.League{
			{ID: "l1", Game: "val", Region: "EMEA", Name: "VCT EMEA", Slug: "vct-emea"},
			{ID: "l2", Game: "lol", Region: "EMEA", Name: "LEC", Slug: "lec"},
		},
		comps: []entities.Comp{
			{
				ID:        "c1",
				Game:      "val",
				Region:    "EMEA",
				League:    "VCT EMEA",
				Patch:     "9.15",
				Map:       "Ascent",
				Side:      "attack",
				SlotsJSON: string(rawSlots),
				PickRate:  &pick,
				WinRate:   &win,
				DeltaWin:  &delta,
			},
		},
	}
}

func (m *memRepo) PatchesByGame(_ context.Context, game string) ([]entities.Patch, error) {
	var outv []entities.Patch
	for _, p := range m.patches {
		if strings.EqualFold(p.Game, game) {
			outv = append(outv, p)
		}
	}
	return outv, nil
}

func (m *memRepo) Changes(_ context.Context, game, version, entityType string) ([]entities.Change, error) {
	// Filtra por game/version usando la tabla de patches (join en memoria)
	patchByID := make(map[string]entities.Patch, len(m.patches))
	for _, p := range m.patches {
		patchByID[p.ID] = p
	}
	var outv []entities.Change
	for _, c := range m.changes {
		p, ok := patchByID[c.PatchID]
		if !ok {
			continue
		}
		if !strings.EqualFold(p.Game, game) {
			continue
		}
		if !strings.EqualFold(p.Version, version) {
			continue
		}
		if entityType != "" && !strings.EqualFold(c.EntityType, entityType) {
			continue
		}
		outv = append(outv, c)
	}
	return outv, nil
}

func (m *memRepo) Leagues(_ context.Context, game, region string) ([]entities.League, error) {
	var outv []entities.League
	for _, l := range m.leagues {
		if !strings.EqualFold(l.Game, game) {
			continue
		}
		if region != "" && !strings.EqualFold(l.Region, region) {
			continue
		}
		outv = append(outv, l)
	}
	return outv, nil
}

func (m *memRepo) Comps(_ context.Context, game, region, league, patch, mapp, side string, limit int) ([]entities.Comp, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	var outv []entities.Comp
	for _, c := range m.comps {
		if !strings.EqualFold(c.Game, game) {
			continue
		}
		if !strings.EqualFold(c.Region, region) {
			continue
		}
		if !strings.EqualFold(c.Patch, patch) {
			continue
		}
		if league != "" && !strings.EqualFold(c.League, league) {
			continue
		}
		if mapp != "" && !strings.EqualFold(c.Map, mapp) {
			continue
		}
		if side != "" && !strings.EqualFold(c.Side, side) {
			continue
		}
		outv = append(outv, c)
		if len(outv) >= limit {
			break
		}
	}
	return outv, nil
}
