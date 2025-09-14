package entities

import (
	"time"

	"github.com/google/uuid"
)

type Patch struct {
	UUID       uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"uuid"`
	Game       string     `gorm:"type:varchar(10);not null;index:idx_patches_game_rel_ver,priority:1" json:"game"` // val|lol
	Version    string     `gorm:"type:varchar(32);not null;index:idx_patches_game_rel_ver,priority:3" json:"version"`
	ReleasedAt *time.Time `gorm:"type:timestamptz;index:idx_patches_game_rel_ver,priority:2,sort:desc" json:"released_at"`

	CreatedAt time.Time `gorm:"type:timestamptz;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null" json:"updated_at"`

	Changes []*PatchChange `gorm:"foreignKey:PatchUUID" json:"changes,omitempty"`
}

func (Patch) TableName() string { return "app.patches" }
