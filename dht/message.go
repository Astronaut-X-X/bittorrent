package dht

import (
	"bittorrent/bencode"
	"bytes"
	"encoding/json"
	"fmt"
)

type Message struct {
	T string `json:"t"`
	Y string `json:"y,omitempty"`
	Q string `json:"q,omitempty"`
	A *A     `json:"a,omitempty"`
	R *R     `json:"r,omitempty"`
}

type A struct {
	Id          string `json:"id"`
	Target      string `json:"targer,omitempty"`
	InfoHash    string `json:"info_hash,omitempty"`
	ImpliedPort int    `json:"implied_port,omitempty"`
	Port        int    `json:"port,omitempty"`
	Token       string `json:"token,omitempty"`
}

type R struct {
	Id     string   `json:"id"`
	Nodes  string   `json:"nodes,omitempty"`
	Token  string   `json:"token,omitempty"`
	Values []string `json:"values,omitempty"`
}

func UnmarshalMessage(data []byte) (*Message, error) {
	msg := &Message{}

	msg_, err := bencode.Decode(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	msg_byte, err := json.Marshal(msg_)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(msg_byte, msg); err != nil {
		return nil, err
	}

	return msg, nil
}

func MarshalMessage(msg *Message) []byte {
	byte_message, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	var map_message map[string]interface{}
	if err := json.Unmarshal(byte_message, &map_message); err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return bencode.Encode(msg)
}
