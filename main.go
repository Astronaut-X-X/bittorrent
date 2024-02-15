package main

import (
	"bittorrent/bencode"
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
)

func main() {
	open, err := os.Open("left-4-dead-2.torrent")
	if err != nil {
		return
	}
	defer open.Close()

	r := bufio.NewReader(open)

	data, err := bencode.Decode(r)
	if err != nil {
		return
	}

	formatData, ok := data.(map[string]interface{})
	if !ok {
		return
	}
	info := formatData["info"].(map[string]interface{})
	pieces := info["pieces"].(string)
	l := len(pieces)
	for i := 0; i < l; i += 20 {
		fmt.Println(hex.EncodeToString([]byte(pieces[i : i+20])))
	}
}
