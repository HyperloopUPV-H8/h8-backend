package observable

import "encoding/json"

type SubscriptionMessage struct {
	Id        string `json:"id"`
	Subscribe bool   `json:"subscribe"`
}

func HandleSubscribe[T any](obs Observable[T], source string, payload []byte, sendMessage func(v T, id string) error) {
	var subcription SubscriptionMessage
	err := json.Unmarshal(payload, &subcription)

	if err != nil {
		return
	}

	if subcription.Subscribe {
		addWsObserver(obs, subcription.Id, source, sendMessage)
	} else {
		obs.Unsubscribe(subcription.Id)
	}
}

func addWsObserver[T any](obs Observable[T], id string, connId string, sendMessage func(v T, id string) error) {
	observer := NewWsObserver(id, connId, func(v T) {
		err := sendMessage(v, connId)

		if err != nil {
			obs.Unsubscribe(id)
		}
	})

	obs.Subscribe(observer)
}

type WsObserver[T any] struct {
	id          string
	connId      string
	sendMessage func(T)
}

func NewWsObserver[T any](id string, connId string, sendMessage func(T)) WsObserver[T] {
	return WsObserver[T]{
		id:          id,
		connId:      connId,
		sendMessage: sendMessage,
	}
}

func (o WsObserver[T]) Id() string {
	return o.id
}

func (o WsObserver[T]) Next(v T) {
	o.sendMessage(v)
}
