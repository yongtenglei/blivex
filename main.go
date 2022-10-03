package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

var startUrlFormat = "https://api.live.bilibili.com/xlive/web-room/v1/index/getDanmuInfo?id=%d"
var startInfo StartInfo

func main() {
	biliClient := NewBiliClient(22495868)
	hostList := biliClient.GetHostList()
	biliClient.Connect(hostList)
	defer biliClient.Disconnect()

	go biliClient.ReceiveMessages()
	go biliClient.ParseMessages()
	biliClient.HeartBeat()
}

type BiliClient struct {
	RoomID     uint32
	HTTPClient *http.Client
	Conn       *websocket.Conn

	Ch chan PacketBody
}

func NewBiliClient(roomID uint32) *BiliClient {
	return &BiliClient{
		RoomID:     roomID,
		HTTPClient: &http.Client{},
		Conn:       nil,

		Ch: make(chan PacketBody, 1024),
	}
}

type StartInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    struct {
		Group            string  `json:"group"`
		BusinessID       int     `json:"business_id"`
		RefreshRowFactor float64 `json:"refresh_row_factor"`
		RefreshRate      int     `json:"refresh_rate"`
		MaxDelay         int     `json:"max_delay"`
		Token            string  `json:"token"`
		HostList         []struct {
			Host    string `json:"host"`
			Port    int    `json:"port"`
			WssPort int    `json:"wss_port"`
			WsPort  int    `json:"ws_port"`
		} `json:"host_list"`
	} `json:"data"`
}

func (bc *BiliClient) GetHostList() []url.URL {
	req, err := http.NewRequest("GET", fmt.Sprintf(startUrlFormat, bc.RoomID), nil)
	if err != nil {
		panic(err)
	}

	resp, err := bc.HTTPClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	bts, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bts, &startInfo)
	if err != nil {
		panic(err)
	}

	hostList := make([]url.URL, 0, 3)

	for _, host := range startInfo.Data.HostList {
		u := url.URL{Scheme: "wss", Host: fmt.Sprintf("%s:%d", host.Host, host.WssPort), Path: "/sub"}
		hostList = append(hostList, u)
	}

	return hostList
}

func (bc *BiliClient) Connect(urls []url.URL) {
	var websocketErr error = errors.New("websocket connection error")

	for _, url := range urls {
		bc.Conn, _, websocketErr = websocket.DefaultDialer.Dial(url.String(), nil)
		if websocketErr != nil {
			fmt.Printf("websocket dial %s failed\n", url.String())
			continue
		}
		fmt.Printf("websocket dial %s successfully\n", url.String())
		break
	}

	if websocketErr != nil {
		panic(websocketErr)
	}

	type handShakeInfo struct {
		UID       uint8  `json:"uid"`
		Roomid    uint32 `json:"roomid"`
		Protover  uint8  `json:"protover"`
		Platform  string `json:"platform"`
		Clientver string `json:"clientver"`
		Type      uint8  `json:"type"`
		Key       string `json:"key"`
	}

	hsInfo := handShakeInfo{
		UID:       0,
		Roomid:    bc.RoomID,
		Protover:  2,
		Platform:  "web",
		Clientver: "1.10.2",
		Type:      2,
		Key:       startInfo.Data.Token,
	}

	b, err := json.Marshal(hsInfo)
	if err != nil {
		panic(err)
	}

	// 把json 发送给服务器
	err = bc.SendPacket(0, 16, 0, 7, 0, b)
	if err != nil {
		panic(err)
	}

	fmt.Printf("connect to room %d successfully\n", bc.RoomID)
}

func (bc *BiliClient) HeartBeat() {
	for {
		err := bc.SendPacket(0, 16, 0, 2, 0, []byte(""))
		if err != nil {
			fmt.Println("HeartBeat err: ", err)
			time.Sleep(5 * time.Second)
			continue
		}
		time.Sleep(20 * time.Second)
	}
}

func (bc *BiliClient) Disconnect() {
	err := bc.Conn.Close()
	if err != nil {
		fmt.Println(err)
	}
}
