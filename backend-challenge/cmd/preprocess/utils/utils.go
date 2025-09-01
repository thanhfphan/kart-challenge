package utils

import (
	"bufio"
	"encoding/binary"
	"io"

	"github.com/cespare/xxhash/v2"
	"github.com/thanhfphan/kart-challenge/cmd/preprocess/types"
)

// WritePair writes a record to the writer.
func WritePair(bw *bufio.Writer, code string) error {
	h := xxhash.Sum64String(code)
	if err := binary.Write(bw, binary.LittleEndian, h); err != nil {
		return err
	}

	if err := binary.Write(bw, binary.LittleEndian, uint16(len(code))); err != nil {
		return err
	}

	_, err := bw.WriteString(code)
	return err
}

// ReadPair reads a record from the reader.
func ReadPair(br *bufio.Reader) (types.Rec, bool, error) {
	var h uint64
	if err := binary.Read(br, binary.LittleEndian, &h); err != nil {
		if err == io.EOF {
			return types.Rec{}, false, nil
		}
		return types.Rec{}, false, err
	}

	var ln uint16
	if err := binary.Read(br, binary.LittleEndian, &ln); err != nil {
		return types.Rec{}, false, err
	}

	b := make([]byte, ln)
	if _, err := io.ReadFull(br, b); err != nil {
		return types.Rec{}, false, err
	}

	return types.Rec{H: h, Code: string(b)}, true, nil
}

// RecLess compares two records.
// Returns true if a < b
func RecLess(a, b types.Rec) bool {
	if a.H != b.H {
		return a.H < b.H
	}

	return a.Code < b.Code
}

// WriteRec writes a record to the writer.
func WriteRec(bw *bufio.Writer, r types.Rec) error {
	if err := binary.Write(bw, binary.LittleEndian, r.H); err != nil {
		return err
	}

	if err := binary.Write(bw, binary.LittleEndian, uint16(len(r.Code))); err != nil {
		return err
	}

	_, err := bw.WriteString(r.Code)
	return err
}
