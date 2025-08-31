package usecases

import (
	"fmt"

	"github.com/thanhfphan/kart-challenge/app/repos"
	"github.com/thanhfphan/kart-challenge/config"
	"github.com/thanhfphan/kart-challenge/env"
)

var _ UseCase = (*useCase)(nil)

type UseCase interface {
	Product() Product
	Order() Order
}

type useCase struct {
	product Product
	order   Order
}

func New(cfg *config.Config, env *env.Env, repos repos.Repo) (UseCase, error) {
	product, err := newProduct(cfg, env, repos)
	if err != nil {
		return nil, fmt.Errorf("new product failed err=%w", err)
	}

	order, err := newOrder(cfg, env, repos)
	if err != nil {
		return nil, fmt.Errorf("new order failed err=%w", err)
	}

	return &useCase{
		product: product,
		order:   order,
	}, nil
}

func (u *useCase) Product() Product {
	return u.product
}

func (u *useCase) Order() Order {
	return u.order
}
