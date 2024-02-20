package utils

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
)

func RandomID() string {
	randomData := make([]byte, 20)
	if _, err := io.ReadFull(rand.Reader, randomData); err != nil {
		fmt.Println(err.Error())
		return ""
	}

	hasher := sha1.New()
	hasher.Write(randomData)
	sha1Hash := hasher.Sum(nil)

	return hex.EncodeToString(sha1Hash)
}

func RandomT() string {
	return RandomID()
}

func XOR(x, y string) int64 {
	a := new(big.Int)
	b := new(big.Int)

	a.SetString(x, 16)
	b.SetString(y, 16)

	return new(big.Int).Xor(a, b).Int64()
}
