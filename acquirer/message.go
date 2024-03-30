package acquirer

import (
	"encoding/binary"
	"io"
)

type messageID uint8

const (
	MsgChoke         messageID = 0
	MsgUnchoke       messageID = 1
	MsgInterested    messageID = 2
	MsgNotInterested messageID = 3
	MsgHave          messageID = 4
	MsgBitfield      messageID = 5
	MsgRequest       messageID = 6
	MsgPiece         messageID = 7
	MsgCancel        messageID = 8
	MsgExtended      messageID = 20

	BitTorrentProtocol = "BitTorrent protocol"
)

type Handshake struct {
	Prefix   string
	InfoHash []byte
	PeerId   []byte
}

func (h *Handshake) Serialize() []byte {
	const firstByte = byte(0x13)
	BitTorrent := []byte(BitTorrentProtocol)
	ReservedBytes := make([]byte, 8)
	data := make([]byte, 0, len(BitTorrent)+49)
	data = append(data, firstByte)
	data = append(data, BitTorrent...)
	data = append(data, ReservedBytes...)
	data = append(data, h.InfoHash...)
	data = append(data, h.PeerId...)
	return data
}

// Message stores ID and payload of a message
type Message struct {
	ID      messageID
	Payload []byte
}

func (m *Message) Serialize() []byte {
	if m == nil {
		return make([]byte, 4)
	}
	length := uint32(len(m.Payload) + 1) // +1 for id
	buf := make([]byte, 4+length)
	binary.BigEndian.PutUint32(buf[0:4], length)
	buf[4] = byte(m.ID)
	copy(buf[5:], m.Payload)
	return buf
}

func Read(r io.Reader) (*Message, error) {
	lengthBuf := make([]byte, 4)
	_, err := io.ReadFull(r, lengthBuf)
	if err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(lengthBuf)

	// keep-alive message
	if length == 0 {
		return nil, nil
	}

	messageBuf := make([]byte, length)
	_, err = io.ReadFull(r, messageBuf)
	if err != nil {
		return nil, err
	}

	m := Message{
		ID:      messageID(messageBuf[0]),
		Payload: messageBuf[1:],
	}

	return &m, nil
}
