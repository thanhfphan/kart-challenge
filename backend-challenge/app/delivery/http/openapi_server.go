package http

//go:generate oapi-codegen --config=../../../oapi-codegen.yaml ../../../openapi.yaml

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhfphan/kart-challenge/app/delivery/http/openapi"
	"github.com/thanhfphan/kart-challenge/app/models"
	"github.com/thanhfphan/kart-challenge/app/usecases"
	"github.com/thanhfphan/kart-challenge/pkg/logging"
)

// OpenAPIServer implements the generated ServerInterface
type OpenAPIServer struct {
	productUC usecases.Product
	// Add other use cases as needed
	// orderUC   usecases.Order
}

// NewOpenAPIServer creates a new OpenAPI server implementation
func NewOpenAPIServer(ucs usecases.UseCase) *OpenAPIServer {
	return &OpenAPIServer{
		productUC: ucs.Product(),
		// orderUC:   ucs.Order(), // Add when order use case is implemented
	}
}

// Ensure OpenAPIServer implements the ServerInterface
var _ openapi.ServerInterface = (*OpenAPIServer)(nil)

// ListProducts implements GET /product
func (s *OpenAPIServer) ListProducts(c *gin.Context) {
	ctx := c.Request.Context()
	log := logging.FromContext(ctx)

	log.Info("Listing all products")

	// TODO: Implement product listing logic
	// For now, return a placeholder response
	// This should call something like: products, err := s.productUC.List(ctx)

	products := []openapi.Product{
		{
			Id:       stringPtr("1"),
			Name:     stringPtr("Sample Product"),
			Price:    float32Ptr(10.99),
			Category: stringPtr("Sample Category"),
		},
	}

	c.JSON(http.StatusOK, products)
}

// GetProduct implements GET /product/{productId}
func (s *OpenAPIServer) GetProduct(c *gin.Context, productId int64) {
	ctx := c.Request.Context()
	log := logging.FromContext(ctx)

	log.Infof("Getting product with ID: %d", productId)

	// Call the existing product use case
	product, err := s.productUC.Get(ctx, productId)
	if err != nil {
		log.Errorf("Failed to get product: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Convert internal model to OpenAPI model
	openAPIProduct := convertToOpenAPIProduct(product)
	c.JSON(http.StatusOK, openAPIProduct)
}

// PlaceOrder implements POST /order
func (s *OpenAPIServer) PlaceOrder(c *gin.Context) {
	ctx := c.Request.Context()
	log := logging.FromContext(ctx)

	var orderReq openapi.OrderReq
	if err := c.ShouldBindJSON(&orderReq); err != nil {
		log.Warnf("Failed to bind order request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	log.Infof("Placing order with %d items", len(orderReq.Items))

	// TODO: Implement order placement logic
	// This should call something like: order, err := s.orderUC.PlaceOrder(ctx, orderReq)

	// For now, return a placeholder response
	order := openapi.Order{
		Id:        stringPtr("order-123"),
		Total:     float32Ptr(99.99),
		Discounts: float32Ptr(0.0),
		Items: &[]struct {
			ProductId *string `json:"productId,omitempty"`
			Quantity  *int    `json:"quantity,omitempty"`
		}{},
		Products: &[]openapi.Product{},
	}

	c.JSON(http.StatusOK, order)
}

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
		Id:   stringPtr(fmt.Sprintf("%d", product.ID)), // Convert int64 ID to string
		Name: stringPtr(product.Name),
		// Note: Internal model doesn't have Price and Category fields
		// You may need to extend the internal model or fetch additional data
		Price:    nil, // TODO: Add price field to internal model or fetch from another source
		Category: nil, // TODO: Add category field to internal model or fetch from another source
		Image: &struct {
			Desktop   *string `json:"desktop,omitempty"`
			Mobile    *string `json:"mobile,omitempty"`
			Tablet    *string `json:"tablet,omitempty"`
			Thumbnail *string `json:"thumbnail,omitempty"`
		}{
			// TODO: Add image fields to internal model or fetch from another source
			Thumbnail: nil,
			Mobile:    nil,
			Tablet:    nil,
			Desktop:   nil,
		},
	}
}
