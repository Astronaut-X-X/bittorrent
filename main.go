package main

import (
	"bittorrent/dht"
	"fmt"
)

func main() {
	d, err := dht.NewDHT(dht.DefaultConfig())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	d.Run()
}
