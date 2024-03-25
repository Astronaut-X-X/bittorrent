package routing

import (
	"testing"
)

func TestNewBucket(t *testing.T) {
	index, size := 10, 10

	bucket := NewBucket(10, 10)

	if bucket.Index != index || bucket.Cap != size {
		t.Error("new bucket err")
	}
}

func TestBucket_Insert(t *testing.T) {
	bucket := NewBucket(10, 3)

	node1 := NewNode("test001-000000000000", "192.168.1.1", 6881)
	node2 := NewNode("test002-000000000000", "192.168.1.1", 6881)
	node3 := NewNode("test003-000000000000", "192.168.1.1", 6881)
	node_ := NewNode("test003-000000000000", "192.168.1.100", 6881)
	node4 := NewNode("test004-000000000000", "192.168.1.100", 6881)

	bucket.Insert(node1)
	bucket.Insert(node2)
	bucket.Insert(node3)

	bucket.Insert(node_)

	bucket.Insert(node4)
}

func TestBucket_GetNode(t *testing.T) {
	bucket := NewBucket(10, 3)

	node1 := NewNode("test001-000000000000", "192.168.1.1", 6881)
	node2 := NewNode("test002-000000000000", "192.168.1.1", 6881)
	node3 := NewNode("test003-000000000000", "192.168.1.1", 6881)

	bucket.Insert(node1)
	bucket.Insert(node2)
	bucket.Insert(node3)

	bucket.GetNode("test001-000000000000")
}

func TestBucket_GetNodes(t *testing.T) {
	bucket := NewBucket(10, 3)

	node1 := NewNode("test001-000000000000", "192.168.1.1", 6881)
	node2 := NewNode("test002-000000000000", "192.168.1.1", 6881)
	node3 := NewNode("test003-000000000000", "192.168.1.1", 6881)

	bucket.Insert(node1)
	bucket.Insert(node2)
	bucket.Insert(node3)

	nodes := bucket.GetNodes()
	for _, node := range nodes {
		t.Log(node.NodeId)
	}
}
