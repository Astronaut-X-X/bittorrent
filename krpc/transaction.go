package krpc

import (
	"bittorrent/config"
	"sync"
	"time"
)

type TransactionManager struct {
	TransactionMap sync.Map
	keepTime       time.Duration
	tickerTime     time.Duration
	timer          *time.Timer
}

type Transaction struct {
	Query        *Message
	NodeQueue    *NodeQueue
	ResponseData []byte
	Response     chan bool
	Time         time.Time
	Timer        time.Timer
}

func NewTransactionManager(config *config.Config) *TransactionManager {
	manager := &TransactionManager{
		TransactionMap: sync.Map{},
		keepTime:       config.TransactionKeepTime,
		tickerTime:     config.TransactionTickerTime,
	}
	manager.timer = time.AfterFunc(config.TransactionTickerTime, manager.clearTransaction)

	return manager
}

func NewTransaction(query *Message) *Transaction {
	return &Transaction{
		Query:        query,
		NodeQueue:    NewNodeQueue(16),
		ResponseData: nil,
		Response:     make(chan bool, 1),
	}
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

func (m *TransactionManager) clearTransaction() {
	m.TransactionMap.Range(func(key, value any) bool {
		transaction := value.(*Transaction)
		if transaction.Time.Add(m.keepTime).Before(time.Now()) {
			m.TransactionMap.Delete(key)
		}
		return true
	})
	m.timer.Reset(m.tickerTime)
}

func (m *TransactionManager) Close() {
	m.timer.Stop()
}
