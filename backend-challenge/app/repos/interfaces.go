package repos

import (
	"context"

	"github.com/thanhfphan/kart-challenge/app/models"
)

//go:generate mockgen -package=repos -destination=interfaces_mock.go -source=interfaces.go

// Repo ...
type Repo interface {
	WithTransaction(ctx context.Context, fn func(Repo) error) error
	Product() Product
}

type Product interface {
	GetByID(ctx context.Context, id int64) (*models.Product, error)
	Create(ctx context.Context, record *models.Product) (*models.Product, error)
	UpdateWithMap(ctx context.Context, record *models.Product, params map[string]interface{}) error
	GetByIDList(ctx context.Context, ids []int64) ([]*models.Product, error)
}
