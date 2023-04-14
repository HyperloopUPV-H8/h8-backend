package excel_adapter

import (
	"path/filepath"

	"github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/internals"
	internalModels "github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/internals/models"
	"github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/models"
	trace "github.com/rs/zerolog/log"
	"github.com/xuri/excelize/v2"
)

type ExcelAdapterConfig struct {
	Download internals.DownloadConfig
	Parse    internals.ParseConfig
}

func FetchBoardsAndGlobalInfo(config ExcelAdapterConfig) (map[string]models.Board, models.GlobalInfo) {
	document := fetchDocument(config.Download, config.Parse)

	globalInfo := getGlobalInfo(document, GlobalInfoConfig{
		AddressTable:    config.Parse.AddressTable,
		BackendEntryKey: config.Parse.BackendEntryKey,
		UnitsTable:      config.Parse.UnitsTable,
		PortsTable:      config.Parse.PortsTable,
		IdsTable:        config.Parse.IdsTable,
	})

	boards := getBoards(document, config.Parse.AddressTable)

	return boards, globalInfo
}

func fetchDocument(downloadConfig internals.DownloadConfig, parseConfig internals.ParseConfig) internalModels.Document {
	trace.Info().Str("id", downloadConfig.Id).Str("path", downloadConfig.Path).Str("name", downloadConfig.Name).Msg("fetch document")

	errDownloading := internals.DownloadFile(downloadConfig)
	if errDownloading != nil {
		trace.Error().Stack().Err(errDownloading).Msg("")
		trace.Warn().Str("id", downloadConfig.Id).Str("path", downloadConfig.Path).Str("name", downloadConfig.Name).Msg("using local document")
	}

	file, err := excelize.OpenFile(filepath.Join(downloadConfig.Path, downloadConfig.Name))
	if err != nil {
		trace.Fatal().Stack().Err(err).Msg("")
	}

	return internals.GetDocument(file, parseConfig)
}

func getBoards(document internalModels.Document, addressTableName string) map[string]models.Board {
	trace.Trace().Msg("get boards")
	boards := make(map[string]models.Board, len(document.BoardSheets)-1)
	for name, sheet := range document.BoardSheets {
		trace.Trace().Str("board", name).Msg("add board")
		boards[name] = models.NewBoard(name, getIP(name, document, addressTableName), sheet)
	}
	return boards
}

func getIP(board string, document internalModels.Document, addressTableName string) string {
	for _, row := range document.Info.Tables[addressTableName].Rows {
		if row[0] == board {
			trace.Trace().Str("board", board).Str("addr", row[1]).Msg("get board ip")
			return row[1]
		}
	}

	trace.Fatal().Str("board", board).Msg("missing board ip")
	return ""
}

type GlobalInfoConfig struct {
	AddressTable    string
	BackendEntryKey string
	UnitsTable      string
	PortsTable      string
	IdsTable        string
}

func getGlobalInfo(document internalModels.Document, config GlobalInfoConfig) models.GlobalInfo {
	trace.Trace().Msg("get global info")
	return models.GlobalInfo{
		BackendIP:        getBackendIP(config.AddressTable, config.BackendEntryKey, document),
		BoardToIP:        getInfoTableToMap(config.AddressTable, document),
		UnitToOperations: getInfoTableToMap(config.UnitsTable, document),
		ProtocolToPort:   getInfoTableToMap(config.PortsTable, document),
		BoardToID:        getInfoTableToMap(config.IdsTable, document),
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
		//TODO: PUT THIS IN TOML
		if row[0] == "Backend" {
			continue
		}
		mapping[row[0]] = row[1]
	}
	trace.Trace().Str("table", tableName).Msg("get info table")
	return mapping
}

func getBackendIP(addressTableName string, backendKey string, document internalModels.Document) string {
	for _, entry := range document.Info.Tables[addressTableName].Rows {
		if entry[0] == backendKey {
			return entry[1]
		}
	}

	trace.Fatal().Msg("Backend IP not found")
	panic("Backend IP not found") // NEVER RUN BECAUSE trace.Fatal() calls os.exit()
}
