package excel_adapter

import (
	"os"
	"path/filepath"

	"github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/internals"
	internalModels "github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/internals/models"
	"github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/models"
	trace "github.com/rs/zerolog/log"
	"github.com/xuri/excelize/v2"
)

func FetchDocument(id string, path string, name string) internalModels.Document {
	trace.Info().Str("id", id).Str("path", path).Str("name", name).Msg("fetch document")

	errDownloading := internals.DownloadFile(id, path, name)
	if errDownloading != nil {
		trace.Error().Stack().Err(errDownloading).Msg("")
		trace.Warn().Str("id", id).Str("path", path).Str("name", name).Msg("using local document")
	}

	file, err := excelize.OpenFile(filepath.Join(path, name))
	if err != nil {
		trace.Fatal().Stack().Err(err).Msg("")
	}

	return internals.GetDocument(file)
}

func getBoards(document internalModels.Document) map[string]models.Board {
	trace.Trace().Msg("get boards")
	boards := make(map[string]models.Board, len(document.BoardSheets)-1)
	for name, sheet := range document.BoardSheets {
		trace.Trace().Str("board", name).Msg("add board")
		boards[name] = models.NewBoard(name, getIP(name, document), sheet)
	}
	return boards
}

func getIP(board string, document internalModels.Document) string {
	for _, row := range document.Info.Tables[os.Getenv("EXCEL_ADAPTER_ADDRESS_TABLE_NAME")].Rows {
		if row[0] == board {
			trace.Trace().Str("board", board).Str("addr", row[1]).Msg("get board ip")
			return row[1]
		}
	}

	trace.Fatal().Str("board", board).Msg("missing board ip")
	return ""
}

func Update(document internalModels.Document, objects ...models.FromDocument) {
	trace.Debug().Msg("update from document")

	trace.Trace().Msg("update global info")
	globalInfo := getGlobalInfo(document)
	for _, object := range objects {
		object.AddGlobal(globalInfo)
	}

	for name, board := range getBoards(document) {
		trace.Trace().Str("board", name).Msg("update board")
		for _, packet := range board.GetPackets() {
			trace.Trace().Str("packet", packet.Description.ID).Msg("update packet")
			for _, object := range objects {
				object.AddPacket(board.Name, packet)
			}
		}
	}
}

func getGlobalInfo(document internalModels.Document) models.GlobalInfo {
	trace.Trace().Msg("get global info")
	return models.GlobalInfo{
		BoardToIP:        getInfoTableToMap(os.Getenv("EXCEL_ADAPTER_ADDRESS_TABLE_NAME"), document),
		UnitToOperations: getInfoTableToMap(os.Getenv("EXCEL_ADAPTER_UNITS_TABLE_NAME"), document),
		ProtocolToPort:   getInfoTableToMap(os.Getenv("EXCEL_ADAPTER_PORTS_TABLE_NAME"), document),
		BoardToID:        getInfoTableToMap(os.Getenv("EXCEL_ADAPTER_IDS_TABLE_NAME"), document),
	}
}

func getInfoTableToMap(tableName string, document internalModels.Document) map[string]string {
	mapping := make(map[string]string)
	table, found := document.Info.Tables[tableName]
	if !found {
		trace.Fatal().Str("table", tableName).Msg("table not found")
		return nil
	}

	for _, row := range table.Rows {
		mapping[row[0]] = row[1]
	}
	trace.Trace().Str("table", tableName).Msg("get info table")
	return mapping
}
