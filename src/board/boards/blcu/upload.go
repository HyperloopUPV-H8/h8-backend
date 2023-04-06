package blcu

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/common"
	"github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
	"github.com/pin/tftp/v3"
)

func (blcu *BLCU) handleUpload(payload json.RawMessage) error {
	var uploadData uploadRequestPayload
	if err := json.Unmarshal(payload, &uploadData); err != nil {
		return err
	}

	if err := blcu.requestUpload(uploadData.Board); err != nil {
		return err
	}

	if err := blcu.WriteTFTP(bytes.NewReader(uploadData.File)); err != nil {
		return err
	}

	return nil
}

type uploadRequestPayload struct {
	Board string `json:"board"`
	File  []byte `json:"file"`
}

type uploadResponsePayload struct {
	Percentage int  `json:"percentage"`
	IsSuccess  bool `json:"success"`
}

func (blcu *BLCU) requestUpload(board string) error {
	uploadOrder := createUploadOrder(board)
	if err := blcu.Request(uploadOrder); err != nil {
		return err
	}

	// TODO: remove hardcoded timeout
	if _, err := common.ReadTimeout(blcu.ackChannel, time.Minute); err != nil {
		return err
	}

	return nil
}

func createUploadOrder(board string) models.Order {
	return models.Order{
		ID: BLCU_UPLOAD_ORDER_ID,
		Fields: map[string]any{
			BLCU_UPLOAD_ORDER_FIELD: board,
		},
	}
}

func (blcu *BLCU) WriteTFTP(reader io.Reader) error {
	client, err := tftp.NewClient(blcu.addr)
	if err != nil {
		return err
	}

	sender, err := client.Send("a.bin", "octet")
	if err != nil {
		return err
	}

	_, err = sender.ReadFrom(reader)
	if err != nil {
		return err
	}

	return nil
}

// the topic BLCU_STATE_TOPIC expects a number between 0 and 100, the idea is in the future to inform about the percentage of the file uploaded
func (blcu *BLCU) notifyUploadFailure() {
	// 0 means failre
	blcu.sendMessage(os.Getenv("BLCU_UPLOAD_TOPIC"), uploadResponsePayload{Percentage: 0, IsSuccess: false})
}

func (blcu *BLCU) notifyUploadSuccess() {
	// 100 means success
	blcu.sendMessage(os.Getenv("BLCU_UPLOAD_TOPIC"), uploadResponsePayload{Percentage: 100, IsSuccess: true})
}
