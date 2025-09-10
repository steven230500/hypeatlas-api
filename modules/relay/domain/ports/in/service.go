package in

import (
	"context"

	"github.com/steven230500/hypeatlas-api/modules/relay/domain/entities"
)

type Service interface {
	ListLiveCoStreams(ctx context.Context, eventID, lang string) ([]entities.CoStream, error)
}
