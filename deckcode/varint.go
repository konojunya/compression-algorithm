package deckcode

import (
	"encoding/binary"
	"errors"
)

func putUvarint(dst []byte, x uint64) []byte {
	var tmp [10]byte
	n := binary.PutUvarint(tmp[:], x)
	return append(dst, tmp[:n]...)
}

func readUvarint(src []byte, cur *int) (uint64, error) {
	if *cur >= len(src) {
		return 0, errors.New("uvarint: eof")
	}
	v, n := binary.Uvarint(src[*cur:])
	if n <= 0 {
		return 0, errors.New("uvarint: invalid")
	}
	*cur += n
	return v, nil
}
