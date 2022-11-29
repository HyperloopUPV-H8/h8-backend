package tcp

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
)

type Pipe struct {
	laddr *net.TCPAddr
	raddr *net.TCPAddr

	conn         *net.TCPConn
	connRMx      sync.Mutex
	connWMx      sync.Mutex
	onConnUpdate func(*net.TCPAddr, bool)

	onRead  func([]byte)
	snaplen int32
}

func getPipes(laddr string, rips []string, rports []uint16, snaplen int32) map[string]*Pipe {
	tcpLAddr, err := net.ResolveTCPAddr("tcp", laddr)
	if err != nil {
		log.Fatalf("tcp: get pipes: %s\n", err)
	}

	pipes := make(map[string]*Pipe, len(rips))
	for i, ip := range rips {
		raddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", ip, rports[i]))
		if err != nil {
			log.Fatalf("tcp: get pipes: %s\n", err)
		}
		pipes[ip] = newPipe(tcpLAddr, raddr, snaplen)
	}
	return pipes
}

func newPipe(laddr *net.TCPAddr, raddr *net.TCPAddr, snaplen int32) *Pipe {
	pipe := &Pipe{
		laddr:        laddr,
		raddr:        raddr,
		conn:         nil,
		connRMx:      sync.Mutex{},
		connWMx:      sync.Mutex{},
		onRead:       func(b []byte) {},
		onConnUpdate: func(t *net.TCPAddr, b bool) {},
		snaplen:      snaplen,
	}

	go pipe.tryConnect()

	return pipe
}

func (pipe *Pipe) setOnRead(action func([]byte)) {
	pipe.onRead = action
}

func (pipe *Pipe) setOnUpdate(action func(*net.TCPAddr, bool)) {
	pipe.onConnUpdate = action
}

func (pipe *Pipe) tryConnect() {
	for pipe.conn == nil {
		if conn, _ := net.DialTCP("tcp", pipe.laddr, pipe.raddr); conn != nil {
			pipe.connect(conn)
		}
	}
}

func (pipe *Pipe) connect(conn *net.TCPConn) {
	pipe.connRMx.Lock()
	pipe.connWMx.Lock()
	defer pipe.connRMx.Unlock()
	defer pipe.connWMx.Unlock()
	pipe.conn = conn
	go pipe.read()
	pipe.onConnUpdate(pipe.raddr, true)
}

func (pipe *Pipe) read() {
	for pipe.conn != nil {
		buf := make([]byte, pipe.snaplen)

		pipe.connRMx.Lock()
		n, err := pipe.conn.Read(buf)
		pipe.connRMx.Unlock()

		if err != nil {
			pipe.close()
		} else {
			pipe.onRead(buf[:n])
		}
	}
}

func (pipe *Pipe) write(payload []byte) error {
	pipe.connWMx.Lock()
	defer pipe.connWMx.Unlock()
	if pipe.conn == nil {
		return errors.New("tcp: write: board is disconnected")
	}

	n, err := pipe.conn.Write(payload)
	if err != nil {
		pipe.close()
		return err
	} else if n != len(payload) {
		return errors.New("tcp: write: couldn't send all data in same packet")
	}

	return nil
}

func (pipe *Pipe) close() {
	pipe.connRMx.Lock()
	pipe.connWMx.Lock()
	defer pipe.connRMx.Unlock()
	defer pipe.connWMx.Unlock()
	if pipe.conn == nil {
		return
	}
	pipe.conn.Close()
	pipe.conn = nil
	go pipe.tryConnect()
	pipe.onConnUpdate(pipe.raddr, false)
}
