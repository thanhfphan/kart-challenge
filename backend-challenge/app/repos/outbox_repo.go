package repos

import (
	"context"
	"errors"
	"time"

	"github.com/thanhfphan/kart-challenge/app/models"
	"github.com/thanhfphan/kart-challenge/pkg/cache"
	"github.com/thanhfphan/kart-challenge/pkg/logging"
	"github.com/thanhfphan/kart-challenge/pkg/xerror"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var _ Outbox = (*outbox)(nil)

type outbox struct {
	*outboxSQL
	cache *cache.RedisRepo
}

func newOutbox(redisClient *redis.Client, db *gorm.DB) *outbox {
	return &outbox{
		outboxSQL: newOutboxSQL(db),
		cache:     cache.NewRedisRepo(redisClient, 10*time.Minute), // Cache outbox events for 10 minutes
	}
}

func (r *outbox) GetByID(ctx context.Context, id int64) (*models.OutboxEvent, error) {
	log := logging.FromContext(ctx)

	record := &models.OutboxEvent{ID: id}
	err := r.cache.GetByCacheKey(ctx, record)
	if err == nil {
		return record, nil
	}

	record, err = r.outboxSQL.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := r.cache.Create(ctx, record); err != nil {
		log.Warnf("create cache key=%s err=%v", record.CacheKey(), err)
	}

	return record, nil
}

func (r *outbox) Create(ctx context.Context, record *models.OutboxEvent) (*models.OutboxEvent, error) {
	log := logging.FromContext(ctx)

	record, err := r.outboxSQL.Create(ctx, record)
	if err != nil {
		return nil, err
	}

	if err := r.cache.Create(ctx, record); err != nil {
		log.Warnf("create cache key=%s err=%v", record.CacheKey(), err)
	}

	return record, nil
}

func (r *outbox) GetUnprocessedEvents(ctx context.Context, limit int) ([]*models.OutboxEvent, error) {
	// Don't cache unprocessed events as they change frequently
	return r.outboxSQL.GetUnprocessedEvents(ctx, limit)
}

func (r *outbox) MarkAsProcessed(ctx context.Context, id int64) error {
	log := logging.FromContext(ctx)

	err := r.outboxSQL.MarkAsProcessed(ctx, id)
	if err != nil {
		return err
	}

	// Update cache
	record := &models.OutboxEvent{ID: id}
	now := time.Now().Unix()
	params := map[string]interface{}{
		"status":       models.OutboxEventStatusProcessed,
		"processed_at": now,
	}

	if err := r.cache.UpdateWithMap(ctx, record, params); err != nil {
		log.Warnf("update cache key=%s err=%v", record.CacheKey(), err)
	}

	return nil
}

func (r *outbox) MarkAsFailed(ctx context.Context, id int64) error {
	log := logging.FromContext(ctx)

	err := r.outboxSQL.MarkAsFailed(ctx, id)
	if err != nil {
		return err
	}

	// Update cache
	record := &models.OutboxEvent{ID: id}
	params := map[string]interface{}{
		"status": models.OutboxEventStatusFailed,
	}

	if err := r.cache.UpdateWithMap(ctx, record, params); err != nil {
		log.Warnf("update cache key=%s err=%v", record.CacheKey(), err)
	}

	return nil
}

func (r *outbox) UpdateWithMap(ctx context.Context, record *models.OutboxEvent, params map[string]interface{}) error {
	log := logging.FromContext(ctx)

	err := r.outboxSQL.UpdateWithMap(ctx, record, params)
	if err != nil {
		return err
	}

	if err := r.cache.UpdateWithMap(ctx, record, params); err != nil {
		log.Warnf("update withmap cache err=%v", err)
	}

	return nil
}

// ************* Outbox SQL *************
type outboxSQL struct {
	db *gorm.DB
}

func newOutboxSQL(db *gorm.DB) *outboxSQL {
	return &outboxSQL{
		db: db,
	}
}

func (r *outboxSQL) dbWithContext(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

func (r *outboxSQL) Create(ctx context.Context, record *models.OutboxEvent) (*models.OutboxEvent, error) {
	err := r.dbWithContext(ctx).Create(record).Error
	if err != nil {
		return nil, err
	}

	return record, nil
}

func (r *outboxSQL) GetByID(ctx context.Context, id int64) (*models.OutboxEvent, error) {
	var record models.OutboxEvent
	err := r.dbWithContext(ctx).Where("id = ?", id).First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerror.ErrRecordNotFound
		}
		return nil, err
	}

	return &record, nil
}

func (r *outboxSQL) GetUnprocessedEvents(ctx context.Context, limit int) ([]*models.OutboxEvent, error) {
	var records []*models.OutboxEvent
	err := r.dbWithContext(ctx).
		Where("status = ?", models.OutboxEventStatusPending).
		Order("created_at ASC").
		Limit(limit).
		Find(&records).Error
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (r *outboxSQL) MarkAsProcessed(ctx context.Context, id int64) error {
	now := time.Now().Unix()
	return r.dbWithContext(ctx).
		Model(&models.OutboxEvent{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       models.OutboxEventStatusProcessed,
			"processed_at": now,
		}).Error
}

func (r *outboxSQL) MarkAsFailed(ctx context.Context, id int64) error {
	return r.dbWithContext(ctx).
		Model(&models.OutboxEvent{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status": models.OutboxEventStatusFailed,
		}).Error
}

func (r *outboxSQL) UpdateWithMap(ctx context.Context, record *models.OutboxEvent, params map[string]interface{}) error {
	return r.dbWithContext(ctx).Model(record).Updates(params).Error
}
