package models

import "fmt"

type Product struct {
	ID           int64   `gorm:"primaryKey;column:id" json:"id"`
	SKU          string  `json:"sku"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	Category     string  `json:"category"`
	ThumbnailURL string  `json:"thumbnail_url"`
	MobileURL    string  `json:"mobile_url"`
	TabletURL    string  `json:"tablet_url"`
	DesktopURL   string  `json:"desktop_url"`
	Description  string  `json:"description"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

func (*Product) TableName() string {
	return "products"
}

func (c *Product) CacheKey() string {
	return fmt.Sprintf("products:%d", c.ID)
}
