package dht

import (
	routingTable "bittorrent/routingTabel"
	"sync"
	"time"
)

const Timeout = time.Second * 60

type TransactionManager struct {
	TransactionMap sync.Map
}

type Transaction struct {
	Id           string
	DHT          *DHT
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

func NewTransaction(id string, d *DHT, query *Message, afterFunc func(*Transaction)) *Transaction {
	t := &Transaction{
		Id:           id,
		DHT:          d,
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
	m.TransactionMap.Store(t.Id, t)
}

func (m *TransactionManager) Load(id string) (*Transaction, bool) {
	value, ok := m.TransactionMap.Load(id)
	if !ok {
		return nil, false
	}

	return value.(*Transaction), true
}

func (m *TransactionManager) Delete(t *Transaction) {
	m.TransactionMap.Delete(t.Id)
}
