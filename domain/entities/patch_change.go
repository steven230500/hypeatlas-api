package entities

import (
	"time"

	"github.com/google/uuid"
)

type PatchChange struct {
	ID          int64     `gorm:"primaryKey"                json:"id"`
	PatchUUID   uuid.UUID `gorm:"type:uuid;not null;index:idx_patch_changes_patch_type,priority:1" json:"patch_uuid"`
	EntityType  string    `gorm:"type:varchar(24);not null;index:idx_patch_changes_patch_type,priority:2" json:"entity_type"` // champion|agent|item|weapon|map
	EntityID    string    `gorm:"type:varchar(80);not null" json:"entity_id"`
	Field       string    `gorm:"type:varchar(80);not null" json:"field"`
	Old         string    `gorm:"type:text"                 json:"old"`
	New         string    `gorm:"type:text"                 json:"new"`
	ImpactScore float64   `gorm:"type:numeric(5,2);not null;default:0.00" json:"impact_score"`

	CreatedAt time.Time `gorm:"type:timestamptz;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null" json:"updated_at"`
}

func (PatchChange) TableName() string { return "app.patch_changes" }
