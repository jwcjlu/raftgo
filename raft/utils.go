package raft

import (
	"fmt"
	"io"
)

func WriteInt(w io.Writer, data int) error {
	_, err := fmt.Fprintf(w, "%8x\n", data)
	return err
}

func ReadInt(r io.Reader) (int, error) {
	var length int
	_, err := fmt.Fscanf(r, "%8x\n", &length)
	if err != nil {
		return 0, err
	}
	return length, err
}
