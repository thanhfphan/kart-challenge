package sorter

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"os"

	"github.com/lanrat/extsort"
	"github.com/thanhfphan/kart-challenge/cmd/preprocess/types"
	"github.com/thanhfphan/kart-challenge/cmd/preprocess/utils"
)

// compareRec implements the CompareGeneric interface for external sorting
func compareRec(a, b types.Rec) int {
	if a.H != b.H {
		if a.H < b.H {
			return -1
		}
		return 1
	}
	if a.Code != b.Code {
		if a.Code < b.Code {
			return -1
		}
		return 1
	}
	return 0
}

// recToBytes serializes a types.Rec to bytes for external sorting
func recToBytes(rec types.Rec) ([]byte, error) {
	// Calculate total size: 8 bytes for hash + 2 bytes for length + string length
	totalSize := 8 + 2 + len(rec.Code)
	buf := make([]byte, totalSize)

	// Write hash (8 bytes)
	binary.LittleEndian.PutUint64(buf[0:8], rec.H)

	// Write string length (2 bytes)
	binary.LittleEndian.PutUint16(buf[8:10], uint16(len(rec.Code)))

	// Write string data
	copy(buf[10:], []byte(rec.Code))

	return buf, nil
}

// recFromBytes deserializes bytes back to types.Rec for external sorting
func recFromBytes(data []byte) (types.Rec, error) {
	if len(data) < 10 {
		return types.Rec{}, fmt.Errorf("invalid data length: %d, expected at least 10 bytes", len(data))
	}

	// Read hash (8 bytes)
	h := binary.LittleEndian.Uint64(data[0:8])

	// Read string length (2 bytes)
	strLen := binary.LittleEndian.Uint16(data[8:10])

	// Validate remaining data length
	if len(data) != 10+int(strLen) {
		return types.Rec{}, fmt.Errorf("invalid data length: %d, expected %d", len(data), 10+int(strLen))
	}

	// Read string data
	code := string(data[10:])

	return types.Rec{H: h, Code: code}, nil
}

// ExternalSortPairs sorts and deduplicates a pairs file using external sorting.
func ExternalSortPairs(ctx context.Context, inputPath, outputPath string, chunkLimit int) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	inputChan := make(chan types.Rec, 100)

	go func() {
		defer close(inputChan)
		br := bufio.NewReader(inputFile)
		for {
			rec, ok, err := utils.ReadPair(br)
			if err != nil {
				fmt.Println("Error reading pair: ", err)
				return
			}
			if !ok {
				break
			}
			inputChan <- rec
		}
	}()

	config := &extsort.Config{
		NumWorkers:         2,
		ChanBuffSize:       10,
		SortedChanBuffSize: 1000,
	}

	sorter, outputChan, errChan := extsort.Generic(
		inputChan,
		recFromBytes,
		recToBytes,
		compareRec,
		config,
	)

	go sorter.Sort(ctx)

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	bw := bufio.NewWriter(outputFile)
	defer bw.Flush()

	var lastRec types.Rec
	var hasLast bool

	for rec := range outputChan {
		// Deduplicate: only write if different from last record
		if !hasLast || rec.H != lastRec.H || rec.Code != lastRec.Code {
			if err := utils.WriteRec(bw, rec); err != nil {
				return fmt.Errorf("failed to write record: %w", err)
			}
			lastRec = rec
			hasLast = true
		}
	}

	if err := <-errChan; err != nil {
		return fmt.Errorf("external sort failed: %w", err)
	}

	return nil
}
