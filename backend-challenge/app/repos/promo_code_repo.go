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
	"gorm.io/gorm/clause"
)

var _ PromoCode = (*promoCode)(nil)

type promoCode struct {
	*promoCodeSQL
	cache *cache.RedisRepo
}

func newPromoCode(redisClient *redis.Client, db *gorm.DB) *promoCode {
	return &promoCode{
		promoCodeSQL: newPromoCodeSQL(db),
		cache:        cache.NewRedisRepo(redisClient, 1*time.Hour), // Cache promo codes for 1 hour
	}
}

func (r *promoCode) GetCode(ctx context.Context, code string) (*models.PromoCode, error) {
	log := logging.FromContext(ctx)

	record := &models.PromoCode{Code: code}
	err := r.cache.GetByCacheKey(ctx, record)
	if err == nil {
		return record, nil
	}

	record, err = r.promoCodeSQL.GetCode(ctx, code)
	if err != nil {
		return nil, err
	}

	if err := r.cache.Create(ctx, record); err != nil {
		log.Warnf("create cache key=%s err=%v", record.CacheKey(), err)
	}

	return record, nil
}

func (r *promoCode) UpdateWithMap(ctx context.Context, record *models.PromoCode, params map[string]interface{}) error {
	log := logging.FromContext(ctx)

	err := r.promoCodeSQL.UpdateWithMap(ctx, record, params)
	if err != nil {
		return err
	}

	if err := r.cache.UpdateWithMap(ctx, record, params); err != nil {
		log.Warnf("update withmap cache err=%v", err)
	}

	return nil
}

func (r *promoCode) BulkUpsert(ctx context.Context, promoCodes []*models.PromoCode) error {
	return r.promoCodeSQL.BulkUpsert(ctx, promoCodes)
}

// ************* PromoCode SQL *************
type promoCodeSQL struct {
	db *gorm.DB
}

func newPromoCodeSQL(db *gorm.DB) *promoCodeSQL {
	return &promoCodeSQL{
		db: db,
	}
}

func (r *promoCodeSQL) dbWithContext(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

func (r *promoCodeSQL) GetCode(ctx context.Context, code string) (*models.PromoCode, error) {
	var record models.PromoCode
	err := r.dbWithContext(ctx).Where("code = ?", code).First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerror.ErrRecordNotFound
		}
		return nil, err
	}

	return &record, nil
}

func (r *promoCodeSQL) UpdateWithMap(ctx context.Context, record *models.PromoCode, params map[string]interface{}) error {
	return r.dbWithContext(ctx).Model(record).Updates(params).Error
}

func (r *promoCodeSQL) BulkUpsert(ctx context.Context, promoCodes []*models.PromoCode) error {
	log := logging.FromContext(ctx)

	if len(promoCodes) == 0 {
		return nil
	}

	result := r.dbWithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "code"}}, // conflict on unique code column
		DoUpdates: clause.AssignmentColumns([]string{
			"description",
			"discount_pct",
			"updated_at",
		}),
	}).CreateInBatches(promoCodes, 1000) // Process in batches of 1000

	if result.Error != nil {
		log.Errorf("Bulk upsert failed: %v", result.Error)
		return result.Error
	}

	return nil
}
