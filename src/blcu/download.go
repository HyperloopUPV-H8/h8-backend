package blcu

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func (blcu *BLCU) download(payload json.RawMessage) (string, DownloadResult) {
	blcu.trace.Debug().Msg("Handling download")
	var request downloadRequest
	if err := json.Unmarshal(payload, &request); err != nil {
		blcu.trace.Error().Err(err).Stack().Msg("Unmarshal payload")
		return "", DownloadResult{Err: err}
	}

	if err := blcu.requestDownload(request.Board); err != nil {
		blcu.trace.Error().Err(err).Stack().Msg("Request download")
		return request.Board, DownloadResult{Err: err}
	}

	download, err := blcu.ReadTFTP()
	if err != nil {
		blcu.trace.Error().Err(err).Stack().Msg("Read payload")
		return request.Board, DownloadResult{Err: err}
	}

	go blcu.consumeDownloadPercentage(download.Percentage)

	return request.Board, <-download.Result
}

func (blcu *BLCU) consumeDownloadPercentage(percentage <-chan float64) {
	for p := range percentage {
		blcu.notifyDownloadProgress(p)
	}
}

func (blcu *BLCU) requestDownload(board string) error {
	blcu.trace.Info().Str("board", board).Msg("Requesting download")

	downloadOrder := blcu.createDownloadOrder(board)
	if err := blcu.sendOrder(downloadOrder); err != nil {
		return err
	}

	// TODO: remove hardcoded timeout
	if _, err := common.ReadTimeout(blcu.ackChannel, time.Second*10); err != nil {
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

func (blcu *BLCU) ReadTFTP() (Download, error) {
	blcu.trace.Info().Msg("Reading TFTP")

	client, err := tftp.NewClient(blcu.addr)
	if err != nil {
		return Download{}, err
	}

	receiver, err := client.Receive("a.bin", "octet")
	if err != nil {
		return Download{}, err
	}

	download := NewDownload(786432) // Size of the flash memory (downlaod size)
	go receiver.WriteTo(&download)

	return download, nil
}

type downloadResponse struct {
	Percentage float64 `json:"percentage"`
	IsSuccess  *bool   `json:"success,omitempty"`
	File       []byte  `json:"file,omitempty"`
}

func (blcu *BLCU) notifyDownloadFailure() {
	blcu.trace.Warn().Msg("Download failed")
	success := new(bool)
	*success = false
	blcu.sendMessage(blcu.config.Topics.Download, downloadResponse{IsSuccess: success, File: nil, Percentage: 0.0})
}

func (blcu *BLCU) notifyDownloadSuccess(data []byte) {
	blcu.trace.Info().Msg("Download success")
	success := new(bool)
	*success = true
	blcu.sendMessage(blcu.config.Topics.Download, downloadResponse{IsSuccess: success, File: data, Percentage: 1.0})
}

func (blcu *BLCU) notifyDownloadProgress(percentage float64) {
	blcu.sendMessage(blcu.config.Topics.Download, downloadResponse{IsSuccess: nil, File: nil, Percentage: percentage})
}

func (blcu *BLCU) writeDownloadFile(board string, data []byte) error {
	blcu.trace.Info().Msg("Creating download file")

	err := os.MkdirAll(blcu.config.DownloadPath, 0777)
	if err != nil {
		return err
	}
	err = os.Chmod(blcu.config.DownloadPath, 0777)
	if err != nil {
		return err
	}

	return os.WriteFile(path.Join(blcu.config.DownloadPath, fmt.Sprintf("%s-%d.bin", board, time.Now().Unix())), data, 0777)
}

type DownloadResult struct {
	Data []byte
	Err  error
}

type Download struct {
	output         bytes.Buffer
	Result         <-chan DownloadResult
	resultChan     chan<- DownloadResult
	Percentage     <-chan float64
	percentageChan chan<- float64
	notify         chan int
	total          int
	writeErr       error
}

func NewDownload(size int) Download {
	resultChan := make(chan DownloadResult)
	percentageChan := make(chan float64, 100)

	download := Download{
		output:         bytes.Buffer{},
		Result:         resultChan,
		resultChan:     resultChan,
		Percentage:     percentageChan,
		percentageChan: percentageChan,
		notify:         make(chan int),
		total:          size,
		writeErr:       nil,
	}

	go download.consumeNotifications()

	return download
}

func (download *Download) Write(p []byte) (n int, err error) {
	if download.writeErr != nil {
		return 0, download.writeErr
	}
	n, err = download.output.Write(p)
	if err != nil {
		download.writeErr = err
		download.abort(err)
	} else {
		download.notify <- n
	}
	return n, err
}

func (download *Download) abort(err error) {
	close(download.notify)
	download.resultChan <- DownloadResult{
		Data: download.output.Bytes(),
		Err:  err,
	}
}

func (download *Download) consumeNotifications() {
	defer close(download.percentageChan)
	current := 0
	for amount := range download.notify {
		current += amount
		percentage := float64(current) / float64(download.total)
		download.percentageChan <- common.Clamp(percentage, 0.0, 1.0)
		if current >= amount {
			download.resultChan <- DownloadResult{
				Data: download.output.Bytes(),
				Err:  nil,
			}
			return
		}
	}
}
