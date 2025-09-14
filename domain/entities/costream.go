package entities

import (
	"time"

	"github.com/google/uuid"
)

type CoStream struct {
	UUID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"uuid"`
	EventUUID   uuid.UUID `gorm:"type:uuid;not null;index"                       json:"event_uuid"`
	CreatorUUID uuid.UUID `gorm:"type:uuid;not null;index"                       json:"creator_uuid"`

	Platform string `gorm:"type:varchar(16);not null;index"                           json:"platform"` // twitch|youtube
	URL      string `gorm:"type:text;not null"                                         json:"url"`
	Lang     string `gorm:"type:varchar(8);not null;index:idx_costreams_event_live_lang_viewers,priority:3" json:"lang"`
	Country  string `gorm:"type:varchar(8)"                                            json:"country"`

	Viewers  int  `gorm:"type:int;default:0;index:idx_costreams_event_live_lang_viewers,priority:4,sort:desc" json:"viewers"`
	Verified bool `gorm:"not null;default:false"                                                               json:"verified"`
	IsLive   bool `gorm:"not null;default:false;index:idx_costreams_event_live_lang_viewers,priority:2;index" json:"is_live"`

	LastSeenAt time.Time `gorm:"type:timestamptz;index" json:"last_seen_at"`

	CreatedAt time.Time `gorm:"type:timestamptz;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null" json:"updated_at"`

	Event   *Event   `gorm:"foreignKey:EventUUID"   json:"event,omitempty"`
	Creator *Creator `gorm:"foreignKey:CreatorUUID" json:"creator,omitempty"`
}

func (CoStream) TableName() string { return "app.co_streams" }
