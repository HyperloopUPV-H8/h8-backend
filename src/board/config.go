package board

import vehicle_models "github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"

type MuxConfig interface {
	UpdateMux(mux *BoardMux)
}

type muxWithOutput struct {
	outputFunc func(vehicle_models.Order) error
}

func (withOutput muxWithOutput) UpdateMux(mux *BoardMux) {
	mux.withOutput(withOutput.outputFunc)
}

func WithOutput(output func(vehicle_models.Order) error) MuxConfig {
	return muxWithOutput{
		outputFunc: output,
	}
}

type muxWithInput struct {
	inputChan <-chan vehicle_models.Update
}

func (withInput muxWithInput) UpdateMux(mux *BoardMux) {
	go mux.withInput(withInput.inputChan)
}

func WithInput(input <-chan vehicle_models.Update) MuxConfig {
	return muxWithInput{
		inputChan: input,
	}
}
