package board

import (
	"log"

	"github.com/HyperloopUPV-H8/Backend-H8/board/models"
	vehicle_models "github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
)

const UNRECOGNIZED_CHAN_BUF_SIZE = 100

type BoardMux struct {
	boards    map[string]models.Board
	idToBoard map[uint16]string

	sendOrder    func(vehicle_models.Order) error
	unrecognized chan vehicle_models.Update
}

func NewMux(config ...MuxConfig) *BoardMux {
	mux := &BoardMux{
		unrecognized: make(chan vehicle_models.Update, UNRECOGNIZED_CHAN_BUF_SIZE),
	}

	for _, conf := range config {
		conf.UpdateMux(mux)
	}

	return mux
}

func (mux *BoardMux) AddBoard(board models.Board) {
	board.Output(mux.sendOrder)
	mux.boards[board.Name()] = board
}

func (mux *BoardMux) AddBoardMapping(idToBoard map[uint16]string) {
	for id, board := range idToBoard {
		mux.idToBoard[id] = board
	}
}

func (mux *BoardMux) withInput(input <-chan vehicle_models.Update) {
	for update := range input {
		boardName, found := mux.idToBoard[update.ID]
		if !found {
			mux.unrecognized <- update
			continue
		}

		board, found := mux.boards[boardName]
		if !found {
			mux.unrecognized <- update
			continue
		}

		board.Input(update)
	}
}

func (mux *BoardMux) withOutput(output func(vehicle_models.Order) error) {
	for _, board := range mux.boards {
		board.Output(output)
	}

	mux.sendOrder = output
}

// Boards must not be added after the first call to Listen
func (mux *BoardMux) Listen(output chan<- vehicle_models.Update) {
	for _, board := range mux.boards {
		go board.Listen(output)
	}

	for update := range mux.unrecognized {
		output <- update
	}
}

func (mux *BoardMux) Request(order vehicle_models.Order) error {
	log.Println("request", order)
	boardName, found := mux.idToBoard[order.ID]
	if !found {
		return mux.sendOrder(order)
	}

	board, found := mux.boards[boardName]
	if !found {
		return mux.sendOrder(order)
	}

	return board.Request(order)
}
