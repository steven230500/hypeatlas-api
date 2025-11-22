package repository

import (
	"context"
	"time"

	"github.com/steven230500/hypeatlas-api/domain/entities"
	out "github.com/steven230500/hypeatlas-api/modules/relay/domain/ports/out"
	"github.com/steven230500/hypeatlas-api/shared/db"
	"gorm.io/gorm"
)

type Repo struct{ db *gorm.DB }

func New(db *gorm.DB) out.Repository { return &Repo{db: db} }

func (r *Repo) FindLiveByEvent(ctx context.Context, eventSlug, lang string) ([]entities.CoStream, error) {
	// Primero encontrar el evento por slug
	var event entities.Event
	result := db.Call(r.db.WithContext(ctx).Where("slug = ?", eventSlug).First(&event))
	if result.Error != nil {
		return nil, result.Error
	}

	// Luego buscar co-streams por event_uuid
	var coStreams []entities.CoStream
	query := r.db.WithContext(ctx).Where("event_uuid = ? AND is_live = true", event.UUID)
	if lang != "" {
		query = query.Where("lang = ?", lang)
	}
	result = db.Call(query.Order("viewers DESC").Find(&coStreams))
	return coStreams, result.Error
}

func (r *Repo) HypeMapLive(ctx context.Context, game, lang string, limit, offset int) ([]entities.HypeMapItem, error) {
	var items []entities.HypeMapItem

	// Asegurar valores por defecto
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	query := r.db.WithContext(ctx).
		Table("app.co_streams c").
		Select(`
			e.slug as event_slug,
			e.title as event_title,
			e.game,
			e.league,
			c.platform,
			cr.handle,
			c.lang,
			c.country,
			c.viewers,
			c.is_live,
			c.viewers as score
		`).
		Joins("JOIN app.events e ON e.uuid = c.event_uuid").
		Joins("JOIN app.creators cr ON cr.uuid = c.creator_uuid").
		Where("c.is_live = ?", true)

	if game != "" {
		query = query.Where("e.game = ?", game)
	}
	if lang != "" {
		query = query.Where("c.lang = ?", lang)
	}

	result := db.Call(query.Order("c.viewers DESC").Limit(limit).Offset(offset).Scan(&items))
	return items, result.Error
}

func (r *Repo) HypeMapSummary(ctx context.Context, game, lang string, limit, offset int) ([]entities.HypeMapSummaryItem, error) {
	var items []entities.HypeMapSummaryItem

	// Asegurar valores por defecto
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	query := r.db.WithContext(ctx).
		Table("app.co_streams c").
		Select(`
			e.slug as event_slug,
			e.title as event_title,
			e.game,
			e.league,
			COUNT(c.uuid) as streamers,
			SUM(c.viewers) as total_viewers,
			MAX(c.last_seen_at) as last_seen_at
		`).
		Joins("JOIN app.events e ON e.uuid = c.event_uuid").
		Joins("JOIN app.creators cr ON cr.uuid = c.creator_uuid").
		Where("c.is_live = ?", true)

	if game != "" {
		query = query.Where("e.game = ?", game)
	}
	if lang != "" {
		query = query.Where("c.lang = ?", lang)
	}

	query = query.Group("e.uuid, e.slug, e.title, e.game, e.league").
		Order("total_viewers DESC").
		Limit(limit).
		Offset(offset)

	result := db.Call(query.Scan(&items))
	return items, result.Error
}

func (r *Repo) UpsertCoStream(ctx context.Context, eventSlug, eventTitle, game, league string, startsAtNullable *string, platform, handle, url, lang, country string, verified bool, viewers int, isLive bool) error {
	// Primero, encontrar o crear el evento
	var event entities.Event
	result := db.Call(r.db.WithContext(ctx).Where("slug = ?", eventSlug).First(&event))
	if result.Error == gorm.ErrRecordNotFound {
		event = entities.Event{
			Slug:   eventSlug,
			Title:  eventTitle,
			Game:   game,
			League: &league,
		}
		// Parse startsAt if provided
		if startsAtNullable != nil {
			if parsed, err := time.Parse(time.RFC3339, *startsAtNullable); err == nil {
				event.StartsAt = &parsed
			}
		}
		if err := db.Call(r.db.WithContext(ctx).Create(&event)).Error; err != nil {
			return err
		}
	} else if result.Error != nil {
		return result.Error
	}

	// Encontrar o crear el creator
	var creator entities.Creator
	result = db.Call(r.db.WithContext(ctx).Where("platform = ? AND handle = ?", platform, handle).First(&creator))
	if result.Error == gorm.ErrRecordNotFound {
		creator = entities.Creator{
			Platform: platform,
			Handle:   handle,
			URL:      url,
			Lang:     lang,
			Country:  country,
			Verified: verified,
		}
		if err := db.Call(r.db.WithContext(ctx).Create(&creator)).Error; err != nil {
			return err
		}
	} else if result.Error != nil {
		return result.Error
	}

	// Luego, upsert el co-stream
	coStream := entities.CoStream{
		EventUUID:   event.UUID,
		CreatorUUID: creator.UUID,
		Platform:    platform,
		URL:         url,
		Lang:        lang,
		Country:     country,
		Viewers:     viewers,
		Verified:    verified,
		IsLive:      isLive,
		LastSeenAt:  time.Now(),
	}

	return db.Call(r.db.WithContext(ctx).Where("event_uuid = ? AND creator_uuid = ?", event.UUID, creator.UUID).
		Assign(coStream).
		FirstOrCreate(&coStream)).Error
}

func (r *Repo) MarkStaleCoStreamsOffline(ctx context.Context, olderThan time.Duration) (int64, error) {
	result := db.Call(r.db.WithContext(ctx).Model(&entities.CoStream{}).
		Where("last_seen_at < ? AND is_live = true", time.Now().Add(-olderThan)).
		Update("is_live", false))
	return result.RowsAffected, result.Error
}

func (r *Repo) LoadStreamRules(ctx context.Context) (map[string]string, error) {
	var rules []entities.EventStreamRule
	result := db.Call(r.db.WithContext(ctx).Find(&rules))
	if result.Error != nil {
		return nil, result.Error
	}
	rulesMap := make(map[string]string)
	for _, rule := range rules {
		key := rule.Platform + ":" + rule.Handle
		rulesMap[key] = rule.EventSlug
	}
	return rulesMap, nil
}

func (r *Repo) ActiveWindows(ctx context.Context, now time.Time) ([]entities.EventWindow, error) {
	var windows []entities.EventWindow
	result := db.Call(r.db.WithContext(ctx).Where("starts_at <= ? AND ends_at >= ?", now, now).Find(&windows))
	return windows, result.Error
}

func (r *Repo) ListCreatorHandles(ctx context.Context, platform string, verified bool) ([]entities.Creator, error) {
	var creators []entities.Creator
	query := r.db.WithContext(ctx).Where("platform = ?", platform)
	if verified {
		query = query.Where("verified = ?", verified)
	}
	result := db.Call(query.Find(&creators))
	return creators, result.Error
}
