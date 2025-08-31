package usecases

import (
	"context"

	"github.com/thanhfphan/kart-challenge/app/models"
	"github.com/thanhfphan/kart-challenge/app/repos"
	"github.com/thanhfphan/kart-challenge/config"
	"github.com/thanhfphan/kart-challenge/env"
	"github.com/thanhfphan/kart-challenge/pkg/logging"
)

var _ Product = (*product)(nil)

type Product interface {
	Get(ctx context.Context, id int64) (*models.Product, error)
	List(ctx context.Context) ([]*models.Product, error)
}

type product struct {
	cfg *config.Config
	env *env.Env

	productRepo repos.Product
}

func newProduct(cfg *config.Config, env *env.Env, repos repos.Repo) (Product, error) {
	return &product{
		cfg:         cfg,
		env:         env,
		productRepo: repos.Product(),
	}, nil
}

func (u *product) Get(ctx context.Context, id int64) (*models.Product, error) {
	log := logging.FromContext(ctx)
	log.Infof("Get product by id=%d", id)

	product, err := u.productRepo.GetByID(ctx, id)
	if err != nil {
		log.Errorf("Get product by id=%d err=%v", id, err)
		return nil, err
	}

	return product, nil
}

func (u *product) List(ctx context.Context) ([]*models.Product, error) {
	log := logging.FromContext(ctx)
	log.Info("Listing all products")

	products, err := u.productRepo.List(ctx)
	if err != nil {
		log.Errorf("List products err=%v", err)
		return nil, err
	}

	log.Infof("Found %d products", len(products))
	return products, nil
}
