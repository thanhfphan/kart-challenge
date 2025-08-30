package models

import "fmt"

type Product struct {
	ID   int64  `gorm:"primaryKey;column:id" json:"id"`
	SKU  string `json:"sku"`
	Name string `json:"name"`

	CreatedAt int64 `json:"created_at"`
	UpdateAt  int64 `json:"updated_at"`
}

func (*Product) TableName() string {
	return "products"
}

func (c *Product) CacheKey() string {
	return fmt.Sprintf("products:%d", c.ID)
}
