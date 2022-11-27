package tcp

import (
	"fmt"
	"net"
	"strings"
)

type Server struct {
	listener *net.TCPListener
	pipes    map[string]*Pipe
	onRead   func([]byte)
}

func Open(config *Config) Server {
	server := Server{
		listener: bindListener(resolvePortAddr(config.LocalPort)),
		pipes:    getPipes(fmt.Sprintf(":%d", config.LocalPort), config.RemoteIPs, config.RemotePorts, config.Snaplen),
		onRead:   func([]byte) {},
	}

	go server.listenConnections()

	return server
}

func (server *Server) SetOnRead(action func([]byte)) {
	server.onRead = action
	for _, pipe := range server.pipes {
		pipe.setOnRead(action)
	}
}

func resolvePortAddr(port uint16) *net.TCPAddr {
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
		pipe, exists := server.pipes[getTCPConnIP(conn)]
		if err != nil && exists {
			pipe.connect(conn)
		} else {
			conn.Close()
		}
	}
}

func getTCPConnIP(conn *net.TCPConn) string {
	return strings.Split(conn.RemoteAddr().String(), ":")[0]
}

func (server *Server) Send(ip string, payload []byte) error {
	return server.pipes[ip].write(payload)
}

func (server *Server) Close() {
	server.listener.Close()
	for _, pipe := range server.pipes {
		pipe.close()
	}
}
