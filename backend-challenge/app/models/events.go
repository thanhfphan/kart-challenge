package models

import "time"

// OrderEventData represents the common structure for order-related events
type OrderEventData struct {
	OrderID    string                 `json:"order_id"`
	Total      float64                `json:"total"`
	Discounts  float64                `json:"discounts"`
	CouponCode string                 `json:"coupon_code,omitempty"`
	Status     string                 `json:"status"`
	Items      []OrderItemEventData   `json:"items"`
	Products   []ProductEventData     `json:"products"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	Timestamp  int64                  `json:"timestamp"`
}

// OrderItemEventData represents order item data in events
type OrderItemEventData struct {
	ID        int64   `json:"id"`
	ProductID int64   `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

// ProductEventData represents product data in events
type ProductEventData struct {
	ID          int64   `json:"id"`
	SKU         string  `json:"sku"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
	Description string  `json:"description,omitempty"`
}

// OrderPlacedEventData represents data for order placed events
type OrderPlacedEventData struct {
	OrderEventData
	PaymentMethod string `json:"payment_method,omitempty"`
	CustomerID    string `json:"customer_id,omitempty"`
}

// OrderCompletedEventData represents data for order completed events
type OrderCompletedEventData struct {
	OrderEventData
	CompletedAt   int64  `json:"completed_at"`
	FulfillmentID string `json:"fulfillment_id,omitempty"`
}

// OrderCancelledEventData represents data for order cancelled events
type OrderCancelledEventData struct {
	OrderEventData
	CancelledAt int64  `json:"cancelled_at"`
	Reason      string `json:"reason,omitempty"`
}

// NewOrderPlacedEventData creates a new OrderPlacedEventData from order and related data
func NewOrderPlacedEventData(order *Order, items []*OrderItem, products []*Product) *OrderPlacedEventData {
	eventItems := make([]OrderItemEventData, len(items))
	for i, item := range items {
		eventItems[i] = OrderItemEventData{
			ID:        item.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	eventProducts := make([]ProductEventData, len(products))
	for i, product := range products {
		eventProducts[i] = ProductEventData{
			ID:          product.ID,
			SKU:         product.SKU,
			Name:        product.Name,
			Price:       product.Price,
			Category:    product.Category,
			Description: product.Description,
		}
	}

	return &OrderPlacedEventData{
		OrderEventData: OrderEventData{
			OrderID:    order.ID,
			Total:      order.Total,
			Discounts:  order.Discounts,
			CouponCode: order.CouponCode,
			Status:     order.Status,
			Items:      eventItems,
			Products:   eventProducts,
			Timestamp:  time.Now().Unix(),
		},
	}
}

// NewOrderCompletedEventData creates a new OrderCompletedEventData from order and related data
func NewOrderCompletedEventData(order *Order, items []*OrderItem, products []*Product) *OrderCompletedEventData {
	eventItems := make([]OrderItemEventData, len(items))
	for i, item := range items {
		eventItems[i] = OrderItemEventData{
			ID:        item.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	eventProducts := make([]ProductEventData, len(products))
	for i, product := range products {
		eventProducts[i] = ProductEventData{
			ID:          product.ID,
			SKU:         product.SKU,
			Name:        product.Name,
			Price:       product.Price,
			Category:    product.Category,
			Description: product.Description,
		}
	}

	return &OrderCompletedEventData{
		OrderEventData: OrderEventData{
			OrderID:    order.ID,
			Total:      order.Total,
			Discounts:  order.Discounts,
			CouponCode: order.CouponCode,
			Status:     order.Status,
			Items:      eventItems,
			Products:   eventProducts,
			Timestamp:  time.Now().Unix(),
		},
		CompletedAt: time.Now().Unix(),
	}
}

// NewOrderCancelledEventData creates a new OrderCancelledEventData from order and related data
func NewOrderCancelledEventData(order *Order, items []*OrderItem, products []*Product, reason string) *OrderCancelledEventData {
	eventItems := make([]OrderItemEventData, len(items))
	for i, item := range items {
		eventItems[i] = OrderItemEventData{
			ID:        item.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	eventProducts := make([]ProductEventData, len(products))
	for i, product := range products {
		eventProducts[i] = ProductEventData{
			ID:          product.ID,
			SKU:         product.SKU,
			Name:        product.Name,
			Price:       product.Price,
			Category:    product.Category,
			Description: product.Description,
		}
	}

	return &OrderCancelledEventData{
		OrderEventData: OrderEventData{
			OrderID:    order.ID,
			Total:      order.Total,
			Discounts:  order.Discounts,
			CouponCode: order.CouponCode,
			Status:     order.Status,
			Items:      eventItems,
			Products:   eventProducts,
			Timestamp:  time.Now().Unix(),
		},
		CancelledAt: time.Now().Unix(),
		Reason:      reason,
	}
}
