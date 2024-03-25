package utils

import (
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"math/big"
	"math/rand"
	"os"
	"time"
)

func RandomID() string {
	num := big.NewInt(int64(rand.Uint32()))
	for i := 1; i < 5; i++ {
		randNum := rand.Int63n(math.MaxUint32)
		num = new(big.Int).Mul(num, big.NewInt(math.MaxUint32))
		num = new(big.Int).Sub(num, big.NewInt(randNum))
	}

	id := string(num.Bytes())
	return GetStoreId(id)
}

func RandomToken() string {
	return RandomID()
}

func RandomInfoHash() string {
	return RandomID()
}

func RandomT() string {
	rand.NewSource(time.Now().Unix())
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 4)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func IndexByXOR(x, y string) int {
	return FirstIndex(XOR(x, y)) - 1 // 0-160
}

func FirstIndex(i *big.Int) int {
	zero := big.NewInt(0)

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

func GetStoreId(id string) string {
	const filename = "_id.case"

	file, err := os.Open(filename)
	if err != nil && os.IsNotExist(err) {
		file, err = os.Create(filename)
		if err != nil {
			fmt.Println("file.Write err", err.Error())
		}

		if _, err = file.Write([]byte(id)); err != nil {
			fmt.Println("file.Write err", err.Error())
		}

		return id
	}

	data, err := io.ReadAll(file)
	if err != nil {

	}
	return string(data)
}
