package main

import (
	"fmt"

	// NO ELIMINAR //"github.com/HyperloopUPV-H8/Backend-H8/dataTransfer"

	podData "github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/application"
	excelParser "github.com/HyperloopUPV-H8/Backend-H8/Shared/ExcelParser/application"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/ExcelParser/domain/board"
	excelRetriever "github.com/HyperloopUPV-H8/Backend-H8/Shared/ExcelParser/domain/document"
	packetAdapter "github.com/HyperloopUPV-H8/Backend-H8/Shared/PacketAdapter/application"

	"github.com/joho/godotenv"
)

var structure = excelRetriever.Document{
	Sheets: map[string]excelRetriever.Sheet{
		"BMS": {
			Name: "BMS",
			Tables: map[string]excelRetriever.Table{
				"Packet Description": {
					Name: "Packet Description",
					Rows: [][]string{
						{"1[0,5]", "Voltages", "200", "Input", "TCP"},
						{"2", "Speeds", "300", "Input", "UDP"},
						{"3", "Currents", "400", "Input", "TCP"},
						{"4", "Airgaps", "500", "Output", "UDP"},
						{"5", "Positions", "600", "Input", "UDP"},
						{"6", "Batteries", "700", "Input", "TCP"},
					},
				},
				"Value Description": {
					Name: "Value Description",
					Rows: [][]string{
						{"Voltage0", "uint8", "cdeg#/100#", "deg##", "[0,10]", "[-10,20]"},
						{"Speed0", "bool", "cdeg#/100#", "deg##", "[0,10]", "[-10,20]"},
						{"Current0", "uint32", "cdeg#/100#", "deg##", "[0,10]", "[-10,20]"},
						{"Airgap0", "uint64", "cdeg#/100#", "deg##", "[0,10]", "[-10,20]"},
						{"Position0", "uint8", "cdeg#/100#", "deg##", "[0,10]", "[-10,20]"},
						{"Battery0", "int16", "cdeg#/100#", "deg##", "[0,10]", "[-10,20]"},
					},
				},
				"Packet Structure": {
					Name: "Packet Structure",
					Rows: [][]string{
						{"Voltages", "Speeds", "Currents", "Airgaps", "Positions", "Batteries"},
						{"Voltage0", "Speed0", "Current0", "Airgap0", "Position0", "Battery0"},
					},
				},
			},
		},
	},
}

func main() {
	godotenv.Load("./.env")

	ips := []string{"127.0.0.1"}
	document := excelParser.GetExcel("excel.xlsx", ".")

	boards := excelParser.GetBoards(document)
	packets := make([]board.Packet, 0)
	for _, board := range boards {
		packets = append(packets, board.GetPackets()...)
	}

	packetAdapter := packetAdapter.New(ips, packets)

	fmt.Println("Starting loop")

	data := podData.New(boards)
	data.Invoke(packetAdapter.ReadData)
}
