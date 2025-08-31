package validation

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Validator wraps the go-playground validator
type Validator struct {
	validator *validator.Validate
}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	v := validator.New()

	// Register custom validation tags
	v.RegisterValidation("productid", validateProductID)

	return &Validator{
		validator: v,
	}
}

// ValidateStruct validates a struct and returns formatted error messages
func (v *Validator) ValidateStruct(s interface{}) []string {
	var errors []string

	err := v.validator.Struct(s)
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
		return fmt.Sprintf("Field '%s' must be at least %s", field, err.Param())
	case "max":
		return fmt.Sprintf("Field '%s' must be at most %s", field, err.Param())
	case "email":
		return fmt.Sprintf("Field '%s' must be a valid email address", field)
	case "productid":
		return fmt.Sprintf("Field '%s' must be a valid product ID", field)
	default:
		return fmt.Sprintf("Field '%s' is invalid", field)
	}
}

// validateProductID validates that a product ID is a valid string representation of an integer
func validateProductID(fl validator.FieldLevel) bool {
	productID := fl.Field().String()
	if productID == "" {
		return false
	}

	_, err := strconv.ParseInt(productID, 10, 64)
	return err == nil
}

// ValidateAPIKey validates the API key from request headers
func ValidateAPIKey(c *gin.Context, key string) bool {
	apiKey := c.GetHeader("api_key")
	if apiKey == "" {
		apiKey = c.GetHeader("Api-Key")
	}
	if apiKey == "" {
		apiKey = c.GetHeader("API-Key")
	}

	return apiKey == key
}

// ValidatePathParameter validates and converts path parameters
func ValidatePathParameter(paramName, paramValue string, paramType string) (interface{}, error) {
	switch paramType {
	case "int64":
		value, err := strconv.ParseInt(paramValue, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parameter '%s' must be a valid integer", paramName)
		}
		if value <= 0 {
			return nil, fmt.Errorf("parameter '%s' must be a positive integer", paramName)
		}
		return value, nil
	case "string":
		if strings.TrimSpace(paramValue) == "" {
			return nil, fmt.Errorf("parameter '%s' cannot be empty", paramName)
		}
		return paramValue, nil
	default:
		return paramValue, nil
	}
}

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Code    int32  `json:"code"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

// SendErrorResponse sends a standardized error response
func SendErrorResponse(c *gin.Context, statusCode int, errorType, message string) {
	response := ErrorResponse{
		Code:    int32(statusCode),
		Type:    errorType,
		Message: message,
	}
	c.JSON(statusCode, response)
	c.Abort()
}

// SendValidationError sends a 422 validation error response
func SendValidationError(c *gin.Context, errors []string) {
	message := "Validation failed"
	if len(errors) > 0 {
		message = strings.Join(errors, "; ")
	}
	SendErrorResponse(c, http.StatusUnprocessableEntity, "validation_error", message)
}

// SendBadRequestError sends a 400 bad request error response
func SendBadRequestError(c *gin.Context, message string) {
	SendErrorResponse(c, http.StatusBadRequest, "bad_request", message)
}

// SendNotFoundError sends a 404 not found error response
func SendNotFoundError(c *gin.Context, message string) {
	SendErrorResponse(c, http.StatusNotFound, "not_found", message)
}

// SendUnauthorizedError sends a 401 unauthorized error response
func SendUnauthorizedError(c *gin.Context, message string) {
	SendErrorResponse(c, http.StatusUnauthorized, "unauthorized", message)
}

// SendForbiddenError sends a 403 forbidden error response
func SendForbiddenError(c *gin.Context, message string) {
	SendErrorResponse(c, http.StatusForbidden, "forbidden", message)
}

// SendInternalServerError sends a 500 internal server error response
func SendInternalServerError(c *gin.Context, message string) {
	SendErrorResponse(c, http.StatusInternalServerError, "internal_error", message)
}

// ValidateJSON validates JSON request body and binds it to the target struct
func ValidateJSON(c *gin.Context, target interface{}) bool {
	if err := c.ShouldBindJSON(target); err != nil {
		SendBadRequestError(c, fmt.Sprintf("Invalid JSON: %s", err.Error()))
		return false
	}

	// Use reflection to check if the target has validation tags
	v := NewValidator()
	if errors := v.ValidateStruct(target); len(errors) > 0 {
		SendValidationError(c, errors)
		return false
	}

	return true
}

// // RequireAPIKey middleware that validates API key for protected endpoints
// func RequireAPIKey() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		if !ValidateAPIKey(c) {
// 			SendUnauthorizedError(c, "Invalid or missing API key")
// 			return
// 		}
// 		c.Next()
// 	}
// }
