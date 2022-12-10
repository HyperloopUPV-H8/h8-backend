package internals

import (
	"log"
	"net"
	"strings"

	"github.com/HyperloopUPV-H8/Backend-H8/transport_controller/models"
)

type PipeHandle struct {
	listener *net.TCPListener
	pipes    map[string]*Pipe
}

func OpenPipes(laddr *net.TCPAddr, raddrs []*net.TCPAddr, config models.Config) *PipeHandle {
	handle := &PipeHandle{
		listener: bindListener(laddr),
		pipes:    getPipes(laddr, raddrs, config),
	}

	go handle.listen()

	return handle
}

func bindListener(addr *net.TCPAddr) *net.TCPListener {
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatalf("PipeHandle: bindListener: %s\n", err)
	}
	return listener
}

func (handle *PipeHandle) listen() {
	for handle.listener != nil {
		conn, err := handle.listener.AcceptTCP()
		if err != nil {
			continue
		}
		pipe, ok := handle.pipes[withoutPort(conn.RemoteAddr())]
		if !ok {
			continue
		}

		pipe.open(conn)
	}
}

func (handle *PipeHandle) Write(addr string, payload []byte) bool {
	pipe, ok := handle.pipes[addr]
	if !ok {
		return false
	}
	return pipe.write(payload)
}

func getPipes(laddr *net.TCPAddr, raddrs []*net.TCPAddr, config models.Config) map[string]*Pipe {
	pipes := make(map[string]*Pipe, len(raddrs))
	for _, raddr := range raddrs {
		pipes[raddr.IP.String()] = createPipe(laddr, raddr, config)
	}
	return pipes
}

func withoutPort(addr net.Addr) string {
	return strings.Split(addr.String(), ":")[0]
}

func (handle *PipeHandle) Close() {
	handle.listener.Close()
	for _, pipe := range handle.pipes {
		pipe.close()
	}
}
