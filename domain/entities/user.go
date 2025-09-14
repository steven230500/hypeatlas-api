package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UUID     uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"uuid"`
	Email    string    `gorm:"type:varchar(120);uniqueIndex;not null"         json:"email"`
	Role     string    `gorm:"type:varchar(16);not null;default:'user'"       json:"role"` // user|admin
	ApiKey   string    `gorm:"type:varchar(64);uniqueIndex"                   json:"api_key,omitempty"`
	Verified bool      `gorm:"not null;default:false"                         json:"verified"`

	CreatedAt time.Time `gorm:"type:timestamptz;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null" json:"updated_at"`
}

func (User) TableName() string { return "app.users" }
