package dht

import "encoding/json"

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
