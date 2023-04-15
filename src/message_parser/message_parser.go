package message_parser

import (
	"bytes"
	"strconv"
	"strings"

	excel_adapter_models "github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/models"
	"github.com/HyperloopUPV-H8/Backend-H8/message_parser/models"
)

type MessageParser struct {
	faultId   uint16
	warningId uint16
	blcuAckId uint16
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
		//TODO: trace
	}

	warningId, warningErr := strconv.Atoi(globalInfo.MessageToId[config.WarningIdKey])

	if warningErr != nil {
		//TODO: trace
	}

	blcuAckId, blcuErr := strconv.Atoi(globalInfo.MessageToId[config.BlcuAckIdKey])

	if blcuErr != nil {
		//TODO: trace
	}

	return MessageParser{
		faultId:   uint16(faultId),
		warningId: uint16(warningId),
		blcuAckId: uint16(blcuAckId),
	}
}

func (parser MessageParser) Parse(raw []byte) interface{} {
	rawStr := bytes.NewBuffer(raw).String()
	messageParts := strings.Split(rawStr, "\n")
	id, err := strconv.Atoi(messageParts[0])

	if err != nil {
		//TODO: error
	}

	if uint16(id) == parser.faultId {
		//TODO: pasar la messageParts sin el primer elemento es confuso,
		// sobre todo dentro de ParseProtectionMessage
		return models.ParseProtectionMessage("fault", messageParts[1:])
	} else if uint16(id) == parser.warningId {
		return models.ParseProtectionMessage("warning", messageParts[1:])
	} else if uint16(id) == parser.blcuAckId {
		return models.BLCU_ACK{}
	} else {
		panic("to remove return error")
	}
}
