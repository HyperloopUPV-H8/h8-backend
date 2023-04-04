package pipe

import (
	"errors"
	"net"
	"os"
	"strconv"

	"github.com/rs/zerolog"
	trace "github.com/rs/zerolog/log"
)

type Pipe struct {
	conn *net.TCPConn

	laddr *net.TCPAddr
	raddr *net.TCPAddr

	isClosed bool
	mtu      int

	output             chan<- []byte
	onConnectionChange func(bool)

	trace zerolog.Logger
}

func New(laddr string, raddr string) (*Pipe, error) {
	trace.Info().Str("laddr", laddr).Str("raddr", raddr).Msg("new pipe")
	localAddr, err := net.ResolveTCPAddr("tcp", laddr)
	if err != nil {
		trace.Error().Str("laddr", laddr).Stack().Err(err).Msg("")
		return nil, err
	}

	remoteAddr, err := net.ResolveTCPAddr("tcp", raddr)
	if err != nil {
		trace.Error().Str("raddr", raddr).Stack().Err(err).Msg("")
		return nil, err
	}

	mtu, err := strconv.ParseInt(os.Getenv("INTERFACE_MTU"), 10, 32)
	if err != nil {
		trace.Fatal().Stack().Err(err).Str("INTERFACE_MTU", os.Getenv("INTERFACE_MTU")).Msg("")
		return nil, err
	}

	pipe := &Pipe{
		laddr: localAddr,
		raddr: remoteAddr,

		isClosed: true,
		mtu:      int(mtu),

		trace: trace.With().Str("component", "pipe").IPAddr("addr", remoteAddr.IP).Logger(),
	}

	go pipe.connect()

	return pipe, nil
}

func (pipe *Pipe) connect() {
	pipe.trace.Debug().Msg("connecting")
	for pipe.isClosed {
		pipe.trace.Trace().Msg("dial")
		if conn, err := net.DialTCP("tcp", pipe.laddr, pipe.raddr); err == nil {
			pipe.open(conn)
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

func (pipe *Pipe) SetOutput(output chan<- []byte) {
	pipe.trace.Debug().Msg("set output")
	pipe.output = output
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
		pipe.output <- buffer[:n]
	}
}

func (pipe *Pipe) Write(data []byte) (int, error) {
	if pipe == nil || pipe.conn == nil {
		err := errors.New("pipe is nil")
		pipe.trace.Error().Stack().Err(err).Msg("")
		return 0, err
	}

	pipe.trace.Trace().Msg("write")
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

func (pipe *Pipe) OnConnectionChange(callback func(bool)) {
	pipe.trace.Debug().Msg("set on connection change")
	pipe.onConnectionChange = callback
}
