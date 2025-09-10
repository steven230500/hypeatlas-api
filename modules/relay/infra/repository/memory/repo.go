package memory

import (
	"context"
	"strings"

	"github.com/steven230500/hypeatlas-api/modules/relay/domain/entities"
	out "github.com/steven230500/hypeatlas-api/modules/relay/domain/ports/out"
)

type memRepo struct{ data []entities.CoStream }

func New() out.Repository {
	return &memRepo{
		data: []entities.CoStream{
			{ID: "c1", EventID: "vct-emea-final", Platform: "twitch", URL: "https://twitch.tv/caster1", Lang: "es", Country: "ES", Viewers: 8200, Verified: true, Live: true},
			{ID: "c2", EventID: "vct-emea-final", Platform: "youtube", URL: "https://youtube.com/@caster2/live", Lang: "en", Country: "GB", Viewers: 4300, Verified: false, Live: true},
		},
	}
}

func (m *memRepo) FindLiveByEvent(_ context.Context, eventID, lang string) ([]entities.CoStream, error) {
	var outv []entities.CoStream
	for _, cs := range m.data {
		if !cs.Live || cs.EventID != eventID {
			continue
		}
		if lang != "" && !strings.EqualFold(lang, cs.Lang) {
			continue
		}
		outv = append(outv, cs)
	}
	return outv, nil
}
