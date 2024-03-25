package routing

import (
	"bittorrent/utils"
	"context"
	"fmt"
	"testing"
	"time"
)

func TestNewRoutingTable(t *testing.T) {
	id := utils.RandomID()
	routing := NewRoutingTable(context.Background(), id, time.Second*60)
	byteId0 := []byte(id)
	bit := byteId0[len(byteId0)-1]
	byteId0[len(byteId0)-1] = bit + 1
	routing.Insert(string(byteId0), "192.168.1.1", 6881)

	byteId1 := []byte(id)
	bit1 := byteId1[len(byteId1)-1]
	byteId1[len(byteId1)-1] = bit1 + 1
	routing.Insert(string(byteId1), "192.168.1.1", 6881)

	node := routing.Find(string(byteId1))
	fmt.Println(node.NodeId)
}
