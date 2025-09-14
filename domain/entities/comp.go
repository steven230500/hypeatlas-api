package entities

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	"gorm.io/datatypes"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Comp struct {
	UUID   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"uuid"`
	Game   string    `gorm:"type:varchar(10);not null;index:idx_comps_full_filter,priority:1;index:idx_comps_filter,priority:1" json:"game"`
	Region string    `gorm:"type:varchar(16);not null;index:idx_comps_full_filter,priority:2;index:idx_comps_filter,priority:2" json:"region"`
	League string    `gorm:"type:varchar(120);not null;default:'';index:idx_comps_full_filter,priority:4;index:idx_comps_filter,priority:3" json:"league"`
	Patch  string    `gorm:"type:varchar(32);not null;index:idx_comps_full_filter,priority:3;index:idx_comps_filter,priority:4" json:"patch"`
	Map    string    `gorm:"type:varchar(64);not null;default:'';index:idx_comps_full_filter,priority:5" json:"map"`
	Side   string    `gorm:"type:varchar(16);not null;default:'';index:idx_comps_full_filter,priority:6" json:"side"`

	Slots   datatypes.JSON `gorm:"type:jsonb;not null" json:"slots"`
	SlotsFP string         `gorm:"type:text;not null;uniqueIndex:uq_comp_fingerprint" json:"-"`

	// UNIQUE lógico: (game,region,league,patch,map,side,slots_fp)
	_ struct{} `gorm:"uniqueIndex:uq_comp_fingerprint,priority:1"` // game
	_ struct{} `gorm:"uniqueIndex:uq_comp_fingerprint,priority:2"` // region
	_ struct{} `gorm:"uniqueIndex:uq_comp_fingerprint,priority:3"` // league
	_ struct{} `gorm:"uniqueIndex:uq_comp_fingerprint,priority:4"` // patch
	_ struct{} `gorm:"uniqueIndex:uq_comp_fingerprint,priority:5"` // map
	_ struct{} `gorm:"uniqueIndex:uq_comp_fingerprint,priority:6"` // side
	// slots_fp es priority:7 por tag en el propio campo

	PickRate  *float64  `gorm:"type:numeric(6,3)" json:"pick_rate,omitempty"`
	WinRate   *float64  `gorm:"type:numeric(6,3)" json:"win_rate,omitempty"`
	DeltaWin  *float64  `gorm:"type:numeric(6,3)" json:"delta_win,omitempty"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;index:idx_comps_created,sort:desc" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null"                                        json:"updated_at"`
}

func (Comp) TableName() string { return "app.comps" }

// Hook: simula columna generada slots_fp = md5(slots::text)
func (c *Comp) BeforeSave(tx *gorm.DB) error {
	if len(c.Slots) > 0 {
		sum := md5.Sum([]byte(c.Slots.String()))
		c.SlotsFP = hex.EncodeToString(sum[:])
	}
	return nil
}

// ChampionRotation representa la rotación semanal de campeones gratuitos
type ChampionRotation struct {
	UUID                  uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"uuid"`
	Game                  string         `gorm:"type:varchar(10);not null;index"                json:"game"`     // lol|val
	Platform              string         `gorm:"type:varchar(16);not null;index"                json:"platform"` // na1, euw1, etc.
	FreeChampionIDs       datatypes.JSON `gorm:"type:jsonb;not null" json:"free_champion_ids"`
	FreeChampionIDsForNew datatypes.JSON `gorm:"type:jsonb;not null" json:"free_champion_ids_for_new"`
	MaxNewPlayerLevel     int            `gorm:"not null" json:"max_new_player_level"`
	RotationStartedAt     *time.Time     `gorm:"type:timestamptz" json:"rotation_started_at"`
	RotationEndsAt        *time.Time     `gorm:"type:timestamptz" json:"rotation_ends_at"`
	CreatedAt             time.Time      `gorm:"type:timestamptz;not null" json:"created_at"`
	UpdatedAt             time.Time      `gorm:"type:timestamptz;not null" json:"updated_at"`
}

func (ChampionRotation) TableName() string { return "app.champion_rotations" }

// LeagueRanking representa estadísticas de una liga específica
type LeagueRanking struct {
	UUID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"uuid"`
	Game         string    `gorm:"type:varchar(10);not null;index"                json:"game"`
	Platform     string    `gorm:"type:varchar(16);not null;index"                json:"platform"`
	Queue        string    `gorm:"type:varchar(32);not null;index"                json:"queue"`    // RANKED_SOLO_5x5, etc.
	Tier         string    `gorm:"type:varchar(16);not null;index"                json:"tier"`     // CHALLENGER, MASTER, etc.
	Division     string    `gorm:"type:varchar(8);default:''"                     json:"division"` // I, II, III, IV
	LeagueID     string    `gorm:"type:varchar(64);not null"                      json:"league_id"`
	LeagueName   string    `gorm:"type:varchar(120);not null"                     json:"league_name"`
	TotalPlayers int       `gorm:"not null" json:"total_players"`
	AvgLP        float64   `gorm:"type:numeric(8,2)" json:"avg_lp"`
	AvgWinRate   float64   `gorm:"type:numeric(5,2)" json:"avg_win_rate"`
	TopLP        int       `gorm:"not null" json:"top_lp"`
	CreatedAt    time.Time `gorm:"type:timestamptz;not null" json:"created_at"`
	UpdatedAt    time.Time `gorm:"type:timestamptz;not null" json:"updated_at"`
}

func (LeagueRanking) TableName() string { return "app.league_rankings" }

// ChampionMasteryStats representa estadísticas de maestría de un campeón
type ChampionMasteryStats struct {
	UUID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"uuid"`
	Game             string    `gorm:"type:varchar(10);not null;index"                json:"game"`
	Platform         string    `gorm:"type:varchar(16);not null;index"                json:"platform"`
	ChampionID       int       `gorm:"not null;index:idx_mastery_champ_platform"      json:"champion_id"`
	ChampionName     string    `gorm:"type:varchar(64);not null"                      json:"champion_name"`
	AvgMasteryPoints int64     `gorm:"not null" json:"avg_mastery_points"`
	AvgMasteryLevel  float64   `gorm:"type:numeric(4,2);not null" json:"avg_mastery_level"`
	TotalPlayers     int       `gorm:"not null" json:"total_players"`
	TopMasteryPoints int64     `gorm:"not null" json:"top_mastery_points"`
	SampleSize       int       `gorm:"not null" json:"sample_size"`
	CreatedAt        time.Time `gorm:"type:timestamptz;not null" json:"created_at"`
	UpdatedAt        time.Time `gorm:"type:timestamptz;not null" json:"updated_at"`
}

func (ChampionMasteryStats) TableName() string { return "app.champion_mastery_stats" }

// MetaGameAnalysis representa un análisis completo de meta-game
type MetaGameAnalysis struct {
	UUID         uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"uuid"`
	Game         string         `gorm:"type:varchar(10);not null;index"                json:"game"`
	Platform     string         `gorm:"type:varchar(16);not null;index"                json:"platform"`
	Patch        string         `gorm:"type:varchar(32);not null;index"                json:"patch"`
	AnalysisType string         `gorm:"type:varchar(32);not null"                      json:"analysis_type"` // weekly, monthly, patch
	TimeFrame    string         `gorm:"type:varchar(16);not null"                      json:"time_frame"`    // 7d, 30d, patch
	Data         datatypes.JSON `gorm:"type:jsonb;not null" json:"data"`
	Insights     datatypes.JSON `gorm:"type:jsonb" json:"insights,omitempty"`
	CreatedAt    time.Time      `gorm:"type:timestamptz;not null" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"type:timestamptz;not null" json:"updated_at"`
}

func (MetaGameAnalysis) TableName() string { return "app.meta_game_analyses" }
