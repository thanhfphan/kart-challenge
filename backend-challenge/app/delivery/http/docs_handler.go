package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhfphan/kart-challenge/app/delivery/http/openapi"
)

// DocsHandler handles API documentation endpoints
type DocsHandler struct{}

// NewDocsHandler creates a new documentation handler
func NewDocsHandler() *DocsHandler {
	return &DocsHandler{}
}

// ServeOpenAPISpec serves the raw OpenAPI specification as JSON
func (h *DocsHandler) ServeOpenAPISpec(c *gin.Context) {
	spec, err := openapi.GetSwagger()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to load OpenAPI specification",
		})
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, spec)
}

// ServeDocumentation serves the HTML documentation page
func (h *DocsHandler) ServeDocumentation(c *gin.Context) {
	// Get the current request scheme and host for the spec URL
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	
	specURL := scheme + "://" + c.Request.Host + "/openapi.json"
	
	html := h.generateScalarHTML(specURL)
	
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, html)
}

// generateScalarHTML generates the HTML page using Scalar API Reference
func (h *DocsHandler) generateScalarHTML(specURL string) string {
	return `<!DOCTYPE html>
<html>
<head>
    <title>API Documentation - Order Food Online</title>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <style>
        body {
            margin: 0;
            padding: 0;
        }
    </style>
</head>
<body>
    <script
        id="api-reference"
        data-url="` + specURL + `"
        data-configuration='{
            "theme": "purple",
            "layout": "modern",
            "defaultHttpClient": {
                "targetKey": "javascript",
                "clientKey": "fetch"
            },
            "spec": {
                "url": "` + specURL + `"
            },
            "customCss": ".scalar-app { font-family: -apple-system, BlinkMacSystemFont, \"Segoe UI\", Roboto, Oxygen, Ubuntu, Cantarell, \"Open Sans\", \"Helvetica Neue\", sans-serif; }"
        }'></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
</body>
</html>`
}
