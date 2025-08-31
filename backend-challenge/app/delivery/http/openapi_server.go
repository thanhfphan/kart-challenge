package http

//go:generate oapi-codegen --config=../../../oapi-codegen.yaml ../../../openapi.yaml

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhfphan/kart-challenge/app/delivery/http/openapi"
	"github.com/thanhfphan/kart-challenge/app/usecases"
	"github.com/thanhfphan/kart-challenge/config"
	"github.com/thanhfphan/kart-challenge/pkg/logging"
	"github.com/thanhfphan/kart-challenge/pkg/validation"
	"github.com/thanhfphan/kart-challenge/pkg/xerror"
)

// Ensure OpenAPIServer implements the ServerInterface
var _ openapi.ServerInterface = (*OpenAPIServer)(nil)

// OpenAPIServer implements the generated ServerInterface
type OpenAPIServer struct {
	productUC usecases.Product
	orderUC   usecases.Order
	security  *config.Security
}

// NewOpenAPIServer creates a new OpenAPI server implementation
func NewOpenAPIServer(cfg *config.Config, ucs usecases.UseCase) *OpenAPIServer {
	return &OpenAPIServer{
		productUC: ucs.Product(),
		orderUC:   ucs.Order(),
		security:  cfg.Security,
	}
}

// ListProducts implements GET /product
func (s *OpenAPIServer) ListProducts(c *gin.Context) {
	ctx := c.Request.Context()
	log := logging.FromContext(ctx)

	log.Info("Listing all products")

	products, err := s.productUC.List(ctx)
	if err != nil {
		log.Errorf("Failed to list products: %v", err)
		validation.SendInternalServerError(c, "Failed to list products")
		return
	}

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

	if productId <= 0 {
		log.Warnf("Invalid product ID: %d", productId)
		validation.SendBadRequestError(c, "Invalid product ID: must be a positive integer")
		return
	}

	log.Infof("Getting product with ID: %d", productId)

	product, err := s.productUC.Get(ctx, productId)
	if err != nil {
		log.Errorf("Failed to get product: %v", err)
		if err == xerror.ErrRecordNotFound {
			validation.SendNotFoundError(c, "Product not found")
		} else {
			validation.SendInternalServerError(c, "Failed to get product")
		}
		return
	}

	c.JSON(http.StatusOK, convertToOpenAPIProduct(product))
}

// PlaceOrder implements POST /order
func (s *OpenAPIServer) PlaceOrder(c *gin.Context) {
	ctx := c.Request.Context()
	log := logging.FromContext(ctx)

	if !validation.ValidateAPIKey(c, s.security.APIKey) {
		log.Warn("Invalid or missing API key for order placement")
		validation.SendUnauthorizedError(c, "Invalid or missing API key")
		return
	}

	var validationReq OrderRequestValidation
	if err := c.ShouldBindJSON(&validationReq); err != nil {
		log.Warnf("Failed to bind order request: %v", err)
		validation.SendBadRequestError(c, fmt.Sprintf("Invalid JSON: %s", err.Error()))
		return
	}

	// Validate request structure and business rules
	if validationErrors := ValidateOrderRequest(&validationReq); len(validationErrors) > 0 {
		log.Warnf("Order request validation failed: %v", validationErrors)
		validation.SendValidationError(c, validationErrors)
		return
	}

	dtoReq := convertValidationToDTO(&validationReq)
	order, err := s.orderUC.PlaceOrder(ctx, dtoReq)
	if err != nil {
		log.Errorf("Failed to place order: %v", err)
		validation.SendBadRequestError(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, convertToOpenAPIOrder(order))
}
