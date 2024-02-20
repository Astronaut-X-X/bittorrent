package dht

import (
	"sync"
	"time"
)

type TransactionManager struct {
	TransactionMap sync.Map
}

type Transaction struct {
	Id           string
	Query        *Message
	ResponseData []byte
	Response     chan bool
}

func NewTransactionManager() *TransactionManager {
	return &TransactionManager{
		TransactionMap: sync.Map{},
	}
}

func NewTransaction(id string, query *Message) *Transaction {
	t := &Transaction{
		Id:           id,
		Query:        query,
		ResponseData: nil,
		Response:     make(chan bool),
	}

	// Transaction timeout
	time.AfterFunc(time.Second*60, func() {
		t.Response <- false
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
