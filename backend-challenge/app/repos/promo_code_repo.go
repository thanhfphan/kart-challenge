package repos

import (
	"context"
	"time"

	"github.com/thanhfphan/kart-challenge/pkg/cache"
	"github.com/thanhfphan/kart-challenge/pkg/logging"

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

func (r *promoCode) ValidateCode(ctx context.Context, code string) (bool, error) {
	log := logging.FromContext(ctx)

	// Basic validation: length between 8 and 10 characters
	if len(code) < 8 || len(code) > 10 {
		log.Warnf("Invalid promo code length: %s (length: %d)", code, len(code))
		return false, nil
	}

	return true, nil
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
