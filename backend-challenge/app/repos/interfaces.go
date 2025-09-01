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
	Order() Order
	OrderItem() OrderItem
	PromoCode() PromoCode
	Outbox() Outbox
}

type Product interface {
	GetByID(ctx context.Context, id int64) (*models.Product, error)
	List(ctx context.Context) ([]*models.Product, error)
	Create(ctx context.Context, record *models.Product) (*models.Product, error)
	UpdateWithMap(ctx context.Context, record *models.Product, params map[string]interface{}) error
	GetByIDList(ctx context.Context, ids []int64) ([]*models.Product, error)
}

type Order interface {
	GetByID(ctx context.Context, id string) (*models.Order, error)
	Create(ctx context.Context, record *models.Order) (*models.Order, error)
	UpdateWithMap(ctx context.Context, record *models.Order, params map[string]interface{}) error
}

type OrderItem interface {
	Create(ctx context.Context, record *models.OrderItem) (*models.OrderItem, error)
	CreateMany(ctx context.Context, records []*models.OrderItem) error
	GetByOrderID(ctx context.Context, orderID string) ([]*models.OrderItem, error)
	GetByID(ctx context.Context, id int64) (*models.OrderItem, error)
}

type PromoCode interface {
	GetCode(ctx context.Context, code string) (*models.PromoCode, error)
	BulkUpsert(ctx context.Context, promoCodes []*models.PromoCode) error
	UpdateWithMap(ctx context.Context, record *models.PromoCode, params map[string]interface{}) error
}

type Outbox interface {
	Create(ctx context.Context, record *models.OutboxEvent) (*models.OutboxEvent, error)
	GetByID(ctx context.Context, id int64) (*models.OutboxEvent, error)
	GetUnprocessedEvents(ctx context.Context, limit int) ([]*models.OutboxEvent, error)
	MarkAsProcessed(ctx context.Context, id int64) error
	MarkAsFailed(ctx context.Context, id int64) error
	UpdateWithMap(ctx context.Context, record *models.OutboxEvent, params map[string]interface{}) error
}
