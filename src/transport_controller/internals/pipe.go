package internals

import (
	"net"

	"github.com/HyperloopUPV-H8/Backend-H8/transport_controller/models"
)

type Pipe struct {
	conn   *net.TCPConn
	laddr  *net.TCPAddr
	raddr  *net.TCPAddr
	config models.Config
}

func (pipe *Pipe) open(conn *net.TCPConn) {
	if conn == nil {
		return
	}

	pipe.conn = conn
	pipe.config.OnConnUpdate(pipe.laddr, true)
	go pipe.read()
}

func (pipe *Pipe) close() error {
	if pipe.conn == nil {
		return nil
	}

	if err := pipe.conn.Close(); err != nil {
		return err
	}

	pipe.conn = nil
	pipe.config.OnConnUpdate(pipe.laddr, false)
	go pipe.connect()

	return nil
}

func (pipe *Pipe) connect() {
	for pipe.conn == nil {
		if conn, err := net.DialTCP("tcp", pipe.laddr, pipe.raddr); err == nil {
			pipe.open(conn)
		}
	}
}

func (pipe *Pipe) read() {
	for pipe.conn != nil {
		buf := make([]byte, pipe.config.Snaplen)
		if n, err := pipe.conn.Read(buf); err == nil {
			pipe.config.Dump <- buf[:n]
		} else {
			pipe.close()
		}
	}
}

func (pipe *Pipe) write(payload []byte) (success bool) {
	if _, err := pipe.conn.Write(payload); err != nil {
		pipe.close()
		return false
	}
	return true
}

func createPipe(laddr *net.TCPAddr, raddr *net.TCPAddr, config models.Config) *Pipe {
	pipe := &Pipe{
		conn:   nil,
		laddr:  laddr,
		raddr:  raddr,
		config: config,
	}

	go pipe.connect()

	return pipe
}
