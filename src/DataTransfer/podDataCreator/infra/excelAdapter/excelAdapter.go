package excelAdapter

import (
	"github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/excelRetreiver"
	podDataCreator "github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/podDataCreator/domain"
	"github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/podDataCreator/infra/excelAdapter/dto"
)

func GetPackets(tables map[string]excelRetreiver.Table) map[dto.Id]*podDataCreator.Packet {
	boardDTO := dto.NewBoardDTO(tables)
	packets := boardDTO.GetPackets()

	return packets
}
