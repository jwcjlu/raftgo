package raft

import (
	"fmt"
	"io"
)

func WriteInt(writer io.Writer, data int) error {
	_, err := fmt.Fprintf(writer, "%8x", data)
	return err
}

func ReadInt(reader io.Reader) (int, error) {
	var length int
	_, err := fmt.Fscanf(reader, "%8x", &length)
	if err != nil {
		return 0, err
	}
	return length, err
}
