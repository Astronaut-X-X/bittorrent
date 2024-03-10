package routingTable

import (
	"context"
	"fmt"
	"testing"
)

func TestNewRoutingTable(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	table := NewRoutingTable("", ctx)
	//id := table.LocalId
	peerId := string([]byte{
		0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000000,
		0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000000, 0b00000001,
	})

	if err := table.Add(peerId, "192.168.1.1", 6881); err != nil {
		t.Error("add peer error")
	}

	peers := table.GetPeers(peerId)
	for _, peer := range peers {
		fmt.Println(peer.Id)
	}
}
