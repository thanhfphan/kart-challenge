package sorter

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thanhfphan/kart-challenge/cmd/preprocess/types"
	"github.com/thanhfphan/kart-challenge/cmd/preprocess/utils"
)

func TestExternalSortPairs(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "sorter_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test input file
	inputPath := filepath.Join(tempDir, "test_input.bin")
	outputPath := filepath.Join(tempDir, "test_output.bin")

	// Create test data - unsorted coupon codes
	testCodes := []string{
		"ZEBRA123",
		"APPLE456",
		"BANANA789",
		"CHERRY012",
		"APPLE456",  // duplicate
		"BANANA789", // duplicate
		"DELTA345",
	}

	// Write test data to input file
	inputFile, err := os.Create(inputPath)
	require.NoError(t, err)
	bw := bufio.NewWriter(inputFile)
	for _, code := range testCodes {
		err := utils.WritePair(bw, code)
		require.NoError(t, err)
	}
	err = bw.Flush()
	require.NoError(t, err)
	err = inputFile.Close()
	require.NoError(t, err)

	// Test the ExternalSortPairs function
	err = ExternalSortPairs(context.Background(), inputPath, outputPath, 1000)
	require.NoError(t, err)

	// Read and verify the output
	outputFile, err := os.Open(outputPath)
	require.NoError(t, err)
	defer outputFile.Close()

	br := bufio.NewReader(outputFile)
	var sortedRecords []types.Rec
	for {
		rec, ok, err := utils.ReadPair(br)
		require.NoError(t, err)
		if !ok {
			break
		}
		sortedRecords = append(sortedRecords, rec)
	}

	// Verify results
	assert.True(t, len(sortedRecords) > 0, "Should have sorted records")
	assert.True(t, len(sortedRecords) <= len(testCodes), "Should not have more records than input")

	// Verify sorting order
	for i := 1; i < len(sortedRecords); i++ {
		prev := sortedRecords[i-1]
		curr := sortedRecords[i]

		// Check that records are sorted by hash first, then by code
		if prev.H == curr.H {
			assert.True(t, prev.Code <= curr.Code, "Records with same hash should be sorted by code")
		} else {
			assert.True(t, prev.H < curr.H, "Records should be sorted by hash")
		}
	}

	// Verify deduplication - no consecutive duplicates
	for i := 1; i < len(sortedRecords); i++ {
		prev := sortedRecords[i-1]
		curr := sortedRecords[i]
		assert.False(t, prev.H == curr.H && prev.Code == curr.Code, "Should not have duplicate records")
	}
}

func TestExternalSortPairsWithSmallChunkSize(t *testing.T) {
	// Test with very small chunk size to ensure external sorting works
	tempDir, err := os.MkdirTemp("", "sorter_test_small")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	inputPath := filepath.Join(tempDir, "test_input_small.bin")
	outputPath := filepath.Join(tempDir, "test_output_small.bin")

	// Create more test data to force multiple chunks
	testCodes := make([]string, 100)
	for i := 0; i < 100; i++ {
		testCodes[i] = fmt.Sprintf("CODE%04d", 99-i) // Reverse order to test sorting
	}

	// Write test data
	inputFile, err := os.Create(inputPath)
	require.NoError(t, err)
	bw := bufio.NewWriter(inputFile)
	for _, code := range testCodes {
		err := utils.WritePair(bw, code)
		require.NoError(t, err)
	}
	err = bw.Flush()
	require.NoError(t, err)
	err = inputFile.Close()
	require.NoError(t, err)

	// Test with small chunk size (10 records per chunk)
	err = ExternalSortPairs(context.Background(), inputPath, outputPath, 10)
	require.NoError(t, err)

	// Verify output exists and is not empty
	info, err := os.Stat(outputPath)
	require.NoError(t, err)
	assert.True(t, info.Size() > 0, "Output file should not be empty")
}

func TestExternalSortPairsEmptyFile(t *testing.T) {
	// Test with empty input file
	tempDir, err := os.MkdirTemp("", "sorter_test_empty")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	inputPath := filepath.Join(tempDir, "empty_input.bin")
	outputPath := filepath.Join(tempDir, "empty_output.bin")

	// Create empty input file
	inputFile, err := os.Create(inputPath)
	require.NoError(t, err)
	err = inputFile.Close()
	require.NoError(t, err)

	// Test sorting empty file
	err = ExternalSortPairs(context.Background(), inputPath, outputPath, 1000)
	require.NoError(t, err)

	// Verify output file exists and is empty
	info, err := os.Stat(outputPath)
	require.NoError(t, err)
	assert.Equal(t, int64(0), info.Size(), "Output file should be empty for empty input")
}
