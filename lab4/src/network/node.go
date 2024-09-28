package network

import (
	"google.golang.org/protobuf/proto"
	pb "lab4/generated"
	"log"
	"net"
)

func (c *MessageContext) getMessageId() int64 {
	id := c.Id
	c.Id++
	return id
}

func sendGameMessage(conn *net.UDPConn, m proto.Message) {
	data, err := proto.Marshal(m)
	if err != nil {
		log.Fatal("Marshaling error: ", err)
	}

	log.Println("Sending Message to multicast ...")
	_, err = conn.Write(data)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Sending Message to multicast ... Done!")
}

func listenGameMessages(node NodeMethods, conn *net.UDPConn) {
	buf := make([]byte, MessageBufferSize)
	for {
		bytes, addr, err := conn.ReadFromUDP(buf)
		log.Println("Got a message!")
		if err != nil {
			log.Fatal(err)
		}

		gameMessage := &pb.GameMessage{}
		err = proto.Unmarshal(buf[:bytes], gameMessage)
		if err != nil {
			log.Fatal(err)
		}

		node.processSenderAddress(addr)
		node.processMsg(gameMessage)
	}
}
