package acquirer

import (
	"bittorrent/bencode"
	"bittorrent/common"
	"bittorrent/config"
	"bittorrent/logger"
	"bittorrent/utils"
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"os"
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
			//fmt.Println("[AcquireManager] run handle")
			if info := m.Pop(); info != nil {
				go m.handle(info)
			}
		}
	}
}

func handle(info *PeerInfo) {
	row := DataStorage.Get(info.InfoHash)
	if row != nil {
		return
	}

	fmt.Println("[handle]", hex.EncodeToString([]byte(info.InfoHash)), " ", info.Ip, " ", info.Port)
	acquirer, err := NewAcquirer(info.InfoHash, info.Ip, info.Port)
	if err != nil {
		return
	}
	defer acquirer.close()
	if err = acquirer.sendHandshake(); err != nil {
		fmt.Println("[handle] sendHandshake", err.Error())
		logger.Println("[handle] sendHandshake", err.Error())
		return
	}
	if err = acquirer.readHandshake(); err != nil {
		fmt.Println("[handle] readHandshake", err.Error())
		logger.Println("[handle] readHandshake", err.Error())
		return
	}
	if err = acquirer.sendExtHandshake(); err != nil {
		fmt.Println("[handle] sendExtHandshake", err.Error())
		logger.Println("[handle] sendExtHandshake", err.Error())
		return
	}
	if err = acquirer.readMessage(); err != nil {
		fmt.Println("[handle] readMessage", err.Error())
		logger.Println("[handle] readMessage", err.Error())
		return
	}
}

type Acquirer struct {
	conn     net.Conn
	done     chan struct{}
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
	if err = conn.(*net.TCPConn).SetReadBuffer(1024 * 1024 * 1); err != nil {
		logger.Println("[Acquirer] tcp SetReadBuffer err: %v", err.Error())
		return nil, err
	}

	acquirer := &Acquirer{
		conn:     conn,
		infoHash: infoHash,
		done:     make(chan struct{}),
		peerId:   utils.RandomID(),
	}

	return acquirer, nil
}

func (a *Acquirer) close() {
	a.done <- struct{}{}
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
		logger.Println("[sendHandshake] SetWriteDeadline", err.Error())
		return err
	}

	n, err := a.conn.Write(data)
	if err != nil {
		logger.Println("[sendHandshake] Write", err.Error())
		return err
	}

	logger.Println("[Acquirer] sendHandshake done ", n)
	return nil
}

func (a *Acquirer) readHandshake() error {
	_, err := ReadHandshake(a.conn)
	return err
}

func (a *Acquirer) sendExtHandshake() error {
	message := Message{}
	message.ID = MsgExtended

	msg := map[string]interface{}{
		"m": map[string]interface{}{
			"ut_metadata": 1,
		},
	}

	message.Payload = append(message.Payload, byte(0))
	message.Payload = append(message.Payload, bencode.Encode(msg)...)
	data := message.Serialize()

	if err := a.conn.SetWriteDeadline(time.Now().Add(time.Second + 15)); err != nil {
		logger.Println("[sendExtHandshake] SetWriteDeadline", err.Error())
		return err
	}

	n, err := a.conn.Write(data)
	if err != nil {
		logger.Println("[sendExtHandshake] Write", err.Error())
		return err
	}

	logger.Println("[Acquirer] sendExtHandshake done : ", string(data[:n]), data[:n])
	return nil
}

func (a *Acquirer) readMessage() error {

	var metadataSize int64 = 0
	var piecesNum int64 = 0
	var pieces [][]byte

	for {
		if err := a.conn.SetWriteDeadline(time.Now().Add(time.Second + 15)); err != nil {
			return err
		}

		message, err := ReadMessage(a.conn)
		if err != nil {
			return err
		}

		switch message.ID {
		case MsgExtended:
			buf := bytes.NewBuffer(message.Payload)
			extendedID, err := buf.ReadByte()
			if err != nil {
				return err
			}
			switch extendedID {
			case 0:
				decode, err := bencode.Decode(buf)
				if err != nil {
					return err
				}
				d := decode.(map[string]interface{})
				metadataSize = d["metadata_size"].(int64)
				m := d["m"].(map[string]interface{})
				utMetadata := m["ut_metadata"].(int64)
				piecesNum = metadataSize / BlockSize
				if metadataSize%BlockSize != 0 {
					piecesNum++
				}
				pieces = make([][]byte, piecesNum)
				go a.sendRequestPieces(utMetadata, piecesNum)

			case 1:
				decode, err := bencode.Decode(buf)
				if err != nil {
					return err
				}
				d := decode.(map[string]interface{})
				msgType := d["msg_type"].(int64)
				piece := d["piece"].(int64)
				totalSize := d["total_size"].(int64)
				if msgType != ExMsgData || totalSize != metadataSize {
					return errors.New("[readMessage] error data")
				}

				l := len(bencode.Encode(decode)) + 1
				pieces[piece] = message.Payload[l:]

				if piece+1 == piecesNum {
					logger.Println("[readMessage] start")
					metadataInfo := bytes.Join(pieces, nil)

					logger.Println("[readMessage] metadataInfo ", string(metadataInfo))

					info := sha1.Sum(metadataInfo)
					if !bytes.Equal([]byte(a.infoHash), info[:]) {
						logger.Println("[readMessage] infoHash err")
						return nil
					}

					writeToFile(hex.EncodeToString([]byte(a.infoHash)), metadataInfo)
					DataStorage.Put(&DBMetaInfo{InfoHash: a.infoHash})
					//if metadataInfo, err := bencode.Decode(bytes.NewBuffer(metadataInfo)); err == nil {
					//	fmt.Println(metadataInfo)
					//}

					logger.Println("[readMessage] done")
					return nil
				}

			default:
				continue
			}
		default:
			continue
		}
	}

}

func (a *Acquirer) sendRequestPieces(utMetadata int64, piecesNum int64) {
	for i := 0; i < int(piecesNum); i++ {
		reqByte := bencode.Encode(map[string]interface{}{
			"msg_type": ExMsgRequest,
			"piece":    i,
		})

		msg := Message{}
		msg.ID = MsgExtended
		msg.Payload = append(msg.Payload, byte(utMetadata))
		msg.Payload = append(msg.Payload, reqByte...)
		data := msg.Serialize()

		if err := a.conn.SetWriteDeadline(time.Now().Add(time.Second * 15)); err != nil {
			logger.Println("[sendRequestPieces] err", err.Error())
			break
		}
		_, err := a.conn.Write(data)
		if err != nil {
			logger.Println("[sendRequestPieces] err", err.Error())
			break
		}
		logger.Println("[sendRequestPieces] done ", string(data), data)
	}
}

func writeToFile(infoHash string, content []byte) {
	file, err := os.Create(infoHash + ".torrent")
	if err != nil {
		logger.Println("[writeToFile]", err.Error())
		return
	}

	_, err = file.Write(content)
	if err != nil {
		logger.Println("[writeToFile]", err.Error())
		return
	}

	if err := file.Sync(); err != nil {
		logger.Println("[writeToFile]", err.Error())
	}
	if err := file.Close(); err != nil {
		logger.Println("[writeToFile]", err.Error())
	}
}
