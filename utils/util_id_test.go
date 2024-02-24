package utils

import (
	"testing"
)

func TestRandomID(t *testing.T) {
	got := RandomID()
	if len(got) != 20 {
		t.Error("error")
	}
}
