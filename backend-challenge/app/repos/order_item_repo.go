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

var _ OrderItem = (*orderItem)(nil)

type orderItem struct {
	*orderItemSQL
	cache *cache.RedisRepo
}

func newOrderItem(redisClient *redis.Client, db *gorm.DB) *orderItem {
	return &orderItem{
		orderItemSQL: newOrderItemSQL(db),
		cache:        cache.NewRedisRepo(redisClient, 5*time.Hour),
	}
}

func (r *orderItem) GetByID(ctx context.Context, id int64) (*models.OrderItem, error) {
	log := logging.FromContext(ctx)

	record := &models.OrderItem{ID: id}
	err := r.cache.GetByCacheKey(ctx, record)
	if err == nil {
		return record, nil
	}

	record, err = r.orderItemSQL.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := r.cache.Create(ctx, record); err != nil {
		log.Warnf("create cache key=%s err=%v", record.CacheKey(), err)
	}

	return record, nil
}

func (r *orderItem) GetByOrderID(ctx context.Context, orderID string) ([]*models.OrderItem, error) {
	var records []*models.OrderItem
	err := r.orderItemSQL.dbWithContext(ctx).Where("order_id = ?", orderID).Find(&records).Error
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (r *orderItem) Create(ctx context.Context, record *models.OrderItem) (*models.OrderItem, error) {
	log := logging.FromContext(ctx)

	record, err := r.orderItemSQL.Create(ctx, record)
	if err != nil {
		return nil, err
	}

	if err := r.cache.Create(ctx, record); err != nil {
		log.Warnf("create cache key=%s err=%v", record.CacheKey(), err)
	}

	return record, nil
}

func (r *orderItem) CreateMany(ctx context.Context, records []*models.OrderItem) error {
	log := logging.FromContext(ctx)

	err := r.orderItemSQL.CreateMany(ctx, records)
	if err != nil {
		return err
	}

	items := make([]cache.Item, len(records))
	for i, record := range records {
		items[i] = record
	}
	if err := r.cache.CreateList(ctx, items); err != nil {
		log.Warnf("create cache list err=%v", err)
	}

	return nil
}

func (r *orderItem) UpdateWithMap(ctx context.Context, record *models.OrderItem, params map[string]interface{}) error {
	log := logging.FromContext(ctx)

	err := r.orderItemSQL.UpdateWithMap(ctx, record, params)
	if err != nil {
		return err
	}

	if err := r.cache.UpdateWithMap(ctx, record, params); err != nil {
		log.Warnf("update withmap cache err=%v", err)
	}

	return nil
}

// ************* OrderItem SQL *************
type orderItemSQL struct {
	db *gorm.DB
}

func newOrderItemSQL(db *gorm.DB) *orderItemSQL {
	return &orderItemSQL{
		db: db,
	}
}

func (r *orderItemSQL) dbWithContext(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

func (r *orderItemSQL) Create(ctx context.Context, record *models.OrderItem) (*models.OrderItem, error) {
	err := r.dbWithContext(ctx).Create(record).Error
	if err != nil {
		return nil, err
	}

	return record, nil
}

func (r *orderItemSQL) CreateMany(ctx context.Context, records []*models.OrderItem) error {
	err := r.dbWithContext(ctx).CreateInBatches(records, 100).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *orderItemSQL) UpdateWithMap(ctx context.Context, record *models.OrderItem, params map[string]interface{}) error {
	return r.dbWithContext(ctx).Model(record).Updates(params).Error
}

func (r *orderItemSQL) GetByID(ctx context.Context, id int64) (*models.OrderItem, error) {
	var record models.OrderItem
	err := r.dbWithContext(ctx).Where("id = ?", id).First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerror.ErrRecordNotFound
		}
		return nil, err
	}

	return &record, nil
}
