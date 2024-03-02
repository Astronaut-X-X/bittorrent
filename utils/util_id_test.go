package utils

import (
	"fmt"
	"math/big"
	"testing"
)

func TestRandomID(t *testing.T) {
	got := RandomID()
	fmt.Println([]byte(got))
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
	for i := 0; i < 20; i++ {
		x := RandomID()
		y := RandomID()

		distance := XOR(x, y)
		i := FirstIndex(distance)
		fmt.Println(i)
	}
}
