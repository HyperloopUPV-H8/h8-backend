package observable

type ReplayObservable[T any] struct {
	last      T
	observers map[string]Observer[T]
}

func NewReplayObservable[T any](initialValue T) ReplayObservable[T] {
	return ReplayObservable[T]{
		last:      initialValue,
		observers: make(map[string]Observer[T], 0),
	}
}

func (o *ReplayObservable[T]) Next(value T) {
	o.last = value
	o.broadcast(value)
}

func (o *ReplayObservable[T]) Subscribe(observer Observer[T]) {
	o.observers[observer.Id()] = observer
	observer.Next(o.last)
}

func (o *ReplayObservable[T]) Unsubscribe(id string) {
	delete(o.observers, id)
}

func (o *ReplayObservable[T]) broadcast(value T) {
	for _, observer := range o.observers {
		observer.Next(value)
	}
}
