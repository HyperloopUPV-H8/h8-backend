package channels

// Multiple Producer Single Consumer channel (mpsc channel)
type MPSC[T any] struct {
	input  []<-chan T
	output chan<- T
}

func NewMPSC[T any](output chan<- T, blocking bool) *MPSC[T] {
	channel := &MPSC[T]{
		input:  make([]<-chan T, 0),
		output: output,
	}

	if blocking {
		go channel.runBlocking()
	} else {
		go channel.run()
	}

	return channel
}

func (mpsc *MPSC[T]) run() {
	for {
		for _, channel := range mpsc.input {
			select {
			case payload := <-channel:
				select {
				case mpsc.output <- payload:
				default:
				}
			default:
			}
		}
	}
}

func (mpsc *MPSC[T]) runBlocking() {
	for {
		for _, channel := range mpsc.input {
			select {
			case payload := <-channel:
				mpsc.output <- payload
			default:
			}
		}
	}
}

func (mpsc *MPSC[T]) Add(input <-chan T) {
	mpsc.input = append(mpsc.input, input)
}

func (mpsc *MPSC[T]) New(bufSize int) chan<- T {
	channel := make(chan T, bufSize)
	mpsc.input = append(mpsc.input, channel)
	return channel
}
