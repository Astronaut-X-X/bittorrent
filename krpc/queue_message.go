package krpc

type QueueMessage struct {
	Message *Message
	Node    *Node
}

func NewQueueMessage(message *Message, node *Node) *QueueMessage {
	return &QueueMessage{
		Message: message,
		Node:    node,
	}
}
