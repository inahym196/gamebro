package main

import (
	"bufio"
	"fmt"
	"io"
)

type BPP2Writer struct {
	bw  *bufio.Writer
	buf [2]byte
}

func NewBPP2Writer(w io.Writer) *BPP2Writer {
	return &BPP2Writer{
		bw: bufio.NewWriter(w),
	}
}

func (w *BPP2Writer) Write(bpp2 uint16) error {
	w.buf[0] = byte(bpp2 >> 8) // hi
	w.buf[1] = byte(bpp2)      // lo
	if _, err := w.bw.Write(w.buf[:]); err != nil {
		return fmt.Errorf("failed to write bpp2: %w", err)
	}
	return nil
}

func (w *BPP2Writer) Flush() error {
	if err := w.bw.Flush(); err != nil {
		return fmt.Errorf("falied to flush writer: %w", err)
	}
	return nil
}
