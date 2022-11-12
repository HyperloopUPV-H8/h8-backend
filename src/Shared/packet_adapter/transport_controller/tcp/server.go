package tcp

import (
	"fmt"
	"net"
)

type IP = string
type Port = uint16
type Payload = []byte

type Server struct {
	listener   *net.TCPListener
	pipes      Pipes
	validAddrs []IP
}

func Open(localPort Port, remoteAddrs []IP) Server {
	server := Server{
		listener:   bindListener(resolvePortAddr(localPort)),
		pipes:      NewPipes(len(remoteAddrs)),
		validAddrs: remoteAddrs,
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
		if err == nil && server.isValidAddr(getTCPConnIP(conn)) {
			conn.SetLinger(0)
			server.pipes.AddConnection(getTCPConnIP(conn), conn)
		}
	}
}

func getTCPConnIP(conn *net.TCPConn) IP {
	connTCPAddr, err := net.ResolveTCPAddr("tcp", conn.RemoteAddr().String())
	if err != nil {
		panic(err)
	}
	return IP(connTCPAddr.IP.String())
}

func (server *Server) Send(ip IP, payload Payload) {
	server.pipes.Send(ip, payload)
}

func (server *Server) ReceiveNext() []Payload {
	for {
		payloads := server.pipes.Receive()
		if len(payloads) != 0 {
			return payloads
		}
	}
}

func (server *Server) ConnectedAddresses() []IP {
	return server.pipes.ConnectedIPs()
}

func (server *Server) isValidAddr(addr IP) bool {
	for _, ip := range server.validAddrs {
		if ip == addr {
			return true
		}
	}

	return false
}

func (server *Server) Close() {
	server.listener.Close()
	server.pipes.Close()
}
