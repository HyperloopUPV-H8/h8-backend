package blcu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/common"
	"github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
	"github.com/pin/tftp/v3"
)

type uploadRequest struct {
	Board string `json:"board"`
	File  string `json:"file"`
}

func (blcu *BLCU) handleUpload(payload json.RawMessage) error {
	blcu.trace.Debug().Msg("Handling upload")

	var request uploadRequest
	if err := json.Unmarshal(payload, &request); err != nil {
		blcu.trace.Error().Err(err).Stack().Msg("Unmarshal payload")
		return err
	}

	id, ok := blcu.boardToId[request.Board]

	if !ok {
		err := fmt.Errorf("id not found for board %s", request.Board)
		blcu.trace.Error().Err(err).Stack().Msg("Board not found")
		return err
	}

	if err := blcu.requestUpload(id); err != nil {
		blcu.trace.Error().Err(err).Stack().Msg("Request upload")
		return err
	}

	if err := blcu.WriteTFTP(bytes.NewReader([]byte(request.File))); err != nil {
		blcu.trace.Error().Err(err).Stack().Msg("Write TFTP")
		return err
	}

	return nil
}

type uploadResponse struct {
	Percentage int  `json:"percentage"`
	IsSuccess  bool `json:"success"`
}

func (blcu *BLCU) requestUpload(board uint16) error {
	blcu.trace.Info().Uint16("board", board).Msg("Requesting upload")

	uploadOrder := blcu.createUploadOrder(float64(board))
	if err := blcu.sendOrder(uploadOrder); err != nil {
		return err
	}

	// TODO: remove hardcoded timeout
	if _, err := common.ReadTimeout(blcu.ackChannel, time.Minute); err != nil {
		return err
	}

	return nil
}

func (blcu *BLCU) createUploadOrder(board float64) models.Order {
	return models.Order{
		ID: blcu.config.Packets.Upload.Id,
		Fields: map[string]models.Field{
			blcu.config.Packets.Upload.Field: {
				Value:     board,
				IsEnabled: true,
			},
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
	blcu.sendMessage(blcu.config.Topics.Download, uploadResponse{Percentage: 0, IsSuccess: false})
}

func (blcu *BLCU) notifyUploadSuccess() {
	blcu.trace.Info().Msg("Upload success")
	blcu.sendMessage(blcu.config.Topics.Download, uploadResponse{Percentage: 100, IsSuccess: true})
}
