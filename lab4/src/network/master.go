package network

import (
	"fmt"
	pb "lab4/generated"
	"log"
	"net"
	"time"
)

func (m *NodeMaster) processMsg(msg *pb.GameMessage) {
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
			m.processAnnouncementMsg(t)
		case *pb.GameMessage_Join:
			log.Println("Received Join Message")
		case *pb.GameMessage_Error:
			fmt.Printf("Error Message: %m\n", *t.Error.ErrorMessage)
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

func (m *NodeMaster) processAnnouncementMsg(msg *pb.GameMessage_Announcement) {
	log.Println(msg.Announcement)
}

func (m *NodeMaster) SendMulticastAnnouncementMsg() {
	mAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", MulticastAddress, MulticastPort))
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp4", nil, mAddr)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	for {
		sendGameMessage(conn, m.getAnnouncementMsg())
		time.Sleep(MulticastDelay)
	}
}

func (m *NodeMaster) getAnnouncementMsg() *pb.GameMessage { // TODO: Delete
	width, height, foodStatic, stateDelayMs := int32(50), int32(50), int32(10), int32(5)

	gameConfig := &pb.GameConfig{
		Width:        &width,
		Height:       &height,
		FoodStatic:   &foodStatic,
		StateDelayMs: &stateDelayMs,
	}

	gamePlayers := &pb.GamePlayers{
		Players: make([]*pb.GamePlayer, 0),
	}

	gameName := "The snake game"

	gameAnnouncement := &pb.GameAnnouncement{
		GameName: &gameName,
		Config:   gameConfig,
		Players:  gamePlayers,
	}

	announcements := []*pb.GameAnnouncement{gameAnnouncement}

	msgSeq := m.getMessageId()

	return &pb.GameMessage{
		MsgSeq: &msgSeq,
		Type: &pb.GameMessage_Announcement{
			Announcement: &pb.GameMessage_AnnouncementMsg{
				Games: announcements,
			},
		},
	}
}
