package downloader

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/thanhfphan/kart-challenge/cmd/preprocess/utils"
)

const (
	MinCodeLength = 8
	MaxCodeLength = 10
)

// ProcessGzipToPairs downloads and processes a gzip file to a pairs binary file.
func ProcessGzipToPairs(url, outputPath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("download: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	gr, err := gzip.NewReader(resp.Body)
	if err != nil {
		return fmt.Errorf("gzip: %w", err)
	}
	defer gr.Close()

	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()
	bw := bufio.NewWriter(f)
	defer bw.Flush()

	sc := bufio.NewScanner(gr)
	sc.Buffer(make([]byte, 64<<10), 1<<20)
	var n, invalid int
	for sc.Scan() {
		s := strings.TrimSpace(sc.Text())
		l := len(s)
		if l < MinCodeLength || l > MaxCodeLength {
			invalid++
			continue
		}
		code := strings.ToUpper(s)
		if err := utils.WritePair(bw, code); err != nil {
			return err
		}
		n++
	}

	if err := sc.Err(); err != nil {
		return err
	}

	log.Printf("Processed %d pairs, skipped %d invalid -> %s", n, invalid, outputPath)

	return nil
}
