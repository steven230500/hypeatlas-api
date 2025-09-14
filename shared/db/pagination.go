package db

import (
	"time"

	"gorm.io/gorm"
)

func Paginate(db *gorm.DB, lastCreatedAt time.Time, limit int, table string) *gorm.DB {
	column := "created_at"
	if table != "" {
		column = table + "." + column
	}
	if limit <= 0 {
		limit = 50
	}
	return db.Order(column+" DESC").Where(column+" < ?", lastCreatedAt).Limit(limit)
}
