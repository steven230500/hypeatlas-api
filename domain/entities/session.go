package entities

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	UUID        uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"uuid"`
	CoStreamID  uuid.UUID  `gorm:"column:co_stream_id;type:uuid;not null;index"   json:"co_stream_id"`
	StartedAt   time.Time  `gorm:"type:timestamptz;not null;index"                json:"started_at"`
	EndedAt     *time.Time `gorm:"type:timestamptz;index"                        json:"ended_at"`
	PeakViewers int        `gorm:"type:int;default:0"                            json:"peak_viewers"`
	Duration    int        `gorm:"type:int;default:0"                            json:"duration"` // in seconds

	CreatedAt time.Time `gorm:"type:timestamptz;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null" json:"updated_at"`

	CoStream *CoStream `gorm:"foreignKey:CoStreamID;references:UUID" json:"co_stream,omitempty"`
}

func (Session) TableName() string { return "app.sessions" }
