package protection_parser

import (
	"strconv"

	excelAdapterModels "github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/models"
	"github.com/rs/zerolog"
	trace "github.com/rs/zerolog/log"
)

func NewProtectionParser(globalInfo excelAdapterModels.GlobalInfo, faultId uint16, warningId uint16, errorId uint16) ProtectionParser {
	parserLogger := trace.With().Str("component", "protection parser").Logger()

	idToBoard := getIdToBoard(globalInfo.BoardToId, parserLogger)

	return ProtectionParser{
		faultId:       faultId,
		errorId:       errorId,
		warningId:     warningId,
		boardIdToName: idToBoard,
		trace:         parserLogger,
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
