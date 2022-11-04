package main

import (

	// NO ELIMINAR //"github.com/HyperloopUPV-H8/Backend-H8/dataTransfer"

	dataTransfer "github.com/HyperloopUPV-H8/Backend-H8/DataTransfer"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/PacketAdapter/transportController"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/excelAdapter"
	excelAdapterDomain "github.com/HyperloopUPV-H8/Backend-H8/Shared/excelAdapter/domain"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/excelRetriever"
	excelRetrieverDomain "github.com/HyperloopUPV-H8/Backend-H8/Shared/excelRetriever/domain"
	"github.com/HyperloopUPV-H8/Backend-H8/cmd/MVP-1/logger"

	"github.com/joho/godotenv"
)

var structure = excelRetrieverDomain.Document{
	Sheets: map[string]excelRetrieverDomain.Sheet{
		"BMS": {
			Name: "BMS",
			Tables: map[string]excelRetrieverDomain.Table{
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
	document := excelRetriever.GetExcel("excel.xlsx", ".")

	boards := excelAdapter.GetBoards(document)
	packets := make([]excelAdapterDomain.PacketDTO, 0)
	for _, board := range boards {
		packets = append(packets, board.GetPackets()...)
	}

	packetAdapter := transportController.New(ips, packets)

	logFile := logger.CreateFile()
	defer logFile.Close()

	dataTransfer := dataTransfer.New(boards)
	dataTransfer.Invoke(packetAdapter.ReceiveData)

	for packet := range dataTransfer.PacketChannel {
		logger.WritePacket(packet, logFile)
	}
}
