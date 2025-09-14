package db

import (
	"log"

	"gorm.io/gorm"

	"github.com/steven230500/hypeatlas-api/domain/entities"
)

func Migrate(g *gorm.DB) {
	err := g.AutoMigrate(
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

		// New entities
		&entities.User{},
		&entities.IngestionLog{},
		&entities.Metric{},
		&entities.Game{},
		&entities.Notification{},
		&entities.Session{},
		&entities.StreamSource{},
		&entities.HypeThreshold{},
		&entities.EventRule{},

		// Meta-game analysis entities
		&entities.ChampionRotation{},
		&entities.LeagueRanking{},
		&entities.ChampionMasteryStats{},
		&entities.MetaGameAnalysis{},
	)

	if err != nil {
		log.Fatalf("auto-migrate failed: %v", err)
	}

	log.Println("Database migration completed successfully - 22 entities migrated")
}
