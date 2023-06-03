package pipe

import (
	"encoding/binary"
	"errors"
	"net"
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/packet"
	"github.com/rs/zerolog"
)

type Pipe struct {
	conn *net.TCPConn

	laddr *net.TCPAddr
	raddr *net.TCPAddr

	isClosed bool
	mtu      int

	output             chan<- packet.Packet
	onConnectionChange func(bool)

	trace zerolog.Logger
}

func (pipe *Pipe) connect() {
	pipe.trace.Debug().Msg("connecting")
	dialer := net.Dialer{
		LocalAddr: pipe.laddr,
	}
	for pipe.isClosed {
		pipe.trace.Trace().Msg("dial")
		dialer.Deadline = time.Now().Add(time.Millisecond * 2500)
		if conn, err := dialer.Dial("tcp", pipe.raddr.String()); err == nil {
			pipe.open(conn.(*net.TCPConn))
		} else {
			pipe.trace.Trace().Stack().Err(err).Msg("dial failed")
		}
	}
	pipe.trace.Info().Msg("connected")

	go pipe.listen()
}

func (pipe *Pipe) open(conn *net.TCPConn) {
	pipe.trace.Debug().Msg("open")
	pipe.conn = conn
	pipe.isClosed = false
	pipe.onConnectionChange(!pipe.isClosed)
}

func (pipe *Pipe) listen() {
	pipe.trace.Info().Msg("start listening")
	for {
		buffer := make([]byte, pipe.mtu)
		n, err := pipe.conn.Read(buffer)
		if err != nil {
			pipe.trace.Error().Stack().Err(err).Msg("")
			pipe.Close(true)
			return
		}

		if pipe.output == nil {
			pipe.trace.Debug().Msg("no output configured")
			continue
		}

		pipe.trace.Trace().Msg("new message")

		raw := pipe.getRaw(buffer[:n])

		pipe.output <- raw
	}
}

var syntheticSeqNum uint32 = 0

func (pipe *Pipe) getRaw(payload []byte) packet.Packet {
	syntheticSeqNum++
	return packet.Packet{
		Metadata: packet.NewMetaData(pipe.raddr.String(), pipe.laddr.String(), binary.LittleEndian.Uint16(payload[0:2]), syntheticSeqNum, time.Now()),
		Payload:  payload[2:],
	}
}

func (pipe *Pipe) Write(data []byte) (int, error) {
	if pipe == nil || pipe.conn == nil {
		err := errors.New("pipe is nil")
		pipe.trace.Error().Stack().Err(err).Msg("")
		return 0, err
	}

	pipe.trace.Trace().Msg("write")
	pipe.conn.SetWriteDeadline(time.Now().Add(time.Millisecond * 2500))
	return pipe.conn.Write(data)
}

func (pipe *Pipe) Close(reconnect bool) error {
	pipe.trace.Warn().Bool("reconnect", reconnect).Msg("close")

	err := pipe.conn.Close()
	pipe.isClosed = err == nil
	pipe.onConnectionChange(!pipe.isClosed)

	if reconnect {
		go pipe.connect()
	}
	return err
}

func (pipe *Pipe) Laddr() string {
	return pipe.laddr.String()
}

func (pipe *Pipe) Raddr() string {
	return pipe.raddr.String()
}
