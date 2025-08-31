package models

import "fmt"

type PromoCode struct {
	ID          int64   `gorm:"primaryKey;column:id" json:"id"`
	Code        string  `gorm:"uniqueIndex;column:code" json:"code"`
	Description string  `json:"description"`
	DiscountPct float64 `gorm:"column:discount_pct" json:"discount_pct"` // Percentage discount (0-100)
	IsActive    bool    `gorm:"column:is_active;default:true" json:"is_active"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

func (*PromoCode) TableName() string {
	return "promo_codes"
}

func (pc *PromoCode) CacheKey() string {
	return fmt.Sprintf("promo_codes:%s", pc.Code)
}
