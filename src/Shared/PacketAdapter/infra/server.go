package PacketAdapter

import (
	"fmt"
	"net"
)

type Server struct {
	listener   net.TCPListener
	pipes      SyncPipes
	validAddrs []IP
}

func OpenServer(localPort Port, remoteAddrs []IP) Server {
	server := Server{
		listener:   bindListener(resolvePortAddr(localPort)),
		pipes:      NewPipes(len(remoteAddrs)),
		validAddrs: remoteAddrs,
	}

	go server.listenConnections()

	return server
}

// Only specifying the port makes go to listen for traffic on all ips
func resolvePortAddr(port Port) net.TCPAddr {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	return *addr
}

func bindListener(addr net.TCPAddr) net.TCPListener {
	listener, err := net.ListenTCP("tcp", &addr)
	if err != nil {
		panic(err)
	}

	return *listener
}

func (server Server) listenConnections() {
	for {
		if conn, err := server.accept(); err == nil && server.isValidAddr(getTCPConnIP(conn)) {
			server.pipes.AddConnection(getTCPConnIP(conn), conn)
		}
	}
}

func (server Server) accept() (net.TCPConn, error) {
	conn, err := server.listener.AcceptTCP()
	return *conn, err
}

func getTCPConnIP(conn net.TCPConn) IP {
	connTCPAddr, err := net.ResolveTCPAddr("tcp", conn.LocalAddr().String())
	if err != nil {
		panic(err)
	}
	return IP(connTCPAddr.IP.String())
}

func (server Server) Send(ip IP, payload Payload) {
	err := server.pipes.Send(ip, payload)
	if err != nil {
		server.pipes.RemoveConnection(ip)
	}
}

func (server Server) Receive() []Payload {
	for {
		payloads := server.pipes.Receive()
		if len(payloads) != 0 {
			return payloads
		}
	}
}

func (server Server) ConnectedAddresses() []IP {
	return server.pipes.ConnectedIPs()
}

func (server Server) isValidAddr(addr IP) bool {
	for _, ip := range server.validAddrs {
		if ip == addr {
			return true
		}
	}

	return false
}
