package message_parser

import (
	"bytes"
	"errors"
	"strconv"
	"strings"

	excel_adapter_models "github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/models"
	"github.com/HyperloopUPV-H8/Backend-H8/message_parser/models"
	"github.com/rs/zerolog"
	trace "github.com/rs/zerolog/log"
)

const COMPONENT_NAME = "messageParser"

type MessageParser struct {
	faultId   uint16
	warningId uint16
	blcuAckId uint16
	trace     zerolog.Logger
}

type MessageParserConfig struct {
	FaultIdKey   string `toml:"fault_id_key"`
	WarningIdKey string `toml:"warning_id_key"`
	BlcuAckIdKey string `toml:"blcu_ack_id"`
}

// 12 CURRENT_1 OUTOFBOUNDS 3 12 34? timestamp
// Primero got, luego want

func New(globalInfo excel_adapter_models.GlobalInfo, config MessageParserConfig) MessageParser {
	faultId, faultErr := strconv.Atoi(globalInfo.MessageToId[config.FaultIdKey])

	if faultErr != nil {
		trace.Fatal().Stack().Err(faultErr).Msg("")
	}

	warningId, warningErr := strconv.Atoi(globalInfo.MessageToId[config.WarningIdKey])

	if warningErr != nil {
		trace.Fatal().Stack().Err(warningErr).Msg("")
	}

	blcuAckId, blcuErr := strconv.Atoi(globalInfo.MessageToId[config.BlcuAckIdKey])

	if blcuErr != nil {
		trace.Fatal().Stack().Err(blcuErr).Msg("")
	}

	return MessageParser{
		faultId:   uint16(faultId),
		warningId: uint16(warningId),
		blcuAckId: uint16(blcuAckId),
		trace:     trace.With().Str("component", COMPONENT_NAME).Logger(),
	}
}

func (parser MessageParser) Parse(raw []byte) (interface{}, error) {
	rawStr := bytes.NewBuffer(raw).String()
	messageParts := strings.Split(rawStr, "\n")
	id, err := strconv.Atoi(messageParts[0])

	if err != nil {
		parser.trace.Error().Err(err).Stack().Msg("")
		return nil, err
	}

	if uint16(id) == parser.faultId {
		//TODO: pasar la messageParts sin el primer elemento es confuso,
		// sobre todo dentro de ParseProtectionMessage
		return models.ParseProtectionMessage("fault", messageParts[1:])
	} else if uint16(id) == parser.warningId {
		return models.ParseProtectionMessage("warning", messageParts[1:])
	} else if uint16(id) == parser.blcuAckId {
		return models.BLCU_ACK{}, nil
	} else {
		err = errors.New("unidentified message")
		trace.Fatal().Err(err).Stack().Msg("")
		return nil, err
	}
}
