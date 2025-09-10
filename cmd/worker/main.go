package main

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	relaypg "github.com/steven230500/hypeatlas-api/modules/relay/infra/repository/postgres"
	signalpg "github.com/steven230500/hypeatlas-api/modules/signal/infra/repository/postgres"
	pg "github.com/steven230500/hypeatlas-api/shared/db"

	twitchprov "github.com/steven230500/hypeatlas-api/providers/twitch"
)

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

	relayRepo := relaypg.NewRaw(pool)
	signalRepo := signalpg.NewRaw(pool)

	// Twitch config
	twID := os.Getenv("TWITCH_CLIENT_ID")
	twSec := os.Getenv("TWITCH_SECRET")
	twHandles := strings.Split(strings.ToLower(os.Getenv("TWITCH_HANDLES")), ",")
	var tw *twitchprov.Client
	if twID != "" && twSec != "" {
		tw = twitchprov.New(twID, twSec)
	}

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

	runOnce(ctx, relayRepo, signalRepo, tw, twHandles)
	for range ticker.C {
		runOnce(ctx, relayRepo, signalRepo, tw, twHandles)
	}
}

func runOnce(
	ctx context.Context,
	relayRepo *relaypg.Repo,
	signalRepo *signalpg.Repo,
	tw *twitchprov.Client,
	_ []string, // ya no usamos TWITCH_HANDLES del .env
) {
	// ====== Resolver básico de evento (si ya tienes LoadStreamRules/ActiveWindows) ======
	rules, _ := relayRepo.LoadStreamRules(ctx) // map["twitch:<login>"]=event_slug
	wins, _ := relayRepo.ActiveWindows(ctx, time.Now())

	resolveEvent := func(platform, login, lang string) (slug, title, game, league string) {
		key := platform + ":" + strings.ToLower(login)
		if e, ok := rules[key]; ok {
			return e, "Mapped by rule", "val", "Community"
		}
		for _, w := range wins {
			if w.Lang == "" || strings.EqualFold(w.Lang, lang) {
				return w.EventSlug, "Window match", "val", "Community"
			}
		}
		return "misc-live", "Community Live", "val", "Community"
	}

	// ====== TWITCH: lee handles desde DB y upserts ======
	if tw != nil {
		creators, err := relayRepo.ListCreatorHandles(ctx, "twitch", true)
		if err != nil {
			log.Error().Err(err).Msg("list twitch creators failed")
		} else if len(creators) > 0 {
			var logins []string
			for _, c := range creators {
				logins = append(logins, c.Handle)
			}
			for _, chunk := range twitchprov.Chunk(logins, 100) {
				streams, err := tw.GetStreamsByLogin(ctx, chunk)
				if err != nil {
					log.Error().Err(err).Msg("twitch GetStreams failed")
					continue
				}
				for login, s := range streams {
					url := "https://twitch.tv/" + login
					eventSlug, eventTitle, game, league := resolveEvent("twitch", login, s.Language)

					if err := relayRepo.UpsertCoStream(ctx,
						eventSlug, eventTitle, game, league, nil,
						"twitch", login, url, s.Language, "", true,
						s.ViewerCount, s.Type == "live",
					); err != nil {
						log.Error().Err(err).Str("login", login).Msg("upsert co-stream failed")
					}
				}
			}
		}
	}

	// ====== META/COMPS (mock para demo) ======
	for _, c := range mockPullComps() {
		raw, _ := json.Marshal(c.Slots)
		if err := signalRepo.UpsertComp(ctx,
			c.Game, c.Region, c.League, c.Patch, c.Map, c.Side, string(raw), c.Pick, c.Win, c.Delta,
		); err != nil {
			log.Error().Err(err).Str("patch", c.Patch).Msg("upsert comp failed")
		}
	}

	// ====== CLEANUP: marca offline streams con más de N minutos sin señal ======
	staleMinutes := 10
	if v := os.Getenv("WORKER_STALE_MINUTES"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			staleMinutes = n
		}
	}
	affected, err := relayRepo.MarkStaleCoStreamsOffline(ctx, time.Duration(staleMinutes)*time.Minute)
	if err != nil {
		log.Error().Err(err).Int("stale_minutes", staleMinutes).Msg("cleanup stale co_streams failed")
	} else if affected > 0 {
		log.Info().Int64("rows", affected).Int("stale_minutes", staleMinutes).Msg("cleanup co_streams marked offline")
	}

	log.Info().Msg("ingest cycle OK")
}

func mockPullComps() []CompSample {
	pick := 25.1
	win := 53.2
	delta := 1.8
	return []CompSample{{
		Game: "val", Region: "EMEA", League: "VCT EMEA", Patch: "9.15", Map: "Ascent", Side: "attack",
		Slots: map[string]any{
			"roles":   []string{"smokes", "initiator", "duelist", "sentinel", "flex"},
			"members": []map[string]string{{"agent": "omen"}, {"agent": "sova"}, {"agent": "jett"}, {"agent": "killjoy"}, {"agent": "skye"}},
		},
		Pick: &pick, Win: &win, Delta: &delta,
	}}
}
