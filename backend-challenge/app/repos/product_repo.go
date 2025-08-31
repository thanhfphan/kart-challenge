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

var _ Product = (*product)(nil)

type product struct {
	*productSQL
	cache *cache.RedisRepo
}

func newProduct(redisClient *redis.Client, db *gorm.DB) *product {
	return &product{
		productSQL: newProductSQL(db),
		cache:      cache.NewRedisRepo(redisClient, 24*time.Hour),
	}
}
func (r *product) GetByID(ctx context.Context, id int64) (*models.Product, error) {
	log := logging.FromContext(ctx)

	record := &models.Product{ID: id}
	err := r.cache.GetByCacheKey(ctx, record)
	if err == nil {
		return record, nil
	}

	record, err = r.productSQL.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := r.cache.Create(ctx, record); err != nil {
		log.Warnf("create cache key=%s err=%v", record.CacheKey(), err)
	}

	return record, nil
}

func (r *product) List(ctx context.Context) ([]*models.Product, error) {
	products, err := r.productSQL.List(ctx)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (r *product) Create(ctx context.Context, record *models.Product) (*models.Product, error) {
	log := logging.FromContext(ctx)

	record, err := r.productSQL.Create(ctx, record)
	if err != nil {
		return nil, err
	}

	if err := r.cache.Create(ctx, record); err != nil {
		log.Warnf("create cache key=%s err=%v", record.CacheKey(), err)
	}

	return record, nil
}

func (r *product) UpdateWithMap(ctx context.Context, record *models.Product, params map[string]interface{}) error {
	log := logging.FromContext(ctx)

	err := r.productSQL.UpdateWithMap(ctx, record, params)
	if err != nil {
		return err
	}

	if err := r.cache.UpdateWithMap(ctx, record, params); err != nil {
		log.Warnf("update withmap cache err=%v", err)
	}

	return nil
}

// ************* Product SQL *************
type productSQL struct {
	db *gorm.DB
}

func newProductSQL(db *gorm.DB) *productSQL {
	return &productSQL{
		db: db,
	}
}

func (r *productSQL) dbWithContext(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

func (r *productSQL) Create(ctx context.Context, record *models.Product) (*models.Product, error) {
	err := r.dbWithContext(ctx).Create(record).Error
	if err != nil {
		return nil, err
	}

	return record, nil
}

func (r *productSQL) UpdateWithMap(ctx context.Context, record *models.Product, params map[string]interface{}) error {
	query := r.dbWithContext(ctx).Model(record).Updates(params)

	return query.Error
}

func (r *productSQL) GetByID(ctx context.Context, orderID int64) (*models.Product, error) {
	var record models.Product
	err := r.dbWithContext(ctx).Where("id = ?", orderID).First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerror.ErrRecordNotFound
		}
		return nil, err
	}

	return &record, nil
}

func (r *productSQL) GetByIDList(ctx context.Context, ids []int64) ([]*models.Product, error) {
	var records []*models.Product
	err := r.dbWithContext(ctx).Where("id IN (?)", ids).Find(&records).Error
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (r *productSQL) List(ctx context.Context) ([]*models.Product, error) {
	var records []*models.Product
	err := r.dbWithContext(ctx).Find(&records).Error
	if err != nil {
		return nil, err
	}

	return records, nil
}
