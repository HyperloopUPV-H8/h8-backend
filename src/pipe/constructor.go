package pipe

import (
	"net"

	"github.com/HyperloopUPV-H8/Backend-H8/common"
	"github.com/HyperloopUPV-H8/Backend-H8/info"
	"github.com/HyperloopUPV-H8/Backend-H8/packet"
	"github.com/rs/zerolog"
	trace "github.com/rs/zerolog/log"
)

func contains(boards []string, board string) bool {
	for _, b := range boards {
		if b == board {
			return true
		}
	}
	return false
}

func CreatePipes(info info.Info, boards []string, dataChan chan<- packet.Packet, onConnectionChange func(string, bool), config Config, readers map[uint16]common.ReaderFrom, trace zerolog.Logger) map[string]*Pipe {
	laddr := net.TCPAddr{
		IP:   info.Addresses.Backend,
		Port: int(info.Ports.TcpClient),
	}

	pipes := make(map[string]*Pipe)

	for board, ip := range info.Addresses.Boards {
		if boards != nil && !contains(boards, board) {
			continue
		}

		raddr := net.TCPAddr{
			IP:   ip,
			Port: int(info.Ports.TcpServer),
		}
		pipe, err := newPipe(laddr, raddr, config.Mtu, dataChan, readers, getOnConnectionChange(board, onConnectionChange))
		if err != nil {
			//TODO: how to handle this error
			trace.Fatal().Stack().Err(err).Msg("error creating pipe")
		}

		pipes[board] = pipe

	}

	return pipes
}

func newPipe(laddr net.TCPAddr, raddr net.TCPAddr, mtu uint, outputChan chan<- packet.Packet, readers map[uint16]common.ReaderFrom, onConnectionChange func(bool)) (*Pipe, error) {
	trace.Info().Str("laddr", laddr.String()).Str("raddr", raddr.String()).Msg("new pipe")

	pipe := &Pipe{
		laddr:  &laddr,
		raddr:  &raddr,
		output: outputChan,

		readers: readers,

		isClosed: true,
		mtu:      int(mtu),

		onConnectionChange: onConnectionChange,

		trace: trace.With().Str("component", "pipe").IPAddr("addr", raddr.IP).Logger(),
	}

	go pipe.connect()

	return pipe, nil
}

func getOnConnectionChange(board string, onConnectionChange func(string, bool)) func(bool) {
	return func(state bool) {
		onConnectionChange(board, state)
	}
}
