package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: client <filepath> <address>")
		return
	}
	filePath := os.Args[1]
	serverAddr := os.Args[2]
	sendFile(filePath, serverAddr)
}

func sendFile(filePath, serverAddr string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("open file error:", err)
		return
	}
	defer file.Close()

	fileName := filepath.Base(filePath)
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("get file info error:", err)
		return
	}

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("dial error:", err)
		return
	}
	defer conn.Close()

	_, err = sendInt64(int64(len(fileName)), conn)
	if err != nil {
		fmt.Println("send error:", err)
		return
	}

	_, err = conn.Write([]byte(fileName))
	if err != nil {
		fmt.Println("send file name error:", err)
		return
	}

	_, err = sendInt64(fileInfo.Size(), conn)
	if err != nil {
		fmt.Println("send file size error:", err)
		return
	}

	buffer := make([]byte, 4096)
	for {
		n, err := file.Read(buffer)
		if n == 0 {
			break
		} else if err != nil {
			fmt.Println("read file error:", err)
			return
		}
		_, err = conn.Write(buffer[:n])
		if err != nil {
			fmt.Println("send file size error:", err)
			return
		}
	}

	response := make([]byte, 8)
	n, err := conn.Read(response)
	if err != nil {
		fmt.Println("read response error:", err)
		return
	}

	if string(response[:n]) == "SUCCESS" {
		fmt.Println("SUCCESS")
	} else {
		fmt.Println("FAIL")
	}
}

func sendInt64(size int64, conn net.Conn) (int, error) {
	sizeBytes := make([]byte, 8)
	for i := 0; i < 8; i++ {
		sizeBytes[7-i] = byte(size >> (8 * i))
	}
	return conn.Write(sizeBytes)
}
