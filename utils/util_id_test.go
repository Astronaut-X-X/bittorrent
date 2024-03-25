package utils

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"
)

func TestRandomID(t *testing.T) {
	got := RandomID()
	t.Log(hex.EncodeToString([]byte(got)))
	if len(got) != 20 {
		t.Error("error")
	}
}

func TestXOR(t *testing.T) {
	x := string([]byte{
		0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000000,
		0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000010,
	})
	y := string([]byte{
		0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000000,
		0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000011,
	})
	if XOR(x, y).Cmp(big.NewInt(1)) != 0 {
		t.Error("ERROR")
	}
}

func Test_toUint(t *testing.T) {
	x := RandomID()
	fmt.Println(toUint(x))
}

func TestParseIdToByte(t *testing.T) {
	id := RandomID()
	byteId := ParseIdToByte(id)
	fmt.Println(byteId)
}

func TestFirstIndex(t *testing.T) {
	x := []byte{0, 0, 0, 79, 19, 76, 198, 91, 152, 71, 105, 119, 147, 133, 163, 86, 188, 170, 238, 74}
	y := []byte{0, 0, 0, 105, 115, 81, 255, 74, 236, 41, 205, 186, 171, 242, 251, 227, 70, 124, 194, 103}
	distance := XOR(string(x), string(y))
	i := FirstIndex(distance)
	t.Log(i)

	x = []byte{0, 0, 0, 79, 19, 76, 198, 91, 152, 71, 105, 119, 147, 133, 163, 86, 188, 170, 238, 74}
	y = []byte{0, 0, 0, 79, 19, 76, 198, 91, 152, 71, 105, 119, 147, 133, 163, 86, 188, 170, 238, 74}
	distance = XOR(string(x), string(y))
	i = FirstIndex(distance)
	t.Log(i)

	x = []byte{170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170}
	y = []byte{85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85}
	distance = XOR(string(x), string(y))
	i = FirstIndex(distance)
	t.Log(i)
}

func TestFirstIndex_(t *testing.T) {
	num := big.NewInt(1024)
	i := FirstIndex(num)
	t.Log(i)
}
