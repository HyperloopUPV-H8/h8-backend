package board

import vehicle_models "github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"

type muxWithOutput struct {
	outputFunc func(vehicle_models.Order) error
}

func (withOutput muxWithOutput) UpdateMux(mux *BoardMux) {
	mux.trace.Trace().Msg("set vehicle output")
	mux.withOutput(withOutput.outputFunc)
}

func WithOutput(output func(vehicle_models.Order) error) MuxOptions {
	return muxWithOutput{
		outputFunc: output,
	}
}

func (mux *BoardMux) withOutput(output func(vehicle_models.Order) error) {
	for _, board := range mux.boards {
		board.Output(output)
	}

	mux.sendOrder = output
}
