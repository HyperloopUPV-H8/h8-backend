package PacketAdapter

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
		connMutex:   &sync.Mutex{},
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
		conn, err := server.accept()
		if err != nil {
			continue
		}

		connIP := getTCPConnIP(conn)
		if !server.filterValidAddr(connIP) {
			continue
		}

		server.addConnection(connIP, conn)
	}
}

func (server *Server) accept() (*net.TCPConn, error) {
	return server.listener.AcceptTCP()
}

func getTCPConnIP(conn *net.TCPConn) IP {
	connTCPAddr, err := net.ResolveTCPAddr("tcp", conn.LocalAddr().String())
	if err != nil {
		panic(err)
	}
	return IP(connTCPAddr.IP.String())
}

func (server *Server) Send(ip IP, payload []byte) {
	conn, exists := server.connections[ip]

	if !exists {
		return
	}

	_, err := conn.Write(payload)

	if err != nil {
		server.removeConnection(ip)
	}
}

func (server *Server) Receive() []byte {
	for {
		for _, conn := range server.connections {
			conn.SetDeadline(time.Now().Add(time.Microsecond))
			buf := make([]byte, packetMaxLength)
			_, err := conn.Read(buf)
			if err != nil {
				continue
			}
			return buf
		}
	}
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

func (server *Server) filterValidAddr(addr IP) bool {
	for _, ip := range server.validAddrs {
		if ip == addr {
			return true
		}
	}

	return false
}

func (server *Server) addConnection(ip IP, conn *net.TCPConn) {
	server.connMutex.Lock()
	defer server.connMutex.Unlock()
	server.connections[ip] = conn
}

func (server *Server) removeConnection(ip IP) {
	server.connMutex.Lock()
	defer server.connMutex.Unlock()
	conn := server.connections[ip]
	conn.Close()
	delete(server.connections, ip)
}
