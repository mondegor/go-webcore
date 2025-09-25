package mrcrypt

import (
	"crypto/rand"
	"errors"
	"math"
)

const (
	charsetDigit  = "0123456789"
	charsetHex    = "0123456789abcdef"
	charsetToken  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	maxCharsetLen = 128
)

// GenerateDigits - comment func.
func GenerateDigits(length int) (string, error) {
	return GenerateSequence([]byte(charsetDigit), length)
}

// GenerateHex - comment func.
func GenerateHex(length int) (string, error) {
	return GenerateSequence([]byte(charsetHex), length)
}

// GenerateToken - comment func.
func GenerateToken(length int) (string, error) {
	return GenerateSequence([]byte(charsetToken), length)
}

// GenerateSequence - comment func.
func GenerateSequence(charset []byte, length int) (string, error) {
	s, err := GenerateBytes(charset, length)
	if err != nil {
		return "", err
	}

	return string(s), nil
}

// GenerateBytes - comment func.
func GenerateBytes(charset []byte, length int) ([]byte, error) {
	if len(charset) > maxCharsetLen {
		return nil, errors.New("charset length exceeds max length")
	}

	if length < 1 {
		return nil, errors.New("length less than 1")
	}

	chunk := make([]byte, length*2)
	sequence := chunk[:0]
	indexes := chunk[length:]

	bits100 := uint64(math.Log2(float64(len(charset))) * 100)
	bits := bits100 / 100

	if bits100%100 != 0 {
		bits++
	}

	mask := uint8(1<<bits) - 1

	for {
		if _, err := rand.Read(indexes); err != nil {
			return nil, err
		}

		read := 0

		for i := 0; i < len(indexes); i++ {
			rnd := indexes[i] & mask

			if int(rnd) < len(charset) {
				sequence = append(sequence, charset[rnd])
				read++
			}
		}

		if read >= len(indexes) {
			return sequence[0:length:length], nil
		}

		indexes = indexes[:len(indexes)-read]
	}
}
