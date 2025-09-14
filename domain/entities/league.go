package entities

import (
	"time"

	"github.com/google/uuid"
)

type League struct {
	UUID   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"uuid"`
	Game   string    `gorm:"type:varchar(10);not null;index"                json:"game"`
	Region string    `gorm:"type:varchar(16);not null;index"                json:"region"`
	Name   string    `gorm:"type:varchar(120);not null"                     json:"name"`
	Slug   string    `gorm:"type:varchar(120);not null;uniqueIndex"         json:"slug"`

	CreatedAt time.Time `gorm:"type:timestamptz;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null" json:"updated_at"`
}

func (League) TableName() string { return "app.leagues" }

// ProfessionalLeague representa una liga profesional de League of Legends
type ProfessionalLeague struct {
	Code        string    `gorm:"type:varchar(10);primaryKey" json:"code"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Region      string    `gorm:"type:varchar(50);not null" json:"region"`
	Platform    string    `gorm:"type:varchar(20);not null" json:"platform"`
	Seasons     string    `gorm:"type:text" json:"seasons"` // JSON array
	Teams       int       `gorm:"not null" json:"teams"`
	Description string    `gorm:"type:text" json:"description"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `gorm:"type:timestamptz;not null" json:"created_at"`
	UpdatedAt   time.Time `gorm:"type:timestamptz;not null" json:"updated_at"`
}

func (ProfessionalLeague) TableName() string { return "app.professional_leagues" }

// LeagueChampionStats representa estadísticas de campeones en una liga profesional
type LeagueChampionStats struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	LeagueCode    string    `gorm:"type:varchar(10);not null;index" json:"league_code"`
	ChampionName  string    `gorm:"type:varchar(50);not null;index" json:"champion_name"`
	PickRate      float64   `gorm:"type:numeric(5,2)" json:"pick_rate"`
	WinRate       float64   `gorm:"type:numeric(5,2)" json:"win_rate"`
	BanRate       float64   `gorm:"type:numeric(5,2)" json:"ban_rate"`
	Position      string    `gorm:"type:varchar(20)" json:"position"`
	Season        string    `gorm:"type:varchar(20);not null" json:"season"`
	GamesAnalyzed int       `gorm:"not null" json:"games_analyzed"`
	LastUpdated   time.Time `gorm:"type:timestamptz;not null" json:"last_updated"`
	CreatedAt     time.Time `gorm:"type:timestamptz;not null" json:"created_at"`
	UpdatedAt     time.Time `gorm:"type:timestamptz;not null" json:"updated_at"`

	// Foreign key relationship
	ProfessionalLeague ProfessionalLeague `gorm:"foreignKey:LeagueCode;references:Code" json:"-"`
}

func (LeagueChampionStats) TableName() string { return "app.league_champion_stats" }
