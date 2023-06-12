package pipe

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/common"
	excel_models "github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/models"
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

func CreatePipes(global excel_models.GlobalInfo, keepaliveInterval, writeTimeout *time.Duration, boards []string, dataChan chan<- packet.Packet, onConnectionChange func(string, bool), config Config, readers map[uint16]common.ReaderFrom, trace zerolog.Logger) map[string]*Pipe {
	pipes := make(map[string]*Pipe)
	i := uint64(0)
	for board, ip := range global.BoardToIP {
		ip := ip
		if boards != nil && !contains(boards, board) {
			continue
		}
		port, err := strconv.ParseUint(global.ProtocolToPort[config.TcpClientTag], 10, 16)
		if err != nil {
			panic(err)
		}
		laddr := common.AddrWithPort(global.BackendIP, fmt.Sprintf("%d", port+i))
		raddr := common.AddrWithPort(ip, global.ProtocolToPort[config.TcpServerTag])
		pipe, err := newPipe(laddr, raddr, keepaliveInterval, writeTimeout, config.Mtu, dataChan, readers, getOnConnectionChange(board, onConnectionChange))
		if err != nil {
			trace.Fatal().Stack().Err(err).Msg("error creating pipe")
		}

		pipes[board] = pipe
		i++
	}
	return pipes
}

func newPipe(laddr string, raddr string, keepaliveInterval, writeTimeout *time.Duration, mtu uint, outputChan chan<- packet.Packet, readers map[uint16]common.ReaderFrom, onConnectionChange func(bool)) (*Pipe, error) {
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

	pipe := &Pipe{
		laddr:  localAddr,
		raddr:  remoteAddr,
		output: outputChan,

		readers: readers,

		isClosed: true,
		mtu:      int(mtu),

		onConnectionChange: onConnectionChange,

		keepaliveInterval: keepaliveInterval,
		writeTiemout:      writeTimeout,

		trace: trace.With().Str("component", "pipe").IPAddr("addr", remoteAddr.IP).Logger(),
	}

	go pipe.connect()

	return pipe, nil
}

func getOnConnectionChange(board string, onConnectionChange func(string, bool)) func(bool) {
	return func(state bool) {
		onConnectionChange(board, state)
	}
}
