package entities

import "time"

type HypeMapItem struct {
	EventSlug  string `json:"event_slug"  gorm:"column:event_slug"`
	EventTitle string `json:"event_title" gorm:"column:event_title"`
	Game       string `json:"game"        gorm:"column:game"`
	League     string `json:"league"      gorm:"column:league"`
	Platform   string `json:"platform"    gorm:"column:platform"`
	Handle     string `json:"handle"      gorm:"column:handle"`
	Lang       string `json:"lang"        gorm:"column:lang"`
	Country    string `json:"country"     gorm:"column:country"`
	Viewers    int    `json:"viewers"     gorm:"column:viewers"`
	IsLive     bool   `json:"is_live"     gorm:"column:is_live"`
	Score      int    `json:"score"       gorm:"column:score"`
}

type HypeMapSummaryItem struct {
	EventSlug    string    `json:"event_slug"     gorm:"column:event_slug"`
	EventTitle   string    `json:"event_title"    gorm:"column:event_title"`
	Game         string    `json:"game"           gorm:"column:game"`
	League       string    `json:"league"         gorm:"column:league"`
	Streamers    int       `json:"streamers"      gorm:"column:streamers"`
	TotalViewers int       `json:"total_viewers"  gorm:"column:total_viewers"`
	LastSeenAt   time.Time `json:"last_seen_at"   gorm:"column:last_seen_at"`
}
