package merger

import (
	"bufio"
	"encoding/binary"
	"os"

	"github.com/thanhfphan/kart-challenge/cmd/preprocess/types"
	"github.com/thanhfphan/kart-challenge/cmd/preprocess/utils"
)

// Merge3PairsToValid merges three sorted pairs files and outputs valid codes.
func Merge3PairsToValid(aSorted, bSorted, cSorted, outTxt, optValidBin string) error {
	ra, fa, err := utils.NewReader(aSorted)
	if err != nil {
		return err
	}
	defer fa.Close()

	rb, fb, err := utils.NewReader(bSorted)
	if err != nil {
		fa.Close()
		return err
	}

	defer fb.Close()
	rc, fc, err := utils.NewReader(cSorted)
	if err != nil {
		fa.Close()
		fb.Close()
		return err
	}
	defer fc.Close()

	out, err := os.Create(outTxt)
	if err != nil {
		return err
	}
	defer out.Close()
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	var bin *os.File
	var bwb *bufio.Writer
	if optValidBin != "" {
		bin, err = os.Create(optValidBin)
		if err != nil {
			return err
		}
		defer bin.Close()
		bwb = bufio.NewWriter(bin)
		defer bwb.Flush()
	}

	for {
		okA, okB, okC := ra.Has(), rb.Has(), rc.Has()
		if !okA && !okB && !okC {
			break
		}

		// pick minimum (hash,code)
		var min types.Rec
		set := false
		pick := func(r types.Rec) {
			if !set || utils.RecLess(r, min) {
				min = r
				set = true
			}
		}
		if okA {
			pick(ra.Peek())
		}
		if okB {
			pick(rb.Peek())
		}
		if okC {
			pick(rc.Peek())
		}

		// count appearances
		cnt := 0
		if okA && ra.Peek().H == min.H && ra.Peek().Code == min.Code {
			ra.Pop()
			cnt++
		}
		if okB && rb.Peek().H == min.H && rb.Peek().Code == min.Code {
			rb.Pop()
			cnt++
		}
		if okC && rc.Peek().H == min.H && rc.Peek().Code == min.Code {
			rc.Pop()
			cnt++
		}

		if cnt >= 2 {
			if _, err := bw.WriteString(min.Code + "\n"); err != nil {
				return err
			}
			if bwb != nil {
				if err := binary.Write(bwb, binary.LittleEndian, min.H); err != nil {
					return err
				}
			}
		}
	}

	// Propagate errors
	for _, pr := range []*utils.Reader{ra, rb, rc} {
		if pr.Err() != nil {
			return pr.Err()
		}
	}
	return nil
}
