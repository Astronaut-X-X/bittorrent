package routingTable

import (
	"testing"
)

func TestNewPeer(t *testing.T) {
	type args struct {
		id   string
		ip   string
		port int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"0", args{id: "AAAAABBBBBCCCCCDDDDD", ip: "192.168.1.10", port: 1234}, false},
		{"0", args{id: "AAAAABBBBBCCCCCDDDDD", ip: "192.168.1.10", port: 123433}, true},
		{"0", args{id: "AAAAABBBBBCCCCC", ip: "192.168.1.10", port: 123433}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewPeer(tt.args.id, tt.args.ip, tt.args.port)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPeer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
