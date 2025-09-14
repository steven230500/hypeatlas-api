package entities

import (
	"time"

	"github.com/google/uuid"
)

type EventRule struct {
	UUID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"uuid"`
	EventID    uuid.UUID `gorm:"column:event_id;type:uuid;not null;index"       json:"event_id"`
	Platform   string    `gorm:"type:varchar(16);not null;index"                json:"platform"` // twitch, youtube
	Handle     string    `gorm:"type:varchar(120);not null;index"               json:"handle"`
	Keyword    string    `gorm:"type:varchar(120);index"                        json:"keyword,omitempty"` // palabra clave en título
	AutoAssign bool      `gorm:"not null;default:true"                          json:"auto_assign"`
	Priority   int       `gorm:"type:int;not null;default:1"                    json:"priority"`

	CreatedAt time.Time `gorm:"type:timestamptz;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null" json:"updated_at"`

	Event *Event `gorm:"foreignKey:EventID;references:UUID" json:"event,omitempty"`
}

func (EventRule) TableName() string { return "app.event_rules" }
