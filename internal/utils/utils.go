package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"io"
)

func Int64ToBE(n uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, n)
	return b

}
func RandBytes(len int) []byte {
	b := make([]byte, len)
	Must(io.ReadFull(rand.Reader, b))
	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Concat[S ~[]E, E any](slices ...S) []E {
	length := 0
	for _, item := range slices {
		length += len(item)
	}
	if length == 0 {
		return nil
	}
	result := make([]E, 0, length)
	for _, item := range slices {
		result = append(result, item...)
	}
	return result
}
func Sum(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}
func Must[T any](any T, err error) T {
	Throw(err)
	return any
}

func Throw(err error) {
	if err != nil {
		panic(err)
	}
}
