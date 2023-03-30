package board

import vehicle_models "github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"

type muxWithInput struct {
	inputChan <-chan vehicle_models.Update
}

func (withInput muxWithInput) UpdateMux(mux *BoardMux) {
	mux.trace.Trace().Msg("set vehicle input")
	go mux.withInput(withInput.inputChan)
}

func WithInput(input <-chan vehicle_models.Update) MuxOptions {
	return muxWithInput{
		inputChan: input,
	}
}

func (mux *BoardMux) withInput(input <-chan vehicle_models.Update) {
	for update := range input {
		boardName, found := mux.idToBoard[update.ID]
		if !found {
			mux.unrecognized <- update
			mux.trace.Trace().Uint16("id", update.ID).Msg("skip input mapping with unknown id")
			continue
		}

		board, found := mux.boards[boardName]
		if !found {
			mux.unrecognized <- update
			mux.trace.Trace().Uint16("id", update.ID).Str("board", boardName).Msg("skip input mapping without custom board logic")
			continue
		}

		mux.trace.Trace().Uint16("id", update.ID).Str("board", boardName).Msg("input mapping")
		board.Input(update)
	}
}
