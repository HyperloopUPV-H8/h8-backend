package blcu

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/common"
	"github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
	"github.com/pin/tftp/v3"
)

func (blcu *BLCU) handleDownload(payload json.RawMessage) ([]byte, error) {
	blcu.trace.Debug().Msg("Handling download")
	var board string
	if err := json.Unmarshal(payload, &board); err != nil {
		blcu.trace.Error().Err(err).Stack().Msg("Unmarshal payload")
		return nil, err
	}

	if err := blcu.requestDownload(board); err != nil {
		blcu.trace.Error().Err(err).Stack().Msg("Request download")
		return nil, err
	}

	var buffer *bytes.Buffer
	if err := blcu.ReadTFTP(buffer); err != nil {
		blcu.trace.Error().Err(err).Stack().Msg("Read TFTP")
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (blcu *BLCU) requestDownload(board string) error {
	blcu.trace.Info().Str("board", board).Msg("Requesting download")

	downloadOrder := blcu.createDownloadOrder(board)
	if err := blcu.sendOrder(downloadOrder); err != nil {
		return err
	}

	// TODO: remove hardcoded timeout
	if _, err := common.ReadTimeout(blcu.ackChannel, time.Minute); err != nil {
		return err
	}

	return nil
}

func (blcu *BLCU) createDownloadOrder(board string) models.Order {
	return models.Order{
		ID: blcu.config.Packets.Download.Id,
		Fields: map[string]models.Field{
			blcu.config.Packets.Download.Field: {
				Value:     board,
				IsEnabled: true,
			},
		},
	}
}

func (blcu *BLCU) ReadTFTP(writer io.Writer) error {
	blcu.trace.Info().Msg("Reading TFTP")

	client, err := tftp.NewClient(blcu.addr)
	if err != nil {
		return err
	}

	receiver, err := client.Receive("a.bin", "octet")
	if err != nil {
		return err
	}

	_, err = receiver.WriteTo(writer)
	if err != nil {
		return err
	}

	return nil
}

type fileResponsePayload struct {
	File      []byte `json:"file"`
	IsSuccess bool   `json:"success"`
}

func (blcu *BLCU) notifyDownloadFailure() {
	blcu.trace.Warn().Msg("Download failed")
	blcu.sendMessage(blcu.config.Topics.Download, fileResponsePayload{IsSuccess: false, File: nil})
}

func (blcu *BLCU) notifyDownloadSuccess(bytes []byte) {
	blcu.trace.Info().Msg("Download success")
	blcu.sendMessage(blcu.config.Topics.Download, fileResponsePayload{IsSuccess: true, File: bytes})
}
