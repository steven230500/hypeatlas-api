package db

import (
	"log"

	"github.com/steven230500/hypeatlas-api/domain/entities"
	"gorm.io/gorm"
)

func ensureSchema(g *gorm.DB) {
	// No falla si ya existe
	_ = g.Exec(`CREATE SCHEMA IF NOT EXISTS app`).Error
	// Si alguna versión antigua quiso borrar una constraint, que no explote:
	_ = g.Exec(`ALTER TABLE app.events DROP CONSTRAINT IF EXISTS "uni_events_slug"`).Error
}

func Migrate(g *gorm.DB) {
	ensureSchema(g)

	if err := g.AutoMigrate(
		// Relay (HypeMap)
		&entities.Event{},
		&entities.Creator{},
		&entities.CoStream{},
		&entities.EventWindow{},
		&entities.EventStreamRule{},
		// Signal (MetaLens)
		&entities.Patch{},
		&entities.PatchChange{},
		&entities.League{},
		&entities.Comp{},
		// Nuevas
		&entities.User{},
		&entities.IngestionLog{},
		&entities.Metric{},
		&entities.Game{},
		&entities.Notification{},
		&entities.Session{},
		&entities.StreamSource{},
		&entities.HypeThreshold{},
		&entities.EventRule{},
		// Meta-game
		&entities.ChampionRotation{},
		&entities.LeagueRanking{},
		&entities.ChampionMasteryStats{},
		&entities.MetaGameAnalysis{},
	); err != nil {
		log.Fatalf("auto-migrate failed: %v", err)
	}

	log.Println("Database migration completed successfully - 22 entities migrated")
}
