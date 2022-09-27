package main

import (
	"bytes"
	"encoding/binary"

	"github.com/gorilla/websocket"
)

type Packet struct {
	PacketLen  int
	HeaderLen  int
	Version    int
	Operation  int
	SequenceID int

	Body []PacketBody
}

type PacketBody struct {
	Cmd     string                 `json:"cmd"`
	Data    map[string]interface{} `json:"data"`
	MsgSelf string                 `json:"msg_self"`
	Info    []interface{}          `json:"info"`
}

func (bc *BiliClient) SendPacket(packetLen uint32, headerLen uint16, version uint16, operation uint32, sequenceID uint32, body []byte) error {
	if packetLen == 0 {
		packetLen = uint32(len(body) + 16)
	}

	header := new(bytes.Buffer)

	var data = []interface{}{
		packetLen,
		headerLen,
		version,
		operation,
		sequenceID,
	}

	for _, v := range data {
		err := binary.Write(header, binary.BigEndian, v)
		if err != nil {
			return err
		}
	}

	socketData := append(header.Bytes(), body...)

	err := bc.Conn.WriteMessage(websocket.TextMessage, socketData)

	return err

}
