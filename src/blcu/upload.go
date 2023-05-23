package blcu

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/common"
	"github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
	"github.com/pin/tftp/v3"
)

type uploadRequest struct {
	Board string `json:"board"`
	File  string `json:"file"`
}

func (blcu *BLCU) upload(payload json.RawMessage) UploadResult {
	blcu.trace.Debug().Msg("Handling upload")

	var request uploadRequest
	if err := json.Unmarshal(payload, &request); err != nil {
		blcu.trace.Error().Err(err).Stack().Msg("Unmarshal payload")
		return UploadResult{Err: err}
	}

	if err := blcu.requestUpload(request.Board); err != nil {
		blcu.trace.Error().Err(err).Stack().Msg("Request upload")
		return UploadResult{Err: err}
	}

	decoded, err := base64.StdEncoding.DecodeString(request.File)
	if err != nil {
		blcu.trace.Error().Err(err).Stack().Msg("Decode payload")
		return UploadResult{Err: err}
	}

	upload, err := blcu.WriteTFTP(bytes.NewReader(decoded))
	if err != nil {
		blcu.trace.Error().Err(err).Stack().Msg("Write payload")
		return UploadResult{Err: err}
	}

	go blcu.consumeUploadPercentage(upload.Percentage)

	return <-upload.Result
}

func (blcu *BLCU) consumeUploadPercentage(percentage <-chan float64) {
	for p := range percentage {
		blcu.notifyUploadProgress(p)
	}
}

type uploadResponse struct {
	Percentage float64 `json:"percentage"`
	IsSuccess  *bool   `json:"success,omitempty"`
}

func (blcu *BLCU) requestUpload(board string) error {
	blcu.trace.Info().Str("board", board).Msg("Requesting upload")

	uploadOrder := blcu.createUploadOrder(board)
	if err := blcu.sendOrder(uploadOrder); err != nil {
		return err
	}

	// TODO: remove hardcoded timeout
	if _, err := common.ReadTimeout(blcu.ackChannel, time.Second*10); err != nil {
		return err
	}

	return nil
}

func (blcu *BLCU) createUploadOrder(board string) models.Order {
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

func (blcu *BLCU) WriteTFTP(reader *bytes.Reader) (Upload, error) {
	blcu.trace.Info().Msg("Writing TFTP")

	client, err := tftp.NewClient(blcu.addr)
	if err != nil {
		return Upload{}, err
	}

	sender, err := client.Send("a.bin", "octet")
	if err != nil {
		return Upload{}, err
	}

	upload := NewUpload(reader, reader.Len())

	go sender.ReadFrom(&upload)

	return upload, nil
}

func (blcu *BLCU) notifyUploadFailure() {
	blcu.trace.Warn().Msg("Upload failed")
	success := new(bool)
	*success = false
	blcu.sendMessage(blcu.config.Topics.Download, uploadResponse{Percentage: 0.0, IsSuccess: success})
}

func (blcu *BLCU) notifyUploadSuccess() {
	blcu.trace.Info().Msg("Upload success")
	success := new(bool)
	*success = true
	blcu.sendMessage(blcu.config.Topics.Download, uploadResponse{Percentage: 1.0, IsSuccess: success})
}

func (blcu *BLCU) notifyUploadProgress(percentage float64) {
	blcu.sendMessage(blcu.config.Topics.Download, uploadResponse{Percentage: percentage, IsSuccess: nil})
}

type UploadResult struct {
	Err error
}

type Upload struct {
	input          *bytes.Reader
	Result         <-chan UploadResult
	resultChan     chan<- UploadResult
	Percentage     <-chan float64
	percentageChan chan<- float64
	notify         chan int
	total          int
	readErr        error
}

func NewUpload(data *bytes.Reader, size int) Upload {
	resultChan := make(chan UploadResult)
	percentageChan := make(chan float64, 100)

	upload := Upload{
		input:          data,
		Result:         resultChan,
		resultChan:     resultChan,
		Percentage:     percentageChan,
		percentageChan: percentageChan,
		notify:         make(chan int),
		total:          size,
		readErr:        nil,
	}

	go upload.consumeNotifications()

	return upload
}

func (upload *Upload) Read(p []byte) (n int, err error) {
	if upload.readErr != nil {
		return 0, upload.readErr
	}
	n, err = upload.input.Read(p)
	if err != nil {
		upload.readErr = err
		upload.abort(err)
	} else {
		upload.notify <- n
	}
	return n, err
}

func (upload *Upload) abort(err error) {
	close(upload.notify)
	upload.resultChan <- UploadResult{
		Err: err,
	}
}

func (upload *Upload) consumeNotifications() {
	defer close(upload.percentageChan)
	current := 0
	for amount := range upload.notify {
		current += amount
		percentage := float64(current) / float64(upload.total)
		upload.percentageChan <- common.Clamp(percentage, 0.0, 1.0)
		if current >= amount {
			upload.resultChan <- UploadResult{
				Err: nil,
			}
			return
		}
	}
}
