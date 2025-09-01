package utils

import (
	"bufio"
	"os"

	"github.com/thanhfphan/kart-challenge/cmd/preprocess/types"
)

// Reader is a generic reader for records.
type Reader struct {
	br   *bufio.Reader
	cur  types.Rec
	ok   bool
	done bool
	err  error
}

func NewReader(name string) (*Reader, *os.File, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, nil, err
	}

	return &Reader{br: bufio.NewReader(f)}, f, nil
}

// Has returns true if there are more records.
func (r *Reader) Has() bool {
	if r.done {
		return false
	}

	if r.ok {
		return true
	}

	rec, ok, err := ReadPair(r.br)
	if err != nil {
		r.err = err
		r.done = true
		return false
	}
	if !ok {
		r.done = true
		return false
	}

	r.cur, r.ok = rec, true
	return true
}

func (r *Reader) Peek() types.Rec { return r.cur }
func (r *Reader) Pop()            { r.ok = false }
func (r *Reader) Err() error      { return r.err }
