package logger

import (
	"fmt"
	"log"
	"os"

	"github.com/HyperloopUPV-H8/Backend-H8/data_transfer/domain"
)

func CreateFile() *os.File {
	logFile, errFile := os.Create(os.Getenv("LOG_FILENAME"))

	if errFile != nil {
		log.Fatalf("Error creating file: %v", errFile)
		return nil
	}
	return logFile
}

func WritePacket(packet domain.Packet, logFile *os.File) {
	titlePacket := fmt.Sprintf(`Id: %v    Name: %v    Count: %v    CycleTime: %v`,
		packet.Id, packet.Name, packet.Count, packet.CycleTime)
	fmt.Fprintln(logFile, titlePacket)

	for _, measurement := range packet.Measurements {
		measuramentString := fmt.Sprintf(`	%v: %v`, measurement.Name, measurement.Value.ToDisplayUnitsString())
		fmt.Fprintln(logFile, measuramentString)
	}
}
