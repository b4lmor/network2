package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"
	"utils"
)

const (
	UploadDir  = "./assets"
	Port       = "8000"
	Delay      = 1 * time.Second
	BufferSize = 4096
)

type ClientInfo struct {
	Addr          net.Addr
	TotalBytes    int64
	StartTime     time.Time
	ReceivedBytes int64
}

var clients = utils.NewSafeMap[net.Addr, *ClientInfo]()

func main() {
	startServer(Port)
}

func startServer(port string) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer ln.Close()
	fmt.Println("Listening on " + port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	addr := conn.RemoteAddr()
	fmt.Printf("Received connection from %s\n", addr)
	clients.Set(addr, &ClientInfo{Addr: addr, StartTime: time.Now()})
	defer clients.Delete(addr)

	fileNameSize, err := readInt64(conn)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return
	}

	fileNameBytes := make([]byte, fileNameSize)
	n, err := conn.Read(fileNameBytes)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		return
	}
	fileName := string(fileNameBytes[:n])

	fileSize, err := readInt64(conn)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		return
	}

	clientInfo, _ := clients.Get(addr)
	clientInfo.TotalBytes = fileSize

	if _, err = os.Stat(UploadDir); os.IsNotExist(err) {
		os.Mkdir(UploadDir, os.ModePerm)
	}

	filePath := filepath.Join(UploadDir, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error creating file: %s\n", err)
		return
	}
	defer file.Close()

	go speedMonitor(addr)

	buffer := make([]byte, BufferSize)
	for {
		if clientInfo.ReceivedBytes >= fileSize {
			break
		}
		n, err = conn.Read(buffer)
		if err != nil {
			fmt.Printf("Error reading file: %s\n", err)
			break
		}

		file.Write(buffer[:n])
		clientInfo.ReceivedBytes += int64(n)
	}

	stat, _ := file.Stat()
	if stat.Size() == fileSize {
		conn.Write([]byte("SUCCESS"))
		fmt.Println("SUCCESS")
	} else {
		conn.Write([]byte("FAIL"))
		fmt.Println("FAIL")
	}
}

func speedMonitor(addr net.Addr) {
	currentReceived := int64(0)
	for {
		time.Sleep(Delay)
		if info, exists := clients.Get(addr); exists {
			elapsed := time.Since(info.StartTime).Seconds()
			if elapsed > 0 {
				instantSpeed := float64(info.ReceivedBytes-currentReceived) / elapsed
				fmt.Printf("Скорость для     %s: %.2f байт/с\n", addr, instantSpeed)
				fmt.Printf("Общий размер для %s: %d / %d байт\n", addr, info.ReceivedBytes, info.TotalBytes)
				info.StartTime = time.Now()
				currentReceived = info.ReceivedBytes
			}
		}
	}
}

func readInt64(conn net.Conn) (int64, error) {
	sizeBytes := make([]byte, 8)
	_, err := conn.Read(sizeBytes)
	fileNameSize := int64(0)
	for _, b := range sizeBytes {
		fileNameSize = fileNameSize<<8 + int64(b)
	}
	return fileNameSize, err
}
