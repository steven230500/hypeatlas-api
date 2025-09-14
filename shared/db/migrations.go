package db

import (
	"log"

	"github.com/steven230500/hypeatlas-api/domain/entities"
	"gorm.io/gorm"
)

func ensureSchema(g *gorm.DB) {
	// No falla si ya existe
	_ = g.Exec(`CREATE SCHEMA IF NOT EXISTS app`).Error

	// Limpiar constraints problemáticas de forma más robusta
	constraints := []string{
		"uni_events_slug",
		"idx_events_slug", // nombre alternativo que GORM podría usar
		"events_slug_key", // nombre que PostgreSQL podría asignar
	}

	for _, constraint := range constraints {
		_ = g.Exec(`ALTER TABLE app.events DROP CONSTRAINT IF EXISTS "` + constraint + `"`).Error
		_ = g.Exec(`DROP INDEX IF EXISTS app.` + constraint).Error
	}

	// Manejar columnas created_at/updated_at existentes con datos NULL
	// Primero agregar como nullable, luego actualizar valores, finalmente hacer NOT NULL
	tables := []string{
		"events", "event_windows", "event_stream_rules",
		"creators", "costreams", "patch_changes", "metrics", "games",
		"event_rules", "leagues", "notifications", "sessions",
		"stream_sources", "ingestion_logs", "patches", "hype_thresholds",
		"users", "comps", "league_rankings", "champion_mastery_stats",
		"meta_game_analyses", "champion_rotations",
	}
	for _, table := range tables {
		// Agregar created_at si no existe
		_ = g.Exec(`ALTER TABLE app.` + table + ` ADD COLUMN IF NOT EXISTS created_at timestamptz`).Error
		// Actualizar valores NULL con NOW()
		_ = g.Exec(`UPDATE app.` + table + ` SET created_at = NOW() WHERE created_at IS NULL`).Error
		// Hacer NOT NULL
		_ = g.Exec(`ALTER TABLE app.` + table + ` ALTER COLUMN created_at SET NOT NULL`).Error

		// Lo mismo para updated_at
		_ = g.Exec(`ALTER TABLE app.` + table + ` ADD COLUMN IF NOT EXISTS updated_at timestamptz`).Error
		_ = g.Exec(`UPDATE app.` + table + ` SET updated_at = NOW() WHERE updated_at IS NULL`).Error
		_ = g.Exec(`ALTER TABLE app.` + table + ` ALTER COLUMN updated_at SET NOT NULL`).Error
	}
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
