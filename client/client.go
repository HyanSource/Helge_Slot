package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"runtime"
	"strconv"
	"time"

	"github.com/HyanSource/Helge_Slot/pb"
	"github.com/golang/protobuf/proto"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	c := NewTcpClient("220.142.99.234", 8124)
	c.Start()
	select {}
}

type Message struct {
	Len   uint32
	MsgId uint32
	Data  []byte
}

//結果
type Result struct {
	Table         []int32
	Paylinesnum   []int32
	Paylinescount []int32
	Odds          int32
}

//初始化client
func NewTcpClient(ip string, port int) *TcpClient {

	addstr := ip + ":" + strconv.Itoa(port)
	conn, err := net.Dial("tcp", addstr)

	if err != nil {
		panic("net dial err:" + err.Error())
	}

	return &TcpClient{
		conn:     conn,
		isOnline: make(chan bool),
	}
}

type TcpClient struct {
	conn     net.Conn
	Id       int32
	isOnline chan bool
}

func (t *TcpClient) Start() {
	//read goroutine
	go func() {
		defer func() { t.isOnline <- false }()

		for {
			//讀取datalen
			headdata := make([]byte, 4)
			if _, err := io.ReadFull(t.conn, headdata); err != nil {
				fmt.Println(err)
				return
			}

			datalen, err := t.Unpack(headdata)

			if err != nil {
				fmt.Println("datalen unpack err:", err)
				return
			}

			headmsgid := make([]byte, 4)
			if _, err := io.ReadFull(t.conn, headmsgid); err != nil {
				fmt.Println(err)
				return
			}

			msgid, err := t.Unpack(headmsgid)

			if err != nil {
				fmt.Println("msgid unpack err:", err)
				return
			}

			bodydata := make([]byte, datalen)

			if _, err := io.ReadFull(t.conn, bodydata); err != nil {
				fmt.Println("bodydata read err:", err)
				return
			}

			m := &Message{
				Len:   uint32(len(bodydata)),
				MsgId: msgid,
				Data:  bodydata,
			}

			fmt.Println(m)

			t.DoMsg(m)
		}
	}()

	//write goroutine
	go func() {
		for {

			t.SendMsg(100, &pb.Play{Bet: 10})
			time.Sleep(10 * time.Second)
		}
	}()
}

//處理收到Message的業務
func (t *TcpClient) DoMsg(msg *Message) {
	r := &pb.Result{}
	err := proto.Unmarshal(msg.Data, r)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(r)
}

func (t *TcpClient) SendMsg(msgid uint32, data proto.Message) {
	//Marshal
	binarydata, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("marshal err:", err)
		return
	}

	senddata, err := t.Pack(msgid, binarydata)
	if err != nil {
		fmt.Println("pack err:", err)
		return
	}

	if _, err := t.conn.Write(senddata); err != nil {
		fmt.Println("conn write err:", err)
		return
	}

}

//解出datalen and msgid
func (t *TcpClient) Unpack(headdata []byte) (data uint32, err error) {

	databuff := bytes.NewBuffer(headdata)

	if err = binary.Read(databuff, binary.LittleEndian, &data); err != nil {
		return 0, err
	}

	return data, nil
}

func (t *TcpClient) Pack(msgid uint32, dataBytes []byte) (out []byte, err error) {

	outbuff := bytes.NewBuffer([]byte{})

	//寫Len
	if err = binary.Write(outbuff, binary.LittleEndian, uint32(len(dataBytes))); err != nil {
		return nil, err
	}

	//寫MsgId
	if err = binary.Write(outbuff, binary.LittleEndian, msgid); err != nil {
		return nil, err
	}

	if err = binary.Write(outbuff, binary.LittleEndian, dataBytes); err != nil {
		return nil, err
	}

	out = outbuff.Bytes()

	return out, nil
}
