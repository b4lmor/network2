package proxy

import (
	"log"
	"net"
)

func Start(port int) {
	listener, err := net.ListenTCP("tcp4", &net.TCPAddr{Port: port})
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	log.Printf("Proxy server is listening on %v\n", listener.Addr())
	for {
		client, err := listener.AcceptTCP()
		if err != nil {
			log.Println(err.Error())
			continue
		}
		log.Printf("Accepted client: %v\n", client.RemoteAddr())
		go handleClient(client)
	}
}

func handleClient(client *net.TCPConn) {
	defer func() {
		client.Close()
		log.Printf("%v: Client connection is closed\n", client.RemoteAddr())
	}()
	// Authenticate client
	authResponse, err := authenticate(client)
	if err != nil {
		log.Printf("%v: %v\n", client.RemoteAddr(), err)
		err = sendAuthReply(client, authResponse)
		if err != nil {
			log.Printf("%v: %v\n", client.RemoteAddr(), err)
		}
		return
	}
	// Send auth reply
	err = sendAuthReply(client, authResponse)
	if err != nil {
		log.Printf("%v: %v\n", client.RemoteAddr(), err)
		return
	}
	log.Printf("%v: Client is authenticated\n", client.RemoteAddr())
	// Execute command
	peer, commandReply, err := connectCommand(client)
	if err != nil {
		log.Printf("%v: %v\n", client.RemoteAddr(), err)
		err = sendCommandReply(client, commandReply)
		if err != nil {
			log.Printf("%v: %v\n", client.RemoteAddr(), err)
		}
		return
	}
	// Send command reply
	err = sendCommandReply(client, commandReply)
	if err != nil {
		log.Printf("%v: %v\n", client.RemoteAddr(), err)
		return
	}
	log.Printf("%v: Proxy server is connected to peer\n", client.RemoteAddr())
	// Transfer data
	transferData(client, peer)
}
