package bencode

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func TestDecode(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"0", args{r: bytes.NewBuffer([]byte("i100e"))}, int64(100), false},
		{"1", args{r: bytes.NewBuffer([]byte("5:hello"))}, "hello", false},
		{"2", args{r: bytes.NewBuffer([]byte("li100e5:helloe"))}, []interface{}{int64(100), "hello"}, false},
		{"3", args{r: bytes.NewBuffer([]byte("d5:helloi100ee"))}, map[string]interface{}{"hello": int64(100)}, false},
		{"3", args{r: bytes.NewBuffer([]byte("d5:helloi100eefdsaljgnfdkl"))}, map[string]interface{}{"hello": int64(100)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decode(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decode() got = %v, want %v", got, tt.want)
			}
		})
	}
}
