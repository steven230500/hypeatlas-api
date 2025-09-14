package entities

import (
	"time"

	"github.com/google/uuid"
)

type HypeThreshold struct {
	UUID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"uuid"`
	EventID    uuid.UUID `gorm:"column:event_id;type:uuid;not null;index"       json:"event_id"`
	Game       string    `gorm:"type:varchar(10);not null;index"                json:"game"`
	MinViewers int       `gorm:"type:int;not null;default:1000"                json:"min_viewers"`
	MaxViewers int       `gorm:"type:int;not null;default:10000"               json:"max_viewers"`
	AlertLevel string    `gorm:"type:varchar(16);not null;default:'medium'"    json:"alert_level"` // low, medium, high
	IsActive   bool      `gorm:"not null;default:true"                         json:"is_active"`

	CreatedAt time.Time `gorm:"type:timestamptz;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null" json:"updated_at"`

	Event *Event `gorm:"foreignKey:EventID;references:UUID" json:"event,omitempty"`
}

func (HypeThreshold) TableName() string { return "app.hype_thresholds" }
