package entities

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	UUID    uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"uuid"`
	EventID uuid.UUID  `gorm:"column:event_id;type:uuid;not null;index"       json:"event_id"`
	Type    string     `gorm:"type:varchar(24);not null;index"                json:"type"`    // hype_spike|event_start|comp_update
	Payload string     `gorm:"type:text;not null"                             json:"payload"` // JSON data
	SentAt  *time.Time `gorm:"type:timestamptz"                               json:"sent_at"`

	CreatedAt time.Time `gorm:"type:timestamptz;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null" json:"updated_at"`

	Event *Event `gorm:"foreignKey:EventID;references:UUID" json:"event,omitempty"`
}

func (Notification) TableName() string { return "app.notifications" }
