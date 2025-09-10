package in

import (
	"context"

	"github.com/steven230500/hypeatlas-api/modules/signal/domain/entities"
)

type Service interface {
	ListPatches(ctx context.Context, game string) ([]entities.Patch, error)
	ListChanges(ctx context.Context, game, version, entityType string) ([]entities.Change, error)
}
