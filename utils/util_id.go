package utils

import (
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"math/rand"
)

func RandomID() string {

	num := big.NewInt(int64(rand.Uint32()))
	for i := 1; i < 5; i++ {
		randNum := rand.Int63n(math.MaxUint32)
		num = new(big.Int).Mul(num, big.NewInt(math.MaxUint32))
		num = new(big.Int).Sub(num, big.NewInt(randNum))
	}

	return string(num.Bytes())
}

func RandomToken() string {
	return RandomID()
}

func RandomT() string {
	return RandomID()
}

func FirstIndex(i *big.Int) int {
	zero := big.NewInt(0)

	fmt.Println("FirstIndex", zero.String())

	c := 0
	for i.Cmp(zero) > 0 {
		i = new(big.Int).Div(i, big.NewInt(2))
		c++
	}

	return c
}

func XOR(x, y string) *big.Int {
	ix := toUint(x)
	iy := toUint(y)

	bytesA := ix.Bytes()
	bytesB := iy.Bytes()

	if len(bytesA) != 20 {
		bytesA = append(make([]byte, 20-len(bytesA)), bytesA...)
	}
	if len(bytesB) != 20 {
		bytesB = append(make([]byte, 20-len(bytesB)), bytesB...)
	}

	xorResult := make([]byte, len(bytesA))
	for i := 0; i < len(xorResult); i++ {
		xorResult[i] = bytesA[i] ^ bytesB[i]
	}

	return new(big.Int).SetBytes(xorResult)
}

func toUint(s string) *big.Int {
	hexS := hex.EncodeToString([]byte(s))
	num := new(big.Int)
	num.SetString(hexS, 16)
	return num
}

func ParseIdToByte(id string) []byte {
	return []byte(id)
}
