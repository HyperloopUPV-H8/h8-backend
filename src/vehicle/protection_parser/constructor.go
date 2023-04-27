package protection_parser

import (
	"strconv"

	"github.com/HyperloopUPV-H8/Backend-H8/common"
	excelAdapterModels "github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/models"
	"github.com/rs/zerolog"
	trace "github.com/rs/zerolog/log"
)

func NewProtectionParser(globalInfo excelAdapterModels.GlobalInfo, config Config) ProtectionParser {
	parserLogger := trace.With().Str("component", "protection parser").Logger()

	faultId := mustGetId(globalInfo.MessageToId, config.FaultIdKey, parserLogger)
	warningId := mustGetId(globalInfo.MessageToId, config.WarningIdKey, parserLogger)
	errorId := mustGetId(globalInfo.MessageToId, config.ErrorIdKey, parserLogger)
	ids := common.NewSet[uint16]()

	ids.Add(faultId)
	ids.Add(warningId)
	ids.Add(errorId)

	idToBoard := getIdToBoard(globalInfo.BoardToId, parserLogger)

	return ProtectionParser{
		Ids:           ids,
		faultId:       faultId,
		errorId:       errorId,
		warningId:     warningId,
		boardIdToName: idToBoard,
		trace:         parserLogger,
	}
}

func mustGetId(kindToId map[string]string, key string, trace zerolog.Logger) uint16 {
	idStr, ok := kindToId[key]

	if !ok {
		trace.Fatal().Str("key", key).Msg("key not found")
	}

	id, err := strconv.ParseUint(idStr, 10, 16)

	if err != nil {
		trace.Fatal().Str("id", idStr).Msg("error parsing id")
	}

	return uint16(id)
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
