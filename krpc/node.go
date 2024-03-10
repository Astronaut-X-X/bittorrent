package krpc

type Node struct {
	Id   string
	Ip   string
	Port int
}

func NewNode(id string, ip string, port int) *Node {
	return &Node{
		Id:   id,
		Ip:   ip,
		Port: port,
	}
}
