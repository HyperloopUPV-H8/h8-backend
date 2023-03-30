package board

import (
	"github.com/HyperloopUPV-H8/Backend-H8/board/models"
	vehicle_models "github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
	"github.com/rs/zerolog"
	trace "github.com/rs/zerolog/log"
)

const UNRECOGNIZED_CHAN_BUF_SIZE = 100

type BoardMux struct {
	boards    map[string]models.Board
	idToBoard map[uint16]string

	sendOrder    func(vehicle_models.Order) error
	unrecognized chan vehicle_models.Update

	trace zerolog.Logger
}

func NewMux(options ...MuxOptions) *BoardMux {
	trace.Info().Msg("create new board mux")

	mux := &BoardMux{
		unrecognized: make(chan vehicle_models.Update, UNRECOGNIZED_CHAN_BUF_SIZE),
		boards:       make(map[string]models.Board),
		idToBoard:    make(map[uint16]string),
		trace:        trace.With().Str("component", "boardMux").Logger(),
	}

	for _, option := range options {
		option.UpdateMux(mux)
	}

	return mux
}

func (mux *BoardMux) AddBoard(board models.Board) {
	board.Output(mux.sendOrder)
	mux.boards[board.Name()] = board
	mux.trace.Debug().Str("board", board.Name()).Msg("add board to board mux")
}

func (mux *BoardMux) AddBoardMapping(idToBoard map[uint16]string) {
	for id, board := range idToBoard {
		mux.idToBoard[id] = board
	}
	mux.trace.Debug().Interface("mappings", idToBoard).Msg("add board mappings to board mux")
}

// Boards must not be added after the first call to Listen
func (mux *BoardMux) Listen(output chan<- vehicle_models.Update) {
	mux.trace.Debug().Msg("started listening")
	for _, board := range mux.boards {
		go board.Listen(output)
		mux.trace.Debug().Str("board", board.Name()).Msg("board started listening")
	}

	for update := range mux.unrecognized {
		output <- update
		mux.trace.Trace().Msg("skip unrecognized update")
	}

	mux.trace.Warn().Msg("stopped listening")
}

func (mux *BoardMux) Request(order vehicle_models.Order) error {
	mux.trace.Debug().Uint16("id", order.ID).Msg("request order")
	boardName, found := mux.idToBoard[order.ID]
	if !found {
		mux.trace.Trace().Uint16("id", order.ID).Msg("default send order with unknown id")
		return mux.sendOrder(order)
	}

	board, found := mux.boards[boardName]
	if !found {
		mux.trace.Trace().Uint16("id", order.ID).Str("board", boardName).Msg("default send order without custom board logic")
		return mux.sendOrder(order)
	}

	mux.trace.Trace().Uint16("id", order.ID).Str("board", boardName).Msg("request board")
	return board.Request(order)
}
