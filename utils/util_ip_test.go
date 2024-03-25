package utils

import (
	"net"
	"reflect"
	"testing"
)

func TestParseAddrToByte(t *testing.T) {

	udpAddr, err := net.ResolveUDPAddr("udp", "192.168.1.1:6881")
	if err != nil {
		t.Error(err.Error())
	}

	type args struct {
		addr *net.UDPAddr
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"0", args{addr: udpAddr}, []byte{192, 168, 1, 1, 26, 225}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseAddrToByte(tt.args.addr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseAddrToByte() = %v, want %v", got, tt.want)
			}
		})
	}
}
