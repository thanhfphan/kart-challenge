package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thanhfphan/kart-challenge/pkg/validation"
)

func TestValidateOrderRequest(t *testing.T) {
	tests := []struct {
		name           string
		request        OrderRequestValidation
		expectedErrors int
		errorContains  []string
	}{
		{
			name: "valid request",
			request: OrderRequestValidation{
				Items: []OrderItemValidation{
					{ProductID: "1", Quantity: 2},
					{ProductID: "2", Quantity: 1},
				},
			},
			expectedErrors: 0,
		},
		{
			name: "valid request with coupon",
			request: OrderRequestValidation{
				CouponCode: stringPtr("HAPPYHRS"),
				Items: []OrderItemValidation{
					{ProductID: "1", Quantity: 1},
				},
			},
			expectedErrors: 0,
		},
		{
			name: "empty items array",
			request: OrderRequestValidation{
				Items: []OrderItemValidation{},
			},
			expectedErrors: 1,
			errorContains:  []string{"items", "1 items"},
		},
		{
			name: "invalid product ID",
			request: OrderRequestValidation{
				Items: []OrderItemValidation{
					{ProductID: "invalid", Quantity: 1},
				},
			},
			expectedErrors: 1,
			errorContains:  []string{"productid", "valid number"},
		},
		{
			name: "zero quantity",
			request: OrderRequestValidation{
				Items: []OrderItemValidation{
					{ProductID: "1", Quantity: 0},
				},
			},
			expectedErrors: 1,
			errorContains:  []string{"quantity", "required"},
		},
		{
			name: "quantity too high",
			request: OrderRequestValidation{
				Items: []OrderItemValidation{
					{ProductID: "1", Quantity: 101},
				},
			},
			expectedErrors: 1,
			errorContains:  []string{"quantity", "most 100"},
		},
		{
			name: "coupon too short",
			request: OrderRequestValidation{
				CouponCode: stringPtr("AB"),
				Items: []OrderItemValidation{
					{ProductID: "1", Quantity: 1},
				},
			},
			expectedErrors: 1,
			errorContains:  []string{"couponcode", "3 characters"},
		},
		{
			name: "coupon too long",
			request: OrderRequestValidation{
				CouponCode: stringPtr("VERYLONGCOUPONCODETHATEXCEEDSLIMIT"),
				Items: []OrderItemValidation{
					{ProductID: "1", Quantity: 1},
				},
			},
			expectedErrors: 1,
			errorContains:  []string{"couponcode", "20 characters"},
		},
		{
			name: "multiple validation errors",
			request: OrderRequestValidation{
				CouponCode: stringPtr("AB"),
				Items: []OrderItemValidation{
					{ProductID: "", Quantity: 0},
					{ProductID: "invalid", Quantity: 101},
				},
			},
			expectedErrors: 5, // coupon + 2 items * 2 errors each
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := ValidateOrderRequest(&tt.request)

			assert.Equal(t, tt.expectedErrors, len(errors), "Expected %d errors, got %d: %v", tt.expectedErrors, len(errors), errors)

			for _, expectedError := range tt.errorContains {
				found := false
				for _, actualError := range errors {
					if contains(actualError, expectedError) {
						found = true
						break
					}
				}
				assert.True(t, found, "Expected error containing '%s' not found in: %v", expectedError, errors)
			}
		})
	}
}

func TestAPIKeyValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		headers     map[string]string
		expectValid bool
	}{
		{
			name:        "valid API key",
			headers:     map[string]string{"api_key": "apitest"},
			expectValid: true,
		},
		{
			name:        "valid API key with different case header",
			headers:     map[string]string{"Api-Key": "apitest"},
			expectValid: true,
		},
		{
			name:        "valid API key with uppercase header",
			headers:     map[string]string{"API-Key": "apitest"},
			expectValid: true,
		},
		{
			name:        "invalid API key",
			headers:     map[string]string{"api_key": "invalid"},
			expectValid: false,
		},
		{
			name:        "missing API key",
			headers:     map[string]string{},
			expectValid: false,
		},
		{
			name:        "empty API key",
			headers:     map[string]string{"api_key": ""},
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			req := httptest.NewRequest("GET", "/", nil)
			for key, value := range tt.headers {
				req.Header.Set(key, value)
			}
			c.Request = req

			result := validation.ValidateAPIKey(c, "apitest")
			assert.Equal(t, tt.expectValid, result)
		})
	}
}

func TestErrorResponseFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		sendError      func(*gin.Context)
		expectedStatus int
		expectedType   string
	}{
		{
			name: "bad request error",
			sendError: func(c *gin.Context) {
				validation.SendBadRequestError(c, "Invalid input")
			},
			expectedStatus: http.StatusBadRequest,
			expectedType:   "bad_request",
		},
		{
			name: "validation error",
			sendError: func(c *gin.Context) {
				validation.SendValidationError(c, []string{"Field 'name' is required"})
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedType:   "validation_error",
		},
		{
			name: "unauthorized error",
			sendError: func(c *gin.Context) {
				validation.SendUnauthorizedError(c, "Invalid API key")
			},
			expectedStatus: http.StatusUnauthorized,
			expectedType:   "unauthorized",
		},
		{
			name: "not found error",
			sendError: func(c *gin.Context) {
				validation.SendNotFoundError(c, "Product not found")
			},
			expectedStatus: http.StatusNotFound,
			expectedType:   "not_found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			tt.sendError(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response validation.ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Equal(t, int32(tt.expectedStatus), response.Code)
			assert.Equal(t, tt.expectedType, response.Type)
			assert.NotEmpty(t, response.Message)
		})
	}
}

func TestPlaceOrderValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		headers        map[string]string
		body           interface{}
		expectedStatus int
		expectedType   string
	}{
		{
			name:    "missing API key",
			headers: map[string]string{},
			body: OrderRequestValidation{
				Items: []OrderItemValidation{{ProductID: "1", Quantity: 1}},
			},
			expectedStatus: http.StatusUnauthorized,
			expectedType:   "unauthorized",
		},
		{
			name:           "invalid JSON",
			headers:        map[string]string{"api_key": "apitest"},
			body:           "invalid json",
			expectedStatus: http.StatusBadRequest,
			expectedType:   "bad_request",
		},
		{
			name:    "validation error",
			headers: map[string]string{"api_key": "apitest"},
			body: OrderRequestValidation{
				Items: []OrderItemValidation{{ProductID: "", Quantity: 0}},
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedType:   "validation_error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			var bodyBytes []byte
			if str, ok := tt.body.(string); ok {
				bodyBytes = []byte(str)
			} else {
				var err error
				bodyBytes, err = json.Marshal(tt.body)
				require.NoError(t, err)
			}

			req := httptest.NewRequest("POST", "/order", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			for key, value := range tt.headers {
				req.Header.Set(key, value)
			}
			c.Request = req

			if !validation.ValidateAPIKey(c, "apitest") {
				validation.SendUnauthorizedError(c, "Invalid or missing API key")
				return
			}

			var validationReq OrderRequestValidation
			if err := c.ShouldBindJSON(&validationReq); err != nil {
				validation.SendBadRequestError(c, "Invalid JSON")
				return
			}

			if validationErrors := ValidateOrderRequest(&validationReq); len(validationErrors) > 0 {
				validation.SendValidationError(c, validationErrors)
				return
			}

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus != http.StatusOK {
				var response validation.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedType, response.Type)
			}
		})
	}
}

// Helper functions
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			bytes.Contains([]byte(s), []byte(substr)))))
}
