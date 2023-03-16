package internals

import (
	"net"

	"github.com/HyperloopUPV-H8/Backend-H8/transport_controller/models"
)

type PipeHandle struct {
	pipes map[string]*Pipe
}

func OpenPipes(laddr *net.TCPAddr, raddrs []*net.TCPAddr, config models.Config) *PipeHandle {
	handle := &PipeHandle{
		pipes: getPipes(laddr, raddrs, config),
	}

	return handle
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

func (handle *PipeHandle) Close() {
	for _, pipe := range handle.pipes {
		pipe.close()
	}
}
