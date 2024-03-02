package routingTable

import (
	"context"
	"fmt"
	"net"
	"reflect"
	"sync"
	"testing"
)

func TestNewRoutingTable(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	table := NewRoutingTable(ctx)
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

func TestRoutingTable_Add(t *testing.T) {
	type fields struct {
		L        sync.Mutex
		context  context.Context
		Bucket   []*Bucket
		LocalId  string
		pingPeer func(addr *net.UDPAddr) bool
	}
	type args struct {
		id   string
		ip   string
		port int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RoutingTable{
				L:        tt.fields.L,
				context:  tt.fields.context,
				Bucket:   tt.fields.Bucket,
				LocalId:  tt.fields.LocalId,
				pingPeer: tt.fields.pingPeer,
			}
			if err := r.Add(tt.args.id, tt.args.ip, tt.args.port); (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRoutingTable_GetBucket(t *testing.T) {
	type fields struct {
		L        sync.Mutex
		context  context.Context
		Bucket   []*Bucket
		LocalId  string
		pingPeer func(addr *net.UDPAddr) bool
	}
	type args struct {
		x string
		y string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Bucket
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RoutingTable{
				L:        tt.fields.L,
				context:  tt.fields.context,
				Bucket:   tt.fields.Bucket,
				LocalId:  tt.fields.LocalId,
				pingPeer: tt.fields.pingPeer,
			}
			if got := r.GetBucket(tt.args.x, tt.args.y); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBucket() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoutingTable_GetPeers(t *testing.T) {
	type fields struct {
		L        sync.Mutex
		context  context.Context
		Bucket   []*Bucket
		LocalId  string
		pingPeer func(addr *net.UDPAddr) bool
	}
	type args struct {
		x string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*Peer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RoutingTable{
				L:        tt.fields.L,
				context:  tt.fields.context,
				Bucket:   tt.fields.Bucket,
				LocalId:  tt.fields.LocalId,
				pingPeer: tt.fields.pingPeer,
			}
			if got := r.GetPeers(tt.args.x); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPeers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoutingTable_PrintRoutingTable(t *testing.T) {
	type fields struct {
		L        sync.Mutex
		context  context.Context
		Bucket   []*Bucket
		LocalId  string
		pingPeer func(addr *net.UDPAddr) bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RoutingTable{
				L:        tt.fields.L,
				context:  tt.fields.context,
				Bucket:   tt.fields.Bucket,
				LocalId:  tt.fields.LocalId,
				pingPeer: tt.fields.pingPeer,
			}
			r.PrintRoutingTable()
		})
	}
}

func TestRoutingTable_RefreshAllBucket(t *testing.T) {
	type fields struct {
		L        sync.Mutex
		context  context.Context
		Bucket   []*Bucket
		LocalId  string
		pingPeer func(addr *net.UDPAddr) bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RoutingTable{
				L:        tt.fields.L,
				context:  tt.fields.context,
				Bucket:   tt.fields.Bucket,
				LocalId:  tt.fields.LocalId,
				pingPeer: tt.fields.pingPeer,
			}
			r.RefreshAllBucket()
		})
	}
}

func TestRoutingTable_RunTimePrint(t *testing.T) {
	type fields struct {
		L        sync.Mutex
		context  context.Context
		Bucket   []*Bucket
		LocalId  string
		pingPeer func(addr *net.UDPAddr) bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RoutingTable{
				L:        tt.fields.L,
				context:  tt.fields.context,
				Bucket:   tt.fields.Bucket,
				LocalId:  tt.fields.LocalId,
				pingPeer: tt.fields.pingPeer,
			}
			r.RunTimePrint()
		})
	}
}

func TestRoutingTable_RunTimeRefresh(t *testing.T) {
	type fields struct {
		L        sync.Mutex
		context  context.Context
		Bucket   []*Bucket
		LocalId  string
		pingPeer func(addr *net.UDPAddr) bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RoutingTable{
				L:        tt.fields.L,
				context:  tt.fields.context,
				Bucket:   tt.fields.Bucket,
				LocalId:  tt.fields.LocalId,
				pingPeer: tt.fields.pingPeer,
			}
			r.RunTimeRefresh()
		})
	}
}

func TestRoutingTable_SetPingPeer(t *testing.T) {
	type fields struct {
		L        sync.Mutex
		context  context.Context
		Bucket   []*Bucket
		LocalId  string
		pingPeer func(addr *net.UDPAddr) bool
	}
	type args struct {
		pingPeer func(addr string) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RoutingTable{
				L:        tt.fields.L,
				context:  tt.fields.context,
				Bucket:   tt.fields.Bucket,
				LocalId:  tt.fields.LocalId,
				pingPeer: tt.fields.pingPeer,
			}
			r.SetPingPeer(tt.args.pingPeer)
		})
	}
}
