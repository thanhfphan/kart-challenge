package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/thanhfphan/kart-challenge/pkg/infras/migrate"
	"github.com/thanhfphan/kart-challenge/pkg/logging"
	"github.com/thanhfphan/kart-challenge/setup"
)

var (
	pathMigration = flag.String("path", "migrations/", "path to migrations folder")
)

func main() {
	flag.Parse()
	ctx, done := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	log := logging.FromContext(ctx)
	ctx = logging.WithLogger(ctx, log)

	defer func() {
		done()
		if r := recover(); r != nil {
			log.Errorf("apllication went wrong. Panic err=%v", r)
		}
	}()

	err := realMain(ctx)
	done()
	if err != nil {
		log.Errorf("realMain has failed with err=%v", err)
		return
	}

	log.Infof("APP shutdown successful")
}

func realMain(ctx context.Context) error {
	log := logging.FromContext(ctx)
	log.Infof("starting migration ...")

	cfg, env, err := setup.LoadFromEnv(ctx)
	if err != nil {
		return fmt.Errorf("load config from environment failed with err=%w", err)
	}
	defer env.Close(ctx)

	dir := fmt.Sprintf("file://%s", *pathMigration)
	tool := migrate.New()
	err = tool.Migrate(dir, cfg.DB.MigrationURL)
	if err != nil {
		return err
	}

	log.Infof("Migration done ...")

	return nil
}
