package app

import (
	"context"
	"fmt"

	appHttp "github.com/thanhfphan/kart-challenge/app/delivery/http"
	appMetric "github.com/thanhfphan/kart-challenge/app/delivery/metrics"
	"github.com/thanhfphan/kart-challenge/app/repos"
	"github.com/thanhfphan/kart-challenge/app/usecases"
	"github.com/thanhfphan/kart-challenge/config"
	"github.com/thanhfphan/kart-challenge/env"
	"github.com/thanhfphan/kart-challenge/pkg/logging"
	"github.com/thanhfphan/kart-challenge/pkg/server"
)

type App interface {
	Start(ctx context.Context) error
}

type app struct {
	cfg *config.Config
	env *env.Env
	ucs usecases.UseCase
}

func New(cfg *config.Config, env *env.Env) (App, error) {
	repos := repos.New(cfg, env, env.Database())

	ucs, err := usecases.New(cfg, env, repos)
	if err != nil {
		return nil, fmt.Errorf("new usecases failed err=%w", err)
	}

	return &app{
		cfg: cfg,
		env: env,
		ucs: ucs,
	}, nil
}

func (a *app) Start(ctx context.Context) error {
	log := logging.FromContext(ctx)

	// HTTP
	app, err := appHttp.New(a.cfg, a.ucs)
	if err != nil {
		return fmt.Errorf("new application failed err=%w", err)
	}

	srv, err := server.New(a.cfg.HTTPPort)
	if err != nil {
		return err
	}

	// Metric
	go func() {
		srvMetric, err := server.New(a.cfg.MetricPort)
		if err != nil {
			log.Warnf("New server metric err=%v", err)
			return
		}

		log.Infof("Metric running on PORT: %s", srvMetric.Port())

		metric, err := appMetric.New()
		if err != nil {
			log.Warnf("New app metric err=%v", err)
			return
		}
		if err := srvMetric.ServeHTTPHandler(ctx, metric.Handler()); err != nil {
			log.Warnf("Serve metric handler err=%v", err)
		}
	}()

	log.Infof("HTTP Server running on PORT: %s", srv.Port())

	// This is block the main goroutine
	return srv.ServeHTTPHandler(ctx, app.Routes(ctx))
}
