package krpc

import (
	"net"
	"reflect"
	"testing"
)

func TestParseNodes(t *testing.T) {

	addr, _ := net.ResolveUDPAddr("udp", "192.168.1.1:6881")
	data := []byte{65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 192, 168, 1, 1, 26, 225}

	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []*Node
		wantErr bool
	}{
		{"0", args{data}, []*Node{{Id: "AAAAAAAAAAAAAAAAAAAA", Addr: addr}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseNodes(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseNodes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseNodes() got = %v, want %v", got, tt.want)
			}
		})
	}
}
