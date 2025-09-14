package entities

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	UUID     uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"uuid"`
	Slug     string     `gorm:"type:varchar(80);uniqueIndex;not null"         json:"slug"`
	Title    string     `gorm:"type:varchar(160);not null"                    json:"title"`
	Game     string     `gorm:"type:varchar(10);not null;index"               json:"game"` // val|lol
	League   *string    `gorm:"type:varchar(80)"                               json:"league"`
	StartsAt *time.Time `gorm:"type:timestamptz;not null"                      json:"starts_at"`
	EndsAt   *time.Time `gorm:"type:timestamptz"                               json:"ends_at"`

	CreatedAt time.Time `gorm:"type:timestamptz;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null" json:"updated_at"`

	CoStreams []*CoStream `gorm:"foreignKey:EventUUID" json:"-"`
}

func (Event) TableName() string { return "app.events" }

// EventWindow (worker: ventanas activas)
type EventWindow struct {
	UUID      uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"uuid"`
	EventSlug string    `gorm:"type:varchar(80);not null;index"                json:"event_slug"`
	StartsAt  time.Time `gorm:"type:timestamptz;not null"                      json:"starts_at"`
	EndsAt    time.Time `gorm:"type:timestamptz;not null"                      json:"ends_at"`
	Region    string    `gorm:"type:varchar(16);default:''"                    json:"region"`
	Lang      string    `gorm:"type:varchar(8);default:''"                     json:"lang"`

	CreatedAt time.Time `gorm:"type:timestamptz;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null" json:"updated_at"`
}

func (EventWindow) TableName() string { return "app.event_windows" }

// EventStreamRule (worker: mapping plataforma/handle → evento)
type EventStreamRule struct {
	UUID      uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"uuid"`
	Platform  string    `gorm:"type:varchar(16);not null;uniqueIndex:uq_platform_handle,priority:1" json:"platform"`
	Handle    string    `gorm:"type:varchar(120);not null;uniqueIndex:uq_platform_handle,priority:2" json:"handle"`
	EventSlug string    `gorm:"type:varchar(80);not null"                                      json:"event_slug"`
	Note      string    `gorm:"type:text;default:''"                                           json:"note"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null"                                      json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null"                                      json:"updated_at"`
}

func (EventStreamRule) TableName() string { return "app.event_stream_rules" }
