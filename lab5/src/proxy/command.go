package proxy

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

func connectCommand(client *net.TCPConn) (*net.TCPConn, byte, error) {
	// Accepting data
	const RequestSize = 4
	request := make([]byte, RequestSize)
	n, err := io.ReadFull(client, request)
	if n != RequestSize {
		return nil, SOCKS_REPLY_CONNECTION_NOT_ALLOWED_BY_RULESET, errors.New("bad request")
	} else if err != nil {
		return nil, SOCKS_REPLY_CONNECTION_NOT_ALLOWED_BY_RULESET, err
	}
	// Check version
	version := request[0]
	if version != SOCKS_VERSION {
		return nil, SOCKS_REPLY_CONNECTION_NOT_ALLOWED_BY_RULESET, errors.New(fmt.Sprintf(
			"Socks version %v is expected, but not %v", SOCKS_VERSION, version))
	}
	// Check command
	command := request[1]
	if command != SOCKS_CMD_CONNECT {
		return nil, SOCKS_REPLY_COMMAND_NOT_SUPPORTED, errors.New(fmt.Sprintf(
			"Unsupported command %v, %v command supported", command, SOCKS_CMD_CONNECT))
	}
	// Check reserved byte
	reservedByte := request[2]
	if reservedByte != SOCKS_RESERVED_BYTE {
		return nil, SOCKS_REPLY_CONNECTION_NOT_ALLOWED_BY_RULESET, errors.New(fmt.Sprintf(
			"Reserved byte must be set to %v, but not %v", SOCKS_RESERVED_BYTE, reservedByte))
	}
	// Check address type
	addressType := request[3]
	switch addressType {
	case SOCKS_ADDR_TYPE_IPV4:
		ipv4, port, err := readIpv4AndPort(client)
		if err != nil {
			return nil, SOCKS_REPLY_CONNECTION_NOT_ALLOWED_BY_RULESET, err
		}
		return ipv4Connect(ipv4, port)
	case SOCKS_ADDR_TYPE_FQDN:
		domainName, port, err := readDomainNameAndPort(client)
		if err != nil {
			return nil, SOCKS_REPLY_CONNECTION_NOT_ALLOWED_BY_RULESET, err
		}
		return domainNameConnect(domainName, port, client)
	default:
		return nil, SOCKS_REPLY_ADDRESS_TYPE_NOT_SUPPORTED, errors.New(fmt.Sprintf(
			"Unsupported address type %v, %v is supported", addressType,
			[]byte{SOCKS_ADDR_TYPE_IPV4, SOCKS_ADDR_TYPE_FQDN}))
	}
}

func readIpv4AndPort(client *net.TCPConn) (net.IP, int, error) {
	// Check ipv4
	ip := make([]byte, 4)
	_, err := io.ReadFull(client, ip)
	if err != nil {
		return nil, -1, errors.New("no ipv4 address")
	}
	// Check port
	portBytes := make([]byte, 2)
	_, err = io.ReadFull(client, portBytes)
	if err != nil {
		return nil, -1, errors.New("no port")
	}
	port := int(binary.BigEndian.Uint16(portBytes))
	return ip, port, nil
}

func ipv4Connect(ipv4 net.IP, port int) (*net.TCPConn, byte, error) {
	// Connect to peer
	tcpAddr := &net.TCPAddr{
		IP:   ipv4,
		Port: port,
	}
	peer, err := net.DialTCP("tcp4", nil, tcpAddr)
	if err != nil {
		var opErr *net.OpError
		if errors.As(err, &opErr) {
			if opErr.Temporary() {
				return nil, SOCKS_REPLY_TTL_EXPIRED, errors.New(err.Error())
			} else if opErr.Err.Error() == "network is unreachable" {
				return nil, SOCKS_REPLY_NETWORK_UNREACHABLE, errors.New(err.Error())
			} else if opErr.Err.Error() == "no route to host" {
				return nil, SOCKS_REPLY_HOST_UNREACHABLE, errors.New(err.Error())
			} else if opErr.Err.Error() == "connection refused" {
				return nil, SOCKS_REPLY_CONNECTION_REFUSED, errors.New(err.Error())
			}
		}
		return nil, SOCKS_REPLY_GENERAL_SOCKS_SERVER_FAILURE, errors.New(err.Error())
	}
	return peer, SOCKS_REPLY_SUCCEEDED, nil
}

func readDomainNameAndPort(client *net.TCPConn) (string, int, error) {
	// Check domain name size
	domainNameSize := make([]byte, 1)
	_, err := io.ReadFull(client, domainNameSize)
	if err != nil {
		return "", -1, errors.New("no domain name size")
	}
	// Check domain name
	domainName := make([]byte, domainNameSize[0])
	_, err = io.ReadFull(client, domainName)
	if err != nil {
		return "", -1, errors.New("no domain name")
	}
	// Check port
	portBytes := make([]byte, 2)
	_, err = io.ReadFull(client, portBytes)
	if err != nil {
		return "", -1, errors.New("no port")
	}
	port := int(binary.BigEndian.Uint16(portBytes))
	return string(domainName), port, nil
}

func domainNameConnect(domainName string, port int, client *net.TCPConn) (*net.TCPConn, byte, error) {
	// Resolve domain name
	ips, err := net.LookupIP(domainName)
	if err != nil {
		return nil, SOCKS_REPLY_HOST_UNREACHABLE, err
	}
	log.Printf("%v: Domain name %v is resolved to %v\n", client.RemoteAddr(), domainName, ips)
	// Try connecting to each ipv4 address
	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			peer, reply, err := ipv4Connect(ipv4, port)
			if err == nil {
				// Found working ip address
				return peer, reply, nil
			}
		}
	}
	// Not found working ipv4 address
	return nil, SOCKS_REPLY_HOST_UNREACHABLE,
		errors.New("no hosts with IPv6 addresses or working IPv4 addresses")
}

func sendCommandReply(client *net.TCPConn, reply byte) error {
	// Create message
	replyMsg := []byte{
		SOCKS_VERSION, reply, SOCKS_RESERVED_BYTE, SOCKS_ADDR_TYPE_IPV4, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	// Send reply
	_, err := client.Write(replyMsg)
	if err != nil {
		return err
	}
	return nil
}
