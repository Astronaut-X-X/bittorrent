package dht

import (
	"bittorrent/bencode"
	"bytes"
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

	msg := mapToMessage(info.(map[string]interface{}))
	return msg, nil
}

func EncodeMessage(msg *Message) []byte {
	msgMap := messageToMap(msg)
	return bencode.Encode(msgMap)
}

func messageToMap(msg *Message) map[string]interface{} {
	msgMap := map[string]interface{}{}
	if msg.T != "" {
		msgMap["t"] = msg.T
	}
	if msg.Y != "" {
		msgMap["y"] = msg.Y
	}
	if msg.Q != "" {
		msgMap["q"] = msg.Q
	}
	if msg.A != nil {
		msgMap["a"] = aToMap(msg.A)
	}
	if msg.R != nil {
		msgMap["r"] = rToMap(msg.R)
	}
	return msgMap
}

func mapToMessage(msgMap map[string]interface{}) *Message {
	msg := &Message{}
	if msgMap["t"] != nil {
		msg.T = msgMap["t"].(string)
	}
	if msgMap["y"] != nil {
		msg.Y = msgMap["y"].(string)
	}
	if msgMap["q"] != nil {
		msg.Q = msgMap["q"].(string)
	}
	if msgMap["a"] != nil {
		msg.A = mapToA(msgMap["a"].(map[string]interface{}))
	}
	if msgMap["r"] != nil {
		msg.R = mapToR(msgMap["r"].(map[string]interface{}))
	}
	return msg
}

func aToMap(A *A) map[string]interface{} {
	aMap := map[string]interface{}{}
	if A.Id != "" {
		aMap["id"] = A.Id
	}
	if A.InfoHash != "" {
		aMap["info_hash"] = A.InfoHash
	}
	if A.Target != "" {
		aMap["targer"] = A.Target
	}
	if A.ImpliedPort != 0 {
		aMap["implied_port"] = A.ImpliedPort
	}
	if A.Port != 0 {
		aMap["port"] = A.Port
	}
	if A.Token != "" {
		aMap["token"] = A.Token
	}
	return aMap
}

func mapToA(aMap map[string]interface{}) *A {
	A := &A{}
	if aMap["id"] != nil {
		A.Id = aMap["id"].(string)
	}
	if aMap["info_hash"] != nil {
		A.InfoHash = aMap["info_hash"].(string)
	}
	if aMap["targer"] != nil {
		A.Target = aMap["targer"].(string)
	}
	if aMap["implied_port"] != nil {
		A.ImpliedPort = aMap["implied_port"].(int)
	}
	if aMap["port"] != nil {
		A.Port = aMap["port"].(int)
	}
	if aMap["token"] != nil {
		A.Token = aMap["token"].(string)
	}
	return A
}

func rToMap(R *R) map[string]interface{} {
	rMap := map[string]interface{}{}
	if R.Id != "" {
		rMap["id"] = R.Id
	}
	if R.Nodes != "" {
		rMap["nodes"] = R.Nodes
	}
	if R.Token != "" {
		rMap["token"] = R.Token
	}
	if len(R.Values) != 0 {
		rMap["values"] = R.Values
	}
	return rMap
}

func mapToR(rMap map[string]interface{}) *R {
	R := &R{}
	if rMap["id"] != nil {
		R.Id = rMap["id"].(string)
	}
	if rMap["nodes"] != nil {
		R.Nodes = rMap["nodes"].(string)
	}
	if rMap["token"] != nil {
		R.Token = rMap["token"].(string)
	}
	if rMap["values"] != nil {
		R.Values = rMap["values"].([]string)
	}
	return R
}
