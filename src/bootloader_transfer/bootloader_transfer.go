package bootloader_transfer

import (
	"bytes"
	"encoding/json"
	"log"
	"time"

	orderTransferModels "github.com/HyperloopUPV-H8/Backend-H8/order_transfer/models"
	ws_models "github.com/HyperloopUPV-H8/Backend-H8/websocket_handle/models"
	"github.com/pin/tftp/v3"
)

var WRITE_ORDER orderTransferModels.Order = orderTransferModels.Order{
	ID:     700,
	Values: map[string]any{},
}

var READ_ORDER orderTransferModels.Order = orderTransferModels.Order{
	ID:     700,
	Values: map[string]any{},
}

type BootloaderTransfer struct {
	client       *tftp.Client
	orderChannel chan<- orderTransferModels.Order
	channel      chan ws_models.MessageTarget
}

func NewBootloaderTransfer(addr string, orderChannel chan<- orderTransferModels.Order) (*BootloaderTransfer, chan ws_models.MessageTarget, error) {
	client, err := tftp.NewClient(addr)
	if err != nil {
		return nil, nil, err
	}

	bootloaderTransfer := &BootloaderTransfer{
		client:       client,
		orderChannel: orderChannel,
		channel:      make(chan ws_models.MessageTarget),
	}

	client.SetTimeout(time.Second * 5)

	go bootloaderTransfer.run()

	return bootloaderTransfer, bootloaderTransfer.channel, nil
}

func (bootloader *BootloaderTransfer) run() {
	for {
		msg := <-bootloader.channel

		data := new([]byte)
		if err := json.Unmarshal(msg.Msg.Msg, data); err != nil {
			log.Printf("BootloaderTransfer: run: Unmarshal: %s\n", err)
			continue
		}

		bootloader.Put("code.bin", *data)
	}
}

func (bootloader *BootloaderTransfer) Put(name string, data []byte) {
	bootloader.orderChannel <- WRITE_ORDER
	writer, err := bootloader.client.Send(name, "octet")
	if err != nil {
		log.Printf("BootloaderTransfer: Put: Send: %s\n", err)
		bootloader.SendResponse("failure")
		return
	}

	buf := bytes.NewBuffer(data)
	for buf.Len() > 0 {
		_, err := writer.ReadFrom(buf)
		if err != nil {
			log.Printf("BootloaderTransfer: Put: Write: %s\n", err)
			bootloader.SendResponse("failure")
			return
		}
	}

	bootloader.SendResponse("success")
}

func (bootloader *BootloaderTransfer) SendResponse(success string) {
	msgPayload, err := json.Marshal("success")
	if err != nil {
		log.Printf("BootloaderTransfer: SendResponse: Marshal: %s\n", err)
	}
	bootloader.channel <- ws_models.NewMessageTargetRaw([]string{}, "bootloader/response", msgPayload)
}
