package krpc

import (
	routingTable "bittorrent/routingTabel"
	"sync"
	"time"
)

const Timeout = time.Second * 20

type TransactionManager struct {
	TransactionMap sync.Map
}

type Transaction struct {
	Query        *Message
	Peers        []*routingTable.Peer
	ResponseData []byte
	Response     chan bool
	timer        *time.Timer
}

func NewTransactionManager() *TransactionManager {
	return &TransactionManager{
		TransactionMap: sync.Map{},
	}
}

func NewTransaction(query *Message, afterFunc func(*Transaction)) *Transaction {
	t := &Transaction{
		Query:        query,
		Peers:        make([]*routingTable.Peer, 0),
		ResponseData: nil,
		Response:     make(chan bool, 1),
		timer:        nil,
	}

	// Transaction timeout
	t.timer = time.AfterFunc(Timeout, func() {
		afterFunc(t)
	})

	return t
}

func (m *TransactionManager) Store(t *Transaction) {
	m.TransactionMap.Store(t.Query.T, t)
}

func (m *TransactionManager) Load(id string) (*Transaction, bool) {
	value, ok := m.TransactionMap.Load(id)
	if !ok {
		return nil, false
	}

	return value.(*Transaction), true
}

func (m *TransactionManager) DeleteById(id string) {
	m.TransactionMap.Delete(id)
}

func (m *TransactionManager) Delete(t *Transaction) {
	m.DeleteById(t.Query.T)
}
