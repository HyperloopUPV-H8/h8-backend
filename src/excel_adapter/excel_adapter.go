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

const GlobalSheet = "Info"
const AddressesTable = "addresses"

func FetchDocument(id string, path string, name string) internalModels.Document {
	internals.DownloadFile(id, path, name)
	file, err := excelize.OpenFile(filepath.Join(path, name))
	if err != nil {
		log.Fatalf("excel adapter: FetchDocument: %s\n", err)
	}
	return internals.GetDocument(file)
}

func getBoards(document internalModels.Document) map[string]models.Board {
	boards := make(map[string]models.Board, len(document.Sheets)-1)
	for name, sheet := range document.Sheets {
		if name != GlobalSheet {
			boards[name] = models.NewBoard(name, getIP(name, document), sheet)
		}
	}
	return boards
}

func getIP(sheet string, document internalModels.Document) string {
	for _, row := range document.Sheets[GlobalSheet].Tables[AddressesTable].Rows {
		if row[0] == sheet {
			return row[1]
		}
	}
	panic(fmt.Sprintf("excel adapter: getIP: Missing board %s IP\n", sheet))
}

func Compile(document internalModels.Document, objects ...models.FromBoards) {
	for _, board := range getBoards(document) {
		for _, packet := range board.GetPackets() {
			for _, object := range objects {
				object.AddPacket(board.Name, board.IP, packet.Description, packet.Values)
			}
		}
	}
}
