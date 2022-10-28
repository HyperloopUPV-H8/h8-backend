package main

import (
	"fmt"

	"github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/podDataCreator"
	packetadapter "github.com/HyperloopUPV-H8/Backend-H8/Shared/PacketAdapter"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/excelAdapter"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/excelAdapter/dto"
	excelRetriever "github.com/HyperloopUPV-H8/Backend-H8/Shared/excelRetriever/domain"
	"github.com/davecgh/go-spew/spew"
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
						{"Voltage1", "uint8", "cdeg#/100#", "deg##", "[0,10]", "[-10,20]"},
						{"Speed1", "uint16", "cdeg#/100#", "deg##", "[0,10]", "[-10,20]"},
						{"Current1", "uint32", "cdeg#/100#", "deg##", "[0,10]", "[-10,20]"},
						{"Airgap1", "uint64", "cdeg#/100#", "deg##", "[0,10]", "[-10,20]"},
						{"Position1", "uint8", "cdeg#/100#", "deg##", "[0,10]", "[-10,20]"},
						{"Battery1", "int16", "cdeg#/100#", "deg##", "[0,10]", "[-10,20]"},
					},
				},
				"Packet Structure": {
					Name: "Packet Structure",
					Rows: [][]string{
						{"Voltages", "Speeds", "Currents", "Airgaps", "Positions", "Batteries"},
						{"Voltage1", "Speed1", "Current1", "Airgap1", "Position1", "Battery1"},
					},
				},
			},
		},
	},
}

func main() {
	godotenv.Load()

	ips := []string{"127.0.0.1"}

	boardDTOs := excelAdapter.GetBoardDTOs(structure)
	packetDTOs := make([]dto.PacketDTO, 0)
	for _, boardDTO := range boardDTOs {
		packetDTOs = append(packetDTOs, boardDTO.GetPacketDTOs()...)
	}

	packetAdapter := packetadapter.New(packetDTOs, ips)
	podData := podDataCreator.Invoke(boardDTOs)

	fmt.Println("Starting loop")
	for {
		update := packetAdapter.GetPacketUpdate()
		podData.UpdatePacket(update)
		packet := podData.GetPacket(update.ID)
		spew.Dump(packet)
	}
}
