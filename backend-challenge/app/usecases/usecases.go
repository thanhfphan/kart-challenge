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
	Hello() Hello
}

type useCase struct {
	product Product
	hello   Hello
}

func New(cfg *config.Config, env *env.Env, repos repos.Repo) (UseCase, error) {
	product, err := newProduct(cfg, env, repos)
	if err != nil {
		return nil, fmt.Errorf("new product failed err=%w", err)
	}
	hello, err := newHello(cfg, env, repos)
	if err != nil {
		return nil, fmt.Errorf("new hello failed err=%w", err)
	}

	return &useCase{
		product: product,
		hello:   hello,
	}, nil
}

func (u *useCase) Product() Product {
	return u.product
}

func (u *useCase) Hello() Hello {
	return u.hello
}
