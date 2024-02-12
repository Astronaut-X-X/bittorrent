package bencode

import (
	"bufio"
	"io"
	"strconv"
)

func Decode(r io.Reader) (interface{}, error) {
	reader := reader{*bufio.NewReader(r)}
	return reader.decodeInterface()
}

type reader struct {
	bufio.Reader
}

func (r *reader) decodeInterface() (interface{}, error) {
	ch, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	switch ch {
	case 'i':
		return r.decodeInt()
	case 'l':
		return r.decodeList()
	case 'd':
		return r.decodeDictionary()
	default:
		if err = r.UnreadByte(); err != nil {
			return nil, err
		}
		return r.decodeString()
	}
}

func (r *reader) decodeInt() (int64, error) {
	i, err := r.ReadSlice('e')
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(string(i[:len(i)-1]), 10, 64)
}

func (r *reader) decodeString() (string, error) {
	length, err := r.ReadSlice(':')
	if err != nil {
		return "", err
	}

	l, err := strconv.ParseInt(string(length[:len(length)-1]), 10, 64)

	buf := make([]byte, l)
	_, err = io.ReadFull(r, buf)
	return string(buf), err
}

func (r *reader) decodeList() ([]interface{}, error) {
	list := make([]interface{}, 0)

	for {
		ch, err := r.ReadByte()
		if err != nil {
			return nil, err
		}

		if ch == 'e' {
			break
		}
		if err = r.UnreadByte(); err != nil {
			return nil, err
		}

		item, err := r.decodeInterface()
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}

	return list, nil
}

func (r *reader) decodeDictionary() (map[string]interface{}, error) {
	dict := make(map[string]interface{})

	for {
		ch, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		if ch == 'e' {
			break
		}
		if err = r.UnreadByte(); err != nil {
			return nil, err
		}

		key, err := r.decodeString()
		if err != nil {
			return nil, err
		}
		item, err := r.decodeInterface()
		if err != nil {
			return nil, err
		}
		dict[key] = item
	}

	return dict, nil
}
