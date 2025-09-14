package entities

import (
	"time"

	"github.com/google/uuid"
)

type IngestionLog struct {
	UUID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"uuid"`
	Source      string    `gorm:"type:varchar(16);not null;index"                json:"source"`      // twitch|api
	EntityType  string    `gorm:"type:varchar(24);not null;index"                json:"entity_type"` // costream|comp|event
	EntityID    string    `gorm:"type:varchar(80);not null;index"                json:"entity_id"`
	Status      string    `gorm:"type:varchar(16);not null;index"                json:"status"` // success|error
	ErrorMsg    string    `gorm:"type:text"                                      json:"error_msg,omitempty"`
	ProcessedAt time.Time `gorm:"type:timestamptz;not null"                      json:"processed_at"`

	CreatedAt time.Time `gorm:"type:timestamptz;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null" json:"updated_at"`
}

func (IngestionLog) TableName() string { return "app.ingestion_logs" }
