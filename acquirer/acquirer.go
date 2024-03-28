package acquirer

import (
	"bittorrent/common"
	"bittorrent/logger"
	"context"
	"time"
)

type PeerInfo struct {
	Ip       string
	Port     int
	InfoHash string
}

func NewPeerInfo(infoHash string, ip string, port int) *PeerInfo {
	return &PeerInfo{
		Ip:       ip,
		Port:     port,
		InfoHash: infoHash,
	}
}

type AcquireManager struct {
	MaxSize      int
	IntervalTime time.Duration
	context      context.Context
	queue        *common.SyncList
	handle       func(*PeerInfo)
}

func NewAcquireManager(ctx context.Context, maxSize int, intervalTime time.Duration) *AcquireManager {
	return &AcquireManager{
		MaxSize:      maxSize,
		IntervalTime: intervalTime,
		context:      ctx,
		queue:        common.NewSyncList(),
	}
}

func (m *AcquireManager) Push(info *PeerInfo) {
	if m.queue.Len() > m.MaxSize {
		return
	}
	m.queue.PushBack(info)
}

func (m *AcquireManager) Pop() *PeerInfo {
	if m.queue.Len() == 0 {
		return nil
	}
	elem := m.queue.Front()
	return elem.Value.(*PeerInfo)
}

func (m *AcquireManager) run() {
	ticker := time.NewTicker(m.IntervalTime)

	for {
		select {
		case <-m.context.Done():
			logger.Println("[AcquireManager] stop")
			return
		case <-ticker.C:
			if info := m.Pop(); info != nil {
				m.handle(info)
			}
		}
	}
}
