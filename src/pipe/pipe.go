package pipe

import (
	"errors"
	"log"
	"net"
)

const READ_BUFFER_SIZE = 1500

type Pipe struct {
	conn *net.TCPConn

	laddr *net.TCPAddr
	raddr *net.TCPAddr

	isClosed bool

	output             chan<- []byte
	onConnectionChange func(bool)
}

func New(laddr string, raddr string) (*Pipe, error) {
	localAddr, err := net.ResolveTCPAddr("tcp", laddr)
	if err != nil {
		return nil, err
	}

	remoteAddr, err := net.ResolveTCPAddr("tcp", raddr)
	if err != nil {
		return nil, err
	}

	pipe := &Pipe{
		laddr:    localAddr,
		raddr:    remoteAddr,
		isClosed: true,
	}

	go pipe.connect()

	return pipe, nil
}

func (pipe *Pipe) connect() {
	for pipe.isClosed {
		if conn, err := net.DialTCP("tcp", pipe.laddr, pipe.raddr); err == nil {
			pipe.open(conn)
		}
	}

	go pipe.listen()
}

func (pipe *Pipe) open(conn *net.TCPConn) {
	log.Println("pipe open")
	pipe.conn = conn
	pipe.isClosed = false
	pipe.onConnectionChange(!pipe.isClosed)
}

func (pipe *Pipe) SetOutput(output chan<- []byte) {
	pipe.output = output
}

func (pipe *Pipe) listen() {
	for {
		buffer := make([]byte, READ_BUFFER_SIZE)
		n, err := pipe.conn.Read(buffer)
		if err != nil {
			pipe.Close(true)
			return
		}

		if pipe.output == nil {
			continue
		}
		pipe.output <- buffer[:n]
	}
}

func (pipe *Pipe) Write(data []byte) (int, error) {
	if pipe == nil || pipe.conn == nil {
		return 0, errors.New("pipe is nil")
	}
	return pipe.conn.Write(data)
}

func (pipe *Pipe) Close(reconnect bool) error {
	log.Println("pipe close")
	err := pipe.conn.Close()
	pipe.isClosed = err == nil
	pipe.onConnectionChange(!pipe.isClosed)

	if reconnect {
		go pipe.connect()
	}
	return err
}

func (pipe *Pipe) OnConnectionChange(callback func(bool)) {
	pipe.onConnectionChange = callback
}
