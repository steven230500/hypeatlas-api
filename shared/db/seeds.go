package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/steven230500/hypeatlas-api/domain/entities"
)

// RunSeeds replica tus seeds SQL en GORM.
func RunSeeds(g *gorm.DB) error {
	// --- Event base
	var ev entities.Event
	err := g.Where("slug = ?", "vct-emea-final").First(&ev).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if ev.UUID == uuid.Nil {
		now := time.Now().UTC().Add(24 * time.Hour)
		ev = entities.Event{
			Slug:     "vct-emea-final",
			Title:    "VCT EMEA Final",
			Game:     "val",
			League:   ptr("VCT EMEA"),
			StartsAt: &now,
		}
		if err := g.Create(&ev).Error; err != nil && !IsDuplicateEntry(err) {
			return err
		}
	}

	// --- Creators KOI y otros
	creators := []entities.Creator{
		{Platform: "twitch", Handle: "ibai", URL: "https://twitch.tv/ibai", Lang: "es", Country: "ES", Verified: true},
		{Platform: "twitch", Handle: "ernesbarbeq", URL: "https://twitch.tv/ernesbarbeq", Lang: "es", Country: "ES", Verified: true},
		{Platform: "twitch", Handle: "tvander", URL: "https://twitch.tv/tvander", Lang: "es", Country: "ES", Verified: true},
		{Platform: "twitch", Handle: "knekro", URL: "https://twitch.tv/knekro", Lang: "es", Country: "ES", Verified: true},
		{Platform: "twitch", Handle: "rioboo", URL: "https://twitch.tv/rioboo", Lang: "es", Country: "ES", Verified: true},
		{Platform: "twitch", Handle: "lakshartnia", URL: "https://twitch.tv/lakshartnia", Lang: "es", Country: "ES", Verified: true},
		{Platform: "twitch", Handle: "mayichi", URL: "https://twitch.tv/mayichi", Lang: "es", Country: "ES", Verified: true},
		{Platform: "twitch", Handle: "karchez", URL: "https://twitch.tv/karchez", Lang: "es", Country: "ES", Verified: true},
		{Platform: "twitch", Handle: "blackelespanolito", URL: "https://twitch.tv/blackelespanolito", Lang: "es", Country: "ES", Verified: true},
		{Platform: "twitch", Handle: "skain", URL: "https://twitch.tv/skain", Lang: "es", Country: "ES", Verified: true},
		{Platform: "twitch", Handle: "sergiofferra", URL: "https://twitch.tv/sergiofferra", Lang: "es", Country: "ES", Verified: true},
		{Platform: "twitch", Handle: "koi", URL: "https://twitch.tv/koi", Lang: "es", Country: "ES", Verified: true},
	}
	for _, c := range creators {
		_ = g.Where("platform = ? AND handle = ?", c.Platform, c.Handle).
			Attrs(c).FirstOrCreate(&entities.Creator{}).Error
	}

	// --- CoStream sample
	var koi entities.Creator
	_ = g.Where("platform = ? AND handle = ?", "twitch", "koi").First(&koi).Error
	if koi.UUID != uuid.Nil && ev.UUID != uuid.Nil {
		var exists int64
		_ = g.Model(&entities.CoStream{}).
			Where("event_uuid = ? AND creator_uuid = ?", ev.UUID, koi.UUID).
			Count(&exists).Error
		if exists == 0 {
			_ = g.Create(&entities.CoStream{
				EventUUID:   ev.UUID,
				CreatorUUID: koi.UUID,
				Platform:    "twitch",
				URL:         "https://twitch.tv/koi",
				Lang:        "es",
				Country:     "ES",
				Viewers:     8200,
				Verified:    true,
				IsLive:      true,
				LastSeenAt:  time.Now().UTC(),
			}).Error
		}
	}

	// --- Event windows & rules
	_ = g.Where("platform=? AND handle=?", "twitch", "koi").
		Attrs(entities.EventStreamRule{EventSlug: "vct-emea-final"}).
		FirstOrCreate(&entities.EventStreamRule{}).Error
	_ = g.Where("platform=? AND handle=?", "twitch", "sergiofferra").
		Attrs(entities.EventStreamRule{EventSlug: "vct-emea-final"}).
		FirstOrCreate(&entities.EventStreamRule{}).Error

	_ = g.Create(&entities.EventWindow{
		EventSlug: "vct-emea-final",
		StartsAt:  time.Now().UTC().Add(-1 * time.Hour),
		EndsAt:    time.Now().UTC().Add(6 * time.Hour),
		Region:    "EMEA",
		Lang:      "es",
	}).Error

	// --- Patches + change
	var p entities.Patch
	_ = g.Where("game=? AND version=?", "val", "9.15").
		Attrs(entities.Patch{ReleasedAt: ptrTime(time.Now().UTC().Add(-8 * 24 * time.Hour))}).
		FirstOrCreate(&p).Error

	_ = g.Where("patch_uuid=? AND entity_type=? AND entity_id=? AND field=?",
		p.UUID, "agent", "sova", "recon bolt cd").
		Attrs(entities.PatchChange{
			Old:         "40s",
			New:         "45s",
			ImpactScore: 0.6,
		}).
		FirstOrCreate(&entities.PatchChange{}).Error

	// --- Leagues
	_ = g.Where("slug=?", "vct-emea").
		Attrs(entities.League{Game: "val", Region: "EMEA", Name: "VCT EMEA"}).
		FirstOrCreate(&entities.League{}).Error
	_ = g.Where("slug=?", "lec").
		Attrs(entities.League{Game: "lol", Region: "EMEA", Name: "LEC"}).
		FirstOrCreate(&entities.League{}).Error

	// --- Comp ejemplo (VAL / EMEA / Ascent)
	slots := datatypes.JSON([]byte(`{
		"roles": ["smokes","initiator","duelist","sentinel","flex"],
		"members":[{"agent":"omen"},{"agent":"sova"},{"agent":"jett"},{"agent":"killjoy"},{"agent":"skye"}]
	}`))
	_ = g.Where("game=? AND region=? AND league=? AND patch=? AND map=? AND side=? AND slots_fp=?",
		"val", "EMEA", "VCT EMEA", "9.15", "Ascent", "attack", "").
		Attrs(entities.Comp{
			Game: "val", Region: "EMEA", League: "VCT EMEA", Patch: "9.15", Map: "Ascent", Side: "attack",
			Slots: slots, PickRate: ptrf(24.300), WinRate: ptrf(52.100), DeltaWin: ptrf(1.600),
		}).
		FirstOrCreate(&entities.Comp{}).Error // BeforeSave completará slots_fp

	// --- New entities seeds

	// Games
	_ = g.Where("slug=?", "val").
		Attrs(entities.Game{Name: "Valorant", Slug: "val", Platforms: `["twitch","youtube"]`}).
		FirstOrCreate(&entities.Game{}).Error
	_ = g.Where("slug=?", "lol").
		Attrs(entities.Game{Name: "League of Legends", Slug: "lol", Platforms: `["twitch","youtube"]`}).
		FirstOrCreate(&entities.Game{}).Error

	// Stream sources
	_ = g.Where("name=?", "twitch").
		Attrs(entities.StreamSource{
			Name: "Twitch", BaseURL: "https://api.twitch.tv/helix", ApiKey: "your-twitch-api-key", IsActive: true,
		}).
		FirstOrCreate(&entities.StreamSource{}).Error

	// Users
	_ = g.Where("email=?", "admin@hypeatlas.com").
		Attrs(entities.User{Email: "admin@hypeatlas.com", Role: "admin", Verified: true}).
		FirstOrCreate(&entities.User{}).Error

	// Hype thresholds
	if ev.UUID != uuid.Nil {
		_ = g.Where("event_id=? AND game=?", ev.UUID, "val").
			Attrs(entities.HypeThreshold{
				EventID: ev.UUID, Game: "val", MinViewers: 5000, MaxViewers: 15000, AlertLevel: "high", IsActive: true,
			}).
			FirstOrCreate(&entities.HypeThreshold{}).Error
	}

	// Event rules
	if ev.UUID != uuid.Nil {
		_ = g.Where("event_id=? AND platform=? AND handle=?", ev.UUID, "twitch", "koi").
			Attrs(entities.EventRule{
				EventID: ev.UUID, Platform: "twitch", Handle: "koi", AutoAssign: true, Priority: 1,
			}).
			FirstOrCreate(&entities.EventRule{}).Error
	}

	return nil
}

func ptr(s string) *string           { return &s }
func ptrTime(t time.Time) *time.Time { return &t }
func ptrf(f float64) *float64        { return &f }
