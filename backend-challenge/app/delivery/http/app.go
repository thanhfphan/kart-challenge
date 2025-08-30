package http

import (
	"context"
	"net/http"

	"github.com/thanhfphan/kart-challenge/app/delivery/http/openapi"
	"github.com/thanhfphan/kart-challenge/app/usecases"
	"github.com/thanhfphan/kart-challenge/config"
	"github.com/thanhfphan/kart-challenge/pkg/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

var (
	_ App = (*app)(nil)
)

type App interface {
	Routes(ctx context.Context) http.Handler
}

type app struct {
	cfg *config.Config

	ucs usecases.UseCase
}

func New(cfg *config.Config, ucs usecases.UseCase) (App, error) {
	return &app{
		cfg: cfg,
		ucs: ucs,
	}, nil
}

func (a *app) Routes(ctx context.Context) http.Handler {
	if a.cfg.Environment == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()

	// middlewares
	r.Use(gin.Recovery())
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(middleware.SetRequestID())
	r.Use(middleware.SetLogger())
	r.Use(middleware.GinRequestProfiler(a.cfg.ServiceName))
	r.Use(middleware.SetupAppContext())
	r.Use(middleware.SetCommonData())

	// cors
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{
		"*",
		"Origin",
		"Content-Length",
		"Content-Type",
		"Authorization",
	}
	r.Use(cors.New(corsConfig))

	// health check
	pingHandler := func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"version": a.cfg.ServiceVersion,
			},
		})
	}

	r.GET("/health-check", pingHandler)
	r.HEAD("/health-check", pingHandler)

	openAPIServer := NewOpenAPIServer(a.ucs)
	apiGroup := r.Group("api")
	openapi.RegisterHandlers(apiGroup, openAPIServer)

	return r
}
