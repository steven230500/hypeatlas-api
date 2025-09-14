package entities

import (
	"time"

	"github.com/google/uuid"
)

type Metric struct {
	UUID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"uuid"`
	CoStreamID uuid.UUID `gorm:"column:co_stream_id;type:uuid;not null;index"   json:"co_stream_id"`
	MetricType string    `gorm:"type:varchar(24);not null;index"                json:"metric_type"` // viewers|chat_messages|follows|sentiment
	Value      float64   `gorm:"type:numeric(10,2);not null"                    json:"value"`
	RecordedAt time.Time `gorm:"type:timestamptz;not null;index"                json:"recorded_at"`

	CreatedAt time.Time `gorm:"type:timestamptz;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null" json:"updated_at"`

	CoStream *CoStream `gorm:"foreignKey:CoStreamID;references:UUID" json:"co_stream,omitempty"`
}

func (Metric) TableName() string { return "app.metrics" }
