package repos

import (
	"context"

	"github.com/thanhfphan/kart-challenge/config"
	"github.com/thanhfphan/kart-challenge/env"
	"github.com/thanhfphan/kart-challenge/pkg/logging"

	"gorm.io/gorm"
)

var _ Repo = (*repo)(nil)

type repo struct {
	cfg *config.Config
	db  *gorm.DB
	env *env.Env

	product   Product
	order     Order
	orderItem OrderItem
	promoCode PromoCode
}

// New returns new instance of Repo
// IMPORTANT: only pass the db instance from outside, NOT use the env.Database() method
func New(cfg *config.Config, env *env.Env, db *gorm.DB) Repo {
	return &repo{
		cfg: cfg,
		db:  db,
		env: env,

		product:   newProduct(env.RedisClient(), db),
		order:     newOrder(env.RedisClient(), db),
		orderItem: newOrderItem(env.RedisClient(), db),
		promoCode: newPromoCode(env.RedisClient(), db),
	}
}

// WithTransaction ...
func (r *repo) WithTransaction(ctx context.Context, fn func(Repo) error) (err error) {
	log := logging.FromContext(ctx)
	log.Info("Starting transaction")

	tx := r.db.Begin()
	tr := New(r.cfg, r.env, tx)

	err = tx.Error
	if err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil { // nolint
			log.Warnf("Transaction failed with panic: %+v", p)
			// a panic occurred, rollback and repanic
			tx.Rollback()
			panic(p)
		} else if err != nil {
			log.Warnf("Transaction failure with error: %+v", err)
			// something went wrong, rollback
			tx.Rollback()
		} else {
			log.Info("Finishing transaction")
			// all good, commit
			err = tx.Commit().Error
		}
	}()

	err = fn(tr)

	return err
}

func (r *repo) Product() Product {
	return r.product
}

func (r *repo) Order() Order {
	return r.order
}

func (r *repo) OrderItem() OrderItem {
	return r.orderItem
}

func (r *repo) PromoCode() PromoCode {
	return r.promoCode
}
