package krpc

import (
	"bittorrent/config"
	"bittorrent/logger"
	"sync"
	"time"
)

type TransactionManager struct {
	client         *Client
	transactionMap sync.Map
	keepTime       time.Duration
	tickerTime     time.Duration
	timer          *time.Timer
}

func NewTransactionManager(client *Client, config *config.Config) *TransactionManager {
	manager := &TransactionManager{
		client:         client,
		transactionMap: sync.Map{},
		keepTime:       config.TransactionKeepTime,
		tickerTime:     config.TransactionTickerTime,
	}
	manager.timer = time.AfterFunc(config.TransactionTickerTime, manager.clearTransaction)

	return manager
}

func (m *TransactionManager) Store(t *Transaction) {
	m.transactionMap.Store(t.Query.T, t)
}

func (m *TransactionManager) Load(id string) (*Transaction, bool) {
	value, ok := m.transactionMap.Load(id)
	if !ok {
		return nil, false
	}

	return value.(*Transaction), true
}

func (m *TransactionManager) DeleteById(id string) {
	m.transactionMap.Delete(id)
}

func (m *TransactionManager) Delete(t *Transaction) {
	m.DeleteById(t.Query.T)
}

func (m *TransactionManager) clearTransaction() {
	transactions := make([]*Transaction, 0)

	m.transactionMap.Range(func(key, value any) bool {
		transaction := value.(*Transaction)
		if transaction.Time.Add(m.keepTime).After(time.Now()) {
			return true
		}

		if transaction.NodeQueue.Len() != 0 {
			transactions = append(transactions, transaction)
		}

		logger.Println("[clearTransaction] clear ", transaction.Query.T)
		m.transactionMap.Delete(key)
		return true
	})

	m.ResendTransactions(transactions)
	m.timer.Reset(m.tickerTime)
}

func (m *TransactionManager) Close() {
	m.timer.Stop()
}

func (m *TransactionManager) ResendTransactions(transactions []*Transaction) {
	for _, transaction := range transactions {
		m.Resend(transaction)
	}
}

func (m *TransactionManager) Resend(transaction *Transaction) {
	switch transaction.Query.Q {
	case get_peers:
		if transaction.Query.A == nil {
			return
		}
		infoHash := transaction.Query.A.InfoHash
		m.client.GetPeersContinuous(transaction.NodeQueue, infoHash)
	default:

	}
}

func (m *TransactionManager) Size() int {
	count := 0
	m.transactionMap.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

type Transaction struct {
	Query        *Message
	NodeQueue    *NodeQueue
	ResponseData []byte
	Response     chan bool
	Time         time.Time
}

func NewTransaction(query *Message) *Transaction {
	return &Transaction{
		Query:        query,
		NodeQueue:    NewNodeQueue(16),
		ResponseData: nil,
		Response:     make(chan bool, 1),
		Time:         time.Now(),
	}
}
