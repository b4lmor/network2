package network

import (
	"google.golang.org/protobuf/proto"
	"lab4/context"
	pb "lab4/generated/external_api"
	"log"
	"net"
	"sync"
	"time"
)

const (
	MulticastSocket   = "239.192.0.4:9192"
	MulticastDelay    = 1 * time.Second
	MessageBufferSize = 1 << 11
)

type NodeMaster struct {
	ctx  *context.Context
	conn *net.UDPConn
	mu   sync.Mutex
}

func (nm *NodeMaster) initConnection() {
	laddr, err := net.ResolveUDPAddr("udp", nm.ctx.GameSocket)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.ListenUDP("udp4", laddr)
	if err != nil {
		log.Fatal(err)
	}
	nm.conn = conn
}

func NewNodeMaster(ctx *context.Context) *NodeMaster {
	nm := NodeMaster{ctx: ctx}
	nm.initConnection()
	return &nm
}

func (nm *NodeMaster) StartListenUnicastMessages() {
	buf := make([]byte, MessageBufferSize)
	for {
		bytes, raddr, err := nm.conn.ReadFromUDP(buf)
		if err != nil {
			log.Println(err)
			continue
		}
		gameMessage := &pb.GameMessage{}
		err = proto.Unmarshal(buf[:bytes], gameMessage)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("Got a message! address: ", raddr.String())
		nm.consumeGameMessage(gameMessage, raddr)
	}
}

func (nm *NodeMaster) StartListenMulticastAnnouncement() {
	log.Println("Start listening multicast ...")
	laddr, err := net.ResolveUDPAddr("udp4", MulticastSocket)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.ListenMulticastUDP("udp4", nil, laddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	buf := make([]byte, MessageBufferSize)
	for {
		bytes, raddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Println(err)
			continue
		}
		if err != nil {
			log.Fatal(err)
		}
		gameMessage := &pb.GameMessage{}
		err = proto.Unmarshal(buf[:bytes], gameMessage)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("Got a multicast message! Sender's address: ", raddr.String())
		nm.consumeGameMessage(gameMessage, raddr)
	}
}

func (nm *NodeMaster) StartSendMulticastAnnouncementMsg() {
	for {
		time.Sleep(MulticastDelay)
		if nm.ctx.CurrentGame != nil {
			nm.sendAnnouncement(-1, MulticastSocket)
		}
	}
}

func (nm *NodeMaster) sendGameMessage(socket string, m proto.Message) {
	raddr, err := net.ResolveUDPAddr("udp4", socket)
	if err != nil {
		log.Fatal(err)
	}
	data, err := proto.Marshal(m)
	if err != nil {
		log.Fatal("Marshaling error: ", err)
	}
	log.Println("Sending Message to ", raddr.String())
	_, err = nm.conn.WriteToUDP(data, raddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Sending Message ... Done!")
}

func newGameMessage(seq int64, senderId, receiverId int32) *pb.GameMessage {
	return &pb.GameMessage{MsgSeq: &seq, SenderId: &senderId, ReceiverId: &receiverId}
}
