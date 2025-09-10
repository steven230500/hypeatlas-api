package main

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"

	relaypg "github.com/steven230500/hypeatlas-api/modules/relay/infra/repository/postgres"
	signalpg "github.com/steven230500/hypeatlas-api/modules/signal/infra/repository/postgres"
	pg "github.com/steven230500/hypeatlas-api/shared/db"
)

type CoStreamSample struct {
	EventSlug   string
	EventTitle  string
	Game        string // "val"|"lol"
	League      string
	Platform    string // "twitch"|"youtube"
	Handle      string
	URL         string
	Lang        string
	Country     string
	Verified    bool
	Viewers     int
	IsLive      bool
	StartsAtISO *string
}

type CompSample struct {
	Game, Region, League, Patch, Map, Side string
	Slots                                  map[string]any
	Pick, Win, Delta                       *float64
}

func main() {
	if os.Getenv("POSTGRES_URL") == "" {
		log.Fatal().Msg("POSTGRES_URL missing")
	}
	pool := pg.MustOpen()
	defer pool.Close()

	// Ojo: usamos NewRaw para obtener *Repo (tipo concreto)
	relayRepo := relaypg.NewRaw(pool)
	signalRepo := signalpg.NewRaw(pool)

	interval := 30 * time.Second
	if v := os.Getenv("WORKER_INTERVAL_SEC"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			interval = time.Duration(n) * time.Second
		}
	}

	ctx := context.Background()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Info().Dur("interval", interval).Msg("worker started")

	// primera corrida inmediata y luego cada tick
	runOnce(ctx, relayRepo, signalRepo)
	for range ticker.C {
		runOnce(ctx, relayRepo, signalRepo)
	}
}

func runOnce(ctx context.Context, relayRepo *relaypg.Repo, signalRepo *signalpg.Repo) {
	// 1) Ingesta de co-streams (mock)
	streams := mockPullCoStreams()
	for _, s := range streams {
		if err := relayRepo.UpsertCoStream(ctx,
			s.EventSlug, s.EventTitle, s.Game, s.League, s.StartsAtISO,
			s.Platform, s.Handle, s.URL, s.Lang, s.Country, s.Verified, s.Viewers, s.IsLive,
		); err != nil {
			log.Error().Err(err).Str("event", s.EventSlug).Str("handle", s.Handle).Msg("upsert co-stream failed")
		}
	}

	// 2) Ingesta de comps/meta (mock)
	comps := mockPullComps()
	for _, c := range comps {
		raw, _ := json.Marshal(c.Slots)
		if err := signalRepo.UpsertComp(ctx,
			c.Game, c.Region, c.League, c.Patch, c.Map, c.Side,
			string(raw), c.Pick, c.Win, c.Delta,
		); err != nil {
			log.Error().Err(err).Str("patch", c.Patch).Msg("upsert comp failed")
		}
	}

	log.Info().Int("streams", len(streams)).Int("comps", len(comps)).Msg("ingest cycle OK")
}

// ========== MOCKS ==========

func mockPullCoStreams() []CoStreamSample {
	return []CoStreamSample{
		{
			EventSlug:  "vct-emea-final",
			EventTitle: "VCT EMEA Final",
			Game:       "val",
			League:     "VCT EMEA",
			Platform:   "twitch",
			Handle:     "caster1",
			URL:        "https://twitch.tv/caster1",
			Lang:       "es",
			Country:    "ES",
			Verified:   true,
			Viewers:    9100,
			IsLive:     true,
		},
	}
}

func mockPullComps() []CompSample {
	pick := 25.1
	win := 53.2
	delta := 1.8
	return []CompSample{
		{
			Game:   "val",
			Region: "EMEA",
			League: "VCT EMEA",
			Patch:  "9.15",
			Map:    "Ascent",
			Side:   "attack",
			Slots: map[string]any{
				"roles": []string{"smokes", "initiator", "duelist", "sentinel", "flex"},
				"members": []map[string]string{
					{"agent": "omen"}, {"agent": "sova"}, {"agent": "jett"},
					{"agent": "killjoy"}, {"agent": "skye"},
				},
			},
			Pick:  &pick,
			Win:   &win,
			Delta: &delta,
		},
	}
}
