package bencode

import (
	"reflect"
	"testing"
)

func TestEncode(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"0", args{v: 100}, []byte("i100e")},
		{"1", args{v: "hello"}, []byte("5:hello")},
		{"2", args{v: []interface{}{100, "hello"}}, []byte("li100e5:helloe")},
		{"3", args{v: map[string]interface{}{"hello": 100}}, []byte("d5:helloi100ee")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Encode(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}
