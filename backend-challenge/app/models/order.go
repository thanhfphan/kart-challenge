package models

import (
	"fmt"

	"github.com/google/uuid"
)

type Order struct {
	ID         string  `gorm:"primaryKey;column:id;type:varchar(36)" json:"id"`
	Total      float64 `json:"total"`
	Discounts  float64 `json:"discounts"`
	CouponCode string  `gorm:"column:coupon_code" json:"coupon_code"`
	Status     string  `json:"status"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

func (*Order) TableName() string {
	return "orders"
}

func (o *Order) CacheKey() string {
	return fmt.Sprintf("orders:%s", o.ID)
}

func (o *Order) BeforeCreate() error {
	if o.ID == "" {
		o.ID = uuid.New().String()
	}
	return nil
}
