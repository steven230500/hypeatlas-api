package out

import (
	"context"

	"github.com/steven230500/hypeatlas-api/modules/relay/domain/entities"
)

type Repository interface {
	FindLiveByEvent(ctx context.Context, eventID, lang string) ([]entities.CoStream, error)
}
