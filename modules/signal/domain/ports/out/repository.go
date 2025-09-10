package out

import (
	"context"

	"github.com/steven230500/hypeatlas-api/modules/signal/domain/entities"
)

type Repository interface {
	PatchesByGame(ctx context.Context, game string) ([]entities.Patch, error)
	Changes(ctx context.Context, game, version, entityType string) ([]entities.Change, error)
}
