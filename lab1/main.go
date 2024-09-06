package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"time"
	"utils"
)

const (
	PORT  = 10000
	DELAY = 2 * time.Second
	TTL   = 3 * time.Second
)

var (
	aliveAddresses = utils.NewSafeMap[string, time.Time]()
	connections    int
)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Usage: lab1 [address]")
		return
	}

	multicastAddr := os.Args[1]

	addr, err := resolveUDPAddr(multicastAddr)
	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := listenUDPConn(addr)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	go receiveMessages(conn)
	go monitorConnections()
	sendMessages(multicastAddr)
}

func sendMessages(multicastAddr string) {
	mAddr, err := resolveUDPAddr(multicastAddr)
	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := dialUDPConn(mAddr)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
		return
	}

	ipAddr, err := exec.Command("hostname", "-I").Output()
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		msg := fmt.Sprintf("Hello from %s %s", hostname, string(ipAddr))
		_, err = conn.Write([]byte(msg))
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(DELAY)
	}
}

func receiveMessages(conn *net.UDPConn) {
	buf := make([]byte, 1024)
	for {
		_, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println(err)
		}
		aliveAddresses.Set(addr.String(), time.Now())
		currentConnections := aliveAddresses.Len()
		if connections != currentConnections {
			connections = currentConnections
			printAliveAddresses()
		}
	}
}

func monitorConnections() {
	for {
		time.Sleep(DELAY)
		for _, addr := range aliveAddresses.GetAllKeys() {
			lastSeen, _ := aliveAddresses.Get(addr)
			if time.Since(lastSeen) > TTL {
				aliveAddresses.Delete(addr)
				connections = aliveAddresses.Len()
				printAliveAddresses()
			}
		}
	}
}

func printAliveAddresses() {
	fmt.Println("Alive copies:")
	for _, addr := range aliveAddresses.GetAllKeys() {
		fmt.Printf("\t\t%s\n", addr)
	}
}

func isIpv4(multicastAddr string) bool {
	return net.ParseIP(multicastAddr).To4() != nil
}

func resolveUDPAddr(multicastAddr string) (*net.UDPAddr, error) {
	if isIpv4(multicastAddr) {
		return net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", multicastAddr, PORT))
	} else {
		return net.ResolveUDPAddr("udp6", fmt.Sprintf("[%s]:%d", multicastAddr, PORT))
	}
}

func dialUDPConn(multicastAddr *net.UDPAddr) (*net.UDPConn, error) {
	if isIpv4(multicastAddr.IP.String()) {
		return net.DialUDP("udp4", nil, multicastAddr)
	} else {
		return net.DialUDP("udp6", nil, multicastAddr)
	}
}

func listenUDPConn(multicastAddr *net.UDPAddr) (*net.UDPConn, error) {
	if isIpv4(multicastAddr.IP.String()) {
		return net.ListenMulticastUDP("udp4", nil, multicastAddr)
	} else {
		return net.ListenMulticastUDP("udp6", nil, multicastAddr)
	}
}
