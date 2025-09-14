package entities

import (
	"time"

	"github.com/google/uuid"
)

type StreamSource struct {
	UUID     uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"uuid"`
	Name     string    `gorm:"type:varchar(80);not null;uniqueIndex"         json:"name"` // twitch, youtube
	BaseURL  string    `gorm:"type:text;not null"                            json:"base_url"`
	ApiKey   string    `gorm:"type:text"                                     json:"api_key,omitempty"`
	IsActive bool      `gorm:"not null;default:true"                        json:"is_active"`

	CreatedAt time.Time `gorm:"type:timestamptz;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null" json:"updated_at"`
}

func (StreamSource) TableName() string { return "app.stream_sources" }
