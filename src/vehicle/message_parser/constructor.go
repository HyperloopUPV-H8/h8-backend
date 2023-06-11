package message_parser

import (
	"github.com/HyperloopUPV-H8/Backend-H8/info"
	"github.com/rs/zerolog"
	trace "github.com/rs/zerolog/log"
)

func NewMessageParser(info info.Info) MessageParser {
	parserLogger := trace.With().Str("component", "protection parser").Logger()

	idToBoard := getIdToBoard(info.BoardIds, parserLogger)

	return MessageParser{
		infoId:             info.MessageIds.Info,
		warningId:          info.MessageIds.Warning,
		faultId:            info.MessageIds.Fault,
		boardIdToName:      idToBoard,
		trace:              parserLogger,
		addStateOrderId:    info.MessageIds.AddStateOrder,
		removeStateOrderId: info.MessageIds.RemoveStateOrder,
	}
}

func getIdToBoard(boardToId map[string]uint16, trace zerolog.Logger) map[uint16]string {
	idToBoard := make(map[uint16]string)

	for board, id := range boardToId {

		idToBoard[id] = board
	}

	return idToBoard
}
