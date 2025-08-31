package models

import "fmt"

type OrderItem struct {
	ID        int64   `gorm:"primaryKey;column:id" json:"id"`
	OrderID   string  `gorm:"column:order_id;type:varchar(36)" json:"order_id"`
	ProductID int64   `gorm:"column:product_id" json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"` // Price at time of order

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

func (*OrderItem) TableName() string {
	return "order_items"
}

func (oi *OrderItem) CacheKey() string {
	return fmt.Sprintf("order_items:%d", oi.ID)
}
