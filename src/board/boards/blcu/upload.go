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
	blcu.trace.Debug().Msg("Handling upload")

	var uploadData uploadRequestPayload
	if err := json.Unmarshal(payload, &uploadData); err != nil {
		blcu.trace.Error().Err(err).Stack().Msg("Unmarshal payload")
		return err
	}

	if err := blcu.requestUpload(uploadData.Board); err != nil {
		blcu.trace.Error().Err(err).Stack().Msg("Request upload")
		return err
	}

	if err := blcu.WriteTFTP(bytes.NewReader(uploadData.File)); err != nil {
		blcu.trace.Error().Err(err).Stack().Msg("Write TFTP")
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
	blcu.trace.Info().Str("board", board).Msg("Requesting upload")

	uploadOrder := blcu.createUploadOrder(board)
	if err := blcu.Request(uploadOrder); err != nil {
		return err
	}

	// TODO: remove hardcoded timeout
	if _, err := common.ReadTimeout(blcu.ackChannel, time.Minute); err != nil {
		return err
	}

	return nil
}

func (blcu *BLCU) createUploadOrder(board string) models.Order {
	return models.Order{
		ID: blcu.config.Packets.Upload.Id,
		Fields: map[string]any{
			blcu.config.Packets.Upload.Field: board,
		},
	}
}

func (blcu *BLCU) WriteTFTP(reader io.Reader) error {
	blcu.trace.Info().Msg("Writing TFTP")

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
	blcu.trace.Warn().Msg("Upload failed")
	blcu.sendMessage(os.Getenv("BLCU_UPLOAD_TOPIC"), uploadResponsePayload{Percentage: 0, IsSuccess: false})
}

func (blcu *BLCU) notifyUploadSuccess() {
	blcu.trace.Info().Msg("Upload success")
	blcu.sendMessage(os.Getenv("BLCU_UPLOAD_TOPIC"), uploadResponsePayload{Percentage: 100, IsSuccess: true})
}
