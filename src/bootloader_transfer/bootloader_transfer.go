package bootloader_transfer

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"pack.ag/tftp"
)

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	programBuffer, err := io.ReadAll(r.Body)

	if err != nil {
		log.Fatal("error reading bootloader request Body")
	}

	//TODO: once firmware is ready, we have to configure these options
	// tftpClient, err := tftp.NewClient(
	// 	tftp.ClientBlocksize(),
	// 	tftp.ClientMode(),
	// 	tftp.ClientRetransmit(),
	// 	tftp.ClientTimeout(),
	// 	tftp.ClientTransferSize(),
	// 	tftp.ClientWindowsize())

	tftpClient, err := tftp.NewClient()

	if err != nil {
		log.Fatal("error creating tftp client")
	}

	tftpClient.Put(fmt.Sprintf("tftp://%s:%s/program", os.Getenv("BOOTLOADER_BOARD_IP"), os.Getenv("BOOTLOADER_BOARD_IP")), bytes.NewBuffer(programBuffer), 0)

}
