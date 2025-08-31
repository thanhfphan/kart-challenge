package http

//go:generate oapi-codegen --config=../../../oapi-codegen.yaml ../../../openapi.yaml

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhfphan/kart-challenge/app/delivery/http/openapi"
	"github.com/thanhfphan/kart-challenge/app/dto"
	"github.com/thanhfphan/kart-challenge/app/models"
	"github.com/thanhfphan/kart-challenge/app/usecases"
	"github.com/thanhfphan/kart-challenge/pkg/logging"
	"github.com/thanhfphan/kart-challenge/pkg/xerror"
)

// OpenAPIServer implements the generated ServerInterface
type OpenAPIServer struct {
	productUC usecases.Product
	orderUC   usecases.Order
}

// NewOpenAPIServer creates a new OpenAPI server implementation
func NewOpenAPIServer(ucs usecases.UseCase) *OpenAPIServer {
	return &OpenAPIServer{
		productUC: ucs.Product(),
		orderUC:   ucs.Order(),
	}
}

// Ensure OpenAPIServer implements the ServerInterface
var _ openapi.ServerInterface = (*OpenAPIServer)(nil)

// ListProducts implements GET /product
func (s *OpenAPIServer) ListProducts(c *gin.Context) {
	ctx := c.Request.Context()
	log := logging.FromContext(ctx)

	log.Info("Listing all products")

	products, err := s.productUC.List(ctx)
	if err != nil {
		log.Errorf("Failed to list products: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list products",
		})
		return
	}

	// Convert to OpenAPI format
	response := make([]openapi.Product, 0, len(products))
	for _, product := range products {
		response = append(response, convertToOpenAPIProduct(product))
	}

	c.JSON(http.StatusOK, response)
}

// GetProduct implements GET /product/{productId}
func (s *OpenAPIServer) GetProduct(c *gin.Context, productId int64) {
	ctx := c.Request.Context()
	log := logging.FromContext(ctx)

	log.Infof("Getting product with ID: %d", productId)

	product, err := s.productUC.Get(ctx, productId)
	if err != nil {
		log.Errorf("Failed to get product: %v", err)
		if err == xerror.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Product not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get product",
			})
		}
		return
	}

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

	// Convert OpenAPI request to DTO
	dtoReq := convertToOrderRequest(&orderReq)

	// Place order using use case
	order, err := s.orderUC.PlaceOrder(ctx, dtoReq)
	if err != nil {
		log.Errorf("Failed to place order: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Convert DTO response to OpenAPI format
	response := convertToOpenAPIOrder(order)
	c.JSON(http.StatusOK, response)
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

// convertToOrderRequest converts OpenAPI OrderReq to DTO OrderRequest
func convertToOrderRequest(req *openapi.OrderReq) *dto.OrderRequest {
	items := make([]dto.OrderItem, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, dto.OrderItem{
			ProductID: item.ProductId,
			Quantity:  item.Quantity,
		})
	}

	couponCode := ""
	if req.CouponCode != nil {
		couponCode = *req.CouponCode
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
				Thumbnail: getImageURL(product.Image, "thumbnail"),
				Mobile:    getImageURL(product.Image, "mobile"),
				Tablet:    getImageURL(product.Image, "tablet"),
				Desktop:   getImageURL(product.Image, "desktop"),
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

// getImageURL safely gets image URL from ProductImage
func getImageURL(image *dto.ProductImage, imageType string) *string {
	if image == nil {
		return nil
	}

	switch imageType {
	case "thumbnail":
		return stringPtr(image.Thumbnail)
	case "mobile":
		return stringPtr(image.Mobile)
	case "tablet":
		return stringPtr(image.Tablet)
	case "desktop":
		return stringPtr(image.Desktop)
	default:
		return nil
	}
}
