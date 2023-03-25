package excel_adapter

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/internals"
	internalModels "github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/internals/models"
	"github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/models"
	"github.com/xuri/excelize/v2"
)

const GLOBAL_SHEET_NAME = "INFO"
const ADDRESSES_TABLE_NAME = "addresses"
const UNITS_TABLE_NAME = "units"
const PORTS_TABLE_NAME = "ports"
const IDS_TABLE_NAME = "ids"

func FetchDocument(id string, path string, name string) internalModels.Document {
	errDownloading := internals.DownloadFile(id, path, name)
	if errDownloading != nil {
		log.Println("USING LOCAL FILE")
	}
	file, err := excelize.OpenFile(filepath.Join(path, name))
	if err != nil {
		log.Fatalf("excel adapter: FetchDocument: %s\n", err)
	}
	return internals.GetDocument(file)
}

func getBoards(document internalModels.Document) map[string]models.Board {
	boards := make(map[string]models.Board, len(document.BoardSheets)-1)
	for name, sheet := range document.BoardSheets {
		boards[name] = models.NewBoard(name, getIP(name, document), sheet)
	}
	return boards
}

func getIP(sheet string, document internalModels.Document) string {
	for _, row := range document.Info.Tables[ADDRESSES_TABLE_NAME].Rows {
		if row[0] == sheet {
			return row[1]
		}
	}
	panic(fmt.Sprintf("excel adapter: getIP: Missing board %s IP\n", sheet))
}

func Update(document internalModels.Document, objects ...models.FromDocument) {
	globalInfo := getGlobalInfo(document)
	for _, object := range objects {
		object.AddGlobal(globalInfo)
	}

	for _, board := range getBoards(document) {
		for _, packet := range board.GetPackets() {
			for _, object := range objects {
				object.AddPacket(board.Name, packet)
			}
		}
	}
}

func getGlobalInfo(document internalModels.Document) models.GlobalInfo {
	return models.GlobalInfo{
		BoardToIP:        getInfoTableToMap(ADDRESSES_TABLE_NAME, document),
		UnitToOperations: getInfoTableToMap(UNITS_TABLE_NAME, document),
		ProtocolToPort:   getInfoTableToMap(PORTS_TABLE_NAME, document),
		BoardToID:        getInfoTableToMap(IDS_TABLE_NAME, document),
	}
}

func getExpandedSection(section models.ControlSection, board models.Board) models.ControlSection {
	expandedSection := make(map[string]interface{})
	for guiName, valueNameInterface := range section {
		valueName := valueNameInterface.(string)
		packetName := board.FindContainingPacket(valueName)
		allIds := internals.GetAllIds(board.Descriptions[packetName].ID)

		if len(allIds) == 1 {
			expandedSection[guiName] = valueName
		} else {
			expandedSection[guiName] = getNamesWithSufix(valueName, len(allIds))
		}

	}
	return expandedSection
}

func getNamesWithSufix(name string, length int) []string {
	namesWithSufix := make([]string, length)
	for i := 0; i < length; i++ {
		namesWithSufix[i] = fmt.Sprintf("%s_%d", name, i)
	}
	return namesWithSufix
}

func getInfoTableToMap(tableName string, document internalModels.Document) map[string]string {
	mapping := make(map[string]string)
	table, found := document.Info.Tables[tableName]
	if !found {
		log.Fatalf("excel adapter: getInfoTableToMap: table %s not found\n", tableName)
	}
	for _, row := range table.Rows {
		mapping[row[0]] = row[1]
	}
	return mapping
}
