package bencode

import (
	"bytes"
	"reflect"
	"strconv"
)

func Encode(v any) []byte {
	var b buffer
	b.encodeInterface(v)
	return b.Bytes()
}

type buffer struct {
	bytes.Buffer
}

func (b *buffer) encodeInterface(v interface{}) {
	switch v := v.(type) {
	case int, int8, int16, int32, int64:
		b.encodeInt(reflect.ValueOf(v).Int())
	case uint, uint8, uint16, uint32, uint64:
		b.encodeUint(reflect.ValueOf(v).Uint())
	case string:
		b.encodeString(reflect.ValueOf(v).String())
	case []interface{}:
		b.encodeList(v)
	case map[string]interface{}:
		b.encodeDictionary(v)
	}
}

func (b *buffer) encodeInt(v int64) {
	b.WriteByte('i')
	b.WriteString(strconv.FormatInt(v, 10))
	b.WriteByte('e')
}

func (b *buffer) encodeUint(v uint64) {
	b.WriteByte('i')
	b.WriteString(strconv.FormatUint(v, 10))
	b.WriteByte('e')
}

func (b *buffer) encodeString(v string) {
	b.WriteString(strconv.Itoa(len(v)))
	b.WriteByte(':')
	b.WriteString(v)
}

func (b *buffer) encodeList(L []interface{}) {
	b.WriteByte('l')
	for _, v := range L {
		b.encodeInterface(v)
	}
	b.WriteByte('e')
}

func (b *buffer) encodeDictionary(D map[string]interface{}) {
	b.WriteByte('d')
	for k, v := range D {
		b.encodeString(k)
		b.encodeInterface(v)
	}
	b.WriteByte('e')
}
