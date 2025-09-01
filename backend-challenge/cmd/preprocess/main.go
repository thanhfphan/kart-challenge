package main

import (
	"flag"
	"log"
	"os"
	"path"
	"sync"

	"github.com/thanhfphan/kart-challenge/cmd/preprocess/downloader"
	"github.com/thanhfphan/kart-challenge/cmd/preprocess/merger"
	"github.com/thanhfphan/kart-challenge/cmd/preprocess/sorter"
)

var flagOutputDir = flag.String("output_dir", "data", "output directory")

const (
	chunkLimit = 1_000_000
)

func main() {
	flag.Parse()

	log.Println("Starting promo code preprocessing pipeline...")

	if err := os.MkdirAll(*flagOutputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

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
			log.Printf("Processing %s...", urls[i])
			if err := downloader.ProcessGzipToPairs(urls[i], pairsFiles[i]); err != nil {
				log.Fatalf("Error processing %s: %v", urls[i], err)
			}
			log.Printf("Finished processing %s", urls[i])
		}(i)
	}
	wg.Wait()

	// Sort each
	for i := range pairsFiles {
		log.Printf("Sorting %s...", pairsFiles[i])
		if err := sorter.ExternalSortPairs(pairsFiles[i], sortedPairsFiles[i], chunkLimit); err != nil {
			log.Fatalf("Error sorting %s: %v", pairsFiles[i], err)
		}
	}

	// Merge
	log.Println("Merging...")
	if err := merger.Merge3PairsToValid(
		sortedPairsFiles[0], sortedPairsFiles[1], sortedPairsFiles[2],
		path.Join(*flagOutputDir, "valid_coupons.txt"),
		"", // no binary output, add binary output path here if needed
	); err != nil {
		log.Fatalf("Error merging: %v", err)
	}

	log.Println("Preprocessing completed successfully!")
}
