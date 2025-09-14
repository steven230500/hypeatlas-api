package entities

import (
	"time"

	"github.com/google/uuid"
)

type Creator struct {
	UUID      uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"uuid"`
	Platform  string    `gorm:"type:varchar(16);not null;uniqueIndex:uq_creator_platform_handle,priority:1;index:idx_creators_platform_verified,priority:1" json:"platform"`
	Handle    string    `gorm:"type:varchar(120);not null;uniqueIndex:uq_creator_platform_handle,priority:2"                                                         json:"handle"`
	URL       string    `gorm:"type:text;not null"                                                                                                                   json:"url"`
	Lang      string    `gorm:"type:varchar(8);not null"                                                                                                             json:"lang"`
	Country   string    `gorm:"type:varchar(8)"                                                                                                                      json:"country"`
	Verified  bool      `gorm:"not null;default:false;index:idx_creators_platform_verified,priority:2"                                                               json:"verified"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null"            json:"updated_at"`

	CoStreams []*CoStream `gorm:"foreignKey:CreatorUUID" json:"-"`
}

func (Creator) TableName() string { return "app.creators" }
