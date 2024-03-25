package routing

import (
	"bittorrent/utils"
	"context"
	"time"
)

const (
	TableSize   = 160
	BucketSize  = 8
	RefreshTime = time.Minute * 15
	PrintTime   = time.Second * 60
)

type IRoutingTable interface {
	Insert(nodeId string, ip string, port int)
	Remove(nodeId string)
	Find(nodeId string) *Node
	Neighbouring(nodeId string) []*Node
}

type Table struct {
	context    context.Context
	NodeId     string
	Bucket     []*Bucket
	Expiration time.Duration
}

func NewRoutingTable(ctx context.Context, nodeId string, expiration time.Duration) IRoutingTable {
	table := &Table{
		context:    ctx,
		Bucket:     make([]*Bucket, TableSize),
		NodeId:     nodeId,
		Expiration: expiration,
	}

	for i := 0; i < TableSize; i++ {
		table.Bucket[i] = NewBucket(i, i+1)
	}

	//go table.RunTimeRefresh()
	//go table.RunTimePrint()

	return table
}

func (r *Table) Insert(nodeId string, ip string, port int) {
	if bucket := r.GetBucket(nodeId); bucket != nil {
		node := NewNode(nodeId, ip, port)
		bucket.Insert(node)
	}
}

func (r *Table) Remove(nodeId string) {
	if bucket := r.GetBucket(nodeId); bucket != nil {
		bucket.Remove(nodeId)
	}
}

func (r *Table) Find(nodeId string) *Node {
	if bucket := r.GetBucket(nodeId); bucket != nil {
		return bucket.GetNode(nodeId)
	}
	return nil
}

func (r *Table) Neighbouring(nodeId string) []*Node {
	li := utils.IndexByXOR(r.NodeId, nodeId)
	ri := li + 1

	nodes := make([]*Node, 0, 8)
	push := func(origin, nodes []*Node) ([]*Node, bool) {
		need, have, full := 8-len(origin), len(nodes), false
		if have >= need {
			have = need
			full = true
		}

		origin = append(origin, nodes[:have]...)
		return origin, full
	}

	for {
		if li > -1 {
			bNodes := r.Bucket[li].GetNodes()
			if nodes, full := push(nodes, bNodes); full {
				return nodes
			}
			li--
		}

		if ri < 160 {
			nodes_ := r.Bucket[li].GetNodes()
			if nodes, full := push(nodes, nodes_); full {
				return nodes
			}
			ri++
		}

		if li < 0 && ri > 159 {
			break
		}
	}

	return nodes
}

func (r *Table) GetBucket(nodeId string) *Bucket {
	i := utils.IndexByXOR(r.NodeId, nodeId)
	if i == -1 {
		return nil
	}
	if i >= TableSize || i < 0 {
		return nil
	}

	return r.Bucket[i]
}
