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
