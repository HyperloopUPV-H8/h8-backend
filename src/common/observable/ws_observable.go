package observable

func HandleSubscribe[T any](obs Observable[T], id string, subscribe bool, sendMessage func(v T, id string) error) {
	if subscribe {
		addWsObserver(obs, id, sendMessage)
	} else {
		obs.Unsubscribe(id)
	}
}

func addWsObserver[T any](obs Observable[T], id string, sendMessage func(v T, id string) error) {
	observer := NewWsObserver(id, func(v T) {
		err := sendMessage(v, id)

		if err != nil {
			obs.Unsubscribe(id)
		}
	})

	obs.Subscribe(observer)
}

type WsObserver[T any] struct {
	id          string
	sendMessage func(T)
}

func NewWsObserver[T any](id string, sendMessage func(T)) WsObserver[T] {
	return WsObserver[T]{
		id:          id,
		sendMessage: sendMessage,
	}
}

func (o WsObserver[T]) Id() string {
	return o.id
}

func (o WsObserver[T]) Next(v T) {
	o.sendMessage(v)
}
