package channels

type SPMC[T any] struct {
	input  <-chan T
	output []chan<- T
}

func NewSPMC[T any](input <-chan T, blocking bool) *SPMC[T] {
	channel := &SPMC[T]{
		input:  input,
		output: make([]chan<- T, 0),
	}

	if blocking {
		go channel.runBlocking()
	} else {
		go channel.run()
	}

	return channel
}

func (spmc *SPMC[T]) run() {
	for {
		payload := <-spmc.input

		for _, channel := range spmc.output {
			select {
			case channel <- payload:
			default:
			}
		}
	}
}

func (spmc *SPMC[T]) runBlocking() {
	for {
		payload := <-spmc.input

		for _, channel := range spmc.output {
			channel <- payload
		}
	}
}

func (spmc *SPMC[T]) Add(output chan<- T) {
	spmc.output = append(spmc.output, output)
}

func (spmc *SPMC[T]) New(bufSize int) <-chan T {
	channel := make(chan T, bufSize)
	spmc.output = append(spmc.output, channel)
	return channel
}
