package sorter

import (
	"bufio"
	"fmt"
	"os"
	"sort"

	"github.com/thanhfphan/kart-challenge/cmd/preprocess/types"
	"github.com/thanhfphan/kart-challenge/cmd/preprocess/utils"
)

// RecLess compares two records.
func RecLess(a, b types.Rec) bool {
	if a.H != b.H {
		return a.H < b.H
	}
	return a.Code < b.Code
}

// ExternalSortPairs sorts and deduplicates a pairs file.
func ExternalSortPairs(inputPath, outputPath string, chunkLimit int) error {
	runs, err := createSortedRuns(inputPath, chunkLimit)
	if err != nil {
		return err
	}
	return mergeRunsByHashCode(runs, outputPath)
}

// createSortedRuns creates sorted runs from the input file.
// chunkLimit is the maximum number of records in a chunk.
func createSortedRuns(input string, chunkLimit int) ([]string, error) {
	f, err := os.Open(input)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	br := bufio.NewReader(f)
	var runs []string
	for {
		buf := make([]types.Rec, 0, chunkLimit)
		for len(buf) < chunkLimit {
			r, ok, err := utils.ReadPair(br)
			if err != nil {
				return nil, err
			}
			if !ok {
				break
			}
			buf = append(buf, r)
		}
		if len(buf) == 0 {
			break
		}
		sort.Slice(buf, func(i, j int) bool { return RecLess(buf[i], buf[j]) })
		// Dedupe in chunk
		w := 0
		for i := range buf {
			if i == 0 || buf[i].H != buf[i-1].H || buf[i].Code != buf[i-1].Code {
				buf[w] = buf[i]
				w++
			}
		}
		buf = buf[:w]
		run := fmt.Sprintf("%s.run.%06d", input, len(runs))
		rf, err := os.Create(run)
		if err != nil {
			return nil, err
		}
		bw := bufio.NewWriter(rf)
		for _, r := range buf {
			if err := utils.WriteRec(bw, r); err != nil {
				rf.Close()
				return nil, err
			}
		}
		if err := bw.Flush(); err != nil {
			rf.Close()
			return nil, err
		}
		if err := rf.Close(); err != nil {
			return nil, err
		}
		runs = append(runs, run)
	}
	return runs, nil
}

func mergeRunsByHashCode(runs []string, outFile string) error {
	if len(runs) == 0 {
		// no data
		f, _ := os.Create(outFile)
		f.Close()
		return nil
	}

	// open all runs
	readers := make([]*utils.Reader, 0, len(runs))
	files := make([]*os.File, 0, len(runs))
	for _, rn := range runs {
		rr, f, err := utils.NewReader(rn)
		if err != nil {
			return err
		}
		readers = append(readers, rr)
		files = append(files, f)
	}
	defer func() {
		for _, f := range files {
			_ = f.Close()
			_ = os.Remove(f.Name())
		}
	}()

	out, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer out.Close()
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	var have bool
	var last types.Rec

	for {
		// pick min among readers
		var mi int = -1
		for i, rr := range readers {
			if rr.Has() {
				if mi == -1 || utils.RecLess(rr.Peek(), readers[mi].Peek()) {
					mi = i
				}
			}
		}
		if mi == -1 {
			break
		} // all done

		m := readers[mi].Peek()
		readers[mi].Pop()

		// dedupe across runs
		if !have || m.H != last.H || m.Code != last.Code {
			if err := utils.WriteRec(bw, m); err != nil {
				return err
			}
			last = m
			have = true
		}
	}
	// bubble up read errors if any
	for _, rr := range readers {
		if rr.Err() != nil {
			return rr.Err()
		}
	}
	return nil
}
