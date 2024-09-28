package network

import (
	pb "lab4/generated"
	"net"
	"time"
	"utils"
)

const (
	MulticastAddress  = "239.192.0.4"
	MulticastPort     = 9192
	MulticastDelay    = 1 * time.Second
	MessageBufferSize = 1024
)

type MessageContext struct {
	Id int64
}

type Socket struct {
	Address string
	Port    int
}

type NodeMethods interface {
	processSenderAddress(addr *net.UDPAddr)
	processMsg(msg *pb.GameMessage)
	processAnnouncementMsg(msg *pb.GameMessage_Announcement)
}

type Node struct {
	MessageContext
	Id     int32
	socket Socket
}

type NodeMaster struct {
	Node
	Slaves utils.SafeMap[Node, time.Time]
}

type NodeSlave struct {
	Node
	MasterId           int32
	Master             Socket
	MasterLastActiveAt time.Time
}
