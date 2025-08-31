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

func (r *promoCode) GetByCode(ctx context.Context, code string) (*models.PromoCode, error) {
	log := logging.FromContext(ctx)

	record := &models.PromoCode{Code: code}
	err := r.cache.GetByCacheKey(ctx, record)
	if err == nil {
		return record, nil
	}

	record, err = r.promoCodeSQL.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	if err := r.cache.Create(ctx, record); err != nil {
		log.Warnf("create cache key=%s err=%v", record.CacheKey(), err)
	}

	return record, nil
}

func (r *promoCode) ValidateCode(ctx context.Context, code string) (bool, float64, error) {
	log := logging.FromContext(ctx)

	// Basic validation: length between 8 and 10 characters
	if len(code) < 8 || len(code) > 10 {
		log.Warnf("Invalid promo code length: %s (length: %d)", code, len(code))
		return false, 0, nil
	}

	// Get promo code from database
	promoCode, err := r.GetByCode(ctx, code)
	if err != nil {
		if errors.Is(err, xerror.ErrRecordNotFound) {
			log.Warnf("Promo code not found: %s", code)
			return false, 0, nil
		}
		return false, 0, err
	}

	// Check if promo code is active
	if !promoCode.IsActive {
		log.Warnf("Promo code is inactive: %s", code)
		return false, 0, nil
	}

	log.Infof("Valid promo code: %s, discount: %.2f%%", code, promoCode.DiscountPct)
	return true, promoCode.DiscountPct, nil
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

func (r *promoCodeSQL) GetByCode(ctx context.Context, code string) (*models.PromoCode, error) {
	var record models.PromoCode
	err := r.dbWithContext(ctx).Where("code = ? AND is_active = ?", code, true).First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerror.ErrRecordNotFound
		}
		return nil, err
	}

	return &record, nil
}
