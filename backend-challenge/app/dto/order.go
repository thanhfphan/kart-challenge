package dto

// OrderRequest represents the API request for placing an order
// Maps to OpenAPI OrderReq schema
type OrderRequest struct {
	CouponCode string      `json:"couponCode,omitempty"`
	Items      []OrderItem `json:"items" binding:"required,min=1"`
}

// OrderItem represents an item in the order request
type OrderItem struct {
	ProductID string `json:"productId" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

// OrderResponse represents the API response for an order
// Maps to OpenAPI Order schema
type OrderResponse struct {
	ID        string              `json:"id"`
	Total     float64             `json:"total"`
	Discounts float64             `json:"discounts"`
	Items     []OrderItemResponse `json:"items"`
	Products  []ProductResponse   `json:"products"`
}

// OrderItemResponse represents an item in the order response
type OrderItemResponse struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}
