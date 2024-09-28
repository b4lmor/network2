package network

import (
	"fmt"
	pb "lab4/generated"
	"log"
	"net"
	"time"
)

func (s *NodeSlave) ListenMulticastAnnouncementMsg() {
	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", MulticastAddress, MulticastPort))
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		log.Fatal(err)
	}

	listenGameMessages(s, conn)
}

func (s *NodeSlave) processMsg(msg *pb.GameMessage) {
	s.MasterId = *msg.SenderId
	if msg.GetType() != nil {
		switch t := msg.GetType().(type) {
		case *pb.GameMessage_Ping:
			log.Println("Received Ping Message")
		case *pb.GameMessage_Steer:
			fmt.Printf("Steer Direction: %v\n", t.Steer.Direction)
		case *pb.GameMessage_Ack:
			log.Println("Received Ack Message")
		case *pb.GameMessage_State:
			log.Println("Received State Message")
		case *pb.GameMessage_Announcement:
			log.Println("Received Announcement Message")
			s.processAnnouncementMsg(t)
		case *pb.GameMessage_Join:
			log.Println("Received Join Message")
		case *pb.GameMessage_Error:
			fmt.Printf("Error Message: %s\n", *t.Error.ErrorMessage)
		case *pb.GameMessage_RoleChange:
			log.Println("Received Role Change Message")
		case *pb.GameMessage_Discover:
			log.Println("Received Discover Message")
		default:
			log.Println("Unknown message type")
		}
	} else {
		log.Println("Message type is nil")
	}
}

func (s *NodeSlave) processAnnouncementMsg(msg *pb.GameMessage_Announcement) {
	log.Println("[SLAVE] :: ", msg.Announcement)
}

func (s *NodeSlave) processSenderAddress(addr *net.UDPAddr) {
	s.Master.Address = addr.IP.String()
	s.Master.Port = addr.Port
	s.MasterLastActiveAt = time.Now()
}
