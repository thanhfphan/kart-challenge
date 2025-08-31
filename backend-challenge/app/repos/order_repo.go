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

var _ Order = (*order)(nil)

type order struct {
	*orderSQL
	cache *cache.RedisRepo
}

func newOrder(redisClient *redis.Client, db *gorm.DB) *order {
	return &order{
		orderSQL: newOrderSQL(db),
		cache:    cache.NewRedisRepo(redisClient, 30*time.Minute), // Cache orders for 30 minutes
	}
}

func (r *order) GetByID(ctx context.Context, id string) (*models.Order, error) {
	log := logging.FromContext(ctx)

	record := &models.Order{ID: id}
	err := r.cache.GetByCacheKey(ctx, record)
	if err == nil {
		return record, nil
	}

	record, err = r.orderSQL.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := r.cache.Create(ctx, record); err != nil {
		log.Warnf("create cache key=%s err=%v", record.CacheKey(), err)
	}

	return record, nil
}

func (r *order) Create(ctx context.Context, record *models.Order) (*models.Order, error) {
	log := logging.FromContext(ctx)

	record, err := r.orderSQL.Create(ctx, record)
	if err != nil {
		return nil, err
	}

	if err := r.cache.Create(ctx, record); err != nil {
		log.Warnf("create cache key=%s err=%v", record.CacheKey(), err)
	}

	return record, nil
}

func (r *order) UpdateWithMap(ctx context.Context, record *models.Order, params map[string]interface{}) error {
	log := logging.FromContext(ctx)

	err := r.orderSQL.UpdateWithMap(ctx, record, params)
	if err != nil {
		return err
	}

	if err := r.cache.UpdateWithMap(ctx, record, params); err != nil {
		log.Warnf("update withmap cache err=%v", err)
	}

	return nil
}

// ************* Order SQL *************
type orderSQL struct {
	db *gorm.DB
}

func newOrderSQL(db *gorm.DB) *orderSQL {
	return &orderSQL{
		db: db,
	}
}

func (r *orderSQL) dbWithContext(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

func (r *orderSQL) Create(ctx context.Context, record *models.Order) (*models.Order, error) {
	err := r.dbWithContext(ctx).Create(record).Error
	if err != nil {
		return nil, err
	}

	return record, nil
}

func (r *orderSQL) UpdateWithMap(ctx context.Context, record *models.Order, params map[string]interface{}) error {
	return r.dbWithContext(ctx).Model(record).Updates(params).Error
}

func (r *orderSQL) GetByID(ctx context.Context, id string) (*models.Order, error) {
	var record models.Order
	err := r.dbWithContext(ctx).Where("id = ?", id).First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerror.ErrRecordNotFound
		}
		return nil, err
	}

	return &record, nil
}
