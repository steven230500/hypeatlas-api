package entities

import (
	"time"

	"github.com/google/uuid"
)

type Game struct {
	UUID      uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"uuid"`
	Name      string    `gorm:"type:varchar(80);not null"                      json:"name"`
	Slug      string    `gorm:"type:varchar(80);uniqueIndex;not null"         json:"slug"`
	Platforms string    `gorm:"type:text;not null"                            json:"platforms"` // JSON array: ["twitch","youtube"]

	CreatedAt time.Time `gorm:"type:timestamptz;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null" json:"updated_at"`
}

func (Game) TableName() string { return "app.games" }
