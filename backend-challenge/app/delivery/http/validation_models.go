package http

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// OrderRequestValidation represents the validation model for order requests
// This mirrors the OpenAPI OrderReq but with validation tags
type OrderRequestValidation struct {
	CouponCode *string               `json:"couponCode,omitempty" validate:"omitempty,min=3,max=20,alphanum"`
	Items      []OrderItemValidation `json:"items" validate:"required,min=1,dive"`
}

// OrderItemValidation represents the validation model for order items
type OrderItemValidation struct {
	ProductID string `json:"productId" validate:"required,numeric"`
	Quantity  int    `json:"quantity" validate:"required,min=1,max=100"`
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateOrderRequest validates an order request using validator/v10 and returns validation errors
func ValidateOrderRequest(req *OrderRequestValidation) []string {
	var errors []string

	err := validate.Struct(req)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, formatValidationError(err))
		}
	}

	return errors
}

// formatValidationError formats a validation error into a human-readable message
func formatValidationError(err validator.FieldError) string {
	field := strings.ToLower(err.Field())

	switch err.Tag() {
	case "required":
		return fmt.Sprintf("Field '%s' is required", field)
	case "min":
		if err.Kind().String() == "slice" {
			return fmt.Sprintf("Field '%s' must contain at least %s items", field, err.Param())
		}
		return fmt.Sprintf("Field '%s' must be at least %s characters long", field, err.Param())
	case "max":
		if err.Kind().String() == "slice" {
			return fmt.Sprintf("Field '%s' must contain at most %s items", field, err.Param())
		}
		return fmt.Sprintf("Field '%s' must be at most %s characters long", field, err.Param())
	case "numeric":
		return fmt.Sprintf("Field '%s' must be a valid number", field)
	case "alphanum":
		return fmt.Sprintf("Field '%s' must contain only alphanumeric characters", field)
	case "email":
		return fmt.Sprintf("Field '%s' must be a valid email address", field)
	default:
		return fmt.Sprintf("Field '%s' is invalid: %s", field, err.Tag())
	}
}
