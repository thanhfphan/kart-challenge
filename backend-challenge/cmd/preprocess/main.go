package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"path"
	"sync"
	"syscall"

	"github.com/thanhfphan/kart-challenge/app/repos"
	"github.com/thanhfphan/kart-challenge/app/usecases"
	"github.com/thanhfphan/kart-challenge/cmd/preprocess/downloader"
	"github.com/thanhfphan/kart-challenge/cmd/preprocess/merger"
	"github.com/thanhfphan/kart-challenge/cmd/preprocess/sorter"
	"github.com/thanhfphan/kart-challenge/pkg/logging"
	"github.com/thanhfphan/kart-challenge/setup"
)

var flagOutputDir = flag.String("output_dir", "data", "output directory")

var chunkLimit = 1_000_000

func main() {
	flag.Parse()

	ctx, done := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer done()

	log := logging.FromContext(ctx)
	ctx = logging.WithLogger(ctx, log)

	log.Info("Starting promo code preprocessing pipeline...")

	if err := os.MkdirAll(*flagOutputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
		return
	}

	if err := realMain(ctx); err != nil {
		log.Fatalf("Preprocessing failed: %v", err)
		return
	}

	log.Info("Preprocessing completed successfully!")
}

func realMain(ctx context.Context) error {
	log := logging.FromContext(ctx)

	urls := []string{
		"https://orderfoodonline-files.s3.ap-southeast-2.amazonaws.com/couponbase1.gz",
		"https://orderfoodonline-files.s3.ap-southeast-2.amazonaws.com/couponbase2.gz",
		"https://orderfoodonline-files.s3.ap-southeast-2.amazonaws.com/couponbase3.gz",
	}
	pairsFiles := []string{
		path.Join(*flagOutputDir, "pairs1.bin"),
		path.Join(*flagOutputDir, "pairs2.bin"),
		path.Join(*flagOutputDir, "pairs3.bin"),
	}
	sortedPairsFiles := []string{
		path.Join(*flagOutputDir, "pairs1.sorted.bin"),
		path.Join(*flagOutputDir, "pairs2.sorted.bin"),
		path.Join(*flagOutputDir, "pairs3.sorted.bin"),
	}

	// Parallel processing for downloads
	var wg sync.WaitGroup
	for i := range urls {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			log.Infof("Processing %s...", urls[i])
			if err := downloader.ProcessGzipToPairs(urls[i], pairsFiles[i]); err != nil {
				log.Errorf("Error processing %s: %v", urls[i], err)
				return
			}
			log.Infof("Finished processing %s", urls[i])
		}(i)
	}
	wg.Wait()

	// Sort each
	for i := range pairsFiles {
		log.Infof("Sorting %s...", pairsFiles[i])
		if err := sorter.ExternalSortPairs(ctx, pairsFiles[i], sortedPairsFiles[i], chunkLimit); err != nil {
			return err
		}
	}

	// Merge
	log.Info("Merging...")
	validCouponsFile := path.Join(*flagOutputDir, "valid_coupons.txt")
	if err := merger.Merge3PairsToValid(
		sortedPairsFiles[0], sortedPairsFiles[1], sortedPairsFiles[2],
		validCouponsFile,
		"", // no binary output, add binary output path here if needed
	); err != nil {
		return err
	}

	log.Info("Data preprocessing completed successfully!")

	// Now process the coupon data and store in database
	return processCouponData(ctx, validCouponsFile)
}

func processCouponData(ctx context.Context, validCouponsFile string) error {
	log := logging.FromContext(ctx)
	log.Info("Starting coupon data processing...")

	cfg, env, err := setup.LoadFromEnv(ctx)
	if err != nil {
		return err
	}
	defer env.Close(ctx)

	repos := repos.New(cfg, env, env.Database())
	ucs, err := usecases.New(cfg, env, repos)
	if err != nil {
		return err
	}

	if err := ucs.PromoCode().ProcessCouponFile(ctx, validCouponsFile); err != nil {
		return err
	}

	log.Info("Coupon data processing completed successfully!")
	return nil
}
