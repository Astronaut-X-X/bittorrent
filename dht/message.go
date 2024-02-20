package dht

import (
	"bittorrent/bencode"
	"bytes"
	"encoding/json"
	"fmt"
)

const (
	q = "q"
	r = "r"

	ping          = "ping"
	find_node     = "find_node"
	get_peers     = "get_peers"
	announce_peer = "announce_peer"
)

type Message struct {
	T string `json:"t"`
	Y string `json:"y,omitempty"` // 'q'|'r'
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

func DecodeMessage(data []byte) (*Message, error) {
	info, err := bencode.Decode(bytes.NewBuffer(data))
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	info_json, err := json.Marshal(info)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	msg := &Message{}
	if err := json.Unmarshal(info_json, msg); err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return msg, nil
}

func EncodeMessage(msg *Message) []byte {
	msg_json, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	var msg_map map[string]interface{}
	if err := json.Unmarshal(msg_json, &msg_map); err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return bencode.Encode(msg_map)
}
