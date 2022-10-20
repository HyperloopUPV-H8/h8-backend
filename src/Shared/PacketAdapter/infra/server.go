package infra

import (
	"fmt"
	"net"
	"time"
)

const (
	packetMaxLength int = 1024
)

type Server struct {
	listener    net.TCPListener
	connections map[IP]net.TCPConn
	validAddrs  []IP
}

func OpenServer(localPort Port, remoteAddrs []IP) Server {
	laddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", localPort))
	if err != nil {
		panic(err)
	}

	server := Server{
		listener:    bindListener(laddr),
		connections: make(map[IP]net.TCPConn, len(remoteAddrs)),
		validAddrs:  remoteAddrs,
	}

	go server.listenConnections()

	return server
}

func bindListener(addr *net.TCPAddr) net.TCPListener {
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}

	return *listener
}

func (server *Server) listenConnections() {
	for {
		conn, err := server.listener.AcceptTCP()
		if err != nil {
			continue
		}

		connTCPAddr, err := net.ResolveTCPAddr("tcp", conn.LocalAddr().String())
		if err != nil {
			panic(err)
		}
		connIP := IP(connTCPAddr.IP.String())

		if !server.filterAddr(connIP) {
			continue
		}

		server.connections[IP(connIP)] = *conn
	}
}

func (server *Server) Send(ip IP, payload []byte) {
	conn, exists := server.connections[ip]

	if !exists {
		return
	}

	_, err := conn.Write(payload)

	if err != nil {
		server.disconnect(ip)
	}
}

func (server *Server) Receive() []byte {
	for _, conn := range server.connections {
		conn.SetDeadline(time.Now())
		buf := make([]byte, packetMaxLength)
		_, err := conn.Read(buf)
		if err != nil {
			continue
		}
		return buf
	}
	return nil
}

func (server *Server) ConnectedAddresses() []string {
	addresses := make([]string, 0)
	for k := range server.connections {
		addresses = append(addresses, string(k))
	}
	return addresses
}

func (server *Server) disconnect(ip IP) {
	conn := server.connections[ip]
	conn.Close()
	delete(server.connections, ip)
}

func (server *Server) filterAddr(addr IP) bool {
	for _, ip := range server.validAddrs {
		if ip == addr {
			return true
		}
	}

	return false
}
