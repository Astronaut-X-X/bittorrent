package dht

import "encoding/json"

type Message struct {
	T string `json:"t"`
	Y string `json:"y"`
	Q string `json:"q"`
	A *A     `json:"a"`
	R *R     `json:"r"`
}

type A struct {
	Id          string `json:"id"`
	Target      string `json:"targer"`
	InfoHash    string `json:"info_hash"`
	ImpliedPort int    `json:"implied_port"`
	Port        int    `json:"port"`
	Token       string `json:"token"`
}

type R struct {
	Id     string   `json:"id"`
	Nodes  string   `json:"nodes"`
	Token  string   `json:"token"`
	Values []string `json:"values"`
}

func UnmarshalMessage(data []byte) *Message {
	msg := &Message{}

	if err := json.Unmarshal(data, msg); err != nil {
		return nil
	}

	return msg
}

func MarshalMessage(msg *Message) ([]byte, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return data, nil
}
