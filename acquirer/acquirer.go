package acquirer

import (
	"bittorrent/bencode"
	"bittorrent/common"
	"bittorrent/config"
	"bittorrent/logger"
	"bittorrent/utils"
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
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
	config       *config.Config
	context      context.Context
	queue        *common.SyncList
	checkMap     sync.Map
	handle       func(*PeerInfo)
}

func NewAcquireManager(ctx context.Context, config *config.Config) *AcquireManager {
	manager := &AcquireManager{
		MaxSize:      config.AcquirerMaxSize,
		IntervalTime: config.AcquirerIntervalTime,
		config:       config,
		context:      ctx,
		queue:        common.NewSyncList(),
	}
	go manager.run()
	manager.SetHandle(handle)

	return manager
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
	m.queue.Remove(elem)
	return elem.Value.(*PeerInfo)
}

func (m *AcquireManager) SetHandle(f func(*PeerInfo)) {
	m.handle = f
}

func (m *AcquireManager) run() {
	ticker := time.NewTicker(m.IntervalTime)

	for {
		select {
		case <-m.context.Done():
			logger.Println("[AcquireManager] stop")
			return
		case <-ticker.C:
			fmt.Println("[AcquireManager] run handle")
			if info := m.Pop(); info != nil {
				go m.handle(info)
			}
		}
	}
}

func handle(info *PeerInfo) {
	fmt.Println("[handle]", info.InfoHash, " ", info.Ip, " ", info.Port)
	acquirer, err := NewAcquirer(info.InfoHash, info.Ip, info.Port)
	if err != nil {
		return
	}
	defer acquirer.close()
	if err = acquirer.sendHandshake(); err != nil {
		return
	}
	if err = acquirer.readHandshake(); err != nil {
		logger.Println(err.Error())
		return
	}
	if err = acquirer.sendExtHandshake(); err != nil {
		logger.Println(err.Error())
		return
	}
	if err = acquirer.readMessage(); err != nil {
		logger.Println(err.Error())
		return
	}
}

type Acquirer struct {
	conn     net.Conn
	infoHash string
	peerId   string
	error    error
}

func NewAcquirer(infoHash string, ip string, port int) (*Acquirer, error) {
	const timeout = time.Second * 15
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), timeout)
	if err != nil {
		logger.Println("[Acquirer] tcp connect err: %v", err.Error())
		return nil, err
	}

	acquirer := &Acquirer{
		conn:     conn,
		infoHash: infoHash,
		peerId:   utils.RandomID(),
	}

	return acquirer, nil
}

func (a *Acquirer) close() {
	if a.conn != nil {
		if err := a.conn.Close(); err != nil {
			logger.Println("[Acquirer] conn close err: %v", err.Error())
		}
	}
}

func (a *Acquirer) sendHandshake() error {
	handshake := Handshake{}
	handshake.InfoHash = []byte(a.infoHash)
	handshake.PeerId = []byte(a.peerId)
	data := handshake.Serialize()

	if err := a.conn.SetWriteDeadline(time.Now().Add(time.Second + 15)); err != nil {
		return err
	}

	n, err := a.conn.Write(data)
	if err != nil {
		return err
	}

	logger.Println("[Acquirer] sendHandshake done ", n)
	return nil
}

func (a *Acquirer) readHandshake() error {
	buf := make([]byte, 1024)
	n, err := a.conn.Read(buf)
	if err != nil {
		return err
	}

	logger.Println("[readHandshake] ", buf[:n])
	lbt := len(BitTorrentProtocol)

	if n < lbt+49 {
		return errors.New("error data length")
	}

	if buf[0] != byte(0x13) {
		return errors.New("error data first bit")
	}

	if string(buf[1:lbt+1]) != BitTorrentProtocol {
		return errors.New("error BitTorrent protocol data")
	}

	if string(buf[lbt+1:lbt+9]) != string([]byte{0, 0, 0, 0, 0, 16, 0, 0}) {
		return errors.New("error reserved bytes")
	}

	//if string(buf[lbt+9:lbt+29]) != a.infoHash {
	//	return errors.New("error infoHash bytes")
	//}

	return nil
}

func (a *Acquirer) sendExtHandshake() error {
	message := Message{}
	message.ID = MsgExtended

	msg := map[string]interface{}{
		"m": map[string]interface{}{
			"ut_metadata": 1,
		},
	}

	message.Payload = append(message.Payload, byte(ExMsgRequest))
	message.Payload = bencode.Encode(msg)
	data := message.Serialize()

	if err := a.conn.SetWriteDeadline(time.Now().Add(time.Second + 15)); err != nil {
		return err
	}

	n, err := a.conn.Write(data)
	if err != nil {
		return err
	}

	logger.Println("[Acquirer] sendExtHandshake done : ", string(data[:n]))
	return nil
}

func (a *Acquirer) readMessage() error {
	message, err := Read(a.conn)
	if err != nil {
		logger.Println("[Read] ", err.Error())
		return err
	}

	logger.Println("[Acquirer] readMessage done : %v", string(message.Payload))

	//for {
	//	switch message.ID {
	//	case MsgExtended:
	//	default:
	//
	//	}
	//}
	return nil
}
