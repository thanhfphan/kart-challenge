package http

import (
	"fmt"

	"github.com/thanhfphan/kart-challenge/app/delivery/http/openapi"
	"github.com/thanhfphan/kart-challenge/app/dto"
	"github.com/thanhfphan/kart-challenge/app/models"
)

// Helper functions for pointer conversions
func stringPtr(s string) *string {
	return &s
}

func float32Ptr(f float32) *float32 {
	return &f
}

func intPtr(i int) *int {
	return &i
}

// convertToOpenAPIProduct converts internal Product model to OpenAPI Product model
func convertToOpenAPIProduct(product *models.Product) openapi.Product {
	return openapi.Product{
		Id:       stringPtr(fmt.Sprintf("%d", product.ID)),
		Name:     stringPtr(product.Name),
		Price:    float32Ptr(float32(product.Price)),
		Category: stringPtr(product.Category),
		Image: &struct {
			Desktop   *string `json:"desktop,omitempty"`
			Mobile    *string `json:"mobile,omitempty"`
			Tablet    *string `json:"tablet,omitempty"`
			Thumbnail *string `json:"thumbnail,omitempty"`
		}{
			Thumbnail: stringPtr(product.ThumbnailURL),
			Mobile:    stringPtr(product.MobileURL),
			Tablet:    stringPtr(product.TabletURL),
			Desktop:   stringPtr(product.DesktopURL),
		},
	}
}

// convertValidationToDTO converts validation model directly to DTO
func convertValidationToDTO(validationReq *OrderRequestValidation) *dto.OrderRequest {
	items := make([]dto.OrderItem, 0, len(validationReq.Items))
	for _, item := range validationReq.Items {
		items = append(items, dto.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	couponCode := ""
	if validationReq.CouponCode != nil {
		couponCode = *validationReq.CouponCode
	}

	return &dto.OrderRequest{
		CouponCode: couponCode,
		Items:      items,
	}
}

// convertToOpenAPIOrder converts DTO OrderResponse to OpenAPI Order
func convertToOpenAPIOrder(order *dto.OrderResponse) openapi.Order {
	items := make([]struct {
		ProductId *string `json:"productId,omitempty"`
		Quantity  *int    `json:"quantity,omitempty"`
	}, 0, len(order.Items))

	for _, item := range order.Items {
		items = append(items, struct {
			ProductId *string `json:"productId,omitempty"`
			Quantity  *int    `json:"quantity,omitempty"`
		}{
			ProductId: stringPtr(item.ProductID),
			Quantity:  intPtr(item.Quantity),
		})
	}

	products := make([]openapi.Product, 0, len(order.Products))
	for _, product := range order.Products {
		products = append(products, openapi.Product{
			Id:       stringPtr(product.ID),
			Name:     stringPtr(product.Name),
			Price:    float32Ptr(float32(product.Price)),
			Category: stringPtr(product.Category),
			Image: &struct {
				Desktop   *string `json:"desktop,omitempty"`
				Mobile    *string `json:"mobile,omitempty"`
				Tablet    *string `json:"tablet,omitempty"`
				Thumbnail *string `json:"thumbnail,omitempty"`
			}{
				Thumbnail: stringPtr(product.Image.Thumbnail),
				Mobile:    stringPtr(product.Image.Mobile),
				Tablet:    stringPtr(product.Image.Tablet),
				Desktop:   stringPtr(product.Image.Desktop),
			},
		})
	}

	return openapi.Order{
		Id:        stringPtr(order.ID),
		Total:     float32Ptr(float32(order.Total)),
		Discounts: float32Ptr(float32(order.Discounts)),
		Items:     &items,
		Products:  &products,
	}
}
