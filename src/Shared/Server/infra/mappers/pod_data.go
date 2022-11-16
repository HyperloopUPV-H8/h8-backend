package mappers

import (
	excelAdapter "github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_adapter/domain"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/server/domain"
)

func NewPodData(boards map[string]excelAdapter.BoardDTO) domain.PodData {
	return domain.PodData{
		Boards: newBoards(boards),
	}
}
