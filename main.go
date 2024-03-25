package main

import (
	"bittorrent/config"
	"bittorrent/dht"
	"fmt"
)

func main() {
	d, err := dht.NewDHT(config.DefaultConfig)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	d.Run()
}
