package infra

import (
	"fmt"
	"net"
	"sync"
	"time"
)

const (
	packetMaxLength int = 1024
)

type Server struct {
	listener    *net.TCPListener
	connections map[IP]*net.TCPConn
	connMutex   *sync.Mutex
	validAddrs  []IP
}

func OpenServer(localPort Port, remoteAddrs []IP) Server {
	server := Server{
		listener:    bindListener(resolvePortAddr(localPort)),
		connections: make(map[IP]*net.TCPConn, len(remoteAddrs)),
		validAddrs:  remoteAddrs,
	}

	go server.listenConnections()

	return server
}

// Only specifying the port makes go to listen for traffic on all ips
func resolvePortAddr(port Port) *net.TCPAddr {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	return addr
}

func bindListener(addr *net.TCPAddr) *net.TCPListener {
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}

	return listener
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

		if !server.filterValidAddr(connIP) {
			continue
		}

		server.connMutex.Lock()
		defer server.connMutex.Unlock()

		server.connections[IP(connIP)] = conn
		defer conn.Close()
	}
}

func (server *Server) Send(ip IP, payload []byte) {
	server.connMutex.Lock()
	defer server.connMutex.Unlock()

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
	server.connMutex.Lock()
	defer server.connMutex.Unlock()

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

func (server *Server) ConnectedAddresses() []IP {
	addresses := make([]IP, 0)

	server.connMutex.Lock()
	defer server.connMutex.Unlock()

	for k := range server.connections {
		addresses = append(addresses, k)
	}
	return addresses
}

func (server *Server) disconnect(ip IP) {
	// No need to lock here because it's already locked everywhere it's called
	conn := server.connections[ip]
	conn.Close()
	delete(server.connections, ip)
}

func (server *Server) filterValidAddr(addr IP) bool {
	for _, ip := range server.validAddrs {
		if ip == addr {
			return true
		}
	}

	return false
}
