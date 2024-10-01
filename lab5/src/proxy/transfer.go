package proxy

import (
	"io"
	"log"
	"net"
	"sync"
)

func transferData(client *net.TCPConn, peer *net.TCPConn) {
	var wg sync.WaitGroup
	wg.Add(2)
	go copyData(client, peer, &wg)
	go copyData(peer, client, &wg)
	wg.Wait()
}

func copyData(dest *net.TCPConn, src *net.TCPConn, wg *sync.WaitGroup) {
	defer wg.Done()
	defer dest.Close()
	_, err := io.Copy(dest, src)
	if err != nil {
		log.Printf("%v: Reading error: %v\n", dest.RemoteAddr(), err)
	}
}
