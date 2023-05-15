package blcu

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path"
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/common"
	"github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
	"github.com/pin/tftp/v3"
)

type downloadRequest struct {
	Board string `json:"board"`
}

type downloadData struct {
	Board   string
	Payload []byte
}

func (blcu *BLCU) handleDownload(payload json.RawMessage) (downloadData, error) {
	blcu.trace.Debug().Msg("Handling download")
	var request downloadRequest
	if err := json.Unmarshal(payload, &request); err != nil {
		blcu.trace.Error().Err(err).Stack().Msg("Unmarshal payload")
		return downloadData{}, err
	}

	if err := blcu.requestDownload(request.Board); err != nil {
		blcu.trace.Error().Err(err).Stack().Msg("Request download")
		return downloadData{}, err
	}

	buffer := bytes.NewBuffer([]byte{})
	if err := blcu.ReadTFTP(buffer); err != nil {
		blcu.trace.Error().Err(err).Stack().Msg("Read TFTP")
		return downloadData{}, err
	}

	return downloadData{
		Board:   request.Board,
		Payload: buffer.Bytes(),
	}, nil
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

type downloadResponse struct {
	IsSuccess bool   `json:"success"`
	File      []byte `json:"file,omitempty"`
}

func (blcu *BLCU) notifyDownloadFailure() {
	blcu.trace.Warn().Msg("Download failed")
	blcu.sendMessage(blcu.config.Topics.Download, downloadResponse{IsSuccess: false, File: nil})
}

func (blcu *BLCU) notifyDownloadSuccess(data downloadData) {
	blcu.trace.Info().Msg("Download success")
	blcu.sendMessage(blcu.config.Topics.Download, downloadResponse{IsSuccess: true, File: data.Payload})
}

func (blcu *BLCU) writeDownloadFile(data downloadData) error {
	blcu.trace.Info().Msg("Creating download file")

	err := os.MkdirAll(blcu.config.DownloadPath, 0777)
	if err != nil {
		return err
	}
	err = os.Chmod(blcu.config.DownloadPath, 0777)
	if err != nil {
		return err
	}

	return os.WriteFile(path.Join(blcu.config.DownloadPath, data.Board+".bin"), data.Payload, 0777)
}
