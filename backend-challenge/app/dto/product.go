package dto

// ProductResponse represents the API response for a product
// Maps to OpenAPI Product schema
type ProductResponse struct {
	ID       string        `json:"id"`
	Name     string        `json:"name"`
	Price    float64       `json:"price"`
	Category string        `json:"category"`
	Image    *ProductImage `json:"image,omitempty"`
}

// ProductImage represents the image URLs for different screen sizes
type ProductImage struct {
	Thumbnail string `json:"thumbnail,omitempty"`
	Mobile    string `json:"mobile,omitempty"`
	Tablet    string `json:"tablet,omitempty"`
	Desktop   string `json:"desktop,omitempty"`
}

// ProductListResponse represents the API response for listing products
type ProductListResponse []ProductResponse
