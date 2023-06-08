package message_parser

import (
	"strconv"

	excelAdapterModels "github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/models"
	"github.com/rs/zerolog"
	trace "github.com/rs/zerolog/log"
)

func NewMessageParser(globalInfo excelAdapterModels.GlobalInfo, infoId uint16, faultId uint16, warningId uint16, addStateOrdersId uint16, removeStateOrdersId uint16) MessageParser {
	parserLogger := trace.With().Str("component", "protection parser").Logger()

	idToBoard := getIdToBoard(globalInfo.BoardToId, parserLogger)

	return MessageParser{
		infoId:             infoId,
		warningId:          warningId,
		faultId:            faultId,
		boardIdToName:      idToBoard,
		trace:              parserLogger,
		addStateOrderId:    addStateOrdersId,
		removeStateOrderId: removeStateOrdersId,
	}
}

func getIdToBoard(boardToId map[string]string, trace zerolog.Logger) map[uint]string {
	idToBoard := make(map[uint]string)

	for board, idStr := range boardToId {
		id, err := strconv.Atoi(idStr)

		if err != nil {
			trace.Error().Err(err).Str("id", idStr).Msg("error parsing board id")
		}

		idToBoard[uint(id)] = board
	}

	return idToBoard
}
