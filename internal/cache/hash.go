package cache

import (
	"encoding/binary"
)

func Uint16Bytes(v uint16) []byte {
	a := make([]byte, 2)
	binary.LittleEndian.PutUint16(a, v)

	return a
}

func Uint32Bytes(v uint32) []byte {
	a := make([]byte, 4)
	binary.LittleEndian.PutUint32(a, v)

	return a
}

func Uint64Bytes(v uint64) []byte {
	a := make([]byte, 8)
	binary.LittleEndian.PutUint64(a, v)

	return a
}
