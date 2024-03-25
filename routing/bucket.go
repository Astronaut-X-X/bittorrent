package routing

import "time"

// Bucket a bucket in the routing table contains a peer
type Bucket struct {
	Index int       // Index the index in routing table
	Cap   int       // Cap the capacity of this bucket
	Len   int       // Len the peer quantity of this bucket
	Nodes *SyncList // Nodes a list of node links
}

func NewBucket(index int, size int) *Bucket {
	return &Bucket{
		Index: index,
		Cap:   size,
		Len:   0,
		Nodes: NewSyncList(),
	}
}

func (b *Bucket) GetNode(nodeId string) *Node {
	for _, elem := range b.Nodes.Elements() {
		node := elem.Value.(*Node)
		if node != nil && node.NodeId == nodeId {
			return node
		}
	}
	return nil
}

func (b *Bucket) GetNodes() []*Node {
	nodes := make([]*Node, 0, b.Len)
	for _, elem := range b.Nodes.Elements() {
		node := elem.Value.(*Node)
		nodes = append(nodes, node)
	}
	return nodes
}

func (b *Bucket) Insert(node *Node) {
	if tempNode := b.GetNode(node.NodeId); tempNode != nil {
		tempNode.Create = time.Now()
		return
	}

	if b.Len > b.Cap {
		// TODO ping
		b.Nodes.Remove(b.Nodes.Front())
		b.Nodes.PushBack(node)
		return
	}

	b.Nodes.PushBack(node)
	b.Len++
}

func (b *Bucket) Remove(nodeId string) {
	for _, elem := range b.Nodes.Elements() {
		node := elem.Value.(*Node)
		if node.NodeId == nodeId {
			b.Nodes.Remove(elem)
			b.Len--
			break
		}
	}
}

func (b *Bucket) RefreshBucket(expiration time.Duration) {
	for _, elem := range b.Nodes.Elements() {
		//	TODO ping
		node := elem.Value.(*Node)
		if node.Create.Add(expiration).Before(time.Now()) {
			b.Nodes.Remove(elem)
			b.Len--
		}
	}
}
