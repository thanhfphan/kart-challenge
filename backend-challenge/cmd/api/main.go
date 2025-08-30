package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/thanhfphan/kart-challenge/app"
	"github.com/thanhfphan/kart-challenge/pkg/logging"
	"github.com/thanhfphan/kart-challenge/setup"
)

func main() {
	ctx, done := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	log := logging.FromContext(ctx)
	ctx = logging.WithLogger(ctx, log)

	defer func() {
		done()
		if r := recover(); r != nil {
			log.Error("Application went wrong. Got panic", slog.Any("error", r))
		}
	}()

	err := realMain(ctx)
	done()
	if err != nil {
		log.Errorf("Start application failed with err=%v", err)
		return
	}

	log.Infof("Application shutdown successful")
}

func realMain(ctx context.Context) error {
	log := logging.FromContext(ctx)
	log.Info("Starting application ...")

	cfg, env, err := setup.LoadFromEnv(ctx)
	if err != nil {
		return fmt.Errorf("Load config failed with err=%v", err)
	}
	defer env.Close(ctx)

	app, err := app.New(cfg, env)
	if err != nil {
		return fmt.Errorf("New application failed with err=%v", err)
	}

	return app.Start(ctx)
}
